package teams

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	appapi "github.com/chatwoot/chatwoot-cli/internal/chatwoot/application"
	chatwoot "github.com/chatwoot/chatwoot-cli/internal/chatwoot"
	"github.com/chatwoot/chatwoot-cli/internal/cli/cmdutil"
	"github.com/chatwoot/chatwoot-cli/internal/contract"
	"github.com/chatwoot/chatwoot-cli/internal/credentials"
	"github.com/spf13/cobra"
)

var membersListCmd = &cobra.Command{
	Use:   "list",
	Short: "List team members",
	RunE:  runMembersList,
}

var membersAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add agents to team",
	RunE:  runMembersAdd,
}

var membersUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Replace team member list",
	RunE:  runMembersUpdate,
}

var membersDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Remove agents from team",
	RunE:  runMembersDelete,
}

func init() {
	membersListCmd.Flags().Int("team-id", 0, "Team ID")
	membersListCmd.MarkFlagRequired("team-id")
	membersCmd.AddCommand(membersListCmd)

	membersAddCmd.Flags().Int("team-id", 0, "Team ID")
	membersAddCmd.MarkFlagRequired("team-id")
	membersAddCmd.Flags().String("agent-ids", "", "Comma-separated agent IDs")
	membersAddCmd.MarkFlagRequired("agent-ids")
	membersCmd.AddCommand(membersAddCmd)

	membersUpdateCmd.Flags().Int("team-id", 0, "Team ID")
	membersUpdateCmd.MarkFlagRequired("team-id")
	membersUpdateCmd.Flags().String("agent-ids", "", "Comma-separated agent IDs")
	membersUpdateCmd.MarkFlagRequired("agent-ids")
	membersCmd.AddCommand(membersUpdateCmd)

	membersDeleteCmd.Flags().Int("team-id", 0, "Team ID")
	membersDeleteCmd.MarkFlagRequired("team-id")
	membersDeleteCmd.Flags().String("agent-ids", "", "Comma-separated agent IDs")
	membersDeleteCmd.MarkFlagRequired("agent-ids")
	membersCmd.AddCommand(membersDeleteCmd)
}

func parseIntList(s string) ([]int, error) {
	parts := strings.Split(s, ",")
	ids := make([]int, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		id, err := strconv.Atoi(p)
		if err != nil {
			return nil, fmt.Errorf("invalid agent ID %q: %w", p, err)
		}
		ids = append(ids, id)
	}
	return ids, nil
}

func runMembersList(cmd *cobra.Command, args []string) error {
	rctx, err := cmdutil.ResolveContext(cmd)
	if err != nil {
		return cmdutil.WriteError(cmd, contract.ErrCodeConfig, err.Error())
	}

	tokenAuth, err := cmdutil.ResolveAuth(rctx.ProfileName, credentials.ModeApplication)
	if err != nil {
		return cmdutil.WriteError(cmd, contract.ErrCodeAuth, err.Error())
	}

	transport := chatwoot.NewClient(rctx.BaseURL, tokenAuth.Token, tokenAuth.HeaderName)
	client := appapi.NewClient(transport, rctx.AccountID)

	teamID, _ := cmd.Flags().GetInt("team-id")
	ctx := context.Background()

	agents, err := client.ListTeamMembers(ctx, teamID)
	if err != nil {
		return cmdutil.WriteError(cmd, contract.ErrCodeServer, err.Error())
	}

	resp := contract.SuccessList(agents, contract.Meta{})
	return contract.Write(cmd.OutOrStdout(), resp, cmdutil.Pretty(cmd))
}

func runMembersAdd(cmd *cobra.Command, args []string) error {
	rctx, err := cmdutil.ResolveContext(cmd)
	if err != nil {
		return cmdutil.WriteError(cmd, contract.ErrCodeConfig, err.Error())
	}

	tokenAuth, err := cmdutil.ResolveAuth(rctx.ProfileName, credentials.ModeApplication)
	if err != nil {
		return cmdutil.WriteError(cmd, contract.ErrCodeAuth, err.Error())
	}

	transport := chatwoot.NewClient(rctx.BaseURL, tokenAuth.Token, tokenAuth.HeaderName)
	client := appapi.NewClient(transport, rctx.AccountID)

	teamID, _ := cmd.Flags().GetInt("team-id")
	agentIDsStr, _ := cmd.Flags().GetString("agent-ids")
	agentIDs, err := parseIntList(agentIDsStr)
	if err != nil {
		return cmdutil.WriteError(cmd, contract.ErrCodeValidation, err.Error())
	}

	ctx := context.Background()
	agents, err := client.AddTeamMember(ctx, teamID, agentIDs)
	if err != nil {
		return cmdutil.WriteError(cmd, contract.ErrCodeServer, err.Error())
	}

	resp := contract.SuccessList(agents, contract.Meta{})
	return contract.Write(cmd.OutOrStdout(), resp, cmdutil.Pretty(cmd))
}

func runMembersUpdate(cmd *cobra.Command, args []string) error {
	rctx, err := cmdutil.ResolveContext(cmd)
	if err != nil {
		return cmdutil.WriteError(cmd, contract.ErrCodeConfig, err.Error())
	}

	tokenAuth, err := cmdutil.ResolveAuth(rctx.ProfileName, credentials.ModeApplication)
	if err != nil {
		return cmdutil.WriteError(cmd, contract.ErrCodeAuth, err.Error())
	}

	transport := chatwoot.NewClient(rctx.BaseURL, tokenAuth.Token, tokenAuth.HeaderName)
	client := appapi.NewClient(transport, rctx.AccountID)

	teamID, _ := cmd.Flags().GetInt("team-id")
	agentIDsStr, _ := cmd.Flags().GetString("agent-ids")
	agentIDs, err := parseIntList(agentIDsStr)
	if err != nil {
		return cmdutil.WriteError(cmd, contract.ErrCodeValidation, err.Error())
	}

	ctx := context.Background()
	agents, err := client.UpdateTeamMembers(ctx, teamID, agentIDs)
	if err != nil {
		return cmdutil.WriteError(cmd, contract.ErrCodeServer, err.Error())
	}

	resp := contract.SuccessList(agents, contract.Meta{})
	return contract.Write(cmd.OutOrStdout(), resp, cmdutil.Pretty(cmd))
}

func runMembersDelete(cmd *cobra.Command, args []string) error {
	rctx, err := cmdutil.ResolveContext(cmd)
	if err != nil {
		return cmdutil.WriteError(cmd, contract.ErrCodeConfig, err.Error())
	}

	tokenAuth, err := cmdutil.ResolveAuth(rctx.ProfileName, credentials.ModeApplication)
	if err != nil {
		return cmdutil.WriteError(cmd, contract.ErrCodeAuth, err.Error())
	}

	transport := chatwoot.NewClient(rctx.BaseURL, tokenAuth.Token, tokenAuth.HeaderName)
	client := appapi.NewClient(transport, rctx.AccountID)

	teamID, _ := cmd.Flags().GetInt("team-id")
	agentIDsStr, _ := cmd.Flags().GetString("agent-ids")
	agentIDs, err := parseIntList(agentIDsStr)
	if err != nil {
		return cmdutil.WriteError(cmd, contract.ErrCodeValidation, err.Error())
	}

	ctx := context.Background()
	err = client.RemoveTeamMember(ctx, teamID, agentIDs)
	if err != nil {
		return cmdutil.WriteError(cmd, contract.ErrCodeServer, err.Error())
	}

	resp := contract.Success(map[string]any{"deleted": true})
	return contract.Write(cmd.OutOrStdout(), resp, cmdutil.Pretty(cmd))
}
