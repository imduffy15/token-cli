[![GitHub license](https://img.shields.io/github/license/imduffy15/token-cli.svg)](https://github.com/imduffy15/token-cli/blob/master/LICENSE)
![GitHub release](https://img.shields.io/github/release/imduffy15/token-cli.svg)

# TokenCLI

_tokenCLI_ is a command line utility for generating tokens from a OpenID identity provider, such as [Keycloak](https://www.keycloak.org/).

_tokenCLI_ uses the [Authorization Code Grant Flow](https://developer.okta.com/authentication-guide/implementing-authentication/auth-code), as such a refresh token is generated and used to automatically renew the access token without browser interaction.

![token-cli - The OpenID token generator](images/walkthrough.gif)

## Installation

### OSX

Install:

```
brew install imduffy15/tap/token-cli
```

Upgrade:

```
brew upgrade token-cli
```


## Alternative Installs (tar.gz, RPM, deb, snap)
Check out the [releases](https://github.com/imduffy15/token-cli/releases) section on Github for alternative binaries.

## Contribute
[Fork token-cli](https://github.com/imduffy15/token-cli) and build a custom version. We welcome any useful pull requests.

## Usage

Create a new target called example-realm:

```
$ token-cli target create example-realm -t http://localhost:8080/auth/realms/example-realm/.well-known/openid-configuration
```

Set example-realm as the active target:

```
$ token-cli target set example-realm
```

Get a token for the client "service-template" with redirection port 9090

```
$ token-cli token get service-template -p 9090
```

## Examples

### Google

Add Google as a target and set it as the active target

```
$ token-cli target create google --openid-configuration-url https://accounts.google.com/.well-known/openid-configuration
$ token-cli target set google
```

Generate a token for client-id `571394967398-j6vs98u325la013f0ho6hehosdi2h2eb.apps.googleusercontent.com` with scope openid

```
$ token-cli token get 571394967398-j6vs98u325la013f0ho6hehosdi2h2eb.apps.googleusercontent.com --scope openid
```

You can register Google clients at https://console.cloud.google.com/apis/credentials

### Microsoft

Add Microsoft as a target and set it as the active target

```
$ token-cli target create microsoft --openid-configuration-url https://login.microsoftonline.com/common/v2.0/.well-known/openid-configuration
$ token-cli target set microsoft
```

Generate a token for client-id `b9951982-9e22-4bb8-8632-436f47b030f4`

```
$ token-cli token get b9951982-9e22-4bb8-8632-436f47b030f4 --client_secret 'uNmN9EdyV0HpcT@Sel.v[4LH2mRs@/bH'
```

You can register Microsoft clients at https://portal.azure.com/#blade/Microsoft_AAD_IAM/ActiveDirectoryMenuBlade/RegisteredApps , this target is for `(Any Azure AD directory - Multitenant) and personal Microsoft accounts (e.g. Skype, Xbox)` authorization only.


### Azure AD

Add Azure AD as a target, where `960a630a-dab4-4fd9-a048-88eddede726d` is my tenant id and set it as the active target

```
$ token-cli target create azure --openid-configuration-url https://login.microsoftonline.com/960a630a-dab4-4fd9-a048-88eddede726d/v2.0/.well-known/openid-configuration
$ token-cli target set azure
```

Generate a token for client-id `90a49166-df3b-46a9-bb20-155f4055ef83`

```
$ token-cli token get 90a49166-df3b-46a9-bb20-155f4055ef83 --client_secret 'pGw]yxRb:fww-fk?X2uskpfSPlHXV559'
```

You can register Azure AD clients at https://portal.azure.com/#blade/Microsoft_AAD_IAM/ActiveDirectoryMenuBlade/RegisteredApps , this target is for `(Default Directory only - Single tenant)` authorization only.

### Okta

Add Okta as a target and set it as the active target

```
$ token-cli target create okta --openid-configuration-url https://dev-351390-admin.oktapreview.com/.well-known/openid-configuration
$ token-cli target set okta
```

Generate a token for client-id `0oaodg6in82hvS2uV0h7`

```
$ token-cli token get 0oaodg6in82hvS2uV0h7 --client_secret '8XojWRVjegSMh8hzmZZ-NIq-9ur6fauRDPk-Rv-k'
```

You can register OKTA clients in your personal dashboard at https://TENANT-ID-admin.oktapreview.com/admin/apps/active , where tenant-id is your tenant id.

### Instagram

Add Instagram as as target and set it as the active target

```
$ token-cli target create instagram --token-url https://api.instagram.com/oauth/access_token --authorization-url https://api.instagram.com/oauth/authorize
$ token-cli target set instagram
```

Generate a token for client-id `c3b3514c9a614b53b6f393b7dc3f7459`

```
$ token-cli token get c3b3514c9a614b53b6f393b7dc3f7459 --client_secret 4dfdef8221284c2480c8d71cea00d0b2  --scope basic
```

You can register Instagram clients at https://www.instagram.com/developer/clients/manage/

### Strava

```
$ token-cli target create strava --token-url https://www.strava.com/oauth/token --authorizaion-url http://www.strava.com/oauth/authorize
$ token-cli target set strava
```

Generate a token for client-id `40638`

```
token-cli token get 40638 --client_secret e36b1089bfe26c8010cd10eabe419c96493c412b --scope read
```

You can register Strava clients at https://www.strava.com/settings/api

## Help

```bash
$ token-cli --help
Token Command Line Interface, version

Usage:
  token-cli [command]

Available Commands:
  help        Help about any command
  target      Configure and view OIDC targets
  token       Configure and view tokens

Flags:
  -h, --help      help for token-cli
  -v, --verbose   See additional info on HTTP requests

Use "token-cli [command] --help" for more information about a command.
```

```bash
$ token-cli target --help
Configure and view OIDC targets

Usage:
  token-cli target [flags]
  token-cli target [command]

Available Commands:
  create      Creates a new target
  delete      Delete the target named TARGET_NAME
  get         View the target named TARGET_NAME
  list        List all targets
  set         sets TARGET_NAME as active

Flags:
  -h, --help                  help for target
  -k, --skip-ssl-validation   Disable security validation on requests to this target

Global Flags:
  -v, --verbose   See additional info on HTTP requests

Use "token-cli target [command] --help" for more information about a command.
```

```bash
$ token-cli token --help
Configure and view tokens

Usage:
  token-cli token [command]

Available Commands:
  get         Obtain a token for the specified CLIENT_ID

Flags:
  -h, --help   help for token

Global Flags:
  -v, --verbose   See additional info on HTTP requests

Use "token-cli token [command] --help" for more information about a command.
```

## License 

Apache License 2.0  
