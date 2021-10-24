package test

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Runner struct {
	T       *testing.T
	cleanup []func() error
	mu      *sync.Mutex
}

func NewRunner(t *testing.T) *Runner {
	return &Runner{
		T:  t,
		mu: &sync.Mutex{},
	}
}

func (tr *Runner) AddCleanupStep(c func() error) {
	tr.mu.Lock()
	defer tr.mu.Unlock()

	tr.cleanup = append(tr.cleanup, c)
}

func (tr *Runner) Clean() {
	tr.T.Helper()

	for _, f := range tr.cleanup {
		assert.NoError(tr.T, f())
	}
}
