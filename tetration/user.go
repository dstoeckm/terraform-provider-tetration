package tetration

import (
	"github.com/hashicorp/terraform/helper/schema"
	client "github.com/tetration-exchange/terraform-go-sdk"
	tetration "github.com/tetration-exchange/terraform-go-sdk"
)

func resourceTetrationUser() *schema.Resource {
	return &schema.Resource{
		Create: resourceTetrationUserCreate,
		Update: nil,
		Read:   resourceTetrationUserRead,
		Delete: resourceTetrationUserDelete,

		SchemaVersion: 1,

		Schema: map[string]*schema.Schema{
			"email": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Email address associated with the user account.",
			},
			"first_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Userʼs first name.",
			},
			"last_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Userʼs last name.",
			},
			"app_scope_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "(Optional) Root scope to which the user belongs.",
			},
			"role_ids": {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "(Optional) A list of roles to be assigned to the user.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"enable_existing": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				ForceNew:    true,
				Description: "If true, and an existing but disabled user with the same email exists they will be enabled.",
			},
			"disabled_at": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "UNIX timestamp indicating when the user account was disabled. Zero or null if not disabled.",
			},
		},
	}
}

func resourceTetrationUserCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(client.Client)
	enableExistingUser := d.Get("enable_existing").(bool)
	if enableExistingUser {
		users, err := client.ListUsers(tetration.ListUsersRequest{
			AppScopeId:      d.Get("app_scope_id").(string),
			IncludeDisabled: true,
		})
		if err != nil {
			return err
		}
		var userExists bool
		var user tetration.User
		for _, existingUser := range users {
			if existingUser.Email == d.Get("email").(string) {
				user = existingUser
				userExists = true
				break
			}
		}
		if userExists {
			user, err = client.EnableUser(user.Id)
			if err != nil {
				return err
			}
			d.SetId(user.Id)
			return nil
		}
	}
	tfRoleIds := d.Get("role_ids").(*schema.Set).List()
	roleIds := make([]string, len(tfRoleIds))
	for _, tfRoleId := range tfRoleIds {
		roleIds = append(roleIds, tfRoleId.(string))
	}
	createUserParams := tetration.CreateUserRequest{
		Email:      d.Get("email").(string),
		FirstName:  d.Get("first_name").(string),
		LastName:   d.Get("last_name").(string),
		AppScopeId: d.Get("app_scope_id").(string),
		RoleIds:    roleIds,
	}
	user, err := client.CreateUser(createUserParams)
	if err != nil {
		return err
	}
	d.SetId(user.Id)
	return nil
}

func resourceTetrationUserRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(client.Client)
	user, err := client.DescribeUser(d.Id())
	if err != nil {
		return err
	}
	d.Set("email", user.Email)
	d.Set("first_name", user.FirstName)
	d.Set("last_name", user.LastName)
	d.Set("app_scope_id", user.AppScopeId)
	d.Set("role_ids", user.RoleIds)
	d.Set("disabled_at", user.DisabledAt)
	return nil
}

func resourceTetrationUserDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(client.Client)
	return client.DeleteUser(d.Id())
}
