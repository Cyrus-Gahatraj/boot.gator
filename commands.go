package main

import (
	"fmt"
	"os"
)

type command struct {
	name string
	args []string
}

type commands struct {
	handler map[string]func(*state, command) error
}

func newCommands() commands {
	c := commands{
		handler: make(map[string]func(*state, command) error),
	}

	c.register("login", handlerLogin)
	c.register("register", register)
	c.register("reset", reset)
	c.register("users", getUsers)
	c.register("agg", aggregator)
	c.register("addfeed", middlewareLoggedIn(addFeed))
	c.register("feeds", feeds)
	c.register("follow", middlewareLoggedIn(follow))
	c.register("following", middlewareLoggedIn(following))
	c.register("unfollow", middlewareLoggedIn(unfollow))

	return c
}

func parseCmd() command {
	if len(os.Args) < 2 {
		fmt.Println("Usage: gator <command> [args...]")
		os.Exit(1)
	}
	cmd := command{
		name: os.Args[1],
		args: os.Args[2:],
	}
	return cmd
}

func (c *commands) run(s *state, cmd command) error {
	if _, ok := c.handler[cmd.name]; !ok {
		return fmt.Errorf("no cmd found")
	}

	if len(os.Args) < 2 {
		return fmt.Errorf("No few than 2 arg")
	}

	if err := c.handler[cmd.name](s, cmd); err != nil {
		return err
	}

	return nil
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.handler[name] = f
}
