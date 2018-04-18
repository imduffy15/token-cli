package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/imduffy15/token-cli/cli"
	"github.com/imduffy15/token-cli/client"
	"github.com/spf13/cobra"
)

const TargetDoesNotExist = "the target you specified does not exist"
const NoActiveTarget = "there is no active target set, please set one"

func MissingArgumentError(argName string) error {
	return MissingArgumentWithExplanationError(argName, "")
}

func MissingArgumentWithExplanationError(argName string, explanation string) error {
	return fmt.Errorf("Missing argument `%v` must be specified. %v", argName, explanation)
}

func EnsureActiveTarget(cfg client.Config) error {
	if cfg.ActiveTargetName == "" {
		return errors.New(NoActiveTarget)
	}
	if !cfg.TargetExists(cfg.ActiveTargetName) {
		return errors.New(TargetDoesNotExist)
	}
	return nil
}

func EnsureTargetInConfig(cfg client.Config, targetName string) error {
	if !cfg.TargetExists(targetName) {
		return errors.New(TargetDoesNotExist)
	}
	return nil
}

func NotifyValidationErrors(err error, cmd *cobra.Command, log cli.Logger) {
	if err != nil {
		log.Error(err.Error())
		cmd.Usage()
		os.Exit(1)
	}
}

func NotifyErrorsWithRetry(err error, log cli.Logger) {
	if err != nil {
		log.Error(err.Error())
		VerboseRetryMsg(GetSavedConfig())
		os.Exit(1)
	}
}

func VerboseRetryMsg(c client.Config) {
	if !c.Verbose {
		log.Info("Retry with --verbose for more information.")
	}
}
