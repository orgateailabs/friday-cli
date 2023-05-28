package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const url = "https://orgateai.ue.r.appspot.com/query"

type postBody struct {
	query     string
	api_key   string
	db_schema string
}

// API key is working as auth bearer.
// TODO: Should not pass API key everytime
func formPostCallBody(query string, apiKey string, dbSchema string) postBody {
	body := postBody{query: query, api_key: apiKey, db_schema: dbSchema}
	return body
}

func RunQuery(query string, apiKey string, dbSchema string) []byte {
	reqBody := formPostCallBody(query, apiKey, dbSchema)
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
	// req.Header.Add("User-Agent", "")
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
