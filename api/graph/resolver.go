package graph

import (
	"sync"

	"github.com/mocsiTeam/mocsiServer/api/graph/model"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	newUsersRoom map[string]chan *model.Room
	sync.Mutex
}

func (r *Resolver) InitResolver() {
	r.newUsersRoom = make(map[string]chan *model.Room)
}
