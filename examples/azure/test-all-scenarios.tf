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
# ❌ NON-COMPLIANT: VNet missing required tags
resource "azurerm_virtual_network" "vnet" {
  name                = "terraship-vnet"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.compliant_rg.location
  resource_group_name = azurerm_resource_group.compliant_rg.name

  tags = {
    Environment = "Production"
    # Missing: Owner, Project, CostCenter
  }
}

# ✅ COMPLIANT: VNet with all required tags
resource "azurerm_virtual_network" "compliant_vnet" {
  name                = "terraship-compliant-vnet"
  address_space       = ["10.1.0.0/16"]
  location            = azurerm_resource_group.compliant_rg.location
  resource_group_name = azurerm_resource_group.compliant_rg.name

  tags = {
    Environment = "Production"
    Owner       = "Network Team"
    Project     = "Infrastructure"
    CostCenter  = "Infrastructure"
    ManagedBy   = "Terraform"
  }
}

# ✅ Private subnet - compliant
resource "azurerm_subnet" "private_subnet" {
  name                 = "private-subnet"
  resource_group_name  = azurerm_resource_group.compliant_rg.name
  virtual_network_name = azurerm_virtual_network.compliant_vnet.name
  address_prefixes     = ["10.1.1.0/24"]
}

# ❌ Public subnet - non-compliant (uses non-compliant vnet)
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

# ❌ NON-COMPLIANT: App Service Plan missing required tags
resource "azurerm_service_plan" "app_plan" {
  name                = "terraship-app-plan"
  resource_group_name = azurerm_resource_group.compliant_rg.name
  location            = azurerm_resource_group.compliant_rg.location
  os_type             = "Linux"
  sku_name            = "B1"

  tags = {
    Name = "App Plan"
    # Missing: Environment, Owner, Project, CostCenter
  }
}

# ✅ COMPLIANT: App Service Plan with full tags
resource "azurerm_service_plan" "compliant_app_plan" {
  name                = "terraship-compliant-app-plan"
  resource_group_name = azurerm_resource_group.compliant_rg.name
  location            = azurerm_resource_group.compliant_rg.location
  os_type             = "Linux"
  sku_name            = "B1"

  tags = {
    Environment = "Production"
    Owner       = "App Team"
    Project     = "Web App"
    CostCenter  = "Engineering"
    ManagedBy   = "Terraform"
  }
}

# ✅ COMPLIANT: App Service with security settings
resource "azurerm_linux_web_app" "compliant_app" {
  name                = "terraship-secure-app"
  resource_group_name = azurerm_resource_group.compliant_rg.name
  location            = azurerm_resource_group.compliant_rg.location
  service_plan_id     = azurerm_service_plan.compliant_app_plan.id

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
  service_plan_id     = azurerm_service_plan.compliant_app_plan.id

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
# SCENARIO 10: NETWORK SECURITY - NSG COMPLIANCE
# Tests: naming-convention, required-tags
# =========================================================================

# ✅ COMPLIANT: Network Security Group with all required tags
resource "azurerm_network_security_group" "compliant_nsg" {
  name                = "terraship-compliant-nsg"
  location            = azurerm_resource_group.compliant_rg.location
  resource_group_name = azurerm_resource_group.compliant_rg.name

  tags = {
    Environment = "Production"
    Owner       = "Network Security Team"
    Project     = "Network Security"
    CostCenter  = "Infrastructure"
    ManagedBy   = "Terraform"
  }
}

# =========================================================================
# SCENARIO 11: NETWORK SECURITY - SSH/RDP RESTRICTIONS
# Tests: security-group-restrict-ssh-rdp
# =========================================================================

# ✅ COMPLIANT: NSG with restricted SSH/RDP access
resource "azurerm_network_security_rule" "compliant_restrict_ssh" {
  name                        = "DenySSHFromInternet"
  priority                    = 100
  direction                   = "Inbound"
  access                      = "Deny"
  protocol                    = "Tcp"
  source_port_range           = "*"
  destination_port_range      = "22"
  source_address_prefix       = "*"
  destination_address_prefix  = "*"
  resource_group_name         = azurerm_resource_group.compliant_rg.name
  network_security_group_name = azurerm_network_security_group.compliant_nsg.name
}

resource "azurerm_network_security_rule" "compliant_restrict_rdp" {
  name                        = "DenyRDPFromInternet"
  priority                    = 101
  direction                   = "Inbound"
  access                      = "Deny"
  protocol                    = "Tcp"
  source_port_range           = "*"
  destination_port_range      = "3389"
  source_address_prefix       = "*"
  destination_address_prefix  = "*"
  resource_group_name         = azurerm_resource_group.compliant_rg.name
  network_security_group_name = azurerm_network_security_group.compliant_nsg.name
}

# ❌ NON-COMPLIANT: NSG allowing SSH/RDP from anywhere
resource "azurerm_network_security_group" "insecure_nsg" {
  name                = "terraship-insecure-nsg"
  location            = azurerm_resource_group.compliant_rg.location
  resource_group_name = azurerm_resource_group.compliant_rg.name

  tags = {
    Environment = "Development"
    Owner       = "Dev Team"
  }
}

resource "azurerm_network_security_rule" "insecure_allow_ssh" {
  name                        = "AllowSSHFromInternet"
  priority                    = 100
  direction                   = "Inbound"
  access                      = "Allow"  # ❌ Allows SSH from anywhere
  protocol                    = "Tcp"
  source_port_range           = "*"
  destination_port_range      = "22"
  source_address_prefix       = "*"
  destination_address_prefix  = "*"
  resource_group_name         = azurerm_resource_group.compliant_rg.name
  network_security_group_name = azurerm_network_security_group.insecure_nsg.name
}

# =========================================================================
# SCENARIO 12: NETWORK SECURITY - NAT GATEWAY (FOR PRIVATE SUBNETS)
# Tests: nat-gateway-for-private-subnets
# =========================================================================

# ✅ COMPLIANT: Public IP and NAT Gateway for private subnets
resource "azurerm_public_ip" "nat_gateway_ip" {
  name                = "terraship-nat-gateway-ip"
  location            = azurerm_resource_group.compliant_rg.location
  resource_group_name = azurerm_resource_group.compliant_rg.name
  allocation_method   = "Static"
  sku                 = "Standard"

  tags = {
    Environment = "Production"
    Owner       = "Network Team"
    Project     = "NAT Gateway"
    CostCenter  = "Infrastructure"
  }
}

resource "azurerm_nat_gateway" "compliant_nat_gateway" {
  name                    = "terraship-nat-gateway"
  location                = azurerm_resource_group.compliant_rg.location
  resource_group_name     = azurerm_resource_group.compliant_rg.name
  public_ip_address_ids   = [azurerm_public_ip.nat_gateway_ip.id]
  idle_timeout_in_minutes = 10

  tags = {
    Environment = "Production"
    Owner       = "Network Team"
    Project     = "NAT Gateway"
    CostCenter  = "Infrastructure"
  }
}

# ✅ Associate NAT Gateway with private subnet
resource "azurerm_subnet_nat_gateway_association" "compliant_nat_assoc" {
  subnet_id      = azurerm_subnet.private_subnet.id
  nat_gateway_id = azurerm_nat_gateway.compliant_nat_gateway.id
}

# ❌ NON-COMPLIANT: Private subnet without NAT Gateway (direct internet access)
resource "azurerm_subnet" "isolated_subnet" {
  name                 = "isolated-subnet"
  resource_group_name  = azurerm_resource_group.compliant_rg.name
  virtual_network_name = azurerm_virtual_network.vnet.name
  address_prefixes     = ["10.0.3.0/24"]
}

# =========================================================================
# SCENARIO 13: APPLICATION GATEWAY WITH WAF
# Tests: waf-enabled-on-apis
# =========================================================================

# ✅ COMPLIANT: Application Gateway with WAF enabled
resource "azurerm_application_gateway" "compliant_app_gateway" {
  name                = "terraship-app-gateway"
  location            = azurerm_resource_group.compliant_rg.location
  resource_group_name = azurerm_resource_group.compliant_rg.name

  sku {
    name     = "WAF_v2"  # ✅ WAF enabled
    tier     = "WAF_v2"
    capacity = 1
  }

  gateway_ip_configuration {
    name      = "gateway-ip-config"
    subnet_id = azurerm_subnet.private_subnet.id
  }

  frontend_port {
    name = "http"
    port = 80
  }

  frontend_ip_configuration {
    name                 = "frontend-ip"
    public_ip_address_id = azurerm_public_ip.nat_gateway_ip.id
  }

  backend_address_pool {
    name = "backend-pool"
  }

  backend_http_settings {
    name            = "http-settings"
    cookie_based_affinity = "Disabled"
    port            = 80
    protocol        = "Http"
  }

  http_listener {
    name                           = "http-listener"
    frontend_ip_configuration_name = "frontend-ip"
    frontend_port_name             = "http"
    protocol                       = "Http"
  }

  request_routing_rule {
    name               = "routing-rule"
    rule_type          = "Basic"
    http_listener_name = "http-listener"
    backend_address_pool_name = "backend-pool"
    backend_http_settings_name = "http-settings"
  }

  tags = {
    Environment = "Production"
    Owner       = "Security Team"
    Project     = "API Protection"
    CostCenter  = "Infrastructure"
  }
}

# ❌ NON-COMPLIANT: Standard gateway without WAF
resource "azurerm_application_gateway" "insecure_app_gateway" {
  name                = "terraship-standard-gateway"
  location            = azurerm_resource_group.compliant_rg.location
  resource_group_name = azurerm_resource_group.compliant_rg.name

  sku {
    name     = "Standard_v2"  # ❌ No WAF
    tier     = "Standard_v2"
    capacity = 1
  }

  gateway_ip_configuration {
    name      = "gateway-ip-config"
    subnet_id = azurerm_subnet.public_subnet.id
  }

  frontend_port {
    name = "http"
    port = 80
  }

  frontend_ip_configuration {
    name                 = "frontend-ip"
    public_ip_address_id = azurerm_public_ip.nat_gateway_ip.id
  }

  backend_address_pool {
    name = "backend-pool"
  }

  backend_http_settings {
    name            = "http-settings"
    cookie_based_affinity = "Disabled"
    port            = 80
    protocol        = "Http"
  }

  http_listener {
    name                           = "http-listener"
    frontend_ip_configuration_name = "frontend-ip"
    frontend_port_name             = "http"
    protocol                       = "Http"
  }

  request_routing_rule {
    name               = "routing-rule"
    rule_type          = "Basic"
    http_listener_name = "http-listener"
    backend_address_pool_name = "backend-pool"
    backend_http_settings_name = "http-settings"
  }

  tags = {
    Environment = "Development"
    Owner       = "Dev Team"
  }
}

# =========================================================================
# SCENARIO 14: DATABASE HARDENING - DELETE PROTECTION
# Tests: database-delete-protection
# =========================================================================

# ✅ COMPLIANT: MSSQL Database with delete protection (backup long-term retention acts as protection)
# Note: Azure doesn't have explicit "delete protection" like AWS RDS, but using backup retention provides protection
resource "azurerm_mssql_database" "protected_db" {
  name      = "terraship-protected-db"
  server_id = azurerm_mssql_server.compliant_sql_server.id

  long_term_retention_policy {
    weekly_retention  = "P4W"   # ✅ 4 weeks retention
    monthly_retention = "P12M"  # ✅ 12 months retention
    yearly_retention  = "P5Y"   # ✅ 5 years retention
    week_of_year      = 1
  }

  tags = {
    Environment = "Production"
    Owner       = "Database Team"
    Project     = "Protected Database"
    CostCenter  = "Engineering"
  }
}

# ❌ NON-COMPLIANT: Database with minimal retention (delete could leave no backup)
resource "azurerm_mssql_database" "unprotected_db" {
  name      = "terraship-unprotected-db"
  server_id = azurerm_mssql_server.compliant_sql_server.id

  short_term_retention_policy {
    retention_days = 7  # ❌ Only 7 days - minimal protection
  }

  tags = {
    Environment = "Development"
    Owner       = "Dev Team"
  }
}

# =========================================================================
# SCENARIO 15: DATABASE HARDENING - ENHANCED MONITORING
# Tests: database-enhanced-monitoring
# =========================================================================

# ✅ COMPLIANT: MSSQL Database with threat detection and auditing
resource "azurerm_mssql_database_security_alert_policy" "compliant_threat_detection" {
  resource_group_name        = azurerm_resource_group.compliant_rg.name
  server_name                = azurerm_mssql_server.compliant_sql_server.name
  database_name              = azurerm_mssql_database.compliant_db.name
  state                      = "Enabled"  # ✅ Threat detection enabled
  retention_days             = 30
  disabled_alerts            = []
  email_notification_admins  = true
}

# ✅ COMPLIANT: Database auditing policy
resource "azurerm_mssql_database_auditing_policy" "compliant_auditing" {
  database_id                 = azurerm_mssql_database.compliant_db.id
  enabled                     = true  # ✅ Auditing enabled
  storage_endpoint            = azurerm_storage_account.compliant_storage.primary_blob_endpoint
  storage_account_access_key  = azurerm_storage_account.compliant_storage.primary_access_key
  storage_account_access_key_is_secondary = false
  retention_in_days           = 30
}

# ❌ NON-COMPLIANT: Database without threat detection or auditing
resource "azurerm_mssql_database_security_alert_policy" "insecure_threat_detection" {
  resource_group_name        = azurerm_resource_group.compliant_rg.name
  server_name                = azurerm_mssql_server.compliant_sql_server.name
  database_name              = azurerm_mssql_database.insecure_db.name
  state                      = "Disabled"  # ❌ Threat detection disabled
  retention_days             = 0
}

# =========================================================================
# SCENARIO 16: AUDIT & COMPLIANCE - LOG RETENTION
# Tests: log-retention-minimum-90-days
# =========================================================================

# ✅ COMPLIANT: Storage account with diagnostic settings (90+ days retention)
resource "azurerm_monitor_diagnostic_setting" "compliant_diagnostic" {
  name                       = "compliant-diagnostic-setting"
  target_resource_id         = azurerm_storage_account.compliant_storage.id
  log_analytics_workspace_id = null  # Would use Log Analytics in production

  enabled_log {
    category = "StorageRead"
    retention_policy {
      enabled = true
      days    = 90  # ✅ 90 days retention
    }
  }

  metric {
    category = "Transaction"
    retention_policy {
      enabled = true
      days    = 90
    }
  }
}

# ❌ NON-COMPLIANT: Short log retention (less than 90 days)
resource "azurerm_monitor_diagnostic_setting" "insecure_diagnostic" {
  name                       = "insecure-diagnostic-setting"
  target_resource_id         = azurerm_storage_account.insecure_storage.id
  log_analytics_workspace_id = null

  enabled_log {
    category = "StorageRead"
    retention_policy {
      enabled = true
      days    = 7  # ❌ Only 7 days - insufficient
    }
  }

  metric {
    category = "Transaction"
    retention_policy {
      enabled = true
      days    = 7
    }
  }
}

# =========================================================================
# SCENARIO 17: RBAC - LEAST PRIVILEGE WITH SERVICE PRINCIPAL
# Tests: iam-least-privilege, role-based-access-control
# =========================================================================

# ✅ COMPLIANT: Minimal role for service principal (Storage Blob Data Reader)
resource "azurerm_role_assignment" "compliant_minimal_role" {
  scope                = azurerm_storage_account.compliant_storage.id
  role_definition_name = "Storage Blob Data Reader"  # ✅ Minimal permission
  principal_id         = data.azurerm_client_config.current.object_id
}

# ❌ NON-COMPLIANT: Overly permissive role (Contributor)
resource "azurerm_role_assignment" "insecure_contributor_role" {
  scope                = azurerm_storage_account.insecure_storage.id
  role_definition_name = "Contributor"  # ❌ Too many permissions
  principal_id         = data.azurerm_client_config.current.object_id
}

# =========================================================================
# SCENARIO 18: COMPUTE - AUTO-SCALING CONFIGURATION
# Tests: compute-auto-scaling
# =========================================================================

# ✅ COMPLIANT: App Service Plan with autoscale settings
resource "azurerm_monitor_autoscale_setting" "compliant_autoscale" {
  name                = "terraship-autoscale"
  resource_group_name = azurerm_resource_group.compliant_rg.name
  location            = azurerm_resource_group.compliant_rg.location
  target_resource_id  = azurerm_service_plan.compliant_app_plan.id

  profile {
    name = "default"

    capacity {
      default = 1
      minimum = 1
      maximum = 5  # ✅ Auto-scaling enabled
    }

    rule {
      metric_trigger {
        metric_name        = "CpuPercentage"
        metric_resource_id = azurerm_service_plan.compliant_app_plan.id
        time_grain         = "PT1M"
        statistic          = "Average"
        time_window        = "PT5M"
        operator           = "GreaterThan"
        threshold          = 70
      }

      scale_action {
        direction = "Increase"
        type      = "ChangeCount"
        value     = 1
        cooldown  = "PT5M"
      }
    }
  }
}

# ❌ NON-COMPLIANT: No auto-scaling (static capacity)
# Note: Non-compliant version is implicit (no autoscale setting created)

# =========================================================================
# SCENARIO 19: ENCRYPTION - KMS/CMK MANDATORY
# Tests: kms-encryption-mandatory
# =========================================================================

# ✅ COMPLIANT: Storage account with customer-managed key encryption
resource "azurerm_key_vault_key" "storage_key" {
  name         = "storage-key"
  key_vault_id = azurerm_key_vault.compliant_kv.id
  key_type     = "RSA"
  key_size     = 2048
  key_opts     = ["decrypt", "encrypt", "sign", "unwrapKey", "verify", "wrapKey"]
}

resource "azurerm_storage_account_customer_managed_key" "compliant_cmk" {
  storage_account_id        = azurerm_storage_account.compliant_storage.id
  key_vault_id              = azurerm_key_vault.compliant_kv.id
  key_name                  = azurerm_key_vault_key.storage_key.name
  user_assigned_identity_id = null
}

# ❌ NON-COMPLIANT: Storage account with default encryption (not customer-managed)
# This is implicit - the insecure_storage doesn't have customer_managed_key configuration

# =========================================================================
# SCENARIO 20: NETWORKING - COMPREHENSIVE TAGGING & GOVERNANCE
# Tests: comprehensive-resource-tagging, naming-convention
# =========================================================================

# ✅ COMPLIANT: All networking resources properly tagged and named
resource "azurerm_network_interface" "compliant_nic_tagged" {
  name                = "terraship-prod-nic-001"  # ✅ Proper naming
  location            = azurerm_resource_group.compliant_rg.location
  resource_group_name = azurerm_resource_group.compliant_rg.name

  ip_configuration {
    name                          = "internal"
    subnet_id                     = azurerm_subnet.private_subnet.id
    private_ip_address_allocation = "Dynamic"
  }

  tags = {
    Environment = "Production"
    Owner       = "Network Team"
    Project     = "Core Infrastructure"
    CostCenter  = "Infrastructure"
    ManagedBy   = "Terraform"
    DataClass   = "Public"
    BackupPolicy = "Daily"
  }
}

# ❌ NON-COMPLIANT: Resource with missing governance tags
resource "azurerm_network_interface" "minimal_nic" {
  name                = "nic-test"  # ❌ Non-standard naming
  location            = azurerm_resource_group.compliant_rg.location
  resource_group_name = azurerm_resource_group.compliant_rg.name

  ip_configuration {
    name                          = "internal"
    subnet_id                     = azurerm_subnet.public_subnet.id
    private_ip_address_allocation = "Dynamic"
  }

  tags = {
    Name = "Test NIC"
    # Missing: Environment, Owner, Project, CostCenter, ManagedBy, DataClass, BackupPolicy
  }
}

# =========================================================================
# OUTPUTS - For Testing and Verification
# =========================================================================

output "test_summary" {
  description = "Summary of test scenarios - COMPREHENSIVE POLICY COVERAGE"
  value = {
    scenario_count = 20
    
    compliant_resources = [
      "azurerm_resource_group.compliant_rg",
      "azurerm_storage_account.compliant_storage",
      "azurerm_storage_container.compliant_container",
      "azurerm_managed_disk.compliant_disk",
      "azurerm_mssql_server.compliant_sql_server",
      "azurerm_mssql_database.compliant_db",
      "azurerm_key_vault.compliant_kv",
      "azurerm_linux_web_app.compliant_app",
      "azurerm_virtual_network.compliant_vnet",
      "azurerm_service_plan.compliant_app_plan",
      "azurerm_network_security_group.compliant_nsg",
      "azurerm_network_security_rule.compliant_restrict_ssh",
      "azurerm_network_security_rule.compliant_restrict_rdp",
      "azurerm_public_ip.nat_gateway_ip",
      "azurerm_nat_gateway.compliant_nat_gateway",
      "azurerm_subnet_nat_gateway_association.compliant_nat_assoc",
      "azurerm_application_gateway.compliant_app_gateway",
      "azurerm_mssql_database.protected_db",
      "azurerm_mssql_database_security_alert_policy.compliant_threat_detection",
      "azurerm_mssql_database_auditing_policy.compliant_auditing",
      "azurerm_monitor_diagnostic_setting.compliant_diagnostic",
      "azurerm_monitor_autoscale_setting.compliant_autoscale",
      "azurerm_storage_account_customer_managed_key.compliant_cmk",
      "azurerm_network_interface.compliant_nic_tagged",
      "azurerm_role_assignment.compliant_reader",
      "azurerm_role_assignment.compliant_minimal_role"
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
      "azurerm_linux_web_app.insecure_app",
      "azurerm_virtual_network.vnet",
      "azurerm_service_plan.app_plan",
      "azurerm_network_security_group.insecure_nsg",
      "azurerm_network_security_rule.insecure_allow_ssh",
      "azurerm_subnet.isolated_subnet",
      "azurerm_application_gateway.insecure_app_gateway",
      "azurerm_mssql_database.unprotected_db",
      "azurerm_mssql_database_security_alert_policy.insecure_threat_detection",
      "azurerm_monitor_diagnostic_setting.insecure_diagnostic",
      "azurerm_role_assignment.insecure_contributor_role",
      "azurerm_network_interface.minimal_nic"
    ]
    
    policy_coverage = {
      total_policies_available = 41
      policies_tested = 25  # Updated from 9 to 25
      coverage_percentage = "61%"
      
      tested_policies = [
        "required-tags",
        "cost-tagging",
        "naming-convention",
        "encryption-at-rest",
        "block-public-access",
        "enable-logging",
        "backup-enabled",
        "use-private-subnet",
        "azure-storage-https-only",
        "iam-least-privilege",
        "security-group-restrict-ssh-rdp",
        "nat-gateway-for-private-subnets",
        "waf-enabled-on-apis",
        "database-delete-protection",
        "database-enhanced-monitoring",
        "log-retention-minimum-90-days",
        "tls-minimum-version-1-2",
        "compute-auto-scaling",
        "kms-encryption-mandatory",
        "comprehensive-resource-tagging",
        "database-backup-retention-14-days",
        "database-threat-detection",
        "database-auditing-enabled",
        "storage-account-customer-managed-key",
        "role-based-access-control"
      ]
    }
  }
}

output "resource_balance" {
  description = "Balance of compliant vs non-compliant resources"
  value = {
    compliant_count = 26
    non_compliant_count = 21
    total_resources = 47
    balance_note = "Expanded to 47 resources covering 25 policies. 61% coverage of all 41 policies."
  }
}

output "validation_command" {
  description = "Command to validate this configuration"
  value       = "terraship validate . --policy ../../policies/terraship-policy.yml --provider azure --output html"
}
