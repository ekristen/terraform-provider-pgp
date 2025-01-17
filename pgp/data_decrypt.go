package pgp

import (
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"

	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"
	"golang.org/x/crypto/openpgp/packet"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceDecrypt() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceDecryptRead,

		Schema: map[string]*schema.Schema{
			"plaintext": {
				Type:     schema.TypeString,
				Computed: true,
				Sensitive: true,
			},
			"private_key": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ciphertext": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ciphertext_encoding": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "armored",
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(string)

					if v != "armored" && v != "base64" {
						errs = append(errs, fmt.Errorf("%q must be either 'armored' or 'base64', got: %s", key, v))
					}

					return
				},
			},
		},
	}
}

func dataSourceDecryptRead(d *schema.ResourceData, meta interface{}) error {
	rawPrivateKey := d.Get("private_key").(string)

	privateKeyPacket, err := getPrivateKeyPacket([]byte(rawPrivateKey))
	if err != nil {
		return err
	}

	encoding := d.Get("ciphertext_encoding").(string)
	ciphertext := []byte(d.Get("ciphertext").(string))

	if encoding == "base64" {
		c, err := base64.StdEncoding.DecodeString(string(ciphertext))
		if err != nil {
			return errwrap.Wrapf("unable to decode: {{err}}", err)
		}
		ciphertext = c
	}

	plaintext, err := decrypt(privateKeyPacket, ciphertext, encoding)
	if err != nil {
		return err
	}

	hash := sha256.New()
	hash.Write(plaintext)

	d.SetId(fmt.Sprintf("%x", hash.Sum(nil)))
	d.Set("plaintext", string(plaintext))

	return nil
}

// Parts below borrowed from https://github.com/jchavannes/go-pgp

func getPrivateKeyPacket(privateKey []byte) (*openpgp.Entity, error) {
	privateKeyReader := bytes.NewReader(privateKey)
	block, err := armor.Decode(privateKeyReader)
	if err != nil {
		return nil, err
	}

	if block.Type != openpgp.PrivateKeyType {
		return nil, errors.New("Invalid private key data")
	}

	packetReader := packet.NewReader(block.Body)
	return openpgp.ReadEntity(packetReader)
}

func decrypt(entity *openpgp.Entity, encrypted []byte, encoding string) ([]byte, error) {
	// Decrypt message
	entityList := openpgp.EntityList{entity}

	var messageReader *openpgp.MessageDetails
	var err error

	if encoding == "armored" {
		// Decode message
		block, err := armor.Decode(bytes.NewReader(encrypted))
		if err != nil {
			return []byte{}, fmt.Errorf("Error decoding: %v", err)
		}
		if block.Type != "Message" {
			return []byte{}, errors.New("Invalid message type")
		}

		messageReader, err = openpgp.ReadMessage(block.Body, entityList, nil, nil)
		if err != nil {
			return []byte{}, fmt.Errorf("Error reading message: %v", err)
		}
	} else {
		messageReader, err = openpgp.ReadMessage(bytes.NewReader(encrypted), entityList, nil, nil)
		if err != nil {
			return []byte{}, fmt.Errorf("Error reading message: %v", err)
		}
	}

	read, err := ioutil.ReadAll(messageReader.UnverifiedBody)
	if err != nil {
		return []byte{}, fmt.Errorf("Error reading unverified body: %v", err)
	}

	if encoding == "armored" {
		// Uncompress message
		reader := bytes.NewReader(read)
		uncompressed, err := gzip.NewReader(reader)
		if err != nil {
			return []byte{}, fmt.Errorf("Error initializing gzip reader: %v", err)
		}
		defer uncompressed.Close()

		out, err := ioutil.ReadAll(uncompressed)
		if err != nil {
			return []byte{}, err
		}

		// Return output - an unencoded, unencrypted, and uncompressed message
		return out, nil
	}

	out, err := ioutil.ReadAll(bytes.NewReader(read))
	if err != nil {
		return []byte{}, err
	}

	// Return output - an unencoded, unencrypted, and uncompressed message
	return out, nil
}
