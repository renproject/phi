package phi

import (
	"github.com/renproject/phi/parallel"
	"github.com/renproject/phi/task"
)

// Package `task` re-exports
type (
	// Message is an interface re-exported from package `task`.
	Message = task.Message

	// Runner is an interface re-exported from package `task`.
	Runner = task.Runner

	// Sender is an interface re-exported from package `task`.
	Sender = task.Sender

	// Task is an interface re-exported from package `task`.
	Task = task.Task

	// Handler is an interface re-exported from package `task`.
	Handler = task.Handler
)

var (
	// New is a function re-exported from package `task`.
	New = task.New
)

// Package `co` re-exports
var (
	// ParBegin is a function re-exported from package `co`.
	ParBegin = parallel.ParBegin

	// ParForAll is a function re-exported from package `co`.
	ParForAll = parallel.ParForAll

	// ForAll is a function re-exported from package `co`.
	ForAll = parallel.ForAll
)
