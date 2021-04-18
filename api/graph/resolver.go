package graph

import (
	"sync"

	"github.com/mocsiTeam/mocsiServer/api/graph/model"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	newUserRoom  map[string]chan *model.Room
	kickUserRoom map[string]chan *model.Room
	sync.Mutex
}

func (r *Resolver) InitResolver() {
	r.newUserRoom = make(map[string]chan *model.Room)
	r.kickUserRoom = make(map[string]chan *model.Room)
}
