package main

import (
	"log"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("invalid count of arguments")
		return
	}
	env, err := ReadDir(os.Args[1])
	if err != nil {
		log.Fatal(err)
		return
	}

	os.Exit(RunCmd(os.Args[2:], env))
}
