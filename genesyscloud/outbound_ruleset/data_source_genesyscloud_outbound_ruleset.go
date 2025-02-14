package outbound_ruleset

import (
	"context"
	"fmt"
	"terraform-provider-genesyscloud/genesyscloud/provider"
	"terraform-provider-genesyscloud/genesyscloud/util"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

/*
   The data_source_genesyscloud_outbound_ruleset.go contains the data source implementation
   for the resource.
*/

// dataSourceOutboundRulesetRead retrieves by name the id in question
func dataSourceOutboundRulesetRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sdkConfig := meta.(*provider.ProviderMeta).ClientConfig
	proxy := newOutboundRulesetProxy(sdkConfig)

	name := d.Get("name").(string)

	return util.WithRetries(ctx, 15*time.Second, func() *retry.RetryError {
		rulesetId, retryable, err := proxy.getOutboundRulesetIdByName(ctx, name)

		if err != nil && !retryable {
			return retry.NonRetryableError(fmt.Errorf("Error ruleset %s: %s", name, err))
		}

		if retryable {
			return retry.RetryableError(fmt.Errorf("No ruleset found with name %s", name))
		}

		d.SetId(rulesetId)
		return nil
	})
}
