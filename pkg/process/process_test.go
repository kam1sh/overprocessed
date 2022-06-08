package process_test

import (
	"context"
	"io/ioutil"
	"path"
	"syscall"
	"testing"
	"time"

	proc "github.com/kam1sh/overprocessed/pkg/process"
	procapi "github.com/kam1sh/overprocessed/pkg/process/api"
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

func TestPipe(t *testing.T) {
	pipe := proc.Pipe{
		From:   proc.NewProcess("echo", "123\n456"),
		Stream: procapi.STDOUT,
		To:     proc.NewProcess("grep", "2"),
	}
	require.NoError(t, pipe.Start())
	err := pipe.Wait()
	require.NoError(t, err)
}

func TestPipeline(t *testing.T) {
	buf := proc.MemoryBuffer()
	grep := proc.Redirect{
		From: &proc.Pipe{
			From:   proc.NewProcess("echo", "123\n456"),
			Stream: procapi.STDOUT,
			To:     proc.NewProcess("grep", "2"),
		},
		Stdout: buf,
	}
	err := grep.Start()
	require.NoError(t, err)
	require.NoError(t, grep.Wait())
	stdout := buf.ReadAll()
	require.NoError(t, err)
	require.Equal(t, "123\n", string(stdout))
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
