terraform {
  required_providers {
    pgvecto-rs-cloud = {
      source = "tensorchord/pgvecto-rs-cloud"
    }
  }
}

provider "pgvecto-rs-cloud" {
  api_key = "pgrs-xxxxxxxxxxxxxxxx"
}

resource "pgvecto-rs-cloud_cluster" "starter_plan_cluster" {
  cluster_name      = "starter-plan-cluster"
  account_id        = "8364ded2-5580-45c4-a394-edfa582e35a0"
  plan              = "Starter"
  server_resource   = "aws-t3-xlarge-4c-16g"
  region            = "us-east-1"
  cluster_provider  = "aws"
  database_name     = "test"
  pg_data_disk_size = "5"
}


resource "pgvecto-rs-cloud_cluster" "enterprise_plan_cluster" {
  account_id        = "8364ded2-5580-45c4-a394-edfa582e35a0"
  cluster_name      = "starter-plan-cluster"
  plan              = "Enterprise"
  server_resource   = "aws-m7i-large-2c-8g"
  region            = "eu-west-1"
  cluster_provider  = "aws"
  database_name     = "test"
  pg_data_disk_size = "10"
}
