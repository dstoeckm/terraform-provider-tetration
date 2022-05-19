package tetration

import (
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"

	client "github.com/tetration-exchange/terraform-go-sdk"
)

// Provider returns a terraform resource provider for managing tetration resources.
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"api_key": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("TETRATION_API_KEY", nil),
				Description: "API key for calculating request signatures for Tetration API calls.",
			},
			"api_secret": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("TETRATION_API_SECRET", nil),
				Description: "API secret for calculating request signatures for Tetration API calls.",
			},
			"api_url": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("TETRATION_API_URL", nil),
				Description: "URL for a Tetration API.",
			},
			"disable_tls_verification": {
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("TETRATION_DISABLE_TLS_VERIFICATION", false),
				Description: "Allow connections to Tetration endpoints without validating their TLS certificate.",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"tetration_filter":      resourceTetrationFilter(),
			"tetration_scope":       resourceTetrationScope(),
			"tetration_tag":         resourceTetrationTag(),
			"tetration_user":        resourceTetrationUser(),
			"tetration_application": resourceTetrationApplication(),
			"tetration_role":        resourceTetrationRole(),
		},
		ConfigureFunc: configureClient,
	}
}

func configureClient(d *schema.ResourceData) (interface{}, error) {
	config := client.Config{
		APIKey:                 d.Get("api_key").(string),
		APISecret:              d.Get("api_secret").(string),
		APIURL:                 d.Get("api_url").(string),
		DisableTLSVerification: d.Get("disable_tls_verification").(bool),
	}
	if err := validate(config); err != nil {
		return nil, err
	}
	client, err := client.New(config)
	if err != nil {
		return nil, err
	}
	return client, nil
}

// validate validates the config needed to initialize a tetration client,
// returning a single error with all validation errors, or nil if no error.
func validate(config client.Config) error {
	var err *multierror.Error
	if config.APIKey == "" {
		err = multierror.Append(err, fmt.Errorf("API Key must be configured for the Tetration provider"))
	}
	if config.APISecret == "" {
		err = multierror.Append(err, fmt.Errorf("API Secret must be configured for the Tetration provider"))
	}
	if config.APIURL == "" {
		err = multierror.Append(err, fmt.Errorf("API URL must be configured for the Tetration provider"))
	}
	return err.ErrorOrNil()
}
