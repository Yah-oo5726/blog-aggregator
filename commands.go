package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/Yah-oo5726/blog-aggregator/internal/config"
	"github.com/Yah-oo5726/blog-aggregator/internal/database"
)

type state struct {
	db  *database.Queries
	cfg *config.Config
}

type command struct {
	name      string
	arguments []string
}

type commands struct {
	functions map[string]func(*state, command) error
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.arguments) == 0 {
		return errors.New("username is required")
	}
	s.cfg.SetUser(cmd.arguments[0])
	fmt.Println("user has been set.")
	return nil
}

func (c *commands) run(s *state, cmd command) error {
	function, exists := c.functions[cmd.name]
	if !exists {
		return errors.New("Function doesn't exist")
	}
	err := function(s, command{name: cmd.name, arguments: cmd.arguments})
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	return nil
}

func (c *commands) register(name string, f func(*state, command) error) error {
	c.functions[name] = f
	return nil
}
