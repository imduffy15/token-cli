package cli

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/imduffy15/token-cli/client"
)

type ClientImpersonator interface {
	Start()
	Authorize()
	Done() chan client.Token
}

type AuthcodeClientImpersonator struct {
	httpClient         *http.Client
	config             client.Config
	ClientID           string
	Scope              string
	Port               int
	Log                Logger
	AuthCallbackServer CallbackServer
	BrowserLauncher    func(string) error
	done               chan client.Token
}

const CallbackCSS = `<style>
	@import url('https://fonts.googleapis.com/css?family=Source+Sans+Pro');
	html {
		background: #f8f8f8;
		font-family: "Source Sans Pro", sans-serif;
	}
</style>`
const AuthcodeCallbackJS = ``
const AuthcodeCallbackHTML = `<body>
	<h1>Token successfully generated</h1>
	<p>The identity provider redirected you to this page with an access token.</p>
	<p> The token has been added to the CLI's active context. You may close this window.</p>
</body>`

func NewAuthcodeClientImpersonator(
	httpClient *http.Client,
	config client.Config,
	clientID string,
	scope string,
	port int,
	log Logger,
	launcher func(string) error) AuthcodeClientImpersonator {

	impersonator := AuthcodeClientImpersonator{
		httpClient:      httpClient,
		config:          config,
		ClientID:        clientID,
		Scope:           scope,
		Port:            port,
		BrowserLauncher: launcher,
		Log:             log,
		done:            make(chan client.Token),
	}

	callbackServer := NewAuthCallbackServer(AuthcodeCallbackHTML, CallbackCSS, AuthcodeCallbackJS, log, port)
	callbackServer.SetHangupFunc(func(done chan url.Values, values url.Values) {
		token := values.Get("code")
		if token != "" {
			done <- values
		}
	})
	impersonator.AuthCallbackServer = callbackServer
	return impersonator
}

func (aci AuthcodeClientImpersonator) Start() {
	go func() {
		urlValues := make(chan url.Values)
		go aci.AuthCallbackServer.Start(urlValues)
		values := <-urlValues
		code := values.Get("code")
		tokenRequester := client.AuthorizationCodeClient{ClientID: aci.ClientID}
		aci.Log.Infof("Calling token endpoint to exchange code %v for an access token", code)
		resp, err := tokenRequester.RequestToken(aci.httpClient, aci.config, code, aci.redirectURI())
		if err != nil {
			aci.Log.Error(err.Error())
			aci.Log.Info("Retry with --verbose for more information.")
			os.Exit(1)
		}
		aci.Done() <- resp
	}()
}
func (aci AuthcodeClientImpersonator) Authorize() {
	requestValues := url.Values{}
	requestValues.Add("response_type", "code")
	requestValues.Add("client_id", aci.ClientID)
	requestValues.Add("redirect_uri", aci.redirectURI())
	requestValues.Add("scope", strings.Replace(aci.Scope, ",", " ", -1))

	authURL, err := url.Parse(aci.config.GetActiveTarget().AuthorizationEndpoint)
	if err != nil {
		aci.Log.Error("Something went wrong while building the authorization URL.")
		os.Exit(1)
	}

	authURL.RawQuery = requestValues.Encode()

	aci.Log.Info("Launching browser window to " + authURL.String() + " where the user should login and grant approvals")
	aci.BrowserLauncher(authURL.String())
}

func (aci AuthcodeClientImpersonator) Done() chan client.Token {
	return aci.done
}

func (aci AuthcodeClientImpersonator) redirectURI() string {
	return fmt.Sprintf("http://localhost:%v", aci.Port)
}
