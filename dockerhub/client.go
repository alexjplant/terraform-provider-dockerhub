package dockerhub

import (
  "net/http"
  "encoding/json"
  "bytes"
  "fmt"
  "errors"
)

type imageTag struct {
	Creator             int    `json:"creator"`
	ID                  int    `json:"id"`
	ImageID             string `json:"image_id"`
	LastUpdated         string `json:"last_updated"`
	LastUpdater         int `json:"last_updater"`
	LastUpdaterUsername string `json:"last_updater_username"`
	Name                string `json:"name"`
	Repository          int    `json:"repository"`
	FullSize            int    `json:"full_size"`
	V2                  bool   `json:"v2"`
	TagStatus           string `json:"tag_status"`
	TagLastPulled       string `json:"tag_last_pulled"`
	TagLastPushed       string `json:"tag_last_pushed"`
}

type Client struct {
  http http.Client
  token string
}

func (c Client) Auth(username string, password string) {
		login, _ := json.Marshal(map[string]string{"username": username, "password": password})

		res, _ := c.http.Post("https://hub.docker.com/v2/users/login/", "application/json", bytes.NewBuffer(login))
		defer res.Body.Close()
		var resJSON map[string]string

		json.NewDecoder(res.Body).Decode(&resJSON)
    c.token = resJSON["token"]
}

func (c Client) GetImageTags(namespace string, repositoryName string, tagName string) (*imageTag, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("https://hub.docker.com/v2/repositories/%s/%s/tags/%s", namespace, repositoryName, tagName), nil)
	if err != nil {
	  return new(imageTag), err
  }

  if c.token != "" {
    bearer := "Bearer " + c.token 
    req.Header.Add("Authorization", bearer)
  }

	r, err := c.http.Do(req)
	if err != nil {
    return new(imageTag), err	
  }

	defer r.Body.Close()

	if r.StatusCode != 200 {
		return new(imageTag), errors.New(fmt.Sprintf("received HTTP %v when invoking HTTP API", r.StatusCode))
	}

	var j imageTag
	err = json.NewDecoder(r.Body).Decode(&j)
  return &j, err
}
