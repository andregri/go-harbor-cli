package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"uuid"
)

const (
	USER_MEMBER_ADMIN = 1
)

// type Member struct {
// 	RoleId int64       `json:"role_id"`
// 	User   MemberUser  `json:"member_user"`
// 	Group  MemberGroup `json:"member_group,omitempty"`
// }

// type MemberUser struct {
// 	Id       int64  `json:"user_id,omitempty"`
// 	Username string `json:"username"`
// }

// type MemberGroup struct {
// 	Id          int64  `json:"id,omitempty"`
// 	Name        string `json:"group_name,omitempty"`
// 	Type        int64  `json:"group_type,omitempty"`
// 	LdapGroupDn string `json:"ldap_group_dn,omitempty"`
// }

func (c Client) AddProjectMember(projectName string, member map[string]any) error {
	client := &http.Client{}

	body, err := json.Marshal(member)
	if err != nil {
		return fmt.Errorf("json.Marshal(): %v", err)
	}

	url := fmt.Sprintf("%s/api/v2.0/projects/%s/members", c.BaseUrl, projectName)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("http.NewRequest() POST %s: %v", url, err)
	}
	req.SetBasicAuth(c.Username, c.Password)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	reqID := uuid.New().String()
	req.Header.Set("X-Request-Id", reqID)

	fmt.Printf("POST %s, payload: %s\n", url, body)

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("client.Do() POST %s: %v", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return parseError(resp)
	}

	location := resp.Header.Get("location")
	fmt.Printf("Added user admin: %s\n", location)

	return nil
}
