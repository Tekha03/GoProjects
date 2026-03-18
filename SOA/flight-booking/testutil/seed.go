package testutil

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
)

func WaitForDB(connString string) {
	for i := 0; i < 10; i++ {
		conn, err := pgx.Connect(context.Background(), connString)
		if err == nil {
			conn.Close(context.Background())
			return
		}
		time.Sleep(2 * time.Second)
	}
	log.Fatal("database not ready")
}

func SeedFlights(connString string) {
	ctx := context.Background()

	conn, err := pgx.Connect(ctx, connString)
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}
	defer conn.Close(ctx)

	// создаем таблицу, если нет
	_, err = conn.Exec(ctx, `
	CREATE TABLE IF NOT EXISTS flights (
		id TEXT PRIMARY KEY,
		flight_number VARCHAR(20),
		origin VARCHAR(255),
		destination VARCHAR(255),
		departure_time TIMESTAMP,
		arrival_time TIMESTAMP,
		total_seats INT,
		available_seats INT,
		price NUMERIC,
		status INT
	)
	`)
	if err != nil {
		log.Fatalf("failed to create flights table: %v", err)
	}

	// чистим таблицу
	_, err = conn.Exec(ctx, `DELETE FROM flights`)
	if err != nil {
		log.Fatalf("failed to clean flights: %v", err)
	}

	// добавляем несколько тестовых рейсов
	flights := []struct {
		ID, FlightNumber, Origin, Destination string
		Seats                                int
		Price                                float64
	}{
		{"flight1", "SU100", "SVO", "LED", 100, 100.0},
		{"flight2", "SU101", "SVO", "JFK", 200, 500.0},
		{"flight3", "SU102", "LED", "SVO", 150, 120.0},
	}

	for _, f := range flights {
		_, err = conn.Exec(ctx, `
		INSERT INTO flights (
			id, flight_number, origin, destination,
			departure_time, arrival_time,
			total_seats, available_seats, price, status
		) VALUES ($1,$2,$3,$4,NOW(), NOW() + INTERVAL '2 hour',$5,$5,$6,0)
		`, f.ID, f.FlightNumber, f.Origin, f.Destination, f.Seats, f.Price)
		if err != nil {
			log.Fatalf("failed to insert flight %s: %v", f.ID, err)
		}
	}
}
