package process_test

import (
	"context"
	"io/ioutil"
	"path"
	"runtime"
	"syscall"
	"testing"
	"time"

	proc "github.com/kam1sh/overprocessed/pkg/process"
	"github.com/stretchr/testify/require"
)

func TestRedirect(t *testing.T) {
	pth := path.Join(t.TempDir(), "output.txt")
	echo := proc.Redirect{
		From: proc.NewProcess("echo", "123"),
		Stdout: &proc.FileDestination{
			Path: pth,
		},
		Stderr: &proc.StderrDestination{},
	}
	err := echo.Start()
	require.NoError(t, err)
	err = echo.Wait()
	require.NoError(t, err)
	require.Equal(t, 0, echo.ExitCode())
	data, err := ioutil.ReadFile(pth)
	require.NoError(t, err)
	require.Equal(t, []byte("123\n"), data)
}

func TestRedirectBuf(t *testing.T) {
	buf := proc.MemoryBuffer()
	echo := proc.Redirect{
		From:   proc.NewProcess("echo", "123"),
		Stdout: buf,
	}
	err := echo.Start()
	require.NoError(t, err)
	err = echo.Wait()
	require.NoError(t, err)
}

func TestTerminate(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	sleep := proc.ProcessContext{
		Proc:         proc.NewProcess("sleep", "100"),
		Ctx:          ctx,
		CancelSignal: syscall.SIGINT,
		Timeout:      1 * time.Second,
		StopSignal:   syscall.SIGKILL,
	}
	sleep.Start()
	done := make(chan struct{})
	var err error
	go func() {
		defer close(done)
		err = sleep.Wait()
	}()
	time.Sleep(2 * time.Second)
	cancel()
	<-done
	require.NoError(t, err)
}

func TestRedirectProcess(t *testing.T) {
	goexec := path.Join(runtime.GOROOT(), "bin", "go")
	scanner := proc.NewScannerDest()
	pipeline := &proc.Redirect{
		From:   proc.NewProcess(goexec, "run", path.Join("testdata", "pipe.go")),
		Stdin:  proc.NewProcess(goexec, "run", path.Join("testdata", "gen.go")),
		Stdout: scanner,
	}
	t.Log(proc.DumpTree(pipeline))
	out := make([]string, 0, 10)
	require.NoError(t, pipeline.Start())
	for scanner.Scan() {
		out = append(out, scanner.Text())
	}
	require.NoError(t, pipeline.Wait())
	require.Len(t, out, 10)
}
