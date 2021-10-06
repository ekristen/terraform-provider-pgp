package pgp

import (
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"

	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"
	"golang.org/x/crypto/openpgp/packet"

	_ "golang.org/x/crypto/ripemd160"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceEncrypt() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceEncryptRead,

		Schema: map[string]*schema.Schema{
			"plaintext": {
				Type:     schema.TypeString,
				Required: true,
			},
			"public_key": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ciphertext": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceEncryptRead(d *schema.ResourceData, meta interface{}) error {
	rawPublicKey := d.Get("public_key").(string)

	publicKeyPacket, err := getPublicKeyPacket([]byte(rawPublicKey))
	if err != nil {
		return err
	}

	ciphertext, err := encrypt(publicKeyPacket, []byte(d.Get("plaintext").(string)))
	if err != nil {
		return err
	}

	hash := sha256.New()
	hash.Write(ciphertext)

	d.SetId(fmt.Sprintf("%x", hash.Sum(nil)))
	d.Set("ciphertext", string(ciphertext))

	return nil
}

// Parts below borrowed from https://github.com/jchavannes/go-pgp

func getPublicKeyPacket(publicKey []byte) (*openpgp.Entity, error) {
	publicKeyReader := bytes.NewReader(publicKey)
	block, err := armor.Decode(publicKeyReader)
	if err != nil {
		return nil, err
	}

	if block.Type != openpgp.PublicKeyType {
		return nil, errors.New("Invalid public key data")
	}

	packetReader := packet.NewReader(block.Body)
	return openpgp.ReadEntity(packetReader)
}

func encrypt(entity *openpgp.Entity, plaintext []byte) ([]byte, error) {
	buf := new(bytes.Buffer)

	// Create encoder
	encoderWriter, err := armor.Encode(buf, "Message", make(map[string]string))
	if err != nil {
		return []byte{}, fmt.Errorf("Error creating OpenPGP armor: %v", err)
	}

	// Create encryptor with encoder
	encryptorWriter, err := openpgp.Encrypt(encoderWriter, []*openpgp.Entity{entity}, nil, nil, nil)
	if err != nil {
		return []byte{}, fmt.Errorf("Error creating entity for encryption: %v", err)
	}

	// Create compressor with encryptor
	compressorWriter, err := gzip.NewWriterLevel(encryptorWriter, gzip.BestCompression)
	if err != nil {
		return []byte{}, fmt.Errorf("Invalid compression level: %v", err)
	}

	// Write message to compressor
	messageReader := bytes.NewReader(plaintext)
	_, err = io.Copy(compressorWriter, messageReader)
	if err != nil {
		return []byte{}, fmt.Errorf("Error writing data to compressor: %v", err)
	}

	compressorWriter.Close()
	encryptorWriter.Close()
	encoderWriter.Close()

	// Return buffer output - an encoded, encrypted, and compressed message
	return buf.Bytes(), nil
}
