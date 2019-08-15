package main

import (
	"context"
	"fmt"
	"log"

	cloudkms "cloud.google.com/go/kms/apiv1"
	"google.golang.org/api/option"
	kmspb "google.golang.org/genproto/googleapis/cloud/kms/v1"
)

const (
	location        string = "global"
	keyRingID       string = "key.ring.id"
	cryptoKeyName   string = "crypto.key.name"
	projectName     string = "gcp.project.name"
	keyRingParent   string = "projects/%s/locations/%s"
	cryptoKeyParent string = "projects/%s/locations/%s/keyRings/%s"
)

func main() {
	ctx := context.Background()
	client, err := cloudkms.NewKeyManagementClient(ctx, option.WithCredentialsFile("kms.json"))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	// first create a key ring
	parentName := fmt.Sprintf(keyRingParent, projectName, location)
	request := &kmspb.CreateKeyRingRequest{
		Parent:    parentName,
		KeyRingId: keyRingID,
	}

	result, err := client.CreateKeyRing(ctx, request)
	if err != nil {
		log.Printf("error creating keyring %s", err)
	}
	log.Printf("Created key ring: %s", result)

	// then add a key to key-ring
	r := &kmspb.CreateCryptoKeyRequest{
		Parent:      fmt.Sprintf(cryptoKeyParent, projectName, location, keyRingID),
		CryptoKeyId: cryptoKeyName,
		CryptoKey: &kmspb.CryptoKey{
			Purpose:         kmspb.CryptoKey_ENCRYPT_DECRYPT,
			VersionTemplate: &kmspb.CryptoKeyVersionTemplate{Algorithm: kmspb.CryptoKeyVersion_GOOGLE_SYMMETRIC_ENCRYPTION},
		},
	}

	// finally create the keyring
	res, err := client.CreateCryptoKey(ctx, r)
	if err != nil {
		log.Println(err)
	}

	// use this reference for encryption-decryption in rest of your project
	log.Printf("created crypto key %s", res)

}
