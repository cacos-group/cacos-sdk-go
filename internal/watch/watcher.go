package watch

import (
	"context"
	"github.com/cacos-group/cacos-sdk-go/internal/model"
	source "github.com/cacos-group/cacos-sdk-go/internal/source"
	"github.com/pkg/errors"
	clientV3 "go.etcd.io/etcd/client/v3"
)

type Watcher interface {
	Next() ([]*model.Event, error)
	Stop()
}

type watcher struct {
	source    source.Sourcer
	ch        clientV3.WatchChan
	closeChan chan struct{}
}

func New(s source.Sourcer, path string, opts ...Option) (*watcher, error) {
	options := &options{
		path:   path,
		prefix: true,
	}

	for _, opt := range opts {
		opt(options)
	}

	if options.path == "" {
		return nil, errors.New("path invalid")
	}

	w := &watcher{
		source:    s,
		closeChan: make(chan struct{}),
	}

	var opOpts []clientV3.OpOption
	if options.prefix {
		opOpts = append(opOpts, clientV3.WithPrefix())
	}

	w.ch = s.Client().Watch(context.Background(), path, opOpts...)

	return w, nil
}

func (s *watcher) Next() ([]*model.Event, error) {
	select {
	case rsp := <-s.ch:
		list := make([]*model.Event, 0, len(rsp.Events))
		for _, event := range rsp.Events {
			newEvent := &model.Event{
				Type: event.Type,
				Kv:   event.Kv,
			}

			list = append(list, newEvent)
		}

		return list, nil
	case <-s.closeChan:
		return nil, nil
	}
}
func (s *watcher) Stop() {
	close(s.closeChan)
}
