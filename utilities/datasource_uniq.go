package utilities

import (
	"context"
	"reflect"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceUniq() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceUniqRead,

		Schema: map[string]*schema.Schema{
			"list": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Required:    true,
				Description: "Provided list of items to run against uniq",
			},
			"fail_on_duplicate": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Force data source failure upon presence of duplicates",
			},
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID",
			},
			"total_duplicates": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total number of duplicate items found in original list",
			},
			"duplicates": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed:    true,
				Description: "List of duplicates found in original list",
			},
			"total_uniques": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total number of unique items found in original list",
			},
			"uniques": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed:    true,
				Description: "List of uniques found in original list",
			},
		},
	}
}

func unique(src interface{}) (int, interface{}, int, interface{}) {
	srcValue := reflect.ValueOf(src)
	uniques := reflect.MakeSlice(srcValue.Type(), 0, 0)
	duplicates := reflect.MakeSlice(srcValue.Type(), 0, 0)
	visited := make(map[interface{}]struct{})

	for i := 0; i < srcValue.Len(); i++ {
		elemValue := srcValue.Index(i)

		if _, ok := visited[elemValue.Interface()]; ok {
			// Append element value to duplicates slice
			duplicates = reflect.Append(duplicates, elemValue)
			continue
		}

		visited[elemValue.Interface()] = struct{}{}

		// Append element value to uniques slice
		uniques = reflect.Append(uniques, elemValue)
	}

	return uniques.Len(), uniques.Interface(), duplicates.Len(), duplicates.Interface()
}

func dataSourceUniqRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	// Get failOnDuplicate
	failOnDuplicate := d.Get("fail_on_duplicate").(bool)

	uniquesLen, uniques, duplicatesLen, duplicates := unique(d.Get("list"))

	// Set the schema uniques value
	if err := d.Set("total_uniques", uniquesLen); err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to set the 'uniq' utility 'total_uniques' value",
			Detail:   "Unable to set the value of the 'uniq' utility 'total_uniques' output.",
		})
	}

	// Set the schema uniques value
	if err := d.Set("uniques", uniques); err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to set the 'uniq' utility 'uniques' value",
			Detail:   "Unable to set the value of the 'uniq' utility 'uniques' output.",
		})
	}

	// Set the schema total_duplicates value
	if err := d.Set("total_duplicates", duplicatesLen); err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to set the 'uniq' utility 'total_duplicates' value",
			Detail:   "Unable to set the value of the 'uniq' utility 'total_duplicates' output.",
		})
	}

	// Set the schema duplicates value
	if err := d.Set("duplicates", duplicates); err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to set the 'uniq' utility 'duplicates' value",
			Detail:   "Unable to set the value of the 'uniq' utility 'duplicates' output.",
		})
	}

	// Always define the schema id value
	d.SetId(uuid.New().String())

	// Check if 'failOnDuplicate' boolean is true and if duplicates were returned
	if failOnDuplicate && (duplicatesLen > 0) {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "'total_duplicates' is greater than 0 while 'fail_on_duplicate' boolean is set to 'true'",
			Detail:   "The total number of duplicates returned by data source is greater than 0 while the 'fail_on_duplicate' control boolean is set to 'true'",
		})
	}

	return diags
}
