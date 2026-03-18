#!/bin/bash
set -e

echo "=== ТЕСТ 1: Поиск рейсов (SearchFlights + кэш) ==="
curl -s "http://localhost:8080/flights?origin=SVO&destination=LED" | jq
echo "→ Первый запрос = CACHE MISS (посмотри логи flight-service)"

sleep 2
curl -s "http://localhost:8080/flights?origin=SVO&destination=LED" | jq
echo "→ Второй запрос = CACHE HIT ✅"

echo "=== ТЕСТ 2: Создание бронирования ==="
BOOKING=$(curl -s -X POST http://localhost:8080/bookings \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "user123",
    "flight_id": "SU1234-2026-04-01",
    "passenger_name": "Иван Иванов",
    "passenger_email": "ivan@test.ru",
    "seat_count": 2
  }')

BOOKING_ID=$(echo $BOOKING | jq -r '.booking_id')
echo "Бронирование создано: $BOOKING_ID"
echo "$BOOKING" | jq

echo "=== ТЕСТ 3: Получение бронирования ==="
curl -s "http://localhost:8080/bookings/$BOOKING_ID" | jq

echo "=== ТЕСТ 4: Отмена бронирования (ReleaseReservation) ==="
curl -s -X POST "http://localhost:8080/bookings/$BOOKING_ID/cancel" | jq

echo "=== ТЕСТ 5: Проверка списка ==="
curl -s "http://localhost:8080/bookings?user_id=user123" | jq

echo "🎉 БАЗОВЫЕ ТЕСТЫ ПРОЙДЕНЫ!"
