---
page_title: "pgp_key Resource - terraform-provider-pgp"
subcategory: ""
description: |-
  The pgp_key allows you to create a PGP keypair.
---

# Resource `pgp_key`

The pgp_key resource creates PGP keypair.

## Example Usage

```terraform

resource "pgp_key" "example" {
  name    = "John Doe"
  email   = "jdoe@exammple.com"
  comment = "Generated PGP Key"
}
```

## Argument Reference

- `name` - (Required) Name for PGP Key.
- `email` - (Required) Email for PGP Key.
- `comment` - (Required) Comment for PGP Key.

## Attributes Reference

In addition to all the arguments above, the following attributes are exported.

- `private_key` - The private key in PGP armored format.
- `private_key_base64` - The private key in base64 format.
- `public_key` - The public key in PGP armored format.
- `public_key_base64` - The public key in base64 format.
