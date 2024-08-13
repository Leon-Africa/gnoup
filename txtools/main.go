package main

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
)

type BlockHeader struct {
	Number     string `json:"number"`
	ParentHash string `json:"parentHash"`
	// Add other fields as needed
}

func main() {
	// Define WebSocket endpoint
	wsURL := "ws://0.0.0.0:8546/ws"

	// Establish WebSocket connection
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		log.Fatal("WebSocket connection failed:", err)
	}
	defer conn.Close()

	// Send subscription request
	subscribeRequest := map[string]interface{}{
		"id":      1,
		"jsonrpc": "2.0",
		"method":  "subscribe",
		"params":  []string{"newHeads"},
	}
	if err := conn.WriteJSON(subscribeRequest); err != nil {
		log.Fatal("Error sending subscription request:", err)
	}

	// Log success message
	log.Println("Subscription request sent successfully")

	// Read and process messages from the WebSocket connection
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Fatal("Error reading WebSocket message:", err)
		}

		var response map[string]interface{}
		if err := json.Unmarshal(msg, &response); err != nil {
			log.Fatal("Error decoding JSON data:", err)
		}

		// Check if it's a subscription message
		params, ok := response["params"].(map[string]interface{})
		if !ok {
			log.Println("Received message:", string(msg))
			continue
		}

		blockHeader, ok := params["result"].(map[string]interface{})
		if !ok {
			log.Println("Received message:", string(msg))
			continue
		}

		// Process block header
		log.Println("Received new block header:")
		log.Println("Block Number:", blockHeader["number"])
		log.Println("Parent Hash:", blockHeader["parentHash"])

	}
}
