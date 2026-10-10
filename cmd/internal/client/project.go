package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"uuid"
)

type Project struct {
	ProjectName  string          `json:"project_name"`
	Public       bool            `json:"public"`
	Metadata     ProjectMetadata `json:"metadata"`
	StorageLimit int64           `json:"storage_limit"`
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

func (c Client) Create(p Project) (int64, error) {
	client := &http.Client{}
	var id int64

	body, err := json.Marshal(p)
	if err != nil {
		return id, fmt.Errorf("json.Marshal(): %v", err)
	}

	fmt.Printf("%v\n", c)

	url := fmt.Sprintf("%s/api/v2.0/projects", c.BaseUrl)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return id, fmt.Errorf("http.NewRequest() POST /projects: %v", err)
	}
	req.SetBasicAuth(c.Username, c.Password)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	reqID := uuid.New().String()
	req.Header.Set("X-Request-Id", reqID)

	fmt.Printf("POST %s, payload: %s\n", url, body)

	resp, err := client.Do(req)
	if err != nil {
		return id, fmt.Errorf("client.Do() POST /projects: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return id, parseError(resp)
	}

	location := resp.Header.Get("location")
	locationItems := strings.Split(location, "/")
	idStr := locationItems[len(locationItems)-1]
	fmt.Printf("Project %s (id: %s) created\n", p.ProjectName, idStr)

	id, err = strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return id, fmt.Errorf("strconv.Atoi(id=%s): %v", idStr, err)
	}

	return id, nil
}
