package main

import (
	"fmt"

	"github.com/alvise88/zero-turnaround-cicd-with-dagger/cmd/calc/cmd"
)

func main() {
	err := cmd.Root().Execute()

	if err != nil {
		fmt.Println("unable to run calc cli")
	}
}
