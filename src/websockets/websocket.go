package websockets

import (
	"log"
	"net/http"
	"sync"
	"tic-tac-toe-game/src/database"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WebSocketService struct {
	db                         *database.InMemoryGameSessionDB
	GameSessionConnectionStore *GameSessionConnectionStore
}

func NewWebSocketService(db *database.InMemoryGameSessionDB) *WebSocketService {
	cs := newConnectionStore()
	gameSessionConnectionStore := newGameSessionConnectionStore(cs)

	return &WebSocketService{
		db:                         db,
		GameSessionConnectionStore: gameSessionConnectionStore,
	}
}

func (wss *WebSocketService) StartWebSocketServer() {
	http.HandleFunc("/game-session", func(w http.ResponseWriter, r *http.Request) {
		gameSessionId := r.URL.Query().Get("gameSessionId")

		if gameSessionId == "" {
			http.Error(w, "Missing gameSessionId query parameter", http.StatusBadRequest)
			return
		}

		wss.GameSessionConnectionStore.Echo(w, r, gameSessionId, wss.db)
	})

	go func() {
		log.Println("WebSocket server started on :8090")
		log.Fatal(http.ListenAndServe(":8090", nil))
	}()
}

type ConnectionStore struct {
	connections map[*websocket.Conn]bool
	mu          sync.Mutex
}

func newConnectionStore() *ConnectionStore {
	return &ConnectionStore{
		connections: make(map[*websocket.Conn]bool),
	}
}

func (cs *ConnectionStore) AddConnection(conn *websocket.Conn) {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	cs.connections[conn] = true
}

func (cs *ConnectionStore) RemoveConnection(conn *websocket.Conn) {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	delete(cs.connections, conn)
}

func (cs *ConnectionStore) SendMessage(conn *websocket.Conn, message string) error {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	if _, ok := cs.connections[conn]; !ok {
		return nil
	}

	err := conn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		log.Printf("Error sending message: %v", err)
		cs.RemoveConnection(conn)
		return err
	}
	return nil
}

type GameSessionConnectionStore struct {
	connections map[string]*websocket.Conn
	mu          sync.Mutex
	cs          *ConnectionStore
}

func newGameSessionConnectionStore(cs *ConnectionStore) *GameSessionConnectionStore {
	return &GameSessionConnectionStore{
		connections: make(map[string]*websocket.Conn),
		cs:          cs,
	}
}

func (ucs *GameSessionConnectionStore) AddGameSessionConnection(gameSessionId string, conn *websocket.Conn) {
	ucs.mu.Lock()
	defer ucs.mu.Unlock()
	ucs.connections[gameSessionId] = conn
	ucs.cs.AddConnection(conn)
}

func (ucs *GameSessionConnectionStore) RemoveGameSessionConnection(gameSessionId string) {
	ucs.mu.Lock()
	defer ucs.mu.Unlock()
	conn, ok := ucs.connections[gameSessionId]
	if ok {
		delete(ucs.connections, gameSessionId)
		ucs.cs.RemoveConnection(conn)
	}
}

func (ucs *GameSessionConnectionStore) SendMessageToGameSession(gameSessionId string, message string) error {
	ucs.mu.Lock()
	defer ucs.mu.Unlock()
	conn, ok := ucs.connections[gameSessionId]
	if ok {
		return ucs.cs.SendMessage(conn, message)
	}
	return nil
}

func (ucs *GameSessionConnectionStore) Echo(w http.ResponseWriter, r *http.Request, gameSessionId string, db *database.InMemoryGameSessionDB) {
	if !checkGameSessionExists(gameSessionId, db) {
		log.Println("Game Session does not exist:", gameSessionId)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}

	ucs.AddGameSessionConnection(gameSessionId, conn)
	defer func() {
		ucs.RemoveGameSessionConnection(gameSessionId)
		err := conn.Close()
		if err != nil {
			return
		}
	}()

	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			break
		}
		log.Printf("Received: %s", message)

		if err := conn.WriteMessage(messageType, []byte("Server does not process messages from the client")); err != nil {
			log.Println("Write error:", err)
			break
		}
	}
}

func checkGameSessionExists(gameSessionId string, db *database.InMemoryGameSessionDB) bool {
	_, err := db.GetSession(gameSessionId)
	return err == nil
}
