# vauth

vauth is a simple Hashicorp Vault login CLI.

vauth has been created to have a ligthweight version of the login client for [Hashicorp Vault](https://vaultproject.io/) to be used in  CICD pipelines and containarised workflow without install the full vault binary.

It supports the core vault authentication methods:

```go
    credAws "github.com/hashicorp/vault/builtin/credential/aws"
    credCert "github.com/hashicorp/vault/builtin/credential/cert"
    credGitHub "github.com/hashicorp/vault/builtin/credential/github"
    credLdap "github.com/hashicorp/vault/builtin/credential/ldap"
    credOkta "github.com/hashicorp/vault/builtin/credential/okta"
    credToken "github.com/hashicorp/vault/builtin/credential/token"
    credUserpass "github.com/hashicorp/vault/builtin/credential/userpass"
```

It's implemented using [spf13/cobra](https://github.com/spf13/cobra).

The help documentation provided by the different login methods are the native vault messages.

