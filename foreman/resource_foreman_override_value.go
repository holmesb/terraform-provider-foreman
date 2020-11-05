package foreman

import (
	"fmt"
	"strconv"

	"github.com/HanseMerkur/terraform-provider-foreman/foreman/api"
	"github.com/HanseMerkur/terraform-provider-utils/autodoc"
	"github.com/HanseMerkur/terraform-provider-utils/log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceForemanOverrideValue() *schema.Resource {
	return &schema.Resource{

		Create: resourceForemanOverrideValueCreate,
		Read:   resourceForemanOverrideValueRead,
		Update: resourceForemanOverrideValueUpdate,
		Delete: resourceForemanOverrideValueDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{

			autodoc.MetaAttribute: &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
				Description: fmt.Sprintf(
					"%s Foreman representation of common_parameter. Global parameters are available for all resources",
					autodoc.MetaSummary,
				),
			},

			// -- Reference --
			"smart_class_parameter_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			// -- Actual Content --
			"match": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"value": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"omit": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}
}

// -----------------------------------------------------------------------------
// Conversion Helpers
// -----------------------------------------------------------------------------

// buildForemanOverrideValue constructs a ForemanOverrideValue reference from a resource data
// reference.  The struct's  members are populated from the data populated in
// the resource data.  Missing members will be left to the zero value for that
// member's type.
func buildForemanOverrideValue(d *schema.ResourceData) *api.ForemanOverrideValue {
	log.Tracef("resource_foreman_common_parameter.go#buildForemanOverrideValue")

	override_value := api.ForemanOverrideValue{}

	obj := buildForemanObject(d)
	override_value.ForemanObject = *obj

	var attr interface{}
	var ok bool

	if attr, ok = d.GetOk("match"); ok {
		override_value.Match = attr.(string)
	}
	if attr, ok = d.GetOk("value"); ok {
		override_value.Value = attr.(string)
	}
	if attr, ok = d.GetOk("omit"); ok {
		override_value.Omit = attr.(bool)
	}
	return &common_parameter
}

// setResourceDataFromForemanOverrideValue sets a ResourceData's attributes from the
// attributes of the supplied ForemanOverrideValue reference
func setResourceDataFromForemanOverrideValue(d *schema.ResourceData, fd *api.ForemanOverrideValue) {
	log.Tracef("resource_foreman_common_parameter.go#setResourceDataFromForemanOverrideValue")

	d.SetId(strconv.Itoa(fd.Id))
	d.Set("match", fd.Match)
	d.Set("value", fd.Value)
	d.Set("omit", fd.Omit)
}

// -----------------------------------------------------------------------------
// Resource CRUD Operations
// -----------------------------------------------------------------------------

func resourceForemanOverrideValueCreate(d *schema.ResourceData, meta interface{}) error {
	log.Tracef("resource_foreman_common_parameter.go#Create")

	client := meta.(*api.Client)
	p := buildForemanOverrideValue(d)

	log.Debugf("ForemanOverrideValue: [%+v]", d)

	createdParam, createErr := client.CreateOverrideValue(p)
	if createErr != nil {
		return createErr
	}

	log.Debugf("Created ForemanOverrideValue: [%+v]", createdParam)

	setResourceDataFromForemanOverrideValue(d, createdParam)

	return nil
}

func resourceForemanOverrideValueRead(d *schema.ResourceData, meta interface{}) error {
	log.Tracef("resource_foreman_common_parameter.go#Read")

	client := meta.(*api.Client)
	common_parameter := buildForemanOverrideValue(d)

	log.Debugf("ForemanOverrideValue: [%+v]", common_parameter)

	readOverrideValue, readErr := client.ReadOverrideValue(common_parameter, common_parameter.Id)
	if readErr != nil {
		return readErr
	}

	log.Debugf("Read ForemanOverrideValue: [%+v]", readOverrideValue)

	setResourceDataFromForemanOverrideValue(d, readOverrideValue)

	return nil
}

func resourceForemanOverrideValueUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Tracef("resource_foreman_common_parameter.go#Update")

	client := meta.(*api.Client)
	p := buildForemanOverrideValue(d)

	log.Debugf("ForemanOverrideValue: [%+v]", p)

	updatedParam, updateErr := client.UpdateOverrideValue(p, p.Id)
	if updateErr != nil {
		return updateErr
	}

	log.Debugf("Updated ForemanOverrideValue: [%+v]", updatedParam)

	setResourceDataFromForemanOverrideValue(d, updatedParam)

	return nil
}

func resourceForemanOverrideValueDelete(d *schema.ResourceData, meta interface{}) error {
	log.Tracef("resource_foreman_common_parameter.go#Delete")

	client := meta.(*api.Client)
	p := buildForemanOverrideValue(d)

	log.Debugf("ForemanOverrideValue: [%+v]", p)

	return client.DeleteOverrideValue(p, p.Id)
}
