package command

import (
	"bytes"
	"context"
	"fmt"
	"github.com/hashicorp/vault/api"
	vt "github.com/mauromedda/vauth/command/token"
	testcontainers "github.com/testcontainers/testcontainers-go"
	"strings"
	"testing"
)

func TestLogin(t *testing.T) {
	checkLogins := func(t *testing.T, got, want string) {
		t.Helper()
		if !strings.Contains(got, want) {
			t.Errorf("got %q want %q", got, want)
		}
	}
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

	loginTest := []struct {
		name   string
		method string
		params map[string]string
		want   string
	}{
		{name: "userpass successful login", method: "userpass", params: map[string]string{"username": "test", "password": "test"}, want: "Success! You are now authenticated."},
		{name: "userpass wrong login", method: "userpass", params: map[string]string{"username": "testWrong", "password": "test"}, want: "invalid username or password"},
		{name: "token successful login", method: "token", params: map[string]string{"token": token}, want: "Success! You are now authenticated."},
		{name: "token wrong login", method: "token", params: map[string]string{"token": "tokenWrong"}, want: "permission denied"},
	}
	tokenHelper := vt.InternalTokenHelper{}
	for _, tt := range loginTest {
		t.Run(tt.name, func(t *testing.T) {
			// Erase the token in the local client
			defer tokenHelper.Erase()
			got := &bytes.Buffer{}
			if err := Login(client, tt.method, tt.params, got); err != nil {
				checkLogins(t, err.Error(), tt.want)
			} else {
				checkLogins(t, got.String(), tt.want)
			}

		})
	}
}
