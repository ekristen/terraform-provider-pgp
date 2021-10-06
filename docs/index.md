---
page_title: "Provider: PGP"
subcategory: ""
description: |-
  Terraform provider for interacting with PGP.
---

# PGP Terraform Provider

**Warning:** Use of this provider will result in secrets being in terraform state in **PLAIN TEXT** (aka **NOT ENCRYPTED**). You've been warned.

There are use cases and situations where you need full access to all values generated within terraform, unfortunately there are some resources that force you to provide a PGP key and it will only encrypt and store those values, then manual commands must be run to decrypt.

This provider allows you to generate a PGP or use an existing one, from there it provides encrypt and decrypt data sources to allow you to get access to the data.
