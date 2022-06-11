package process

import (
	"bufio"
	"os/exec"

	"github.com/kam1sh/overprocessed/pkg/process/api"
)

type ScannerDestination struct {
	*bufio.Scanner
}

func NewScannerDest() *ScannerDestination {
	return &ScannerDestination{}
}

func (s *ScannerDestination) Intercept(cmd *exec.Cmd, stream string) error {
	switch stream {
	case api.STDOUT:
		out, err := cmd.StdoutPipe()
		if err != nil {
			return err
		}
		s.Scanner = bufio.NewScanner(out)
	case api.STDERR:
		out, err := cmd.StderrPipe()
		if err != nil {
			return err
		}
		s.Scanner = bufio.NewScanner(out)
	}
	return nil
}

func (s *ScannerDestination) String() string {
	return "[Scanner]"
}
