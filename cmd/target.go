package cmd

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/imduffy15/token-cli/cli"
	"github.com/imduffy15/token-cli/client"
	"github.com/imduffy15/token-cli/config"
	"github.com/imduffy15/token-cli/utils"
	"github.com/spf13/cobra"
)

type targetStatus struct {
	AuthorizationEndpoint string
	TokenEndpoint         string
	Name                  string
	SkipSSLValidation     bool
}

type openidConfiguration struct {
	AuthorizationEndpoint string `json:"authorization_endpoint"`
	TokenEndpoint         string `json:"token_endpoint"`
}

func printTarget(target client.Target) {
	cli.NewJSONPrinter(log).Print(targetStatus{
		AuthorizationEndpoint: target.AuthorizationEndpoint,
		TokenEndpoint:         target.TokenEndpoint,
		Name:                  target.Name,
		SkipSSLValidation:     target.SkipSSLValidation,
	})
}

func UpdateTargetCmd(cfg client.Config, httpClient *http.Client, openIDConfigurationURL string, name string, log cli.Logger) error {

	res, err := httpClient.Get(openIDConfigurationURL)
	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	var openidConfiguration openidConfiguration
	json.Unmarshal(body, &openidConfiguration)

	target := client.Target{
		Name:                  name,
		SkipSSLValidation:     skipSSLValidation,
		AuthorizationEndpoint: openidConfiguration.AuthorizationEndpoint,
		TokenEndpoint:         openidConfiguration.TokenEndpoint,
	}

	cfg.AddTarget(target)

	config.Write(cfg)
	log.Info("Successfully added target " + utils.Emphasize(name))
	printTarget(target)
	return nil
}

var openIDConfigurationURL string

var targetCmd = &cobra.Command{
	Use:   "target",
	Short: "Configure and view OIDC targets",
	PreRun: func(cmd *cobra.Command, args []string) {
		cfg := GetSavedConfig()
		NotifyValidationErrors(targetCmdArgumentValidation(cfg, args), cmd, log)
	},
	Run: func(cmd *cobra.Command, args []string) {
		cfg := GetSavedConfig()
		log.Robots(cfg.GetActiveTarget().Name)
	},
}

var getCmd = &cobra.Command{
	Use:   "get TARGET_NAME",
	Short: "View the target named TARGET_NAME",
	PreRun: func(cmd *cobra.Command, args []string) {
		cfg := GetSavedConfig()
		NotifyValidationErrors(actionCmdArgumentValidation(cfg, args), cmd, log)
	},
	Run: func(cmd *cobra.Command, args []string) {
		cfg := GetSavedConfig()
		printTarget(*cfg.GetTarget(args[0]))
	},
}

var setCmd = &cobra.Command{
	Use:   "set TARGET_NAME",
	Short: "sets TARGET_NAME as active",
	PreRun: func(cmd *cobra.Command, args []string) {
		cfg := GetSavedConfig()
		NotifyValidationErrors(actionCmdArgumentValidation(cfg, args), cmd, log)
	},
	Run: func(cmd *cobra.Command, args []string) {
		cfg := GetSavedConfig()
		cfg.ActiveTargetName = args[0]
		config.Write(cfg)
		log.Info("Successfully set target to " + utils.Emphasize(args[0]))
	},
}

var deleteCmd = &cobra.Command{
	Use:   "delete TARGET_NAME",
	Short: "Delete the target named TARGET_NAME",
	PreRun: func(cmd *cobra.Command, args []string) {
		cfg := GetSavedConfig()
		NotifyValidationErrors(actionCmdArgumentValidation(cfg, args), cmd, log)
	},
	Run: func(cmd *cobra.Command, args []string) {
		cfg := GetSavedConfig()
		cfg.DeleteTarget(args[0])
		config.Write(cfg)
	},
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all targets",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := GetSavedConfig()
		cli.NewJSONPrinter(log).Print(cfg.ListTargets())
	},
}

var createCmd = &cobra.Command{
	Use:   "create TARGET_NAME",
	Short: "Creates a new target",
	PreRun: func(cmd *cobra.Command, args []string) {
		NotifyValidationErrors(createCmdArgumentValidation(args), cmd, log)
	},
	Run: func(cmd *cobra.Command, args []string) {
		cfg := GetSavedConfig()
		NotifyErrorsWithRetry(UpdateTargetCmd(cfg, HTTPClient(), openIDConfigurationURL, args[0], log), log)
	},
}

func createCmdArgumentValidation(args []string) error {
	if len(args) < 1 {
		return MissingArgumentError("target_name")
	}
	return nil
}

func actionCmdArgumentValidation(cfg client.Config, args []string) error {
	if len(args) < 1 {
		return MissingArgumentError("target_name")
	}
	if err := EnsureTargetInConfig(cfg, args[0]); err != nil {
		return err
	}
	return nil
}

func targetCmdArgumentValidation(cfg client.Config, args []string) error {
	if err := EnsureActiveTarget(cfg); err != nil {
		return err
	}
	return nil
}

func init() {

	createCmd.Flags().StringVarP(&openIDConfigurationURL, "openid-configuration-url", "t", "", "OpenID Configuration URL")
	createCmd.MarkFlagRequired("openid-configuration-url")
	createCmd.Flags().BoolVarP(&skipSSLValidation, "skip-ssl-validation", "k", false, "Disable security validation on requests to this target")

	targetCmd.AddCommand(getCmd)
	targetCmd.AddCommand(listCmd)
	targetCmd.AddCommand(createCmd)
	targetCmd.AddCommand(deleteCmd)
	targetCmd.AddCommand(setCmd)

	RootCmd.AddCommand(targetCmd)
}
