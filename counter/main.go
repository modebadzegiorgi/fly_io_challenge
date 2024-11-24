package main

import (
	"context"
	"encoding/json"
	"log"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

const key = "Total"

func main() {

	n := maelstrom.NewNode()
	kv := maelstrom.NewSeqKV(n)

	n.Handle("add", func(msg maelstrom.Message) error {
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

		for _, id := range n.NodeIDs() {
			if id == n.ID() {
				continue
			}

			n.Send(id, map[string]any{
				"type":  "broadcast_add",
				"delta": delta,
			})

		}

		body["type"] = "add_ok"

		delete(body, "delta")

		return n.Reply(msg, body)
	})

	n.Handle("broadcast_add", func(msg maelstrom.Message) error {
		ctx := context.Background()
		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		delta := int(body["delta"].(float64))

		// Update the local KV store
		currentTotal, err := kv.ReadInt(ctx, key)
		if err != nil {
			currentTotal = 0
		}

		kv.CompareAndSwap(ctx, key, currentTotal, currentTotal+delta, true)

		return nil // No reply needed for broadcasts
	})

	n.Handle("read", func(msg maelstrom.Message) error {
		ctx := context.Background()
		var body map[string]any

		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		val, _ := kv.ReadInt(ctx, key)
		body["type"] = "read_ok"
		body["value"] = val

		return n.Reply(msg, body)
	})

	if err := n.Run(); err != nil {
		log.Fatal(err)
	}
}
