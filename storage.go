package main

import "sync"

// Storage provides thread-safe in-memory storage for messages.
type Storage struct {
	mu       sync.Mutex
	messages map[string]*Message
}

// NewStorage initializes a new Storage instance.
func NewStorage() *Storage {
	return &Storage{
		messages: make(map[string]*Message),
	}
}

// AddMessage stores a new message in memory.
func (s *Storage) AddMessage(msg *Message) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.messages[msg.ID] = msg
}

// GetMessage retrieves a message by its ID.
func (s *Storage) GetMessage(id string) (*Message, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	msg, exists := s.messages[id]
	return msg, exists
}

// GetAllMessages retrieves all messages.
func (s *Storage) GetAllMessages() []*Message {
	s.mu.Lock()
	defer s.mu.Unlock()
	messages := make([]*Message, 0, len(s.messages))
	for _, msg := range s.messages {
		messages = append(messages, msg)
	}
	return messages
}

// UpdateMessage updates an existing message in memory.
func (s *Storage) UpdateMessage(msg *Message) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.messages[msg.ID] = msg
}
