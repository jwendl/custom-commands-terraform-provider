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
			"name": &schema.Schema{
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},

			"location": &schema.Schema{
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},

			"resource_group_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"subscription_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"custom_commands_speech_key": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"custom_commands_speech_luis_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"custom_commands_speech_luis_key": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"custom_commands_speech_luis_location": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"app_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func customCommandsProjectCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[LOG] Starting Custom Commands Project")

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	subscriptionID := d.Get("subscription_id").(string)
	location := d.Get("location").(string)
	basePath := fmt.Sprintf("https://%s.commands.speech.microsoft.com/apps", location)

	type CustomCommandsProjectRequest struct {
		Name                    string `json:"name"`
		Stage                   string `json:"stage"`
		Culture                 string `json:"culture"`
		Description             string `json:"description"`
		SkillEnabled            string `json:"skillEnabled"`
		LuisAuthoringResourceID string `json:"luisAuthoringResourceId"`
		LuisAuthoringKey        string `json:"luisAuthoringKey"`
		LuisAuthoringRegion     string `json:"luisAuthoringRegion"`
	}

	var customCommandsProjectRequest CustomCommandsProjectRequest
	customCommandsProjectRequest.Name = name
	customCommandsProjectRequest.Stage = "default"
	customCommandsProjectRequest.Culture = "en-us"
	customCommandsProjectRequest.Description = "New Speech Project"
	customCommandsProjectRequest.SkillEnabled = "true"
	customCommandsProjectRequest.LuisAuthoringResourceID = d.Get("custom_commands_speech_luis_id").(string)
	customCommandsProjectRequest.LuisAuthoringKey = d.Get("custom_commands_speech_luis_key").(string)
	customCommandsProjectRequest.LuisAuthoringRegion = d.Get("custom_commands_speech_luis_location").(string)

	customCommandsJSON, errorMessage := json.Marshal(customCommandsProjectRequest)
	if errorMessage != nil {
		return fmt.Errorf("An error happened with json.Marshal: %+v", errorMessage)
	}

	log.Printf("[LOG] Json object is: %+v", customCommandsProjectRequest)

	apiKey := d.Get("custom_commands_speech_key").(string)
	var response *http.Response
	tries := 0
	for true {
		innerResponse, errorMessage := CallWebService(basePath, http.MethodPost, apiKey, 200, customCommandsJSON)
		if innerResponse.StatusCode == 200 {
			response = innerResponse
			break
		} else {
			log.Printf("[LOG] Failed attempt %d because of error: %+v", tries, errorMessage)
			if errorMessage != nil && tries > 3 {
				return fmt.Errorf("An error happened with CallWebService after 3 tries: %+v", errorMessage)
			}
		}

		tries++
		time.Sleep(30 * time.Second)
	}

	responseData, errorMessage := ioutil.ReadAll(response.Body)
	if errorMessage != nil {
		return fmt.Errorf("An error happened with ioutil.ReadAll: %+v", errorMessage)
	}

	type CustomCommandsProjectResponse struct {
		AppID string `json:"appId"`
	}
	var customCommandsProjectResponse CustomCommandsProjectResponse
	json.Unmarshal([]byte(responseData), &customCommandsProjectResponse)
	log.Printf("[LOG] AppId from Speech Service: %+v", customCommandsProjectResponse.AppID)

	idString := fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.CognitiveServices/accounts/%s", subscriptionID, resourceGroup, name)
	d.SetId(idString)
	d.Set("app_id", customCommandsProjectResponse.AppID)

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
