package model

import (
	"go.etcd.io/etcd/api/v3/mvccpb"
)

type Event struct {
	// type is the kind of event. If type is a PUT, it indicates
	// new data has been stored to the key. If type is a DELETE,
	// it indicates the key was deleted.
	Type mvccpb.Event_EventType
	// kv holds the KeyValue for the event.
	// A PUT event contains current kv pair.
	// A PUT event with kv.Version=1 indicates the creation of a key.
	// A DELETE/EXPIRE event contains the deleted key with
	// its modification revision set to the revision of deletion.
	Kv *mvccpb.KeyValue
}
