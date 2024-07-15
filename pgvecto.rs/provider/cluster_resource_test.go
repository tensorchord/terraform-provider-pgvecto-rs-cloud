package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccClusterResource(t *testing.T) {
	rName := acctest.RandomWithPrefix("tftest")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccCheckAPIKeyConfigBasic() + testAccClusterResourceConfig(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("pgvecto-rs-cloud_cluster.starter_plan_cluster", "id"),
					resource.TestCheckResourceAttrSet("pgvecto-rs-cloud_cluster.starter_plan_cluster", "last_updated"),
					resource.TestCheckResourceAttr("pgvecto-rs-cloud_cluster.starter_plan_cluster", "cluster_name", rName),
					resource.TestCheckResourceAttr("pgvecto-rs-cloud_cluster.starter_plan_cluster", "account_id", "5c3cb62b-d00b-4dda-85e6-2c0452d50138"),
					resource.TestCheckResourceAttr("pgvecto-rs-cloud_cluster.starter_plan_cluster", "plan", "Enterprise"),
					resource.TestCheckResourceAttr("pgvecto-rs-cloud_cluster.starter_plan_cluster", "server_resource", "aws-m7i-large-2c-8g"),
					resource.TestCheckResourceAttr("pgvecto-rs-cloud_cluster.starter_plan_cluster", "region", "us-east-1"),
					resource.TestCheckResourceAttr("pgvecto-rs-cloud_cluster.starter_plan_cluster", "cluster_provider", "aws"),
					resource.TestCheckResourceAttr("pgvecto-rs-cloud_cluster.starter_plan_cluster", "database_name", "test"),
					resource.TestCheckResourceAttr("pgvecto-rs-cloud_cluster.starter_plan_cluster", "pg_data_disk_size", "5"),
					resource.TestCheckResourceAttr("pgvecto-rs-cloud_cluster.starter_plan_cluster", "status", "Ready"),
				),
			},
			{
				Config: testAccCheckAPIKeyConfigBasic() + testAccClusterResourceUpdateConfig(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("pgvecto-rs-cloud_cluster.starter_plan_cluster", "id"),
					resource.TestCheckResourceAttr("pgvecto-rs-cloud_cluster.starter_plan_cluster", "cluster_name", rName),
					resource.TestCheckResourceAttr("pgvecto-rs-cloud_cluster.starter_plan_cluster", "account_id", "5c3cb62b-d00b-4dda-85e6-2c0452d50138"),
					resource.TestCheckResourceAttr("pgvecto-rs-cloud_cluster.starter_plan_cluster", "plan", "Enterprise"),
					resource.TestCheckResourceAttr("pgvecto-rs-cloud_cluster.starter_plan_cluster", "server_resource", "aws-m7i-large-2c-8g"),
					resource.TestCheckResourceAttr("pgvecto-rs-cloud_cluster.starter_plan_cluster", "region", "us-east-1"),
					resource.TestCheckResourceAttr("pgvecto-rs-cloud_cluster.starter_plan_cluster", "cluster_provider", "aws"),
					resource.TestCheckResourceAttr("pgvecto-rs-cloud_cluster.starter_plan_cluster", "database_name", "test"),
					resource.TestCheckResourceAttr("pgvecto-rs-cloud_cluster.starter_plan_cluster", "pg_data_disk_size", "10"),
					resource.TestCheckResourceAttr("pgvecto-rs-cloud_cluster.starter_plan_cluster", "status", "Ready"),
				),
			},
		},
	})
}

func testAccClusterResourceConfig(name string) string {
	return fmt.Sprintf(`
resource "pgvecto-rs-cloud_cluster" "starter_plan_cluster" {
	cluster_name      = %q
	account_id = "5c3cb62b-d00b-4dda-85e6-2c0452d50138"
	plan              = "Enterprise"
	server_resource   = "aws-m7i-large-2c-8g"
	region            = "us-east-1"
	cluster_provider  = "aws"
	database_name    = "test"
	pg_data_disk_size = "5"
}
`, name)
}

func testAccClusterResourceUpdateConfig(name string) string {
	return fmt.Sprintf(`
resource "pgvecto-rs-cloud_cluster" "starter_plan_cluster" {
	cluster_name      = %q
	account_id = "5c3cb62b-d00b-4dda-85e6-2c0452d50138"
	plan              = "Enterprise"
	server_resource   = "aws-m7i-large-2c-8g"
	region            = "us-east-1"
	cluster_provider  = "aws"
	database_name    = "test"
	pg_data_disk_size = "10"
}
`, name)
}
