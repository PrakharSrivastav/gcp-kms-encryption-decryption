package main

import (
	"context"
	"fmt"
	"log"

	"cloud.google.com/go/datastore"
	cloudkms "cloud.google.com/go/kms/apiv1"
	"google.golang.org/api/option"
	kmspb "google.golang.org/genproto/googleapis/cloud/kms/v1"
)

const (
	KeyName            string = "my-key"         // replace this as per your project.
	Location           string = "global"         // replace this as per your project.
	KeyRingName        string = "my-key-ring"    // replace this as per your project.
	ProjectID          string = "my-gcp-project" // replace this as per your project.
	DatastoreAdminPath string = "../datastore-user.json"
	DatastoreKind      string = "users"
	DatastoreKeyName   string = "user-1"
	KMSEndDecPath      string = "../kms-enc-dec.json"
	CryptoKeyParent    string = "projects/%s/locations/%s/keyRings/%s/cryptoKeys/%s"
)

func main() {
	ctx := context.Background()
	// create KMS client
	client, err := cloudkms.NewKeyManagementClient(ctx, option.WithCredentialsFile(KMSEndDecPath))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	// Encrypt the password field
	req := &kmspb.EncryptRequest{
		Name:      fmt.Sprintf(CryptoKeyParent, ProjectID, Location, KeyRingName, KeyName),
		Plaintext: []byte("Hello 123 World"),
	}
	resp, err := client.Encrypt(ctx, req)
	if err != nil {
		log.Fatal(err)
	}

	// replace the plaintext password with the encrypted password
	t := User{Username: "123", Password: resp.GetCiphertext()}

	// save the user kind in database with encrypted password
	var ds *datastore.Client
	if ds, err = datastore.NewClient(ctx, ProjectID, option.WithCredentialsFile(DatastoreAdminPath)); err != nil {
		log.Fatal(err)
	}
	defer ds.Close()
	key := datastore.NameKey(DatastoreKind, DatastoreKeyName, nil)
	if _, err = ds.Put(ctx, key, &t); err != nil {
		log.Fatalf("ds err %v", err)
	}

}

type User struct {
	Username string `json:"username"`
	Password []byte `json:"password"`
}

// "projects/planday-api-dev/locations/global/keyRings/test-key-ring/cryptoKeys/test-key-id"
