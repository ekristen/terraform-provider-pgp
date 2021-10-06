# Terraform Provider PGP

**Warning:** Use of this provider will result in secrets being in terraform state in **PLAIN TEXT** (aka **NOT ENCRYPTED**). You've been warned.

There are use cases and situations where you need full access to all values generated within terraform, unfortunately there are some resources that force you to provide a PGP key and it will only encrypt and store those values, then manual commands must be run to decrypt.

This provider allows you to generate a PGP or use an existing one, from there it provides encrypt and decrypt data sources to allow you to get access to the data.

## Build provider

Run the following command to build the provider

```shell
$ go build -o terraform-provider-pgp
```

## Local release build

```shell
$ go install github.com/goreleaser/goreleaser@latest
```

```shell
$ make release
```

You will find the releases in the `/dist` directory. You will need to rename the provider binary to `terraform-provider-gpg` and move the binary into [the appropriate subdirectory within the user plugins directory](https://learn.hashicorp.com/tutorials/terraform/provider-use?in=terraform/providers#install-hashicups-provider).

## Test sample configuration

First, build and install the provider.

```shell
$ make install
```

Then, navigate to the `examples` directory.

```shell
$ cd examples
```

Run the following command to initialize the workspace and apply the sample configuration.

```shell
$ terraform init && terraform apply
```

**Note:** you might have to remove the `.terraform.lock.hcl` file.
