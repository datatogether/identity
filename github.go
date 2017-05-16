package main

import (
	"context"
	"encoding/json"
	"fmt"
	"golang.org/x/oauth2"
	"net/http"
)

type Github struct {
	client *http.Client
}

func NewGithub(token *oauth2.Token) Github {
	return Github{
		client: githubOAuth.Client(context.Background(), token),
	}
}

func (g Github) ExtractUser() (*User, error) {
	info, err := g.CurrentUserInfo()
	if err != nil {
		return nil, err
	}

	return &User{
		Username: info["username"].(string),
		// TODO - interpret github "type" field
		Type:        UserTypeUser,
		Name:        info["name"].(string),
		Description: info["bio"].(string),
		Email:       info["email"].(string),
	}, nil
}

func (g Github) CurrentUserInfo() (map[string]interface{}, error) {
	res, err := g.client.Get(g.endpoint("/user"))
	if err != nil {
		return nil, err
	}

	info := map[string]interface{}{}
	if err := json.NewDecoder(res.Body).Decode(&info); err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		logger.Println(info)
		return nil, fmt.Errorf("invalid response status code fetching User Info: %d", res.StatusCode)
	}
	return info, nil
}

func (g Github) RepoPermission(org, repo, username string) (string, error) {
	logger.Println(fmt.Sprintf("/repos/%s/%s/collaborators/%s/permission", org, repo, username))
	res, err := g.client.Get(g.endpoint(fmt.Sprintf("/repos/%s/%s/collaborators/%s/permission", org, repo, username)))
	if err != nil {
		return "", err
	}

	perm := map[string]interface{}{}
	if err := json.NewDecoder(res.Body).Decode(&perm); err != nil {
		return "", err
	}

	if perm["permission"] == nil {
		logger.Println(perm)
		return "", err
	}

	return perm["permission"].(string), nil
}

func (g Github) endpoint(path string) string {
	return fmt.Sprintf("https://api.github.com%s", path)
}
