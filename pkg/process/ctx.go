package process

import (
	"context"
	"os"
	"os/exec"
	"time"

	"github.com/kam1sh/overprocessed/pkg/process/api"
)

type ProcessContext struct {
	Proc         api.Process
	Ctx          context.Context
	CancelSignal os.Signal
	Timeout      time.Duration
	StopSignal   os.Signal
}

func (c *ProcessContext) Start() error {
	return c.Proc.Start()
}

func (c *ProcessContext) Wait() error {
	return <-c.WaitCh()
}

func (c *ProcessContext) WaitCh() chan error {
	ch := make(chan error)
	done := make(chan struct{})
	go func() {
		defer close(ch)
		defer close(done)
		err := c.Proc.Wait()
		ch <- err
	}()
	go func() {
		select {
		case <-done:
		case <-c.Ctx.Done():
			process := c.Cmd().Process
			process.Signal(c.CancelSignal)
			select {
			case <-done:
			case <-time.After(c.Timeout):
				process.Signal(c.StopSignal)
			}
		}
	}()
	return ch
}

func (c *ProcessContext) ExitCode() int {
	return c.Proc.ExitCode()
}

func (c *ProcessContext) Cmd() *exec.Cmd {
	return c.Proc.Cmd()
}
