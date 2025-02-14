package task_management_workitem

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
   The data_source_genesyscloud_task_management_workitem.go contains the data source implementation
   for the resource.
*/

// dataSourceTaskManagementWorkitemRead retrieves by name the id in question
func dataSourceTaskManagementWorkitemRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sdkConfig := meta.(*provider.ProviderMeta).ClientConfig
	proxy := newTaskManagementWorkitemProxy(sdkConfig)

	name := d.Get("name").(string)
	workbinId := d.Get("workbin_id").(string)
	worktypeId := d.Get("worktype_id").(string)

	return util.WithRetries(ctx, 15*time.Second, func() *retry.RetryError {
		workitemId, retryable, err := proxy.getTaskManagementWorkitemIdByName(ctx, name, workbinId, worktypeId)

		if err != nil && !retryable {
			return retry.NonRetryableError(fmt.Errorf("error searching task management workitem %s: %s", name, err))
		}

		if retryable {
			return retry.RetryableError(fmt.Errorf("no task management workitem found with name %s", name))
		}

		d.SetId(workitemId)
		return nil
	})
}
