package list

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/MakeNowJust/heredoc"
	"github.com/moondev/cli/v2/internal/config"
	"github.com/moondev/cli/v2/pkg/httpmock"
	"github.com/moondev/cli/v2/pkg/iostreams"
)

func TestListRun(t *testing.T) {
	tests := []struct {
		name       string
		opts       ListOptions
		isTTY      bool
		wantStdout string
		wantStderr string
		wantErr    bool
	}{
		{
			name: "list tty",
			opts: ListOptions{
				HTTPClient: func() (*http.Client, error) {
					createdAt := time.Now().Add(time.Duration(-24) * time.Hour)
					reg := &httpmock.Registry{}
					reg.Register(
						httpmock.REST("GET", "user/keys"),
						httpmock.StringResponse(fmt.Sprintf(`[
							{
								"id": 1234,
								"key": "ssh-rsa AAAABbBB123",
								"title": "Mac",
								"created_at": "%[1]s"
							},
							{
								"id": 5678,
								"key": "ssh-rsa EEEEEEEK247",
								"title": "hubot@Windows",
								"created_at": "%[1]s"
							}
						]`, createdAt.Format(time.RFC3339))),
					)
					return &http.Client{Transport: reg}, nil
				},
			},
			isTTY: true,
			wantStdout: heredoc.Doc(`
				TITLE          ID    KEY                  ADDED
				Mac            1234  ssh-rsa AAAABbBB123  1d
				hubot@Windows  5678  ssh-rsa EEEEEEEK247  1d
			`),
			wantStderr: "",
		},
		{
			name: "list non-tty",
			opts: ListOptions{
				HTTPClient: func() (*http.Client, error) {
					createdAt, _ := time.Parse(time.RFC3339, "2020-08-31T15:44:24+02:00")
					reg := &httpmock.Registry{}
					reg.Register(
						httpmock.REST("GET", "user/keys"),
						httpmock.StringResponse(fmt.Sprintf(`[
							{
								"id": 1234,
								"key": "ssh-rsa AAAABbBB123",
								"title": "Mac",
								"created_at": "%[1]s"
							},
							{
								"id": 5678,
								"key": "ssh-rsa EEEEEEEK247",
								"title": "hubot@Windows",
								"created_at": "%[1]s"
							}
						]`, createdAt.Format(time.RFC3339))),
					)
					return &http.Client{Transport: reg}, nil
				},
			},
			isTTY: false,
			wantStdout: heredoc.Doc(`
				Mac	ssh-rsa AAAABbBB123	2020-08-31T15:44:24+02:00	1234
				hubot@Windows	ssh-rsa EEEEEEEK247	2020-08-31T15:44:24+02:00	5678
			`),
			wantStderr: "",
		},
		{
			name: "no keys tty",
			opts: ListOptions{
				HTTPClient: func() (*http.Client, error) {
					reg := &httpmock.Registry{}
					reg.Register(
						httpmock.REST("GET", "user/keys"),
						httpmock.StringResponse(`[]`),
					)
					return &http.Client{Transport: reg}, nil
				},
			},
			wantStdout: "",
			wantStderr: "",
			wantErr:    true,
			isTTY:      true,
		},
		{
			name: "no keys non-tty",
			opts: ListOptions{
				HTTPClient: func() (*http.Client, error) {
					reg := &httpmock.Registry{}
					reg.Register(
						httpmock.REST("GET", "user/keys"),
						httpmock.StringResponse(`[]`),
					)
					return &http.Client{Transport: reg}, nil
				},
			},
			wantStdout: "",
			wantStderr: "",
			wantErr:    true,
			isTTY:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ios, _, stdout, stderr := iostreams.Test()
			ios.SetStdinTTY(tt.isTTY)
			ios.SetStdoutTTY(tt.isTTY)
			ios.SetStderrTTY(tt.isTTY)

			opts := tt.opts
			opts.IO = ios
			opts.Config = func() (config.Config, error) { return config.NewBlankConfig(), nil }

			err := listRun(&opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("listRun() return error: %v", err)
				return
			}

			if stdout.String() != tt.wantStdout {
				t.Errorf("wants stdout %q, got %q", tt.wantStdout, stdout.String())
			}
			if stderr.String() != tt.wantStderr {
				t.Errorf("wants stderr %q, got %q", tt.wantStderr, stderr.String())
			}
		})
	}
}
