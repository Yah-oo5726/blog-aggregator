package main

import _ "github.com/lib/pq"

import (
	"database/sql"
	"fmt"
	"github.com/Yah-oo5726/blog-aggregator/internal/config"
	"github.com/Yah-oo5726/blog-aggregator/internal/database"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("not enough arguments passed in")
		os.Exit(1)
	}
	arguments := os.Args[1:]
	configFile := config.Read()
	db, err := sql.Open("postgres", "postgres://postgres:Postgres%2352@localhost:5432/gator?sslmode=disable")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	dbQueries := database.New(db)
	program_state := state{cfg: &configFile, db: dbQueries}
	program_commands := commands{functions: make(map[string]func(*state, command) error)}
	program_commands.register("login", handlerLogin)
	program_commands.register("register", handlerRegister)
	program_commands.register("reset", handlerReset)
	program_commands.register("users", handlerGetUsers)
	program_commands.register("agg", handlerAgg)
	err = program_commands.run(&program_state, command{name: arguments[0], arguments: arguments[1:]})
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
