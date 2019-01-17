package battles_test

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/bitterbattles/api/battles"
	"github.com/bitterbattles/api/battles/mocks"
	"github.com/bitterbattles/api/core/errors"
	"github.com/bitterbattles/api/core/tests"
)

const sort = "recent"

var index *mocks.Index
var repository *mocks.Repository
var manager *battles.Manager
var controller *battles.Controller

func TestMain(m *testing.M) {
	index = mocks.NewIndex()
	repository = mocks.NewRepository()
	manager = battles.NewManager(index, repository)
	controller = battles.NewController(manager)
	os.Exit(m.Run())
}

func TestGetFullPage(t *testing.T) {
	reset()
	addBattles(sort, 3)
	response, err := get(sort, 1, 2)
	tests.AssertNil(t, err)
	expected := `[{"id":"id0","title":"title0","description":"description0","votesFor":0,"votesAgainst":0,"createdOn":0},{"id":"id1","title":"title1","description":"description1","votesFor":1,"votesAgainst":2,"createdOn":3}]`
	tests.AssertEquals(t, toJSON(response), expected)
}

func TestGetLastPage(t *testing.T) {
	reset()
	addBattles(sort, 3)
	response, err := get(sort, 2, 2)
	tests.AssertNil(t, err)
	expected := `[{"id":"id2","title":"title2","description":"description2","votesFor":2,"votesAgainst":4,"createdOn":6}]`
	tests.AssertEquals(t, toJSON(response), expected)
}

func TestGetBeyondLastPage(t *testing.T) {
	reset()
	addBattles(sort, 3)
	response, err := get(sort, 3, 2)
	tests.AssertNil(t, err)
	expected := `[]`
	tests.AssertEquals(t, toJSON(response), expected)
}

func TestGetNoPagination(t *testing.T) {
	reset()
	addBattles(sort, 50)
	response, err := get(sort, 0, 0)
	tests.AssertNil(t, err)
	tests.AssertEquals(t, len(response), 50)
}

func TestGetTooLargePage(t *testing.T) {
	reset()
	addBattles(sort, 100)
	response, err := get(sort, 0, 101)
	tests.AssertNil(t, err)
	tests.AssertEquals(t, len(response), 100)
}

func TestGetNoSort(t *testing.T) {
	reset()
	addBattles(sort, 3)
	response, err := get(sort, 1, 1)
	tests.AssertNil(t, err)
	expected := `[{"id":"id0","title":"title0","description":"description0","votesFor":0,"votesAgainst":0,"createdOn":0}]`
	tests.AssertEquals(t, toJSON(response), expected)
}

func TestGetBadSort(t *testing.T) {
	reset()
	index.Add(sort, "bad")
	addBattles(sort, 3)
	response, err := get(sort, 1, 2)
	tests.AssertNil(t, err)
	expected := `[{"id":"id0","title":"title0","description":"description0","votesFor":0,"votesAgainst":0,"createdOn":0}]`
	tests.AssertEquals(t, toJSON(response), expected)
}

func TestPost(t *testing.T) {
	reset()
	title := "title"
	description := "description"
	err := post(title, description)
	tests.AssertNil(t, err)
	battle := repository.GetLastAdded()
	tests.AssertNotNil(t, battle)
	tests.AssertEquals(t, battle.Title, title)
	tests.AssertEquals(t, battle.Description, description)
}

func TestPostTooShortTitle(t *testing.T) {
	reset()
	err := post("", "description")
	tests.AssertHTTPError(t, err, errors.BadRequestCode)
	tests.AssertNil(t, repository.GetLastAdded())
}

func TestPostTooLongTitle(t *testing.T) {
	reset()
	err := post("loooooooooooooooooooooooooooooooooooooooooooooooooongtitle", "description")
	tests.AssertHTTPError(t, err, errors.BadRequestCode)
	tests.AssertNil(t, repository.GetLastAdded())
}

func TestPostTooShortDescription(t *testing.T) {
	reset()
	err := post("title", "")
	tests.AssertHTTPError(t, err, errors.BadRequestCode)
	tests.AssertNil(t, repository.GetLastAdded())
}

func TestPostTooLongDescription(t *testing.T) {
	reset()
	err := post("title", "loooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooongdescription")
	tests.AssertHTTPError(t, err, errors.BadRequestCode)
	tests.AssertNil(t, repository.GetLastAdded())
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
			CreatedOn:    int64(i * 3),
		}
		repository.Add(battle)
		index.Add(indexName, battle.ID)
	}
}

func reset() {
	index.Reset()
	repository.Reset()
}

func get(sort string, page int, pageSize int) ([]battles.GetResponse, error) {
	request := battles.GetRequest{
		Sort:     sort,
		Page:     page,
		PageSize: pageSize,
	}
	return controller.HandleGet(request)
}

func post(title string, description string) error {
	request := battles.PostRequest{
		Title:       title,
		Description: description,
	}
	return controller.HandlePost(request)
}

func toJSON(obj interface{}) string {
	bytes, _ := json.Marshal(obj)
	return string(bytes)
}
