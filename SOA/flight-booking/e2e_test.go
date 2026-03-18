package e2e

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"os"
	"testing"
	"time"

	"flight-booking/testutil"
)

const baseURL = "http://localhost:8080"

type Flight struct {
	ID             string  `json:"id"`
	FlightNumber   string  `json:"flight_number"`
	Origin         string  `json:"origin"`
	Destination    string  `json:"destination"`
	AvailableSeats int32   `json:"available_seats"`
	Price          float64 `json:"price"`
}

type BookingResponse struct {
	ID string `json:"id"`
}

func TestMain(m *testing.M) {
	// берем URL базы из переменной окружения, или дефолт
	conn := os.Getenv("DATABASE_URL")
	if conn == "" {
		conn = "postgres://user:password@localhost:5434/flight?sslmode=disable"
	}

	// ждем готовности базы
	testutil.WaitForDB(conn)

	// создаем таблицу и заливаем тестовые данные
	testutil.SeedFlights(conn)

	code := m.Run()
	os.Exit(code)
}

func TestBookingFlowE2E(t *testing.T) {
	client := &http.Client{Timeout: 5 * time.Second}

	// 1. Search flights
	resp, err := client.Get(baseURL + "/flights?origin=SVO&destination=LED")
	if err != nil {
		t.Fatalf("failed to search flights: %v", err)
	}
	defer resp.Body.Close()

	var respData struct {
		Flights []Flight `json:"flights"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&respData); err != nil {
		t.Fatalf("decode error: %v", err)
	}

	if len(respData.Flights) == 0 {
		t.Fatal("no flights found — убедись, что сервер использует ту же базу, что и SeedFlights")
	}

	flight := respData.Flights[0]
	initialSeats := flight.AvailableSeats

	// 2. Create booking
	reqBody := map[string]interface{}{
		"user_id":         "e2e-user",
		"flight_id":       flight.ID,
		"passenger_name":  "Test User",
		"passenger_email": "test@example.com",
		"seat_count":      1,
	}

	body, _ := json.Marshal(reqBody)
	resp, err = client.Post(baseURL+"/bookings", "application/json", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("failed to create booking: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("unexpected status: %d", resp.StatusCode)
	}

	var booking BookingResponse
	if err := json.NewDecoder(resp.Body).Decode(&booking); err != nil {
		t.Fatalf("decode booking error: %v", err)
	}

	if booking.ID == "" {
		t.Fatal("empty booking id")
	}

	// 3. Check seats decreased
	resp, err = client.Get(baseURL + "/flights/" + flight.ID)
	if err != nil {
		t.Fatalf("failed to get flight: %v", err)
	}
	defer resp.Body.Close()

	var updatedFlight Flight
	if err := json.NewDecoder(resp.Body).Decode(&updatedFlight); err != nil {
		t.Fatalf("decode updated flight error: %v", err)
	}

	if updatedFlight.AvailableSeats != initialSeats-1 {
		t.Fatalf("seats not decreased: expected %d, got %d",
			initialSeats-1, updatedFlight.AvailableSeats)
	}

	// 4. Cancel booking
	req, _ := http.NewRequestWithContext(context.Background(), "POST", baseURL+"/bookings/"+booking.ID+"/cancel", nil)
	resp, err = client.Do(req)
	if err != nil {
		t.Fatalf("cancel failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("cancel failed status: %d", resp.StatusCode)
	}

	// 5. Check seats restored
	resp, err = client.Get(baseURL + "/flights/" + flight.ID)
	if err != nil {
		t.Fatalf("failed to get flight: %v", err)
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&updatedFlight); err != nil {
		t.Fatalf("decode updated flight error: %v", err)
	}

	if updatedFlight.AvailableSeats != initialSeats {
		t.Fatalf("seats not restored: expected %d, got %d",
			initialSeats, updatedFlight.AvailableSeats)
	}
}
