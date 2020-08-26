package mattermost

import (
	"fmt"
	"log"
	"runtime"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/mattermost/mattermost-server/v5/model"
)

// Provider represents a resource provider in Terraform
func Provider() terraform.ResourceProvider {
	p := &schema.Provider{
		Schema: map[string]*schema.Schema{
			"token": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("MM_TOKEN", nil),
			},
			"url": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("MM_URL", nil),
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"mattermost_channel": resourceMattermostChannel(),
			"mattermost_team":    resourceMattermostTeam(),
		},
	}

	p.ConfigureFunc = func(d *schema.ResourceData) (interface{}, error) {
		terraformVersion := p.TerraformVersion
		if terraformVersion == "" {
			// Terraform 0.12 introduced this field to the protocol
			// We can therefore assume that if it's missing it's 0.10 or 0.11
			terraformVersion = "0.11+compatible"
		}
		return providerConfigure(d, terraformVersion)
	}
	return p
}

func providerConfigure(data *schema.ResourceData, terraformVersion string) (interface{}, error) {
	token := data.Get("token").(string)
	url := data.Get("url").(string)

	log.Println("[INFO] Initializing Mattermost client")
	client := model.NewAPIv4Client(url)
	userAgent := fmt.Sprintf("(%s %s) Terraform/%s", runtime.GOOS, runtime.GOARCH, terraformVersion)
	client.HttpHeader = map[string]string{"User-Agent": userAgent}

	client.AuthType = model.HEADER_BEARER
	client.AuthToken = token

	_, response := client.GetMe("")
	if response.Error != nil {
		return nil, response.Error
	}

	return client, nil
}
