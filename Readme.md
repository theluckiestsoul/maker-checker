# Maker-Checker Approval Process Service

This is a simple Go application implementing a Maker-Checker approval process. The service allows messages to be created and then either approved or rejected by a checker. Approved messages are marked accordingly, while rejected messages are flagged and not sent.

The application runs an HTTP server on port 9999. If you want to change the port, you can set the `PORT` environment variable.

## Features

- **Maker**: Creates a new message (pending approval).
- **Checker**: Approves or rejects a pending message.
- **View**: Retrieve all messages with their status (Pending, Approved, Rejected).

## Endpoints

### 1. Create a Message (Maker)

**URL**: `/messages`  
**Method**: `POST`

#### Request Payload

http://localhost:9999/messages

```json
{
  "content": "Hello, this is a test message",
  "recipient": "test@example.com"
}
```

#### Response

```json
{
  "id": "3aada4d8-4b1d-4443-933b-8cfbf004fb4b",
  "content": "Hello, this is a test message",
  "recipient": "test@example.com",
  "status": "Pending",
  "created_at": "2024-12-12T06:08:27.490306+05:30"
}
```

### 2. Approve a Message (Checker)

**URL**: `/messages/{id}/approve`  
**Method**: `PATCH`

### Request Payload

```json
http://localhost:9999/messages/3aada4d8-4b1d-4443-933b-8cfbf004fb4b/approve
```

#### Response

```json
{
  "id": "3aada4d8-4b1d-4443-933b-8cfbf004fb4b",
  "content": "Hello, this is a test message",
  "recipient": "test@example.com",
  "status": "Approved",
  "created_at": "2024-12-12T06:08:27.490306+05:30",
  "approved_at": "2024-12-12T06:19:13.252364+05:30"
}
```

### 3. Reject a Message (Checker)

**URL**: `/messages/{id}/reject`  
**Method**: `PATCH`

#### Request Payload

```json
http://localhost:9999/messages/dc79f491-8fdd-403b-9fd0-ef5ea5647889/reject
```

#### Response

```json
{
  "id": "dc79f491-8fdd-403b-9fd0-ef5ea5647889",
  "content": "Hello, this is a test message",
  "recipient": "test@example.com",
  "status": "Pending",
  "created_at": "2024-12-12T06:19:46.821784+05:30"
}
```

### 4. View All Messages

**URL**: `/messages`  
**Method**: `GET`

#### Response

```json
[
  {
    "id": "3aada4d8-4b1d-4443-933b-8cfbf004fb4b",
    "content": "Hello, this is a test message",
    "recipient": "test@example.com",
    "status": "Approved",
    "created_at": "2024-12-12T06:08:27.490306+05:30",
    "approved_at": "2024-12-12T06:19:13.252364+05:30"
  },
  {
    "id": "dc79f491-8fdd-403b-9fd0-ef5ea5647889",
    "content": "Hello, this is a test message",
    "recipient": "test@example.com",
    "status": "Rejected",
    "created_at": "2024-12-12T06:19:46.821784+05:30",
    "rejected_at": "2024-12-12T06:21:19.745511+05:30"
  },
  {
    "id": "d3eb30ef-38ba-475f-9175-420b177bbff4",
    "content": "Hello, this is a test message",
    "recipient": "test@example.com",
    "status": "Pending",
    "created_at": "2024-12-12T06:21:27.02618+05:30"
  }
]
```

## How to Run

### Prerequisites

- **Go**: Ensure you have Go installed (version 1.18 or higher).

### Steps

1. Clone the repository:

   ```bash
   git clone github.com/theluckiestsoul/maker-checker
   cd maker-checker
   ```

2. Run the application:

   ```bash
   go run .
   ```

3. The server will start on `http://localhost:9999`.

## Testing the Endpoints

You can use tools like **Postman**, **curl**, or write your own HTTP client to test the endpoints.

### Example with `curl`

#### Create a Message:

```bash
curl -X POST http://localhost:9999/messages \
  -H "Content-Type: application/json" \
  -d '{"content": "Test Message", "recipient": "user@example.com"}'
```

#### Approve a Message:

```bash
curl -X PATCH "http://localhost:9999/messages/{id}/approve"
```

#### Reject a Message:

```bash
curl -X PATCH "http://localhost:9999/messages/{id}/reject"
```

#### View All Messages:

```bash
curl -X GET http://localhost:9999/messages
```

## Project Structure

- `main.go`: Contains the application logic and HTTP handlers.
- `storage.go`: Implements an in-memory storage for messages.

## Future Improvements

- Add persistent storage (e.g., PostgreSQL).
- Enhance validation and error handling.
- Add authentication and authorization.
- Implement logging and monitoring.

## License

This project is licensed under the MIT License.
