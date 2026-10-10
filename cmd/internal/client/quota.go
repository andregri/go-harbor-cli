package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"uuid"
)

type QuotaListItem struct {
	Id        int64            `json:"id"`
	Reference Reference        `json:"ref"`
	Hard      map[string]int64 `json:"hard"`
}

type Quota struct {
	Hard map[string]int64 `json:"hard"`
}

type Reference struct {
	Name string `json:"name"`
	Id   int64  `json:"id"`
}

func (c Client) GetQuota(projectId int64) (QuotaListItem, error) {
	client := &http.Client{}
	var quota QuotaListItem

	url, err := url.Parse(fmt.Sprintf("%s/api/v2.0/quotas", c.BaseUrl))
	if err != nil {
		return quota, fmt.Errorf("url.Parse() %s: %v", c.BaseUrl, err)
	}
	query := url.Query()
	query.Add("reference", "project")
	query.Add("reference_id", strconv.FormatInt(projectId, 10))
	encodedUrl := fmt.Sprintf("%s/api/v2.0/quotas?%s", c.BaseUrl, query.Encode())

	req, err := http.NewRequest("GET", encodedUrl, nil)
	if err != nil {
		return quota, fmt.Errorf("http.NewRequest() GET %s: %v", encodedUrl, err)
	}
	req.SetBasicAuth(c.Username, c.Password)

	reqID := uuid.New().String()
	req.Header.Set("X-Request-Id", reqID)
	req.Header.Set("Accept", "application/json")

	fmt.Printf("GET %s\n", encodedUrl)

	resp, err := client.Do(req)
	if err != nil {
		return quota, fmt.Errorf("client.Do() GET %s: %v", encodedUrl, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return quota, parseError(resp)
	}

	quotas := []QuotaListItem{}
	decoder := json.NewDecoder(resp.Body)
	if err = decoder.Decode(&quotas); err != nil {
		return quota, fmt.Errorf("decoder.Decode(response): %v", err)
	}

	if len(quotas) > 1 {
		return quota, fmt.Errorf("received more than 1 quota for project id %s", strconv.FormatInt(projectId, 10))
	}

	return quotas[0], nil
}

func (c Client) PutQuota(id int64, values map[string]int64) error {
	client := &http.Client{}

	quota := Quota{
		Hard: values,
	}

	body, err := json.Marshal(quota)
	if err != nil {
		return fmt.Errorf("Error json.Marshal(): %v", err)
	}

	url := fmt.Sprintf("%s/api/v2.0/quotas/%s", c.BaseUrl, strconv.FormatInt(id, 10))
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("Error http.NewRequest() PUT %s: %v", url, err)
	}
	req.SetBasicAuth(c.Username, c.Password)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	reqID := uuid.New().String()
	req.Header.Set("X-Request-Id", reqID)

	fmt.Printf("PUT %s, payload: %s\n", url, body)

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("Error client.Do() PUT %s: %v", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return parseError(resp)
	}

	quotaBytes, err := json.MarshalIndent(quota, "", "  ")
	if err != nil {
		return fmt.Errorf("json.MarshalIndent(quota): %w", err)
	}
	fmt.Printf("Quota %d updated:\n%s\n", id, string(quotaBytes))
	return nil
}
