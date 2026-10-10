package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"uuid"
)

type Project struct {
	ProjectName  string          `json:"project_name"`
	Public       bool            `json:"public"`
	Metadata     ProjectMetadata `json:"metadata"`
	StorageLimit int             `json:"storage_limit"`
}

type ProjectMetadata struct {
	AutoScan string `json:"auto_scan"`
}

type ErrorsResponse struct {
	Errors []ErrorResponse `json:"errors"`
}

type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type Client struct {
	BaseUrl  string
	Username string
	Password string
}

func (c Client) Create(p Project) (string, error) {
	client := &http.Client{}

	body, err := json.Marshal(p)
	if err != nil {
		return "", fmt.Errorf("Error json.Marshal(): %v", err)
	}

	url := fmt.Sprintf("%s/api/v2.0/projects", c.BaseUrl)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return "", fmt.Errorf("Error http.NewRequest() POST /projects: %v", err)
	}
	req.SetBasicAuth(c.Username, c.Password)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	reqID := uuid.New().String()
	req.Header.Set("X-Request-Id", reqID)

	fmt.Printf("POST %s, payload: %s\n", url, body)

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("Error client.Do() POST /projects: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		fmt.Printf("status code: %d, status message: %s\n", resp.StatusCode, resp.Status)

		errorsResp := ErrorsResponse{}
		decoder := json.NewDecoder(resp.Body)
		if err := decoder.Decode(&errorsResp); err != nil {
			return "", fmt.Errorf("decoder.Decode(response): %v", err)
		}
		for _, e := range errorsResp.Errors {
			return "", fmt.Errorf("%s - %s\n", e.Code, e.Message)
		}
	}

	location := resp.Header.Get("location")
	locationItems := strings.Split(location, "/")
	id := locationItems[len(locationItems)-1]
	fmt.Printf("Project %s (id: %s) created\n", p.ProjectName, id)
	return id, nil
}
