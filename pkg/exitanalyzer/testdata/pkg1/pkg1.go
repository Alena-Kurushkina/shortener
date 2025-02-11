package main

import (
	"os"
)

func Exit(){
	os.Exit(3)
}

func main() {
	Exit()
    os.Exit(3)           // want "call os.Exit"
} 