---
page_title: "pgp_decrypt Data Source - terraform-provider-pgp"
subcategory: ""
description: |-
  The pgp_decrypt data source allows you to decrypt PGP encrypted data.
---

# Data Source `pgp_decrypt`

The pgp_decrypt data source allows you to decrypt PGP encrypted data.

## Example Usage

```terraform
data "pgp_decrypt" "example" {
  ciphertext  = "ENCRYPTED..."
  private_key = pgp_key.example.private_key
}
```

## Argument Reference

- `private_key` - (Required) PGP Private Key in Armored Format
- `ciphertext` - (Required) Ciphertext to be decrypted by the Private Key

## Attributes Reference

In addition to all the arguments above, the following attributes are exported.

- `plaintext` - The decrypted data.
