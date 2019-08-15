package main

import (
	"context"
	"fmt"
	"log"

	cloudkms "cloud.google.com/go/kms/apiv1"
	"google.golang.org/api/option"
	kmspb "google.golang.org/genproto/googleapis/cloud/kms/v1"
)

func main() {
	ctx := context.Background()
	client, err := cloudkms.NewKeyManagementClient(ctx, option.WithCredentialsFile("kms.json"))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()
	 parentName := fmt.Sprintf("projects/%s/locations/%s", "planday-api-dev", "global")
	// Call the API.
	request := &kmspb.CreateKeyRingRequest{
		Parent:    parentName,
		KeyRingId: "test-key-ring",
	}

	result, err := client.CreateKeyRing(ctx, request)
	if err != nil {
		log.Printf("error creating keyring %s", err)
	}
	log.Printf("Created key ring: %s", result)

	/*req := &kmspb.ListKeyRingsRequest{
		Parent: parentName,
	}
	// Query the API.
	it := client.ListKeyRings(ctx, req)

	// Iterate and print results.
	for {
		resp, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Failed to list key rings: %v", err)
		}
		log.Printf("KeyRing: %q\n", resp.Name)
	}*/

	r := &kmspb.CreateCryptoKeyRequest{
		Parent:      fmt.Sprintf("projects/%s/locations/%s/keyRings/%s", "planday-api-dev", "global", "test-key-ring"),
		CryptoKeyId: "test-key-id",
		CryptoKey: &kmspb.CryptoKey{
			Purpose: kmspb.CryptoKey_ENCRYPT_DECRYPT,
			VersionTemplate: &kmspb.CryptoKeyVersionTemplate{
				Algorithm: kmspb.CryptoKeyVersion_GOOGLE_SYMMETRIC_ENCRYPTION,
			},
		},
	}

	res, err := client.CreateCryptoKey(ctx, r)
	if err != nil {
		log.Println(err)
	}
	log.Printf("created crypto key %s", res)


}
