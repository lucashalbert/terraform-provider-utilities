package utilities

import (
	"context"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceFail() *schema.Resource {
	return &schema.Resource{
		Description: "An autofail resource that always fails with an error",
		ReadContext: dataSourceFailRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    false,
				Optional:    true,
				Description: "Name of the fail resource.",
			},
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID",
			},
		},
	}
}

func dataSourceFailRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	diags = append(diags, diag.Diagnostic{
		Severity: diag.Error,
		Summary:  "Predetermined failure data source",
		Detail:   "A predetermined failure is expected from this data source",
	})

	// Always define the schema id value
	d.SetId(uuid.New().String())

	return diags
}
