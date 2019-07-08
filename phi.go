package phi

import (
	"github.com/renproject/phi/co"
	"github.com/renproject/phi/task"
)

// Package `task` re-exports
type (
	// Message is an interface re-exported from package `task`.
	Message = task.Message

	// Messages is struct re-exported from package `task`.
	Messages = task.Messages

	// Runner is an interface re-exported from package `task`.
	Runner = task.Runner

	// Sender is an interface re-exported from package `task`.
	Sender = task.Sender

	// Task is an interface re-exported from package `task`.
	Task = task.Task

	// Options is a struct re-exported from package `task`.
	Options = task.Options

	// Reducer is an interface re-exported from package `task`.
	Reducer = task.Reducer

	// Router is an interface re-exported from package `task`.
	Router = task.Router
)

var (
	// New is a function re-exported from package `task`.
	New = task.New

	// NewRouter is a function re-exported from package `task`.
	NewRouter = task.NewRouter
)

// Package `co` re-exports
var (
	// ParBegin is a function re-exported from package `co`.
	ParBegin = co.ParBegin

	// ParForAll is a function re-exported from package `co`.
	ParForAll = co.ParForAll

	// ForAll is a function re-exported from package `co`.
	ForAll = co.ForAll
)
