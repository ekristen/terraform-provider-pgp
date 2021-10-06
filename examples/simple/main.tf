terraform {
  required_providers {
    pgp = {
      version = "0.2.0"
      source  = "hashicorp.com/ekristen/pgp"
    }
  }
}

provider "pgp" {}

resource "pgp_key" "testing" {
  name    = "testing"
  email   = "testing@testing.com"
  comment = "testing"
}

data "pgp_encrypt" "testing" {
  plaintext  = "thisisasecret"
  public_key = pgp_key.testing.public_key
}

data "pgp_decrypt" "testing" {
  ciphertext  = data.pgp_encrypt.testing.ciphertext
  private_key = pgp_key.testing.private_key
}

output "public_key" {
  value = pgp_key.testing.public_key
}

output "private_key" {
  value = pgp_key.testing.private_key
}

output "private_key_base64" {
  value = pgp_key.testing.private_key_base64
}

output "ciphertext" {
  value = data.pgp_encrypt.testing.ciphertext
}
output "plaintext" {
  value = data.pgp_decrypt.testing.plaintext
}
