package battles

// GetRequest represents a GET request
type GetRequest struct {
	Sort     string `json:"sort"`
	Page     int    `json:"page"`
	PageSize int    `json:"pageSize"`
}

// GetResponse represents a GET response
type GetResponse struct {
	ID           string `json:"id"`
	UserID       string `json:"-"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	VotesFor     int    `json:"votesFor"`
	VotesAgainst int    `json:"votesAgainst"`
	CreatedOn    uint64 `json:"createdOn"`
}

// PostRequest represents a POST request
type PostRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

// Controller is used to handle requests against the "battles" resource
type Controller struct {
	manager ManagerInterface
}

// NewController creates a new Controller instance
func NewController(manager ManagerInterface) *Controller {
	return &Controller{manager}
}

// HandleGet handles GET requests
func (controller *Controller) HandleGet(request GetRequest) ([]GetResponse, error) {
	battles, err := controller.manager.GetPage(request.Sort, request.Page, request.PageSize)
	if err != nil {
		return nil, err
	}
	return controller.toGetResponse(battles), nil
}

// HandlePost handles POST requests
func (controller *Controller) HandlePost(request PostRequest) error {
	_, err := controller.manager.Create(request.Title, request.Description)
	return err
}

func (*Controller) toGetResponse(battles []*Battle) []GetResponse {
	results := make([]GetResponse, 0, len(battles))
	for _, battle := range battles {
		results = append(results, GetResponse(*battle))
	}
	return results
}
