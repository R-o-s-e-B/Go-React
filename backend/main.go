package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Player struct {
	ID string `json:"id"`
	X  float64 `json:"x"`
	Y  float64 `json:"y"`
}


var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {return true},
}


func reader(conn *websocket.Conn) {
	for {
		messageType, p, err := conn.ReadMessage()
		if err!= nil {
			log.Println(err)
			return
		}

		fmt.Println(string(p))

		if err := conn.WriteMessage(messageType, p); err != nil {
			log.Println(err)
			return
		}
	}
}

var clients = make(map[*websocket.Conn]string) // Connected clients
var mu sync.Mutex
var players = make(map[string]Player)   

func handleMessages(conn *websocket.Conn) {
	defer func() {
			mu.Lock()
			playerID, exists := clients[conn]
			log.Println("Before removing: playerID =", playerID, "exists =", exists)
			
			if exists {
					log.Println("Player disconnected:", playerID)
					delete(players, playerID)
					delete(clients, conn)
					log.Println("Players after removal:", players)
					go broadcast() 
			}
			mu.Unlock()

			// Send WebSocket close message
			deadline := time.Now().Add(5 * time.Second)
			err := conn.WriteControl(websocket.CloseMessage, 
					websocket.FormatCloseMessage(websocket.CloseNormalClosure, "Goodbye"), 
					deadline)
			if err != nil {
					log.Println("Error sending close message:", err)
			}

			conn.Close() // Finally close the connection
	}()

	for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
					log.Println("Read error:", err)
					break
			}

			var player Player
			if err := json.Unmarshal(msg, &player); err != nil {
					log.Println("JSON Unmarshal error:", err, "Received data:", string(msg))
					continue
			}

			mu.Lock()
			if _, exists := clients[conn]; !exists || clients[conn] == "" {
					clients[conn] = player.ID // Assign player ID when first received
			}
			players[player.ID] = player
			mu.Unlock()

			broadcast()
	}
}




func broadcast() {
	mu.Lock()
	defer mu.Unlock()

	playerList, err := json.Marshal(players) // Convert to JSON
	if err != nil {
			log.Println("Error marshalling players:", err)
			return
	}


	for conn := range clients {
			err := conn.WriteMessage(websocket.TextMessage, playerList)
			if err != nil {
					log.Println("Broadcast error:", err)
					conn.Close()
					delete(clients, conn)
			}
	}
}






func serveWs(w http.ResponseWriter, r *http.Request) {
	fmt.Println("New connection from:", r.Host)

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
			log.Println("WebSocket upgrade error:", err)
			return
	}

	go handleMessages(ws) // Move player assignment inside handleMessages
}


func setupRoutes(){
	http.HandleFunc("/", func(w http.ResponseWriter, r*http.Request){
		fmt.Println(w, "Simple server")
	})

	http.HandleFunc("/ws", serveWs)
}

func main(){
	fmt.Println("Chat app v0.01")
	setupRoutes()
	http.ListenAndServe(":8080", nil)
}