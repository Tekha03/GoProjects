package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	pb "flight-booking/booking-service/proto/flight/v1"

	"github.com/jackc/pgx/v5/pgxpool"
)

type CircuitBreaker struct {
	state          string
	failures       int
	lastFail       time.Time
	threshold      int
	timeout        time.Duration
	resetTimeout   time.Duration
}

var cb = &CircuitBreaker{
	state: "CLOSED", threshold: 5, timeout: 30 * time.Second, resetTimeout: 10 * time.Second,
}

func (c *CircuitBreaker) Execute(fn func() error) error {
	if c.state == "OPEN" {
		if time.Since(c.lastFail) > c.resetTimeout {
			c.state = "HALF_OPEN"
			log.Println("CB → HALF_OPEN")
		} else {
			return status.Error(codes.Unavailable, "503 Service Unavailable (circuit breaker open)")
		}
	}

	err := fn()
	if err != nil {
		c.failures++
		c.lastFail = time.Now()
		if c.failures >= c.threshold {
			c.state = "OPEN"
			log.Println("CB → OPEN")
		}
		return err
	}

	c.failures = 0
	if c.state == "HALF_OPEN" {
		c.state = "CLOSED"
		log.Println("CB → CLOSED")
	}
	return nil
}

func retryInterceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	max := 3
	for attempt := 0; attempt < max; attempt++ {
		err := invoker(ctx, method, req, reply, cc, opts...)
		if err == nil { return nil }
		st, ok := status.FromError(err)
		if !ok { return err }
		if st.Code() != codes.Unavailable && st.Code() != codes.DeadlineExceeded {
			return err
		}
		backoff := time.Duration(100*(1<<attempt)) * time.Millisecond
		time.Sleep(backoff)
		log.Printf("Retry %d for %s", attempt+1, method)
	}
	return status.Error(codes.Unavailable, "max retries exceeded")
}

func authInterceptor(apiKey string) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		md := metadata.Pairs("x-api-key", apiKey)
		ctx = metadata.NewOutgoingContext(ctx, md)
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

func initDB() *pgxpool.Pool {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))
	db, _ := pgxpool.New(context.Background(), connStr)

	result, err := db.Exec(context.Background(), `
		CREATE TABLE IF NOT EXISTS bookings (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			user_id TEXT NOT NULL,
			flight_id TEXT NOT NULL,
			passenger_name TEXT NOT NULL,
			passenger_email TEXT NOT NULL,
			seat_count INTEGER NOT NULL,
			total_price NUMERIC NOT NULL,
			status TEXT DEFAULT 'CONFIRMED',
			created_at TIMESTAMPTZ DEFAULT NOW()
		);
	`)

	if err != nil {
		log.Fatalf("Migration failed: %v", err)
	}
	log.Printf("Migration executed, affected rows: %d", result.RowsAffected())
	return db
}

var flightClient pb.FlightServiceClient

func main() {
	db := initDB()
	defer db.Close()

	conn, _ := grpc.Dial(os.Getenv("FLIGHT_GRPC_ADDR"),
		grpc.WithInsecure(),
		grpc.WithChainUnaryInterceptor(
			retryInterceptor,
			authInterceptor(os.Getenv("FLIGHT_API_KEY")),
		),
	)
	flightClient = pb.NewFlightServiceClient(conn)

	r := gin.Default()

	r.GET("/flights", func(c *gin.Context) {
		resp, err := flightClient.SearchFlights(c.Request.Context(), &pb.SearchFlightsRequest{
			Origin: c.Query("origin"), Destination: c.Query("destination"),
		})
		if err != nil { c.JSON(500, gin.H{"error": err.Error()}); return }
		c.JSON(200, resp.Flights)
	})

	r.POST("/bookings", func(c *gin.Context) {
		type Req struct {
			UserID        string `json:"user_id"`
			FlightID      string `json:"flight_id"`
			PassengerName string `json:"passenger_name"`
			PassengerEmail string `json:"passenger_email"`
			SeatCount     int32  `json:"seat_count"`
		}
		var req Req
		c.BindJSON(&req)
		bookingID := uuid.NewString()

		// 1. GetFlight
		flight, err := flightClient.GetFlight(c.Request.Context(), &pb.GetFlightRequest{Id: req.FlightID})
		if err != nil { c.JSON(404, gin.H{"error": "flight not found"}); return }

		// 2. ReserveSeats (с Circuit Breaker + retry)
		err = cb.Execute(func() error {
			_, err := flightClient.ReserveSeats(c.Request.Context(), &pb.ReserveSeatsRequest{
				FlightId: req.FlightID, SeatCount: req.SeatCount, BookingId: bookingID,
			})
			return err
		})
		if err != nil { c.JSON(400, gin.H{"error": err.Error()}); return }

		// 3. Создать бронь
		total := float64(req.SeatCount) * flight.Price
		_, _ = db.Exec(context.Background(), `INSERT INTO bookings
			(user_id, flight_id, passenger_name, passenger_email, seat_count, total_price, status)
			VALUES ($1,$2,$3,$4,$5,$6,'CONFIRMED')`,
			req.UserID, req.FlightID, req.PassengerName, req.PassengerEmail, req.SeatCount, total)

		c.JSON(201, gin.H{"booking_id": bookingID, "total_price": total})
	})

	r.POST("/bookings/:id/cancel", func(c *gin.Context) {
		id := c.Param("id")
		_, _ = db.Exec(context.Background(), `UPDATE bookings SET status = 'CANCELLED' WHERE id = $1`, id)
		flightClient.ReleaseReservation(c.Request.Context(), &pb.ReleaseReservationRequest{BookingId: id})
		c.JSON(200, gin.H{"status": "cancelled"})
	})

	r.GET("/bookings/:id", func(c *gin.Context) {
		id := c.Param("id")
		var b struct {
			ID             string  `json:"id"`
			UserID         string  `json:"user_id"`
			FlightID       string  `json:"flight_id"`
			PassengerName  string  `json:"passenger_name"`
			TotalPrice     float64 `json:"total_price"`
			Status         string  `json:"status"`
		}
		err := db.QueryRow(context.Background(), `SELECT id, user_id, flight_id, passenger_name, total_price, status
			FROM bookings WHERE id = $1`, id).Scan(&b.ID, &b.UserID, &b.FlightID, &b.PassengerName, &b.TotalPrice, &b.Status)
		if err != nil {
			c.JSON(404, gin.H{"error": "booking not found"})
			return
		}
		c.JSON(200, b)
	})

	r.GET("/bookings", func(c *gin.Context) {
		userID := c.Query("user_id")
		rows, err := db.Query(context.Background(), `
			SELECT id, flight_id, passenger_name, total_price, status
			FROM bookings WHERE user_id = $1`, userID)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		var bookings []map[string]interface{}
		for rows.Next() {
			var (
				id            string
				flightID      string
				passengerName string
				totalPrice    float64
				status        int
			)
			if err := rows.Scan(&id, &flightID, &passengerName, &totalPrice, &status); err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}

			b := map[string]interface{}{
				"id":             id,
				"flight_id":      flightID,
				"passenger_name": passengerName,
				"total_price":    totalPrice,
				"status":         status,
			}
			bookings = append(bookings, b)
		}

		c.JSON(200, bookings)
	})

	log.Println("Booking Service ready on :8080")
	r.Run(":" + os.Getenv("PORT"))
}
