package chat

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/RJPearson94/terraform-provider-twilio/twilio/common"
	"github.com/RJPearson94/terraform-provider-twilio/twilio/utils"
	"github.com/RJPearson94/twilio-sdk-go/service/chat/v2/service/channel/webhook"
	"github.com/RJPearson94/twilio-sdk-go/service/chat/v2/service/channel/webhooks"
	sdkUtils "github.com/RJPearson94/twilio-sdk-go/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceChatChannelStudioWebhook() *schema.Resource {
	return &schema.Resource{
		DeprecationMessage: "Programmable Chat API will reach the end of life on 25th July 2022 (except for Flex applications), please see https://www.twilio.com/changelog/programmable-chat-end-of-life for more information",

		CreateContext: resourceChatChannelStudioWebhookCreate,
		ReadContext:   resourceChatChannelStudioWebhookRead,
		UpdateContext: resourceChatChannelStudioWebhookUpdate,
		DeleteContext: resourceChatChannelStudioWebhookDelete,

		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				format := "/Services/(.*)/Channels/(.*)/Webhooks/(.*)"
				regex := regexp.MustCompile(format)
				match := regex.FindStringSubmatch(d.Id())

				if len(match) != 4 {
					return nil, fmt.Errorf("The imported ID (%s) does not match the format (%s)", d.Id(), format)
				}

				d.Set("service_sid", match[1])
				d.Set("channel_sid", match[2])
				d.Set("sid", match[3])
				d.SetId(match[3])
				return []*schema.ResourceData{d}, nil
			},
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"sid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"account_sid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"service_sid": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"channel_sid": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"flow_sid": {
				Type:     schema.TypeString,
				Required: true,
			},
			"retry_count": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"date_created": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"date_updated": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"url": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceChatChannelStudioWebhookCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*common.TwilioClient).Chat

	createInput := &webhooks.CreateChannelWebhookInput{
		Type: "studio",
		Configuration: &webhooks.CreateChannelWebhookConfigurationInput{
			FlowSid:    utils.OptionalString(d, "flow_sid"),
			RetryCount: utils.OptionalInt(d, "retry_count"),
		},
	}

	createResult, err := client.Service(d.Get("service_sid").(string)).Channel(d.Get("channel_sid").(string)).Webhooks.CreateWithContext(ctx, createInput)
	if err != nil {
		return diag.Errorf("Failed to create chat channel webhook: %s", err.Error())
	}

	d.SetId(createResult.Sid)
	return resourceChatChannelStudioWebhookRead(ctx, d, meta)
}

func resourceChatChannelStudioWebhookRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*common.TwilioClient).Chat

	getResponse, err := client.Service(d.Get("service_sid").(string)).Channel(d.Get("channel_sid").(string)).Webhook(d.Id()).FetchWithContext(ctx)
	if err != nil {
		if twilioError, ok := err.(*sdkUtils.TwilioError); ok {
			// currently programmable chat returns a 403 if the service instance does not exist
			if (twilioError.Status == 403 && twilioError.Message == "Service instance not found") || twilioError.IsNotFoundError() {
				d.SetId("")
				return nil
			}
		}
		return diag.Errorf("Failed to read chat channel webhook: %s", err.Error())
	}

	d.Set("sid", getResponse.Sid)
	d.Set("account_sid", getResponse.AccountSid)
	d.Set("service_sid", getResponse.ServiceSid)
	d.Set("channel_sid", getResponse.ChannelSid)
	d.Set("type", getResponse.Type)
	d.Set("flow_sid", getResponse.Configuration.FlowSid)
	d.Set("retry_count", getResponse.Configuration.RetryCount)
	d.Set("date_created", getResponse.DateCreated.Format(time.RFC3339))

	if getResponse.DateUpdated != nil {
		d.Set("date_updated", getResponse.DateUpdated.Format(time.RFC3339))
	}

	d.Set("url", getResponse.URL)

	return nil
}

func resourceChatChannelStudioWebhookUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*common.TwilioClient).Chat

	updateInput := &webhook.UpdateChannelWebhookInput{
		Configuration: &webhook.UpdateChannelWebhookConfigurationInput{
			FlowSid:    utils.OptionalString(d, "flow_sid"),
			RetryCount: utils.OptionalInt(d, "retry_count"),
		},
	}

	updateResp, err := client.Service(d.Get("service_sid").(string)).Channel(d.Get("channel_sid").(string)).Webhook(d.Id()).UpdateWithContext(ctx, updateInput)
	if err != nil {
		return diag.Errorf("Failed to update chat channel webhook: %s", err.Error())
	}

	d.SetId(updateResp.Sid)
	return resourceChatChannelStudioWebhookRead(ctx, d, meta)
}

func resourceChatChannelStudioWebhookDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*common.TwilioClient).Chat

	if err := client.Service(d.Get("service_sid").(string)).Channel(d.Get("channel_sid").(string)).Webhook(d.Id()).DeleteWithContext(ctx); err != nil {
		return diag.Errorf("Failed to delete chat channel webhook: %s", err.Error())
	}
	d.SetId("")
	return nil
}
