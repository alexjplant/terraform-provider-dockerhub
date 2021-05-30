package dockerhub

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"net/http"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"username": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("DOCKERHUB_USERNAME", nil),
			},
			"password": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("DOCKERHUB_PASSWORD", nil),
			},
		},
		ResourcesMap: map[string]*schema.Resource{},
		DataSourcesMap: map[string]*schema.Resource{
			"dockerhub_image_tag": dataSourceImageTags(),
		},
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	var diags diag.Diagnostics

	username, usernameSet := d.GetOk("username")
	password, passwordSet := d.GetOk("password")

	if usernameSet && passwordSet {
		login, _ := json.Marshal(map[string]string{"username": username.(string), "password": password.(string)})

		res, _ := http.Post("https://hub.docker.com/v2/users/login/", "application/json", bytes.NewBuffer(login))
		defer res.Body.Close()
		var resJSON map[string]string

		// json.Unmarshal(res.Body, &resJSON)
		json.NewDecoder(res.Body).Decode(&resJSON)
		return resJSON["token"], diags
	} else {
		return nil, diags
	}
}
