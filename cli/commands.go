package cli

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/hashicorp/vault/sdk/version"
	"github.com/mitchellh/cli"
	"github.com/hashicorp/vault/command"
	/*
		The builtinplugins package is initialized here because it, in turn,
		initializes the database plugins.
		They register multiple database drivers for the "database/sql" package.
	*/
	_ "github.com/hashicorp/vault/helper/builtinplugins"

	credAliCloud "github.com/hashicorp/vault-plugin-auth-alicloud"
	credCentrify "github.com/hashicorp/vault-plugin-auth-centrify"
	credGcp "github.com/hashicorp/vault-plugin-auth-gcp/plugin"
	credOIDC "github.com/hashicorp/vault-plugin-auth-jwt"
	credAws "github.com/hashicorp/vault/builtin/credential/aws"
	credCert "github.com/hashicorp/vault/builtin/credential/cert"
	credGitHub "github.com/hashicorp/vault/builtin/credential/github"
	credLdap "github.com/hashicorp/vault/builtin/credential/ldap"
	credOkta "github.com/hashicorp/vault/builtin/credential/okta"
	credToken "github.com/hashicorp/vault/builtin/credential/token"
	credUserpass "github.com/hashicorp/vault/builtin/credential/userpass"
)

const (
	// EnvVaultCLINoColor is an env var that toggles colored UI output.
	EnvVaultCLINoColor = `VAULT_CLI_NO_COLOR`
	// EnvVaultFormat is the output format
	EnvVaultFormat = `VAULT_FORMAT`

	// flagNameAddress is the flag used in the base command to read in the
	// address of the Vault server.
	flagNameAddress = "address"
	// flagnameCACert is the flag used in the base command to read in the CA
	// cert.
	flagNameCACert = "ca-cert"
	// flagnameCAPath is the flag used in the base command to read in the CA
	// cert path.
	flagNameCAPath = "ca-path"
	//flagNameClientCert is the flag used in the base command to read in the
	//client key
	flagNameClientKey = "client-key"
	//flagNameClientCert is the flag used in the base command to read in the
	//client cert
	flagNameClientCert = "client-cert"
	// flagNameTLSSkipVerify is the flag used in the base command to read in
	// the option to ignore TLS certificate verification.
	flagNameTLSSkipVerify = "tls-skip-verify"
	// flagNameAuditNonHMACRequestKeys is the flag name used for auth/secrets enable
	flagNameAuditNonHMACRequestKeys = "audit-non-hmac-request-keys"
	// flagNameAuditNonHMACResponseKeys is the flag name used for auth/secrets enable
	flagNameAuditNonHMACResponseKeys = "audit-non-hmac-response-keys"
	// flagNameDescription is the flag name used for tuning the secret and auth mount description parameter
	flagNameDescription = "description"
	// flagListingVisibility is the flag to toggle whether to show the mount in the UI-specific listing endpoint
	flagNameListingVisibility = "listing-visibility"
	// flagNamePassthroughRequestHeaders is the flag name used to set passthrough request headers to the backend
	flagNamePassthroughRequestHeaders = "passthrough-request-headers"
	// flagNameAllowedResponseHeaders is used to set allowed response headers from a plugin
	flagNameAllowedResponseHeaders = "allowed-response-headers"
	// flagNameTokenType is the flag name used to force a specific token type
	flagNameTokenType = "token-type"
)

// Commands is the mapping of all the available commands.
var Commands map[string]cli.CommandFactory

func initCommands(ui, serverCmdUi cli.Ui, runOpts *RunOptions) {
	loginHandlers := map[string]LoginHandler{
		"alicloud": &credAliCloud.CLIHandler{},
		"aws":      &credAws.CLIHandler{},
		"centrify": &credCentrify.CLIHandler{},
		"cert":     &credCert.CLIHandler{},
		"gcp":      &credGcp.CLIHandler{},
		"github":   &credGitHub.CLIHandler{},
		"ldap":     &credLdap.CLIHandler{},
		"oidc":     &credOIDC.CLIHandler{},
		"okta":     &credOkta.CLIHandler{},
		"radius": &credUserpass.CLIHandler{
			DefaultMount: "radius",
		},
		"token": &credToken.CLIHandler{},
		"userpass": &credUserpass.CLIHandler{
			DefaultMount: "userpass",
		},
	}

	getBaseCommand := func() *command.BaseCommand {
		return &command.BaseCommand{
			UI:          ui,
			tokenHelper: runOpts.TokenHelper,
			flagAddress: runOpts.Address,
			client:      runOpts.Client,
		}
	}

	Commands = map[string]cli.CommandFactory{
		"login": func() (cli.Command, error) {
			return &command.LoginCommand{
				BaseCommand: getBaseCommand(),
				Handlers:    loginHandlers,
			}, nil
		},
		"token": func() (cli.Command, error) {
			return &command.TokenCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"token create": func() (cli.Command, error) {
			return &command.TokenCreateCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"token capabilities": func() (cli.Command, error) {
			return &command.TokenCapabilitiesCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"token lookup": func() (cli.Command, error) {
			return &command.TokenLookupCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"token renew": func() (cli.Command, error) {
			return &command.TokenRenewCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"token revoke": func() (cli.Command, error) {
			return &command.TokenRevokeCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"unwrap": func() (cli.Command, error) {
			return &command.UnwrapCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"version": func() (cli.Command, error) {
			return command.&VersionCommand{
				VersionInfo: version.GetVersion(),
				BaseCommand: getBaseCommand(),
			}, nil
		},
	}
}

// MakeShutdownCh returns a channel that can be used for shutdown
// notifications for commands. This channel will send a message for every
// SIGINT or SIGTERM received.
func MakeShutdownCh() chan struct{} {
	resultCh := make(chan struct{})

	shutdownCh := make(chan os.Signal, 4)
	signal.Notify(shutdownCh, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-shutdownCh
		close(resultCh)
	}()
	return resultCh
}

// MakeSighupCh returns a channel that can be used for SIGHUP
// reloading. This channel will send a message for every
// SIGHUP received.
func MakeSighupCh() chan struct{} {
	resultCh := make(chan struct{})

	signalCh := make(chan os.Signal, 4)
	signal.Notify(signalCh, syscall.SIGHUP)
	go func() {
		for {
			<-signalCh
			resultCh <- struct{}{}
		}
	}()
	return resultCh
}
