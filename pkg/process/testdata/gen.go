package main

import (
	"fmt"
	"os"
)

func main() {
	for i := 0; i < 10; i++ {
		fmt.Fprintln(os.Stdout, i)
		fmt.Fprintln(os.Stderr, i)
	}
}