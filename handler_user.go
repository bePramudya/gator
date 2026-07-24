package main

import (
	"context"
	"fmt"
	"time"

	"github.com/bePramudya/gator/internal/database"
	"github.com/google/uuid"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("input username to login: login <username>")
	}

	_, err := s.db.GetUser(context.Background(), cmd.args[0])
	if err != nil {
		fmt.Println("user doesn't exist")
		return err
	}

	if err := s.cfg.SetUser(cmd.args[0]); err != nil {
		return err
	}

	fmt.Println("user has been set")

	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("input username: register <username>")
	}

	var args = database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.args[0],
	}

	user, err := s.db.CreateUser(context.Background(), args)
	if err != nil {
		return err
	}

	s.cfg.SetUser(user.Name)

	fmt.Printf("user '%s' has been set \n", user.Name)
	fmt.Println(user)

	return nil
}

func handlerReset(s *state, cmd command) error {
	if err := s.db.DeleteUsers(context.Background()); err != nil {
		fmt.Println("error deleting users data")
		return err
	}

	fmt.Println("success delete users data")
	return nil
}

func handlerGetUsers(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return err
	}

	for _, user := range users {
		if s.cfg.CurrentUserName == user.Name {
			fmt.Printf(" * %v (current)\n", user.Name)
			continue
		}
		fmt.Printf(" * %v\n", user.Name)
	}

	return nil
}
