data "azurerm_subscription" "current" {
}

resource "azurerm_resource_group" "rg" {
    name     = var.resource_group_name
    location = var.resource_group_location
}

resource "azurerm_cognitive_account" "speech" {
    name                = var.custom_speech_account_name
    location            = azurerm_resource_group.rg.location
    resource_group_name = azurerm_resource_group.rg.name
    kind                = "SpeechServices"
    sku_name            = var.custom_speech_account_sku
}

resource "azurerm_cognitive_account" "luis_prediction" {
    name                = var.luis_prediction_name
    location            = var.luis_prediction_location
    resource_group_name = azurerm_resource_group.rg.name
    kind                = "LUIS"
    sku_name            = var.luis_prediction_sku
}

resource "azurerm_cognitive_account" "luis_authoring" {
    name                = var.luis_authoring_name
    location            = var.luis_authoring_location
    resource_group_name = azurerm_resource_group.rg.name
    kind                = "LUIS.Authoring"
    sku_name            = var.luis_authoring_sku
}

resource "azurecc_custom_commands_project" "ccp" {
    name                                 = var.custom_commands_project_name
    location                             = azurerm_resource_group.rg.location
    resource_group_name                  = azurerm_resource_group.rg.name
    subscription_id                      = data.azurerm_subscription.current.subscription_id
    custom_commands_speech_key           = azurerm_cognitive_account.speech.primary_access_key
    custom_commands_speech_luis_id       = azurerm_cognitive_account.luis_authoring.id
    custom_commands_speech_luis_key      = azurerm_cognitive_account.luis_authoring.primary_access_key
    custom_commands_speech_luis_location = azurerm_cognitive_account.luis_authoring.location
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
