package main

import (
	"github.com/hashicorp/terraform/helper/schema"
)

// Provider the main provider
func Provider() *schema.Provider {
	return &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{
			"azurecc_custom_commands_project": customCommandsProject(),
			"azurecc_custom_commands_skills":  customCommandsSkills(),
			"azurecc_custom_commands_publish": customCommandsPublishSkills(),
		},
	}
}
