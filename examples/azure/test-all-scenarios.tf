# =========================================================================
# Terraship Azure Test Scenarios - Comprehensive Policy Validation
# =========================================================================
# This file contains test cases for all Terraship policy rules on Azure
# It includes both COMPLIANT and NON-COMPLIANT resources to demonstrate
# how Terraship validates infrastructure against policies.
# 
# NOTE: Provider configuration is in main.tf (don't duplicate here)
# =========================================================================

# =========================================================================
# SCENARIO 1: TAGGING COMPLIANCE
# Tests: required-tags, cost-tagging, naming-convention
# =========================================================================

# ✅ COMPLIANT: Resource Group with all required tags
resource "azurerm_resource_group" "compliant_rg" {
  name     = "terraship-compliant-rg"
  location = "East US"

  tags = {
    Environment = "Production"
    Owner       = "DevOps Team"
    Project     = "Terraship Testing"
    CostCenter  = "Engineering"
    ManagedBy   = "Terraform"
  }
}

# ❌ NON-COMPLIANT: Missing required tags (Environment, Owner, Project)
resource "azurerm_resource_group" "missing_tags_rg" {
  name     = "terraship-notags-rg"
  location = "East US"

  tags = {
    Description = "This resource is missing required tags"
  }
}

# ❌ NON-COMPLIANT: Missing CostCenter tag (cost-tagging rule)
resource "azurerm_resource_group" "missing_costcenter_rg" {
  name     = "terraship-nocost-rg"
  location = "East US"

  tags = {
    Environment = "Development"
    Owner       = "Test Team"
    Project     = "Testing"
    # Missing: CostCenter
  }
}

# ❌ NON-COMPLIANT: Invalid naming convention (uppercase letters)
resource "azurerm_resource_group" "BAD_NAMING_RG" {
  name     = "TerraShip-BAD-NAMING-RG"
  location = "East US"

  tags = {
    Environment = "Development"
    Owner       = "Test Team"
    Project     = "Testing"
    CostCenter  = "Engineering"
  }
}

# =========================================================================
# SCENARIO 2: STORAGE ACCOUNT - ENCRYPTION & SECURITY
# Tests: encryption-at-rest, block-public-access, enable-logging, 
#        azure-storage-https-only
# =========================================================================

# ✅ COMPLIANT: Storage Account with all security best practices
resource "azurerm_storage_account" "compliant_storage" {
  name                     = "terrashipcompliant"
  resource_group_name      = azurerm_resource_group.compliant_rg.name
  location                 = azurerm_resource_group.compliant_rg.location
  account_tier             = "Standard"
  account_replication_type = "GRS"
  
  # Security configurations
  https_traffic_only_enabled      = true
  min_tls_version                 = "TLS1_2"
  allow_nested_items_to_be_public = false
  
  # Encryption (enabled by default in Azure, but explicitly set)
  infrastructure_encryption_enabled = true

  # Network security
  public_network_access_enabled = false

  tags = {
    Environment = "Production"
    Owner       = "Security Team"
    Project     = "Secure Storage"
    CostCenter  = "Security"
  }
}

# ❌ NON-COMPLIANT: Storage Account - Multiple violations
resource "azurerm_storage_account" "insecure_storage" {
  name                     = "terrashipinsecure"
  resource_group_name      = azurerm_resource_group.compliant_rg.name
  location                 = azurerm_resource_group.compliant_rg.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
  
  # ❌ HTTPS not enforced
  https_traffic_only_enabled = false
  
  # ❌ Public access allowed
  allow_nested_items_to_be_public = true
  public_network_access_enabled   = true
  
  # ❌ Missing required tags
  tags = {
    Name = "Insecure Storage"
  }
}

# =========================================================================
# SCENARIO 3: BLOB CONTAINERS - PUBLIC ACCESS & LOGGING
# Tests: block-public-access, enable-logging
# =========================================================================

# ✅ COMPLIANT: Private blob container with logging
resource "azurerm_storage_container" "compliant_container" {
  name                  = "compliant-data"
  storage_account_name  = azurerm_storage_account.compliant_storage.name
  container_access_type = "private"
}

# ❌ NON-COMPLIANT: Public blob container
resource "azurerm_storage_container" "public_container" {
  name                  = "public-data"
  storage_account_name  = azurerm_storage_account.insecure_storage.name
  container_access_type = "blob"  # ❌ Public access
}

# =========================================================================
# SCENARIO 4: MANAGED DISKS - ENCRYPTION
# Tests: encryption-at-rest
# =========================================================================

# ✅ COMPLIANT: Managed disk with encryption
resource "azurerm_managed_disk" "compliant_disk" {
  name                 = "terraship-encrypted-disk"
  location             = azurerm_resource_group.compliant_rg.location
  resource_group_name  = azurerm_resource_group.compliant_rg.name
  storage_account_type = "Premium_LRS"
  create_option        = "Empty"
  disk_size_gb         = 128
  
  # ✅ Encryption enabled by default in Azure
  # Note: encryption_settings is deprecated, encryption is automatic

  tags = {
    Environment = "Production"
    Owner       = "Infrastructure Team"
    Project     = "Secure Compute"
    CostCenter  = "Infrastructure"
  }
}

# ❌ NON-COMPLIANT: Managed disk without encryption
resource "azurerm_managed_disk" "unencrypted_disk" {
  name                 = "terraship-plain-disk"
  location             = azurerm_resource_group.compliant_rg.location
  resource_group_name  = azurerm_resource_group.compliant_rg.name
  storage_account_type = "Standard_LRS"
  create_option        = "Empty"
  disk_size_gb         = 64
  
  # ❌ No encryption settings
  
  tags = {
    Environment = "Development"
    Owner       = "Dev Team"
    Project     = "Testing"
    CostCenter  = "Engineering"
  }
}

# =========================================================================
# SCENARIO 5: SQL DATABASE - ENCRYPTION, BACKUP, PUBLIC ACCESS
# Tests: encryption-at-rest, backup-enabled, block-public-access
# =========================================================================

# SQL Server (required for database)
resource "azurerm_mssql_server" "compliant_sql_server" {
  name                         = "terraship-compliant-sql"
  resource_group_name          = azurerm_resource_group.compliant_rg.name
  location                     = azurerm_resource_group.compliant_rg.location
  version                      = "12.0"
  administrator_login          = "sqladmin"
  administrator_login_password = "P@ssw0rd123!ComplexPassword"
  
  # ✅ Public access disabled
  public_network_access_enabled = false

  tags = {
    Environment = "Production"
    Owner       = "Database Team"
    Project     = "Application DB"
    CostCenter  = "Engineering"
  }
}

# ✅ COMPLIANT: SQL Database with encryption and backup
resource "azurerm_mssql_database" "compliant_db" {
  name      = "terraship-secure-db"
  server_id = azurerm_mssql_server.compliant_sql_server.id
  
  # ✅ Encryption enabled (TDE - Transparent Data Encryption)
  transparent_data_encryption_enabled = true
  
  # ✅ Backup retention configured
  short_term_retention_policy {
    retention_days = 7
  }
  
  long_term_retention_policy {
    weekly_retention  = "P1W"
    monthly_retention = "P1M"
    yearly_retention  = "P1Y"
    week_of_year      = 1
  }

  tags = {
    Environment = "Production"
    Owner       = "Database Team"
    Project     = "Application DB"
    CostCenter  = "Engineering"
  }
}

# ❌ NON-COMPLIANT: SQL Server with public access
resource "azurerm_mssql_server" "insecure_sql_server" {
  name                         = "terraship-public-sql"
  resource_group_name          = azurerm_resource_group.compliant_rg.name
  location                     = azurerm_resource_group.compliant_rg.location
  version                      = "12.0"
  administrator_login          = "sqladmin"
  administrator_login_password = "P@ssw0rd123!ComplexPassword"
  
  # ❌ Public access enabled
  public_network_access_enabled = true

  tags = {
    Environment = "Development"
    Owner       = "Dev Team"
    Project     = "Testing"
    CostCenter  = "Engineering"
  }
}

# ❌ NON-COMPLIANT: Database without backup configuration and missing tags
resource "azurerm_mssql_database" "insecure_db" {
  name      = "terraship-plain-db"
  server_id = azurerm_mssql_server.insecure_sql_server.id
  
  # Note: TDE cannot be disabled on regular Azure SQL databases
  # transparent_data_encryption_enabled = true (default)
  
  # ❌ No backup configuration (violation)
  # ❌ Missing required tags (violation)
  
  tags = {
    Name = "Test Database"
  }
}

# =========================================================================
# SCENARIO 6: VIRTUAL MACHINE - NETWORK SECURITY
# Tests: use-private-subnet, encryption-at-rest, naming-convention
# =========================================================================

# Virtual Network and Subnets
resource "azurerm_virtual_network" "vnet" {
  name                = "terraship-vnet"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.compliant_rg.location
  resource_group_name = azurerm_resource_group.compliant_rg.name

  tags = {
    Environment = "Production"
    Owner       = "Network Team"
    Project     = "Infrastructure"
    CostCenter  = "Infrastructure"
  }
}

# ✅ Private subnet (compliant)
resource "azurerm_subnet" "private_subnet" {
  name                 = "private-subnet"
  resource_group_name  = azurerm_resource_group.compliant_rg.name
  virtual_network_name = azurerm_virtual_network.vnet.name
  address_prefixes     = ["10.0.1.0/24"]
}

# ❌ Public subnet (non-compliant)
resource "azurerm_subnet" "public_subnet" {
  name                 = "public-subnet"
  resource_group_name  = azurerm_resource_group.compliant_rg.name
  virtual_network_name = azurerm_virtual_network.vnet.name
  address_prefixes     = ["10.0.2.0/24"]
}

# Network Interface for compliant VM
resource "azurerm_network_interface" "compliant_nic" {
  name                = "terraship-compliant-nic"
  location            = azurerm_resource_group.compliant_rg.location
  resource_group_name = azurerm_resource_group.compliant_rg.name

  ip_configuration {
    name                          = "internal"
    subnet_id                     = azurerm_subnet.private_subnet.id
    private_ip_address_allocation = "Dynamic"
  }

  tags = {
    Environment = "Production"
    Owner       = "Compute Team"
    Project     = "Application Servers"
    CostCenter  = "Engineering"
  }
}

# ✅ COMPLIANT: VM in private subnet with encryption
# Note: VM resource commented out due to SSH key validation complexity
# In real scenarios, use actual SSH public keys or password authentication
# resource "azurerm_linux_virtual_machine" "compliant_vm" {
#   name                = "terraship-secure-vm"
#   resource_group_name = azurerm_resource_group.compliant_rg.name
#   location            = azurerm_resource_group.compliant_rg.location
#   size                = "Standard_B2s"
#   admin_username      = "adminuser"
#   admin_password      = "P@ssw0rd1234!ComplexPassword"
#   disable_password_authentication = false
#   
#   network_interface_ids = [
#     azurerm_network_interface.compliant_nic.id,
#   ]
#
#   os_disk {
#     caching              = "ReadWrite"
#     storage_account_type = "Premium_LRS"
#   }
#
#   source_image_reference {
#     publisher = "Canonical"
#     offer     = "0001-com-ubuntu-server-focal"
#     sku       = "20_04-lts"
#     version   = "latest"
#   }
#
#   tags = {
#     Environment = "Production"
#     Owner       = "Operations Team"
#     Project     = "Web Application"
#     CostCenter  = "Engineering"
#   }
# }

# =========================================================================
# SCENARIO 7: KEY VAULT - SECURITY & COMPLIANCE
# Tests: encryption-at-rest, enable-logging, block-public-access
# =========================================================================

data "azurerm_client_config" "current" {}

# ✅ COMPLIANT: Key Vault with security best practices
resource "azurerm_key_vault" "compliant_kv" {
  name                       = "terraship-secure-kv"
  location                   = azurerm_resource_group.compliant_rg.location
  resource_group_name        = azurerm_resource_group.compliant_rg.name
  tenant_id                  = data.azurerm_client_config.current.tenant_id
  sku_name                   = "premium"
  
  # ✅ Soft delete and purge protection
  soft_delete_retention_days = 7
  purge_protection_enabled   = true
  
  # ✅ Network security
  public_network_access_enabled = false
  
  # ✅ RBAC enabled
  enable_rbac_authorization = true

  tags = {
    Environment = "Production"
    Owner       = "Security Team"
    Project     = "Secrets Management"
    CostCenter  = "Security"
  }
}

# ❌ NON-COMPLIANT: Key Vault with weak security
resource "azurerm_key_vault" "insecure_kv" {
  name                       = "terraship-weak-kv"
  location                   = azurerm_resource_group.compliant_rg.location
  resource_group_name        = azurerm_resource_group.compliant_rg.name
  tenant_id                  = data.azurerm_client_config.current.tenant_id
  sku_name                   = "standard"
  
  # ❌ No purge protection
  purge_protection_enabled = false
  
  # ❌ Public access enabled
  public_network_access_enabled = true

  tags = {
    Name = "Test Key Vault"
  }
}

# =========================================================================
# SCENARIO 8: ROLE ASSIGNMENTS - IAM BEST PRACTICES
# Tests: iam-least-privilege
# =========================================================================

# ✅ COMPLIANT: Specific role assignment (Reader)
resource "azurerm_role_assignment" "compliant_reader" {
  scope                = azurerm_resource_group.compliant_rg.id
  role_definition_name = "Reader"
  principal_id         = data.azurerm_client_config.current.object_id
}

# ❌ NON-COMPLIANT: Overly permissive role (Owner/Contributor)
resource "azurerm_role_assignment" "overpermissive_owner" {
  scope                = azurerm_resource_group.compliant_rg.id
  role_definition_name = "Owner"  # ❌ Too broad
  principal_id         = data.azurerm_client_config.current.object_id
}

# =========================================================================
# SCENARIO 9: APP SERVICE - SECURITY CONFIGURATION
# Tests: azure-storage-https-only, enable-logging
# =========================================================================

resource "azurerm_service_plan" "app_plan" {
  name                = "terraship-app-plan"
  resource_group_name = azurerm_resource_group.compliant_rg.name
  location            = azurerm_resource_group.compliant_rg.location
  os_type             = "Linux"
  sku_name            = "B1"

  tags = {
    Environment = "Production"
    Owner       = "App Team"
    Project     = "Web App"
    CostCenter  = "Engineering"
  }
}

# ✅ COMPLIANT: App Service with security settings
resource "azurerm_linux_web_app" "compliant_app" {
  name                = "terraship-secure-app"
  resource_group_name = azurerm_resource_group.compliant_rg.name
  location            = azurerm_resource_group.compliant_rg.location
  service_plan_id     = azurerm_service_plan.app_plan.id

  site_config {
    # ✅ HTTPS only
    minimum_tls_version = "1.2"
    ftps_state          = "FtpsOnly"
  }
  
  # ✅ HTTPS redirect
  https_only = true

  tags = {
    Environment = "Production"
    Owner       = "Development Team"
    Project     = "Public Website"
    CostCenter  = "Engineering"
  }
}

# ❌ NON-COMPLIANT: App Service without HTTPS enforcement
resource "azurerm_linux_web_app" "insecure_app" {
  name                = "terraship-insecure-app"
  resource_group_name = azurerm_resource_group.compliant_rg.name
  location            = azurerm_resource_group.compliant_rg.location
  service_plan_id     = azurerm_service_plan.app_plan.id

  site_config {
    minimum_tls_version = "1.0"  # ❌ Old TLS version
  }
  
  # ❌ HTTP allowed
  https_only = false

  tags = {
    Name = "Test App"
  }
}

# =========================================================================
# OUTPUTS - For Testing and Verification
# =========================================================================

output "test_summary" {
  description = "Summary of test scenarios"
  value = {
    compliant_resources = [
      "azurerm_resource_group.compliant_rg",
      "azurerm_storage_account.compliant_storage",
      "azurerm_storage_container.compliant_container",
      "azurerm_managed_disk.compliant_disk",
      "azurerm_mssql_server.compliant_sql_server",
      "azurerm_mssql_database.compliant_db",
      "azurerm_linux_virtual_machine.compliant_vm",
      "azurerm_key_vault.compliant_kv",
      "azurerm_linux_web_app.compliant_app"
    ]
    
    non_compliant_resources = [
      "azurerm_resource_group.missing_tags_rg",
      "azurerm_resource_group.missing_costcenter_rg",
      "azurerm_resource_group.BAD_NAMING_RG",
      "azurerm_storage_account.insecure_storage",
      "azurerm_storage_container.public_container",
      "azurerm_managed_disk.unencrypted_disk",
      "azurerm_mssql_server.insecure_sql_server",
      "azurerm_mssql_database.insecure_db",
      "azurerm_key_vault.insecure_kv",
      "azurerm_role_assignment.overpermissive_owner",
      "azurerm_linux_web_app.insecure_app"
    ]
    
    policy_rules_tested = [
      "required-tags",
      "cost-tagging",
      "naming-convention",
      "encryption-at-rest",
      "block-public-access",
      "enable-logging",
      "backup-enabled",
      "use-private-subnet",
      "azure-storage-https-only",
      "iam-least-privilege"
    ]
  }
}

output "validation_command" {
  description = "Command to validate this configuration"
  value       = "terraship validate . --policy ../../policies/sample-policy.yml --provider azure --output human"
}
