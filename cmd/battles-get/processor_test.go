package main_test

import (
	"encoding/json"
	"fmt"
	"testing"

	. "github.com/bitterbattles/api/cmd/battles-get"
	"github.com/bitterbattles/api/pkg/battles"
	battlesMocks "github.com/bitterbattles/api/pkg/battles/mocks"
	"github.com/bitterbattles/api/pkg/http"
	indexMocks "github.com/bitterbattles/api/pkg/index/mocks"
	"github.com/bitterbattles/api/pkg/lambda/api"
	. "github.com/bitterbattles/api/pkg/tests"
	"github.com/bitterbattles/api/pkg/users"
	usersMocks "github.com/bitterbattles/api/pkg/users/mocks"
	"github.com/bitterbattles/api/pkg/votes"
	votesMocks "github.com/bitterbattles/api/pkg/votes/mocks"
)

const testSort = "recent"

func TestProcessorFullPage(t *testing.T) {
	indexRepository := indexMocks.NewRepository()
	battlesRepository := battlesMocks.NewRepository()
	addBattles(indexRepository, battlesRepository, testSort, false, 3)
	expectedResponse := `[{"id":"id0","username":"UserID0","title":"title0","description":"description0","canVote":false,"votesFor":0,"votesAgainst":0,"createdOn":0},{"id":"id1","username":"[Deleted]","title":"title1","description":"description1","canVote":false,"votesFor":1,"votesAgainst":2,"createdOn":3}]`
	testProcessor(t, indexRepository, battlesRepository, nil, testSort, "1", "2", expectedResponse)
}

func TestProcessorFullPageLoggedIn(t *testing.T) {
	indexRepository := indexMocks.NewRepository()
	battlesRepository := battlesMocks.NewRepository()
	addBattles(indexRepository, battlesRepository, testSort, false, 3)
	expectedResponse := `[{"id":"id0","username":"UserID0","title":"title0","description":"description0","canVote":false,"votesFor":0,"votesAgainst":0,"createdOn":0},{"id":"id1","username":"[Deleted]","title":"title1","description":"description1","canVote":false,"votesFor":1,"votesAgainst":2,"createdOn":3}]`
	testProcessor(t, indexRepository, battlesRepository, nil, testSort, "1", "2", expectedResponse)
}

func TestProcessorLastPage(t *testing.T) {
	indexRepository := indexMocks.NewRepository()
	battlesRepository := battlesMocks.NewRepository()
	addBattles(indexRepository, battlesRepository, testSort, false, 3)
	expectedResponse := `[{"id":"id2","username":"[Deleted]","title":"title2","description":"description2","canVote":false,"votesFor":2,"votesAgainst":4,"createdOn":6}]`
	testProcessor(t, indexRepository, battlesRepository, nil, testSort, "2", "2", expectedResponse)
}

func TestProcessorBeyondLastPage(t *testing.T) {
	indexRepository := indexMocks.NewRepository()
	battlesRepository := battlesMocks.NewRepository()
	addBattles(indexRepository, battlesRepository, testSort, false, 3)
	expectedResponse := `[]`
	testProcessor(t, indexRepository, battlesRepository, nil, testSort, "3", "2", expectedResponse)
}

func TestProcessorNoPagination(t *testing.T) {
	indexRepository := indexMocks.NewRepository()
	battlesRepository := battlesMocks.NewRepository()
	addBattles(indexRepository, battlesRepository, testSort, false, 50)
	expectedResponse := `[{"id":"id0","username":"UserID0","title":"title0","description":"description0","canVote":false,"votesFor":0,"votesAgainst":0,"createdOn":0},{"id":"id1","username":"[Deleted]","title":"title1","description":"description1","canVote":false,"votesFor":1,"votesAgainst":2,"createdOn":3},{"id":"id2","username":"[Deleted]","title":"title2","description":"description2","canVote":false,"votesFor":2,"votesAgainst":4,"createdOn":6},{"id":"id3","username":"[Deleted]","title":"title3","description":"description3","canVote":false,"votesFor":3,"votesAgainst":6,"createdOn":9},{"id":"id4","username":"[Deleted]","title":"title4","description":"description4","canVote":false,"votesFor":4,"votesAgainst":8,"createdOn":12},{"id":"id5","username":"[Deleted]","title":"title5","description":"description5","canVote":false,"votesFor":5,"votesAgainst":10,"createdOn":15},{"id":"id6","username":"[Deleted]","title":"title6","description":"description6","canVote":false,"votesFor":6,"votesAgainst":12,"createdOn":18},{"id":"id7","username":"[Deleted]","title":"title7","description":"description7","canVote":false,"votesFor":7,"votesAgainst":14,"createdOn":21},{"id":"id8","username":"[Deleted]","title":"title8","description":"description8","canVote":false,"votesFor":8,"votesAgainst":16,"createdOn":24},{"id":"id9","username":"[Deleted]","title":"title9","description":"description9","canVote":false,"votesFor":9,"votesAgainst":18,"createdOn":27},{"id":"id10","username":"[Deleted]","title":"title10","description":"description10","canVote":false,"votesFor":10,"votesAgainst":20,"createdOn":30},{"id":"id11","username":"[Deleted]","title":"title11","description":"description11","canVote":false,"votesFor":11,"votesAgainst":22,"createdOn":33},{"id":"id12","username":"[Deleted]","title":"title12","description":"description12","canVote":false,"votesFor":12,"votesAgainst":24,"createdOn":36},{"id":"id13","username":"[Deleted]","title":"title13","description":"description13","canVote":false,"votesFor":13,"votesAgainst":26,"createdOn":39},{"id":"id14","username":"[Deleted]","title":"title14","description":"description14","canVote":false,"votesFor":14,"votesAgainst":28,"createdOn":42},{"id":"id15","username":"[Deleted]","title":"title15","description":"description15","canVote":false,"votesFor":15,"votesAgainst":30,"createdOn":45},{"id":"id16","username":"[Deleted]","title":"title16","description":"description16","canVote":false,"votesFor":16,"votesAgainst":32,"createdOn":48},{"id":"id17","username":"[Deleted]","title":"title17","description":"description17","canVote":false,"votesFor":17,"votesAgainst":34,"createdOn":51},{"id":"id18","username":"[Deleted]","title":"title18","description":"description18","canVote":false,"votesFor":18,"votesAgainst":36,"createdOn":54},{"id":"id19","username":"[Deleted]","title":"title19","description":"description19","canVote":false,"votesFor":19,"votesAgainst":38,"createdOn":57},{"id":"id20","username":"[Deleted]","title":"title20","description":"description20","canVote":false,"votesFor":20,"votesAgainst":40,"createdOn":60},{"id":"id21","username":"[Deleted]","title":"title21","description":"description21","canVote":false,"votesFor":21,"votesAgainst":42,"createdOn":63},{"id":"id22","username":"[Deleted]","title":"title22","description":"description22","canVote":false,"votesFor":22,"votesAgainst":44,"createdOn":66},{"id":"id23","username":"[Deleted]","title":"title23","description":"description23","canVote":false,"votesFor":23,"votesAgainst":46,"createdOn":69},{"id":"id24","username":"[Deleted]","title":"title24","description":"description24","canVote":false,"votesFor":24,"votesAgainst":48,"createdOn":72},{"id":"id25","username":"[Deleted]","title":"title25","description":"description25","canVote":false,"votesFor":25,"votesAgainst":50,"createdOn":75},{"id":"id26","username":"[Deleted]","title":"title26","description":"description26","canVote":false,"votesFor":26,"votesAgainst":52,"createdOn":78},{"id":"id27","username":"[Deleted]","title":"title27","description":"description27","canVote":false,"votesFor":27,"votesAgainst":54,"createdOn":81},{"id":"id28","username":"[Deleted]","title":"title28","description":"description28","canVote":false,"votesFor":28,"votesAgainst":56,"createdOn":84},{"id":"id29","username":"[Deleted]","title":"title29","description":"description29","canVote":false,"votesFor":29,"votesAgainst":58,"createdOn":87},{"id":"id30","username":"[Deleted]","title":"title30","description":"description30","canVote":false,"votesFor":30,"votesAgainst":60,"createdOn":90},{"id":"id31","username":"[Deleted]","title":"title31","description":"description31","canVote":false,"votesFor":31,"votesAgainst":62,"createdOn":93},{"id":"id32","username":"[Deleted]","title":"title32","description":"description32","canVote":false,"votesFor":32,"votesAgainst":64,"createdOn":96},{"id":"id33","username":"[Deleted]","title":"title33","description":"description33","canVote":false,"votesFor":33,"votesAgainst":66,"createdOn":99},{"id":"id34","username":"[Deleted]","title":"title34","description":"description34","canVote":false,"votesFor":34,"votesAgainst":68,"createdOn":102},{"id":"id35","username":"[Deleted]","title":"title35","description":"description35","canVote":false,"votesFor":35,"votesAgainst":70,"createdOn":105},{"id":"id36","username":"[Deleted]","title":"title36","description":"description36","canVote":false,"votesFor":36,"votesAgainst":72,"createdOn":108},{"id":"id37","username":"[Deleted]","title":"title37","description":"description37","canVote":false,"votesFor":37,"votesAgainst":74,"createdOn":111},{"id":"id38","username":"[Deleted]","title":"title38","description":"description38","canVote":false,"votesFor":38,"votesAgainst":76,"createdOn":114},{"id":"id39","username":"[Deleted]","title":"title39","description":"description39","canVote":false,"votesFor":39,"votesAgainst":78,"createdOn":117},{"id":"id40","username":"[Deleted]","title":"title40","description":"description40","canVote":false,"votesFor":40,"votesAgainst":80,"createdOn":120},{"id":"id41","username":"[Deleted]","title":"title41","description":"description41","canVote":false,"votesFor":41,"votesAgainst":82,"createdOn":123},{"id":"id42","username":"[Deleted]","title":"title42","description":"description42","canVote":false,"votesFor":42,"votesAgainst":84,"createdOn":126},{"id":"id43","username":"[Deleted]","title":"title43","description":"description43","canVote":false,"votesFor":43,"votesAgainst":86,"createdOn":129},{"id":"id44","username":"[Deleted]","title":"title44","description":"description44","canVote":false,"votesFor":44,"votesAgainst":88,"createdOn":132},{"id":"id45","username":"[Deleted]","title":"title45","description":"description45","canVote":false,"votesFor":45,"votesAgainst":90,"createdOn":135},{"id":"id46","username":"[Deleted]","title":"title46","description":"description46","canVote":false,"votesFor":46,"votesAgainst":92,"createdOn":138},{"id":"id47","username":"[Deleted]","title":"title47","description":"description47","canVote":false,"votesFor":47,"votesAgainst":94,"createdOn":141},{"id":"id48","username":"[Deleted]","title":"title48","description":"description48","canVote":false,"votesFor":48,"votesAgainst":96,"createdOn":144},{"id":"id49","username":"[Deleted]","title":"title49","description":"description49","canVote":false,"votesFor":49,"votesAgainst":98,"createdOn":147}]`
	testProcessor(t, indexRepository, battlesRepository, nil, testSort, "", "", expectedResponse)
}

func TestProcessorTooLargePage(t *testing.T) {
	indexRepository := indexMocks.NewRepository()
	battlesRepository := battlesMocks.NewRepository()
	addBattles(indexRepository, battlesRepository, testSort, false, 101)
	expectedResponse := `[{"id":"id0","username":"UserID0","title":"title0","description":"description0","canVote":false,"votesFor":0,"votesAgainst":0,"createdOn":0},{"id":"id1","username":"[Deleted]","title":"title1","description":"description1","canVote":false,"votesFor":1,"votesAgainst":2,"createdOn":3},{"id":"id2","username":"[Deleted]","title":"title2","description":"description2","canVote":false,"votesFor":2,"votesAgainst":4,"createdOn":6},{"id":"id3","username":"[Deleted]","title":"title3","description":"description3","canVote":false,"votesFor":3,"votesAgainst":6,"createdOn":9},{"id":"id4","username":"[Deleted]","title":"title4","description":"description4","canVote":false,"votesFor":4,"votesAgainst":8,"createdOn":12},{"id":"id5","username":"[Deleted]","title":"title5","description":"description5","canVote":false,"votesFor":5,"votesAgainst":10,"createdOn":15},{"id":"id6","username":"[Deleted]","title":"title6","description":"description6","canVote":false,"votesFor":6,"votesAgainst":12,"createdOn":18},{"id":"id7","username":"[Deleted]","title":"title7","description":"description7","canVote":false,"votesFor":7,"votesAgainst":14,"createdOn":21},{"id":"id8","username":"[Deleted]","title":"title8","description":"description8","canVote":false,"votesFor":8,"votesAgainst":16,"createdOn":24},{"id":"id9","username":"[Deleted]","title":"title9","description":"description9","canVote":false,"votesFor":9,"votesAgainst":18,"createdOn":27},{"id":"id10","username":"[Deleted]","title":"title10","description":"description10","canVote":false,"votesFor":10,"votesAgainst":20,"createdOn":30},{"id":"id11","username":"[Deleted]","title":"title11","description":"description11","canVote":false,"votesFor":11,"votesAgainst":22,"createdOn":33},{"id":"id12","username":"[Deleted]","title":"title12","description":"description12","canVote":false,"votesFor":12,"votesAgainst":24,"createdOn":36},{"id":"id13","username":"[Deleted]","title":"title13","description":"description13","canVote":false,"votesFor":13,"votesAgainst":26,"createdOn":39},{"id":"id14","username":"[Deleted]","title":"title14","description":"description14","canVote":false,"votesFor":14,"votesAgainst":28,"createdOn":42},{"id":"id15","username":"[Deleted]","title":"title15","description":"description15","canVote":false,"votesFor":15,"votesAgainst":30,"createdOn":45},{"id":"id16","username":"[Deleted]","title":"title16","description":"description16","canVote":false,"votesFor":16,"votesAgainst":32,"createdOn":48},{"id":"id17","username":"[Deleted]","title":"title17","description":"description17","canVote":false,"votesFor":17,"votesAgainst":34,"createdOn":51},{"id":"id18","username":"[Deleted]","title":"title18","description":"description18","canVote":false,"votesFor":18,"votesAgainst":36,"createdOn":54},{"id":"id19","username":"[Deleted]","title":"title19","description":"description19","canVote":false,"votesFor":19,"votesAgainst":38,"createdOn":57},{"id":"id20","username":"[Deleted]","title":"title20","description":"description20","canVote":false,"votesFor":20,"votesAgainst":40,"createdOn":60},{"id":"id21","username":"[Deleted]","title":"title21","description":"description21","canVote":false,"votesFor":21,"votesAgainst":42,"createdOn":63},{"id":"id22","username":"[Deleted]","title":"title22","description":"description22","canVote":false,"votesFor":22,"votesAgainst":44,"createdOn":66},{"id":"id23","username":"[Deleted]","title":"title23","description":"description23","canVote":false,"votesFor":23,"votesAgainst":46,"createdOn":69},{"id":"id24","username":"[Deleted]","title":"title24","description":"description24","canVote":false,"votesFor":24,"votesAgainst":48,"createdOn":72},{"id":"id25","username":"[Deleted]","title":"title25","description":"description25","canVote":false,"votesFor":25,"votesAgainst":50,"createdOn":75},{"id":"id26","username":"[Deleted]","title":"title26","description":"description26","canVote":false,"votesFor":26,"votesAgainst":52,"createdOn":78},{"id":"id27","username":"[Deleted]","title":"title27","description":"description27","canVote":false,"votesFor":27,"votesAgainst":54,"createdOn":81},{"id":"id28","username":"[Deleted]","title":"title28","description":"description28","canVote":false,"votesFor":28,"votesAgainst":56,"createdOn":84},{"id":"id29","username":"[Deleted]","title":"title29","description":"description29","canVote":false,"votesFor":29,"votesAgainst":58,"createdOn":87},{"id":"id30","username":"[Deleted]","title":"title30","description":"description30","canVote":false,"votesFor":30,"votesAgainst":60,"createdOn":90},{"id":"id31","username":"[Deleted]","title":"title31","description":"description31","canVote":false,"votesFor":31,"votesAgainst":62,"createdOn":93},{"id":"id32","username":"[Deleted]","title":"title32","description":"description32","canVote":false,"votesFor":32,"votesAgainst":64,"createdOn":96},{"id":"id33","username":"[Deleted]","title":"title33","description":"description33","canVote":false,"votesFor":33,"votesAgainst":66,"createdOn":99},{"id":"id34","username":"[Deleted]","title":"title34","description":"description34","canVote":false,"votesFor":34,"votesAgainst":68,"createdOn":102},{"id":"id35","username":"[Deleted]","title":"title35","description":"description35","canVote":false,"votesFor":35,"votesAgainst":70,"createdOn":105},{"id":"id36","username":"[Deleted]","title":"title36","description":"description36","canVote":false,"votesFor":36,"votesAgainst":72,"createdOn":108},{"id":"id37","username":"[Deleted]","title":"title37","description":"description37","canVote":false,"votesFor":37,"votesAgainst":74,"createdOn":111},{"id":"id38","username":"[Deleted]","title":"title38","description":"description38","canVote":false,"votesFor":38,"votesAgainst":76,"createdOn":114},{"id":"id39","username":"[Deleted]","title":"title39","description":"description39","canVote":false,"votesFor":39,"votesAgainst":78,"createdOn":117},{"id":"id40","username":"[Deleted]","title":"title40","description":"description40","canVote":false,"votesFor":40,"votesAgainst":80,"createdOn":120},{"id":"id41","username":"[Deleted]","title":"title41","description":"description41","canVote":false,"votesFor":41,"votesAgainst":82,"createdOn":123},{"id":"id42","username":"[Deleted]","title":"title42","description":"description42","canVote":false,"votesFor":42,"votesAgainst":84,"createdOn":126},{"id":"id43","username":"[Deleted]","title":"title43","description":"description43","canVote":false,"votesFor":43,"votesAgainst":86,"createdOn":129},{"id":"id44","username":"[Deleted]","title":"title44","description":"description44","canVote":false,"votesFor":44,"votesAgainst":88,"createdOn":132},{"id":"id45","username":"[Deleted]","title":"title45","description":"description45","canVote":false,"votesFor":45,"votesAgainst":90,"createdOn":135},{"id":"id46","username":"[Deleted]","title":"title46","description":"description46","canVote":false,"votesFor":46,"votesAgainst":92,"createdOn":138},{"id":"id47","username":"[Deleted]","title":"title47","description":"description47","canVote":false,"votesFor":47,"votesAgainst":94,"createdOn":141},{"id":"id48","username":"[Deleted]","title":"title48","description":"description48","canVote":false,"votesFor":48,"votesAgainst":96,"createdOn":144},{"id":"id49","username":"[Deleted]","title":"title49","description":"description49","canVote":false,"votesFor":49,"votesAgainst":98,"createdOn":147},{"id":"id50","username":"[Deleted]","title":"title50","description":"description50","canVote":false,"votesFor":50,"votesAgainst":100,"createdOn":150},{"id":"id51","username":"[Deleted]","title":"title51","description":"description51","canVote":false,"votesFor":51,"votesAgainst":102,"createdOn":153},{"id":"id52","username":"[Deleted]","title":"title52","description":"description52","canVote":false,"votesFor":52,"votesAgainst":104,"createdOn":156},{"id":"id53","username":"[Deleted]","title":"title53","description":"description53","canVote":false,"votesFor":53,"votesAgainst":106,"createdOn":159},{"id":"id54","username":"[Deleted]","title":"title54","description":"description54","canVote":false,"votesFor":54,"votesAgainst":108,"createdOn":162},{"id":"id55","username":"[Deleted]","title":"title55","description":"description55","canVote":false,"votesFor":55,"votesAgainst":110,"createdOn":165},{"id":"id56","username":"[Deleted]","title":"title56","description":"description56","canVote":false,"votesFor":56,"votesAgainst":112,"createdOn":168},{"id":"id57","username":"[Deleted]","title":"title57","description":"description57","canVote":false,"votesFor":57,"votesAgainst":114,"createdOn":171},{"id":"id58","username":"[Deleted]","title":"title58","description":"description58","canVote":false,"votesFor":58,"votesAgainst":116,"createdOn":174},{"id":"id59","username":"[Deleted]","title":"title59","description":"description59","canVote":false,"votesFor":59,"votesAgainst":118,"createdOn":177},{"id":"id60","username":"[Deleted]","title":"title60","description":"description60","canVote":false,"votesFor":60,"votesAgainst":120,"createdOn":180},{"id":"id61","username":"[Deleted]","title":"title61","description":"description61","canVote":false,"votesFor":61,"votesAgainst":122,"createdOn":183},{"id":"id62","username":"[Deleted]","title":"title62","description":"description62","canVote":false,"votesFor":62,"votesAgainst":124,"createdOn":186},{"id":"id63","username":"[Deleted]","title":"title63","description":"description63","canVote":false,"votesFor":63,"votesAgainst":126,"createdOn":189},{"id":"id64","username":"[Deleted]","title":"title64","description":"description64","canVote":false,"votesFor":64,"votesAgainst":128,"createdOn":192},{"id":"id65","username":"[Deleted]","title":"title65","description":"description65","canVote":false,"votesFor":65,"votesAgainst":130,"createdOn":195},{"id":"id66","username":"[Deleted]","title":"title66","description":"description66","canVote":false,"votesFor":66,"votesAgainst":132,"createdOn":198},{"id":"id67","username":"[Deleted]","title":"title67","description":"description67","canVote":false,"votesFor":67,"votesAgainst":134,"createdOn":201},{"id":"id68","username":"[Deleted]","title":"title68","description":"description68","canVote":false,"votesFor":68,"votesAgainst":136,"createdOn":204},{"id":"id69","username":"[Deleted]","title":"title69","description":"description69","canVote":false,"votesFor":69,"votesAgainst":138,"createdOn":207},{"id":"id70","username":"[Deleted]","title":"title70","description":"description70","canVote":false,"votesFor":70,"votesAgainst":140,"createdOn":210},{"id":"id71","username":"[Deleted]","title":"title71","description":"description71","canVote":false,"votesFor":71,"votesAgainst":142,"createdOn":213},{"id":"id72","username":"[Deleted]","title":"title72","description":"description72","canVote":false,"votesFor":72,"votesAgainst":144,"createdOn":216},{"id":"id73","username":"[Deleted]","title":"title73","description":"description73","canVote":false,"votesFor":73,"votesAgainst":146,"createdOn":219},{"id":"id74","username":"[Deleted]","title":"title74","description":"description74","canVote":false,"votesFor":74,"votesAgainst":148,"createdOn":222},{"id":"id75","username":"[Deleted]","title":"title75","description":"description75","canVote":false,"votesFor":75,"votesAgainst":150,"createdOn":225},{"id":"id76","username":"[Deleted]","title":"title76","description":"description76","canVote":false,"votesFor":76,"votesAgainst":152,"createdOn":228},{"id":"id77","username":"[Deleted]","title":"title77","description":"description77","canVote":false,"votesFor":77,"votesAgainst":154,"createdOn":231},{"id":"id78","username":"[Deleted]","title":"title78","description":"description78","canVote":false,"votesFor":78,"votesAgainst":156,"createdOn":234},{"id":"id79","username":"[Deleted]","title":"title79","description":"description79","canVote":false,"votesFor":79,"votesAgainst":158,"createdOn":237},{"id":"id80","username":"[Deleted]","title":"title80","description":"description80","canVote":false,"votesFor":80,"votesAgainst":160,"createdOn":240},{"id":"id81","username":"[Deleted]","title":"title81","description":"description81","canVote":false,"votesFor":81,"votesAgainst":162,"createdOn":243},{"id":"id82","username":"[Deleted]","title":"title82","description":"description82","canVote":false,"votesFor":82,"votesAgainst":164,"createdOn":246},{"id":"id83","username":"[Deleted]","title":"title83","description":"description83","canVote":false,"votesFor":83,"votesAgainst":166,"createdOn":249},{"id":"id84","username":"[Deleted]","title":"title84","description":"description84","canVote":false,"votesFor":84,"votesAgainst":168,"createdOn":252},{"id":"id85","username":"[Deleted]","title":"title85","description":"description85","canVote":false,"votesFor":85,"votesAgainst":170,"createdOn":255},{"id":"id86","username":"[Deleted]","title":"title86","description":"description86","canVote":false,"votesFor":86,"votesAgainst":172,"createdOn":258},{"id":"id87","username":"[Deleted]","title":"title87","description":"description87","canVote":false,"votesFor":87,"votesAgainst":174,"createdOn":261},{"id":"id88","username":"[Deleted]","title":"title88","description":"description88","canVote":false,"votesFor":88,"votesAgainst":176,"createdOn":264},{"id":"id89","username":"[Deleted]","title":"title89","description":"description89","canVote":false,"votesFor":89,"votesAgainst":178,"createdOn":267},{"id":"id90","username":"[Deleted]","title":"title90","description":"description90","canVote":false,"votesFor":90,"votesAgainst":180,"createdOn":270},{"id":"id91","username":"[Deleted]","title":"title91","description":"description91","canVote":false,"votesFor":91,"votesAgainst":182,"createdOn":273},{"id":"id92","username":"[Deleted]","title":"title92","description":"description92","canVote":false,"votesFor":92,"votesAgainst":184,"createdOn":276},{"id":"id93","username":"[Deleted]","title":"title93","description":"description93","canVote":false,"votesFor":93,"votesAgainst":186,"createdOn":279},{"id":"id94","username":"[Deleted]","title":"title94","description":"description94","canVote":false,"votesFor":94,"votesAgainst":188,"createdOn":282},{"id":"id95","username":"[Deleted]","title":"title95","description":"description95","canVote":false,"votesFor":95,"votesAgainst":190,"createdOn":285},{"id":"id96","username":"[Deleted]","title":"title96","description":"description96","canVote":false,"votesFor":96,"votesAgainst":192,"createdOn":288},{"id":"id97","username":"[Deleted]","title":"title97","description":"description97","canVote":false,"votesFor":97,"votesAgainst":194,"createdOn":291},{"id":"id98","username":"[Deleted]","title":"title98","description":"description98","canVote":false,"votesFor":98,"votesAgainst":196,"createdOn":294},{"id":"id99","username":"[Deleted]","title":"title99","description":"description99","canVote":false,"votesFor":99,"votesAgainst":198,"createdOn":297}]`
	testProcessor(t, indexRepository, battlesRepository, nil, testSort, "1", "101", expectedResponse)
}

func TestProcessorNoSort(t *testing.T) {
	indexRepository := indexMocks.NewRepository()
	battlesRepository := battlesMocks.NewRepository()
	addBattles(indexRepository, battlesRepository, testSort, false, 3)
	expectedResponse := `[{"id":"id0","username":"UserID0","title":"title0","description":"description0","canVote":false,"votesFor":0,"votesAgainst":0,"createdOn":0}]`
	testProcessor(t, indexRepository, battlesRepository, nil, "", "1", "1", expectedResponse)
}

func TestProcessorBadSort(t *testing.T) {
	indexRepository := indexMocks.NewRepository()
	battlesRepository := battlesMocks.NewRepository()
	addBattles(indexRepository, battlesRepository, testSort, false, 3)
	expectedResponse := `[{"id":"id0","username":"UserID0","title":"title0","description":"description0","canVote":false,"votesFor":0,"votesAgainst":0,"createdOn":0}]`
	testProcessor(t, indexRepository, battlesRepository, nil, "bad", "1", "1", expectedResponse)
}

func TestProcessorBadIndexEntry(t *testing.T) {
	indexRepository := indexMocks.NewRepository()
	key := fmt.Sprintf("battleIds:%s", testSort)
	indexRepository.SetScore(key, "badId", 0)
	battlesRepository := battlesMocks.NewRepository()
	addBattles(indexRepository, battlesRepository, testSort, false, 1)
	expectedResponse := `[{"id":"id0","username":"UserID0","title":"title0","description":"description0","canVote":false,"votesFor":0,"votesAgainst":0,"createdOn":0}]`
	testProcessor(t, indexRepository, battlesRepository, nil, testSort, "1", "2", expectedResponse)
}

func TestProcessorDeletedIndexEntry(t *testing.T) {
	indexRepository := indexMocks.NewRepository()
	battlesRepository := battlesMocks.NewRepository()
	addBattles(indexRepository, battlesRepository, testSort, true, 1)
	expectedResponse := `[]`
	testProcessor(t, indexRepository, battlesRepository, nil, testSort, "1", "2", expectedResponse)
}

func TestCanVoteFalse(t *testing.T) {
	authContext := &api.AuthContext{
		UserID: "userId1",
	}
	indexRepository := indexMocks.NewRepository()
	battlesRepository := battlesMocks.NewRepository()
	addBattles(indexRepository, battlesRepository, testSort, false, 1)
	expectedResponse := `[{"id":"id0","username":"UserID0","title":"title0","description":"description0","canVote":false,"votesFor":0,"votesAgainst":0,"createdOn":0}]`
	testProcessor(t, indexRepository, battlesRepository, authContext, testSort, "1", "2", expectedResponse)
}

func TestCanVoteTrue(t *testing.T) {
	authContext := &api.AuthContext{
		UserID: "userId2",
	}
	indexRepository := indexMocks.NewRepository()
	battlesRepository := battlesMocks.NewRepository()
	addBattles(indexRepository, battlesRepository, testSort, false, 1)
	expectedResponse := `[{"id":"id0","username":"UserID0","title":"title0","description":"description0","canVote":true,"votesFor":0,"votesAgainst":0,"createdOn":0}]`
	testProcessor(t, indexRepository, battlesRepository, authContext, testSort, "1", "2", expectedResponse)
}

func addBattles(indexRepository *indexMocks.Repository, battlesRepository *battlesMocks.Repository, sort string, isDeleted bool, count int) {
	key := fmt.Sprintf("battleIds:%s", sort)
	state := battles.Active
	if isDeleted {
		state = battles.Deleted
	}
	for i := 0; i < count; i++ {
		battle := battles.Battle{
			ID:           fmt.Sprintf("id%d", i),
			UserID:       fmt.Sprintf("userId%d", i),
			Title:        fmt.Sprintf("title%d", i),
			Description:  fmt.Sprintf("description%d", i),
			VotesFor:     i,
			VotesAgainst: i * 2,
			CreatedOn:    int64(i * 3),
			State:        state,
		}
		battlesRepository.Add(&battle)
		indexRepository.SetScore(key, battle.ID, float64(i))
	}
}

func testProcessor(t *testing.T, indexRepository *indexMocks.Repository, battlesRepository *battlesMocks.Repository, authContext *api.AuthContext, testSort string, page string, pageSize string, expectedResponseBody string) {
	usersRepository := usersMocks.NewRepository()
	usersRepository.Add(&users.User{
		ID:              "userId0",
		DisplayUsername: "UserID0",
		State:           users.Active,
	})
	usersRepository.Add(&users.User{
		ID:              "userId1",
		DisplayUsername: "UserID1",
		State:           users.Unknown,
	})
	votesRepository := votesMocks.NewRepository()
	votesRepository.Add(&votes.Vote{
		UserID:   "userId1",
		BattleID: "id0",
	})
	queryParams := make(map[string]string)
	if testSort != "" {
		queryParams["sort"] = testSort
	}
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
	indexer := battles.NewIndexer(indexRepository)
	processor := NewProcessor(indexer, battlesRepository, usersRepository, votesRepository)
	output, err := processor.Process(input)
	AssertNil(t, err)
	AssertNotNil(t, output)
	AssertEquals(t, output.StatusCode, http.OK)
	AssertNotNil(t, output.ResponseBody)
	responseBody, _ := json.Marshal(output.ResponseBody)
	AssertEquals(t, string(responseBody), expectedResponseBody)
}
