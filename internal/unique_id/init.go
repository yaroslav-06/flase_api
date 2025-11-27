package uniqueid

import (
	"fmt"
	"math/rand"
)

type Generator struct {
	st *OnSet
}

type OnSet interface {
	Exists(id string) bool
}

func NewGenerator(st OnSet) *Generator {
	gr := &Generator{st: &st}
	return gr
}

func (gr *Generator) GetNewId() string {
	for {
		new_id := fmt.Sprint(rand.Int63())
		if !(*gr.st).Exists(new_id) {
			return new_id
		}
	}
}
