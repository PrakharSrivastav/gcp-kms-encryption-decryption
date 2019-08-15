package main

import (
	"context"
	"log"

	"cloud.google.com/go/datastore"
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

	req := &kmspb.EncryptRequest{
		Name:      "projects/planday-api-dev/locations/global/keyRings/test-key-ring/cryptoKeys/test-key-id",
		Plaintext: []byte("Hello 123 World"),
	}
	resp, err := client.Encrypt(ctx, req)
	if err != nil {
		log.Fatal(err)
	}

	var ds *datastore.Client
	if ds, err = datastore.NewClient(ctx, "planday-api-dev", option.WithCredentialsFile("datastore.json")); err != nil {
		log.Fatal(err)
	}
	defer ds.Close()
	key := datastore.NameKey("testpassword", "testname", nil)

	t := TestPassword{
		Username: "123",
		Password: resp.GetCiphertext(),
	}

	if _, err = ds.Put(ctx, key, &t); err != nil {
		log.Fatalf("ds err %v", err)
	}

}

type TestPassword struct {
	Username string `json:"username"`
	Password []byte `json:"password"`
}

// "projects/planday-api-dev/locations/global/keyRings/test-key-ring/cryptoKeys/test-key-id"
