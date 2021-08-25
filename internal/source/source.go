package source

import (
	"context"
	"errors"
	"github.com/cacos-group/cacos-sdk-go/internal/model"
	"go.etcd.io/etcd/api/v3/mvccpb"
	"time"

	clientV3 "go.etcd.io/etcd/client/v3"
)

type Sourcer interface {
	Load() ([]*model.Event, error)
	Client() *clientV3.Client
}

type source struct {
	client  *clientV3.Client
	options *options
}

func New(client *clientV3.Client, path string, opts ...Option) (*source, error) {
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

	return &source{
		client:  client,
		options: options,
	}, nil
}

// KeyValue is config key value.
type KeyValue struct {
	Key   string
	Value []byte
}

// Load return the config values
func (s *source) Load() ([]*model.Event, error) {
	var opts []clientV3.OpOption
	if s.options.prefix {
		opts = append(opts, clientV3.WithPrefix())
	}

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	rsp, err := s.client.Get(ctx, s.options.path, opts...)
	if err != nil {
		return nil, err
	}

	var kvs []*model.Event
	for _, item := range rsp.Kvs {
		newEvent := &model.Event{
			Type: mvccpb.PUT,
			Kv:   item,
		}
		kvs = append(kvs, newEvent)
	}
	return kvs, nil
}

func (s *source) Client() *clientV3.Client {
	return s.client
}
