# Custom Commands Speech Service Terraform Provider

This terraform provider will enable the ability to upload and publish new speech models using Terraform to the [Custom Commands Speech Service](https://docs.microsoft.com/en-us/azure/cognitive-services/speech-service/custom-commands)

> Note: This is a temporary home for this functionality until we can identify where it should land.

## Example Usage

``` hcl
data "azurerm_subscription" "current" {
}

resource "azurerm_cognitive_account" "speech" {
  name                = var.custom_speech_account_name
  location            = var.cc_resource_group_location
  resource_group_name = var.cc_resource_group_name
  kind                = "SpeechServices"
  sku_name            = var.custom_speech_account_sku
}

resource "azurerm_cognitive_account" "luis_prediction" {
  name                = var.luis_prediction_name
  location            = var.luis_prediction_location
  resource_group_name = var.cc_resource_group_name
  kind                = "LUIS"
  sku_name            = var.luis_prediction_sku
}

resource "azurerm_cognitive_account" "luis_authoring" {
  name                = var.luis_authoring_name
  location            = var.luis_authoring_location
  resource_group_name = var.cc_resource_group_name
  kind                = "LUIS.Authoring"
  sku_name            = var.luis_authoring_sku
}

resource "azurecc_custom_commands_project" "ccp" {
  name                                  = var.custom_commands_project_name
  location                              = var.cc_resource_group_location
  resource_group_name                   = var.cc_resource_group_name
  subscription_id                       = data.azurerm_subscription.current.subscription_id
  custom_commands_speech_key            = azurerm_cognitive_account.speech.primary_access_key
  custom_commands_speech_luisa_id       = azurerm_cognitive_account.luis_authoring.id
  custom_commands_speech_luisa_key      = azurerm_cognitive_account.luis_authoring.primary_access_key
  custom_commands_speech_luisa_location = azurerm_cognitive_account.luis_authoring.location
  custom_commands_speech_luisp_id       = azurerm_cognitive_account.luis_prediction.id
  custom_commands_speech_luisp_key      = azurerm_cognitive_account.luis_prediction.primary_access_key
  custom_commands_speech_luisp_location = azurerm_cognitive_account.luis_prediction.location
}

resource "azurecc_custom_commands_skills" "ccs" {
  location                      = azurecc_custom_commands_project.ccp.location
  custom_commands_speech_key    = azurerm_cognitive_account.speech.primary_access_key
  custom_commands_speech_app_id = azurecc_custom_commands_project.ccp.app_id
  skills_file_path              = var.custom_commands_skills_file_path
  skills_file_md5               = filemd5(var.custom_commands_skills_file_path)
}

resource "azurecc_custom_commands_publish" "ccpub" {
  location                      = azurecc_custom_commands_skills.ccs.location
  custom_commands_speech_key    = azurerm_cognitive_account.speech.primary_access_key
  custom_commands_speech_app_id = azurecc_custom_commands_project.ccp.app_id
  skills_file_md5               = filemd5(var.custom_commands_skills_file_path)
}

```

## Project Information and Structure

This project was built using GOLANG version 1.15.

The goal of this was to utilize the REST calls from the following [PowerShell script](https://github.com/Azure-Samples/Cognitive-Services-Voice-Assistant/blob/master/custom-commands/demos/scripts/deployCustomCommands.ps1) as a Terraform module instead. This buys us the ability to use Infrastructure as Code techniques and makes it easier to add to the overall infrastructure deployment. 

## Project Structure

| File                       | Description                                            |
| -------------------------- | ------------------------------------------------------ |
| main.go                    | The main entry point for the provider                  |
| provider.go                | Describes what commands are available for the provider |
| web_service.go             | A centralized spot to handle the REST calls            |
| custom_commands_project.go | Creates the custom commands project                    |
| custom_commands_skills.go  | Uploads the skills.json to the custom commands project |
| custom_commands_publish.go | Trains and Publishes the model                         |
 
 ## Building Locally

The following commands are how I built the project locally

``` bash
go build -o build/terraform-provider-azurecc
```

Once built, it needs to be put into the proper folder for terraform to use. This folder would be terraform/root/terraform.d/plugins/registry.terraform.io/hashicorp/azurecc/0.0.1/linux_amd64. Additionally, the file needs to be named a specific format. Similar to terraform-provider-azurecc_v0.0.1_x1 inside that folder. 

To create all the folders needed, I ran the following commands:

``` bash
mkdir build
mkdir ../../../deploy/terraform/root/terraform.d/
mkdir ../../../deploy/terraform/root/terraform.d/plugins/
mkdir ../../../deploy/terraform/root/terraform.d/plugins/registry.terraform.io/
mkdir ../../../deploy/terraform/root/terraform.d/plugins/registry.terraform.io/hashicorp/
mkdir ../../../deploy/terraform/root/terraform.d/plugins/registry.terraform.io/hashicorp/azurecc/
mkdir ../../../deploy/terraform/root/terraform.d/plugins/registry.terraform.io/hashicorp/azurecc/0.0.1/
mkdir ../../../deploy/terraform/root/terraform.d/plugins/registry.terraform.io/hashicorp/azurecc/0.0.1/linux_amd64
```

So the command I ran from this folder is:

``` bash
go build -o build/terraform-provider-azurecc && cp build/terraform-provider-azurecc ../../../deploy/terraform/root/terraform.d/plugins/registry.terraform.io/hashicorp/azurecc/0.0.1/linux_amd64
```

## Q&A

q: Why not use the PowerShell from terraform?
a: This was mainly because when you run a PowerShell script inside terraform, it doesn't track whether it's created the resource or not and doesn't support updates vs. creates. So it would always perform a "create" and thus making it not friendly to run in a pipeline.

q: Why is it calling REST calls?
a: Because the code to manage these resources hasn't been updated in the [Azure GO SDK](https://github.com/Azure/azure-sdk-for-go) yet. Once all of the functionality is available there, then we can create an issue with the [Azure Terraform Provider](https://github.com/terraform-providers/terraform-provider-azurerm) project and have it built as an offically supported module.
