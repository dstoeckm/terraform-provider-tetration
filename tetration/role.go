package tetration

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"

	client "github.com/tetration-exchange/terraform-go-sdk"
	tetration "github.com/tetration-exchange/terraform-go-sdk"
)

var (
	ValidAbilities        = []string{"SCOPE_READ", "SCOPE_WRITE", "EXECUTE", "ENFORCE", "SCOPE_OWNER", "DEVELOPER"}
	AccessTypeDescription = fmt.Sprintf("The type of access to grant the role to the `access_app_scope_id` scope.\n Valid values are [%s]", strings.Join(ValidAbilities, ", "))
)

func resourceTetrationRole() *schema.Resource {
	return &schema.Resource{
		Create: resourceTetrationRoleCreate,
		Update: nil,
		Read:   resourceTetrationRoleRead,
		Delete: resourceTetrationRoleDelete,

		SchemaVersion: 1,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    false,
				Computed:    true,
				Description: "(Optional) User-specified name for the role.",
			},
			"description": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The role's description",
			},
			"access_app_scope_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The scope to which this role will be given access",
			},
			"app_scope_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The scope in which this role will be created",
			},
			"access_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  AccessTypeDescription,
				ValidateFunc: validation.StringInSlice(ValidAbilities, true),
			},
			"user_ids": {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The users to which this role will be assigned",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceTetrationRoleCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(client.Client)
	tfUserIds := d.Get("user_ids").(*schema.Set).List()
	userIds := []string{}
	for _, tfUserId := range tfUserIds {
		if tfUserId != "" {
			userIds = append(userIds, tfUserId.(string))
		}
	}

	createScopedRoleForUsersParams := tetration.CreateScopedRoleForUsersRequest{
		CreateScopedRoleRequest: tetration.CreateScopedRoleRequest{
			Name:                d.Get("name").(string),
			Description:         d.Get("description").(string),
			AppScopeId:          d.Get("app_scope_id").(string),
			AbilitiesAppScopeId: d.Get("access_app_scope_id").(string),
			Ability:             d.Get("access_type").(string),
		},
		Users: userIds,
	}

	response, err := client.CreateScopedRoleForUsers(createScopedRoleForUsersParams)
	if err != nil {
		return err
	}
	d.SetId(response.RoleId)
	return nil
}

func resourceTetrationRoleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(client.Client)
	role, err := client.GetRole(d.Id())
	if err != nil {
		return err
	}
	d.Set("app_scope_id", role.AppScopeId)
	d.Set("name", role.Name)
	d.Set("description", role.Description)
	return nil
}
func resourceTetrationRoleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(client.Client)
	return client.DeleteRole(d.Id())
}
