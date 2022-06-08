package process

import (
	"bytes"
	"io"
	"os"
	"os/exec"

	"github.com/kam1sh/overprocessed/pkg/process/api"
	"github.com/kam1sh/overprocessed/pkg/util"
)

type Redirect struct {
	From   api.Process
	Stdin  Source
	Stdout Destination
	Stderr Destination
}

func (r *Redirect) Start() error {
	err := r.setIO()
	if err != nil {
		return err
	}
	err = r.From.Start()
	if err != nil {
		return err
	}
	return nil
}

func (r *Redirect) Wait() error {
	errs := make([]error, 0)
	errs = append(errs, r.From.Wait())
	for _, e := range r.closeIO() {
		if e != nil {
			errs = append(errs, e)
		}
	}
	return util.MergeErrs(errs)
}

func (r *Redirect) ExitCode() int {
	return r.From.ExitCode()
}

func (r *Redirect) Cmd() *exec.Cmd {
	return r.From.Cmd()
}

func (r *Redirect) closeIO() []error {
	streams := []interface{}{r.Stdin, r.Stdout, r.Stderr}
	errs := make([]error, 3)
	for i, stream := range streams {
		closer, ok := stream.(io.Closer)
		if ok {
			errs[i] = closer.Close()
		}
	}
	return errs
}

func (r *Redirect) setIO() error {
	c := r.Cmd()
	if r.Stdin != nil {
		reader, err := r.Stdin.Reader()
		if err != nil {
			return err
		}
		c.Stdin = reader
	}
	if r.Stdout != nil {
		writer, err := r.Stdout.Writer()
		if err != nil {
			return err
		}
		c.Stdout = writer
	}
	if r.Stderr != nil {
		writer, err := r.Stderr.Writer()
		if err != nil {
			return err
		}
		c.Stderr = writer
	}
	return nil
}

///////////////
/// sources ///
///////////////

type Source interface {
	Reader() (io.Reader, error)
}

type FileSource struct {
	Path string
	file *os.File
}

func (s *FileSource) Reader() (io.Reader, error) {
	fd, err := os.Open(s.Path)
	if err != nil {
		return nil, err
	}
	s.file = fd
	return fd, nil
}

func (s *FileSource) Close() error {
	return s.file.Close()
}

type ParentSource struct{}

func (s *ParentSource) Reader() (io.Reader, error) {
	return os.Stdin, nil
}

////////////////////
/// destinations ///
////////////////////

type Destination interface {
	Writer() (io.Writer, error)
}

type FileDestination struct {
	Path string
	file *os.File
}

func (d *FileDestination) Writer() (io.Writer, error) {
	fd, err := os.Create(d.Path)
	if err != nil {
		return nil, err
	}
	d.file = fd
	return fd, nil
}

func (d *FileDestination) Close() error {
	return d.file.Close()
}

type StdoutDestination struct{}

func (s *StdoutDestination) Writer() (io.Writer, error) {
	return os.Stdout, nil
}

type StderrDestination struct{}

func (s *StderrDestination) Writer() (io.Writer, error) {
	return os.Stderr, nil
}

type BufferedWriter struct {
	buf bytes.Buffer
}

func MemoryBuffer() *BufferedWriter {
	return &BufferedWriter{}
}

func (w *BufferedWriter) Writer() (io.Writer, error) {
	return &w.buf, nil
}

func (w *BufferedWriter) ReadAll() []byte {
	return w.buf.Bytes()
}
