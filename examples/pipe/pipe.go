package main

import (
	"bytes"
	"io"
	"log"
	"os"
	"os/exec"

	"github.com/tyndyll/bash"
)

func main() {
	pipe := &bash.Pipe{}
	pipe.Commands = []*exec.Cmd{
		exec.Command("ls", "/Users/tyndyll/Downloads"),
		exec.Command("grep", "as"),
		exec.Command("sort", "-r"),
	}
	var b bytes.Buffer
	if err := pipe.Execute(&b, &b); err != nil {
		log.Fatalln(err)
	}
	io.Copy(os.Stdout, &b)
}
