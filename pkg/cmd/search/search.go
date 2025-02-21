package search

import (
	"github.com/moondev/cli/v2/pkg/cmdutil"
	"github.com/spf13/cobra"

	searchCommitsCmd "github.com/moondev/cli/v2/pkg/cmd/search/commits"
	searchIssuesCmd "github.com/moondev/cli/v2/pkg/cmd/search/issues"
	searchPrsCmd "github.com/moondev/cli/v2/pkg/cmd/search/prs"
	searchReposCmd "github.com/moondev/cli/v2/pkg/cmd/search/repos"
)

func NewCmdSearch(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "search <command>",
		Short: "Search for repositories, issues, and pull requests",
		Long:  "Search across all of GitHub.",
	}

	cmd.AddCommand(searchCommitsCmd.NewCmdCommits(f, nil))
	cmd.AddCommand(searchIssuesCmd.NewCmdIssues(f, nil))
	cmd.AddCommand(searchPrsCmd.NewCmdPrs(f, nil))
	cmd.AddCommand(searchReposCmd.NewCmdRepos(f, nil))

	return cmd
}
