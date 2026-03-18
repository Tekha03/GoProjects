package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "flight-booking/flight-service/proto/flight/v1"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type FlightServer struct {
	pb.UnimplementedFlightServiceServer
	db     *pgxpool.Pool
	rdb    *redis.ClusterClient
	apiKey string
}

func (s *FlightServer) authenticate(ctx context.Context) error {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return status.Error(codes.Unauthenticated, "no metadata")
	}
	keys := md.Get("x-api-key")
	if len(keys) == 0 || keys[0] != s.apiKey {
		return status.Error(codes.Unauthenticated, "invalid api key")
	}
	return nil
}

func (s *FlightServer) getCache(key string) (string, bool) {
	val, err := s.rdb.Get(context.Background(), key).Result()
	if err == redis.Nil {
		log.Printf("CACHE MISS %s", key)
		return "", false
	}
	if err != nil {
		return "", false
	}
	log.Printf("CACHE HIT %s", key)
	return val, true
}

func (s *FlightServer) setCache(key string, value interface{}, ttl time.Duration) {
	data, _ := json.Marshal(value)
	s.rdb.Set(context.Background(), key, data, ttl)
}

func (s *FlightServer) invalidateFlight(flightID string) {
	s.rdb.Del(context.Background(), "flight:"+flightID)
}

func initDB() *pgxpool.Pool {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))
	db, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		log.Fatal(err)
	}

	result, err := db.Exec(context.Background(), `
		CREATE TABLE IF NOT EXISTS flights (
			id TEXT PRIMARY KEY,
			flight_number TEXT NOT NULL,
			origin TEXT NOT NULL,
			destination TEXT NOT NULL,
			departure_time TIMESTAMPTZ NOT NULL,
			arrival_time TIMESTAMPTZ NOT NULL,
			total_seats INTEGER CHECK (total_seats > 0),
			available_seats INTEGER CHECK (available_seats >= 0),
			price NUMERIC CHECK (price > 0),
			status TEXT DEFAULT 'SCHEDULED',
			UNIQUE (flight_number, departure_time)
		);

		CREATE TABLE IF NOT EXISTS seat_reservations (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			flight_id TEXT REFERENCES flights(id),
			booking_id TEXT UNIQUE NOT NULL,
			seat_count INTEGER CHECK (seat_count > 0),
			status TEXT DEFAULT 'ACTIVE',
			created_at TIMESTAMPTZ DEFAULT NOW()
		);
	`)

	if err != nil {
		log.Fatalf("Migration failed: %v", err)
	}
	log.Printf("Migration executed, affected rows: %d", result.RowsAffected())
	return db
}

func (s *FlightServer) SearchFlights(ctx context.Context, req *pb.SearchFlightsRequest) (*pb.SearchFlightsResponse, error) {
    log.Printf("🔍 SearchFlights: origin=%s, destination=%s, date=%v",
        req.Origin, req.Destination, req.Date)

    if err := s.authenticate(ctx); err != nil {
        log.Printf("Auth failed: %v", err)
        return nil, err
    }

    var dbCheck int
    err := s.db.QueryRow(ctx, "SELECT 1").Scan(&dbCheck)
    if err != nil {
        log.Printf("DB connection error: %v", err)
    } else {
        log.Printf("DB connection OK")
    }

    var totalFlights int
    s.db.QueryRow(ctx, "SELECT COUNT(*) FROM flights").Scan(&totalFlights)
    log.Printf("Total flights in DB: %d", totalFlights)

    dateFilter := ""
    if req.Date != nil {
        dateUTC := req.Date.AsTime().UTC()
        dateFilter = dateUTC.Format("2006-01-02")
        log.Printf("Date filter: %s (original: %v)", dateFilter, req.Date.AsTime())
    }

    cacheKey := fmt.Sprintf("search:%s:%s:%s", req.Origin, req.Destination, dateFilter)
    log.Printf("Cache key: %s", cacheKey)

    if val, ok := s.getCache(cacheKey); ok {
        log.Printf("CACHE HIT %s", cacheKey)
        var flights []*pb.Flight
        if err := json.Unmarshal([]byte(val), &flights); err == nil {
            return &pb.SearchFlightsResponse{Flights: flights}, nil
        }
    }
    log.Printf("CACHE MISS %s", cacheKey)

    query := `
        SELECT id, flight_number, origin, destination, departure_time, arrival_time,
               total_seats, available_seats, price, status
        FROM flights
        WHERE origin = $1 AND destination = $2 AND status = 'SCHEDULED'
    `
    args := []interface{}{req.Origin, req.Destination}

    if dateFilter != "" {
        query += " AND departure_time::date = $3::date"
        args = append(args, dateFilter)
    }

    log.Printf("📝 Executing query: %s with args: %v", query, args)

    rows, err := s.db.Query(ctx, query, args...)
    if err != nil {
        log.Printf("Query error: %v", err)
        return nil, status.Error(codes.Internal, err.Error())
    }
    defer rows.Close()

    var flights []*pb.Flight
    for rows.Next() {
        var f pb.Flight
        var dep, arr time.Time
        var price float64
        var statusStr string

        if err := rows.Scan(
            &f.Id, &f.FlightNumber, &f.Origin, &f.Destination,
            &dep, &arr, &f.TotalSeats, &f.AvailableSeats, &price, &statusStr,
        ); err != nil {
            log.Printf("Row scan error: %v", err)
            continue
        }

        switch statusStr {
        case "SCHEDULED":
            f.Status = pb.FlightStatus_SCHEDULED
        case "DEPARTED":
            f.Status = pb.FlightStatus_DEPARTED
        case "CANCELLED":
            f.Status = pb.FlightStatus_CANCELLED
        case "COMPLETED":
            f.Status = pb.FlightStatus_COMPLETED
        default:
            log.Printf("Unknown status string: %s, using SCHEDULED", statusStr)
            f.Status = pb.FlightStatus_SCHEDULED
        }

        f.DepartureTime = timestamppb.New(dep)
        f.ArrivalTime = timestamppb.New(arr)
        f.Price = price
        flights = append(flights, &f)
        log.Printf("Found flight: %s %s->%s [status=%s]",
            f.FlightNumber, f.Origin, f.Destination, statusStr)
    }

    log.Printf("Found %d flights", len(flights))

    if len(flights) == 0 {
        debugRows, _ := s.db.Query(ctx, "SELECT id, flight_number, origin, destination, status FROM flights LIMIT 5")
        defer debugRows.Close()
        log.Printf("🔍 All flights in DB (first 5):")
        for debugRows.Next() {
            var id, fn, origin, dest, status string
            debugRows.Scan(&id, &fn, &origin, &dest, &status)
            log.Printf("   - %s: %s %s->%s [%s]", id, fn, origin, dest, status)
        }
    }

    if data, err := json.Marshal(flights); err == nil {
        s.rdb.Set(ctx, cacheKey, data, 5*time.Minute)
    }

    return &pb.SearchFlightsResponse{Flights: flights}, nil
}

func (s *FlightServer) GetFlight(ctx context.Context, req *pb.GetFlightRequest) (*pb.Flight, error) {
    log.Printf("🔍 GetFlight: id=%s", req.Id)

    if err := s.authenticate(ctx); err != nil {
        log.Printf("Auth failed: %v", err)
        return nil, err
    }

    cacheKey := "flight:" + req.Id
    if val, ok := s.getCache(cacheKey); ok {
        log.Printf("🎯 CACHE HIT %s", cacheKey)
        var f pb.Flight
        if err := json.Unmarshal([]byte(val), &f); err == nil {
            return &f, nil
        }
    }
    log.Printf("💫 CACHE MISS %s", cacheKey)

    var f pb.Flight
    var dep, arr time.Time
    var price float64
    var statusStr string

    err := s.db.QueryRow(ctx, `SELECT id, flight_number, origin, destination, departure_time, arrival_time,
        total_seats, available_seats, price, status FROM flights WHERE id = $1`, req.Id).
        Scan(&f.Id, &f.FlightNumber, &f.Origin, &f.Destination, &dep, &arr,
             &f.TotalSeats, &f.AvailableSeats, &price, &statusStr)

    if err != nil {
        log.Printf("Flight not found: %v", err)
        return nil, status.Error(codes.NotFound, "flight not found")
    }

    switch statusStr {
    case "SCHEDULED":
        f.Status = pb.FlightStatus_SCHEDULED
    case "DEPARTED":
        f.Status = pb.FlightStatus_DEPARTED
    case "CANCELLED":
        f.Status = pb.FlightStatus_CANCELLED
    case "COMPLETED":
        f.Status = pb.FlightStatus_COMPLETED
    default:
        log.Printf("Unknown status string: %s, using SCHEDULED", statusStr)
        f.Status = pb.FlightStatus_SCHEDULED
    }

    f.DepartureTime = timestamppb.New(dep)
    f.ArrivalTime = timestamppb.New(arr)
    f.Price = price

    log.Printf("Found flight: %s %s->%s [status=%s]",
        f.FlightNumber, f.Origin, f.Destination, statusStr)

    if data, err := json.Marshal(f); err == nil {
        s.rdb.Set(ctx, cacheKey, data, 7*time.Minute)
    }

    return &f, nil
}

func (s *FlightServer) ReserveSeats(ctx context.Context, req *pb.ReserveSeatsRequest) (*pb.ReserveSeatsResponse, error) {
	if err := s.authenticate(ctx); err != nil { return nil, err }
	s.invalidateFlight(req.FlightId)

	tx, err := s.db.Begin(ctx)
	if err != nil { return nil, status.Error(codes.Internal, err.Error()) }
	defer tx.Rollback(ctx)

	var avail int32
	err = tx.QueryRow(ctx, `SELECT available_seats FROM flights WHERE id = $1 FOR UPDATE`, req.FlightId).Scan(&avail)
	if err != nil { return nil, status.Error(codes.NotFound, "flight not found") }
	if avail < req.SeatCount {
		return nil, status.Error(codes.ResourceExhausted, "not enough seats")
	}

	_, err = tx.Exec(ctx, `INSERT INTO seat_reservations (flight_id, booking_id, seat_count, status)
		VALUES ($1, $2, $3, 'ACTIVE') ON CONFLICT (booking_id) DO NOTHING`,
		req.FlightId, req.BookingId, req.SeatCount)
	if err != nil { return nil, status.Error(codes.Internal, err.Error()) }

	_, err = tx.Exec(ctx, `UPDATE flights SET available_seats = available_seats - $1 WHERE id = $2`, req.SeatCount, req.FlightId)
	if err != nil { return nil, err }

	tx.Commit(ctx)
	log.Printf("RESERVED %d seats for booking %s", req.SeatCount, req.BookingId)
	return &pb.ReserveSeatsResponse{}, nil
}

func (s *FlightServer) ReleaseReservation(ctx context.Context, req *pb.ReleaseReservationRequest) (*pb.ReleaseReservationResponse, error) {
	if err := s.authenticate(ctx); err != nil { return nil, err }

	tx, _ := s.db.Begin(ctx)
	defer tx.Rollback(ctx)

	var flightID string
	var seatCount int32
	err := tx.QueryRow(ctx, `SELECT flight_id, seat_count FROM seat_reservations
		WHERE booking_id = $1 AND status = 'ACTIVE'`, req.BookingId).Scan(&flightID, &seatCount)
	if err != nil {
		return nil, status.Error(codes.NotFound, "reservation not found")
	}

	_, _ = tx.Exec(ctx, `UPDATE seat_reservations SET status = 'RELEASED' WHERE booking_id = $1`, req.BookingId)
	_, _ = tx.Exec(ctx, `UPDATE flights SET available_seats = available_seats + $1 WHERE id = $2`, seatCount, flightID)
	tx.Commit(ctx)

	s.invalidateFlight(flightID)
	log.Printf("RELEASED booking %s", req.BookingId)
	return &pb.ReleaseReservationResponse{}, nil
}

func main() {
	db := initDB()
	defer db.Close()

	rdb := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:    []string{"redis-node-0:6379", "redis-node-1:6379", "redis-node-2:6379", "redis-node-3:6379", "redis-node-4:6379", "redis-node-5:6379"},
		Password: os.Getenv("REDIS_PASSWORD"),
	})
	defer rdb.Close()

	lis, _ := net.Listen("tcp", ":"+os.Getenv("PORT"))
	srv := grpc.NewServer()
	pb.RegisterFlightServiceServer(srv, &FlightServer{
		db:     db,
		rdb:    rdb,
		apiKey: os.Getenv("API_KEY"),
	})

	log.Println("Flight Service ready on :50051")
	srv.Serve(lis)
}
