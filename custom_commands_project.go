package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
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
	basePath := fmt.Sprintf("https://%s.commands.speech.microsoft.com", location)
	apiKey := d.Get("custom_commands_speech_key").(string)

	type Details struct {
		Name         string `json:"name"`
		Description  string `json:"description"`
		BaseLanguage string `json:"baseLanguage"`
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

	type CustomCommandsAppRequest struct {
		Details Details `json:"details"`
		Slots   Slots   `json:"slots"`
	}

	var customCommandsAppRequest CustomCommandsAppRequest
	customCommandsAppRequest.Details.Name = name
	customCommandsAppRequest.Details.Description = "New Speech Project"
	customCommandsAppRequest.Details.SkillEnabled = true
	customCommandsAppRequest.Details.BaseLanguage = "en-us"
	customCommandsAppRequest.Slots.Default.Languages.EnUs.LuisResources.AuthoringResourceId = d.Get("custom_commands_speech_luisa_id").(string)
	customCommandsAppRequest.Slots.Default.Languages.EnUs.LuisResources.AuthoringRegion = d.Get("custom_commands_speech_luisa_location").(string)
	customCommandsAppRequest.Slots.Default.Languages.EnUs.LuisResources.PredictionResourceId = d.Get("custom_commands_speech_luisp_id").(string)
	customCommandsAppRequest.Slots.Default.Languages.EnUs.LuisResources.PredictionRegion = d.Get("custom_commands_speech_luisp_location").(string)

	customCommandsProjectJSON, errorMessage := json.Marshal(customCommandsAppRequest)
	if errorMessage != nil {
		return fmt.Errorf("An error happened with json.Marshal: %+v", errorMessage)
	}

	log.Printf("[LOG] Json object is: %s", string(customCommandsProjectJSON))

	appId := uuid.New()
	fullPath := fmt.Sprintf("%s/v1.0/apps/%s", basePath, appId)
	_, appErrorMessage := submitRequest(http.MethodPut, apiKey, fullPath, customCommandsProjectJSON, 201)
	if appErrorMessage != nil {
		return fmt.Errorf("An error happened with ioutil.ReadAll: %+v", appErrorMessage)
	}

	appIdString := fmt.Sprintf("%s", appId)
	d.SetId(appIdString)
	d.Set("app_id", appIdString)

	log.Printf("[LOG] Custom commands project created with app_id: %s", appIdString)
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

func submitRequest(method string, apiKey string, path string, json []byte, expectedStatusCode int) ([]byte, error) {
	var response *http.Response
	tries := 0
	for true {
		innerResponse, errorMessage := CallWebService(path, method, apiKey, expectedStatusCode, json)
		if innerResponse.StatusCode == expectedStatusCode {
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
