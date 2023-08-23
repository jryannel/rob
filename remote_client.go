package rob

import (
	"log"

	"github.com/nats-io/nats.go"
)

type IRemoteClient interface {
	OnValue(name string, fn func([]byte)) (func() error, error)
	InvokeMethod(name string, data []byte) ([]byte, error)
	OnSignal(name string, fn func([]byte)) (func() error, error)
}

type RemoteClient struct {
	nc *nats.Conn
	id string
}

func NewRemoteClient(nc *nats.Conn, id string) *RemoteClient {
	return &RemoteClient{
		nc: nc,
		id: id,
	}
}

var _ IRemoteClient = (*RemoteClient)(nil)

func (c *RemoteClient) OnValue(name string, fn func([]byte)) (func() error, error) {
	log.Println("OnValue", name)
	subj := c.id + "." + name
	sub, err := c.nc.Subscribe(subj, func(m *nats.Msg) {
		fn(m.Data)
	})
	if err != nil {
		return nil, err
	}
	return sub.Unsubscribe, nil
}

func (c *RemoteClient) InvokeMethod(name string, data []byte) ([]byte, error) {
	log.Println("InvokeMethod", name)
	subj := c.id + "." + name
	msg, err := c.nc.Request(subj, data, nats.DefaultTimeout)
	if err != nil {
		return nil, err
	}
	return msg.Data, nil
}

func (c *RemoteClient) OnSignal(name string, fn func([]byte)) (func() error, error) {
	log.Println("OnSignal", name)
	subj := c.id + "." + name
	sub, err := c.nc.Subscribe(subj, func(m *nats.Msg) {
		fn(m.Data)
	})
	if err != nil {
		return nil, err
	}
	return sub.Unsubscribe, nil
}
