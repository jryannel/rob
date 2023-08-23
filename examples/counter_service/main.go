package main

import (
	"encoding/json"
	"log"

	"github.com/jryannel/rob"
	"github.com/nats-io/nats.go"
)

type ICounterService interface {
	DoIncrement() error
	DoDecrement() error
	DoReset() error
	SetCount(v int) error
}

type CounterService struct {
	count int
	so    rob.IRemoteService
}

func NewCounterService(nc *nats.Conn, id string) *CounterService {
	cs := &CounterService{
		so: rob.NewRemoteService(nc, id),
	}
	err := cs.init()
	if err != nil {
		log.Fatal(err)
	}
	return cs
}

func (cs *CounterService) init() error {
	log.Println("init service")
	err := cs.so.ProvideValue("count", func() ([]byte, error) {
		data, err := json.Marshal(cs.count)
		if err != nil {
			return nil, err
		}
		return data, nil
	})
	if err != nil {
		return err
	}
	err = cs.so.RegisterMethod("increment", func(data []byte) ([]byte, error) {
		cs.DoIncrement()
		return nil, nil
	})
	if err != nil {
		return err
	}
	err = cs.so.RegisterMethod("decrement", func(data []byte) ([]byte, error) {
		cs.DoDecrement()
		return nil, nil
	})
	if err != nil {
		return err
	}
	err = cs.so.RegisterMethod("reset", func(data []byte) ([]byte, error) {
		cs.DoReset()
		return nil, nil
	})
	if err != nil {
		return err
	}
	return nil
}

var _ ICounterService = (*CounterService)(nil)

func (cs *CounterService) DoIncrement() error {
	cs.count++
	cs.SetCount(cs.count)
	return nil
}

func (cs *CounterService) DoDecrement() error {
	cs.count--
	cs.SetCount(cs.count)
	return nil
}

func (cs *CounterService) DoReset() error {
	cs.count = 0
	cs.SetCount(cs.count)
	return nil
}

func (cs *CounterService) SetCount(v int) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	cs.so.PublishValue("count", data)
	return nil
}

func main() {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		panic(err)
	}
	id := "counter"
	NewCounterService(nc, id)
	select {}
}
