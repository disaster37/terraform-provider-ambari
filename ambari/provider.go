package ambari

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/jinzhu/configor"
	"github.com/pkg/errors"
)

type CLIConfig struct {
	URL      string
	Login    string
	Password string
}

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"api_url": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("AMBARI_URL", ""),
				Description: "The URL to the Ambari API, must include version uri (v1)",
			},
			"login": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("AMBARI_LOGIN", ""),
				Description: "The login used to authenticate with the Ambari server",
			},
			"password": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("AMBARI_PASSWORD", ""),
				Description: "The password used to authenticate with the Ambari server",
			},
			"config": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("RANCHER_CLIENT_CONFIG", ""),
				Description: "Path to the Ambari client cli.yaml config file",
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"ambari_privilege": resourceAmbariPrivilege(),
			"ambari_cluster":   resourceAmbariCluster(),
		},

		DataSourcesMap: map[string]*schema.Resource{},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {

	client := NewClient()
	client.APIURL = d.Get("api_url").(string)
	client.Login = d.Get("login").(string)
	client.Password = d.Get("password").(string)

	if configFile := d.Get("config").(string); configFile != "" {
		config := new(CLIConfig)
		err := configor.Load(config, configFile)
		if err != nil {
			return config, err
		}

		if client.APIURL == "" {
			client.APIURL = config.URL
		}

		if client.Login == "" {
			client.Login = config.Login
		}

		if client.Password == "" {
			client.Password = config.Password
		}
	}

	if client.APIURL == "" {
		return nil, errors.New("No api_url provided")
	}
	if client.Login == "" {
		return nil, errors.New("No login provided")
	}
	if client.Password == "" {
		return nil, errors.New("No password provided")
	}

	// Init client
	client.Client()

	return client, nil
}
