# gcp-enc-dec

This is a demo application to support [this](https://www.prakharsrivastav.com/posts/gcp-using-kms-to-manage-secrets/) blog post. It demonstrates how to use GCP Cloud KMS to easily manage secrets.

## setup

Create 3 service accounts with the below permissions
- [Cloud KMS Admin](docs/kms-admin-2.png)
- [Cloud KMS CryptoKey Encrypter/Decrypter](docs/kms-encdec-2.png)
- [Cloud Datastore User](docs/datastore-user.png)

1. Save `Cloud KMS Admin` service account key as kms-admin.json in the project root.
2. Save `Cloud KMS CryptoKey Encrypter/Decrypter` service account key as kms-enc-dec.json in the project root.
3. Save `Cloud Datastore User` service account key as datastore-user.json in the project root.


Replace the below constants in all the files with your gcp-project settings 
```go
Location        string = "global"              // replace this as per your project.
KeyRingID       string = "my-key-ring"         // replace this as per your project.
CryptoKeyName   string = "my-key"              // replace this as per your project.
ProjectName     string = "my-gcp-project-name" // replace this as per your project.
```

## running

- create the key and the key-ring by running `create/main.go`
- encrypt password and save to datastore by running `write/main.go`
- read from datastore and decrypt the password by running `read/main.go`


**Note**: 
- Make sure to enable Cloud KMS API and DataStore API for your project.
- Make sure to cleanup Google Cloud Project after you run the demo.
