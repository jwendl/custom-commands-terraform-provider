package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
)

func customCommandsPublishSkills() *schema.Resource {
	return &schema.Resource{
		Create: customCommandsPublishSkillsCreate,
		Read:   customCommandsPublishSkillsRead,
		Update: customCommandsPublishSkillsUpdate,
		Delete: customCommandsPublishSkillsDelete,

		Schema: map[string]*schema.Schema{
			"location": {
				Type:     schema.TypeString,
				Required: true,
			},

			"custom_commands_speech_key": {
				Type:     schema.TypeString,
				Required: true,
			},

			"custom_commands_speech_app_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"skills_file_md5": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
		},
	}
}

func customCommandsPublishSkillsCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[LOG] Starting Custom Commands Publish Skills")

	appID := d.Get("custom_commands_speech_app_id").(string)
	location := d.Get("location").(string)
	basePath := fmt.Sprintf("https://%s.commands.speech.microsoft.com/v1.0/apps/%s/slots/default/languages/en-us/train?force=true", location, appID)

	log.Printf("[LOG] Training the Model")
	customCommandsSkillsJSON, errorMessage := json.Marshal("")
	if errorMessage != nil {
		return fmt.Errorf("An error happened with json.Marshal: %+v", errorMessage)
	}

	apiKey := d.Get("custom_commands_speech_key").(string)
	postResponse, errorMessage := CallWebService(basePath, http.MethodPost, apiKey, 201, customCommandsSkillsJSON)
	if errorMessage != nil {
		return fmt.Errorf("An error happened with CallWebService: %+v", errorMessage)
	}

	statusLocation := postResponse.Header.Get("operation-location")
	type CustomCommandsSkillsResponse struct {
		Status string `json:"status"`
	}

	var customCommandsSkillsResponse CustomCommandsSkillsResponse
	tries := 0
	for true {
		innerResponse, errorMessage := CallWebService(statusLocation, http.MethodGet, apiKey, 200, customCommandsSkillsJSON)
		if innerResponse.StatusCode != 200 && tries > 3 {
			return fmt.Errorf("An error happened with CallWebService after 3 tries: %+v", errorMessage)
		}

		responseData, errorMessage := ioutil.ReadAll(innerResponse.Body)
		if errorMessage != nil {
			return fmt.Errorf("An error happened with ioutil.ReadAll: %+v", errorMessage)
		}

		json.Unmarshal([]byte(responseData), &customCommandsSkillsResponse)
		log.Printf("[LOG] Status from Speech Service: %+v", customCommandsSkillsResponse.Status)

		if customCommandsSkillsResponse.Status == "Succeeded" {
			break
		} else {
			if errorMessage != nil {
				return fmt.Errorf("An error happened with CallWebService: %+v", errorMessage)
			}
		}

		tries++
		time.Sleep(1 * time.Second)
	}

	log.Printf("[LOG] Publishing the Model")

	publishPath := strings.Replace(statusLocation, "/train/", "/publish/", -1)
	putResponse, errorMessage := CallWebService(publishPath, http.MethodPut, apiKey, 204, customCommandsSkillsJSON)
	if errorMessage != nil {
		return fmt.Errorf("An error happened with CallWebService: %+v", errorMessage)
	}

	log.Printf("[LOG] Put Response: %+v", putResponse)
	d.SetId(appID)

	return nil
}

func customCommandsPublishSkillsRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func customCommandsPublishSkillsUpdate(d *schema.ResourceData, m interface{}) error {
	return nil
}

func customCommandsPublishSkillsDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}
