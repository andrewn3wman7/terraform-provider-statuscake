package statuscake

import (
    "fmt"
    "log"
    "strconv"

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
            "alert_smaller": {
                Type:     schema.TypeInt,
                Optional: true,
                Default:  0,
            },

            "alert_bigger": {
                Type:     schema.TypeInt,
                Optional: true,
                Default:  0,
            },

            "alert_slower": {
                Type:     schema.TypeInt,
                Optional: true,
                Default:  0,
            },

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

    newPagespeed := &statuscake.PageSpeed{
        AlertBigger:  d.Get("alert_bigger").(int),
        AlertSmaller: d.Get("alert_smaller").(int),
        AlertSlower: d.Get("alert_slower").(int),
        Checkrate:    d.Get("check_rate").(int),
        Location_iso: d.Get("location_iso").(string),
        Name:         d.Get("name").(string),
        Website_url:  d.Get("website_url").(string),
    }

    if v, ok := d.GetOk("contact_group"); ok {
        newPagespeed.ContactGroup = castSetToSliceStrings(v.(*schema.Set).List())   
    }

    log.Printf("[DEBUG] Creating new StatusCake Pagespeed: %s", d.Get("name").(string))

    response, err := statuscake.NewPageSpeeds(client).Create(newPagespeed)
    if err != nil {
        return fmt.Errorf("Error creating StatusCake Pagespeed: %s", err.Error())
    }

    d.Set("pagespeed_id", response.ID)
    d.SetId(fmt.Sprintf("%d", response.ID))
    
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
    d.Set("alert_smaller", response.AlertSmaller)
    d.Set("alert_bigger", response.AlertBigger)
    d.Set("alert_slower", response.AlertSlower)
    d.Set("name", response.Name)
    d.Set("website_url", response.Website_url)
    d.Set("location_iso", response.Location_iso)
    d.Set("check_rate", response.Checkrate)
    d.Set("contact_group", response.ContactGroup)

    return nil
}

func getStatusCakePagespeedInput(d *schema.ResourceData) *statuscake.PageSpeed {
    pagespeedId, parseErr := strconv.Atoi(d.Id())
    if parseErr != nil {
        log.Printf("[DEBUG] Error Parsing StatusCake Id: %s", d.Id())
    }
    pagespeed := &statuscake.PageSpeed{
        ID: pagespeedId,
    }

    if v, ok := d.GetOk("alert_smaller"); ok {
        pagespeed.AlertSmaller = v.(int)
    }

    if v, ok := d.GetOk("alert_bigger"); ok {
        pagespeed.AlertBigger = v.(int)
    }

    if v, ok := d.GetOk("alert_slower"); ok {
        pagespeed.AlertSlower = v.(int)
    }

    if v, ok := d.GetOk("check_rate"); ok {
        pagespeed.Checkrate = v.(int)
    }

    if v, ok := d.GetOk("contact_group"); ok {
        pagespeed.ContactGroup = castSetToSliceStrings(v.(*schema.Set).List())
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

