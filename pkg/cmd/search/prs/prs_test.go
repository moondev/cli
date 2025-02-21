package prs

import (
	"bytes"
	"testing"

	"github.com/google/shlex"
	"github.com/moondev/cli/v2/pkg/cmd/search/shared"
	"github.com/moondev/cli/v2/pkg/cmdutil"
	"github.com/moondev/cli/v2/pkg/iostreams"
	"github.com/moondev/cli/v2/pkg/search"
	"github.com/stretchr/testify/assert"
)

func TestNewCmdPrs(t *testing.T) {
	var trueBool = true
	tests := []struct {
		name    string
		input   string
		output  shared.IssuesOptions
		wantErr bool
		errMsg  string
	}{
		{
			name:    "no arguments",
			input:   "",
			wantErr: true,
			errMsg:  "specify search keywords or flags",
		},
		{
			name:  "keyword arguments",
			input: "some search terms",
			output: shared.IssuesOptions{
				Query: search.Query{
					Keywords:   []string{"some", "search", "terms"},
					Kind:       "issues",
					Limit:      30,
					Qualifiers: search.Qualifiers{Type: "pr"},
				},
			},
		},
		{
			name:  "web flag",
			input: "--web",
			output: shared.IssuesOptions{
				Query: search.Query{
					Keywords:   []string{},
					Kind:       "issues",
					Limit:      30,
					Qualifiers: search.Qualifiers{Type: "pr"},
				},
				WebMode: true,
			},
		},
		{
			name:  "limit flag",
			input: "--limit 10",
			output: shared.IssuesOptions{
				Query: search.Query{
					Keywords:   []string{},
					Kind:       "issues",
					Limit:      10,
					Qualifiers: search.Qualifiers{Type: "pr"},
				},
			},
		},
		{
			name:    "invalid limit flag",
			input:   "--limit 1001",
			wantErr: true,
			errMsg:  "`--limit` must be between 1 and 1000",
		},
		{
			name:  "order flag",
			input: "--order asc",
			output: shared.IssuesOptions{
				Query: search.Query{
					Keywords:   []string{},
					Kind:       "issues",
					Limit:      30,
					Order:      "asc",
					Qualifiers: search.Qualifiers{Type: "pr"},
				},
			},
		},
		{
			name:    "invalid order flag",
			input:   "--order invalid",
			wantErr: true,
			errMsg:  "invalid argument \"invalid\" for \"--order\" flag: valid values are {asc|desc}",
		},
		{
			name:  "app flag",
			input: "--app dependabot",
			output: shared.IssuesOptions{
				Query: search.Query{
					Keywords:   []string{},
					Kind:       "issues",
					Limit:      30,
					Qualifiers: search.Qualifiers{Type: "pr", Author: "app/dependabot"},
				},
			},
		},
		{
			name:    "invalid author and app flags",
			input:   "--author test --app dependabot",
			wantErr: true,
			errMsg:  "specify only `--author` or `--app`",
		},
		{
			name: "qualifier flags",
			input: `
      --archived
      --assignee=assignee
      --author=author
      --closed=closed
      --commenter=commenter
      --created=created
      --match=title,body
      --language=language
      --locked
      --merged
      --no-milestone
      --updated=updated
      --visibility=public
      `,
			output: shared.IssuesOptions{
				Query: search.Query{
					Keywords: []string{},
					Kind:     "issues",
					Limit:    30,
					Qualifiers: search.Qualifiers{
						Archived:  &trueBool,
						Assignee:  "assignee",
						Author:    "author",
						Closed:    "closed",
						Commenter: "commenter",
						Created:   "created",
						In:        []string{"title", "body"},
						Is:        []string{"public", "locked", "merged"},
						Language:  "language",
						No:        []string{"milestone"},
						Type:      "pr",
						Updated:   "updated",
					},
				},
			},
		},
		{
			name:  "review requested flag with user",
			input: "--review-requested user",
			output: shared.IssuesOptions{
				Query: search.Query{
					Keywords:   []string{},
					Kind:       "issues",
					Limit:      30,
					Qualifiers: search.Qualifiers{Type: "pr", ReviewRequested: "user"},
				},
			},
		},
		{
			name:  "review requested flag with team",
			input: "--review-requested org/team",
			output: shared.IssuesOptions{
				Query: search.Query{
					Keywords:   []string{},
					Kind:       "issues",
					Limit:      30,
					Qualifiers: search.Qualifiers{Type: "pr", TeamReviewRequested: "org/team"},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ios, _, _, _ := iostreams.Test()
			f := &cmdutil.Factory{
				IOStreams: ios,
			}
			argv, err := shlex.Split(tt.input)
			assert.NoError(t, err)
			var gotOpts *shared.IssuesOptions
			cmd := NewCmdPrs(f, func(opts *shared.IssuesOptions) error {
				gotOpts = opts
				return nil
			})
			cmd.SetArgs(argv)
			cmd.SetIn(&bytes.Buffer{})
			cmd.SetOut(&bytes.Buffer{})
			cmd.SetErr(&bytes.Buffer{})

			_, err = cmd.ExecuteC()
			if tt.wantErr {
				assert.EqualError(t, err, tt.errMsg)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.output.Query, gotOpts.Query)
			assert.Equal(t, tt.output.WebMode, gotOpts.WebMode)
		})
	}
}
