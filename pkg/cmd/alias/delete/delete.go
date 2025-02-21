package delete

import (
	"fmt"

	"github.com/moondev/cli/v2/internal/config"
	"github.com/moondev/cli/v2/pkg/cmdutil"
	"github.com/moondev/cli/v2/pkg/iostreams"
	"github.com/spf13/cobra"
)

type DeleteOptions struct {
	Config func() (config.Config, error)
	IO     *iostreams.IOStreams

	Name string
}

func NewCmdDelete(f *cmdutil.Factory, runF func(*DeleteOptions) error) *cobra.Command {
	opts := &DeleteOptions{
		IO:     f.IOStreams,
		Config: f.Config,
	}

	cmd := &cobra.Command{
		Use:   "delete <alias>",
		Short: "Delete an alias",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Name = args[0]

			if runF != nil {
				return runF(opts)
			}
			return deleteRun(opts)
		},
	}

	return cmd
}

func deleteRun(opts *DeleteOptions) error {
	cfg, err := opts.Config()
	if err != nil {
		return err
	}

	aliasCfg := cfg.Aliases()

	expansion, err := aliasCfg.Get(opts.Name)
	if err != nil {
		return fmt.Errorf("no such alias %s", opts.Name)

	}

	err = aliasCfg.Delete(opts.Name)
	if err != nil {
		return fmt.Errorf("failed to delete alias %s: %w", opts.Name, err)
	}

	err = cfg.Write()
	if err != nil {
		return err
	}

	if opts.IO.IsStdoutTTY() {
		cs := opts.IO.ColorScheme()
		fmt.Fprintf(opts.IO.ErrOut, "%s Deleted alias %s; was %s\n", cs.SuccessIconWithColor(cs.Red), opts.Name, expansion)
	}

	return nil
}
