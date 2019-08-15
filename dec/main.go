package main

import (
	"context"
	"log"

	cloudkms "cloud.google.com/go/kms/apiv1"
	kmspb "google.golang.org/genproto/googleapis/cloud/kms/v1"

	"cloud.google.com/go/datastore"
	"google.golang.org/api/option"
)

func main() {
	var ds *datastore.Client
	var err error
	ctx := context.Background()
	if ds, err = datastore.NewClient(ctx, "planday-api-dev", option.WithCredentialsFile("datastore.json")); err != nil {
		log.Fatal(err)
	}
	defer ds.Close()
	key := datastore.NameKey("testpassword", "testname", nil)

	t := TestPassword{}

	if err = ds.Get(ctx, key, &t); err != nil {
		log.Fatalf("ds err %v", err)
	}

	// log.Printf("hello %v", t.Password)

	client, err := cloudkms.NewKeyManagementClient(ctx, option.WithCredentialsFile("kms.json"))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	r := &kmspb.DecryptRequest{
		Name:       "projects/planday-api-dev/locations/global/keyRings/test-key-ring/cryptoKeys/test-key-id",
		Ciphertext: t.Password,
	}
	// Call the API.
	resp2, err := client.Decrypt(ctx, r)
	if err != nil {
		log.Println(err)
	}
	log.Println(string(resp2.Plaintext))
}

type TestPassword struct {
	Username string `json:"username"`
	Password []byte `json:"password"`
}
