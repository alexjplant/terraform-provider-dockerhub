package dockerhub

import (
	"context"
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


func dataSourceImageTagRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

  var client Client 

  j, err := client.GetImageTags(d.Get("namespace").(string), d.Get("repository_name").(string), d.Get("tag_name").(string))

  if err != nil {
   return diag.FromErr(err) 
  }

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
