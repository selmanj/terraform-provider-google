package google

import "github.com/hashicorp/terraform/helper/schema"

// ResourceLabels returns a string map representing key/value pairs from a resource.
func ResourceLabels(d *schema.ResourceData) map[string]string {
	labels := map[string]string{}
	if v, ok := d.GetOk("labels"); ok {
		labelMap := v.(map[string]interface{})
		for k, v := range labelMap {
			labels[k] = v.(string)
		}
	}
	return labels
}
