package main

import (
	"context"
	"fmt"
	"log"

	cloudkms "cloud.google.com/go/kms/apiv1"
	kmspb "google.golang.org/genproto/googleapis/cloud/kms/v1"

	"cloud.google.com/go/datastore"
	"google.golang.org/api/option"
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
	var ds *datastore.Client
	var err error
	u := User{}
	ctx := context.Background()

	// read the user data from datastore entity
	if ds, err = datastore.NewClient(ctx, ProjectID, option.WithCredentialsFile(DatastoreAdminPath)); err != nil {
		log.Fatal(err)
	}
	defer ds.Close()
	key := datastore.NameKey(DatastoreKind, DatastoreKeyName, nil)
	if err = ds.Get(ctx, key, &u); err != nil {
		log.Fatalf("ds err %v", err)
	}

	// initialize the kms client using kms-enc-dec service account
	client, err := cloudkms.NewKeyManagementClient(ctx, option.WithCredentialsFile(KMSEndDecPath))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()


	// Call the API to decrypt password field
	r := &kmspb.DecryptRequest{Name: fmt.Sprintf(CryptoKeyParent, ProjectID, Location, KeyRingName, KeyName), Ciphertext: u.Password}
	resp, err := client.Decrypt(ctx, r)
	if err != nil {
		log.Println(err)
	}
	log.Println(string(resp.Plaintext))
}

type User struct {
	Username string `json:"username"`
	Password []byte `json:"password"`
}
