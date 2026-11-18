package main

import (
	"context"
	"fmt"
	"slices"
	"time"

	"github.com/Cyrus-Gahatraj/gator/internal/database"
	"github.com/google/uuid"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("username is required")
	}

	name := cmd.args[0]
	if _, err := s.db.GetUser(context.Background(), name); err != nil {
		return fmt.Errorf("No user found '%v'", err)
	}

	if err := s.cfg.SetUser(name); err != nil {
		return err
	}
	fmt.Printf("The username has been set to %s\n", s.cfg.CurrentUserName)

	return nil
}

func register(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("username is required")
	}

	currentTime := time.Now()
	userParams := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: currentTime,
		UpdatedAt: currentTime,
		Name:      cmd.args[0],
	}

	user, err := s.db.CreateUser(context.Background(), userParams)
	if err != nil {
		return fmt.Errorf("Error while creating the user: %v\n", err)
	}

	if err := s.cfg.SetUser(user.Name); err != nil {
		return err
	}

	return nil
}

func reset(s *state, cmd command) error {
	if len(cmd.args) != 0 {
		return fmt.Errorf("reset need no arguments")
	}

	if err := s.db.ResetAllUser(context.Background()); err != nil {
		return fmt.Errorf("Error: %v\n", err)
	}

	return nil
}

func getUsers(s *state, cmd command) error {
	if len(cmd.args) != 0 {
		return fmt.Errorf("users need no arguments")
	}

	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("error: %v\n", err)
	}

	currentUser := s.cfg.CurrentUserName
	if ok := slices.Contains(users, currentUser); ok {
		users[slices.Index(users, currentUser)] = fmt.Sprintf("%v (current)", currentUser)
	}

	for _, user := range users {
		fmt.Printf("* %v\n", user)
	}

	return nil
}
