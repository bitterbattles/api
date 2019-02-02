package main_test

import (
	"encoding/json"
	"fmt"
	"testing"

	. "github.com/bitterbattles/api/cmd/battles-get"
	"github.com/bitterbattles/api/pkg/battles"
	battlesMocks "github.com/bitterbattles/api/pkg/battles/mocks"
	. "github.com/bitterbattles/api/pkg/common/tests"
	ranksMocks "github.com/bitterbattles/api/pkg/ranks/mocks"
)

const sort = battles.RecentSort

func TestHandlerFullPage(t *testing.T) {
	ranksRepository := ranksMocks.NewRepository()
	battlesRepository := battlesMocks.NewRepository()
	addBattles(ranksRepository, battlesRepository, sort, 3)
	expectedResponse := `[{"id":"id0","title":"title0","description":"description0","votesFor":0,"votesAgainst":0,"createdOn":0},{"id":"id1","title":"title1","description":"description1","votesFor":1,"votesAgainst":2,"createdOn":3}]`
	testHandler(t, ranksRepository, battlesRepository, sort, 1, 2, expectedResponse)
}

func TestHandlerLastPage(t *testing.T) {
	ranksRepository := ranksMocks.NewRepository()
	battlesRepository := battlesMocks.NewRepository()
	addBattles(ranksRepository, battlesRepository, sort, 3)
	expectedResponse := `[{"id":"id2","title":"title2","description":"description2","votesFor":2,"votesAgainst":4,"createdOn":6}]`
	testHandler(t, ranksRepository, battlesRepository, sort, 2, 2, expectedResponse)
}

func TestHandlerBeyondLastPage(t *testing.T) {
	ranksRepository := ranksMocks.NewRepository()
	battlesRepository := battlesMocks.NewRepository()
	addBattles(ranksRepository, battlesRepository, sort, 3)
	expectedResponse := `[]`
	testHandler(t, ranksRepository, battlesRepository, sort, 3, 2, expectedResponse)
}

func TestHandlerNoPagination(t *testing.T) {
	ranksRepository := ranksMocks.NewRepository()
	battlesRepository := battlesMocks.NewRepository()
	addBattles(ranksRepository, battlesRepository, sort, 50)
	expectedResponse := `[{"id":"id0","title":"title0","description":"description0","votesFor":0,"votesAgainst":0,"createdOn":0},{"id":"id1","title":"title1","description":"description1","votesFor":1,"votesAgainst":2,"createdOn":3},{"id":"id2","title":"title2","description":"description2","votesFor":2,"votesAgainst":4,"createdOn":6},{"id":"id3","title":"title3","description":"description3","votesFor":3,"votesAgainst":6,"createdOn":9},{"id":"id4","title":"title4","description":"description4","votesFor":4,"votesAgainst":8,"createdOn":12},{"id":"id5","title":"title5","description":"description5","votesFor":5,"votesAgainst":10,"createdOn":15},{"id":"id6","title":"title6","description":"description6","votesFor":6,"votesAgainst":12,"createdOn":18},{"id":"id7","title":"title7","description":"description7","votesFor":7,"votesAgainst":14,"createdOn":21},{"id":"id8","title":"title8","description":"description8","votesFor":8,"votesAgainst":16,"createdOn":24},{"id":"id9","title":"title9","description":"description9","votesFor":9,"votesAgainst":18,"createdOn":27},{"id":"id10","title":"title10","description":"description10","votesFor":10,"votesAgainst":20,"createdOn":30},{"id":"id11","title":"title11","description":"description11","votesFor":11,"votesAgainst":22,"createdOn":33},{"id":"id12","title":"title12","description":"description12","votesFor":12,"votesAgainst":24,"createdOn":36},{"id":"id13","title":"title13","description":"description13","votesFor":13,"votesAgainst":26,"createdOn":39},{"id":"id14","title":"title14","description":"description14","votesFor":14,"votesAgainst":28,"createdOn":42},{"id":"id15","title":"title15","description":"description15","votesFor":15,"votesAgainst":30,"createdOn":45},{"id":"id16","title":"title16","description":"description16","votesFor":16,"votesAgainst":32,"createdOn":48},{"id":"id17","title":"title17","description":"description17","votesFor":17,"votesAgainst":34,"createdOn":51},{"id":"id18","title":"title18","description":"description18","votesFor":18,"votesAgainst":36,"createdOn":54},{"id":"id19","title":"title19","description":"description19","votesFor":19,"votesAgainst":38,"createdOn":57},{"id":"id20","title":"title20","description":"description20","votesFor":20,"votesAgainst":40,"createdOn":60},{"id":"id21","title":"title21","description":"description21","votesFor":21,"votesAgainst":42,"createdOn":63},{"id":"id22","title":"title22","description":"description22","votesFor":22,"votesAgainst":44,"createdOn":66},{"id":"id23","title":"title23","description":"description23","votesFor":23,"votesAgainst":46,"createdOn":69},{"id":"id24","title":"title24","description":"description24","votesFor":24,"votesAgainst":48,"createdOn":72},{"id":"id25","title":"title25","description":"description25","votesFor":25,"votesAgainst":50,"createdOn":75},{"id":"id26","title":"title26","description":"description26","votesFor":26,"votesAgainst":52,"createdOn":78},{"id":"id27","title":"title27","description":"description27","votesFor":27,"votesAgainst":54,"createdOn":81},{"id":"id28","title":"title28","description":"description28","votesFor":28,"votesAgainst":56,"createdOn":84},{"id":"id29","title":"title29","description":"description29","votesFor":29,"votesAgainst":58,"createdOn":87},{"id":"id30","title":"title30","description":"description30","votesFor":30,"votesAgainst":60,"createdOn":90},{"id":"id31","title":"title31","description":"description31","votesFor":31,"votesAgainst":62,"createdOn":93},{"id":"id32","title":"title32","description":"description32","votesFor":32,"votesAgainst":64,"createdOn":96},{"id":"id33","title":"title33","description":"description33","votesFor":33,"votesAgainst":66,"createdOn":99},{"id":"id34","title":"title34","description":"description34","votesFor":34,"votesAgainst":68,"createdOn":102},{"id":"id35","title":"title35","description":"description35","votesFor":35,"votesAgainst":70,"createdOn":105},{"id":"id36","title":"title36","description":"description36","votesFor":36,"votesAgainst":72,"createdOn":108},{"id":"id37","title":"title37","description":"description37","votesFor":37,"votesAgainst":74,"createdOn":111},{"id":"id38","title":"title38","description":"description38","votesFor":38,"votesAgainst":76,"createdOn":114},{"id":"id39","title":"title39","description":"description39","votesFor":39,"votesAgainst":78,"createdOn":117},{"id":"id40","title":"title40","description":"description40","votesFor":40,"votesAgainst":80,"createdOn":120},{"id":"id41","title":"title41","description":"description41","votesFor":41,"votesAgainst":82,"createdOn":123},{"id":"id42","title":"title42","description":"description42","votesFor":42,"votesAgainst":84,"createdOn":126},{"id":"id43","title":"title43","description":"description43","votesFor":43,"votesAgainst":86,"createdOn":129},{"id":"id44","title":"title44","description":"description44","votesFor":44,"votesAgainst":88,"createdOn":132},{"id":"id45","title":"title45","description":"description45","votesFor":45,"votesAgainst":90,"createdOn":135},{"id":"id46","title":"title46","description":"description46","votesFor":46,"votesAgainst":92,"createdOn":138},{"id":"id47","title":"title47","description":"description47","votesFor":47,"votesAgainst":94,"createdOn":141},{"id":"id48","title":"title48","description":"description48","votesFor":48,"votesAgainst":96,"createdOn":144},{"id":"id49","title":"title49","description":"description49","votesFor":49,"votesAgainst":98,"createdOn":147}]`
	testHandler(t, ranksRepository, battlesRepository, sort, 0, 0, expectedResponse)
}

func TestHandlerTooLargePage(t *testing.T) {
	ranksRepository := ranksMocks.NewRepository()
	battlesRepository := battlesMocks.NewRepository()
	addBattles(ranksRepository, battlesRepository, sort, 101)
	expectedResponse := `[{"id":"id0","title":"title0","description":"description0","votesFor":0,"votesAgainst":0,"createdOn":0},{"id":"id1","title":"title1","description":"description1","votesFor":1,"votesAgainst":2,"createdOn":3},{"id":"id2","title":"title2","description":"description2","votesFor":2,"votesAgainst":4,"createdOn":6},{"id":"id3","title":"title3","description":"description3","votesFor":3,"votesAgainst":6,"createdOn":9},{"id":"id4","title":"title4","description":"description4","votesFor":4,"votesAgainst":8,"createdOn":12},{"id":"id5","title":"title5","description":"description5","votesFor":5,"votesAgainst":10,"createdOn":15},{"id":"id6","title":"title6","description":"description6","votesFor":6,"votesAgainst":12,"createdOn":18},{"id":"id7","title":"title7","description":"description7","votesFor":7,"votesAgainst":14,"createdOn":21},{"id":"id8","title":"title8","description":"description8","votesFor":8,"votesAgainst":16,"createdOn":24},{"id":"id9","title":"title9","description":"description9","votesFor":9,"votesAgainst":18,"createdOn":27},{"id":"id10","title":"title10","description":"description10","votesFor":10,"votesAgainst":20,"createdOn":30},{"id":"id11","title":"title11","description":"description11","votesFor":11,"votesAgainst":22,"createdOn":33},{"id":"id12","title":"title12","description":"description12","votesFor":12,"votesAgainst":24,"createdOn":36},{"id":"id13","title":"title13","description":"description13","votesFor":13,"votesAgainst":26,"createdOn":39},{"id":"id14","title":"title14","description":"description14","votesFor":14,"votesAgainst":28,"createdOn":42},{"id":"id15","title":"title15","description":"description15","votesFor":15,"votesAgainst":30,"createdOn":45},{"id":"id16","title":"title16","description":"description16","votesFor":16,"votesAgainst":32,"createdOn":48},{"id":"id17","title":"title17","description":"description17","votesFor":17,"votesAgainst":34,"createdOn":51},{"id":"id18","title":"title18","description":"description18","votesFor":18,"votesAgainst":36,"createdOn":54},{"id":"id19","title":"title19","description":"description19","votesFor":19,"votesAgainst":38,"createdOn":57},{"id":"id20","title":"title20","description":"description20","votesFor":20,"votesAgainst":40,"createdOn":60},{"id":"id21","title":"title21","description":"description21","votesFor":21,"votesAgainst":42,"createdOn":63},{"id":"id22","title":"title22","description":"description22","votesFor":22,"votesAgainst":44,"createdOn":66},{"id":"id23","title":"title23","description":"description23","votesFor":23,"votesAgainst":46,"createdOn":69},{"id":"id24","title":"title24","description":"description24","votesFor":24,"votesAgainst":48,"createdOn":72},{"id":"id25","title":"title25","description":"description25","votesFor":25,"votesAgainst":50,"createdOn":75},{"id":"id26","title":"title26","description":"description26","votesFor":26,"votesAgainst":52,"createdOn":78},{"id":"id27","title":"title27","description":"description27","votesFor":27,"votesAgainst":54,"createdOn":81},{"id":"id28","title":"title28","description":"description28","votesFor":28,"votesAgainst":56,"createdOn":84},{"id":"id29","title":"title29","description":"description29","votesFor":29,"votesAgainst":58,"createdOn":87},{"id":"id30","title":"title30","description":"description30","votesFor":30,"votesAgainst":60,"createdOn":90},{"id":"id31","title":"title31","description":"description31","votesFor":31,"votesAgainst":62,"createdOn":93},{"id":"id32","title":"title32","description":"description32","votesFor":32,"votesAgainst":64,"createdOn":96},{"id":"id33","title":"title33","description":"description33","votesFor":33,"votesAgainst":66,"createdOn":99},{"id":"id34","title":"title34","description":"description34","votesFor":34,"votesAgainst":68,"createdOn":102},{"id":"id35","title":"title35","description":"description35","votesFor":35,"votesAgainst":70,"createdOn":105},{"id":"id36","title":"title36","description":"description36","votesFor":36,"votesAgainst":72,"createdOn":108},{"id":"id37","title":"title37","description":"description37","votesFor":37,"votesAgainst":74,"createdOn":111},{"id":"id38","title":"title38","description":"description38","votesFor":38,"votesAgainst":76,"createdOn":114},{"id":"id39","title":"title39","description":"description39","votesFor":39,"votesAgainst":78,"createdOn":117},{"id":"id40","title":"title40","description":"description40","votesFor":40,"votesAgainst":80,"createdOn":120},{"id":"id41","title":"title41","description":"description41","votesFor":41,"votesAgainst":82,"createdOn":123},{"id":"id42","title":"title42","description":"description42","votesFor":42,"votesAgainst":84,"createdOn":126},{"id":"id43","title":"title43","description":"description43","votesFor":43,"votesAgainst":86,"createdOn":129},{"id":"id44","title":"title44","description":"description44","votesFor":44,"votesAgainst":88,"createdOn":132},{"id":"id45","title":"title45","description":"description45","votesFor":45,"votesAgainst":90,"createdOn":135},{"id":"id46","title":"title46","description":"description46","votesFor":46,"votesAgainst":92,"createdOn":138},{"id":"id47","title":"title47","description":"description47","votesFor":47,"votesAgainst":94,"createdOn":141},{"id":"id48","title":"title48","description":"description48","votesFor":48,"votesAgainst":96,"createdOn":144},{"id":"id49","title":"title49","description":"description49","votesFor":49,"votesAgainst":98,"createdOn":147},{"id":"id50","title":"title50","description":"description50","votesFor":50,"votesAgainst":100,"createdOn":150},{"id":"id51","title":"title51","description":"description51","votesFor":51,"votesAgainst":102,"createdOn":153},{"id":"id52","title":"title52","description":"description52","votesFor":52,"votesAgainst":104,"createdOn":156},{"id":"id53","title":"title53","description":"description53","votesFor":53,"votesAgainst":106,"createdOn":159},{"id":"id54","title":"title54","description":"description54","votesFor":54,"votesAgainst":108,"createdOn":162},{"id":"id55","title":"title55","description":"description55","votesFor":55,"votesAgainst":110,"createdOn":165},{"id":"id56","title":"title56","description":"description56","votesFor":56,"votesAgainst":112,"createdOn":168},{"id":"id57","title":"title57","description":"description57","votesFor":57,"votesAgainst":114,"createdOn":171},{"id":"id58","title":"title58","description":"description58","votesFor":58,"votesAgainst":116,"createdOn":174},{"id":"id59","title":"title59","description":"description59","votesFor":59,"votesAgainst":118,"createdOn":177},{"id":"id60","title":"title60","description":"description60","votesFor":60,"votesAgainst":120,"createdOn":180},{"id":"id61","title":"title61","description":"description61","votesFor":61,"votesAgainst":122,"createdOn":183},{"id":"id62","title":"title62","description":"description62","votesFor":62,"votesAgainst":124,"createdOn":186},{"id":"id63","title":"title63","description":"description63","votesFor":63,"votesAgainst":126,"createdOn":189},{"id":"id64","title":"title64","description":"description64","votesFor":64,"votesAgainst":128,"createdOn":192},{"id":"id65","title":"title65","description":"description65","votesFor":65,"votesAgainst":130,"createdOn":195},{"id":"id66","title":"title66","description":"description66","votesFor":66,"votesAgainst":132,"createdOn":198},{"id":"id67","title":"title67","description":"description67","votesFor":67,"votesAgainst":134,"createdOn":201},{"id":"id68","title":"title68","description":"description68","votesFor":68,"votesAgainst":136,"createdOn":204},{"id":"id69","title":"title69","description":"description69","votesFor":69,"votesAgainst":138,"createdOn":207},{"id":"id70","title":"title70","description":"description70","votesFor":70,"votesAgainst":140,"createdOn":210},{"id":"id71","title":"title71","description":"description71","votesFor":71,"votesAgainst":142,"createdOn":213},{"id":"id72","title":"title72","description":"description72","votesFor":72,"votesAgainst":144,"createdOn":216},{"id":"id73","title":"title73","description":"description73","votesFor":73,"votesAgainst":146,"createdOn":219},{"id":"id74","title":"title74","description":"description74","votesFor":74,"votesAgainst":148,"createdOn":222},{"id":"id75","title":"title75","description":"description75","votesFor":75,"votesAgainst":150,"createdOn":225},{"id":"id76","title":"title76","description":"description76","votesFor":76,"votesAgainst":152,"createdOn":228},{"id":"id77","title":"title77","description":"description77","votesFor":77,"votesAgainst":154,"createdOn":231},{"id":"id78","title":"title78","description":"description78","votesFor":78,"votesAgainst":156,"createdOn":234},{"id":"id79","title":"title79","description":"description79","votesFor":79,"votesAgainst":158,"createdOn":237},{"id":"id80","title":"title80","description":"description80","votesFor":80,"votesAgainst":160,"createdOn":240},{"id":"id81","title":"title81","description":"description81","votesFor":81,"votesAgainst":162,"createdOn":243},{"id":"id82","title":"title82","description":"description82","votesFor":82,"votesAgainst":164,"createdOn":246},{"id":"id83","title":"title83","description":"description83","votesFor":83,"votesAgainst":166,"createdOn":249},{"id":"id84","title":"title84","description":"description84","votesFor":84,"votesAgainst":168,"createdOn":252},{"id":"id85","title":"title85","description":"description85","votesFor":85,"votesAgainst":170,"createdOn":255},{"id":"id86","title":"title86","description":"description86","votesFor":86,"votesAgainst":172,"createdOn":258},{"id":"id87","title":"title87","description":"description87","votesFor":87,"votesAgainst":174,"createdOn":261},{"id":"id88","title":"title88","description":"description88","votesFor":88,"votesAgainst":176,"createdOn":264},{"id":"id89","title":"title89","description":"description89","votesFor":89,"votesAgainst":178,"createdOn":267},{"id":"id90","title":"title90","description":"description90","votesFor":90,"votesAgainst":180,"createdOn":270},{"id":"id91","title":"title91","description":"description91","votesFor":91,"votesAgainst":182,"createdOn":273},{"id":"id92","title":"title92","description":"description92","votesFor":92,"votesAgainst":184,"createdOn":276},{"id":"id93","title":"title93","description":"description93","votesFor":93,"votesAgainst":186,"createdOn":279},{"id":"id94","title":"title94","description":"description94","votesFor":94,"votesAgainst":188,"createdOn":282},{"id":"id95","title":"title95","description":"description95","votesFor":95,"votesAgainst":190,"createdOn":285},{"id":"id96","title":"title96","description":"description96","votesFor":96,"votesAgainst":192,"createdOn":288},{"id":"id97","title":"title97","description":"description97","votesFor":97,"votesAgainst":194,"createdOn":291},{"id":"id98","title":"title98","description":"description98","votesFor":98,"votesAgainst":196,"createdOn":294},{"id":"id99","title":"title99","description":"description99","votesFor":99,"votesAgainst":198,"createdOn":297}]`
	testHandler(t, ranksRepository, battlesRepository, sort, 1, 101, expectedResponse)
}

func TestHandlerNoSort(t *testing.T) {
	ranksRepository := ranksMocks.NewRepository()
	battlesRepository := battlesMocks.NewRepository()
	addBattles(ranksRepository, battlesRepository, sort, 3)
	expectedResponse := `[{"id":"id0","title":"title0","description":"description0","votesFor":0,"votesAgainst":0,"createdOn":0}]`
	testHandler(t, ranksRepository, battlesRepository, "", 1, 1, expectedResponse)
}

func TestHandlerBadSort(t *testing.T) {
	ranksRepository := ranksMocks.NewRepository()
	battlesRepository := battlesMocks.NewRepository()
	addBattles(ranksRepository, battlesRepository, sort, 3)
	expectedResponse := `[{"id":"id0","title":"title0","description":"description0","votesFor":0,"votesAgainst":0,"createdOn":0}]`
	testHandler(t, ranksRepository, battlesRepository, "bad", 1, 1, expectedResponse)
}

func TestHandlerBadRankEntry(t *testing.T) {
	ranksRepository := ranksMocks.NewRepository()
	ranksRepository.Upsert(sort, "badId", 0)
	battlesRepository := battlesMocks.NewRepository()
	addBattles(ranksRepository, battlesRepository, sort, 1)
	expectedResponse := `[{"id":"id0","title":"title0","description":"description0","votesFor":0,"votesAgainst":0,"createdOn":0}]`
	testHandler(t, ranksRepository, battlesRepository, sort, 1, 2, expectedResponse)
}

func addBattles(ranksRepository *ranksMocks.Repository, battlesRepository *battlesMocks.Repository, category string, count int) {
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
		battlesRepository.Add(battle)
		ranksRepository.Upsert(category, battle.ID, uint64(i))
	}
}

func testHandler(t *testing.T, ranksRepository *ranksMocks.Repository, battlesRepository *battlesMocks.Repository, sort string, page int, pageSize int, expectedResponse string) {
	handler := NewHandler(ranksRepository, battlesRepository)
	request := Request{
		Sort:     sort,
		Page:     page,
		PageSize: pageSize,
	}
	requestBytes, _ := json.Marshal(request)
	responseBytes, err := handler.Invoke(nil, requestBytes)
	AssertNil(t, err)
	AssertEquals(t, string(responseBytes), expectedResponse)
}
