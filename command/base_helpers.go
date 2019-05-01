package command

import (
	"fmt"
	"github.com/hashicorp/vault/api"
	kvbuilder "github.com/hashicorp/vault/helper/kv-builder"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"io"
	"os"
)

// NewClient return a new vault client and an error
func NewClient(config *api.Config) (*api.Client, error) {
	if config == nil {
		config = api.DefaultConfig()
	}
	if err := config.ReadEnvironment(); err != nil {
		return nil, fmt.Errorf("%s failed to read environment", err)
	}
	client, err := api.NewClient(config)
	return client, err
}

// parseArgsData parses the given args in the format key=value into a map of
// the provided arguments. The given reader can also supply key=value pairs.
func parseArgsData(stdin io.Reader, args []string) (map[string]interface{}, error) {
	builder := &kvbuilder.Builder{Stdin: stdin}
	if err := builder.Add(args...); err != nil {
		return nil, err
	}

	return builder.Map(), nil
}

// parseArgsDataString parses the args data and returns the values as strings.
// If the values cannot be represented as strings, an error is returned.
func parseArgsDataString(stdin io.Reader, args []string) (map[string]string, error) {
	raw, err := parseArgsData(stdin, args)
	if err != nil {
		return nil, err
	}

	var result map[string]string
	if err := mapstructure.WeakDecode(raw, &result); err != nil {
		return nil, errors.Wrap(err, "failed to convert values to strings")
	}
	if result == nil {
		result = make(map[string]string)
	}
	return result, nil
}

// parseArgsDataStringLists parses the args data and returns the values as
// string lists. If the values cannot be represented as strings, an error is
// returned.
func parseArgsDataStringLists(stdin io.Reader, args []string) (map[string][]string, error) {
	raw, err := parseArgsData(stdin, args)
	if err != nil {
		return nil, err
	}

	var result map[string][]string
	if err := mapstructure.WeakDecode(raw, &result); err != nil {
		return nil, errors.Wrap(err, "failed to convert values to strings")
	}
	return result, nil
}

// UsernameFromEnv returns the username stored into LOGNAME or USER if defined or
// the empty string
func UsernameFromEnv() string {
	if logname := os.Getenv("LOGNAME"); logname != "" {
		return logname
	}
	if user := os.Getenv("USER"); user != "" {
		return user
	}
	return ""
}

// PasswordFromEnv returns the password stored into PASSWORD or
// the empty string
func PasswordFromEnv() string {
	if password := os.Getenv("PASSWORD"); password != "" {
		return password
	}
	return ""
}
