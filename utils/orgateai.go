package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const url = ""

type postBody struct {
	query  string
	apiKey string
}

// API key is working as auth bearer.
// TODO: Should not pass API key everytime
func formPostCallBody(query string, apiKey string) postBody {
	body := postBody{query: query, apiKey: apiKey}
	return body
}

func RunQuery(query string, apiKey string) []byte {
	reqBody := formPostCallBody(query, apiKey)
	body, err := json.Marshal(reqBody)
	if err != nil {
		panic(err)
	}

	fmt.Println(body)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		panic(err)
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("User-Agent", "")
	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		panic(err)
	}

	respBytes, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		panic(err)
	}

	return respBytes
}
