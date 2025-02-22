---
page_title: "Twilio Serverless Deployment"
subcategory: "Serverless"
---

# twilio_serverless_deployment Resource

Manages a Serverless deployment. See the [API docs](https://www.twilio.com/docs/runtime/functions-assets-api/api/deployment) for more information

For more information on Serverless (also known as Runtime), see the product [page](https://www.twilio.com/runtime)

~> Serverless deployments cannot be removed, they can only be superseded. To allow a build to be deleted, on the destruction of the resource, the provider will check if the `build_sid` is deployed to the environment. If the `build_sid` matches the environment config, a new deployment will be created without a `build_sid` to remove the active deployment. Once the deployment has been completed or if the `build_sid` doesn't match the environment, the state is removed and the deployment is orphaned.

~> To allow terraform to correctly manage the lifecycle of the deployment, it is recommended that use the lifecycle meta-argument `create_before_destroy` with this resource. The docs can be found [here](https://www.terraform.io/docs/configuration/resources.html#create_before_destroy)

!> This API used to manage this resource is currently in beta and is subject to change

## Example Usage

```hcl
resource "twilio_serverless_service" "service" {
  unique_name   = "twilio-test"
  friendly_name = "twilio-test"
}

resource "twilio_serverless_function" "function" {
  service_sid   = twilio_serverless_service.service.sid
  friendly_name = "test"

  content           = <<EOF
exports.handler = function (context, event, callback) {
  callback(null, "Hello World");
};
EOF
  content_type      = "application/javascript"
  content_file_name = "helloWorld.js"
  path              = "/test-function"
  visibility        = "private"
}

resource "twilio_serverless_build" "build" {
  service_sid = twilio_serverless_service.service.sid

  function_version {
    sid = twilio_serverless_function.function.latest_version_sid
  }

  dependencies = {
    "twilio" : "3.6.3"
    "fs"                      = "0.0.1-security"
    "lodash"                  = "4.17.11"
    "util"                    = "0.11.0"
    "xmldom"                  = "0.1.27"
    "@twilio/runtime-handler" = "1.0.1"
  }

  polling {
    enabled = true
  }

  lifecycle {
    create_before_destroy = true
  }
}

resource "twilio_serverless_environment" "environment" {
  service_sid = twilio_serverless_service.service.sid
  unique_name = "test"
}

resource "twilio_serverless_deployment" "deployment" {
  service_sid     = twilio_serverless_service.service.sid
  environment_sid = twilio_serverless_environment.environment.sid
  build_sid       = twilio_serverless_build.build.sid

  lifecycle {
    create_before_destroy = true
  }
}
```

## Argument Reference

The following arguments are supported:

- `service_sid` - (Mandatory) The serverless service SID to associate the deployment with. Changing this forces a new resource to be created
- `environment_sid` - (Mandatory) The serverless environment SID to associate the deployment with. Changing this forces a new resource to be created
- `build_sid` - (Optional) The build SID to be deployed to the environment. Changing this forces a new resource to be created
- `triggers` - (Optional) A map of key-value pairs which can be used to determine if changes have occurred and redeployment is necessary. Changing this forces a new resource to be created
  ~> An alternative strategy is to use the [taint](https://www.terraform.io/docs/commands/taint.html) functionality of Terraform.

## Attributes Reference

The following attributes are exported:

- `id` - The ID of the deployment (Same as the `sid`)
- `sid` - The SID of the deployment (Same as the `id`)
- `account_sid` - The account SID associated with the deployment
- `service_sid` - The service SID associated with the deployment
- `environment_sid` - The environment SID associated with the deployment
- `build_sid` - The build SID to be deployed to the environment
- `is_latest_deployment` - Determine whether this deployment is the latest
  ~> This caters for when deployments are made and Terraform state is not aware of them
- `triggers` - A map of key-value pairs which can be used to determine if changes have occurred and redeployment is necessary.
- `date_created` - The date in RFC3339 format that the deployment was created
- `date_updated` - The date in RFC3339 format that the deployment was updated
- `url` - The URL of the deployment

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

- `create` - (Defaults to 10 minutes) Used when creating the deployment
- `read` - (Defaults to 5 minutes) Used when retrieving the deployment
- `delete` - (Defaults to 10 minutes) Used when deleting the deployment

## Import

A deployment can be imported using the `/Services/{serviceSid}/Environments/{environmentSid}/Deployments/{sid}` format, e.g.

```shell
terraform import twilio_serverless_deployment.deployment /Services/ZSXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX/Environments/ZEXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX/Deployments/ZDXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX
```

!> "triggers" cannot be imported
