package acceptance

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/hashicorp/go-azure-helpers/authentication"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var AzureADProvider *schema.Provider
var ProviderFactories map[string]func() (*schema.Provider, error)
var SupportedProviders map[string]*schema.Provider // TODO deprecated

func PreCheck(t *testing.T) {
	variables := []string{
		"ARM_CLIENT_ID",
		"ARM_CLIENT_SECRET",
		"ARM_TENANT_ID",
		"ARM_TEST_LOCATION",
		"ARM_TEST_LOCATION_ALT",
	}

	for _, variable := range variables {
		value := os.Getenv(variable)
		if value == "" {
			t.Fatalf("`%s` must be set for acceptance tests!", variable)
		}
	}
}

func EnvironmentName() string {
	envName, exists := os.LookupEnv("ARM_ENVIRONMENT")
	if !exists {
		envName = "public"
	}
	return envName
}

func Environment() (*azure.Environment, error) {
	envName := EnvironmentName()
	metadataHost := os.Getenv("ARM_METADATA_HOST")
	return authentication.AzureEnvironmentByNameFromEndpoint(context.TODO(), metadataHost, envName)
}

func RequiresImportError(resourceName string) *regexp.Regexp {
	message := "To be managed via Terraform, this resource needs to be imported into the State. Please see the resource documentation for %q for more information."
	message = strings.Replace(message, " ", "\\s+", -1)
	return regexp.MustCompile(fmt.Sprintf(message, resourceName))
}
