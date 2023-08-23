package main

import (
	"log"
	"runtime"

	"github.com/nats-io/nats.go"
)

func main() {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal(err)
	}
	id := "counter"
	cc := NewCounterClient(nc, id)
	NewCounterService(nc, id)
	go func() {
		for i := 0; i < 10000; i++ {
			err = cc.Increment()
			if err != nil {
				log.Println("error invoking method:", err)
			}
		}
	}()
	go func() {
		for i := 0; i < 10000; i++ {
			err = cc.Decrement()
			if err != nil {
				log.Println("error invoking method:", err)
			}
		}
	}()

	cc.OnValueChange(func(v int) {
		log.Println("value changed:", v)
	})

	runtime.Goexit()
}
