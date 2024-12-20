# Health Reservation System

A Twirp-based API for managing health provider availability and client reservations. This system supports provider availability management, slot reservations, confirmation, and cleanup of expired reservations.

## Features

- Providers can set their availability.
- Clients can reserve available slots and confirm reservations.
- Retrieve reserved slots by provider or client.
- Automatic cleanup of expired reservations.
- SQLite database backend with GORM.

## Prerequisites

- Go 1.21 or higher
- SQLite database
- `protoc` for generating gRPC and Twirp code (if modifying the `.proto` file)

## Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/manueldelreal/health-reservation-system.git
   cd health-reservation-system
   ```

2. Use the `Makefile` for streamlined setup and operations.

## Using the Makefile

### Setup
Install dependencies:
```bash
make setup
```

### Build
Compile the server:
```bash
make build
```

### Run
Start the server:
```bash
make run
```

The server will start on `http://localhost:8080`.

### Testing
Run all tests:
```bash
make test
```

### Generate Proto Code
Generate Twirp and Go code from `.proto` files:
```bash
make generate
```

### Cleanup
Remove build artifacts:
```bash
make clean
```

### Dependency Check
Ensure `protoc` is installed:
```bash
make protoc-check
```

## Database Configuration

The system uses SQLite with GORM for database management. The database file is `health_reservation.db`, and migrations are applied automatically.

## API Endpoints

### Base URL

Twirp service base URL: `http://localhost:8080/twirp/reservation.ReservationService/`

### RPC Methods

#### 1. **SetAvailability**

- **Description:** Allows a provider to set availability.
- **Endpoint:** `SetAvailability`
- **Request:**
  ```json
  {
    "provider_id": "provider_123",
    "time_slots": [
      { "start_time": "2024-12-20T08:00:00Z", "end_time": "2024-12-20T09:00:00Z" }
    ]
  }
  ```
- **Response:**
  ```json
  { "message": "Availability set successfully" }
  ```

#### 2. **GetAvailableSlots**

- **Description:** Retrieves available slots for a provider.
- **Endpoint:** `GetAvailableSlots`
- **Request:**
  ```json
  {
    "provider_id": "provider_123",
    "date": "2024-12-20"
  }
  ```
- **Response:**
  ```json
  {
    "slots": [
      {
        "id": "slot_123",
        "start_time": "2024-12-20T08:00:00Z",
        "end_time": "2024-12-20T08:15:00Z",
        "status": "Available"
      }
    ]
  }
  ```

#### 3. **ReserveSlot**

- **Description:** Reserves an available slot.
- **Endpoint:** `ReserveSlot`
- **Request:**
  ```json
  {
    "slot_id": "slot_123",
    "client_id": "client_456"
  }
  ```
- **Response:**
  ```json
  {
    "reservation_id": "slot_123",
    "message": "Slot reserved successfully"
  }
  ```

#### 4. **ConfirmReservation**

- **Description:** Confirms a reservation.
- **Endpoint:** `ConfirmReservation`
- **Request:**
  ```json
  {
    "reservation_id": "reservation_123"
  }
  ```
- **Response:**
  ```json
  { "message": "Reservation confirmed" }
  ```

#### 5. **GetReservedSlotsByProvider**

- **Description:** Retrieves reservations for a provider, optionally filtered by date.
- **Endpoint:** `GetReservedSlotsByProvider`
- **Request:**
  ```json
  {
    "provider_id": "provider_123",
    "date": "2024-12-20"
  }
  ```
- **Response:**
  ```json
  {
    "reservations": [
      {
        "reservation_id": "reservation_123",
        "client_id": "client_456",
        "provider_id": "provider_123",
        "status": "Confirmed",
        "start_time": "2024-12-20T08:00:00Z",
        "end_time": "2024-12-20T08:15:00Z"
      }
    ]
  }
  ```

#### 6. **GetReservedSlotsByClient**

- **Description:** Retrieves reservations for a client, optionally filtered by date.
- **Endpoint:** `GetReservedSlotsByClient`
- **Request:**
  ```json
  {
    "client_id": "client_456",
    "date": "2024-12-20"
  }
  ```
- **Response:**
  ```json
  {
    "reservations": [
      {
        "reservation_id": "reservation_123",
        "client_id": "client_456",
        "provider_id": "provider_123",
        "status": "Confirmed",
        "start_time": "2024-12-20T08:00:00Z",
        "end_time": "2024-12-20T08:15:00Z"
      }
    ]
  }
  ```

## Cleanup Task

The server includes an automated task to clean up expired reservations every minute. Expired reservations are marked as "Available" and moved back to the slots table.

## License

MIT License

