package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func main() {
	go startGraphics()
	rand.Seed(time.Now().UTC().UnixNano())
	world := NewWorld()
	world.GenerateRandomPrey(1000)
	var wg sync.WaitGroup
	for iteration := 0; iteration < 10000; iteration++ {
		for i, _ := range world.prey {
			wg.Add(1)
			go runAgentWG(world.prey[i], world, &wg)
		}
		wg.Wait()
		for i, _ := range world.prey {
			wg.Add(1)
			go stepAgentWG(world.prey[i], &wg)
		}
		wg.Wait()
		for i, _ := range world.prey {
			if i == 0 {
				gfx <- Graphic{
					location:  *world.prey[i].GetLocation(),
					direction: *world.prey[i].GetDirection(),
					cls:       true,
				}
			} else {
				gfx <- Graphic{
					location:  *world.prey[i].GetLocation(),
					direction: *world.prey[i].GetDirection(),
					cls:       false,
				}
			}
		}
		fmt.Println(iteration)
		time.Sleep(10 * time.Millisecond)
	}
}
