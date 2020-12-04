package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
)

// CallWebService - calls a web service
func CallWebService(basePath string, method string, apiKey string, expectedStatusCode int, data []byte) (*http.Response, error) {
	client := http.Client{}

	var emptyResponse = &http.Response{}
	request, errorMessage := http.NewRequest(method, basePath, bytes.NewBuffer(data))
	if errorMessage != nil {
		return emptyResponse, fmt.Errorf("An error happened with http.NewRequest: %+v", errorMessage)
	}

	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Ocp-Apim-Subscription-Key", apiKey)
	request.Header.Add("User-Agent", "golang client")

	requestDump, errorMessage := httputil.DumpRequest(request, true)
	if errorMessage != nil {
		return emptyResponse, fmt.Errorf("An error happened with httputil.DumpRequest: %+v", errorMessage)
	}

	log.Printf("[LOG] Request object %q", requestDump)

	response, errorMessage := client.Do(request)
	if errorMessage != nil {
		return emptyResponse, fmt.Errorf("An error happened with client.Do: %+v", errorMessage)
	}

	responseDump, errorMessage := httputil.DumpResponse(response, true)
	if errorMessage != nil {
		return emptyResponse, fmt.Errorf("An error happened with responseDump: %+v", errorMessage)
	}

	log.Printf("[LOG] Response object %q", responseDump)
	defer response.Body.Close()

	if response.StatusCode != expectedStatusCode {
		return emptyResponse, fmt.Errorf("Bad status code in response: %+v", response)
	}

	return response, nil
}
