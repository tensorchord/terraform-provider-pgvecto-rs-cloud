# Terraform Provider for PGVecto.rs Cloud 

The Terraform Provider for PGVecto.rs Cloud allows Terraform to manage PGVecto.rs Cloud resources. To learn more or you have found a security issue in The Terraform Provider, please [Contact Us](https://discord.gg/KqswhpVgdU).
> [!NOTE]  
> PGVecto.rs is a Postgres extension that enables scalable vector search, allowing you to build powerful similarity-based applications on top of your Postgres database.


## Table of Contents

- [User Guide](#user-guide)
- [Requirements](#requirements)
- [Building The Provider](#building-the-provider)


## User Guide

If you're building the provider, follow the instructions to [install it as a plugin](https://www.terraform.io/docs/plugins/basics.html#installing-plugins). After placing it into your plugins directory, run `terraform init` to initialize the provider.

Your PGVecto.rs Cloud API Key is required to use the Terraform Provider. You can obtain an API Key by signing up for an account at [PGVecto.rs Cloud](https://cloud.pgvecto.rs).

```hcl
provider "pgvecto-rs-cloud" {
  api_key = "<your_api_key>"
}
```
Remember, your API Key should be a protected secret. See how to protect sensitive input variables when setting your API Key this way.

See [PGVecto.rs Cloud Terraform Integration Overview](https://docs.pgvecto.rs/cloud/manage/terraform.html) for more information.

## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.21

## Building The Provider

Clone the repository

```shell
$ git clone git@github.com:<your_org>/terraform-provider-pgvecto-rs-cloud.git
```

Enter the repository directory

```shell
$ cd terraform-provider-pgvecto-rs-cloud
```

Build the provider using the Go `install` command, the install directory depends on the GOPATH environment variable.

```shell
go install .
```

## Adding Dependencies

This provider uses [Go modules](https://github.com/golang/go/wiki/Modules).
Please see the Go documentation for the most up to date information about using Go modules.

To add a new dependency `github.com/author/dependency` to your Terraform provider:

```shell
go get github.com/author/dependency
go mod tidy
```

Then commit the changes to `go.mod` and `go.sum`.

## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (see [Requirements](#requirements) above).

To compile the provider, run `go install`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

To generate or update documentation, run `go generate`.

In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* Acceptance tests create real resources, and often cost money to run.

```shell
PGVECTORS_CLOUD_API_KEY=pgrs-xxxxxxxxxxxx PGVECTORS_CLOUD_API_URL=https://cloud.pgvecto.rs/api/v1 make testacc
```
