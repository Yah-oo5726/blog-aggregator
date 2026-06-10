package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
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
		return err
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
		return err
	}
	s.cfg.SetUser(cmd.arguments[0])
	fmt.Println("user was created and logged into")
	fmt.Printf("uuid %v created at %v updated at %v name %v\n", user.ID, user.CreatedAt, user.UpdatedAt, user.Name)
	return nil
}

func handlerReset(s *state, cmd command) error {
	err := s.db.DeleteUsers(context.Background())
	if err != nil {
		return err
	}
	fmt.Println("successful reset")
	return nil
}

func handlerGetUsers(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return err
	}
	for _, user := range users {
		fmt.Printf("* %s", user.Name)
		if user.Name == s.cfg.CurrentUserName {
			fmt.Print(" (current)")
		}
		fmt.Print("\n")
	}
	return nil
}

func handlerAgg(s *state, cmd command) error {
	output, err := fetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return err
	}
	fmt.Println(output)
	return nil
}

func handlerAddFeed(s *state, cmd command) error {
	if len(cmd.arguments) < 2 {
		return fmt.Errorf("Not enough arguments")
	}
	user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		return err
	}
	feed, err := s.db.AddFeed(context.Background(), database.AddFeedParams{ID: uuid.New(), CreatedAt: time.Now(), UpdatedAt: time.Now(), Name: cmd.arguments[0], Url: cmd.arguments[1], UserID: user.ID})
	if err != nil {
		return err
	}
	fmt.Println(feed)
	err = handlerFollow(s, command{name: "follow", arguments: cmd.arguments[1:]})
	return nil
}

func handlerFeeds(s *state, cmd command) error {
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return err
	}
	for _, feed := range feeds {
		info, err := s.db.GetFeedInfo(context.Background(), feed.ID)
		if err != nil {
			return err
		}
		_, err = fmt.Printf("* %v\n  - %v\n  - %v\n", info.Name, info.Url, info.Name_2)
		if err != nil {
			return err
		}
	}
	return nil
}

func handlerFollow(s *state, cmd command) error {
	if len(cmd.arguments) == 0 {
		return errors.New("url is required")
	}
	user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		return err
	}
	feed, err := s.db.GetByUrl(context.Background(), cmd.arguments[0])
	if err != nil {
		return err
	}
	_, err = s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{ID: uuid.New(), CreatedAt: sql.NullTime{Time: time.Now(), Valid: true}, UpdatedAt: sql.NullTime{Time: time.Now(), Valid: true}, UserID: user.ID, FeedID: feed.ID})
	if err != nil {
		return err
	}
	fmt.Printf("%s followed %s\n", user.Name, feed.Name)
	return nil
}

func handlerFollowing(s *state, cmd command) error {
	user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		return err
	}
	response, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return err
	}
	fmt.Printf("%s is following\n", s.cfg.CurrentUserName)
	for _, follow := range response {
		fmt.Println(follow.FeedName)
	}
	return nil
}

func (c *commands) run(s *state, cmd command) error {
	function, exists := c.functions[cmd.name]
	if !exists {
		return errors.New("Function doesn't exist")
	}
	err := function(s, command{name: cmd.name, arguments: cmd.arguments})
	if err != nil {
		return err
	}
	return nil
}

func (c *commands) register(name string, f func(*state, command) error) error {
	c.functions[name] = f
	return nil
}
