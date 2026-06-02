package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/Yah-oo5726/blog-aggregator/internal/config"
	"github.com/Yah-oo5726/blog-aggregator/internal/database"
	"github.com/google/uuid"
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
	if _, err := s.db.GetUser(context.Background(), cmd.arguments[0]); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	s.cfg.SetUser(cmd.arguments[0])
	fmt.Println("user has been set.")
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.arguments) == 0 {
		return errors.New("username is required")
	}
	time := time.Now()
	user, err := s.db.CreateUser(context.Background(), database.CreateUserParams{ID: uuid.New(), CreatedAt: time, UpdatedAt: time, Name: cmd.arguments[0]})
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	s.cfg.SetUser(cmd.arguments[0])
	fmt.Println("user was created and logged into")
	fmt.Printf("uuid %v created at %v updated at %v name %v\n", user.ID, user.CreatedAt, user.UpdatedAt, user.Name)
	return nil
}

func handlerReset(s *state, cmd command) error {
	err := s.db.DeleteUsers(context.Background())
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	fmt.Println("successful reset")
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
