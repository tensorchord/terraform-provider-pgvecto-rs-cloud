package provider

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

var (
	// testAccProtoV6ProviderFactories are used to instantiate a provider during
	// acceptance testing. The factory function will be invoked for every Terraform
	// CLI command executed to create a provider server to which the CLI can
	// reattach.
	testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
		"pgvecto-rs-cloud": providerserver.NewProtocol6WithError(New("test")()),
	}
)

// Test the provider configuration with variables set.
func TestProvider(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAPIKeyConfigBasic(),
			},
		},
	})
}
func testAccCheckAPIKeyConfigBasic() string {
	apiKey := os.Getenv("PGVECTORS_CLOUD_API_KEY")
	apiURL := os.Getenv("PGVECTORS_CLOUD_API_URL")

	return fmt.Sprintf(`
provider "pgvecto-rs-cloud" {
  api_key = "%s"
  api_url = "%s"
}
`, apiKey, apiURL)
}
