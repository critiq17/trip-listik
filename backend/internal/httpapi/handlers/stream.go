package handlers

import (
	"bufio"
	"encoding/json"
	"fmt"
	"time"

	"github.com/critiq17/tripListik/internal/realtime"
	"github.com/gofiber/fiber/v2"
)

type StreamHandler struct {
	Hub *realtime.Hub
}

func (h *StreamHandler) TripStream(c *fiber.Ctx) error {
	tripID := c.Params("id")
	if tripID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "invalid trip id")
	}

	c.Set("Content-Type", "text/event-stream")
	c.Set("Cache-Control", "no-cache")
	c.Set("Connection", "keep-alive")

	key := fmt.Sprintf("trip:%s", tripID)
	ch := h.Hub.Subscribe(key)
	defer h.Hub.Unsubscribe(key, ch)

	c.Context().SetBodyStreamWriter(func(w *bufio.Writer) {
		ping := time.NewTicker(25 * time.Second)
		defer ping.Stop()

		for {
			select {
			case <-c.Context().Done():
				return
			case ev, ok := <-ch:
				if !ok {
					return
				}
				payload, _ := json.Marshal(ev)
				fmt.Fprintf(w, "event: %s\n", ev.Type)
				fmt.Fprintf(w, "data: %s\n\n", payload)
				w.Flush()
			case <-ping.C:
				fmt.Fprintf(w, "event: ping\ndata: {}\n\n")
				w.Flush()
			}
		}
	})

	return nil
}
