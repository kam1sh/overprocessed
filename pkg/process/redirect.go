package process

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"

	"github.com/kam1sh/overprocessed/pkg/process/api"
)

type Redirect struct {
	From        api.Process
	Stdin       PipeInterceptor
	Stdout      PipeInterceptor
	Stderr      PipeInterceptor
	beforeFuncs []func() error
	afterFuncs  []func() error
}

func (r *Redirect) Start() error {
	r.beforeFuncs = make([]func() error, 0)
	r.afterFuncs = make([]func() error, 0)
	err := r.setIO()
	if err != nil {
		return err
	}
	if err = r.From.Start(); err != nil {
		return err
	}
	if err = r.before(); err != nil {
		return err
	}
	return nil
}

func (r *Redirect) before() error {
	for _, v := range r.beforeFuncs {
		if err := v(); err != nil {
			return err
		}
	}
	return nil
}

func (r *Redirect) Wait() error {
	err := r.closeIO()
	if err != nil {
		return err
	}
	return r.From.Wait()
}

func (r *Redirect) ExitCode() int {
	return r.From.ExitCode()
}

func (r *Redirect) Cmd() *exec.Cmd {
	return r.From.Cmd()
}

func (r *Redirect) closeIO() error {
	for _, v := range r.afterFuncs {
		if err := v(); err != nil {
			return err
		}
	}
	return nil
}

func (r *Redirect) redirectIO(interceptor PipeInterceptor, stream string) error {
	if interceptor == nil {
		return nil
	}
	log.Println("Setting", interceptor, "to", stream)
	if err := interceptor.Intercept(r.From.Cmd(), stream); err != nil {
		return err
	}
	cmd, ok := interceptor.(api.Process)
	if ok {
		r.beforeFuncs = append(r.beforeFuncs, cmd.Start)
		r.afterFuncs = append(r.afterFuncs, cmd.Wait)
	}
	closer, ok := interceptor.(io.Closer)
	if ok {
		r.afterFuncs = append(r.afterFuncs, closer.Close)
	}
	return nil
}

func (r *Redirect) setIO() (err error) {
	if err = r.redirectIO(r.Stdin, api.STDIN); err != nil {
		return
	}
	if err = r.redirectIO(r.Stdout, api.STDOUT); err != nil {
		return
	}
	if err = r.redirectIO(r.Stderr, api.STDERR); err != nil {
		return
	}
	return nil
}

type PipeInterceptor interface {
	Intercept(cmd *exec.Cmd, stream string) error
}

///////////////
/// sources ///
///////////////

type FileSource struct {
	Path string
	file *os.File
}

func (s *FileSource) Intercept(cmd *exec.Cmd, stream string) error {
	if stream != api.STDIN {
		return fmt.Errorf("FileSource does not accept %v stream", stream)
	}
	fd, err := os.Open(s.Path)
	if err != nil {
		return err
	}
	s.file = fd
	cmd.Stdin = fd
	return nil
}

func (s *FileSource) Close() error {
	return s.file.Close()
}

type ParentSource struct{}

func (s *ParentSource) Intercept(cmd *exec.Cmd, stream string) error {
	if stream != api.STDIN {
		return fmt.Errorf("ParentSource does not accept %v stream", stream)
	}
	cmd.Stdin = os.Stdin
	return nil
}

////////////////////
/// destinations ///
////////////////////

type FileDestination struct {
	Path string
	file *os.File
}

func (d *FileDestination) Intercept(cmd *exec.Cmd, stream string) error {
	fd, err := os.Create(d.Path)
	if err != nil {
		return err
	}
	d.file = fd
	switch stream {
	case api.STDOUT:
		cmd.Stdout = fd
	case api.STDERR:
		cmd.Stderr = fd
	default:
		return fmt.Errorf("FileDestination does not accept %v stream", stream)
	}
	return nil
}

func (d *FileDestination) Close() error {
	return d.file.Close()
}

type StdoutDestination struct{}

func (s *StdoutDestination) Intercept(cmd *exec.Cmd, stream string) error {
	switch stream {
	case api.STDOUT:
		cmd.Stdout = os.Stdout
	case api.STDERR:
		cmd.Stderr = os.Stdout
	default:
		return fmt.Errorf("StdoutDestination does not accept %v stream", stream)
	}
	return nil
}

type StderrDestination struct{}

func (s *StderrDestination) Intercept(cmd *exec.Cmd, stream string) error {
	switch stream {
	case api.STDOUT:
		cmd.Stdout = os.Stderr
	case api.STDERR:
		cmd.Stderr = os.Stderr
	default:
		return fmt.Errorf("StderrDestination does not accept %v stream", stream)
	}
	return nil
}

type BufferedWriter struct {
	buf bytes.Buffer
}

func MemoryBuffer() *BufferedWriter {
	return &BufferedWriter{}
}

func (w *BufferedWriter) Intercept(cmd *exec.Cmd, stream string) error {
	switch stream {
	case api.STDOUT:
		cmd.Stdout = &w.buf
	case api.STDERR:
		cmd.Stderr = &w.buf
	}
	return nil
}

func (w *BufferedWriter) ReadAll() ([]byte, error) {
	return w.buf.Bytes(), nil
}
