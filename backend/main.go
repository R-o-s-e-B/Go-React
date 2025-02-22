package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type Player struct {
	ID string `json:"id"`
	X  float64 `json:"x"`
	Y  float64 `json:"y"`
}


var upgrader = websocket.Upgrader{
	//ReadBufferSize: 1024,
	//WriteBufferSize: 1024,

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
var players = make(map[string]Player)   

func handleMessages(conn *websocket.Conn) {
	defer func() {
			playerID, exists := clients[conn]
			if exists {
					delete(players, playerID) // Remove player from the list
					delete(clients, conn)     // Remove connection from clients
					broadcast()               // Notify others
			}
			conn.Close()
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

			// Update player movement
			players[player.ID] = player
			broadcast()
	}
}


func broadcast() {
	playerList, err := json.Marshal(players) // Convert players map to JSON
	//log.Println("This is the player list", playerList)
	log.Println("Broadcasting players:", string(playerList)) // Debugging
	if err != nil {
		log.Println("Error marshalling players:", err)
		return
	}

	for conn := range clients {
		if err := conn.WriteMessage(websocket.TextMessage, playerList); err != nil {
			log.Println("Error sending message:", err)
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

	// Store the connection, but don't create a player yet
	clients[ws] = "" // No player assigned yet

	go handleMessages(ws)
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