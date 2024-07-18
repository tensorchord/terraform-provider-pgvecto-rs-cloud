package provider

/*
func TestAccClusterDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccCheckAPIKeyConfigBasic() + testAccClusterDataConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.pgvecto-rs-cloud_cluster.starter_plan_cluster", "id"),
					resource.TestCheckResourceAttrSet("data.pgvecto-rs-cloud_cluster.starter_plan_cluster", "last_updated"),
				),
			},
		},
	})
}

func testAccClusterDataConfig() string {
	return `
data "pgvecto-rs-cloud_cluster" "starter_plan_cluster" {
	account_id = "5c3cb62b-d00b-4dda-85e6-2c0452d50138"
	id = "b955f953-b802-4e37-9303-608e382f3317"
}
`
}
*/
