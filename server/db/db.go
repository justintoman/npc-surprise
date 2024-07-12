package db

import (
	"fmt"

	"github.com/supabase-community/supabase-go"
)

func New(url string, key string) Db {
	client, err := supabase.NewClient(url, key, nil)
	if err != nil {
		fmt.Println("cannot initalize client", err)
	}

	db := Db{
		Character: CharacterTable{client: client},
		Action:    ActionTable{client: client},
		Player:    PlayerTable{client: client},
	}
	return db
}

type Db struct {
	Character CharacterTable
	Action    ActionTable
	Player    PlayerTable
}
