package main

import (
	"encoding/json"
	"log"

	"github.com/jryannel/rob"

	"github.com/nats-io/nats.go"
)

type ICounterClient interface {
	Increment() error
	Decrement() error
	Reset() error
	Count() (int, error)
}

type CounterClient struct {
	count int
	co    rob.IRemoteClient
}

func NewCounterClient(nc *nats.Conn, id string) *CounterClient {
	cc := &CounterClient{
		co: rob.NewRemoteClient(nc, id),
	}
	err := cc.init()
	if err != nil {
		log.Fatal(err)
	}
	return cc
}

func (cc *CounterClient) init() error {
	cc.co.OnValue("count", func(data []byte) {
		var v int
		err := json.Unmarshal(data, &v)
		if err != nil {
			log.Println("error un-marshalling value:", err)
			return
		}
		cc.count = v
	})
	return nil
}

var _ ICounterClient = (*CounterClient)(nil)

func (cc *CounterClient) Increment() error {
	_, err := cc.co.InvokeMethod("increment", nil)
	return err
}

func (cc *CounterClient) Decrement() error {
	_, err := cc.co.InvokeMethod("decrement", nil)
	return err
}

func (cc *CounterClient) Reset() error {
	_, err := cc.co.InvokeMethod("reset", nil)
	return err
}

func (cc *CounterClient) Count() (int, error) {
	return cc.count, nil
}

func (co *CounterClient) OnValueChange(fn func(int)) error {
	co.co.OnValue("count", func(data []byte) {
		var v int
		err := json.Unmarshal(data, &v)
		if err != nil {
			log.Println("error un-marshalling value:", err)
			return
		}
		fn(v)
	})
	return nil
}

func main() {
	log.Println("starting counter client")
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal(err)
	}
	cc := NewCounterClient(nc, "counter")
	log.Println("waiting for value changes")
	cc.OnValueChange(func(v int) {
		log.Println("value changed:", v)
	})

	log.Println("increment 10000 times")

	go func() {
		for i := 0; i < 10000; i++ {
			err = cc.Increment()
			if err != nil {
				log.Println("error invoking method:", err)
			}
		}
	}()
	go func() {
		log.Println("decrement 10000 times")
		for i := 0; i < 10000; i++ {
			err = cc.Decrement()
			if err != nil {
				log.Println("error invoking method:", err)
			}
		}
	}()
	select {}
}
