package api

import "os/exec"

var (
	STDIN  = "STDIN"
	STDOUT = "STDOUT"
	STDERR = "STDERR"
)

type Process interface {
	Start() error
	Wait() error
	ExitCode() int
	Cmd() *exec.Cmd
}
