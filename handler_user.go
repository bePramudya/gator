package main

import "fmt"

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("input username to login: login <username>")
	}

	if err := s.Config.SetUser(cmd.args[0]); err != nil {
		return err
	}

	fmt.Println("user has been set")

	return nil
}
