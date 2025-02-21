package pr

import (
	"github.com/MakeNowJust/heredoc"
	cmdLock "github.com/moondev/cli/v2/pkg/cmd/issue/lock"
	cmdCheckout "github.com/moondev/cli/v2/pkg/cmd/pr/checkout"
	cmdChecks "github.com/moondev/cli/v2/pkg/cmd/pr/checks"
	cmdClose "github.com/moondev/cli/v2/pkg/cmd/pr/close"
	cmdComment "github.com/moondev/cli/v2/pkg/cmd/pr/comment"
	cmdCreate "github.com/moondev/cli/v2/pkg/cmd/pr/create"
	cmdDiff "github.com/moondev/cli/v2/pkg/cmd/pr/diff"
	cmdEdit "github.com/moondev/cli/v2/pkg/cmd/pr/edit"
	cmdList "github.com/moondev/cli/v2/pkg/cmd/pr/list"
	cmdMerge "github.com/moondev/cli/v2/pkg/cmd/pr/merge"
	cmdReady "github.com/moondev/cli/v2/pkg/cmd/pr/ready"
	cmdReopen "github.com/moondev/cli/v2/pkg/cmd/pr/reopen"
	cmdReview "github.com/moondev/cli/v2/pkg/cmd/pr/review"
	cmdStatus "github.com/moondev/cli/v2/pkg/cmd/pr/status"
	cmdView "github.com/moondev/cli/v2/pkg/cmd/pr/view"
	"github.com/moondev/cli/v2/pkg/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmdPR(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pr <command>",
		Short: "Manage pull requests",
		Long:  "Work with GitHub pull requests.",
		Example: heredoc.Doc(`
			$ gh pr checkout 353
			$ gh pr create --fill
			$ gh pr view --web
		`),
		Annotations: map[string]string{
			"help:arguments": heredoc.Doc(`
				A pull request can be supplied as argument in any of the following formats:
				- by number, e.g. "123";
				- by URL, e.g. "https://github.com/OWNER/REPO/pull/123"; or
				- by the name of its head branch, e.g. "patch-1" or "OWNER:patch-1".
			`),
		},
		GroupID: "core",
	}

	cmdutil.EnableRepoOverride(cmd, f)

	cmdutil.AddGroup(cmd, "General commands",
		cmdList.NewCmdList(f, nil),
		cmdCreate.NewCmdCreate(f, nil),
		cmdStatus.NewCmdStatus(f, nil),
	)

	cmdutil.AddGroup(cmd, "Targeted commands",
		cmdView.NewCmdView(f, nil),
		cmdDiff.NewCmdDiff(f, nil),
		cmdCheckout.NewCmdCheckout(f, nil),
		cmdChecks.NewCmdChecks(f, nil),
		cmdReview.NewCmdReview(f, nil),
		cmdMerge.NewCmdMerge(f, nil),
		cmdReady.NewCmdReady(f, nil),
		cmdComment.NewCmdComment(f, nil),
		cmdClose.NewCmdClose(f, nil),
		cmdReopen.NewCmdReopen(f, nil),
		cmdEdit.NewCmdEdit(f, nil),
		cmdLock.NewCmdLock(f, cmd.Name(), nil),
		cmdLock.NewCmdUnlock(f, cmd.Name(), nil),
	)

	return cmd
}
