package cli

import (
	"fmt"
	"log/slog"
	"os"

	cliapp "github.com/chatwoot/chatwoot-cli/internal/cli/application"
	cliauth "github.com/chatwoot/chatwoot-cli/internal/cli/auth"
	"github.com/spf13/cobra"
)

var (
	prettyFlag    bool
	verboseFlag   bool
	profileFlag   string
	baseURLFlag   string
	accountIDFlag int
)

var rootCmd = &cobra.Command{
	Use:           "chatwoot",
	Short:         "Chatwoot CLI — machine-friendly command-line interface for Chatwoot",
	SilenceUsage:  true,
	SilenceErrors: true,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Apply env var defaults only when the flag was not explicitly set.
		if !cmd.Flags().Changed("pretty") {
			env := os.Getenv("CHATWOOT_PRETTY")
			if env == "1" || env == "true" {
				prettyFlag = true
			}
		}
		if !cmd.Flags().Changed("verbose") {
			env := os.Getenv("CHATWOOT_VERBOSE")
			if env == "1" || env == "true" {
				verboseFlag = true
			}
		}
		if verboseFlag {
			slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug})))
		}
	},
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&prettyFlag, "pretty", false, "Indent JSON output with 2 spaces")
	rootCmd.PersistentFlags().BoolVar(&verboseFlag, "verbose", false, "Enable diagnostic logging on stderr")
	rootCmd.PersistentFlags().StringVar(&profileFlag, "profile", "", "Select named profile")
	rootCmd.PersistentFlags().StringVar(&baseURLFlag, "base-url", "", "Override base URL")
	rootCmd.PersistentFlags().IntVar(&accountIDFlag, "account-id", 0, "Override account ID")
	rootCmd.AddCommand(cliauth.Cmd)
	rootCmd.AddCommand(cliapp.Cmd)
}

// Execute runs the root command and returns an exit code.
// 0 = success, 1 = runtime/API error, 2 = usage error.
func Execute() int {
	err := rootCmd.Execute()
	if err == nil {
		return 0
	}

	// Detect usage errors (unknown command, bad flags, missing args).
	// Cobra sets SilenceUsage to prevent printing usage on RunE errors,
	// but for unknown-command/flag errors Cobra itself returns an error
	// before RunE executes. We detect this by checking if the error
	// message starts with "unknown command" or "unknown flag".
	errMsg := err.Error()
	if isUsageError(errMsg) {
		fmt.Fprintln(os.Stderr, "Error:", errMsg)
		return 2
	}

	return 1
}

func isUsageError(msg string) bool {
	prefixes := []string{
		"unknown command",
		"unknown flag",
		"unknown shorthand flag",
		"required flag",
		"accepts ",
	}
	for _, p := range prefixes {
		if len(msg) >= len(p) && msg[:len(p)] == p {
			return true
		}
	}
	return false
}

// Pretty returns whether pretty-printing is enabled.
func Pretty() bool {
	return prettyFlag
}
