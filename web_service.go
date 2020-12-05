package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
)

// CallWebService - calls a web service
func CallWebService(basePath string, method string, apiKey string, expectedStatusCode int, data []byte) (*http.Response, error) {
	client := http.Client{}

	var emptyResponse = &http.Response{}

	credentialOptions := azidentity.DefaultAzureCLICredentialOptions()
	tokenProvider, errorMessage := azidentity.NewAzureCLICredential(&credentialOptions)
	if errorMessage != nil {
		return emptyResponse, fmt.Errorf("An error happened getting instance of NewAzureCLICredential: %+v", errorMessage)
	}

	tokenRequestOptions := azcore.TokenRequestOptions{Scopes: []string{"https://management.core.windows.net/"}}
	accessToken, errorMessage := tokenProvider.GetToken(context.TODO(), tokenRequestOptions)
	if errorMessage != nil {
		return emptyResponse, fmt.Errorf("An error happened fetching access token: %+v", errorMessage)
	}

	request, errorMessage := http.NewRequest(method, basePath, bytes.NewBuffer(data))
	if errorMessage != nil {
		return emptyResponse, fmt.Errorf("An error happened with http.NewRequest: %+v", errorMessage)
	}

	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Ocp-Apim-Subscription-Key", apiKey)
	request.Header.Add("Arm-Token", accessToken.Token)
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
