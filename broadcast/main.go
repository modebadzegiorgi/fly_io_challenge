package main

import (
	"encoding/json"
	"log"
	"log/slog"
	"sync"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

type localData struct {
	sync.Mutex
	store map[float64]bool
}

var neighbors []string

func main() {

	m := localData{
		store: make(map[float64]bool),
	}
	n := maelstrom.NewNode()

	n.Handle("broadcast", func(msg maelstrom.Message) error {

		var body map[string]any

		parseMessageBody(msg, &body)

		m.Lock()
		message := body["message"]
		m.store[message.(float64)] = true

		for _, neighbor := range neighbors {
			if neighbor != n.ID() {
				n.Send(
					neighbor,
					map[string]any{
						"type":    "broadcast",
						"message": message,
					},
				)
			}
		}
		m.Unlock()
		delete(body, "message")

		body["type"] = "broadcast_ok"
		return n.Reply(msg, body)
	})

	n.Handle("broadcast_ok", func(msg maelstrom.Message) error {
		// Log the acknowledgment (optional)
		slog.Info("Received broadcast_ok", "from", msg.Src)
		return nil
	})

	n.Handle("read", func(msg maelstrom.Message) error {

		var body map[string]any

		parseMessageBody(msg, &body)

		m.Lock()
		tempMessage := []float64{}
		for key := range m.store {
			tempMessage = append(tempMessage, key)
		}
		body["messages"] = tempMessage
		m.Unlock()
		body["type"] = "read_ok"

		return n.Reply(msg, body)
	})

	n.Handle("topology", func(msg maelstrom.Message) error {

		var body map[string]any

		parseMessageBody(msg, &body)

		// Extract the topology map
		topology := body["topology"].(map[string]any)

		// // Get the neighbors for this node
		m.Lock()
		neighbors = []string{}
		for _, neighbor := range topology[n.ID()].([]any) {
			neighbors = append(neighbors, neighbor.(string))
		}
		m.Unlock()

		delete(body, "topology")
		body["type"] = "topology_ok"
		return n.Reply(msg, body)
	})

	if err := n.Run(); err != nil {
		log.Fatal(err)
	}
}

func parseMessageBody(msg maelstrom.Message, body *map[string]any) {

	if err := json.Unmarshal(msg.Body, body); err != nil {

		slog.Error("Could not Parse body of the message", "err", err)
	}

}
