package main

import (
	"context"
	"fmt"
	"time"

	"github.com/ChipsAhoyEnjoyer/gator/internal/database"
	"github.com/google/uuid"
)

type command struct {
	name string
	args []string
}

type commands struct {
	registry map[string]func(*state, command) error
}

func newCommands() *commands {
	c := commands{
		registry: make(map[string]func(*state, command) error),
	}
	c.register("login", handlerLogin)
	c.register("register", handlerRegister)
	c.register("reset", handlerReset)
	c.register("users", handlerUsers)
	c.register("agg", handlerAgg)
	c.register("addfeed", handlerAddFeed)
	c.register("feeds", handlerFeeds)
	c.register("follow", handlerFollow)
	return &c
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.registry[name] = f
}

func (c *commands) run(s *state, cmd command) error {
	if function, ok := c.registry[cmd.name]; !ok {
		return fmt.Errorf("error: command '%v' does not exist", cmd.name)
	} else {
		err := function(s, cmd)
		if err != nil {
			return err
		}
	}
	return nil
}

func handlerFollow(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("usage: gator follow [link]")
	}

	return nil
}

func handlerFeeds(s *state, cmd command) error {
	if len(cmd.args) > 0 {
		return fmt.Errorf("error: cli <feed> [no args]")
	}
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return err
	}

	fmt.Println("==============================FEED==============================")
	for _, feed := range feeds {
		user, err := s.db.GetUserById(context.Background(), feed.UserID)
		if err != nil {
			return err
		}
		fmt.Printf("ID:        %v\n", feed.ID)
		fmt.Printf("Created:   %v\n", feed.CreatedAt)
		fmt.Printf("UpdatedAt: %v\n", feed.UpdatedAt)
		fmt.Printf("Title:     %v\n", feed.Name)
		fmt.Printf("Link:      %v\n", feed.Url)
		fmt.Printf("User:      %v\n", user.Name)
		fmt.Println()
		fmt.Println("================================================================")

	}
	return nil
}

func handlerAddFeed(s *state, cmd command) error {
	if len(cmd.args) != 2 {
		return fmt.Errorf("error: addfeed 'name' 'url'")
	}
	name := cmd.args[0]
	url := cmd.args[1]
	user, err := s.db.GetUser(
		context.Background(),
		s.cfg.CurrentUsername,
	)
	if err != nil {
		return fmt.Errorf("error fetching user: \n%v", err)
	}
	row, err := s.db.PostFeed(
		context.Background(),
		database.PostFeedParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Name:      name,
			Url:       url,
			UserID:    user.ID,
		},
	)
	if err != nil {
		return fmt.Errorf("error posting feed: \n%v", err)
	}
	fmt.Println("Feed posted successfully!")
	fmt.Println("================================================================")
	fmt.Printf("ID:        %v\n", row.ID)
	fmt.Printf("Created:   %v\n", row.CreatedAt)
	fmt.Printf("Updated:   %v\n", row.UpdatedAt)
	fmt.Printf("Title:     %v\n", row.Name)
	fmt.Printf("Link:      %v\n", row.Url)
	fmt.Printf("Posted by: %v / %v\n", user.Name, row.UserID)
	fmt.Println("================================================================")
	return nil
}

func handlerAgg(s *state, cmd command) error {
	feed, err := fetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return err
	}
	fmt.Printf("Title:       %v\n", feed.Channel.Title)
	fmt.Printf("Description: %v\n", feed.Channel.Description)
	fmt.Printf("Link:        %v\n\n", feed.Channel.Link)
	fmt.Println("============================CONTENT=============================")
	for i := range feed.Channel.Item {
		fmt.Printf(" - Title : %v \n", feed.Channel.Item[i].Title)
		fmt.Printf(" - Description: %v\n\n", feed.Channel.Item[i].Description)
		fmt.Println("================================================================")
		fmt.Println()
	}
	return nil
}

func handlerUsers(s *state, cmd command) error {
	if len(cmd.args) > 0 {
		return fmt.Errorf("error: too many arguments given; users expects zero arguments")
	}
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("error: could not retreive users from db \n%v", err)
	}
	fmt.Println("=============================USERS=============================")
	for _, user := range users {
		if user == s.cfg.CurrentUsername {
			fmt.Printf("* %v (current)\n", user)
		} else {
			fmt.Printf("* %v\n", user)
		}
	}
	return nil
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) >= 2 {
		return fmt.Errorf("error: too many arguments given; login expects one(username) argument")
	}
	username := cmd.args[0]
	if !userExists(s, username) {
		return fmt.Errorf("error: user not registered")
	}
	err := s.cfg.SetUser(username)
	if err != nil {
		return err
	}
	fmt.Printf("Now logged in as %v\n", username)
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) >= 2 {
		return fmt.Errorf("error: too many arguments given; register expects one(username) argument")
	}
	username := cmd.args[0]
	if userExists(s, username) {
		return fmt.Errorf("error: user exists")
	}
	u, err := s.db.CreateUser(
		context.Background(),
		database.CreateUserParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Name:      username,
		},
	)
	if err != nil {
		return fmt.Errorf("error: could not register user to database\n%v", err)
	}
	err = s.cfg.SetUser(u.Name)
	if err != nil {
		return fmt.Errorf("error: user registered but not logged in\n%v", err)
	}
	fmt.Println("================================================================")
	fmt.Printf("User '%v' created! \n", u.Name)
	fmt.Printf("Created: %v \n", u.CreatedAt)
	fmt.Printf("Updated: %v \n", u.UpdatedAt)
	fmt.Printf("ID:      %v \n", u.ID)
	fmt.Println("================================================================")
	return nil
}

func handlerReset(s *state, cmd command) error {
	if len(cmd.args) > 0 {
		return fmt.Errorf("error: too many arguments given; reset expects zero arguments")
	}
	usersDeleted, err := s.db.ResetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("error: users table reset unsuccessful \n%v", err)
	}
	fmt.Printf("Deleted %v user(s)\n", usersDeleted)
	return nil
}

func userExists(s *state, name string) bool {
	_, err := s.db.GetUser(
		context.Background(),
		name,
	)
	return err == nil
}
