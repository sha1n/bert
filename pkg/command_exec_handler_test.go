package pkg

import (
	"os"
	"os/exec"
	"reflect"
	"runtime"
	"testing"
	"time"

	"github.com/sha1n/benchy/api"
	"github.com/stretchr/testify/assert"
)

func TestRunCommandFnFor(t *testing.T) {
	type args struct {
		ce *commandExecutor
	}
	tests := []struct {
		name string
		args args
		want RunCommandFn
	}{
		{name: "with stdout piping", args: args{ce: &commandExecutor{pipeStdout: true, pipeStderr: false}}, want: RunCommand},
		{name: "without stdout piping", args: args{ce: &commandExecutor{pipeStdout: false, pipeStderr: false}}, want: RunCommandWithProgressIndicator},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RunCommandFnFor(tt.args.ce)
			assert.Equal(t, functionName(tt.want), functionName(got))
		})
	}
}

func TestRunCommand(t *testing.T) {
	type args struct {
		cmd *exec.Cmd
		ctx api.IOContext
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "with expected error", args: args{cmd: exec.Command("non-existing-command"), ctx: NewIOContextTty(false)}, wantErr: true},
		{name: "with no error", args: args{cmd: exec.Command("go", "version"), ctx: NewIOContextTty(false)}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := RunCommand(tt.args.cmd, tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("RunCommand() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRunCommandWithProgressIndicator(t *testing.T) {
	type args struct {
		cmd *exec.Cmd
		ctx api.IOContext
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "with expected error", args: args{cmd: exec.Command("non-existing-command"), ctx: NewIOContextTty(false)}, wantErr: true},
		{name: "with no error", args: args{cmd: exec.Command("go", "version"), ctx: NewIOContextTty(false)}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := RunCommandWithProgressIndicator(tt.args.cmd, tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("RunCommandWithProgressIndicator() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRegisterInterruptGuard(t *testing.T) {
	execCmd := exec.Command("test", "--command")
	call := make(chan bool)
	_, c := registerInterruptGuard(execCmd, func(c *exec.Cmd, s os.Signal) {
		call <- true
	})

	c <- os.Interrupt
	assert.Eventually(t, func() bool { return <-call }, time.Second*10, time.Millisecond)
}

func TestRegisterInterruptGuardCancellation(t *testing.T) {
	expectPanic := func() {
		v := recover()
		assert.NotNil(t, v)
	}
	defer expectPanic()

	cancel, c := registerInterruptGuard(aCommand(), func(c *exec.Cmd, s os.Signal) {})
	cancel()

	c <- os.Interrupt // this should fail panic because the channel is closed
}

func functionName(f interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
}

func NewIOContextTty(tty bool) api.IOContext {
	ctx := api.NewIOContext()
	ctx.Tty = tty

	return ctx
}
