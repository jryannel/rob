package rob

import (
	"log"

	"github.com/nats-io/nats.go"
)

type IRemoteService interface {
	ProvideValue(name string, fn func() ([]byte, error)) error
	RegisterMethod(name string, fn func(data []byte) ([]byte, error)) error
	EmitSignal(name string, data []byte) error
	PublishValue(name string, data []byte) error
}

type RemoteService struct {
	nc *nats.Conn
	id string
}

func NewRemoteService(nc *nats.Conn, id string) *RemoteService {
	return &RemoteService{
		nc: nc,
		id: id,
	}
}

var _ IRemoteService = (*RemoteService)(nil)

func (s *RemoteService) ProvideValue(name string, fn func() ([]byte, error)) error {
	log.Println("ProvideValue", name)
	subj := s.id + "." + name
	s.nc.Subscribe(subj, func(m *nats.Msg) {
		data, err := fn()
		if err != nil {
			log.Println("error invoking function")
		}
		s.nc.Publish(m.Reply, data)
	})
	return nil
}

func (s *RemoteService) RegisterMethod(name string, fn func(data []byte) ([]byte, error)) error {
	log.Println("RegisterMethod", name)
	subj := s.id + "." + name
	s.nc.Subscribe(subj, func(m *nats.Msg) {
		data, err := fn(m.Data)
		if err != nil {
			log.Println("error invoking function")
		}
		s.nc.Publish(m.Reply, data)
	})
	return nil
}

func (s *RemoteService) EmitSignal(name string, data []byte) error {
	log.Println("EmitSignal", name)
	subj := s.id + "." + name
	return s.nc.Publish(subj, data)
}

func (s *RemoteService) PublishValue(name string, data []byte) error {
	log.Println("PublishValue", name)
	subj := s.id + "." + name
	return s.nc.Publish(subj, data)
}
