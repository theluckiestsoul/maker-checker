package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestHandleMessages(t *testing.T) {
	storage = NewStorage()
	router := mux.NewRouter()
	router.HandleFunc("/messages", handleMessages).Methods("POST")

	t.Run("Valid Request", func(t *testing.T) {
		reqBody := `{"content":"Hello, World!","recipient":"user@example.com"}`
		req, _ := http.NewRequest("POST", "/messages", bytes.NewBufferString(reqBody))
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		var msg Message
		err := json.NewDecoder(rr.Body).Decode(&msg)
		assert.NoError(t, err)
		assert.Equal(t, "Hello, World!", msg.Content)
		assert.Equal(t, "user@example.com", msg.Recipient)
		assert.Equal(t, Pending, msg.Status)
		assert.NotEmpty(t, msg.ID)
		assert.WithinDuration(t, time.Now(), msg.CreatedAt, time.Second)
	})

	t.Run("Invalid Request Body", func(t *testing.T) {
		reqBody := `{"content":"Hello, World!"`
		req, _ := http.NewRequest("POST", "/messages", bytes.NewBufferString(reqBody))
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("Missing Content or Recipient", func(t *testing.T) {
		reqBody := `{"content":"","recipient":""}`
		req, _ := http.NewRequest("POST", "/messages", bytes.NewBufferString(reqBody))
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})
}

func TestHandleApproveMessage(t *testing.T) {
	storage = NewStorage()
	router := mux.NewRouter()
	router.HandleFunc("/messages/{id}/approve", handleApproveMessage).Methods("PATCH")

	msg := &Message{
		ID:        generateID(),
		Content:   "Hello, World!",
		Recipient: "user@example.com",
		Status:    Pending,
		CreatedAt: time.Now(),
	}
	storage.AddMessage(msg)

	t.Run("Valid Approval", func(t *testing.T) {
		req, _ := http.NewRequest("PATCH", "/messages/"+msg.ID+"/approve", nil)
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		var updatedMsg Message
		err := json.NewDecoder(rr.Body).Decode(&updatedMsg)
		assert.NoError(t, err)
		assert.Equal(t, Approved, updatedMsg.Status)
		assert.NotNil(t, updatedMsg.ApprovedAt)
		assert.WithinDuration(t, time.Now(), *updatedMsg.ApprovedAt, time.Second)
	})

	t.Run("Message Not Found", func(t *testing.T) {
		req, _ := http.NewRequest("PATCH", "/messages/"+uuid.New().String()+"/approve", nil)
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusNotFound, rr.Code)
	})

	t.Run("Message Not Pending", func(t *testing.T) {
		msg.Status = Approved
		storage.UpdateMessage(msg)

		req, _ := http.NewRequest("PATCH", "/messages/"+msg.ID+"/approve", nil)
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusConflict, rr.Code)
	})
}

func TestHandleRejectMessage(t *testing.T) {
	storage = NewStorage()
	router := mux.NewRouter()
	router.HandleFunc("/messages/{id}/reject", handleRejectMessage).Methods("PATCH")

	msg := &Message{
		ID:        generateID(),
		Content:   "Hello, World!",
		Recipient: "user@example.com",
		Status:    Pending,
		CreatedAt: time.Now(),
	}
	storage.AddMessage(msg)

	t.Run("Valid Rejection", func(t *testing.T) {
		req, _ := http.NewRequest("PATCH", "/messages/"+msg.ID+"/reject", nil)
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		var updatedMsg Message
		err := json.NewDecoder(rr.Body).Decode(&updatedMsg)
		assert.NoError(t, err)
		assert.Equal(t, Rejected, updatedMsg.Status)
		assert.NotNil(t, updatedMsg.RejectedAt)
		assert.WithinDuration(t, time.Now(), *updatedMsg.RejectedAt, time.Second)
	})

	t.Run("Message Not Found", func(t *testing.T) {
		req, _ := http.NewRequest("PATCH", "/messages/"+uuid.New().String()+"/reject", nil)
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusNotFound, rr.Code)
	})

	t.Run("Message Not Pending", func(t *testing.T) {
		msg.Status = Approved
		storage.UpdateMessage(msg)

		req, _ := http.NewRequest("PATCH", "/messages/"+msg.ID+"/reject", nil)
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusConflict, rr.Code)
	})
}

func TestHandleViewMessages(t *testing.T) {
	storage = NewStorage()
	router := mux.NewRouter()
	router.HandleFunc("/messages", handleViewMessages).Methods("GET")

	msg1 := &Message{
		ID:        generateID(),
		Content:   "Hello, World!",
		Recipient: "user1@example.com",
		Status:    Pending,
		CreatedAt: time.Now(),
	}
	msg2 := &Message{
		ID:         generateID(),
		Content:    "Goodbye, World!",
		Recipient:  "user2@example.com",
		Status:     Approved,
		CreatedAt:  time.Now(),
		ApprovedAt: ptr(time.Now()),
	}
	storage.AddMessage(msg1)
	storage.AddMessage(msg2)

	t.Run("View All Messages", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/messages", nil)
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		var messages []*Message
		err := json.NewDecoder(rr.Body).Decode(&messages)
		assert.NoError(t, err)
		assert.Len(t, messages, 2)
		assert.Equal(t, "Hello, World!", messages[0].Content)
		assert.Equal(t, "Goodbye, World!", messages[1].Content)
	})
}

func TestStorage(t *testing.T) {
	storage := NewStorage()

	t.Run("Add and Get Message", func(t *testing.T) {
		msg := &Message{
			ID:        generateID(),
			Content:   "Hello, World!",
			Recipient: "user@example.com",
			Status:    Pending,
			CreatedAt: time.Now(),
		}
		storage.AddMessage(msg)

		retrievedMsg, exists := storage.GetMessage(msg.ID)
		assert.True(t, exists)
		assert.Equal(t, msg, retrievedMsg)
	})

	t.Run("Get Non-Existent Message", func(t *testing.T) {
		_, exists := storage.GetMessage(uuid.New().String())
		assert.False(t, exists)
	})

	t.Run("Update Message", func(t *testing.T) {
		msg := &Message{
			ID:        generateID(),
			Content:   "Hello, World!",
			Recipient: "user@example.com",
			Status:    Pending,
			CreatedAt: time.Now(),
		}
		storage.AddMessage(msg)

		msg.Status = Approved
		msg.ApprovedAt = ptr(time.Now())
		storage.UpdateMessage(msg)

		retrievedMsg, exists := storage.GetMessage(msg.ID)
		assert.True(t, exists)
		assert.Equal(t, Approved, retrievedMsg.Status)
		assert.NotNil(t, retrievedMsg.ApprovedAt)
	})

	t.Run("Get All Messages", func(t *testing.T) {
		storage = NewStorage() // Reset storage
		msg1 := &Message{
			ID:        generateID(),
			Content:   "Hello, World!",
			Recipient: "user1@example.com",
			Status:    Pending,
			CreatedAt: time.Now(),
		}
		msg2 := &Message{
			ID:         generateID(),
			Content:    "Goodbye, World!",
			Recipient:  "user2@example.com",
			Status:     Approved,
			CreatedAt:  time.Now(),
			ApprovedAt: ptr(time.Now()),
		}
		storage.AddMessage(msg1)
		storage.AddMessage(msg2)

		messages := storage.GetAllMessages()
		assert.Len(t, messages, 2)
	})
}
