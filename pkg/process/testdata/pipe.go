package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	scan := bufio.NewScanner(os.Stdin)
	for scan.Scan() {
		fmt.Println(scan.Text())
	}
}
