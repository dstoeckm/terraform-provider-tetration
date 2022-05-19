package tetration

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	client "github.com/tetration-exchange/terraform-go-sdk"
	tetration "github.com/tetration-exchange/terraform-go-sdk"
)

const (
	TagIdDelimter = ":"
)

func resourceTetrationTag() *schema.Resource {
	return &schema.Resource{
		Create: resourceTetrationTagCreate,
		Update: resourceTetrationTagCreate,
		Read:   resourceTetrationTagRead,
		Delete: resourceTetrationTagDelete,

		SchemaVersion: 1,

		Schema: map[string]*schema.Schema{
			"tenant_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				ForceNew:    true,
				Description: "Tetration root app scope name.",
			},
			"ip": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "IPv4/IPv6 address or subnet.",
			},
			"attributes": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: false,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Key/value map for tagging matching flows and inventory items.",
			},
		},
	}
}

var requiredCreateTagParams = []string{"ip", "attributes"}

func resourceTetrationTagCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(client.Client)
	for _, param := range requiredCreateTagParams {
		if d.Get(param) == "" {
			return fmt.Errorf("%s is required but was not provided", param)
		}
	}
	tenantName := d.Get("tenant_name").(string)
	if tenantName == "" {
		tenantURL := client.Config.APIURL
		// strip protocol and extract the tenant name/subdomain from the url
		// e.g. https://acme.tetrationpreview.com => acme
		tenantName = strings.Split(strings.Split(tenantURL, "://")[1], ".")[0]
	}
	attributes := d.Get("attributes").(map[string]interface{})
	createTagParams := tetration.CreateTagRequest{
		RootScopeName: tenantName,
		Ip:            d.Get("ip").(string),
		Attributes:    attributes,
	}
	tag, err := client.CreateTag(createTagParams)
	if err != nil {
		return err
	}
	d.SetId(fmt.Sprintf("%s%s%s", createTagParams.RootScopeName, TagIdDelimter, tag.Ip))
	return nil
}

func resourceTetrationTagRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(client.Client)
	tagIdComponents := strings.Split(d.Id(), TagIdDelimter)
	describeTagRequest := tetration.DescribeTagRequest{
		RootAppScopeName: tagIdComponents[0],
		Ip:               tagIdComponents[1],
	}
	attributes := make(map[string]string)
	err := client.DescribeTag(describeTagRequest, &attributes)
	if err != nil {
		return err
	}
	d.Set("tenant_name", describeTagRequest.RootAppScopeName)
	d.Set("ip", describeTagRequest.Ip)
	d.Set("attributes", attributes)
	return nil
}

func resourceTetrationTagDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(client.Client)
	tagIdComponents := strings.Split(d.Id(), TagIdDelimter)
	deleteTagRequest := tetration.DeleteTagRequest{
		RootAppScopeName: tagIdComponents[0],
		Ip:               tagIdComponents[1],
	}
	return client.DeleteTag(deleteTagRequest)
}
