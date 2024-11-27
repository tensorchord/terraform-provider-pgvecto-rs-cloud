package provider

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccClusterResource(t *testing.T) {
	rName := acctest.RandomWithPrefix("tftest")
	backupID := os.Getenv("BACKUP_ID")
	testBackup := true
	if backupID == "" {
		testBackup = false
	}

	testPITR := true
	clusterID := os.Getenv("CLUSTER_ID")
	targetTime := os.Getenv("TARGET_TIME")
	if clusterID == "" {
		testPITR = false
	}
	if targetTime == "" {
		testPITR = false
	}

	steps := []resource.TestStep{
		// Read testing
		{
			Config: testAccCheckAPIKeyConfigBasic() + testAccClusterResourceConfig(rName),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrSet("pgvecto-rs-cloud_cluster.enterprise_plan_cluster", "id"),
				resource.TestCheckResourceAttrSet("pgvecto-rs-cloud_cluster.enterprise_plan_cluster", "last_updated"),
				resource.TestCheckResourceAttr("pgvecto-rs-cloud_cluster.enterprise_plan_cluster", "cluster_name", rName),
				resource.TestCheckResourceAttr("pgvecto-rs-cloud_cluster.enterprise_plan_cluster", "account_id", "5c3cb62b-d00b-4dda-85e6-2c0452d50138"),
				resource.TestCheckResourceAttr("pgvecto-rs-cloud_cluster.enterprise_plan_cluster", "plan", "Enterprise"),
				resource.TestCheckResourceAttr("pgvecto-rs-cloud_cluster.enterprise_plan_cluster", "server_resource", "aws-m7i-large-2c-8g"),
				resource.TestCheckResourceAttr("pgvecto-rs-cloud_cluster.enterprise_plan_cluster", "region", "us-east-1"),
				resource.TestCheckResourceAttr("pgvecto-rs-cloud_cluster.enterprise_plan_cluster", "cluster_provider", "aws"),
				resource.TestCheckResourceAttr("pgvecto-rs-cloud_cluster.enterprise_plan_cluster", "database_name", "test"),
				resource.TestCheckResourceAttr("pgvecto-rs-cloud_cluster.enterprise_plan_cluster", "pg_data_disk_size", "5"),
				resource.TestCheckResourceAttr("pgvecto-rs-cloud_cluster.enterprise_plan_cluster", "status", "Ready"),
				resource.TestCheckResourceAttr("pgvecto-rs-cloud_cluster.enterprise_plan_cluster", "enable_pooler", "true"),
			),
		},
		{
			Config: testAccCheckAPIKeyConfigBasic() + testAccClusterResourceUpdateConfig(rName),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrSet("pgvecto-rs-cloud_cluster.enterprise_plan_cluster", "id"),
				resource.TestCheckResourceAttr("pgvecto-rs-cloud_cluster.enterprise_plan_cluster", "cluster_name", rName),
				resource.TestCheckResourceAttr("pgvecto-rs-cloud_cluster.enterprise_plan_cluster", "account_id", "5c3cb62b-d00b-4dda-85e6-2c0452d50138"),
				resource.TestCheckResourceAttr("pgvecto-rs-cloud_cluster.enterprise_plan_cluster", "plan", "Enterprise"),
				resource.TestCheckResourceAttr("pgvecto-rs-cloud_cluster.enterprise_plan_cluster", "server_resource", "aws-m7i-large-2c-8g"),
				resource.TestCheckResourceAttr("pgvecto-rs-cloud_cluster.enterprise_plan_cluster", "region", "us-east-1"),
				resource.TestCheckResourceAttr("pgvecto-rs-cloud_cluster.enterprise_plan_cluster", "cluster_provider", "aws"),
				resource.TestCheckResourceAttr("pgvecto-rs-cloud_cluster.enterprise_plan_cluster", "database_name", "test"),
				resource.TestCheckResourceAttr("pgvecto-rs-cloud_cluster.enterprise_plan_cluster", "pg_data_disk_size", "10"),
				resource.TestCheckResourceAttr("pgvecto-rs-cloud_cluster.enterprise_plan_cluster", "status", "Ready"),
				resource.TestCheckResourceAttr("pgvecto-rs-cloud_cluster.enterprise_plan_cluster", "enable_pooler", "true"),
			),
		},
	}

	if testBackup {
		steps = append(steps, resource.TestStep{
			Config: testAccCheckAPIKeyConfigBasic() + testAccCheckResourceWithRestore(fmt.Sprintf("%s-restore", rName), backupID),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrSet("pgvecto-rs-cloud_cluster.enterprise_plan_cluster_restore_backup", "id"),
				resource.TestCheckResourceAttr("pgvecto-rs-cloud_cluster.enterprise_plan_cluster_restore_backup", "cluster_name", rName),
				resource.TestCheckResourceAttr("pgvecto-rs-cloud_cluster.enterprise_plan_cluster_restore_backup", "account_id", "5c3cb62b-d00b-4dda-85e6-2c0452d50138"),
				resource.TestCheckResourceAttr("pgvecto-rs-cloud_cluster.enterprise_plan_cluster_restore_backup", "plan", "Enterprise"),
				resource.TestCheckResourceAttr("pgvecto-rs-cloud_cluster.enterprise_plan_cluster_restore_backup", "server_resource", "aws-m7i-large-2c-8g"),
				resource.TestCheckResourceAttr("pgvecto-rs-cloud_cluster.enterprise_plan_cluster_restore_backup", "region", "us-east-1"),
				resource.TestCheckResourceAttr("pgvecto-rs-cloud_cluster.enterprise_plan_cluster_restore_backup", "cluster_provider", "aws"),
				resource.TestCheckResourceAttr("pgvecto-rs-cloud_cluster.enterprise_plan_cluster_restore_backup", "database_name", "test"),
				resource.TestCheckResourceAttr("pgvecto-rs-cloud_cluster.enterprise_plan_cluster_restore_backup", "pg_data_disk_size", "5"),
				resource.TestCheckResourceAttr("pgvecto-rs-cloud_cluster.enterprise_plan_cluster_restore_backup", "status", "Ready"),
				resource.TestCheckResourceAttr("pgvecto-rs-cloud_cluster.enterprise_plan_cluster_restore_backup", "enable_pooler", "true"),
				resource.TestCheckResourceAttr("pgvecto-rs-cloud_cluster.enterprise_plan_cluster_restore_backup", "enable_restore", "true"),
				resource.TestCheckResourceAttr("pgvecto-rs-cloud_cluster.enterprise_plan_cluster_restore_backup", "backup_id", backupID),
			),
		})
	}

	if testPITR {
		steps = append(steps, resource.TestStep{
			Config: testAccCheckAPIKeyConfigBasic() + testAccCheckResourcePITR(fmt.Sprintf("%s-pitr", rName), clusterID, targetTime),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrSet("pgvecto-rs-cloud_cluster.enterprise_plan_cluster_pitr", "id"),
				resource.TestCheckResourceAttr("pgvecto-rs-cloud_cluster.enterprise_plan_cluster_pitr", "cluster_name", rName),
				resource.TestCheckResourceAttr("pgvecto-rs-cloud_cluster.enterprise_plan_cluster_pitr", "account_id", "5c3cb62b-d00b-4dda-85e6-2c0452d50138"),
				resource.TestCheckResourceAttr("pgvecto-rs-cloud_cluster.enterprise_plan_cluster_pitr", "plan", "Enterprise"),
				resource.TestCheckResourceAttr("pgvecto-rs-cloud_cluster.enterprise_plan_cluster_pitr", "server_resource", "aws-m7i-large-2c-8g"),
				resource.TestCheckResourceAttr("pgvecto-rs-cloud_cluster.enterprise_plan_cluster_pitr", "region", "us-east-1"),
				resource.TestCheckResourceAttr("pgvecto-rs-cloud_cluster.enterprise_plan_cluster_pitr", "cluster_provider", "aws"),
				resource.TestCheckResourceAttr("pgvecto-rs-cloud_cluster.enterprise_plan_cluster_pitr", "database_name", "test"),
				resource.TestCheckResourceAttr("pgvecto-rs-cloud_cluster.enterprise_plan_cluster_pitr", "pg_data_disk_size", "5"),
				resource.TestCheckResourceAttr("pgvecto-rs-cloud_cluster.enterprise_plan_cluster_pitr", "status", "Ready"),
				resource.TestCheckResourceAttr("pgvecto-rs-cloud_cluster.enterprise_plan_cluster_pitr", "enable_pooler", "true"),
				resource.TestCheckResourceAttr("pgvecto-rs-cloud_cluster.enterprise_plan_cluster_pitr", "enable_restore", "true"),
				resource.TestCheckResourceAttr("pgvecto-rs-cloud_cluster.enterprise_plan_cluster_pitr", "target_cluster_id", clusterID),
				resource.TestCheckResourceAttr("pgvecto-rs-cloud_cluster.enterprise_plan_cluster_pitr", "target_time", targetTime),
			),
		})
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps:                    steps})
}

func testAccClusterResourceConfig(name string) string {
	return fmt.Sprintf(`
resource "pgvecto-rs-cloud_cluster" "enterprise_plan_cluster" {
	cluster_name      = %q
	account_id = "5c3cb62b-d00b-4dda-85e6-2c0452d50138"
	plan              = "Enterprise"
	image			 = "16-v0.4.0-extensions-exts"
	server_resource   = "aws-m7i-large-2c-8g"
	region            = "us-east-1"
	cluster_provider  = "aws"
	database_name    = "test"
	pg_data_disk_size = "5"
	enable_pooler     = true
}
`, name)
}

func testAccClusterResourceUpdateConfig(name string) string {
	return fmt.Sprintf(`
resource "pgvecto-rs-cloud_cluster" "enterprise_plan_cluster" {
	cluster_name      = %q
	account_id = "5c3cb62b-d00b-4dda-85e6-2c0452d50138"
	plan              = "Enterprise"
	image			 = "16-v0.4.0-extensions-exts"
	server_resource   = "aws-m7i-large-2c-8g"
	region            = "us-east-1"
	cluster_provider  = "aws"
	database_name    = "test"
	pg_data_disk_size = "10"
	enable_pooler     = true
}
`, name)
}

func testAccCheckResourceWithRestore(name string, backupID string) string {
	return fmt.Sprintf(`
resource "pgvecto-rs-cloud_cluster" "enterprise_plan_cluster_restore_backup" {
	cluster_name      = %q
	account_id = "5c3cb62b-d00b-4dda-85e6-2c0452d50138"
	plan              = "Enterprise"
	image			 = "16-v0.4.0-extensions-exts"
	server_resource   = "aws-m7i-large-2c-8g"
	region            = "us-east-1"
	cluster_provider  = "aws"
	database_name    = "test"
	pg_data_disk_size = "5"
	enable_pooler     = true
	enable_restore    = true
	backup_id         = %q
}
`, name, backupID)
}

func testAccCheckResourcePITR(name, clusterID, targetTime string) string {
	return fmt.Sprintf(`
resource "pgvecto-rs-cloud_cluster" "enterprise_plan_cluster_pitr" {
	cluster_name      = %q
	account_id = "5c3cb62b-d00b-4dda-85e6-2c0452d50138"
	plan              = "Enterprise"
	image			 = "16-v0.4.0-extensions-exts"
	server_resource   = "aws-m7i-large-2c-8g"
	region            = "us-east-1"
	cluster_provider  = "aws"
	database_name    = "test"
	pg_data_disk_size = "5"
	enable_pooler     = true
	enable_restore    = true
	target_cluster_id = %q
	target_time       = %q
}
`, name, clusterID, targetTime)
}
