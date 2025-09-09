terraform {
    required_providers {
        altr = {
        source = "altr"
        }
    }
}

provider "altr" {}

resource "altr_data_source" "example_snowflake" {
  name                 = "example-datasource"
  database_type       = "snowflake"
  friendly_database_name = "example_db"
  hostname            = "example.snowflakecomputing.com"
  database_username   = "example_user"
  database_password   = "example_password"
  database_port       = 443
  warehouse           = "example_warehouse"
  should_classify     = true
  classification_type = "full"
}

resource "altr_data_source" "example_databricks" {
  name                 = "example-datasource-databricks"
  database_type       = "databricks"
  database_port       = 1521
  should_classify     = true
  classification_type = "full"
}

