package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/hashicorp/terraform/helper/schema"
)

func customCommandsSkills() *schema.Resource {
	return &schema.Resource{
		Create: customCommandsSkillsCreate,
		Read:   customCommandsSkillsRead,
		Update: customCommandsSkillsUpdate,
		Delete: customCommandsSkillsDelete,

		Schema: map[string]*schema.Schema{
			"location": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"custom_commands_speech_key": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"custom_commands_speech_app_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"skills_file_path": &schema.Schema{
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"skills_file_md5": &schema.Schema{
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
		},
	}
}

func customCommandsSkillsCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[LOG] Starting Custom Commands Skills")

	appID := d.Get("custom_commands_speech_app_id").(string)
	location := d.Get("location").(string)
	basePath := fmt.Sprintf("https://%s.commands.speech.microsoft.com/v1.0/apps/%s/slots/default/languages/en-us/model", location, appID)

	skillsFilePath := d.Get("skills_file_path").(string)
	fileBytes, errorMessage := ioutil.ReadFile(skillsFilePath)
	if errorMessage != nil {
		return fmt.Errorf("An error happened with ioutil.ReadFile: %+v", errorMessage)
	}

	apiKey := d.Get("custom_commands_speech_key").(string)
	response, errorMessage := CallWebService(basePath, http.MethodPut, apiKey, 200, fileBytes)
	if errorMessage != nil {
		return fmt.Errorf("An error happened with CallWebService: %+v", errorMessage)
	}

	log.Printf("[LOG] Response: %+v", response)
	d.SetId(appID)

	return nil
}

func customCommandsSkillsRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func customCommandsSkillsUpdate(d *schema.ResourceData, m interface{}) error {
	return nil
}

func customCommandsSkillsDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}
