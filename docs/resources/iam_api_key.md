---
page_title: "Twilio API Keys"
subcategory: "IAM"
---

# twilio_iam_api_key Resource

Manages an API Key for a Twilio Account. See the [API docs](https://www.twilio.com/docs/iam/keys/api-key-resource) for more information

!> Only Standard API Keys can be created via the API. If you require a Master API Key then you will need to create this manually in the Twilio console

## Example Usage

```hcl
resource "twilio_iam_api_key" "api_key" {
  account_sid   = "ACXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
  friendly_name = "Test API Key"
}
```

## Argument Reference

The following arguments are supported:

- `account_sid` - (Mandatory) The Account SID associated with the API Key. Changing this forces a new resource to be created
- `friendly_name` - (Optional) The name of the API Key. The default value is an empty string/ no configuration specified

## Attributes Reference

The following attributes are exported:

- `id` - The ID of the API Key (Same as the `sid`)
- `sid` - The SID of the API Key (Same as the `id`)
- `account_sid` - The Account SID associated with the API Key
- `friendly_name` - The name of the API Key
- `secret` - The API Key Secret
- `date_created` - The date in RFC3339 format that the API Key was created
- `date_updated` - The date in RFC3339 format that the API Key was updated

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

- `create` - (Defaults to 10 minutes) Used when creating the API Key
- `update` - (Defaults to 10 minutes) Used when updating the API Key
- `read` - (Defaults to 5 minutes) Used when retrieving the API Key
- `delete` - (Defaults to 10 minutes) Used when deleting the API Key

## Import

Not supported
