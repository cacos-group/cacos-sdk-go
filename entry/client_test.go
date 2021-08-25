package entry

import (
	"context"
	"fmt"
	clientV3 "go.etcd.io/etcd/client/v3"
	"testing"
	"time"
)

func TestSource_Load(t *testing.T) {
	p, err := NewClient(&Config{
		Endpoints: []string{"127.0.0.1:2379"},
		Path:      "/",
		Prefix:    true,
		Username:  "root",
		Password:  "root",
	})
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Println(p)
}

func TestNew(t *testing.T) {
	c, err := clientV3.New(clientV3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: 5 * time.Second,
	})

	if err != nil {
		t.Error(err)
		return
	}

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err = c.KV.Put(ctx, "/new/sourcer/project/key7", "val2")
	if err != nil {
		t.Error(err)
		return
	}
}

func TestSource_Delete(t *testing.T) {
	c, err := clientV3.New(clientV3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: 5 * time.Second,
		Username:    "user",
		Password:    "password",
	})

	if err != nil {
		t.Error(err)
		return
	}

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err = c.KV.Delete(ctx, "/sourcer/project/key7")
	if err != nil {
		t.Error(err)
		return
	}
}
