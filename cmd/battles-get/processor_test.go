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
	"github.com/bitterbattles/api/pkg/time"
	"github.com/bitterbattles/api/pkg/votes"
	votesMocks "github.com/bitterbattles/api/pkg/votes/mocks"
)

const testSort = "recent"

func TestProcessorFullPage(t *testing.T) {
	indexRepository := indexMocks.NewRepository()
	battlesRepository := battlesMocks.NewRepository()
	addBattles(indexRepository, battlesRepository, testSort, false, 3)
	expectedResponse := `[{"id":"id0","createdOn":0,"username":"username0","title":"title0","description":"description0","canVote":false,"votesFor":0,"votesAgainst":0,"comments":0,"verdict":4},{"id":"id1","createdOn":3,"username":"username1","title":"title1","description":"description1","canVote":false,"votesFor":1,"votesAgainst":2,"comments":3,"verdict":3}]`
	testProcessor(t, indexRepository, battlesRepository, nil, testSort, "1", "2", expectedResponse)
}

func TestProcessorFullPageLoggedIn(t *testing.T) {
	indexRepository := indexMocks.NewRepository()
	battlesRepository := battlesMocks.NewRepository()
	addBattles(indexRepository, battlesRepository, testSort, false, 3)
	expectedResponse := `[{"id":"id0","createdOn":0,"username":"username0","title":"title0","description":"description0","canVote":false,"votesFor":0,"votesAgainst":0,"comments":0,"verdict":4},{"id":"id1","createdOn":3,"username":"username1","title":"title1","description":"description1","canVote":false,"votesFor":1,"votesAgainst":2,"comments":3,"verdict":3}]`
	testProcessor(t, indexRepository, battlesRepository, nil, testSort, "1", "2", expectedResponse)
}

func TestProcessorLastPage(t *testing.T) {
	indexRepository := indexMocks.NewRepository()
	battlesRepository := battlesMocks.NewRepository()
	addBattles(indexRepository, battlesRepository, testSort, false, 3)
	expectedResponse := `[{"id":"id2","createdOn":6,"username":"username2","title":"title2","description":"description2","canVote":false,"votesFor":2,"votesAgainst":4,"comments":6,"verdict":3}]`
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
	expectedResponse := `[{"id":"id0","createdOn":0,"username":"username0","title":"title0","description":"description0","canVote":false,"votesFor":0,"votesAgainst":0,"comments":0,"verdict":4},{"id":"id1","createdOn":3,"username":"username1","title":"title1","description":"description1","canVote":false,"votesFor":1,"votesAgainst":2,"comments":3,"verdict":3},{"id":"id2","createdOn":6,"username":"username2","title":"title2","description":"description2","canVote":false,"votesFor":2,"votesAgainst":4,"comments":6,"verdict":3},{"id":"id3","createdOn":9,"username":"username3","title":"title3","description":"description3","canVote":false,"votesFor":3,"votesAgainst":6,"comments":9,"verdict":3},{"id":"id4","createdOn":12,"username":"username4","title":"title4","description":"description4","canVote":false,"votesFor":4,"votesAgainst":8,"comments":12,"verdict":3},{"id":"id5","createdOn":15,"username":"username5","title":"title5","description":"description5","canVote":false,"votesFor":5,"votesAgainst":10,"comments":15,"verdict":3},{"id":"id6","createdOn":18,"username":"username6","title":"title6","description":"description6","canVote":false,"votesFor":6,"votesAgainst":12,"comments":18,"verdict":3},{"id":"id7","createdOn":21,"username":"username7","title":"title7","description":"description7","canVote":false,"votesFor":7,"votesAgainst":14,"comments":21,"verdict":3},{"id":"id8","createdOn":24,"username":"username8","title":"title8","description":"description8","canVote":false,"votesFor":8,"votesAgainst":16,"comments":24,"verdict":3},{"id":"id9","createdOn":27,"username":"username9","title":"title9","description":"description9","canVote":false,"votesFor":9,"votesAgainst":18,"comments":27,"verdict":3},{"id":"id10","createdOn":30,"username":"username10","title":"title10","description":"description10","canVote":false,"votesFor":10,"votesAgainst":20,"comments":30,"verdict":3},{"id":"id11","createdOn":33,"username":"username11","title":"title11","description":"description11","canVote":false,"votesFor":11,"votesAgainst":22,"comments":33,"verdict":3},{"id":"id12","createdOn":36,"username":"username12","title":"title12","description":"description12","canVote":false,"votesFor":12,"votesAgainst":24,"comments":36,"verdict":3},{"id":"id13","createdOn":39,"username":"username13","title":"title13","description":"description13","canVote":false,"votesFor":13,"votesAgainst":26,"comments":39,"verdict":3},{"id":"id14","createdOn":42,"username":"username14","title":"title14","description":"description14","canVote":false,"votesFor":14,"votesAgainst":28,"comments":42,"verdict":3},{"id":"id15","createdOn":45,"username":"username15","title":"title15","description":"description15","canVote":false,"votesFor":15,"votesAgainst":30,"comments":45,"verdict":3},{"id":"id16","createdOn":48,"username":"username16","title":"title16","description":"description16","canVote":false,"votesFor":16,"votesAgainst":32,"comments":48,"verdict":3},{"id":"id17","createdOn":51,"username":"username17","title":"title17","description":"description17","canVote":false,"votesFor":17,"votesAgainst":34,"comments":51,"verdict":3},{"id":"id18","createdOn":54,"username":"username18","title":"title18","description":"description18","canVote":false,"votesFor":18,"votesAgainst":36,"comments":54,"verdict":3},{"id":"id19","createdOn":57,"username":"username19","title":"title19","description":"description19","canVote":false,"votesFor":19,"votesAgainst":38,"comments":57,"verdict":3},{"id":"id20","createdOn":60,"username":"username20","title":"title20","description":"description20","canVote":false,"votesFor":20,"votesAgainst":40,"comments":60,"verdict":3},{"id":"id21","createdOn":63,"username":"username21","title":"title21","description":"description21","canVote":false,"votesFor":21,"votesAgainst":42,"comments":63,"verdict":3},{"id":"id22","createdOn":66,"username":"username22","title":"title22","description":"description22","canVote":false,"votesFor":22,"votesAgainst":44,"comments":66,"verdict":3},{"id":"id23","createdOn":69,"username":"username23","title":"title23","description":"description23","canVote":false,"votesFor":23,"votesAgainst":46,"comments":69,"verdict":3},{"id":"id24","createdOn":72,"username":"username24","title":"title24","description":"description24","canVote":false,"votesFor":24,"votesAgainst":48,"comments":72,"verdict":3},{"id":"id25","createdOn":75,"username":"username25","title":"title25","description":"description25","canVote":false,"votesFor":25,"votesAgainst":50,"comments":75,"verdict":3},{"id":"id26","createdOn":78,"username":"username26","title":"title26","description":"description26","canVote":false,"votesFor":26,"votesAgainst":52,"comments":78,"verdict":3},{"id":"id27","createdOn":81,"username":"username27","title":"title27","description":"description27","canVote":false,"votesFor":27,"votesAgainst":54,"comments":81,"verdict":3},{"id":"id28","createdOn":84,"username":"username28","title":"title28","description":"description28","canVote":false,"votesFor":28,"votesAgainst":56,"comments":84,"verdict":3},{"id":"id29","createdOn":87,"username":"username29","title":"title29","description":"description29","canVote":false,"votesFor":29,"votesAgainst":58,"comments":87,"verdict":3},{"id":"id30","createdOn":90,"username":"username30","title":"title30","description":"description30","canVote":false,"votesFor":30,"votesAgainst":60,"comments":90,"verdict":3},{"id":"id31","createdOn":93,"username":"username31","title":"title31","description":"description31","canVote":false,"votesFor":31,"votesAgainst":62,"comments":93,"verdict":3},{"id":"id32","createdOn":96,"username":"username32","title":"title32","description":"description32","canVote":false,"votesFor":32,"votesAgainst":64,"comments":96,"verdict":3},{"id":"id33","createdOn":99,"username":"username33","title":"title33","description":"description33","canVote":false,"votesFor":33,"votesAgainst":66,"comments":99,"verdict":3},{"id":"id34","createdOn":102,"username":"username34","title":"title34","description":"description34","canVote":false,"votesFor":34,"votesAgainst":68,"comments":102,"verdict":3},{"id":"id35","createdOn":105,"username":"username35","title":"title35","description":"description35","canVote":false,"votesFor":35,"votesAgainst":70,"comments":105,"verdict":3},{"id":"id36","createdOn":108,"username":"username36","title":"title36","description":"description36","canVote":false,"votesFor":36,"votesAgainst":72,"comments":108,"verdict":3},{"id":"id37","createdOn":111,"username":"username37","title":"title37","description":"description37","canVote":false,"votesFor":37,"votesAgainst":74,"comments":111,"verdict":3},{"id":"id38","createdOn":114,"username":"username38","title":"title38","description":"description38","canVote":false,"votesFor":38,"votesAgainst":76,"comments":114,"verdict":3},{"id":"id39","createdOn":117,"username":"username39","title":"title39","description":"description39","canVote":false,"votesFor":39,"votesAgainst":78,"comments":117,"verdict":3},{"id":"id40","createdOn":120,"username":"username40","title":"title40","description":"description40","canVote":false,"votesFor":40,"votesAgainst":80,"comments":120,"verdict":3},{"id":"id41","createdOn":123,"username":"username41","title":"title41","description":"description41","canVote":false,"votesFor":41,"votesAgainst":82,"comments":123,"verdict":3},{"id":"id42","createdOn":126,"username":"username42","title":"title42","description":"description42","canVote":false,"votesFor":42,"votesAgainst":84,"comments":126,"verdict":3},{"id":"id43","createdOn":129,"username":"username43","title":"title43","description":"description43","canVote":false,"votesFor":43,"votesAgainst":86,"comments":129,"verdict":3},{"id":"id44","createdOn":132,"username":"username44","title":"title44","description":"description44","canVote":false,"votesFor":44,"votesAgainst":88,"comments":132,"verdict":3},{"id":"id45","createdOn":135,"username":"username45","title":"title45","description":"description45","canVote":false,"votesFor":45,"votesAgainst":90,"comments":135,"verdict":3},{"id":"id46","createdOn":138,"username":"username46","title":"title46","description":"description46","canVote":false,"votesFor":46,"votesAgainst":92,"comments":138,"verdict":3},{"id":"id47","createdOn":141,"username":"username47","title":"title47","description":"description47","canVote":false,"votesFor":47,"votesAgainst":94,"comments":141,"verdict":3},{"id":"id48","createdOn":144,"username":"username48","title":"title48","description":"description48","canVote":false,"votesFor":48,"votesAgainst":96,"comments":144,"verdict":3},{"id":"id49","createdOn":147,"username":"username49","title":"title49","description":"description49","canVote":false,"votesFor":49,"votesAgainst":98,"comments":147,"verdict":3}]`
	testProcessor(t, indexRepository, battlesRepository, nil, testSort, "", "", expectedResponse)
}

func TestProcessorTooLargePage(t *testing.T) {
	indexRepository := indexMocks.NewRepository()
	battlesRepository := battlesMocks.NewRepository()
	addBattles(indexRepository, battlesRepository, testSort, false, 101)
	expectedResponse := `[{"id":"id0","createdOn":0,"username":"username0","title":"title0","description":"description0","canVote":false,"votesFor":0,"votesAgainst":0,"comments":0,"verdict":4},{"id":"id1","createdOn":3,"username":"username1","title":"title1","description":"description1","canVote":false,"votesFor":1,"votesAgainst":2,"comments":3,"verdict":3},{"id":"id2","createdOn":6,"username":"username2","title":"title2","description":"description2","canVote":false,"votesFor":2,"votesAgainst":4,"comments":6,"verdict":3},{"id":"id3","createdOn":9,"username":"username3","title":"title3","description":"description3","canVote":false,"votesFor":3,"votesAgainst":6,"comments":9,"verdict":3},{"id":"id4","createdOn":12,"username":"username4","title":"title4","description":"description4","canVote":false,"votesFor":4,"votesAgainst":8,"comments":12,"verdict":3},{"id":"id5","createdOn":15,"username":"username5","title":"title5","description":"description5","canVote":false,"votesFor":5,"votesAgainst":10,"comments":15,"verdict":3},{"id":"id6","createdOn":18,"username":"username6","title":"title6","description":"description6","canVote":false,"votesFor":6,"votesAgainst":12,"comments":18,"verdict":3},{"id":"id7","createdOn":21,"username":"username7","title":"title7","description":"description7","canVote":false,"votesFor":7,"votesAgainst":14,"comments":21,"verdict":3},{"id":"id8","createdOn":24,"username":"username8","title":"title8","description":"description8","canVote":false,"votesFor":8,"votesAgainst":16,"comments":24,"verdict":3},{"id":"id9","createdOn":27,"username":"username9","title":"title9","description":"description9","canVote":false,"votesFor":9,"votesAgainst":18,"comments":27,"verdict":3},{"id":"id10","createdOn":30,"username":"username10","title":"title10","description":"description10","canVote":false,"votesFor":10,"votesAgainst":20,"comments":30,"verdict":3},{"id":"id11","createdOn":33,"username":"username11","title":"title11","description":"description11","canVote":false,"votesFor":11,"votesAgainst":22,"comments":33,"verdict":3},{"id":"id12","createdOn":36,"username":"username12","title":"title12","description":"description12","canVote":false,"votesFor":12,"votesAgainst":24,"comments":36,"verdict":3},{"id":"id13","createdOn":39,"username":"username13","title":"title13","description":"description13","canVote":false,"votesFor":13,"votesAgainst":26,"comments":39,"verdict":3},{"id":"id14","createdOn":42,"username":"username14","title":"title14","description":"description14","canVote":false,"votesFor":14,"votesAgainst":28,"comments":42,"verdict":3},{"id":"id15","createdOn":45,"username":"username15","title":"title15","description":"description15","canVote":false,"votesFor":15,"votesAgainst":30,"comments":45,"verdict":3},{"id":"id16","createdOn":48,"username":"username16","title":"title16","description":"description16","canVote":false,"votesFor":16,"votesAgainst":32,"comments":48,"verdict":3},{"id":"id17","createdOn":51,"username":"username17","title":"title17","description":"description17","canVote":false,"votesFor":17,"votesAgainst":34,"comments":51,"verdict":3},{"id":"id18","createdOn":54,"username":"username18","title":"title18","description":"description18","canVote":false,"votesFor":18,"votesAgainst":36,"comments":54,"verdict":3},{"id":"id19","createdOn":57,"username":"username19","title":"title19","description":"description19","canVote":false,"votesFor":19,"votesAgainst":38,"comments":57,"verdict":3},{"id":"id20","createdOn":60,"username":"username20","title":"title20","description":"description20","canVote":false,"votesFor":20,"votesAgainst":40,"comments":60,"verdict":3},{"id":"id21","createdOn":63,"username":"username21","title":"title21","description":"description21","canVote":false,"votesFor":21,"votesAgainst":42,"comments":63,"verdict":3},{"id":"id22","createdOn":66,"username":"username22","title":"title22","description":"description22","canVote":false,"votesFor":22,"votesAgainst":44,"comments":66,"verdict":3},{"id":"id23","createdOn":69,"username":"username23","title":"title23","description":"description23","canVote":false,"votesFor":23,"votesAgainst":46,"comments":69,"verdict":3},{"id":"id24","createdOn":72,"username":"username24","title":"title24","description":"description24","canVote":false,"votesFor":24,"votesAgainst":48,"comments":72,"verdict":3},{"id":"id25","createdOn":75,"username":"username25","title":"title25","description":"description25","canVote":false,"votesFor":25,"votesAgainst":50,"comments":75,"verdict":3},{"id":"id26","createdOn":78,"username":"username26","title":"title26","description":"description26","canVote":false,"votesFor":26,"votesAgainst":52,"comments":78,"verdict":3},{"id":"id27","createdOn":81,"username":"username27","title":"title27","description":"description27","canVote":false,"votesFor":27,"votesAgainst":54,"comments":81,"verdict":3},{"id":"id28","createdOn":84,"username":"username28","title":"title28","description":"description28","canVote":false,"votesFor":28,"votesAgainst":56,"comments":84,"verdict":3},{"id":"id29","createdOn":87,"username":"username29","title":"title29","description":"description29","canVote":false,"votesFor":29,"votesAgainst":58,"comments":87,"verdict":3},{"id":"id30","createdOn":90,"username":"username30","title":"title30","description":"description30","canVote":false,"votesFor":30,"votesAgainst":60,"comments":90,"verdict":3},{"id":"id31","createdOn":93,"username":"username31","title":"title31","description":"description31","canVote":false,"votesFor":31,"votesAgainst":62,"comments":93,"verdict":3},{"id":"id32","createdOn":96,"username":"username32","title":"title32","description":"description32","canVote":false,"votesFor":32,"votesAgainst":64,"comments":96,"verdict":3},{"id":"id33","createdOn":99,"username":"username33","title":"title33","description":"description33","canVote":false,"votesFor":33,"votesAgainst":66,"comments":99,"verdict":3},{"id":"id34","createdOn":102,"username":"username34","title":"title34","description":"description34","canVote":false,"votesFor":34,"votesAgainst":68,"comments":102,"verdict":3},{"id":"id35","createdOn":105,"username":"username35","title":"title35","description":"description35","canVote":false,"votesFor":35,"votesAgainst":70,"comments":105,"verdict":3},{"id":"id36","createdOn":108,"username":"username36","title":"title36","description":"description36","canVote":false,"votesFor":36,"votesAgainst":72,"comments":108,"verdict":3},{"id":"id37","createdOn":111,"username":"username37","title":"title37","description":"description37","canVote":false,"votesFor":37,"votesAgainst":74,"comments":111,"verdict":3},{"id":"id38","createdOn":114,"username":"username38","title":"title38","description":"description38","canVote":false,"votesFor":38,"votesAgainst":76,"comments":114,"verdict":3},{"id":"id39","createdOn":117,"username":"username39","title":"title39","description":"description39","canVote":false,"votesFor":39,"votesAgainst":78,"comments":117,"verdict":3},{"id":"id40","createdOn":120,"username":"username40","title":"title40","description":"description40","canVote":false,"votesFor":40,"votesAgainst":80,"comments":120,"verdict":3},{"id":"id41","createdOn":123,"username":"username41","title":"title41","description":"description41","canVote":false,"votesFor":41,"votesAgainst":82,"comments":123,"verdict":3},{"id":"id42","createdOn":126,"username":"username42","title":"title42","description":"description42","canVote":false,"votesFor":42,"votesAgainst":84,"comments":126,"verdict":3},{"id":"id43","createdOn":129,"username":"username43","title":"title43","description":"description43","canVote":false,"votesFor":43,"votesAgainst":86,"comments":129,"verdict":3},{"id":"id44","createdOn":132,"username":"username44","title":"title44","description":"description44","canVote":false,"votesFor":44,"votesAgainst":88,"comments":132,"verdict":3},{"id":"id45","createdOn":135,"username":"username45","title":"title45","description":"description45","canVote":false,"votesFor":45,"votesAgainst":90,"comments":135,"verdict":3},{"id":"id46","createdOn":138,"username":"username46","title":"title46","description":"description46","canVote":false,"votesFor":46,"votesAgainst":92,"comments":138,"verdict":3},{"id":"id47","createdOn":141,"username":"username47","title":"title47","description":"description47","canVote":false,"votesFor":47,"votesAgainst":94,"comments":141,"verdict":3},{"id":"id48","createdOn":144,"username":"username48","title":"title48","description":"description48","canVote":false,"votesFor":48,"votesAgainst":96,"comments":144,"verdict":3},{"id":"id49","createdOn":147,"username":"username49","title":"title49","description":"description49","canVote":false,"votesFor":49,"votesAgainst":98,"comments":147,"verdict":3},{"id":"id50","createdOn":150,"username":"username50","title":"title50","description":"description50","canVote":false,"votesFor":50,"votesAgainst":100,"comments":150,"verdict":3},{"id":"id51","createdOn":153,"username":"username51","title":"title51","description":"description51","canVote":false,"votesFor":51,"votesAgainst":102,"comments":153,"verdict":3},{"id":"id52","createdOn":156,"username":"username52","title":"title52","description":"description52","canVote":false,"votesFor":52,"votesAgainst":104,"comments":156,"verdict":3},{"id":"id53","createdOn":159,"username":"username53","title":"title53","description":"description53","canVote":false,"votesFor":53,"votesAgainst":106,"comments":159,"verdict":3},{"id":"id54","createdOn":162,"username":"username54","title":"title54","description":"description54","canVote":false,"votesFor":54,"votesAgainst":108,"comments":162,"verdict":3},{"id":"id55","createdOn":165,"username":"username55","title":"title55","description":"description55","canVote":false,"votesFor":55,"votesAgainst":110,"comments":165,"verdict":3},{"id":"id56","createdOn":168,"username":"username56","title":"title56","description":"description56","canVote":false,"votesFor":56,"votesAgainst":112,"comments":168,"verdict":3},{"id":"id57","createdOn":171,"username":"username57","title":"title57","description":"description57","canVote":false,"votesFor":57,"votesAgainst":114,"comments":171,"verdict":3},{"id":"id58","createdOn":174,"username":"username58","title":"title58","description":"description58","canVote":false,"votesFor":58,"votesAgainst":116,"comments":174,"verdict":3},{"id":"id59","createdOn":177,"username":"username59","title":"title59","description":"description59","canVote":false,"votesFor":59,"votesAgainst":118,"comments":177,"verdict":3},{"id":"id60","createdOn":180,"username":"username60","title":"title60","description":"description60","canVote":false,"votesFor":60,"votesAgainst":120,"comments":180,"verdict":3},{"id":"id61","createdOn":183,"username":"username61","title":"title61","description":"description61","canVote":false,"votesFor":61,"votesAgainst":122,"comments":183,"verdict":3},{"id":"id62","createdOn":186,"username":"username62","title":"title62","description":"description62","canVote":false,"votesFor":62,"votesAgainst":124,"comments":186,"verdict":3},{"id":"id63","createdOn":189,"username":"username63","title":"title63","description":"description63","canVote":false,"votesFor":63,"votesAgainst":126,"comments":189,"verdict":3},{"id":"id64","createdOn":192,"username":"username64","title":"title64","description":"description64","canVote":false,"votesFor":64,"votesAgainst":128,"comments":192,"verdict":3},{"id":"id65","createdOn":195,"username":"username65","title":"title65","description":"description65","canVote":false,"votesFor":65,"votesAgainst":130,"comments":195,"verdict":3},{"id":"id66","createdOn":198,"username":"username66","title":"title66","description":"description66","canVote":false,"votesFor":66,"votesAgainst":132,"comments":198,"verdict":3},{"id":"id67","createdOn":201,"username":"username67","title":"title67","description":"description67","canVote":false,"votesFor":67,"votesAgainst":134,"comments":201,"verdict":3},{"id":"id68","createdOn":204,"username":"username68","title":"title68","description":"description68","canVote":false,"votesFor":68,"votesAgainst":136,"comments":204,"verdict":3},{"id":"id69","createdOn":207,"username":"username69","title":"title69","description":"description69","canVote":false,"votesFor":69,"votesAgainst":138,"comments":207,"verdict":3},{"id":"id70","createdOn":210,"username":"username70","title":"title70","description":"description70","canVote":false,"votesFor":70,"votesAgainst":140,"comments":210,"verdict":3},{"id":"id71","createdOn":213,"username":"username71","title":"title71","description":"description71","canVote":false,"votesFor":71,"votesAgainst":142,"comments":213,"verdict":3},{"id":"id72","createdOn":216,"username":"username72","title":"title72","description":"description72","canVote":false,"votesFor":72,"votesAgainst":144,"comments":216,"verdict":3},{"id":"id73","createdOn":219,"username":"username73","title":"title73","description":"description73","canVote":false,"votesFor":73,"votesAgainst":146,"comments":219,"verdict":3},{"id":"id74","createdOn":222,"username":"username74","title":"title74","description":"description74","canVote":false,"votesFor":74,"votesAgainst":148,"comments":222,"verdict":3},{"id":"id75","createdOn":225,"username":"username75","title":"title75","description":"description75","canVote":false,"votesFor":75,"votesAgainst":150,"comments":225,"verdict":3},{"id":"id76","createdOn":228,"username":"username76","title":"title76","description":"description76","canVote":false,"votesFor":76,"votesAgainst":152,"comments":228,"verdict":3},{"id":"id77","createdOn":231,"username":"username77","title":"title77","description":"description77","canVote":false,"votesFor":77,"votesAgainst":154,"comments":231,"verdict":3},{"id":"id78","createdOn":234,"username":"username78","title":"title78","description":"description78","canVote":false,"votesFor":78,"votesAgainst":156,"comments":234,"verdict":3},{"id":"id79","createdOn":237,"username":"username79","title":"title79","description":"description79","canVote":false,"votesFor":79,"votesAgainst":158,"comments":237,"verdict":3},{"id":"id80","createdOn":240,"username":"username80","title":"title80","description":"description80","canVote":false,"votesFor":80,"votesAgainst":160,"comments":240,"verdict":3},{"id":"id81","createdOn":243,"username":"username81","title":"title81","description":"description81","canVote":false,"votesFor":81,"votesAgainst":162,"comments":243,"verdict":3},{"id":"id82","createdOn":246,"username":"username82","title":"title82","description":"description82","canVote":false,"votesFor":82,"votesAgainst":164,"comments":246,"verdict":3},{"id":"id83","createdOn":249,"username":"username83","title":"title83","description":"description83","canVote":false,"votesFor":83,"votesAgainst":166,"comments":249,"verdict":3},{"id":"id84","createdOn":252,"username":"username84","title":"title84","description":"description84","canVote":false,"votesFor":84,"votesAgainst":168,"comments":252,"verdict":3},{"id":"id85","createdOn":255,"username":"username85","title":"title85","description":"description85","canVote":false,"votesFor":85,"votesAgainst":170,"comments":255,"verdict":3},{"id":"id86","createdOn":258,"username":"username86","title":"title86","description":"description86","canVote":false,"votesFor":86,"votesAgainst":172,"comments":258,"verdict":3},{"id":"id87","createdOn":261,"username":"username87","title":"title87","description":"description87","canVote":false,"votesFor":87,"votesAgainst":174,"comments":261,"verdict":3},{"id":"id88","createdOn":264,"username":"username88","title":"title88","description":"description88","canVote":false,"votesFor":88,"votesAgainst":176,"comments":264,"verdict":3},{"id":"id89","createdOn":267,"username":"username89","title":"title89","description":"description89","canVote":false,"votesFor":89,"votesAgainst":178,"comments":267,"verdict":3},{"id":"id90","createdOn":270,"username":"username90","title":"title90","description":"description90","canVote":false,"votesFor":90,"votesAgainst":180,"comments":270,"verdict":3},{"id":"id91","createdOn":273,"username":"username91","title":"title91","description":"description91","canVote":false,"votesFor":91,"votesAgainst":182,"comments":273,"verdict":3},{"id":"id92","createdOn":276,"username":"username92","title":"title92","description":"description92","canVote":false,"votesFor":92,"votesAgainst":184,"comments":276,"verdict":3},{"id":"id93","createdOn":279,"username":"username93","title":"title93","description":"description93","canVote":false,"votesFor":93,"votesAgainst":186,"comments":279,"verdict":3},{"id":"id94","createdOn":282,"username":"username94","title":"title94","description":"description94","canVote":false,"votesFor":94,"votesAgainst":188,"comments":282,"verdict":3},{"id":"id95","createdOn":285,"username":"username95","title":"title95","description":"description95","canVote":false,"votesFor":95,"votesAgainst":190,"comments":285,"verdict":3},{"id":"id96","createdOn":288,"username":"username96","title":"title96","description":"description96","canVote":false,"votesFor":96,"votesAgainst":192,"comments":288,"verdict":3},{"id":"id97","createdOn":291,"username":"username97","title":"title97","description":"description97","canVote":false,"votesFor":97,"votesAgainst":194,"comments":291,"verdict":3},{"id":"id98","createdOn":294,"username":"username98","title":"title98","description":"description98","canVote":false,"votesFor":98,"votesAgainst":196,"comments":294,"verdict":3},{"id":"id99","createdOn":297,"username":"username99","title":"title99","description":"description99","canVote":false,"votesFor":99,"votesAgainst":198,"comments":297,"verdict":3}]`
	testProcessor(t, indexRepository, battlesRepository, nil, testSort, "1", "101", expectedResponse)
}

func TestProcessorNoSort(t *testing.T) {
	indexRepository := indexMocks.NewRepository()
	battlesRepository := battlesMocks.NewRepository()
	addBattles(indexRepository, battlesRepository, testSort, false, 3)
	expectedResponse := `[{"id":"id0","createdOn":0,"username":"username0","title":"title0","description":"description0","canVote":false,"votesFor":0,"votesAgainst":0,"comments":0,"verdict":4}]`
	testProcessor(t, indexRepository, battlesRepository, nil, "", "1", "1", expectedResponse)
}

func TestProcessorBadSort(t *testing.T) {
	indexRepository := indexMocks.NewRepository()
	battlesRepository := battlesMocks.NewRepository()
	addBattles(indexRepository, battlesRepository, testSort, false, 3)
	expectedResponse := `[{"id":"id0","createdOn":0,"username":"username0","title":"title0","description":"description0","canVote":false,"votesFor":0,"votesAgainst":0,"comments":0,"verdict":4}]`
	testProcessor(t, indexRepository, battlesRepository, nil, "bad", "1", "1", expectedResponse)
}

func TestProcessorBadIndexEntry(t *testing.T) {
	indexRepository := indexMocks.NewRepository()
	key := fmt.Sprintf("battleIds:%s", testSort)
	indexRepository.SetScore(key, "badId", 0)
	battlesRepository := battlesMocks.NewRepository()
	addBattles(indexRepository, battlesRepository, testSort, false, 1)
	expectedResponse := `[{"id":"id0","createdOn":0,"username":"username0","title":"title0","description":"description0","canVote":false,"votesFor":0,"votesAgainst":0,"comments":0,"verdict":4}]`
	testProcessor(t, indexRepository, battlesRepository, nil, testSort, "1", "2", expectedResponse)
}

func TestProcessorDeletedIndexEntry(t *testing.T) {
	indexRepository := indexMocks.NewRepository()
	battlesRepository := battlesMocks.NewRepository()
	addBattles(indexRepository, battlesRepository, testSort, true, 1)
	expectedResponse := `[]`
	testProcessor(t, indexRepository, battlesRepository, nil, testSort, "1", "2", expectedResponse)
}

func TestProcessorCanVoteFalse(t *testing.T) {
	authContext := &api.AuthContext{
		UserID: "userId1",
	}
	indexRepository := indexMocks.NewRepository()
	battlesRepository := battlesMocks.NewRepository()
	addBattles(indexRepository, battlesRepository, testSort, false, 1)
	expectedResponse := `[{"id":"id0","createdOn":0,"username":"username0","title":"title0","description":"description0","canVote":false,"votesFor":0,"votesAgainst":0,"comments":0,"verdict":4}]`
	testProcessor(t, indexRepository, battlesRepository, authContext, testSort, "1", "2", expectedResponse)
}

func TestProcessorCanVoteTrue(t *testing.T) {
	authContext := &api.AuthContext{
		UserID: "userId2",
	}
	indexRepository := indexMocks.NewRepository()
	battlesRepository := battlesMocks.NewRepository()
	addBattles(indexRepository, battlesRepository, testSort, false, 1)
	expectedResponse := `[{"id":"id0","createdOn":0,"username":"username0","title":"title0","description":"description0","canVote":true,"votesFor":0,"votesAgainst":0,"comments":0,"verdict":4}]`
	testProcessor(t, indexRepository, battlesRepository, authContext, testSort, "1", "2", expectedResponse)
}

func TestProcessorVerdicts(t *testing.T) {
	key := fmt.Sprintf("battleIds:%s", testSort)
	now := time.NowUnix()
	indexRepository := indexMocks.NewRepository()
	battlesRepository := battlesMocks.NewRepository()
	verdictNone := battles.Battle{
		ID:           "id0",
		UserID:       "userId0",
		Username:     "username0",
		Title:        "title0",
		Description:  "description0",
		VotesFor:     100,
		VotesAgainst: 0,
		CreatedOn:    now,
		State:        battles.Active,
	}
	battlesRepository.Add(&verdictNone)
	indexRepository.SetScore(key, verdictNone.ID, 0)
	verdictFor := battles.Battle{
		ID:           "id1",
		UserID:       "userId1",
		Username:     "username1",
		Title:        "title1",
		Description:  "description1",
		VotesFor:     100,
		VotesAgainst: 0,
		CreatedOn:    1,
		State:        battles.Active,
	}
	battlesRepository.Add(&verdictFor)
	indexRepository.SetScore(key, verdictFor.ID, 1)
	verdictAgainst := battles.Battle{
		ID:           "id2",
		UserID:       "userId2",
		Username:     "username2",
		Title:        "title2",
		Description:  "description2",
		VotesFor:     0,
		VotesAgainst: 100,
		CreatedOn:    2,
		State:        battles.Active,
	}
	battlesRepository.Add(&verdictAgainst)
	indexRepository.SetScore(key, verdictAgainst.ID, 2)
	verdictNoDecision1 := battles.Battle{
		ID:           "id3",
		UserID:       "userId3",
		Username:     "username3",
		Title:        "title3",
		Description:  "description3",
		VotesFor:     0,
		VotesAgainst: 0,
		CreatedOn:    3,
		State:        battles.Active,
	}
	battlesRepository.Add(&verdictNoDecision1)
	indexRepository.SetScore(key, verdictNoDecision1.ID, 3)
	verdictNoDecision2 := battles.Battle{
		ID:           "id4",
		UserID:       "userId4",
		Username:     "username4",
		Title:        "title4",
		Description:  "description4",
		VotesFor:     100,
		VotesAgainst: 95,
		CreatedOn:    4,
		State:        battles.Active,
	}
	battlesRepository.Add(&verdictNoDecision2)
	indexRepository.SetScore(key, verdictNoDecision2.ID, 4)
	expectedResponse := `[{"id":"id0","createdOn":` + fmt.Sprintf("%d", now) + `,"username":"username0","title":"title0","description":"description0","canVote":false,"votesFor":100,"votesAgainst":0,"comments":0,"verdict":1},{"id":"id1","createdOn":1,"username":"username1","title":"title1","description":"description1","canVote":false,"votesFor":100,"votesAgainst":0,"comments":0,"verdict":2},{"id":"id2","createdOn":2,"username":"username2","title":"title2","description":"description2","canVote":false,"votesFor":0,"votesAgainst":100,"comments":0,"verdict":3},{"id":"id3","createdOn":3,"username":"username3","title":"title3","description":"description3","canVote":false,"votesFor":0,"votesAgainst":0,"comments":0,"verdict":4},{"id":"id4","createdOn":4,"username":"username4","title":"title4","description":"description4","canVote":false,"votesFor":100,"votesAgainst":95,"comments":0,"verdict":4}]`
	testProcessor(t, indexRepository, battlesRepository, nil, testSort, "1", "5", expectedResponse)
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
			Username:     fmt.Sprintf("username%d", i),
			Title:        fmt.Sprintf("title%d", i),
			Description:  fmt.Sprintf("description%d", i),
			VotesFor:     i,
			VotesAgainst: i * 2,
			Comments:     i * 3,
			CreatedOn:    int64(i * 3),
			State:        state,
		}
		battlesRepository.Add(&battle)
		indexRepository.SetScore(key, battle.ID, float64(i))
	}
}

func testProcessor(t *testing.T, indexRepository *indexMocks.Repository, battlesRepository *battlesMocks.Repository, authContext *api.AuthContext, testSort string, page string, pageSize string, expectedResponseBody string) {
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
	processor := NewProcessor(indexer, battlesRepository, votesRepository)
	output, err := processor.Process(input)
	AssertNil(t, err)
	AssertNotNil(t, output)
	AssertEquals(t, output.StatusCode, http.OK)
	AssertNotNil(t, output.ResponseBody)
	responseBody, _ := json.Marshal(output.ResponseBody)
	AssertEquals(t, string(responseBody), expectedResponseBody)
}
