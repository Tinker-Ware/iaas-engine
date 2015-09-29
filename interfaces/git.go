package interfaces

import (
	"encoding/json"
	"fmt"
	"github.com/iaas-engine/domain"
	"io"
	"io/ioutil"
	"bytes"
	"net/http"
	"net/url"
	"errors"
)


type Git struct {
	baseUrl string
	auth    Auth
}

type Auth struct {
	Username string
	ApiToken string
}

func NewGit(baseUrl string) *Git {
	return &Git{
		baseUrl: baseUrl,
	}
}

func (git *Git) SetAuth(user, token string) {
	auth := Auth {
		Username: user,
		ApiToken: token,
	}

	git.auth = auth
}

func (git Git) GetRepo(name string) error {
	return git.get(name, git.auth.ApiToken, git.auth.Username, nil)
}

func (git Git) CreateRepo(name, org string, private bool) error {
	repoRequest := domain.RepoRequest {
		Owner: git.auth.Username,
		Name: name,
		Private: private,
		Org: org,
	}
	
	repoRequestJSON, err := json.Marshal(repoRequest)
	if err != nil {
		fmt.Println(err)
	}
	
	reader := bytes.NewReader(repoRequestJSON)
	
	return git.postJSON("repos", git.auth.ApiToken, git.auth.Username, reader, nil)
}

func (git Git) get(path, token, user string, body interface{}) error {
	client := &http.Client{}

		
	requestUrl, err := url.Parse(git.baseUrl)
	if err != nil {
		return err
	}

	requestUrl.Path = fmt.Sprintf("v1/github/user/%s/%s", user, path)
	req, err := http.NewRequest("GET", requestUrl.String(), nil)
	if err != nil {
		return err
	}
	
	req.Header.Add("token", git.auth.ApiToken)
	
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return errors.New(fmt.Sprintf("error: HTTP GET returned status code returned:%d", resp.StatusCode))
	}

	return git.parseJsonResponse(resp, body)

}

func (git Git) postJSON(path, token, user string, jsonBody io.Reader, body interface{}) error {
	client := &http.Client{}
	
	requestUrl, err := url.Parse(git.baseUrl)
	if err != nil {
		return  err
	}

	requestUrl.Path = fmt.Sprintf("v1/github/user/%s/%s", user, path)
	req, err := http.NewRequest("POST", requestUrl.String(), jsonBody)
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("token", token)

	fmt.Println(requestUrl.String())
	fmt.Println(user)
	fmt.Println(token)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return errors.New(fmt.Sprintf("error: HTTP POST returned status code returned:%d", resp.StatusCode))
	}

	return git.parseJsonResponse(resp, body)

}

func (git Git) parseJsonResponse(resp *http.Response, body interface{}) (err error) {
	defer resp.Body.Close()

	if body == nil {
		return
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	return json.Unmarshal(data, body)
}
