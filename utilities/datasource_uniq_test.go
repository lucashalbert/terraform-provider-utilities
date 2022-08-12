package utilities

import (
	"fmt"
	"regexp"
	"testing"

	r "github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// Mock utilities_uniq data source
const uniqTest = `
data "utilities_uniq" "cidrs" {
    list = [
		"10.10.0.0/16",
		"10.11.0.0/24",
		"10.12.0.0/24",
		"10.11.1.0/24",
		"10.10.0.0/16"
	]
}
`

// Expected duplicates data structure from mock utilities_uniq data source
func getUniqExpectedDuplicates() []string {
	return []string{
		"10.10.0.0/16",
	}
}

// Expected uniques data structure from mock utilities_uniq data source
func getUniqExpectedUniques() []string {
	return []string{
		"10.10.0.0/16",
		"10.11.0.0/24",
		"10.12.0.0/24",
		"10.11.1.0/24",
	}
}

func getUniqExpectedlist() []string {
	return []string{
		"10.10.0.0/16",
		"10.11.0.0/24",
		"10.12.0.0/24",
		"10.11.1.0/24",
		"10.10.0.0/16",
	}
}

func TestUniq(t *testing.T) {
	r.UnitTest(t, r.TestCase{
		ProviderFactories: testProviderFactories,
		Steps: []r.TestStep{
			{
				Config: uniqTest,
				Check: r.ComposeTestCheckFunc(
					// --- Check duplicates Data Structure --- //
					// Check length of expected duplicates data structure
					r.TestCheckResourceAttr("data.utilities_uniq.cidrs", "duplicates.#", fmt.Sprint(len(getUniqExpectedDuplicates()))),

					// Validate all elements of the expected uniques data structure
					r.TestCheckResourceAttr("data.utilities_uniq.cidrs", "duplicates.0", "10.10.0.0/16"),

					// --- Check uniques Data Structure --- //
					// Check length of expected uniques data structure
					r.TestCheckResourceAttr("data.utilities_uniq.cidrs", "uniques.#", fmt.Sprint(len(getUniqExpectedUniques()))),

					// Validate all elements of the expected uniques data structure
					r.TestCheckResourceAttr("data.utilities_uniq.cidrs", "uniques.0", "10.10.0.0/16"),
					r.TestCheckResourceAttr("data.utilities_uniq.cidrs", "uniques.1", "10.11.0.0/24"),
					r.TestCheckResourceAttr("data.utilities_uniq.cidrs", "uniques.2", "10.12.0.0/24"),
					r.TestCheckResourceAttr("data.utilities_uniq.cidrs", "uniques.3", "10.11.1.0/24"),

					// --- Check list Data Structure --- //
					// Check length of expected list data structure
					r.TestCheckResourceAttr("data.utilities_uniq.cidrs", "list.#", fmt.Sprint(len(getUniqExpectedlist()))),

					// Validate all elements of the expected list data structure
					r.TestCheckResourceAttr("data.utilities_uniq.cidrs", "list.0", "10.10.0.0/16"),
					r.TestCheckResourceAttr("data.utilities_uniq.cidrs", "list.1", "10.11.0.0/24"),
					r.TestCheckResourceAttr("data.utilities_uniq.cidrs", "list.2", "10.12.0.0/24"),
					r.TestCheckResourceAttr("data.utilities_uniq.cidrs", "list.3", "10.11.1.0/24"),
					r.TestCheckResourceAttr("data.utilities_uniq.cidrs", "list.4", "10.10.0.0/16"),

					// --- Check ID Data Structure --- //
					// Check that a valid UUID is set for id
					r.TestMatchResourceAttr("data.utilities_uniq.cidrs", "id", regexp.MustCompile(`^\w{8}-\w{4}-\w{4}-\w{4}-\w{12}$`)),
				),
			},
		},
	})
}
