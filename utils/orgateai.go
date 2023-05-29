package utils

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

const queryUrl = "https://orgateai.ue.r.appspot.com/query"

type resBody struct {
	Data   string `json:"data"`
	Status string `json:"status"`
}

func RunQuery(query string, apiKey string, dbSchema string) (string, error) {
	reqBody, err := json.Marshal(map[string]string{
		"api_key":   apiKey,
		"db_schema": dbSchema,
		"query":     query,
	})

	if err != nil {
		return "", err
	}

	resp, err := http.Post(queryUrl, "application/json", bytes.NewBuffer(reqBody))

	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	respBytes, err := ioutil.ReadAll(resp.Body)

	respBody := &resBody{}
	err = json.Unmarshal(respBytes, respBody)
	if err != nil {
		return "", err
	}
	return respBody.Data, nil
}

// API key is working as auth bearer.
// TODO: Should not pass API key everytime
// func formPostCallBody(query string, apiKey string, dbSchema string) postBody {
// 	body := postBody{query: query, api_key: apiKey, db_schema: dbSchema}
// 	return body
// }
