package cmd

import (
	"time"

	"github.com/imduffy15/token-cli/cli"
	"github.com/imduffy15/token-cli/client"
	"github.com/imduffy15/token-cli/config"
	"github.com/skratchdot/open-golang/open"
	"github.com/spf13/cobra"
)

func addAuthcodeTokenToContext(clientID string, token client.Token, log cli.Logger) {
	ctx := client.ClientContext{
		ClientID: clientID,
		Token:    token,
	}

	SaveContext(ctx, log)
}

func AuthcodeTokenArgumentValidation(cfg client.Config, args []string, port int) error {
	if err := EnsureActiveTarget(cfg); err != nil {
		return err
	}
	if len(args) < 1 {
		return MissingArgumentError("client_id")
	}
	return nil
}

func SaveContext(context client.ClientContext, log cli.Logger) {
	c := GetSavedConfig()
	c.AddContext(context)
	config.Write(c)
	log.Robots(context.Token.AccessToken)
}

func AuthcodeTokenCommandRun(doneRunning chan bool, clientID string, authCodeImp cli.ClientImpersonator, log cli.Logger) {
	authCodeImp.Start()
	authCodeImp.Authorize()
	token := <-authCodeImp.Done()
	SaveContext(client.ClientContext{
		ClientID: clientID,
		Token:    token,
	}, log)
	doneRunning <- true
}

func refreshContext(contextName string, cfg client.Config, log cli.Logger) error {
	context, err := cfg.GetContext(contextName)
	if err != nil {
		return err
	}
	refreshClient := client.RefreshTokenClient{
		ClientID: context.ClientID,
	}
	token, err := refreshClient.RequestToken(HTTPClient(), cfg, context.Token.RefreshToken)
	if err != nil {
		return err
	}
	context.Token = token
	SaveContext(context, log)
	return nil
}

var tokenCmd = &cobra.Command{
	Use:   "token",
	Short: "Configure and view tokens",
}

var getAuthcodeToken = &cobra.Command{
	Use:   "get CLIENT_ID --port REDIRECT_URI_PORT",
	Short: "Obtain a token for the specified CLIENT_ID",
	PreRun: func(cmd *cobra.Command, args []string) {
		cfg := GetSavedConfig()
		NotifyValidationErrors(AuthcodeTokenArgumentValidation(cfg, args, port), cmd, log)
	},
	Run: func(cmd *cobra.Command, args []string) {
		done := make(chan bool)
		cfg := GetSavedConfig()
		if exists := cfg.GetActiveTarget().ClientContextExists(args[0]); exists && !force {
			if val, err := cfg.GetContext(args[0]); err != nil {
				NotifyErrorsWithRetry(err, log)
			} else {
				if time.Unix(val.Token.ExpiresAt, 0).Sub(time.Now()) >= time.Minute*5 {
					log.Robots(val.Token.AccessToken)
				} else {
					NotifyErrorsWithRetry(refreshContext(val.ClientID, cfg, log), log)
				}
			}
		} else {
			authCodeImp := cli.NewAuthcodeClientImpersonator(HTTPClient(), cfg, args[0], scope, port, log, open.Run)
			go AuthcodeTokenCommandRun(done, args[0], authCodeImp, log)
			<-done
		}
	},
}

func init() {
	getAuthcodeToken.Flags().IntVarP(&port, "port", "p", 8080, "port on which to run local callback server")
	getAuthcodeToken.Flags().StringVarP(&scope, "scope", "s", "openid offline_access", "comma-separated scopes to request in token")
	getAuthcodeToken.Flags().BoolVarP(&force, "force", "f", false, "Forces a new token")
	tokenCmd.AddCommand(getAuthcodeToken)
	RootCmd.AddCommand(tokenCmd)
}
