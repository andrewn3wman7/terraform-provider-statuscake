package statuscake

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/andrewn3wman7/statuscake"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceStatusCakePageSpeed() *schema.Resource {
	return &schema.Resource{
		Create: CreatePagespeed,
		Update: UpdatePagespeed,
		Delete: DeletePagespeed,
		Read:   ReadPagespeed,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"check_rate": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  30,
			},

			"contact_group": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				Set:      schema.HashString,
			},

			"location_iso": {
				Type:     schema.TypeString,
				Required: true,
			},

			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"pagespeed_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"website_url": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func CreatePagespeed(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*statuscake.Client)

	newPagespeed := &statuscake.PartialPageSpeed{
		Checkrate:    strconv.Itoa(d.Get("check_rate").(int)),
		Location_iso: d.Get("location_iso").(string),
		Name:         d.Get("name").(string),
		Website_url:  d.Get("website_url").(string),
	}

	if v, ok := d.GetOk("contact_group"); ok {
		newPagespeed.ContactGroupsC = strings.Join(castSetToSliceStrings(v.(*schema.Set).List()), ",")
	}

	log.Printf("[DEBUG] Creating new StatusCake Pagespeed: %s", d.Get("name").(string))

	response, err := statuscake.NewPageSpeeds(client).Create(newPagespeed)
	if err != nil {
		return fmt.Errorf("Error creating StatusCake Pagespeed: %s", err.Error())
	}

	d.Set("check_rate", response.Checkrate)
	d.Set("location_iso", response.Location_iso)
	d.Set("name", response.Name)
	d.Set("pagespeed_id", response.ID)
	d.Set("website_url", response.Website_url)

	s := strconv.Itoa(response.ID)
	d.SetId(s)

	return ReadPagespeed(d, meta)
}

func UpdatePagespeed(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*statuscake.Client)

	params := getStatusCakePagespeedInput(d)

	log.Printf("[DEBUG] params", params)

	log.Printf("[DEBUG] StatusCake Pagespeed Update for %s", d.Id())
	_, err := statuscake.NewPageSpeeds(client).Update(params)
	if err != nil {
		return fmt.Errorf("Error Updating StatusCake Pagespeed: %s", err.Error())
	}
	return nil
}

func DeletePagespeed(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*statuscake.Client)

	log.Printf("[DEBUG] Deleting StatusCake Pagespeed: %s", d.Id())
	pagespeedId, _ := strconv.Atoi(d.Id())
	err := statuscake.NewPageSpeeds(client).Delete(pagespeedId)

	return err
}

func ReadPagespeed(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*statuscake.Client)
	pagespeedId, _ := strconv.Atoi(d.Id())
	response, err := statuscake.NewPageSpeeds(client).Detail(pagespeedId)
	if err != nil {
		return fmt.Errorf("Error Getting StatusCake Pagespeed Details for %s: Error: %s", d.Id(), err)
	}
	d.Set("name", response.Name)
	d.Set("website_url", response.Website_url)
	d.Set("location_iso", response.Location_iso)
	d.Set("check_rate", response.Checkrate)
	s := strconv.Itoa(response.ID)
	d.SetId(s)

	return nil
}

func getStatusCakePagespeedInput(d *schema.ResourceData) *statuscake.PartialPageSpeed {
	pagespeedId, parseErr := strconv.Atoi(d.Id())
	if parseErr != nil {
		log.Printf("[DEBUG] Error Parsing StatusCake Id: %s", d.Id())
	}
	pagespeed := &statuscake.PartialPageSpeed{
		ID: pagespeedId,
	}

	if v, ok := d.GetOk("check_rate"); ok {
		pagespeed.Checkrate = strconv.Itoa(v.(int))
	}

	if v, ok := d.GetOk("contact_group"); ok {
		pagespeed.ContactGroupsC = strings.Join(castSetToSliceStrings(v.(*schema.Set).List()), ",")
	}

	if v, ok := d.GetOk("location_iso"); ok {
		pagespeed.Location_iso = v.(string)
	}

	if v, ok := d.GetOk("name"); ok {
		pagespeed.Name = v.(string)
	}

	if v, ok := d.GetOk("website_url"); ok {
		pagespeed.Website_url = v.(string)
	}

	return pagespeed
}
