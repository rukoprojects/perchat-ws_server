package services

import (
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/gorilla/websocket"
)

type Message struct {
	RecipientID      string `json:"recipientID"`
	SenderID         string `json:"senderID"`
	EncryptedContent string `json:"encryptedContent"`
}

type Server struct {
	clients   map[string]*websocket.Conn // Map userID to connection
	mu        sync.RWMutex
	broadcast chan Message
	logger    *log.Logger
}

var upgrader = websocket.Upgrader{}

func newLogger() (*log.Logger, error) {
	logFile, err := os.OpenFile("server.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}
	return log.New(logFile, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile), nil
}

func NewServer() *Server {
	logger, err := newLogger()
	if err != nil {
		log.Fatal(err)
	}

	return &Server{
		clients:   make(map[string]*websocket.Conn),
		broadcast: make(chan Message),
		mu:        sync.RWMutex{},
		logger:    logger,
	}
}

func (s *Server) addClient(userID string, conn *websocket.Conn) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.clients[userID] = conn
}

func (s *Server) removeClient(userID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.clients, userID)
}

func (s *Server) HandleClient(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Couldn't upgrade to websocket connection!", http.StatusBadRequest)
		s.logger.Printf("Error upgrading connection: %v", err)
		return
	}
	defer conn.Close()

	var initialMsg Message
	if err := conn.ReadJSON(&initialMsg); err != nil {
		s.logger.Printf("Error reading initial message: %v", err)
		return
	}

	s.addClient(initialMsg.SenderID, conn)
	defer s.removeClient(initialMsg.SenderID)

	go s.handleMessages()

	for {
		var msg Message
		if err := conn.ReadJSON(&msg); err != nil {
			s.logger.Printf("Error reading message: %v", err)
			return
		}

		s.broadcast <- msg
	}
}

func (s *Server) handleMessages() {
	for msg := range s.broadcast {
		s.mu.RLock()
		for clientID, clientConn := range s.clients {
			if err := clientConn.WriteJSON(msg); err != nil {
				s.logger.Printf("Error sending message to %s: %v", clientID, err)
			}
		}
		s.mu.RUnlock()
	}
}
