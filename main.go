package main

import (
	"go-training-backend/cmd"
	"os"
)

func main() {
	if len(os.Args) > 0 {
		cmd.Cmd()
	}
}
