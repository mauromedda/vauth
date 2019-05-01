package command

import (
	"fmt"
	"github.com/hashicorp/vault/api"
	credAws "github.com/hashicorp/vault/builtin/credential/aws"
	credCert "github.com/hashicorp/vault/builtin/credential/cert"
	credGitHub "github.com/hashicorp/vault/builtin/credential/github"
	credLdap "github.com/hashicorp/vault/builtin/credential/ldap"
	credOkta "github.com/hashicorp/vault/builtin/credential/okta"
	credToken "github.com/hashicorp/vault/builtin/credential/token"
	credUserpass "github.com/hashicorp/vault/builtin/credential/userpass"
	vt "github.com/mauromedda/vauth/command/token"
	"github.com/spf13/cobra"
	"io"
	"os"
)

// LoginHandler is the interface that any auth handlers must implement to enable
// auth via the CLI.
type LoginHandler interface {
	Auth(*api.Client, map[string]string) (*api.Secret, error)
	Help() string
}

// LoginHandlers is an k:v datatype with authentication method type and
// the related vault Handler
var LoginHandlers = map[string]LoginHandler{
	"aws":    &credAws.CLIHandler{},
	"cert":   &credCert.CLIHandler{},
	"github": &credGitHub.CLIHandler{},
	"ldap":   &credLdap.CLIHandler{},
	"okta":   &credOkta.CLIHandler{},
	"radius": &credUserpass.CLIHandler{
		DefaultMount: "radius",
	},
	"token": &credToken.CLIHandler{},
	"userpass": &credUserpass.CLIHandler{
		DefaultMount: "userpass",
	},
}

// Login function returns a vault token or an error
func Login(cmd *cobra.Command, method string, loginConfig map[string]string) error {
	config := api.DefaultConfig()
	if err := config.ReadEnvironment(); err != nil {
		return fmt.Errorf("%s failed to read environment", err)
	}
	client, err := api.NewClient(config)
	clih, ok := LoginHandlers[method]
	if !ok {
		return fmt.Errorf("%s method not supported", method)
	}
	if method == "userpass" || method == "ldap" {
		username, ok := loginConfig["username"]
		if !ok {
			username = UsernameFromEnv()
			if username == "" {
				return fmt.Errorf("'username' not supplied and neither 'LOGNAME' nor 'USER' env vars set")
			}
		}
		password, ok := loginConfig["password"]
		if !ok {
			password = PasswordFromEnv()
		}
		authConfig := map[string]string{
			"username": username,
			"password": password,
			"method":   method,
		}
		for k, v := range authConfig {
			loginConfig[k] = v
		}
	}
	sec, err := clih.Auth(client, loginConfig)
	if err != nil {
		cmd.SilenceUsage = true
		return fmt.Errorf("%s\n%s", err, clih.Help())
	}

	tokenID, err := sec.TokenID()
	if err != nil {
		return fmt.Errorf("No token available")
	}

	// Store the token in the local client
	tokenHelper := vt.InternalTokenHelper{}
	tokenHelper.PopulateTokenPath()

	if err := tokenHelper.Store(tokenID); err != nil {
		fmt.Printf("Error storing token: %s", err)
		return fmt.Errorf(
			"Authentication was successful, but the token was not persisted. The " +
				"resulting token is shown below for your records.\n")
	}
	fmt.Printf("Token ID: %s\n", tokenID)
	return nil
}

func init() {
	rootCmd.AddCommand(loginCmd)
}

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Authenticate clients against Vault",
	Long: `This subcommand authenticate the client to Vault using the provided method.
The login sub-command and the related methods accept the same parameter of the mainstream Hashicorp Vault CLI.

Valid methods are: aws, ldap, token, userpass, radius, github, okta and cert.
`,
	SilenceUsage: false,
	RunE: func(cmd *cobra.Command, args []string) error {
		if !cmd.Flags().Changed("method") {
			return fmt.Errorf("No authentication method provided")
		}
		method, err := cmd.Flags().GetString("method")
		if err != nil {
			return err
		}
		// Pull the Hashicorp Vault fake stdin if needed
		stdin := (io.Reader)(os.Stdin)
		authConfig, err := parseArgsDataString(stdin, cmd.Flags().Args())
		if err != nil {
			return err
		}
		return Login(cmd, method, authConfig)

	},
}

// Method identify the login method used by the LoginCommand function
var Method string

func init() {
	loginCmd.Flags().StringVarP(&Method, "method", "m", "", "Authentication method for Vault")
}
