# How to Run the Example

This example creates a speech account and uploads a custom json file to the custom commands service.

## Terraform

``` bash
terraform init
terraform plan -out=tf.plan -var-file=variables.tfvars
terraform apply tf.plan
```
