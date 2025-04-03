package main

import (
	"context"
	"fmt"
	"time"

	"github.com/ChipsAhoyEnjoyer/gator/internal/database"
	"github.com/google/uuid"
)

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 2 {
		return fmt.Errorf("usage: cli addfeed [name] [url]")
	}
	name := cmd.args[0]
	url := cmd.args[1]
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
	_, err = s.db.CreateFeedFollow(
		context.Background(),
		database.CreateFeedFollowParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			UserID:    user.ID,
			FeedID:    row.ID,
		},
	)
	if err != nil {
		fmt.Println()
		return fmt.Errorf("error: feed not added to user's following \n%v", err)
	}
	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 0 {
		return fmt.Errorf("usage: cli following")
	}
	follows, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return err
	}
	fmt.Println("=============================FOLLOWS============================")
	fmt.Println()
	for _, feed := range follows {
		fmt.Printf("Name: %v\n", feed.FeedName)
	}
	fmt.Println()
	fmt.Println("================================================================")
	return nil
}

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("usage: cli follow [link]")
	}
	url := cmd.args[0]
	feed, err := s.db.GetFeedByURL(
		context.Background(),
		url,
	)
	if err != nil {
		fmt.Println("feed not registered")
		fmt.Println("use 'cli addfeed [name] [url]' to add feed")
		fmt.Println("use 'cli feeds' see existing feeds")
		return err
	}
	row, err := s.db.CreateFeedFollow(
		context.Background(),
		database.CreateFeedFollowParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			UserID:    user.ID,
			FeedID:    feed.ID,
		},
	)
	if err != nil {
		return err
	}
	fmt.Println("================================================================")
	fmt.Printf("%v now following %v!\n", row.UserName, row.FeedName)
	fmt.Println("================================================================")
	return nil
}

func handlerFeeds(s *state, cmd command) error {
	if len(cmd.args) > 0 {
		return fmt.Errorf("usage: cli feed")
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
	if len(cmd.args) != 1 {
		return fmt.Errorf("usage: cli [username]")
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
	if len(cmd.args) != 1 {
		return fmt.Errorf("usage: cli [username]")
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
