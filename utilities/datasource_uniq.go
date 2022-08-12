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
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID",
			},
			"duplicates": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed:    true,
				Description: "List of duplicates found in original list",
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

func unique(src interface{}) (interface{}, interface{}) {
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
	return uniques.Interface(), duplicates.Interface()
}

func dataSourceUniqRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	uniques, duplicates := unique(d.Get("list"))

	// Set the schema uniques value
	if err := d.Set("uniques", uniques); err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to set the 'uniq' utility 'uniques' value",
			Detail:   "Unable to set the value of the 'uniq' utility 'uniques' output.",
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

	return diags
}
