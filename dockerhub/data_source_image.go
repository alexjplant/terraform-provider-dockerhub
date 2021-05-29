package dockerhub

import (
	"context"
	"net/http"
	"time"
	//  "log"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceImageTags() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceImageTagRead,
		Schema: map[string]*schema.Schema{
			"repository_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"namespace": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"tag_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"creator": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"image_id": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
				Optional: true,
			},
			"last_updated": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"last_updater": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
				Optional: true,
			},
			"last_updater_username": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"repository": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"full_size": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"v2": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
			"tag_status": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"tag_last_pulled": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"tag_last_pushed": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
		},
	}
}

type ApiImageTag struct {
	Creator             int    `json:"creator"`
	ID                  int    `json:"id"`
	ImageID             string `json:"image_id"`
	LastUpdated         string `json:"last_updated"`
	LastUpdater         string `json:"last_updater"`
	LastUpdaterUsername string `json:"last_updater_username"`
	Name                string `json:"name"`
	Repository          int    `json:"repository"`
	FullSize            int    `json:"full_size"`
	V2                  bool   `json:"v2"`
	TagStatus           string `json:"tag_status"`
	TagLastPulled       string `json:"tag_last_pulled"`
	TagLastPushed       string `json:"tag_last_pushed"`
}

func dataSourceImageTagRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	client := &http.Client{Timeout: 10 * time.Second}

	var diags diag.Diagnostics

	req, err := http.NewRequest("GET", fmt.Sprintf("https://hub.docker.com/v2/repositories/%s/%s/tags/%s", d.Get("namespace"), d.Get("repository_name"), d.Get("tag_name")), nil)
	if err != nil {
		return diag.FromErr(err)
	}

	r, err := client.Do(req)
	if err != nil {
		return diag.FromErr(err)
	}
	if r.StatusCode != 200 {
		var d diag.Diagnostic
		d.Severity = diag.Error
		d.Summary = "Image tag not found"
		d.Detail = fmt.Sprintf("Received HTTP %v when invoking HTTP API", r.StatusCode)
		diags = append(diags, d)
	}

	defer r.Body.Close()

	var j ApiImageTag
	err = json.NewDecoder(r.Body).Decode(&j)

	d.SetId(strconv.Itoa(j.ID))
	d.Set("creator", j.Creator)
	d.Set("image_id", j.ImageID)
	d.Set("last_updated", j.LastUpdated)
	d.Set("last_updater", j.LastUpdater)
	d.Set("last_updater_username", j.LastUpdaterUsername)
	d.Set("repository", j.Repository)
	d.Set("full_size", j.FullSize)
	d.Set("v2", j.V2)
	d.Set("tag_status", j.TagStatus)
	d.Set("tag_last_pulled", j.TagLastPulled)
	d.Set("tag_last_pushed", j.TagLastPushed)

	return diags
}
