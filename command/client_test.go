package command

import (
	"context"
	"fmt"
	"github.com/hashicorp/vault/api"
	testcontainers "github.com/testcontainers/testcontainers-go"
	"os"
	"strings"
	"testing"
)

func init() {
	// Ensure our special envvars are not present
	os.Setenv("VAULT_ADDR", "")
	os.Setenv("VAULT_TOKEN", "")
}

func TestDefaultConfig_envvar(t *testing.T) {
	os.Setenv("VAULT_ADDR", "https://vault.mycompany.com")
	defer os.Setenv("VAULT_ADDR", "")

	config := api.DefaultConfig()
	if config.Address != "https://vault.mycompany.com" {
		t.Fatalf("bad: %s", config.Address)
	}

	os.Setenv("VAULT_TOKEN", "testing")
	defer os.Setenv("VAULT_TOKEN", "")

	client, err := NewClient(config)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if token := client.Token(); token != "testing" {
		t.Fatalf("bad: %s", token)
	}
}

func TestClientNilConfig(t *testing.T) {
	client, err := NewClient(nil)
	if err != nil {
		t.Fatal(err)
	}
	if client == nil {
		t.Fatal("expected a non-nil client")
	}
}

func TestClientSetAddress(t *testing.T) {
	client, err := NewClient(nil)
	if err != nil {
		t.Fatal(err)
	}
	if err := client.SetAddress("http://172.168.2.1:8300"); err != nil {
		t.Fatal(err)
	}
	if client.Address() != "http://172.168.2.1:8300" {
		t.Fatalf("bad: expected: 'http://172.168.2.1:8300' actual: %q", client.Address())
	}
}

func TestClientToken(t *testing.T) {
	tokenValue := "foo"

	var config *api.Config

	client, err := NewClient(config)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	client.SetToken(tokenValue)

	// Verify the token is set
	if v := client.Token(); v != tokenValue {
		t.Fatalf("bad: %s", v)
	}

	client.ClearToken()

	if v := client.Token(); v != "" {
		t.Fatalf("bad: %s", v)
	}
}

func TestClientBadToken(t *testing.T) {
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

	client.SetToken(token)
	_, err = client.RawRequest(client.NewRequest("PUT", "/"))
	if err != nil {
		t.Fatal(err)
	}

	client.SetToken("foo\u007f")
	_, err = client.RawRequest(client.NewRequest("PUT", "/"))
	if err == nil || !strings.Contains(err.Error(), "printable") {
		t.Fatalf("expected error due to bad token")
	}
}
