package main_test

import (
	"encoding/json"
	"fmt"
	"testing"

	. "github.com/bitterbattles/api/cmd/comments-me-get"
	"github.com/bitterbattles/api/pkg/comments"
	commentsMocks "github.com/bitterbattles/api/pkg/comments/mocks"
	"github.com/bitterbattles/api/pkg/http"
	indexMocks "github.com/bitterbattles/api/pkg/index/mocks"
	"github.com/bitterbattles/api/pkg/lambda/api"
	. "github.com/bitterbattles/api/pkg/tests"
	"github.com/bitterbattles/api/pkg/users"
	usersMocks "github.com/bitterbattles/api/pkg/users/mocks"
)

const testUserID = "userId0"

func TestProcessorFullPage(t *testing.T) {
	indexRepository := indexMocks.NewRepository()
	commentsRepository := commentsMocks.NewRepository()
	addComments(indexRepository, commentsRepository, testUserID, false, 3)
	expectedResponse := `[{"id":"id0","battleId":"battleId0","createdOn":0,"username":"username0","comment":"comment0"},{"id":"id1","battleId":"battleId1","createdOn":0,"username":"username0","comment":"comment1"}]`
	testProcessor(t, indexRepository, commentsRepository, "1", "2", expectedResponse)
}

func TestProcessorLastPage(t *testing.T) {
	indexRepository := indexMocks.NewRepository()
	commentsRepository := commentsMocks.NewRepository()
	addComments(indexRepository, commentsRepository, testUserID, false, 3)
	expectedResponse := `[{"id":"id2","battleId":"battleId2","createdOn":0,"username":"username0","comment":"comment2"}]`
	testProcessor(t, indexRepository, commentsRepository, "2", "2", expectedResponse)
}

func TestProcessorBeyondLastPage(t *testing.T) {
	indexRepository := indexMocks.NewRepository()
	commentsRepository := commentsMocks.NewRepository()
	addComments(indexRepository, commentsRepository, testUserID, false, 3)
	expectedResponse := `[]`
	testProcessor(t, indexRepository, commentsRepository, "3", "2", expectedResponse)
}

func TestProcessorNoPagination(t *testing.T) {
	indexRepository := indexMocks.NewRepository()
	commentsRepository := commentsMocks.NewRepository()
	addComments(indexRepository, commentsRepository, testUserID, false, 50)
	expectedResponse := `[{"id":"id0","battleId":"battleId0","createdOn":0,"username":"username0","comment":"comment0"},{"id":"id1","battleId":"battleId1","createdOn":0,"username":"username0","comment":"comment1"},{"id":"id2","battleId":"battleId2","createdOn":0,"username":"username0","comment":"comment2"},{"id":"id3","battleId":"battleId3","createdOn":0,"username":"username0","comment":"comment3"},{"id":"id4","battleId":"battleId4","createdOn":0,"username":"username0","comment":"comment4"},{"id":"id5","battleId":"battleId5","createdOn":0,"username":"username0","comment":"comment5"},{"id":"id6","battleId":"battleId6","createdOn":0,"username":"username0","comment":"comment6"},{"id":"id7","battleId":"battleId7","createdOn":0,"username":"username0","comment":"comment7"},{"id":"id8","battleId":"battleId8","createdOn":0,"username":"username0","comment":"comment8"},{"id":"id9","battleId":"battleId9","createdOn":0,"username":"username0","comment":"comment9"},{"id":"id10","battleId":"battleId10","createdOn":0,"username":"username0","comment":"comment10"},{"id":"id11","battleId":"battleId11","createdOn":0,"username":"username0","comment":"comment11"},{"id":"id12","battleId":"battleId12","createdOn":0,"username":"username0","comment":"comment12"},{"id":"id13","battleId":"battleId13","createdOn":0,"username":"username0","comment":"comment13"},{"id":"id14","battleId":"battleId14","createdOn":0,"username":"username0","comment":"comment14"},{"id":"id15","battleId":"battleId15","createdOn":0,"username":"username0","comment":"comment15"},{"id":"id16","battleId":"battleId16","createdOn":0,"username":"username0","comment":"comment16"},{"id":"id17","battleId":"battleId17","createdOn":0,"username":"username0","comment":"comment17"},{"id":"id18","battleId":"battleId18","createdOn":0,"username":"username0","comment":"comment18"},{"id":"id19","battleId":"battleId19","createdOn":0,"username":"username0","comment":"comment19"},{"id":"id20","battleId":"battleId20","createdOn":0,"username":"username0","comment":"comment20"},{"id":"id21","battleId":"battleId21","createdOn":0,"username":"username0","comment":"comment21"},{"id":"id22","battleId":"battleId22","createdOn":0,"username":"username0","comment":"comment22"},{"id":"id23","battleId":"battleId23","createdOn":0,"username":"username0","comment":"comment23"},{"id":"id24","battleId":"battleId24","createdOn":0,"username":"username0","comment":"comment24"},{"id":"id25","battleId":"battleId25","createdOn":0,"username":"username0","comment":"comment25"},{"id":"id26","battleId":"battleId26","createdOn":0,"username":"username0","comment":"comment26"},{"id":"id27","battleId":"battleId27","createdOn":0,"username":"username0","comment":"comment27"},{"id":"id28","battleId":"battleId28","createdOn":0,"username":"username0","comment":"comment28"},{"id":"id29","battleId":"battleId29","createdOn":0,"username":"username0","comment":"comment29"},{"id":"id30","battleId":"battleId30","createdOn":0,"username":"username0","comment":"comment30"},{"id":"id31","battleId":"battleId31","createdOn":0,"username":"username0","comment":"comment31"},{"id":"id32","battleId":"battleId32","createdOn":0,"username":"username0","comment":"comment32"},{"id":"id33","battleId":"battleId33","createdOn":0,"username":"username0","comment":"comment33"},{"id":"id34","battleId":"battleId34","createdOn":0,"username":"username0","comment":"comment34"},{"id":"id35","battleId":"battleId35","createdOn":0,"username":"username0","comment":"comment35"},{"id":"id36","battleId":"battleId36","createdOn":0,"username":"username0","comment":"comment36"},{"id":"id37","battleId":"battleId37","createdOn":0,"username":"username0","comment":"comment37"},{"id":"id38","battleId":"battleId38","createdOn":0,"username":"username0","comment":"comment38"},{"id":"id39","battleId":"battleId39","createdOn":0,"username":"username0","comment":"comment39"},{"id":"id40","battleId":"battleId40","createdOn":0,"username":"username0","comment":"comment40"},{"id":"id41","battleId":"battleId41","createdOn":0,"username":"username0","comment":"comment41"},{"id":"id42","battleId":"battleId42","createdOn":0,"username":"username0","comment":"comment42"},{"id":"id43","battleId":"battleId43","createdOn":0,"username":"username0","comment":"comment43"},{"id":"id44","battleId":"battleId44","createdOn":0,"username":"username0","comment":"comment44"},{"id":"id45","battleId":"battleId45","createdOn":0,"username":"username0","comment":"comment45"},{"id":"id46","battleId":"battleId46","createdOn":0,"username":"username0","comment":"comment46"},{"id":"id47","battleId":"battleId47","createdOn":0,"username":"username0","comment":"comment47"},{"id":"id48","battleId":"battleId48","createdOn":0,"username":"username0","comment":"comment48"},{"id":"id49","battleId":"battleId49","createdOn":0,"username":"username0","comment":"comment49"}]`
	testProcessor(t, indexRepository, commentsRepository, "", "", expectedResponse)
}

func TestProcessorTooLargePage(t *testing.T) {
	indexRepository := indexMocks.NewRepository()
	commentsRepository := commentsMocks.NewRepository()
	addComments(indexRepository, commentsRepository, testUserID, false, 101)
	expectedResponse := `[{"id":"id0","battleId":"battleId0","createdOn":0,"username":"username0","comment":"comment0"},{"id":"id1","battleId":"battleId1","createdOn":0,"username":"username0","comment":"comment1"},{"id":"id2","battleId":"battleId2","createdOn":0,"username":"username0","comment":"comment2"},{"id":"id3","battleId":"battleId3","createdOn":0,"username":"username0","comment":"comment3"},{"id":"id4","battleId":"battleId4","createdOn":0,"username":"username0","comment":"comment4"},{"id":"id5","battleId":"battleId5","createdOn":0,"username":"username0","comment":"comment5"},{"id":"id6","battleId":"battleId6","createdOn":0,"username":"username0","comment":"comment6"},{"id":"id7","battleId":"battleId7","createdOn":0,"username":"username0","comment":"comment7"},{"id":"id8","battleId":"battleId8","createdOn":0,"username":"username0","comment":"comment8"},{"id":"id9","battleId":"battleId9","createdOn":0,"username":"username0","comment":"comment9"},{"id":"id10","battleId":"battleId10","createdOn":0,"username":"username0","comment":"comment10"},{"id":"id11","battleId":"battleId11","createdOn":0,"username":"username0","comment":"comment11"},{"id":"id12","battleId":"battleId12","createdOn":0,"username":"username0","comment":"comment12"},{"id":"id13","battleId":"battleId13","createdOn":0,"username":"username0","comment":"comment13"},{"id":"id14","battleId":"battleId14","createdOn":0,"username":"username0","comment":"comment14"},{"id":"id15","battleId":"battleId15","createdOn":0,"username":"username0","comment":"comment15"},{"id":"id16","battleId":"battleId16","createdOn":0,"username":"username0","comment":"comment16"},{"id":"id17","battleId":"battleId17","createdOn":0,"username":"username0","comment":"comment17"},{"id":"id18","battleId":"battleId18","createdOn":0,"username":"username0","comment":"comment18"},{"id":"id19","battleId":"battleId19","createdOn":0,"username":"username0","comment":"comment19"},{"id":"id20","battleId":"battleId20","createdOn":0,"username":"username0","comment":"comment20"},{"id":"id21","battleId":"battleId21","createdOn":0,"username":"username0","comment":"comment21"},{"id":"id22","battleId":"battleId22","createdOn":0,"username":"username0","comment":"comment22"},{"id":"id23","battleId":"battleId23","createdOn":0,"username":"username0","comment":"comment23"},{"id":"id24","battleId":"battleId24","createdOn":0,"username":"username0","comment":"comment24"},{"id":"id25","battleId":"battleId25","createdOn":0,"username":"username0","comment":"comment25"},{"id":"id26","battleId":"battleId26","createdOn":0,"username":"username0","comment":"comment26"},{"id":"id27","battleId":"battleId27","createdOn":0,"username":"username0","comment":"comment27"},{"id":"id28","battleId":"battleId28","createdOn":0,"username":"username0","comment":"comment28"},{"id":"id29","battleId":"battleId29","createdOn":0,"username":"username0","comment":"comment29"},{"id":"id30","battleId":"battleId30","createdOn":0,"username":"username0","comment":"comment30"},{"id":"id31","battleId":"battleId31","createdOn":0,"username":"username0","comment":"comment31"},{"id":"id32","battleId":"battleId32","createdOn":0,"username":"username0","comment":"comment32"},{"id":"id33","battleId":"battleId33","createdOn":0,"username":"username0","comment":"comment33"},{"id":"id34","battleId":"battleId34","createdOn":0,"username":"username0","comment":"comment34"},{"id":"id35","battleId":"battleId35","createdOn":0,"username":"username0","comment":"comment35"},{"id":"id36","battleId":"battleId36","createdOn":0,"username":"username0","comment":"comment36"},{"id":"id37","battleId":"battleId37","createdOn":0,"username":"username0","comment":"comment37"},{"id":"id38","battleId":"battleId38","createdOn":0,"username":"username0","comment":"comment38"},{"id":"id39","battleId":"battleId39","createdOn":0,"username":"username0","comment":"comment39"},{"id":"id40","battleId":"battleId40","createdOn":0,"username":"username0","comment":"comment40"},{"id":"id41","battleId":"battleId41","createdOn":0,"username":"username0","comment":"comment41"},{"id":"id42","battleId":"battleId42","createdOn":0,"username":"username0","comment":"comment42"},{"id":"id43","battleId":"battleId43","createdOn":0,"username":"username0","comment":"comment43"},{"id":"id44","battleId":"battleId44","createdOn":0,"username":"username0","comment":"comment44"},{"id":"id45","battleId":"battleId45","createdOn":0,"username":"username0","comment":"comment45"},{"id":"id46","battleId":"battleId46","createdOn":0,"username":"username0","comment":"comment46"},{"id":"id47","battleId":"battleId47","createdOn":0,"username":"username0","comment":"comment47"},{"id":"id48","battleId":"battleId48","createdOn":0,"username":"username0","comment":"comment48"},{"id":"id49","battleId":"battleId49","createdOn":0,"username":"username0","comment":"comment49"},{"id":"id50","battleId":"battleId50","createdOn":0,"username":"username0","comment":"comment50"},{"id":"id51","battleId":"battleId51","createdOn":0,"username":"username0","comment":"comment51"},{"id":"id52","battleId":"battleId52","createdOn":0,"username":"username0","comment":"comment52"},{"id":"id53","battleId":"battleId53","createdOn":0,"username":"username0","comment":"comment53"},{"id":"id54","battleId":"battleId54","createdOn":0,"username":"username0","comment":"comment54"},{"id":"id55","battleId":"battleId55","createdOn":0,"username":"username0","comment":"comment55"},{"id":"id56","battleId":"battleId56","createdOn":0,"username":"username0","comment":"comment56"},{"id":"id57","battleId":"battleId57","createdOn":0,"username":"username0","comment":"comment57"},{"id":"id58","battleId":"battleId58","createdOn":0,"username":"username0","comment":"comment58"},{"id":"id59","battleId":"battleId59","createdOn":0,"username":"username0","comment":"comment59"},{"id":"id60","battleId":"battleId60","createdOn":0,"username":"username0","comment":"comment60"},{"id":"id61","battleId":"battleId61","createdOn":0,"username":"username0","comment":"comment61"},{"id":"id62","battleId":"battleId62","createdOn":0,"username":"username0","comment":"comment62"},{"id":"id63","battleId":"battleId63","createdOn":0,"username":"username0","comment":"comment63"},{"id":"id64","battleId":"battleId64","createdOn":0,"username":"username0","comment":"comment64"},{"id":"id65","battleId":"battleId65","createdOn":0,"username":"username0","comment":"comment65"},{"id":"id66","battleId":"battleId66","createdOn":0,"username":"username0","comment":"comment66"},{"id":"id67","battleId":"battleId67","createdOn":0,"username":"username0","comment":"comment67"},{"id":"id68","battleId":"battleId68","createdOn":0,"username":"username0","comment":"comment68"},{"id":"id69","battleId":"battleId69","createdOn":0,"username":"username0","comment":"comment69"},{"id":"id70","battleId":"battleId70","createdOn":0,"username":"username0","comment":"comment70"},{"id":"id71","battleId":"battleId71","createdOn":0,"username":"username0","comment":"comment71"},{"id":"id72","battleId":"battleId72","createdOn":0,"username":"username0","comment":"comment72"},{"id":"id73","battleId":"battleId73","createdOn":0,"username":"username0","comment":"comment73"},{"id":"id74","battleId":"battleId74","createdOn":0,"username":"username0","comment":"comment74"},{"id":"id75","battleId":"battleId75","createdOn":0,"username":"username0","comment":"comment75"},{"id":"id76","battleId":"battleId76","createdOn":0,"username":"username0","comment":"comment76"},{"id":"id77","battleId":"battleId77","createdOn":0,"username":"username0","comment":"comment77"},{"id":"id78","battleId":"battleId78","createdOn":0,"username":"username0","comment":"comment78"},{"id":"id79","battleId":"battleId79","createdOn":0,"username":"username0","comment":"comment79"},{"id":"id80","battleId":"battleId80","createdOn":0,"username":"username0","comment":"comment80"},{"id":"id81","battleId":"battleId81","createdOn":0,"username":"username0","comment":"comment81"},{"id":"id82","battleId":"battleId82","createdOn":0,"username":"username0","comment":"comment82"},{"id":"id83","battleId":"battleId83","createdOn":0,"username":"username0","comment":"comment83"},{"id":"id84","battleId":"battleId84","createdOn":0,"username":"username0","comment":"comment84"},{"id":"id85","battleId":"battleId85","createdOn":0,"username":"username0","comment":"comment85"},{"id":"id86","battleId":"battleId86","createdOn":0,"username":"username0","comment":"comment86"},{"id":"id87","battleId":"battleId87","createdOn":0,"username":"username0","comment":"comment87"},{"id":"id88","battleId":"battleId88","createdOn":0,"username":"username0","comment":"comment88"},{"id":"id89","battleId":"battleId89","createdOn":0,"username":"username0","comment":"comment89"},{"id":"id90","battleId":"battleId90","createdOn":0,"username":"username0","comment":"comment90"},{"id":"id91","battleId":"battleId91","createdOn":0,"username":"username0","comment":"comment91"},{"id":"id92","battleId":"battleId92","createdOn":0,"username":"username0","comment":"comment92"},{"id":"id93","battleId":"battleId93","createdOn":0,"username":"username0","comment":"comment93"},{"id":"id94","battleId":"battleId94","createdOn":0,"username":"username0","comment":"comment94"},{"id":"id95","battleId":"battleId95","createdOn":0,"username":"username0","comment":"comment95"},{"id":"id96","battleId":"battleId96","createdOn":0,"username":"username0","comment":"comment96"},{"id":"id97","battleId":"battleId97","createdOn":0,"username":"username0","comment":"comment97"},{"id":"id98","battleId":"battleId98","createdOn":0,"username":"username0","comment":"comment98"},{"id":"id99","battleId":"battleId99","createdOn":0,"username":"username0","comment":"comment99"}]`
	testProcessor(t, indexRepository, commentsRepository, "1", "101", expectedResponse)
}

func TestProcessorBadIndexEntry(t *testing.T) {
	indexRepository := indexMocks.NewRepository()
	key := fmt.Sprintf("commentIds:forAuthor:%s", testUserID)
	indexRepository.SetScore(key, "badId", 0)
	commentsRepository := commentsMocks.NewRepository()
	addComments(indexRepository, commentsRepository, testUserID, false, 1)
	expectedResponse := `[{"id":"id0","battleId":"battleId0","createdOn":0,"username":"username0","comment":"comment0"}]`
	testProcessor(t, indexRepository, commentsRepository, "1", "2", expectedResponse)
}

func TestProcessorDeletedIndexEntry(t *testing.T) {
	indexRepository := indexMocks.NewRepository()
	commentsRepository := commentsMocks.NewRepository()
	addComments(indexRepository, commentsRepository, testUserID, true, 1)
	expectedResponse := `[]`
	testProcessor(t, indexRepository, commentsRepository, "1", "2", expectedResponse)
}

func addComments(indexRepository *indexMocks.Repository, commentsRepository *commentsMocks.Repository, userID string, isDeleted bool, count int) {
	key := fmt.Sprintf("commentIds:forAuthor:%s", userID)
	state := comments.Active
	if isDeleted {
		state = comments.Deleted
	}
	for i := 0; i < count; i++ {
		comment := comments.Comment{
			ID:       fmt.Sprintf("id%d", i),
			BattleID: fmt.Sprintf("battleId%d", i),
			UserID:   userID,
			Comment:  fmt.Sprintf("comment%d", i),
			State:    state,
		}
		commentsRepository.Add(&comment)
		indexRepository.SetScore(key, comment.ID, float64(i))
	}
}

func testProcessor(t *testing.T, indexRepository *indexMocks.Repository, commentsRepository *commentsMocks.Repository, page string, pageSize string, expectedResponseBody string) {
	authContext := &api.AuthContext{
		UserID: testUserID,
	}
	queryParams := make(map[string]string)
	if page != "" {
		queryParams["page"] = page
	}
	if pageSize != "" {
		queryParams["pageSize"] = pageSize
	}
	input := &api.Input{
		AuthContext: authContext,
		QueryParams: queryParams,
	}
	indexer := comments.NewIndexer(indexRepository)
	usersRepository := usersMocks.NewRepository()
	user := &users.User{
		ID:       testUserID,
		Username: "username0",
	}
	usersRepository.Add(user)
	processor := NewProcessor(indexer, commentsRepository, usersRepository)
	output, err := processor.Process(input)
	AssertNil(t, err)
	AssertNotNil(t, output)
	AssertEquals(t, output.StatusCode, http.OK)
	AssertNotNil(t, output.ResponseBody)
	responseBody, _ := json.Marshal(output.ResponseBody)
	AssertEquals(t, string(responseBody), expectedResponseBody)
}
