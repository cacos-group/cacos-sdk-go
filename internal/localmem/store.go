package localmem

import (
	"fmt"
	"github.com/cacos-group/cacos-sdk-go/internal/model"
	"go.etcd.io/etcd/api/v3/mvccpb"
	"sync"
)

type Storer interface {
	HandleEvent(event *model.Event)
	Get(key string) (value interface{}, ok bool)
}

type store struct {
	m sync.Map
}

func New() Storer {
	return &store{
		m: sync.Map{},
	}
}

func (s *store) HandleEvent(event *model.Event) {
	switch event.Type {
	case mvccpb.PUT:
		s.m.Store(string(event.Kv.Key), event.Kv.Value)
	case mvccpb.DELETE:
		s.m.Delete(string(event.Kv.Key))
	}

	fmt.Println(event)
}

func (s *store) Get(key string) (value interface{}, ok bool) {
	return s.m.Load(key)
}
