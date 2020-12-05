package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
)

func customCommandsProject() *schema.Resource {
	return &schema.Resource{
		Create: customCommandsProjectCreate,
		Read:   customCommandsProjectRead,
		Update: customCommandsProjectUpdate,
		Delete: customCommandsProjectDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},

			"location": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},

			"resource_group_name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"subscription_id": {
				Type:     schema.TypeString,
				Required: true,
			},

			"custom_commands_speech_key": {
				Type:     schema.TypeString,
				Required: true,
			},

			"custom_commands_speech_luisa_id": {
				Type:     schema.TypeString,
				Required: true,
			},

			"custom_commands_speech_luisa_key": {
				Type:     schema.TypeString,
				Required: true,
			},

			"custom_commands_speech_luisa_location": {
				Type:     schema.TypeString,
				Required: true,
			},

			"custom_commands_speech_luisp_id": {
				Type:     schema.TypeString,
				Required: true,
			},

			"custom_commands_speech_luisp_key": {
				Type:     schema.TypeString,
				Required: true,
			},

			"custom_commands_speech_luisp_location": {
				Type:     schema.TypeString,
				Required: true,
			},

			"app_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func customCommandsProjectCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[LOG] Starting Custom Commands Project")

	name := d.Get("name").(string)
	location := d.Get("location").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	subscriptionID := d.Get("subscription_id").(string)
	basePath := fmt.Sprintf("https://%s.commands.speech.microsoft.com", location)
	apiKey := d.Get("custom_commands_speech_key").(string)

	type CustomCommandsRequest struct {
		Name                    string `json:"name"`
		Stage                   string `json:"stage"`
		Culture                 string `json:"culture"`
		Description             string `json:"description"`
		SkillEnabled            string `json:"skillEnabled"`
		LuisAuthoringResourceID string `json:"luisAuthoringResourceId"`
		LuisAuthoringKey        string `json:"luisAuthoringKey"`
		LuisAuthoringRegion     string `json:"luisAuthoringRegion"`
	}

	var customCommandsRequest CustomCommandsRequest
	customCommandsRequest.Name = name
	customCommandsRequest.Stage = "default"
	customCommandsRequest.Culture = "en-us"
	customCommandsRequest.Description = "New Speech Project"
	customCommandsRequest.SkillEnabled = "true"
	customCommandsRequest.LuisAuthoringResourceID = d.Get("custom_commands_speech_luisa_id").(string)
	customCommandsRequest.LuisAuthoringKey = d.Get("custom_commands_speech_luisa_key").(string)
	customCommandsRequest.LuisAuthoringRegion = d.Get("custom_commands_speech_luisa_location").(string)

	customCommandsJSON, errorMessage := json.Marshal(customCommandsRequest)
	if errorMessage != nil {
		return fmt.Errorf("An error happened with json.Marshal: %+v", errorMessage)
	}

	log.Printf("[LOG] Json object is: %s", string(customCommandsJSON))

	baseFullPath := fmt.Sprintf("%s/apps", basePath)
	baseResponseData, baseErrorMessage := submitRequest(http.MethodPost, apiKey, baseFullPath, customCommandsJSON)
	if baseErrorMessage != nil {
		return fmt.Errorf("An error happened with ioutil.ReadAll: %+v", baseErrorMessage)
	}

	type CustomCommandsResponse struct {
		AppID string `json:"appId"`
	}
	var customCommandsResponse CustomCommandsResponse
	json.Unmarshal([]byte(baseResponseData), &customCommandsResponse)
	log.Printf("[LOG] AppId from Speech Service: %+v", customCommandsResponse.AppID)

	type Details struct {
		Name         string `json:"name"`
		Description  string `json:"description"`
		SkillEnabled bool   `json:"skillEnabled"`
	}

	type LuisResources struct {
		AuthoringResourceId  string `json:"authoringResourceId"`
		AuthoringRegion      string `json:"authoringRegion"`
		PredictionResourceId string `json:"predictionResourceId"`
		PredictionRegion     string `json:"predictionRegion"`
	}

	type DefaultLanguage struct {
		LuisResources LuisResources `json:"luisResources"`
		DialogModel   []string      `json:"dialogModel"`
	}

	type Languages struct {
		EnUs DefaultLanguage `json:"en-us"`
	}

	type DefaultSlot struct {
		Languages Languages `json:"languages"`
	}

	type Slots struct {
		Default DefaultSlot `json:"default"`
	}

	type CustomCommandsProjectRequest struct {
		Details Details `json:"details"`
		Slots   Slots   `json:"slots"`
	}

	var customCommandsProjectRequest CustomCommandsProjectRequest
	customCommandsProjectRequest.Details.Name = name
	customCommandsProjectRequest.Details.Description = "New Speech Project"
	customCommandsProjectRequest.Details.SkillEnabled = true
	customCommandsProjectRequest.Slots.Default.Languages.EnUs.LuisResources.AuthoringResourceId = d.Get("custom_commands_speech_luisa_id").(string)
	customCommandsProjectRequest.Slots.Default.Languages.EnUs.LuisResources.AuthoringRegion = d.Get("custom_commands_speech_luisa_location").(string)
	customCommandsProjectRequest.Slots.Default.Languages.EnUs.LuisResources.PredictionResourceId = d.Get("custom_commands_speech_luisp_id").(string)
	customCommandsProjectRequest.Slots.Default.Languages.EnUs.LuisResources.PredictionRegion = d.Get("custom_commands_speech_luisp_location").(string)

	customCommandsProjectJSON, errorMessage := json.Marshal(customCommandsProjectRequest)
	if errorMessage != nil {
		return fmt.Errorf("An error happened with json.Marshal: %+v", errorMessage)
	}

	log.Printf("[LOG] Json object is: %s", string(customCommandsProjectJSON))

	appsFullPath := fmt.Sprintf("%s/v1.0/apps/%s", basePath, customCommandsResponse.AppID)
	_, appErrorMessage := submitRequest(http.MethodPut, apiKey, appsFullPath, customCommandsProjectJSON)
	if appErrorMessage != nil {
		return fmt.Errorf("An error happened with ioutil.ReadAll: %+v", appErrorMessage)
	}

	idString := fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.CognitiveServices/accounts/%s", subscriptionID, resourceGroup, name)
	d.SetId(idString)
	d.Set("app_id", customCommandsResponse.AppID)

	type CustomCommandsLuisRequest struct {
		AuthoringResourceId  string `json:"authoringResourceId"`
		AuthoringRegion      string `json:"authoringRegion"`
		PredictionResourceId string `json:"predictionResourceId"`
		PredictionRegion     string `json:"predictionRegion"`
	}

	var customCommandsLuisRequest CustomCommandsLuisRequest
	customCommandsLuisRequest.AuthoringResourceId = d.Get("custom_commands_speech_luisa_id").(string)
	customCommandsLuisRequest.AuthoringRegion = d.Get("custom_commands_speech_luisa_location").(string)
	customCommandsLuisRequest.PredictionResourceId = d.Get("custom_commands_speech_luisp_id").(string)
	customCommandsLuisRequest.PredictionRegion = d.Get("custom_commands_speech_luisp_location").(string)

	customCommandsLuisJSON, errorMessage := json.Marshal(customCommandsLuisRequest)
	if errorMessage != nil {
		return fmt.Errorf("An error happened with json.Marshal: %+v", errorMessage)
	}

	log.Printf("[LOG] Json object is: %s", string(customCommandsLuisJSON))

	luisFullPath := fmt.Sprintf("%s/v1.0/apps/%s/slots/default/languages/en-us/luisResources", basePath, customCommandsResponse.AppID)
	_, luisErrorMessage := submitRequest(http.MethodPut, apiKey, luisFullPath, customCommandsLuisJSON)
	if luisErrorMessage != nil {
		return fmt.Errorf("An error happened with ioutil.ReadAll: %+v", luisErrorMessage)
	}

	return nil
}

func customCommandsProjectRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func customCommandsProjectUpdate(d *schema.ResourceData, m interface{}) error {
	return nil
}

func customCommandsProjectDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}

func submitRequest(method string, apiKey string, path string, json []byte) ([]byte, error) {
	var response *http.Response
	tries := 0
	for true {
		innerResponse, errorMessage := CallWebService(path, method, apiKey, 200, json)
		if innerResponse.StatusCode == 200 {
			response = innerResponse
			break
		} else {
			log.Printf("[LOG] Failed attempt %d because of error: %+v", tries, errorMessage)
			if errorMessage != nil && tries > 3 {
				return nil, fmt.Errorf("An error happened with CallWebService after 3 tries: %+v", errorMessage)
			}
		}

		tries++
		time.Sleep(30 * time.Second)
	}

	return ioutil.ReadAll(response.Body)
}
