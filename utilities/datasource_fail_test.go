package utilities

import (
	"regexp"
	"testing"

	r "github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// Mock utilities_uniq data source
const failTest = `
data "utilities_fail" "fail" {}
`

func TestFail(t *testing.T) {
	r.UnitTest(t, r.TestCase{
		ProviderFactories: testProviderFactories,
		Steps: []r.TestStep{
			{
				Config:      failTest,
				ExpectError: regexp.MustCompile(`Error: Predetermined failure data source`),
			},
		},
	})
}
