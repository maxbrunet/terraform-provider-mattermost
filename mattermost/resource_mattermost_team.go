package mattermost

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/mattermost/mattermost-server/v5/model"
)

func resourceMattermostTeam() *schema.Resource {
	return &schema.Resource{
		Create: resourceMattermostTeamCreate,
		Read:   resourceMattermostTeamRead,
		Update: resourceMattermostTeamUpdate,
		Delete: resourceMattermostTeamDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"display_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"email": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"O", "I"}, false),
			},
			"company_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"allowed_domains": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"invite_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"allow_open_invite": &schema.Schema{
				Type:     schema.TypeBool,
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
			"deletion": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"soft", "permanent"}, false),
				Default:      "soft",
			},
		},
	}
}

func buildMattermostTeamStruct(d *schema.ResourceData) *model.Team {
	schemeID := d.Get("scheme_id").(string)
	groupConstrained := d.Get("group_constrained").(bool)

	team := &model.Team{
		DisplayName:      d.Get("display_name").(string),
		Name:             d.Get("name").(string),
		Description:      d.Get("description").(string),
		Email:            d.Get("email").(string),
		Type:             d.Get("type").(string),
		CompanyName:      d.Get("company_name").(string),
		AllowedDomains:   d.Get("allowed_domains").(string),
		AllowOpenInvite:  d.Get("allow_open_invite").(bool),
		SchemeId:         &schemeID,
		GroupConstrained: &groupConstrained,
	}

	return team
}

func resourceMattermostTeamCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*model.Client4)
	team := buildMattermostTeamStruct(d)

	log.Printf("[INFO] Creating Mattermost team: %s", team.Name)

	team, response := client.CreateTeam(team)
	if response.Error != nil {
		return response.Error
	}

	d.SetId(team.Id)

	return resourceMattermostTeamRead(d, meta)
}

func resourceMattermostTeamRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*model.Client4)

	log.Printf("[INFO] Reading Mattermost team: %s", d.Id())

	team, response := client.GetTeam(d.Id(), "")
	if response.Error != nil {
		return response.Error
	}

	d.Set("display_name", team.DisplayName)
	d.Set("name", team.Name)
	d.Set("description", team.Description)
	d.Set("email", team.Email)
	d.Set("type", team.Type)
	d.Set("company_name", team.CompanyName)
	d.Set("allowed_domains", team.AllowedDomains)
	d.Set("invite_id", team.InviteId)
	d.Set("allow_open_invite", team.AllowOpenInvite)
	d.Set("scheme_id", team.SchemeId)
	d.Set("group_constrained", team.GroupConstrained)
    d.Set("deletion", d.Get("deletion").(string))

	return nil
}

func resourceMattermostTeamUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*model.Client4)
	team := buildMattermostTeamStruct(d)
	team.Id = d.Id()

	log.Printf("[INFO] Updating Mattermost team: %s", team.Id)

	if _, response := client.UpdateTeam(team); response.Error != nil {
		return response.Error
	}

	return resourceMattermostTeamRead(d, meta)
}

func resourceMattermostTeamDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*model.Client4)

	var response *model.Response
	if d.Get("deletion").(string) == "permanent" {
		log.Printf("[INFO] Permanently deleting Mattermost team: %s", d.Id())
		_, response = client.PermanentDeleteTeam(d.Id())
	} else {
		log.Printf("[INFO] Archiving Mattermost team: %s", d.Id())
		_, response = client.SoftDeleteTeam(d.Id())
	}
	if response.Error != nil {
		return response.Error
	}

	d.SetId("")

	return nil
}
