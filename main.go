package main

import (
	"fmt"
	config "github.com/Yah-oo5726/blog-aggregator/internal/config"
)

func main() {
	configFile := config.Read()
	configFile.SetUser("Yah_oo5726")
	fmt.Printf("Username: %s, database url: %s\n", configFile.CurrentUserName, configFile.DbURL)
}
