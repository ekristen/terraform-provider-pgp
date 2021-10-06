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

resource "aws_iam_user" "testing" {
  name          = "example-test"
  force_destroy = true
}

resource "aws_iam_user_login_profile" "testing" {
  user    = aws_iam_user.testing.name
  pgp_key = pgp_key.testing.public_key_base64
}

data "pgp_decrypt" "testing" {
  private_key         = pgp_key.testing.private_key
  ciphertext          = aws_iam_user_login_profile.testing.encrypted_password
  ciphertext_encoding = "base64"
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

output "password" {
  value = data.pgp_decrypt.testing.plaintext
}
