package foreman

import (
	"fmt"

	"github.com/HanseMerkur/terraform-provider-foreman/foreman/api"
	"github.com/HanseMerkur/terraform-provider-utils/autodoc"
	"github.com/HanseMerkur/terraform-provider-utils/helper"
	"github.com/HanseMerkur/terraform-provider-utils/log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceForemanOverrideValue() *schema.Resource {
	// copy attributes from resource definition
	r := resourceForemanOverrideValue()
	ds := helper.DataSourceSchemaFromResourceSchema(r.Schema)

	// define searchable attributes for the data source
	ds["smart_class_parameter_id"] = &schema.Schema{
		Type:     schema.TypeString,
		Required: true,
		Description: fmt.Sprintf(
			"The name of the override_value - the full DNS override_value name. "+
				"%s \"dev.dc1.company.com\"",
			autodoc.MetaExample,
		),
	}

	return &schema.Resource{

		Read: dataSourceForemanOverrideValueRead,

		// NOTE(ALL): See comments in the corresponding resource file
		Schema: ds,
	}
}

func dataSourceForemanOverrideValueRead(d *schema.ResourceData, meta interface{}) error {
	log.Tracef("data_source_foreman_override_value.go#Read")

	client := meta.(*api.Client)
	override_value := buildForemanOverrideValue(d)

	log.Debugf("ForemanOverrideValue: [%+v]", override_value)

	queryResponse, queryErr := client.QueryOverrideValue(override_value)
	if queryErr != nil {
		return queryErr
	}

	if queryResponse.Subtotal == 0 {
		return fmt.Errorf("Data source override_value returned no results")
	} else if queryResponse.Subtotal > 1 {
		return fmt.Errorf("Data source override_value returned more than 1 result")
	}

	var queryOverrideValue api.ForemanOverrideValue
	var ok bool
	if queryOverrideValue, ok = queryResponse.Results[0].(api.ForemanOverrideValue); !ok {
		return fmt.Errorf(
			"Data source results contain unexpected type. Expected "+
				"[api.ForemanOverrideValue], got [%T]",
			queryResponse.Results[0],
		)
	}
	override_value = &queryOverrideValue

	log.Debugf("ForemanOverrideValue: [%+v]", override_value)

	setResourceDataFromForemanOverrideValue(d, override_value)

	return nil
}
