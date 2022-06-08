package process

import (
	"fmt"
	"os/exec"
	"sync"

	"github.com/kam1sh/overprocessed/pkg/process/api"
	"github.com/kam1sh/overprocessed/pkg/util"
)

type Pipe struct {
	From   api.Process
	Stream string
	To     api.Process
	Errs   []error
}

func (p *Pipe) Start() error {
	if p.Stream == api.STDOUT {
		reader, err := p.From.Cmd().StdoutPipe()
		if err != nil {
			return err
		}
		p.To.Cmd().Stdin = reader
	} else if p.Stream == api.STDERR {
		reader, err := p.From.Cmd().StderrPipe()
		if err != nil {
			return err
		}
		p.To.Cmd().Stdin = reader
	} else {
		return fmt.Errorf("invalid Stream provided")
	}
	err := p.To.Start()
	if err != nil {
		return err
	}
	err = p.From.Start()
	return err
}

func (p *Pipe) Wait() error {
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		p.Errs = append(p.Errs, p.To.Wait())
	}()
	go func() {
		defer wg.Done()
		p.Errs = append(p.Errs, p.From.Wait())
	}()
	wg.Wait()
	return util.MergeErrs(p.Errs)
}

func (p *Pipe) ExitCode() int {
	return p.To.ExitCode()
}

func (p *Pipe) Cmd() *exec.Cmd {
	return p.To.Cmd()
}
