# Terraform stuff

This directory contains the terraform artifacts to build the project infrastructure.

## Project structure.

We are going to configure the directory structure based on the technical components of the infrastructure and within each one we are going to organize each application and service that requires these instruments.

```log
infra/terraform
├── README.md
├── messaging
│   └── README.md
├── monitoring
│   └── README.md
├── networking
│   └── README.md
├── security
│   └── README.md
└── storage
    └── README.md
```

## How to test locally?

To test the solution locally, we are mainly using [localstack](https://docs.localstack.cloud/overview/). Below what you need to do to make it work (you can also follow this [tutorial](https://docs.localstack.cloud/integrations/terraform/)).

0. Install terraform as suggested [here](https://learn.hashicorp.com/tutorials/terraform/install-cli). Personally I didn't use `brew` because manual installation is pretty straightforward.

1. specify mock credentials for the aws provider.

```yml
[default]
aws_access_key_id     = "test"
aws_secret_access_key = "test"
region                = us-west-2
```

2. Second, localstack people say 

> we need to avoid issues with routing and authentication (as we do not need it). Therefore we need to supply some general parameters:

so, in the `versions.tf` file, we are adding this:

```yml
...
provider "aws" {

  access_key = "test"
  secret_key = "test"
  region     = "us-west-2"


  # only required for non virtual hosted-style endpoint use case.
  # https://registry.terraform.io/providers/hashicorp/aws/latest/docs#s3_force_path_style
  s3_use_path_style           = true
  skip_credentials_validation = true
  skip_metadata_api_check     = true
  skip_requesting_account_id  = true

  endpoints {
    s3  = "http://s3.localhost.localstack.cloud:4566"
    sns = "http://localhost:4566"
    sqs = "http://localhost:4566"
  }
}
```

3. Run `localstack`.

```sh
docker run --rm -it -p 4566:4566 -p 4571:4571 localstack/localstack
```

4. Go to `./infra/terraform/messaging/core` and run `terraform init`.

> Initializing a configuration directory downloads and installs the providers defined in the configuration, which in this case is the aws provider.

you will see something like this

```log
Terraform has been successfully initialized!

You may now begin working with Terraform. Try running "terraform plan" to see
any changes that are required for your infrastructure. All Terraform commands
should now work.

If you ever set or change modules or backend configuration for Terraform,
rerun this command to reinitialize your working directory. If you forget, other
commands will detect it and remind you to do so if necessary.
```

5. now is the time for the `fmt` command. It is recommended by `terraform` in order to update configurations in the current directory for readability and consistency.

```log
➜  terraform fmt
main.tf
variables.tf
```

6. Now, let's run the `terraform validate` command.

you will see this:

```log
➜  terraform validate
Success! The configuration is valid.
```

7. Next, run the `terraform plan` command.

```log
...
Plan: 1 to add, 0 to change, 0 to destroy.

Changes to Outputs:
  + sns_arn = (known after apply)
  + sns_id  = (known after apply)
```

8. Then, the `terraform apply` to create the infrastructure.

```log
...
Changes to Outputs:
  + sns_arn = (known after apply)
  + sns_id  = (known after apply)

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes

aws_sns_topic.orders-fifo: Creating...
aws_sns_topic.orders-fifo: Creation complete after 0s [id=arn:aws:sns:us-west-2:000000000000:orders.fifo]

Apply complete! Resources: 1 added, 0 changed, 0 destroyed.

Outputs:

sns_arn = "arn:aws:sns:us-west-2:000000000000:orders.fifo"
sns_id = "arn:aws:sns:us-west-2:000000000000:orders.fifo"
```

9. Let's see in `localstack` if the `sns topic` was created.

```log
➜  aws sns list-topics --endpoint=http://localhost:4566 --region us-west-2
{
    "Topics": [
        {
            "TopicArn": "arn:aws:sns:us-west-2:000000000000:orders.fifo"
        }
    ]
}
```

10. Finally, just for love at art. the `destroy` command.

```sh
➜  terraform destroy
aws_sns_topic.orders-fifo: Refreshing state... [id=arn:aws:sns:us-west-2:000000000000:orders.fifo]
...
Do you really want to destroy all resources?
  Terraform will destroy all your managed infrastructure, as shown above.
  There is no undo. Only 'yes' will be accepted to confirm.

  Enter a value: yes

aws_sns_topic.orders-fifo: Destroying... [id=arn:aws:sns:us-west-2:000000000000:orders.fifo]
aws_sns_topic.orders-fifo: Destruction complete after 0s

Destroy complete! Resources: 1 destroyed.
```

checking

```sh
➜  aws sns list-topics --endpoint=http://localhost:4566 --region us-west-2
{
    "Topics": []
}
```