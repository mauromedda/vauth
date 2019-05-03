package command

import (
	"bytes"
	"context"
	"fmt"
	"github.com/hashicorp/vault/api"
	testcontainers "github.com/testcontainers/testcontainers-go"
	"strings"
	"testing"
)

func TestLoginSuccess(t *testing.T) {
	token := "s6gjRs4pYBO4pyDGyp73e8Zmt"
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "vault",
		ExposedPorts: []string{"8200/tcp"},
		Cmd:          "server -dev",
		Env: map[string]string{
			"VAULT_DEV_ROOT_TOKEN_ID":  token,
			"VAULT_DEV_LISTEN_ADDRESS": "0.0.0.0:8200",
		},
	}
	vaultC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Error(err)
	}
	// At the end of the test remove the container
	defer vaultC.Terminate(ctx)
	// Retrieve the container IP
	ip, err := vaultC.Host(ctx)
	if err != nil {
		t.Error(err)
	}
	// Retrieve the port mapped to port 8200
	port, err := vaultC.MappedPort(ctx, "8200")
	if err != nil {
		t.Error(err)
	}

	address := fmt.Sprintf("http://%s:%s", ip, port.Port())
	client, _ := NewClient(nil)
	client.SetAddress(address)
	client.SetToken(token)
	authOpts := &api.EnableAuthOptions{
		Type: "userpass",
		Config: api.AuthConfigInput{
			DefaultLeaseTTL: "600",
			MaxLeaseTTL:     "800",
		},
	}
	// Enable the userpass authentication
	if err := client.Sys().EnableAuthWithOptions("userpass", authOpts); err != nil {
		t.Fatal(err)
	}
	if _, err := client.Logical().Write("auth/userpass/users/test", map[string]interface{}{
		"password": "test",
		"policies": "default",
	}); err != nil {
		t.Fatal(err)
	}
	got := &bytes.Buffer{}
	if err := Login(client, "userpass", map[string]string{"username": "test", "password": "test"}, got); err != nil {
		t.Fatal(err)
	}
	want := "Success! You are now authenticated."
	if !strings.Contains(got.String(), want) {
		t.Errorf("expected %q to be %q", got, want)
	}

	if err := Login(client, "userpass", map[string]string{"username": "testWrong", "password": "test"}, got); err == nil {
		t.Fatal("expected login failed but success")
	}
}
