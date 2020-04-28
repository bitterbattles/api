package elasticsearch

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// Client is an ElasticSearch client
type Client struct {
	baseURL string
}

// NewClient creates a new Client instance
func NewClient(baseURL string) *Client {
	return &Client{
		baseURL: baseURL,
	}
}

// Search runs a query against an index
func (client *Client) Search(index string, requestBody string) ([]string, error) {
	url := fmt.Sprintf("%s/%s/_search", client.baseURL, index)
	request, err := http.NewRequest("GET", url, strings.NewReader(requestBody))
	if err != nil {
		return nil, err
	}
	request.Header.Add("Content-Type", "application/json")
	http := &http.Client{}
	response, err := http.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	responseBodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	result := &Result{}
	err = json.Unmarshal(responseBodyBytes, result)
	if err != nil {
		return nil, err
	}
	if result.Hits == nil || result.Hits.Total == nil || result.Hits.Hits == nil {
		return nil, errors.New("unexpected response from Elasticsearch")
	}
	count := result.Hits.Total.Value
	ids := make([]string, 0, count)
	for _, hit := range result.Hits.Hits {
		if hit != nil {
			ids = append(ids, hit.ID)
		}
	}
	return ids, nil
}

// Update updates a document in an index
func (client *Client) Update(index string, id string, requestBody string) error {
	url := fmt.Sprintf("%s/%s/_update/%s", client.baseURL, index, id)
	request, err := http.NewRequest("POST", url, strings.NewReader(requestBody))
	request.Header.Add("Content-Type", "application/json")
	http := &http.Client{}
	response, err := http.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	return nil
}

// Delete deletes a document from an index
func (client *Client) Delete(index string, id string) error {
	url := fmt.Sprintf("%s/%s/_doc/%s", client.baseURL, index, id)
	request, err := http.NewRequest("DELETE", url, nil)
	http := &http.Client{}
	response, err := http.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	return nil
}
