package main

import (
	"time"

	"github.com/oskanberg/go-ws-draw"
)

type Graphic struct {
	location  Vector2D
	direction Vector2D
	cls       bool
}

var gfx = make(chan Graphic, 5)

func startGraphics() {
	gowebsockets.OnConnection(func(c *gowebsockets.Connection) {
		gowebsockets.SetSize(512, 512)
	})

	go gowebsockets.ServeForever(8080)

	for msg := range gfx {
		if msg.cls {
			gowebsockets.Clear()
		}

		params := gowebsockets.ArcParameters{
			X: msg.location.x,
			Y: msg.location.y,
			R: 2,
		}
		gowebsockets.Arc(params)
		destination := msg.location.Add(msg.direction.Multiplied(5))
		lineParams := gowebsockets.LineParameters{
			Start: gowebsockets.Vector{
				X: msg.location.x,
				Y: msg.location.y,
			},
			End: gowebsockets.Vector{
				X: destination.x,
				Y: destination.y,
			},
		}
		gowebsockets.Line(lineParams)
		time.Sleep(0 * time.Millisecond)
	}
}
