package main

import (
	"os"
)

func exit() {
	os.Exit(3)
}

func main() {
	exit()
	os.Exit(3) // want "call os.Exit"
}
