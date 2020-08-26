package mattermost

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/mattermost/mattermost-server/v5/model"
)

func resourceMattermostChannel() *schema.Resource {
	return &schema.Resource{
		Create: resourceMattermostChannelCreate,
		Read:   resourceMattermostChannelRead,
		Update: resourceMattermostChannelUpdate,
		Delete: resourceMattermostChannelDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"team_id": &schema.Schema{
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"type": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"O", "P"}, false),
			},
			"display_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"header": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"purpose": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"scheme_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"group_constrained": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}
}

func buildMattermostChannelStruct(d *schema.ResourceData) *model.Channel {
	schemeID := d.Get("scheme_id").(string)
	groupConstrained := d.Get("group_constrained").(bool)

	channel := &model.Channel{
		TeamId:           d.Get("team_id").(string),
		Type:             d.Get("type").(string),
		DisplayName:      d.Get("display_name").(string),
		Name:             d.Get("name").(string),
		Header:           d.Get("header").(string),
		Purpose:          d.Get("purpose").(string),
		SchemeId:         &schemeID,
		GroupConstrained: &groupConstrained,
	}

	return channel
}

func resourceMattermostChannelCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*model.Client4)
	channel := buildMattermostChannelStruct(d)

	log.Printf("[INFO] Creating Mattermost channel: %s", channel.Name)

	channel, response := client.CreateChannel(channel)
	if response.Error != nil {
		return response.Error
	}

	d.SetId(channel.Id)

	return resourceMattermostChannelRead(d, meta)
}

func resourceMattermostChannelRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*model.Client4)

	log.Printf("[INFO] Reading Mattermost channel: %s", d.Id())

	channel, response := client.GetChannel(d.Id(), "")
	if response.Error != nil {
		return response.Error
	}

	d.Set("team_id", channel.TeamId)
	d.Set("type", channel.Type)
	d.Set("display_name", channel.DisplayName)
	d.Set("name", channel.Name)
	d.Set("header", channel.Header)
	d.Set("purpose", channel.Purpose)
	d.Set("scheme_id", channel.SchemeId)
	d.Set("group_constrained", channel.GroupConstrained)

	return nil
}

func resourceMattermostChannelUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*model.Client4)
	channel := buildMattermostChannelStruct(d)
	channel.Id = d.Id()

	log.Printf("[INFO] Updating Mattermost channel: %s", channel.Id)

	if _, response := client.UpdateChannel(channel); response.Error != nil {
		return response.Error
	}

	return resourceMattermostChannelRead(d, meta)
}

func resourceMattermostChannelDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*model.Client4)

	log.Printf("[INFO] Deleting Mattermost channel: %s", d.Id())

	if _, response := client.DeleteChannel(d.Id()); response.Error != nil {
		return response.Error
	}

	d.SetId("")
	return nil
}
