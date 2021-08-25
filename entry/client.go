package entry

import (
	"github.com/cacos-group/cacos-sdk-go/internal/localmem"
	"github.com/cacos-group/cacos-sdk-go/internal/source"
	"github.com/cacos-group/cacos-sdk-go/internal/watch"
	clientV3 "go.etcd.io/etcd/client/v3"
	"time"
)

type Client interface {
	Get(key string) (value interface{}, ok bool)
}

type client struct {
	clientV3 *clientV3.Client
	config   *Config

	sourcer source.Sourcer
	storer  localmem.Storer

	watcher    watch.Watcher
	watcherMap map[string]watch.Watcher
}

func NewClient(conf *Config) (Client, error) {
	c := &client{
		config: conf,
		storer: localmem.New(),
	}

	cc, err := clientV3.New(clientV3.Config{
		Endpoints:   conf.Endpoints,
		DialTimeout: 5 * time.Second,
		Username:    conf.Username,
		Password:    conf.Password,
	})
	if err != nil {
		return nil, err
	}

	c.clientV3 = cc

	err = c.init()
	if err != nil {
		return nil, err
	}
	err = c.daemon()
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (c *client) init() error {
	var opts []source.Option
	opts = append(opts, source.WithPrefix(c.config.Prefix))

	sourcer, err := source.New(c.clientV3, c.config.Path, opts...)
	if err != nil {
		return err
	}

	c.sourcer = sourcer

	events, err := c.sourcer.Load()
	if err != nil {
		return err
	}

	for _, event := range events {
		c.storer.HandleEvent(event)
	}

	return nil
}

func (c *client) daemon() error {
	var watcherOpts []watch.Option
	watcher, err := watch.New(c.sourcer, c.config.Path, watcherOpts...)
	if err != nil {
		return err
	}
	for true {
		events, err := watcher.Next()
		if err != nil {
			continue
		}
		for _, event := range events {
			c.storer.HandleEvent(event)
		}
	}

	return nil
}

func (c *client) Get(key string) (value interface{}, ok bool) {
	return c.storer.Get(key)
}
