package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

type Logs struct {
	Results []Log `json:"results"`
}

type Log struct {
	Id        string `json:"id"`
	CreatedAt string `json:"created_at"`
	Message   string `json:"message"`
}

func ListApplicationLogs(lastLines int, projectId string, environmentId string, applicationId string) Logs {
	logs := Logs{}

	if projectId == "" || environmentId == "" || applicationId == "" {
		return logs
	}

	CheckAuthenticationOrQuitWithMessage()

	url := RootURL + "/project/" + projectId + "/environment/" + environmentId + "/application/" + applicationId +
		"/log?size=" + strconv.Itoa(lastLines)

	req, _ := http.NewRequest(http.MethodGet, url, nil)

	req.Header.Set(headerAuthorization, headerValueBearer+GetAuthorizationToken())

	client := http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return logs
	}

	err = CheckHTTPResponse(resp)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	body, _ := ioutil.ReadAll(resp.Body)

	_ = json.Unmarshal(body, &logs)

	return logs
}

func ListApplicationTailLogs(lastLogId string, projectId string, environmentId string, applicationId string) Logs {
	logs := Logs{}

	if projectId == "" || environmentId == "" || applicationId == "" {
		return logs
	}

	CheckAuthenticationOrQuitWithMessage()

	url := RootURL + "/project/" + projectId + "/environment/" + environmentId + "/application/" + applicationId +
		"/log?last_id=" + lastLogId

	req, _ := http.NewRequest(http.MethodGet, url, nil)

	req.Header.Set(headerAuthorization, headerValueBearer+GetAuthorizationToken())

	client := http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return logs
	}

	err = CheckHTTPResponse(resp)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	body, _ := ioutil.ReadAll(resp.Body)

	_ = json.Unmarshal(body, &logs)

	return logs
}
