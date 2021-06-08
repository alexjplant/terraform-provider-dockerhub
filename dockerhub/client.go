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
  Images              []image `json:"images"`
}

type image struct {
  Architecture  string  `json:"architecture"`
  Features  string `json:"features"`
  Variant string `json:"variant"`
  Digest string `json:"string"`
  OS string `json:"os"`
  OSFeatures string `json:"os_features"`
  OSVersion string `json:"os_version"`
  Size  int `json:"size"`
  Status  string   `json:"status"`
  LastPulled string `json:"last_pulled"`
  LastPushed  string  `json:"last_pushed"`
}

type Client struct {
  http *http.Client
  token string
}

func NewClient() *Client {
  c := new (Client)
  c.http = &http.Client{}
  return c
}

func (c Client) Auth(username string, password string) error {
		login, _ := json.Marshal(map[string]string{"username": username, "password": password})

		res, err := c.http.Post("https://hub.docker.com/v2/users/login/", "application/json", bytes.NewBuffer(login))
		if err != nil {
      return fmt.Errorf("error invoking login API: %w", err)
    }

    defer res.Body.Close()
		var resJSON map[string]string

		json.NewDecoder(res.Body).Decode(&resJSON)
    c.token = resJSON["token"]

    if c.token == "" {
      return errors.New("unable to parse token from login")
    }

    return nil
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
