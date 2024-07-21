package db

import (
	"fmt"
	"strconv"

	"github.com/supabase-community/postgrest-go"
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

func filterById(filterBuilder *postgrest.FilterBuilder, id int) *postgrest.FilterBuilder {
	return filterBuilder.Filter("id", "eq", strconv.Itoa(id)).Single()
}

func filterByCharacterId(filterBuilder *postgrest.FilterBuilder, characterId int) *postgrest.FilterBuilder {
	return filterBuilder.Filter("characterId", "eq", strconv.Itoa(characterId))
}

func filterByPlayerId(filterBuilder *postgrest.FilterBuilder, playerId int) *postgrest.FilterBuilder {
	return filterBuilder.Filter("playerId", "eq", strconv.Itoa(playerId))
}

func selectAll(queryBuilder *postgrest.QueryBuilder) *postgrest.FilterBuilder {
	return queryBuilder.Select("*", "exact", false)
}

func insertSingle(queryBuilder *postgrest.QueryBuilder, payload interface{}) *postgrest.FilterBuilder {
	return queryBuilder.Insert(payload, true, "", "", "exact").Single()
}

func deleteSingle(filterBuilder *postgrest.QueryBuilder) *postgrest.FilterBuilder {
	return filterBuilder.Delete("", "").Single()
}

func orderById(queryBuilder *postgrest.FilterBuilder) *postgrest.FilterBuilder {
	return queryBuilder.Order("id", &postgrest.OrderOpts{
		Ascending: true,
	})
}
