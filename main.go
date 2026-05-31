package main

import (
	"fmt"
	config "github.com/Yah-oo5726/blog-aggregator/internal/config"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("not enough arguments passed in")
		os.Exit(1)
	}
	arguments := os.Args[1:]
	configFile := config.Read()
	program_state := state{configPtr: &configFile}
	program_commands := commands{functions: make(map[string]func(*state, command) error)}
	program_commands.register("login", handlerLogin)
	err := program_commands.run(&program_state, command{name: arguments[0], arguments: arguments[1:]})
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
