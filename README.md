# vauth

[![Build Status](https://travis-ci.org/mauromedda/vauth.svg?branch=master)](https://travis-ci.org/mauromedda/vauth)

vauth is a simple Hashicorp Vault login CLI wrapper.

> Inspired by the Hashicorp vault project.

> Thanks to the [Vault](https://github.com/hashicorp/vault) community and to [hashicorp](https://hashicorp.com/)

vauth has been created to have a ligthweight version of the login client for [Hashicorp Vault](https://vaultproject.io/) to be used in  CICD pipelines and containarised workflow without install the full vault binary.

It supports the following core Hashicorp Vault [authentication methods](https://www.vaultproject.io/docs/auth/):

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

## Installation

### From source

```
go get -u github.com/mauromedda/vauth
```

### From binary

1. Go to the [releases page](https://github.com/mauromedda/vauth/releases)
2. Download the binary for your system and put into your `$PATH`

## Usage

### Sample workthrough

```bash
# Installation
$ curl -sSLfo vauth https://github.com/mauromedda/vauth/releases/download/v0.1.1/vauth_linux_amd64
$ chmod 0755 vauth
$ sudo mv vauth /usr/local/bin/

# Log against vault
$ export VAULT_ADDR=https://<your_vault_server_address>:8200
$ vauth login -m userpass username=test password=test
Success! You are now authenticated. The token information displayed
below is already stored in the token helper. You do NOT need to run
"vauth login" again. Future Vault requests will automatically use this token.
TokenID: s.oXsX8GqsYxyvXmtkjpT8fLhU
```

### From Docker image

```bash
# Pull the docker image

$ docker pull mauromedda/vauth

$ docker run -e VAULT_ADDR=https://<your_vault_server_address>:8200 -it mauromedda/vauth login -m userpass username=test password=test
Success! You are now authenticated. The token information displayed
below is already stored in the token helper. You do NOT need to run
"vauth login" again. Future Vault requests will automatically use this token.
TokenID: s.oXsX8GqsYxyvXmtkjpT8fLhU
```

## Authors

Currently maintained by [Mauro Medda](https://github.com/mauromedda).