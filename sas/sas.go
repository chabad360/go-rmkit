package sas

import (
	"os/exec"
	"strconv"
	"sync"
)

var id = 1

func ID() string {
	id++
	return strconv.Itoa(id)
}

type App[State any] interface {
	Render(State) (Scene, error)
}

type Simple struct {
	lock   *sync.Mutex
	simple *exec.Cmd
}

type Scene struct {
	Widgets []Widget
}

type Size int

const (
	Same Size = -2
	Step      = -1
)

type Widget struct {
	id string
	x  int
	y  int
	w  int
	h  int
}
