package main

import (
	"context"
	"encoding/json"
	"log"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func main() {

	n := maelstrom.NewNode()
	kv := maelstrom.NewSeqKV(n)

	n.Handle("add", func(msg maelstrom.Message) error {
		key := n.ID()

		ctx := context.Background()
		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		delta := int(body["delta"].(float64))

		currentTotal, err := kv.ReadInt(ctx, key)
		if err != nil {
			currentTotal = 0
		}
		kv.CompareAndSwap(ctx, key, currentTotal, currentTotal+delta, true)
		// Reply to the client
		body["type"] = "add_ok"
		delete(body, "delta")

		return n.Reply(msg, body)
	})

	n.Handle("read", func(msg maelstrom.Message) error {
		ctx := context.Background()
		var body map[string]any

		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		total := 0
		for _, id := range n.NodeIDs() {
			val, _ := kv.ReadInt(ctx, id)
			total += val
		}

		body["type"] = "read_ok"
		body["value"] = total

		return n.Reply(msg, body)
	})

	if err := n.Run(); err != nil {
		log.Fatal(err)
	}
}
