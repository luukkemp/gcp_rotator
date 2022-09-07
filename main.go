package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	iam "google.golang.org/api/iam/v1"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("Usage: %v filename.json", os.Args[0])
		return
	}
	keyFile := os.Args[1]
	content, err := ioutil.ReadFile(keyFile)
	if err != nil {
		fmt.Println(err)
	}

	var g GcpJson

	err = json.Unmarshal([]byte(content), &g)

	g.RenewKey()

	newJson, err := json.MarshalIndent(g, "", "  ")
	if err != nil {
		panic(fmt.Sprintf("failed to marshal key to json: %v", err))
	}

	err = ioutil.WriteFile(keyFile, newJson, 0644)
	if err != nil {
		panic(fmt.Sprintf("failed to write key to file: %v", err))
	}

	fmt.Printf("Generated new key for %v, new key ID: %v\n", g.ClientEmail, g.PrivateKeyId)
	fmt.Printf("Saved to file: %v\n", keyFile)

}

type GcpJson struct {
	Type                    string `json:"type"`
	ProjectId               string `json:"project_id"`
	PrivateKeyId            string `json:"private_key_id"`
	PrivateKey              string `json:"private_key"`
	ClientEmail             string `json:"client_email"`
	ClientId                string `json:"client_id"`
	AuthUri                 string `json:"auth_uri"`
	TokenUri                string `json:"token_uri"`
	AuthProviderX509CertUrl string `json:"auth_provider_x509_cert_url"`
	ClientX509CertUrl       string `json:"client_x509_cert_url"`
}

func (g *GcpJson) RenewKey() {
	ctx := context.Background()

	service, err := iam.NewService(ctx)
	if err != nil {
		fmt.Printf("iam.NewService: %v\n", err)
	}

	resource := "projects/-/serviceAccounts/" + g.ClientEmail
	request := &iam.CreateServiceAccountKeyRequest{}
	newKey, err := service.Projects.ServiceAccounts.Keys.Create(resource, request).Do()
	if err != nil {
		panic(fmt.Sprintf("Projects.ServiceAccounts.Keys.Create: %v\n", err))
	}

	oldKeyName := fmt.Sprintf("projects/%v/serviceAccounts/%v/keys/%v", g.ProjectId, g.ClientEmail, g.PrivateKeyId)
	_, err = service.Projects.ServiceAccounts.Keys.Delete(oldKeyName).Do()
	if err != nil {
		panic(fmt.Sprintf("Projects.ServiceAccounts.Keys.Delete: %v\n", err))
	}

	decodedBytes, _ := base64.StdEncoding.DecodeString(newKey.PrivateKeyData)
	json.Unmarshal(decodedBytes, &g)
}
