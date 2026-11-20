package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/Cyrus-Gahatraj/boot.gator/internal/config"
	"github.com/Cyrus-Gahatraj/boot.gator/internal/database"
	_ "github.com/lib/pq"
)

type state struct {
	cfg *config.Config
	db  *database.Queries
}

func main() {
	s := state{}
	cfg := config.Read()

	db, err := sql.Open("postgres", cfg.DBURL)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		os.Exit(1)
	}

	s.cfg = &cfg
	s.db = database.New(db)
	cmd := parseCmd()

	c := newCommands()
	if err := c.run(&s, cmd); err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
}
