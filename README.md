# EDC Terraform Provider

The EDC provider allows Terraform to manage EDC resources.

## How to use it

- Setup local development environment

    1- Find your `GOBIN` using the following command:
    ```bash
    go env GOBIN
    ```
    If nothing prompt, you can use the standard format, which is `$HOME/go/bin`.
    2- Update your `$HOME/.terraformrc` file with this configuration:
    ```
    provider_installation {
        dev_overrides {
            "think-it-labs/edc" = "/Users/<Username>/go/bin"
        }
        direct {}
    }
    ```
    3- In the project directory, run the following command:
    ```bash
    go install .
    ```

    4- Deploy the EDC connectors using `docker-compose` using the following command:

    ```bash
    docker-compose up -d
    ```

- The terraform provider in action

    1- Add the provider declaration in the `main.tf` file, that you should create:

    ```hcl
    terraform {
        required_providers {
            edc = {
                source = "Think-iT-Labs/edc"
            }
        }
    }

    provider "edc" {
        token = "1234"
        addresses = {
            default    = "http://localhost:29193/api"
            management = "http://localhost:29193/api/v1/data"
            protocol   = "http://localhost:29193/api/v1/ids"
            public     = "http://localhost:29193/public"
            control    = "http://localhost:29193/control"
        }
    }
    ```
    2- Add the following block in the `main.tf` file:
    ```hcl
    resource "edc_asset" "s3" {
        asset = {
            "asset:prop:name" : "S3 with generated id",
            "asset:prop:contenttype" : "application/json",
        }

        data = {
            s3 = {
                name              = "test"
                bucket_name       = "test"
                access_key_id     = "dummy_key"
                secret_access_key = "dummy_key"
            }
        }
    }

    resource "edc_policy" "policy" {
        id = "abcdPolicy"
        policy = {
            uid = "231802-bb34-11ec-8422-0242ac120002",
            permissions = [
                {
                    edctype = "dataspaceconnector:permission",
                    target  = "assetId",
                    action = {
                        type = "USE"
                    },
                }
            ]
        }
    }


    resource "edc_contract_definition" "name" {
        access_policy_id   = edc_policy.policy.id
        contract_policy_id = edc_policy.policy.id
        validity           = 600
        criteria = [
            {
                operand_left  = "asset:prop:id"
                operator      = "="
                operand_right = edc_asset.s3.id
            }
        ]
    }
    ```
    3- Run the example by running the following command:
    ```
    terraform apply
    ```
## Docs
For more information, you can check the `docs` folder to see more detailed examples.
