terraform {
  required_version = ">= 1.0"
  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "~> 3.0"
    }
  }
}

provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "terraship-example-rg"
  location = "East US"
  
  tags = {
    Environment = "dev"
    Owner       = "platform-team"
  }
}

resource "azurerm_storage_account" "example" {
  name                     = "terrashipexample123"
  resource_group_name      = azurerm_resource_group.example.name
  location                 = azurerm_resource_group.example.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
  
  tags = {
    Environment = "dev"
    Owner       = "platform-team"
  }
}
