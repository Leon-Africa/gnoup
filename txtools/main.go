// package main

// import (
// 	"encoding/json"
// 	"log"

// 	"github.com/gorilla/websocket"
// )

// type BlockHeader struct {
// 	Number     string `json:"number"`
// 	ParentHash string `json:"parentHash"`
// 	// Add other fields as needed
// }

// func main() {
// 	// Define WebSocket endpoint
// 	wsURL := "ws://0.0.0.0:8546/ws"

// 	// Establish WebSocket connection
// 	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
// 	if err != nil {
// 		log.Fatal("WebSocket connection failed:", err)
// 	}
// 	defer conn.Close()

// 	// Send subscription request
// 	subscribeRequest := map[string]interface{}{
// 		"id":      1,
// 		"jsonrpc": "2.0",
// 		"method":  "subscribe",
// 		"params":  []string{"newHeads"},
// 	}
// 	if err := conn.WriteJSON(subscribeRequest); err != nil {
// 		log.Fatal("Error sending subscription request:", err)
// 	}

// 	// Log success message
// 	log.Println("Subscription request sent successfully")

// 	// Read and process messages from the WebSocket connection
// 	for {
// 		_, msg, err := conn.ReadMessage()
// 		if err != nil {
// 			log.Fatal("Error reading WebSocket message:", err)
// 		}

// 		var response map[string]interface{}
// 		if err := json.Unmarshal(msg, &response); err != nil {
// 			log.Fatal("Error decoding JSON data:", err)
// 		}

// 		// Check if it's a subscription message
// 		params, ok := response["params"].(map[string]interface{})
// 		if !ok {
// 			log.Println("Received message:", string(msg))
// 			continue
// 		}

// 		blockHeader, ok := params["result"].(map[string]interface{})
// 		if !ok {
// 			log.Println("Received message:", string(msg))
// 			continue
// 		}

// 		// Process block header
// 		log.Println("Received new block header:")
// 		log.Println("Block Number:", blockHeader["number"])
// 		log.Println("Parent Hash:", blockHeader["parentHash"])

// 	}
// }

package main

import (
	"fmt"
	"log"

	graphql "github.com/hasura/go-graphql-client"
)

func getServerEndpoint() string {
	return fmt.Sprintf("http://0.0.0.0:%d/graphql/query", 8546)
}

func startSubscription() error {

	client := graphql.NewSubscriptionClient(getServerEndpoint()).
		WithConnectionParams(map[string]interface{}{
			"headers": map[string]string{
				"foo": "bar",
			},
		}).WithLog(log.Println).
		WithoutLogTypes(graphql.GQLData, graphql.GQLConnectionKeepAlive).
		OnError(func(sc *graphql.SubscriptionClient, err error) error {
			log.Print("err", err)
			return err
		})

	defer client.Close()

	// var sub struct {
	// 	Blocks struct {
	// 		Height int `graphql:"height"`
	// 	} `graphql:"blocks(filter: {from_height: 1})"`
	// }

	var sub struct {
		Transactions struct {
			Index       int    `graphql:"index"`
			Hash        string `graphql:"hash"`
			BlockHeight int    `graphql:"block_height"`
			GasWanted   int    `graphql:"gas_wanted"`
			GasUsed     int    `graphql:"gas_used"`
			ContentRaw  string `graphql:"content_raw"`
			Messages    []struct {
				TypeURL string `graphql:"typeUrl"`
				Route   string `graphql:"route"`
				Value   struct {
					TypeName string `graphql:"__typename"`
				} `graphql:"value"`
			} `graphql:"messages"`
			Memo string `graphql:"memo"`
		} `graphql:"transactions(filter: {from_block_height: 1})"`
	}

	//subId
	_, err := client.Subscribe(sub, nil, func(data []byte, err error) error {

		if err != nil {
			log.Println(err)
			return nil
		}

		if data == nil {
			return nil
		}
		log.Println(string(data))
		return nil
	})

	if err != nil {
		panic(err)
	}

	// automatically unsubscribe after 10 seconds
	// go func() {
	// 	time.Sleep(10000 * time.Second)
	// 	_ = client.Unsubscribe(subId)
	// }()

	return client.Run()
}

func main() {
	// Call the startSubscription function
	if err := startSubscription(); err != nil {
		log.Fatal(err)
	}
}
