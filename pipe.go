package bash

import (
	"bytes"
	"io"

	"os/exec"
)

type Pipe struct {
	Commands []*exec.Cmd
}

func (p *Pipe) Execute(stdout *bytes.Buffer, stderr *bytes.Buffer) error {
	if len(p.Commands) == 0 {
		return nil
	}
	pipe_stack := make([]*io.PipeWriter, len(p.Commands)-1)
	i := 0
	for ; i < len(p.Commands)-1; i++ {
		stdin_pipe, stdout_pipe := io.Pipe()
		p.Commands[i].Stdout = stdout_pipe
		p.Commands[i].Stderr = stderr
		p.Commands[i+1].Stdin = stdin_pipe
		pipe_stack[i] = stdout_pipe
	}
	p.Commands[i].Stdout = stdout
	p.Commands[i].Stderr = stderr

	return p.call(p.Commands, pipe_stack)
}

func (p *Pipe) call(stack []*exec.Cmd, pipes []*io.PipeWriter) (err error) {
	if stack[0].Process == nil {
		if err = stack[0].Start(); err != nil {
			return err
		}
	}
	if len(stack) > 1 {
		if err = stack[1].Start(); err != nil {
			return err
		}
		defer func() {
			if err == nil {
				pipes[0].Close()
				err = p.call(stack[1:], pipes[1:])
			}
		}()
	}
	return stack[0].Wait()
}
