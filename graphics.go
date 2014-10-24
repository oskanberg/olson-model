package main

// import (
// 	"fmt"
// 	"time"

// 	"github.com/oskanberg/gowebsockets"
// )

// func main() {
// 	gowebsockets.OnConnection(func(c *gowebsockets.Connection) {
// 		gowebsockets.SetSize(500, 500)
// 		params := gowebsockets.ArcParameters{X: 100, Y: 100, R: 50}
// 		for {
// 			fmt.Println("drawing")
// 			gowebsockets.Arc(params)
// 			time.Sleep(3 * time.Second)
// 		}
// 	})
// 	gowebsockets.ServeForever()
// }
