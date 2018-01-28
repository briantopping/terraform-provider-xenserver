package xenserver

import (
	"github.com/hashicorp/terraform/helper/schema"
	"log"
	"time"
)

const (
	vmNetworksVmUUID       = "vm_uuid"
	vmNetworksIp           = "ip"
	vmNetworksIpv6         = "ipv6"
	vmNetworksStartupDelay = "startup_delay"
)

func dataSourceVmNetworks() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVmNetworksRead,
		Schema: map[string]*schema.Schema{
			vmNetworksVmUUID: &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			vmNetworksIp: &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeList,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
			},
			vmNetworksIpv6: &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeList,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
			},
			vmNetworksStartupDelay: &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
			},
		},
	}
}

func dataSourceVmNetworksRead(d *schema.ResourceData, meta interface{}) (err error) {
	c := meta.(*Connection)

	vm := &VMDescriptor{
		UUID: d.Get(vmNetworksVmUUID).(string),
	}
	if err = vm.Load(c); err != nil {
		return err
	}

	if delay, ok := d.GetOk(vmNetworksStartupDelay); ok {
		delay := time.Duration(delay.(int)) * time.Second
		log.Printf("[DEBUG] Delaying %s\n", delay.String())
		//var vmmetrics *VMMetrics
		//if vmmetrics, err = vm.Metrics(c); err != nil {
		//	return err
		//}
		//
		//log.Printf("[DEBUG] Start time was %s\n", vmmetrics.StartTime.String())
		//now := time.Now()
		//diff := now.Sub(vmmetrics.StartTime)
		//
		//log.Printf("[DEBUG] Difference is %s\n", (delay - diff).String())
		//if delay > diff {
			time.Sleep(delay)
		//}
	}

	var metrics VMGuestMetrics

	if metrics, err = vm.GuestMetrics(c); err != nil {
		return err
	}

	d.SetId(metrics.UUID)

	log.Printf("[DEBUG] Id is %s\n", d.Id())
	log.Println("[DEBUG] Networks: ", metrics.Networks)

	ipNetworks := make([][]string, 0)
	ipv6Networks := make([][]string, 0)

	for _, network := range metrics.Networks {
		ipNetworks = append(ipNetworks, network["ip"])
		ipv6Networks = append(ipv6Networks, network["ipv6"])
	}

	d.Set(vmNetworksIp, ipNetworks)
	d.Set(vmNetworksIpv6, ipv6Networks)

	return nil
}
