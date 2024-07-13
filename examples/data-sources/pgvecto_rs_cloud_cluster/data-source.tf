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


data "pgvecto-rs-cloud_cluster" "test" {
  id         = "7d3c88ec-8147-45b0-a79c-5568a9fd31db"
  account_id = "8364ded2-5580-45c4-a394-edfa582e35a0"
}

output "output" {
  value = data.pgvecto-rs-cloud.test
}
