package task_management_worktype

import (
	"context"
	"fmt"
	"terraform-provider-genesyscloud/genesyscloud/provider"
	"terraform-provider-genesyscloud/genesyscloud/util"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

/*
   The data_source_genesyscloud_task_management_worktype.go contains the data source implementation
   for the resource.
*/

// dataSourceTaskManagementWorktypeRead retrieves by name the id in question
func dataSourceTaskManagementWorktypeRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sdkConfig := meta.(*provider.ProviderMeta).ClientConfig
	proxy := newTaskManagementWorktypeProxy(sdkConfig)

	name := d.Get("name").(string)

	return util.WithRetries(ctx, 15*time.Second, func() *retry.RetryError {
		worktypeId, retryable, err := proxy.getTaskManagementWorktypeIdByName(ctx, name)

		if err != nil && !retryable {
			return retry.NonRetryableError(fmt.Errorf("error searching task management worktype %s: %s", name, err))
		}

		if retryable {
			return retry.RetryableError(fmt.Errorf("no task management worktype found with name %s", name))
		}

		d.SetId(worktypeId)
		return nil
	})
}
