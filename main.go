package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// MessageStatus represents the status of a message.
type MessageStatus string

const (
	Pending  MessageStatus = "Pending"
	Approved MessageStatus = "Approved"
	Rejected MessageStatus = "Rejected"
)

// Message represents a message entity with approval status.
type Message struct {
	ID         string        `json:"id"`
	Content    string        `json:"content"`
	Recipient  string        `json:"recipient"`
	Status     MessageStatus `json:"status"`
	CreatedAt  time.Time     `json:"created_at"`
	ApprovedAt *time.Time    `json:"approved_at,omitempty"`
	RejectedAt *time.Time    `json:"rejected_at,omitempty"`
}

// CreateMessageRequest represents the request payload for creating a message.
type CreateMessageRequest struct {
	Content   string `json:"content"`
	Recipient string `json:"recipient"`
}

var storage = NewStorage()

func main() {
	port := "9999"
	envPort := os.Getenv("PORT")
	if envPort != "" {
		port = envPort
	}
	r := mux.NewRouter()
	r.HandleFunc("/messages", handleMessages).Methods("POST")
	r.HandleFunc("/messages/{id}/approve", handleApproveMessage).Methods("PATCH")
	r.HandleFunc("/messages/{id}/reject", handleRejectMessage).Methods("PATCH")
	r.HandleFunc("/messages", handleViewMessages).Methods("GET")

	log.Printf("Server started on %s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}

// handleMessages handles message creation (Maker).
func handleMessages(w http.ResponseWriter, r *http.Request) {
	var req CreateMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Content == "" || req.Recipient == "" {
		http.Error(w, "Content and Recipient are required", http.StatusBadRequest)
		return
	}

	msg := &Message{
		ID:        generateID(),
		Content:   req.Content,
		Recipient: req.Recipient,
		Status:    Pending,
		CreatedAt: time.Now(),
	}
	storage.AddMessage(msg)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(msg)
}

// handleApproveMessage handles message approval (Checker).
func handleApproveMessage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		http.Error(w, "Message ID is required", http.StatusBadRequest)
		return
	}

	msg, exists := storage.GetMessage(id)
	if !exists {
		http.Error(w, "Message not found", http.StatusNotFound)
		return
	}

	if msg.Status != Pending {
		http.Error(w, "Message is not in a pending state", http.StatusConflict)
		return
	}

	msg.Status = Approved
	msg.ApprovedAt = ptr(time.Now())
	storage.UpdateMessage(msg)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(msg)
}

// handleRejectMessage handles message rejection (Checker).
func handleRejectMessage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		http.Error(w, "Message ID is required", http.StatusBadRequest)
		return
	}

	msg, exists := storage.GetMessage(id)
	if !exists {
		http.Error(w, "Message not found", http.StatusNotFound)
		return
	}

	if msg.Status != Pending {
		http.Error(w, "Message is not in a pending state", http.StatusConflict)
		return
	}

	msg.Status = Rejected
	msg.RejectedAt = ptr(time.Now())
	storage.UpdateMessage(msg)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(msg)
}

// handleViewMessages handles viewing all messages.
func handleViewMessages(w http.ResponseWriter, r *http.Request) {
	messages := storage.GetAllMessages()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}

// generateID creates a unique ID for a message.
func generateID() string {
	return uuid.New().String()
}

// ptr returns a pointer to the given value.
func ptr[T any](v T) *T {
	return &v
}
