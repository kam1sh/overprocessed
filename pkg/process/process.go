package process

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/kam1sh/overprocessed/pkg/process/api"
)

type Process struct {
	cmd      *exec.Cmd
	exitCode int
}

func NewProcess(cmd string, args ...string) *Process {
	c := exec.Command(cmd, args...)
	return &Process{
		cmd:      c,
		exitCode: -1,
	}
}

func NewProcessEnv(cmd string, args []string, env map[string]string, includeSystem bool) api.Process {
	c := exec.Command(cmd, args...)
	if includeSystem {
		for _, v := range os.Environ() {
			name, value, _ := strings.Cut(v, "=")
			if _, exists := env[name]; !exists {
				env[name] = value
			}
		}
	}
	envArr := make([]string, len(env))
	for k, v := range env {
		envArr = append(envArr, fmt.Sprintf("%v=%v", k, v))
	}
	c.Env = envArr
	return &Process{
		cmd:      c,
		exitCode: -1,
	}
}

func (p *Process) Wait() error {
	err := p.cmd.Wait()
	if _, ok := err.(*exec.ExitError); ok {
		err = nil
	}
	p.exitCode = p.cmd.ProcessState.ExitCode()
	return err
}

func (p *Process) ExitCode() int {
	return p.exitCode
}

func (p *Process) Start() error {
	log.Println("Starting", p.cmd.Args)
	err := p.cmd.Start()
	return err
}

func (p *Process) Cmd() *exec.Cmd {
	return p.cmd
}

func (p *Process) Intercept(cmd *exec.Cmd, stream string) (err error) {
	switch stream {
	case api.STDIN:
		cmd.Stdin, err = p.cmd.StdoutPipe()
		if err != nil {
			return
		}
	case api.STDERR:
		p.cmd.Stdin, err = cmd.StderrPipe()
		if err != nil {
			return
		}
	case api.STDOUT:
		p.cmd.Stdin, err = cmd.StdoutPipe()
		if err != nil {
			return
		}
	}
	return nil
}
