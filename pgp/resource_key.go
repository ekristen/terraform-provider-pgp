package pgp

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"fmt"

	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceKey() *schema.Resource {
	return &schema.Resource{
		Create: resourceKeyCreate,
		Read:   schema.Noop,
		Delete: schema.RemoveFromState,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(0, 256),
			},
			"comment": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"email": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},

			"public_key": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Armored PGP Public Key",
			},
			"public_key_base64": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Base64 Encoded Public Key",
			},
			"private_key": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Armored PGP Private Key",
			},
			"private_key_base64": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Base64 Encoded Private Key",
			},
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceKeyCreate(d *schema.ResourceData, meta interface{}) error {
	name := d.Get("name").(string)
	comment := d.Get("comment").(string)
	email := d.Get("email").(string)

	e, err := openpgp.NewEntity(name, comment, email, nil)
	if err != nil {
		return errwrap.Wrapf("error generating pgp: {{err}}", err)
	}

	for _, id := range e.Identities {
		err := id.SelfSignature.SignUserId(id.UserId.Id, e.PrimaryKey, e.PrivateKey, nil)
		if err != nil {
			return errwrap.Wrapf("error signing pgp keys: {{err}}", err)
		}
	}

	// PublicKey
	b64buf := new(bytes.Buffer)
	b64w := bufio.NewWriter(b64buf)

	buf := new(bytes.Buffer)
	w, err := armor.Encode(buf, openpgp.PublicKeyType, nil)
	if err != nil {
		return errwrap.Wrapf("error armor pgp keys: {{err}}", err)
	}

	e.Serialize(w)
	e.Serialize(b64w)

	w.Close()
	b64w.Flush()

	base64PubKey := base64.StdEncoding.EncodeToString(b64buf.Bytes())
	armoredPubKey := buf.String()

	// PrivateKey
	b64buf = new(bytes.Buffer)
	b64w = bufio.NewWriter(b64buf)

	buf = new(bytes.Buffer)
	w, err = armor.Encode(buf, openpgp.PrivateKeyType, nil)
	if err != nil {
		return errwrap.Wrapf("error armor pgp keys: {{err}}", err)
	}

	e.SerializePrivate(w, nil)
	e.SerializePrivate(b64w, nil)

	w.Close()
	b64w.Flush()

	base64PrivateKey := base64.StdEncoding.EncodeToString(b64buf.Bytes())
	armoredPrivateKey := buf.String()

	d.SetId(fmt.Sprintf("%x", e.PrimaryKey.Fingerprint))

	d.Set("public_key", armoredPubKey)
	d.Set("public_key_base64", base64PubKey)

	d.Set("private_key", armoredPrivateKey)
	d.Set("private_key_base64", base64PrivateKey)

	return nil
}
