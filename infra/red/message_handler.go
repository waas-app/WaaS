package red

import (
	"context"
	"reflect"
	"runtime"
	"sync"
)

type MessageHandlerFunc func(context.Context, *Message) error

type handlerStore struct {
	store map[string][]MessageHandlerFunc
	sync.Map
}

func newHandlerStore() *handlerStore {
	h := handlerStore{
		make(map[string][]MessageHandlerFunc),
		sync.Map{},
	}
	return &h
}

func (f MessageHandlerFunc) GetName() string {
	v := reflect.ValueOf(f)
	if v.Kind() == reflect.Func {
		if rf := runtime.FuncForPC(v.Pointer()); rf != nil {
			return rf.Name()
		}
	}
	return v.String()
}
