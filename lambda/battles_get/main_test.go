package main

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/bitterbattles/api/battles"
	"github.com/bitterbattles/api/mocks"
)

const sort = "recent"

var controller *battles.Controller
var manager *battles.Manager
var table *mocks.MockTable
var index *mocks.MockIndex

func TestMain(m *testing.M) {
	index = mocks.NewMockIndex()
	table = mocks.NewMockTable()
	manager = battles.NewManager(index, table)
	controller = battles.NewController(manager)
	os.Exit(m.Run())
}

func TestGetFullPage(t *testing.T) {
	reset()
	addBattles(sort, 3)
	response := get(t, sort, 1, 2)
	expected := `[{"id":"id0","title":"title0","description":"description0","votesFor":0,"votesAgainst":0,"createdOn":0},{"id":"id1","title":"title1","description":"description1","votesFor":1,"votesAgainst":2,"createdOn":3}]`
	verify(t, toJSON(response), expected)
}

func TestGetLastPage(t *testing.T) {
	reset()
	addBattles(sort, 3)
	response := get(t, sort, 2, 2)
	expected := `[{"id":"id2","title":"title2","description":"description2","votesFor":2,"votesAgainst":4,"createdOn":6}]`
	verify(t, toJSON(response), expected)
}

func TestGetBeyondLastPage(t *testing.T) {
	reset()
	addBattles(sort, 3)
	response := get(t, sort, 3, 2)
	expected := `[]`
	verify(t, toJSON(response), expected)
}

func TestGetNoPagination(t *testing.T) {
	reset()
	addBattles(sort, 50)
	response := get(t, sort, 0, 0)
	verify(t, len(response), 50)
}

func TestGetTooLargePage(t *testing.T) {
	reset()
	addBattles(sort, 100)
	response := get(t, sort, 0, 101)
	verify(t, len(response), 100)
}

func TestGetNoSort(t *testing.T) {
	reset()
	addBattles(sort, 3)
	response := get(t, sort, 1, 1)
	expected := `[{"id":"id0","title":"title0","description":"description0","votesFor":0,"votesAgainst":0,"createdOn":0}]`
	verify(t, toJSON(response), expected)
}

func TestGetBadIndex(t *testing.T) {
	reset()
	index.Add(sort, "bad")
	addBattles(sort, 3)
	response := get(t, sort, 1, 2)
	expected := `[{"id":"id0","title":"title0","description":"description0","votesFor":0,"votesAgainst":0,"createdOn":0}]`
	verify(t, toJSON(response), expected)
}

func addBattles(indexName string, count int) {
	for i := 0; i < count; i++ {
		battle := battles.Battle{
			ID:           fmt.Sprintf("id%d", i),
			UserID:       fmt.Sprintf("userId%d", i),
			Title:        fmt.Sprintf("title%d", i),
			Description:  fmt.Sprintf("description%d", i),
			VotesFor:     i,
			VotesAgainst: i * 2,
			CreatedOn:    uint64(i * 3),
		}
		table.Add(battle)
		index.Add(indexName, battle.ID)
	}
}

func reset() {
	index.Reset()
	table.Reset()
}

func get(t *testing.T, sort string, page int, pageSize int) []battles.GetResponse {
	request := battles.GetRequest{
		Sort:     sort,
		Page:     page,
		PageSize: pageSize,
	}
	response, err := controller.HandleGet(request)
	if err != nil {
		t.Error("Unexected error:", err.Error())
	}
	return response
}

func toJSON(obj interface{}) string {
	bytes, _ := json.Marshal(obj)
	return string(bytes)
}

func verify(t *testing.T, actual interface{}, expected interface{}) {
	if actual != expected {
		t.Error("Unexpected value.\nExpected:", expected, "\nActual:", actual)
	}
}
