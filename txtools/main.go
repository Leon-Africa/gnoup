package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	graphql "github.com/hasura/go-graphql-client"
)

var (
	txIndex = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "transaction_index",
		Help: "Index of the transaction",
	})
	txBlockHeight = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "transaction_block_height",
		Help: "Block height of the transaction",
	})
	txGasWanted = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "transaction_gas_wanted",
		Help: "Gas wanted of the transaction",
	})
	txGasUsed = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "transaction_gas_used",
		Help: "Gas used of the transaction",
	})
	txHash = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "transaction_hash",
		Help: "Hash of the transaction",
	}, []string{"hash"})
)

func init() {
	prometheus.MustRegister(txIndex)
	prometheus.MustRegister(txBlockHeight)
	prometheus.MustRegister(txGasWanted)
	prometheus.MustRegister(txGasUsed)
	prometheus.MustRegister(txHash)
}

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

	var sub struct {
		Transactions struct {
			Index       int `graphql:"index"`
			BlockHeight int `graphql:"block_height"`
			GasWanted   int `graphql:"gas_wanted"`
			GasUsed     int `graphql:"gas_used"`
			Hash        string
		} `graphql:"transactions(filter: {from_block_height: 1})"`
	}

	_, err := client.Subscribe(sub, nil, func(data []byte, err error) error {
		if err != nil {
			log.Println(err)
			return nil
		}

		if data == nil {
			return nil
		}

		var response struct {
			Transactions struct {
				Index       int    `json:"index"`
				BlockHeight int    `json:"block_height"`
				GasWanted   int    `json:"gas_wanted"`
				GasUsed     int    `json:"gas_used"`
				Hash        string `json:"hash"`
			} `json:"transactions"`
		}
		if err := json.Unmarshal(data, &response); err != nil {
			log.Println("Failed to unmarshal data:", err)
			return nil
		}

		// Update Prometheus metrics with transaction data
		txIndex.Set(float64(response.Transactions.Index))
		txBlockHeight.Set(float64(response.Transactions.BlockHeight))
		txGasWanted.Set(float64(response.Transactions.GasWanted))
		txGasUsed.Set(float64(response.Transactions.GasUsed))
		txHash.WithLabelValues(response.Transactions.Hash).Set(1)

		log.Printf("Transaction data: %+v\n", response.Transactions)

		return nil
	})

	if err != nil {
		panic(err)
	}

	// Start HTTP server to expose Prometheus metrics
	http.Handle("/metrics", promhttp.Handler())
	go func() {
		if err := http.ListenAndServe(":7770", nil); err != nil {
			log.Fatal("Failed to start HTTP server:", err)
		}
	}()

	return client.Run()
}

func main() {
	if err := startSubscription(); err != nil {
		log.Fatal(err)
	}
}
