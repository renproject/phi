# `φ phi`

[![GoDoc](https://godoc.org/github.com/renproject/phi?status.svg)](https://godoc.org/github.com/renproject/phi)
[![CircleCI](https://circleci.com/gh/renproject/phi/tree/master.svg?style=shield)](https://circleci.com/gh/renproject/phi/tree/master)
![Go Report](https://goreportcard.com/badge/github.com/renproject/phi)
[![Coverage Status](https://coveralls.io/repos/github/renproject/phi/badge.svg?branch=master)](https://coveralls.io/github/renproject/phi?branch=master)
[![License: MIT](https://img.shields.io/badge/license-MIT-green.svg)](https://opensource.org/licenses/MIT)

A simple message passing framework for Golang, free from deadlocks. It support both synchronous and asynchronous message passing.

## Data Parallelism

Phi offers two types of data parallelism: implicit and explicit. Implicit parallelism means that Cogo will control how many goroutines are used to introduce parallelism, in contrast to explicit parallelism which gives control to the user.

### Implicit Parallelism

In the following example, we use the `parallel.ForAll` function to loop over different iterators. Iterators are any value that makes sense to loop over: arrays, slices, maps, and integers. Cogo will use one goroutine per CPU core available and we cannot make any assumptions about which iteration will run on which goroutine. Calling `parallel.ForAll` will block until all iterations have finished running.

```go
// Fill an array of integers with random values
xs := [10]int{}
parallel.ForAll(xs, func(i int) {
    xs[i] = rand.Intn(10)
})

// Map those random values to booleans
ys := [10]bool{}
parallel.ForAll(10, func(i int) {
    ys[i] = xs[i] > 5
})
```

In the following example, we use the `parallel.Begin` function to run distinct tasks. As before, Cogo will use one goroutine per CPU core available and map the different tasks over these goroutines. Calling `parallel.Begin` will block until all tasks have finished running.

```go
parallel.Begin(
    func() {
        log.Info("[task 1] when will this print?")
    },
    func() {
        log.Info("[task 2] who knows?")
    },
    func() {
        log.Info("[task 3] implicit parallelism is great!")
    })
```

### Explicit Parallelism

In the following example, we use the `parallel.ParForAll` function to loop over different iterators. As with implicitly parallel loops, iterators are any value that makes sense to loop over: arrays, slices, maps, and integers. Unlike implicitly parallel loops, Cogo will use one goroutine per iteration, regardless of the number of CPU cores available. Calling `parallel.ParForAll` will block until all iterations have finished running.

```go
// Fill an array of integers with random values
xs := [10]int{}
parallel.ParForAll(xs, func(i int) {
    xs[i] = rand.Intn(10)
})

// Map those random values to booleans
ys := [10]bool{}
parallel.ParForAll(10, func(i int) {
    ys[i] = xs[i] > 5
})
```

In the following example, we use the `parallel.ParBegin` function to run distinct tasks. Similar to the `parallel.ParForAll` function, Cogo will use one goroutine per task. Calling `parallel.ParBegin` will block until all tasks have finished running.

```go
parallel.ParBegin(
    func() {
        log.Info("[task 1] when will this print?")
    },
    func() {
        log.Info("[task 2] who knows?")
    },
    func() {
        log.Info("[task 3] explicit parallelism is great!")
    })
```

## Contribution

Built with ❤ by Ren.
