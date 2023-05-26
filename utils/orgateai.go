package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// postBody, _ := []byte()

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

func main(query string, apiKey string) {
	reqBody := formPostCallBody(query, apiKey)
	body, err := json.Marshal(reqBody)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(body)
	res, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	fmt.Print("Coming here")
}
