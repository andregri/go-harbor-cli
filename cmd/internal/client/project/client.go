package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Project struct {
	Name         string
	QuotaSizeGiB int
}

type ProjectResponse struct {
	Id      string
	QuotaId string
}

type Client struct {
	BaseUrl  string
	Username string
	Password string
}

func (c Client) Create(p Project) error {
	client := &http.Client{}

	body, err := json.Marshal(p)
	if err != nil {
		return fmt.Errorf("Error json.Marshal(): %v", err)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v2/projects", c.BaseUrl), bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("Error http.NewRequest() POST /projects: %v", err)
	}
	req.SetBasicAuth(c.Username, c.Password)

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("Error client.Do() POST /projects: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusCreated {
		projResp := &ProjectResponse{}
		decoder := json.NewDecoder(resp.Body)
		if err = decoder.Decode(&projResp); err != nil {
			return fmt.Errorf("Error decoder.Decode() POST /projects: %v", err)
		}
	}

	return nil
}
