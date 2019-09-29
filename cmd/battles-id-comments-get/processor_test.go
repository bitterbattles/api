package main_test

import (
	"encoding/json"
	"fmt"
	"testing"

	. "github.com/bitterbattles/api/cmd/battles-id-comments-get"
	"github.com/bitterbattles/api/pkg/comments"
	commentsMocks "github.com/bitterbattles/api/pkg/comments/mocks"
	"github.com/bitterbattles/api/pkg/http"
	indexMocks "github.com/bitterbattles/api/pkg/index/mocks"
	"github.com/bitterbattles/api/pkg/lambda/api"
	. "github.com/bitterbattles/api/pkg/tests"
	"github.com/bitterbattles/api/pkg/users"
	usersMocks "github.com/bitterbattles/api/pkg/users/mocks"
)

const testBattleID = "battleId0"

func TestProcessorFullPage(t *testing.T) {
	indexRepository := indexMocks.NewRepository()
	commentsRepository := commentsMocks.NewRepository()
	usersRepository := usersMocks.NewRepository()
	addComments(indexRepository, commentsRepository, usersRepository, testBattleID, false, 3)
	expectedResponse := `[{"id":"id0","battleId":"battleId0","createdOn":0,"username":"username0","comment":"comment0"},{"id":"id1","battleId":"battleId0","createdOn":0,"username":"username1","comment":"comment1"}]`
	testProcessor(t, indexRepository, commentsRepository, usersRepository, testBattleID, "1", "2", expectedResponse)
}

func TestProcessorLastPage(t *testing.T) {
	indexRepository := indexMocks.NewRepository()
	commentsRepository := commentsMocks.NewRepository()
	usersRepository := usersMocks.NewRepository()
	addComments(indexRepository, commentsRepository, usersRepository, testBattleID, false, 3)
	expectedResponse := `[{"id":"id2","battleId":"battleId0","createdOn":0,"username":"username2","comment":"comment2"}]`
	testProcessor(t, indexRepository, commentsRepository, usersRepository, testBattleID, "2", "2", expectedResponse)
}

func TestProcessorBeyondLastPage(t *testing.T) {
	indexRepository := indexMocks.NewRepository()
	commentsRepository := commentsMocks.NewRepository()
	usersRepository := usersMocks.NewRepository()
	addComments(indexRepository, commentsRepository, usersRepository, testBattleID, false, 3)
	expectedResponse := `[]`
	testProcessor(t, indexRepository, commentsRepository, usersRepository, testBattleID, "3", "2", expectedResponse)
}

func TestProcessorNoPagination(t *testing.T) {
	indexRepository := indexMocks.NewRepository()
	commentsRepository := commentsMocks.NewRepository()
	usersRepository := usersMocks.NewRepository()
	addComments(indexRepository, commentsRepository, usersRepository, testBattleID, false, 50)
	expectedResponse := `[{"id":"id0","battleId":"battleId0","createdOn":0,"username":"username0","comment":"comment0"},{"id":"id1","battleId":"battleId0","createdOn":0,"username":"username1","comment":"comment1"},{"id":"id2","battleId":"battleId0","createdOn":0,"username":"username2","comment":"comment2"},{"id":"id3","battleId":"battleId0","createdOn":0,"username":"username3","comment":"comment3"},{"id":"id4","battleId":"battleId0","createdOn":0,"username":"username4","comment":"comment4"},{"id":"id5","battleId":"battleId0","createdOn":0,"username":"username5","comment":"comment5"},{"id":"id6","battleId":"battleId0","createdOn":0,"username":"username6","comment":"comment6"},{"id":"id7","battleId":"battleId0","createdOn":0,"username":"username7","comment":"comment7"},{"id":"id8","battleId":"battleId0","createdOn":0,"username":"username8","comment":"comment8"},{"id":"id9","battleId":"battleId0","createdOn":0,"username":"username9","comment":"comment9"},{"id":"id10","battleId":"battleId0","createdOn":0,"username":"username10","comment":"comment10"},{"id":"id11","battleId":"battleId0","createdOn":0,"username":"username11","comment":"comment11"},{"id":"id12","battleId":"battleId0","createdOn":0,"username":"username12","comment":"comment12"},{"id":"id13","battleId":"battleId0","createdOn":0,"username":"username13","comment":"comment13"},{"id":"id14","battleId":"battleId0","createdOn":0,"username":"username14","comment":"comment14"},{"id":"id15","battleId":"battleId0","createdOn":0,"username":"username15","comment":"comment15"},{"id":"id16","battleId":"battleId0","createdOn":0,"username":"username16","comment":"comment16"},{"id":"id17","battleId":"battleId0","createdOn":0,"username":"username17","comment":"comment17"},{"id":"id18","battleId":"battleId0","createdOn":0,"username":"username18","comment":"comment18"},{"id":"id19","battleId":"battleId0","createdOn":0,"username":"username19","comment":"comment19"},{"id":"id20","battleId":"battleId0","createdOn":0,"username":"username20","comment":"comment20"},{"id":"id21","battleId":"battleId0","createdOn":0,"username":"username21","comment":"comment21"},{"id":"id22","battleId":"battleId0","createdOn":0,"username":"username22","comment":"comment22"},{"id":"id23","battleId":"battleId0","createdOn":0,"username":"username23","comment":"comment23"},{"id":"id24","battleId":"battleId0","createdOn":0,"username":"username24","comment":"comment24"},{"id":"id25","battleId":"battleId0","createdOn":0,"username":"username25","comment":"comment25"},{"id":"id26","battleId":"battleId0","createdOn":0,"username":"username26","comment":"comment26"},{"id":"id27","battleId":"battleId0","createdOn":0,"username":"username27","comment":"comment27"},{"id":"id28","battleId":"battleId0","createdOn":0,"username":"username28","comment":"comment28"},{"id":"id29","battleId":"battleId0","createdOn":0,"username":"username29","comment":"comment29"},{"id":"id30","battleId":"battleId0","createdOn":0,"username":"username30","comment":"comment30"},{"id":"id31","battleId":"battleId0","createdOn":0,"username":"username31","comment":"comment31"},{"id":"id32","battleId":"battleId0","createdOn":0,"username":"username32","comment":"comment32"},{"id":"id33","battleId":"battleId0","createdOn":0,"username":"username33","comment":"comment33"},{"id":"id34","battleId":"battleId0","createdOn":0,"username":"username34","comment":"comment34"},{"id":"id35","battleId":"battleId0","createdOn":0,"username":"username35","comment":"comment35"},{"id":"id36","battleId":"battleId0","createdOn":0,"username":"username36","comment":"comment36"},{"id":"id37","battleId":"battleId0","createdOn":0,"username":"username37","comment":"comment37"},{"id":"id38","battleId":"battleId0","createdOn":0,"username":"username38","comment":"comment38"},{"id":"id39","battleId":"battleId0","createdOn":0,"username":"username39","comment":"comment39"},{"id":"id40","battleId":"battleId0","createdOn":0,"username":"username40","comment":"comment40"},{"id":"id41","battleId":"battleId0","createdOn":0,"username":"username41","comment":"comment41"},{"id":"id42","battleId":"battleId0","createdOn":0,"username":"username42","comment":"comment42"},{"id":"id43","battleId":"battleId0","createdOn":0,"username":"username43","comment":"comment43"},{"id":"id44","battleId":"battleId0","createdOn":0,"username":"username44","comment":"comment44"},{"id":"id45","battleId":"battleId0","createdOn":0,"username":"username45","comment":"comment45"},{"id":"id46","battleId":"battleId0","createdOn":0,"username":"username46","comment":"comment46"},{"id":"id47","battleId":"battleId0","createdOn":0,"username":"username47","comment":"comment47"},{"id":"id48","battleId":"battleId0","createdOn":0,"username":"username48","comment":"comment48"},{"id":"id49","battleId":"battleId0","createdOn":0,"username":"username49","comment":"comment49"}]`
	testProcessor(t, indexRepository, commentsRepository, usersRepository, testBattleID, "", "", expectedResponse)
}

func TestProcessorTooLargePage(t *testing.T) {
	indexRepository := indexMocks.NewRepository()
	commentsRepository := commentsMocks.NewRepository()
	usersRepository := usersMocks.NewRepository()
	addComments(indexRepository, commentsRepository, usersRepository, testBattleID, false, 101)
	expectedResponse := `[{"id":"id0","battleId":"battleId0","createdOn":0,"username":"username0","comment":"comment0"},{"id":"id1","battleId":"battleId0","createdOn":0,"username":"username1","comment":"comment1"},{"id":"id2","battleId":"battleId0","createdOn":0,"username":"username2","comment":"comment2"},{"id":"id3","battleId":"battleId0","createdOn":0,"username":"username3","comment":"comment3"},{"id":"id4","battleId":"battleId0","createdOn":0,"username":"username4","comment":"comment4"},{"id":"id5","battleId":"battleId0","createdOn":0,"username":"username5","comment":"comment5"},{"id":"id6","battleId":"battleId0","createdOn":0,"username":"username6","comment":"comment6"},{"id":"id7","battleId":"battleId0","createdOn":0,"username":"username7","comment":"comment7"},{"id":"id8","battleId":"battleId0","createdOn":0,"username":"username8","comment":"comment8"},{"id":"id9","battleId":"battleId0","createdOn":0,"username":"username9","comment":"comment9"},{"id":"id10","battleId":"battleId0","createdOn":0,"username":"username10","comment":"comment10"},{"id":"id11","battleId":"battleId0","createdOn":0,"username":"username11","comment":"comment11"},{"id":"id12","battleId":"battleId0","createdOn":0,"username":"username12","comment":"comment12"},{"id":"id13","battleId":"battleId0","createdOn":0,"username":"username13","comment":"comment13"},{"id":"id14","battleId":"battleId0","createdOn":0,"username":"username14","comment":"comment14"},{"id":"id15","battleId":"battleId0","createdOn":0,"username":"username15","comment":"comment15"},{"id":"id16","battleId":"battleId0","createdOn":0,"username":"username16","comment":"comment16"},{"id":"id17","battleId":"battleId0","createdOn":0,"username":"username17","comment":"comment17"},{"id":"id18","battleId":"battleId0","createdOn":0,"username":"username18","comment":"comment18"},{"id":"id19","battleId":"battleId0","createdOn":0,"username":"username19","comment":"comment19"},{"id":"id20","battleId":"battleId0","createdOn":0,"username":"username20","comment":"comment20"},{"id":"id21","battleId":"battleId0","createdOn":0,"username":"username21","comment":"comment21"},{"id":"id22","battleId":"battleId0","createdOn":0,"username":"username22","comment":"comment22"},{"id":"id23","battleId":"battleId0","createdOn":0,"username":"username23","comment":"comment23"},{"id":"id24","battleId":"battleId0","createdOn":0,"username":"username24","comment":"comment24"},{"id":"id25","battleId":"battleId0","createdOn":0,"username":"username25","comment":"comment25"},{"id":"id26","battleId":"battleId0","createdOn":0,"username":"username26","comment":"comment26"},{"id":"id27","battleId":"battleId0","createdOn":0,"username":"username27","comment":"comment27"},{"id":"id28","battleId":"battleId0","createdOn":0,"username":"username28","comment":"comment28"},{"id":"id29","battleId":"battleId0","createdOn":0,"username":"username29","comment":"comment29"},{"id":"id30","battleId":"battleId0","createdOn":0,"username":"username30","comment":"comment30"},{"id":"id31","battleId":"battleId0","createdOn":0,"username":"username31","comment":"comment31"},{"id":"id32","battleId":"battleId0","createdOn":0,"username":"username32","comment":"comment32"},{"id":"id33","battleId":"battleId0","createdOn":0,"username":"username33","comment":"comment33"},{"id":"id34","battleId":"battleId0","createdOn":0,"username":"username34","comment":"comment34"},{"id":"id35","battleId":"battleId0","createdOn":0,"username":"username35","comment":"comment35"},{"id":"id36","battleId":"battleId0","createdOn":0,"username":"username36","comment":"comment36"},{"id":"id37","battleId":"battleId0","createdOn":0,"username":"username37","comment":"comment37"},{"id":"id38","battleId":"battleId0","createdOn":0,"username":"username38","comment":"comment38"},{"id":"id39","battleId":"battleId0","createdOn":0,"username":"username39","comment":"comment39"},{"id":"id40","battleId":"battleId0","createdOn":0,"username":"username40","comment":"comment40"},{"id":"id41","battleId":"battleId0","createdOn":0,"username":"username41","comment":"comment41"},{"id":"id42","battleId":"battleId0","createdOn":0,"username":"username42","comment":"comment42"},{"id":"id43","battleId":"battleId0","createdOn":0,"username":"username43","comment":"comment43"},{"id":"id44","battleId":"battleId0","createdOn":0,"username":"username44","comment":"comment44"},{"id":"id45","battleId":"battleId0","createdOn":0,"username":"username45","comment":"comment45"},{"id":"id46","battleId":"battleId0","createdOn":0,"username":"username46","comment":"comment46"},{"id":"id47","battleId":"battleId0","createdOn":0,"username":"username47","comment":"comment47"},{"id":"id48","battleId":"battleId0","createdOn":0,"username":"username48","comment":"comment48"},{"id":"id49","battleId":"battleId0","createdOn":0,"username":"username49","comment":"comment49"},{"id":"id50","battleId":"battleId0","createdOn":0,"username":"username50","comment":"comment50"},{"id":"id51","battleId":"battleId0","createdOn":0,"username":"username51","comment":"comment51"},{"id":"id52","battleId":"battleId0","createdOn":0,"username":"username52","comment":"comment52"},{"id":"id53","battleId":"battleId0","createdOn":0,"username":"username53","comment":"comment53"},{"id":"id54","battleId":"battleId0","createdOn":0,"username":"username54","comment":"comment54"},{"id":"id55","battleId":"battleId0","createdOn":0,"username":"username55","comment":"comment55"},{"id":"id56","battleId":"battleId0","createdOn":0,"username":"username56","comment":"comment56"},{"id":"id57","battleId":"battleId0","createdOn":0,"username":"username57","comment":"comment57"},{"id":"id58","battleId":"battleId0","createdOn":0,"username":"username58","comment":"comment58"},{"id":"id59","battleId":"battleId0","createdOn":0,"username":"username59","comment":"comment59"},{"id":"id60","battleId":"battleId0","createdOn":0,"username":"username60","comment":"comment60"},{"id":"id61","battleId":"battleId0","createdOn":0,"username":"username61","comment":"comment61"},{"id":"id62","battleId":"battleId0","createdOn":0,"username":"username62","comment":"comment62"},{"id":"id63","battleId":"battleId0","createdOn":0,"username":"username63","comment":"comment63"},{"id":"id64","battleId":"battleId0","createdOn":0,"username":"username64","comment":"comment64"},{"id":"id65","battleId":"battleId0","createdOn":0,"username":"username65","comment":"comment65"},{"id":"id66","battleId":"battleId0","createdOn":0,"username":"username66","comment":"comment66"},{"id":"id67","battleId":"battleId0","createdOn":0,"username":"username67","comment":"comment67"},{"id":"id68","battleId":"battleId0","createdOn":0,"username":"username68","comment":"comment68"},{"id":"id69","battleId":"battleId0","createdOn":0,"username":"username69","comment":"comment69"},{"id":"id70","battleId":"battleId0","createdOn":0,"username":"username70","comment":"comment70"},{"id":"id71","battleId":"battleId0","createdOn":0,"username":"username71","comment":"comment71"},{"id":"id72","battleId":"battleId0","createdOn":0,"username":"username72","comment":"comment72"},{"id":"id73","battleId":"battleId0","createdOn":0,"username":"username73","comment":"comment73"},{"id":"id74","battleId":"battleId0","createdOn":0,"username":"username74","comment":"comment74"},{"id":"id75","battleId":"battleId0","createdOn":0,"username":"username75","comment":"comment75"},{"id":"id76","battleId":"battleId0","createdOn":0,"username":"username76","comment":"comment76"},{"id":"id77","battleId":"battleId0","createdOn":0,"username":"username77","comment":"comment77"},{"id":"id78","battleId":"battleId0","createdOn":0,"username":"username78","comment":"comment78"},{"id":"id79","battleId":"battleId0","createdOn":0,"username":"username79","comment":"comment79"},{"id":"id80","battleId":"battleId0","createdOn":0,"username":"username80","comment":"comment80"},{"id":"id81","battleId":"battleId0","createdOn":0,"username":"username81","comment":"comment81"},{"id":"id82","battleId":"battleId0","createdOn":0,"username":"username82","comment":"comment82"},{"id":"id83","battleId":"battleId0","createdOn":0,"username":"username83","comment":"comment83"},{"id":"id84","battleId":"battleId0","createdOn":0,"username":"username84","comment":"comment84"},{"id":"id85","battleId":"battleId0","createdOn":0,"username":"username85","comment":"comment85"},{"id":"id86","battleId":"battleId0","createdOn":0,"username":"username86","comment":"comment86"},{"id":"id87","battleId":"battleId0","createdOn":0,"username":"username87","comment":"comment87"},{"id":"id88","battleId":"battleId0","createdOn":0,"username":"username88","comment":"comment88"},{"id":"id89","battleId":"battleId0","createdOn":0,"username":"username89","comment":"comment89"},{"id":"id90","battleId":"battleId0","createdOn":0,"username":"username90","comment":"comment90"},{"id":"id91","battleId":"battleId0","createdOn":0,"username":"username91","comment":"comment91"},{"id":"id92","battleId":"battleId0","createdOn":0,"username":"username92","comment":"comment92"},{"id":"id93","battleId":"battleId0","createdOn":0,"username":"username93","comment":"comment93"},{"id":"id94","battleId":"battleId0","createdOn":0,"username":"username94","comment":"comment94"},{"id":"id95","battleId":"battleId0","createdOn":0,"username":"username95","comment":"comment95"},{"id":"id96","battleId":"battleId0","createdOn":0,"username":"username96","comment":"comment96"},{"id":"id97","battleId":"battleId0","createdOn":0,"username":"username97","comment":"comment97"},{"id":"id98","battleId":"battleId0","createdOn":0,"username":"username98","comment":"comment98"},{"id":"id99","battleId":"battleId0","createdOn":0,"username":"username99","comment":"comment99"}]`
	testProcessor(t, indexRepository, commentsRepository, usersRepository, testBattleID, "1", "101", expectedResponse)
}

func TestProcessorBadIndexEntry(t *testing.T) {
	indexRepository := indexMocks.NewRepository()
	key := fmt.Sprintf("commentIds:forBattle:%s", testBattleID)
	indexRepository.SetScore(key, "badId", 0)
	commentsRepository := commentsMocks.NewRepository()
	usersRepository := usersMocks.NewRepository()
	addComments(indexRepository, commentsRepository, usersRepository, testBattleID, false, 1)
	expectedResponse := `[{"id":"id0","battleId":"battleId0","createdOn":0,"username":"username0","comment":"comment0"}]`
	testProcessor(t, indexRepository, commentsRepository, usersRepository, testBattleID, "1", "2", expectedResponse)
}

func TestProcessorDeletedIndexEntry(t *testing.T) {
	indexRepository := indexMocks.NewRepository()
	commentsRepository := commentsMocks.NewRepository()
	usersRepository := usersMocks.NewRepository()
	addComments(indexRepository, commentsRepository, usersRepository, testBattleID, true, 1)
	expectedResponse := `[{"id":"id0","battleId":"battleId0","createdOn":0,"username":"username0","comment":"[Deleted]"}]`
	testProcessor(t, indexRepository, commentsRepository, usersRepository, testBattleID, "1", "2", expectedResponse)
}

func addComments(indexRepository *indexMocks.Repository, commentsRepository *commentsMocks.Repository, usersRepository *usersMocks.Repository, battleID string, isDeleted bool, count int) {
	key := fmt.Sprintf("commentIds:forBattle:%s", battleID)
	state := comments.Active
	if isDeleted {
		state = comments.Deleted
	}
	for i := 0; i < count; i++ {
		userID := fmt.Sprintf("userId%d", i)
		comment := &comments.Comment{
			ID:       fmt.Sprintf("id%d", i),
			BattleID: battleID,
			UserID:   userID,
			Comment:  fmt.Sprintf("comment%d", i),
			State:    state,
		}
		commentsRepository.Add(comment)
		indexRepository.SetScore(key, comment.ID, float64(i))
		user := &users.User{
			ID:       userID,
			Username: fmt.Sprintf("username%d", i),
		}
		usersRepository.Add(user)
	}
}

func testProcessor(t *testing.T, indexRepository *indexMocks.Repository, commentsRepository *commentsMocks.Repository, usersRepository *usersMocks.Repository, battleID string, page string, pageSize string, expectedResponseBody string) {
	pathParams := make(map[string]string)
	pathParams["id"] = battleID
	queryParams := make(map[string]string)
	if page != "" {
		queryParams["page"] = page
	}
	if pageSize != "" {
		queryParams["pageSize"] = pageSize
	}
	input := &api.Input{
		PathParams:  pathParams,
		QueryParams: queryParams,
	}
	indexer := comments.NewIndexer(indexRepository)
	processor := NewProcessor(indexer, commentsRepository, usersRepository)
	output, err := processor.Process(input)
	AssertNil(t, err)
	AssertNotNil(t, output)
	AssertEquals(t, output.StatusCode, http.OK)
	AssertNotNil(t, output.ResponseBody)
	responseBody, _ := json.Marshal(output.ResponseBody)
	AssertEquals(t, string(responseBody), expectedResponseBody)
}
