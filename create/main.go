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
	Location        string = "global"              // replace this as per your project.
	KeyRingID       string = "my-key-ring"         // replace this as per your project.
	CryptoKeyName   string = "my-key"              // replace this as per your project.
	ProjectName     string = "my-gcp-project-name" // replace this as per your project.
	KeyRingParent   string = "projects/%s/locations/%s"
	CryptoKeyParent string = "projects/%s/locations/%s/keyRings/%s"
	KMSAdminPath    string = "../kms-admin.json"
)

func main() {
	ctx := context.Background()
	client, err := cloudkms.NewKeyManagementClient(ctx, option.WithCredentialsFile(KMSAdminPath))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	// first create a key ring in your project and location
	parentName := fmt.Sprintf(KeyRingParent, ProjectName, Location)
	request := &kmspb.CreateKeyRingRequest{
		Parent:    parentName,
		KeyRingId: KeyRingID,
	}

	result, err := client.CreateKeyRing(ctx, request)
	if err != nil {
		log.Printf("error creating keyring %s", err)
	}
	log.Printf("Created key ring: %s", result)

	// then add a key to key-ring
	r := &kmspb.CreateCryptoKeyRequest{
		Parent:      fmt.Sprintf(CryptoKeyParent, ProjectName, Location, KeyRingID),
		CryptoKeyId: CryptoKeyName,
		CryptoKey: &kmspb.CryptoKey{
			Purpose: kmspb.CryptoKey_ENCRYPT_DECRYPT,
			VersionTemplate: &kmspb.CryptoKeyVersionTemplate{
				Algorithm: kmspb.CryptoKeyVersion_GOOGLE_SYMMETRIC_ENCRYPTION,
			},
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
