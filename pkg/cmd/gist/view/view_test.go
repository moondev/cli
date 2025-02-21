package view

import (
	"bytes"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/google/shlex"
	"github.com/moondev/cli/v2/internal/config"
	"github.com/moondev/cli/v2/pkg/cmd/gist/shared"
	"github.com/moondev/cli/v2/pkg/cmdutil"
	"github.com/moondev/cli/v2/pkg/httpmock"
	"github.com/moondev/cli/v2/pkg/iostreams"
	"github.com/moondev/cli/v2/pkg/prompt"
	"github.com/stretchr/testify/assert"
)

func TestNewCmdView(t *testing.T) {
	tests := []struct {
		name  string
		cli   string
		wants ViewOptions
		tty   bool
	}{
		{
			name: "tty no arguments",
			tty:  true,
			cli:  "123",
			wants: ViewOptions{
				Raw:       false,
				Selector:  "123",
				ListFiles: false,
			},
		},
		{
			name: "nontty no arguments",
			cli:  "123",
			wants: ViewOptions{
				Raw:       true,
				Selector:  "123",
				ListFiles: false,
			},
		},
		{
			name: "filename passed",
			cli:  "-fcool.txt 123",
			tty:  true,
			wants: ViewOptions{
				Raw:       false,
				Selector:  "123",
				Filename:  "cool.txt",
				ListFiles: false,
			},
		},
		{
			name: "files passed",
			cli:  "--files 123",
			tty:  true,
			wants: ViewOptions{
				Raw:       false,
				Selector:  "123",
				ListFiles: true,
			},
		},
		{
			name: "tty no ID supplied",
			cli:  "",
			tty:  true,
			wants: ViewOptions{
				Raw:       false,
				Selector:  "",
				ListFiles: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ios, _, _, _ := iostreams.Test()
			ios.SetStdoutTTY(tt.tty)

			f := &cmdutil.Factory{
				IOStreams: ios,
			}

			argv, err := shlex.Split(tt.cli)
			assert.NoError(t, err)

			var gotOpts *ViewOptions
			cmd := NewCmdView(f, func(opts *ViewOptions) error {
				gotOpts = opts
				return nil
			})
			cmd.SetArgs(argv)
			cmd.SetIn(&bytes.Buffer{})
			cmd.SetOut(&bytes.Buffer{})
			cmd.SetErr(&bytes.Buffer{})

			_, err = cmd.ExecuteC()
			assert.NoError(t, err)

			assert.Equal(t, tt.wants.Raw, gotOpts.Raw)
			assert.Equal(t, tt.wants.Selector, gotOpts.Selector)
			assert.Equal(t, tt.wants.Filename, gotOpts.Filename)
		})
	}
}

func Test_viewRun(t *testing.T) {
	tests := []struct {
		name         string
		opts         *ViewOptions
		wantOut      string
		gist         *shared.Gist
		wantErr      bool
		mockGistList bool
	}{
		{
			name: "no such gist",
			opts: &ViewOptions{
				Selector:  "1234",
				ListFiles: false,
			},
			wantErr: true,
		},
		{
			name: "one file",
			opts: &ViewOptions{
				Selector:  "1234",
				ListFiles: false,
			},
			gist: &shared.Gist{
				Files: map[string]*shared.GistFile{
					"cicada.txt": {
						Content: "bwhiizzzbwhuiiizzzz",
						Type:    "text/plain",
					},
				},
			},
			wantOut: "bwhiizzzbwhuiiizzzz\n",
		},
		{
			name: "one file, no ID supplied",
			opts: &ViewOptions{
				Selector:  "",
				ListFiles: false,
			},
			mockGistList: true,
			gist: &shared.Gist{
				Files: map[string]*shared.GistFile{
					"cicada.txt": {
						Content: "test interactive mode",
						Type:    "text/plain",
					},
				},
			},
			wantOut: "test interactive mode\n",
		},
		{
			name: "filename selected",
			opts: &ViewOptions{
				Selector:  "1234",
				Filename:  "cicada.txt",
				ListFiles: false,
			},
			gist: &shared.Gist{
				Files: map[string]*shared.GistFile{
					"cicada.txt": {
						Content: "bwhiizzzbwhuiiizzzz",
						Type:    "text/plain",
					},
					"foo.md": {
						Content: "# foo",
						Type:    "application/markdown",
					},
				},
			},
			wantOut: "bwhiizzzbwhuiiizzzz\n",
		},
		{
			name: "filename selected, raw",
			opts: &ViewOptions{
				Selector:  "1234",
				Filename:  "cicada.txt",
				Raw:       true,
				ListFiles: false,
			},
			gist: &shared.Gist{
				Files: map[string]*shared.GistFile{
					"cicada.txt": {
						Content: "bwhiizzzbwhuiiizzzz",
						Type:    "text/plain",
					},
					"foo.md": {
						Content: "# foo",
						Type:    "application/markdown",
					},
				},
			},
			wantOut: "bwhiizzzbwhuiiizzzz\n",
		},
		{
			name: "multiple files, no description",
			opts: &ViewOptions{
				Selector:  "1234",
				ListFiles: false,
			},
			gist: &shared.Gist{
				Files: map[string]*shared.GistFile{
					"cicada.txt": {
						Content: "bwhiizzzbwhuiiizzzz",
						Type:    "text/plain",
					},
					"foo.md": {
						Content: "# foo",
						Type:    "application/markdown",
					},
				},
			},
			wantOut: "cicada.txt\n\nbwhiizzzbwhuiiizzzz\n\nfoo.md\n\n\n  # foo                                                                       \n\n",
		},
		{
			name: "multiple files, trailing newlines",
			opts: &ViewOptions{
				Selector:  "1234",
				ListFiles: false,
			},
			gist: &shared.Gist{
				Files: map[string]*shared.GistFile{
					"cicada.txt": {
						Content: "bwhiizzzbwhuiiizzzz\n",
						Type:    "text/plain",
					},
					"foo.txt": {
						Content: "bar\n",
						Type:    "text/plain",
					},
				},
			},
			wantOut: "cicada.txt\n\nbwhiizzzbwhuiiizzzz\n\nfoo.txt\n\nbar\n",
		},
		{
			name: "multiple files, description",
			opts: &ViewOptions{
				Selector:  "1234",
				ListFiles: false,
			},
			gist: &shared.Gist{
				Description: "some files",
				Files: map[string]*shared.GistFile{
					"cicada.txt": {
						Content: "bwhiizzzbwhuiiizzzz",
						Type:    "text/plain",
					},
					"foo.md": {
						Content: "- foo",
						Type:    "application/markdown",
					},
				},
			},
			wantOut: "some files\n\ncicada.txt\n\nbwhiizzzbwhuiiizzzz\n\nfoo.md\n\n\n                                                                              \n  • foo                                                                       \n\n",
		},
		{
			name: "multiple files, raw",
			opts: &ViewOptions{
				Selector:  "1234",
				Raw:       true,
				ListFiles: false,
			},
			gist: &shared.Gist{
				Description: "some files",
				Files: map[string]*shared.GistFile{
					"cicada.txt": {
						Content: "bwhiizzzbwhuiiizzzz",
						Type:    "text/plain",
					},
					"foo.md": {
						Content: "- foo",
						Type:    "application/markdown",
					},
				},
			},
			wantOut: "some files\n\ncicada.txt\n\nbwhiizzzbwhuiiizzzz\n\nfoo.md\n\n- foo\n",
		},
		{
			name: "one file, list files",
			opts: &ViewOptions{
				Selector:  "1234",
				Raw:       false,
				ListFiles: true,
			},
			gist: &shared.Gist{
				Description: "some files",
				Files: map[string]*shared.GistFile{
					"cicada.txt": {
						Content: "bwhiizzzbwhuiiizzzz",
						Type:    "text/plain",
					},
				},
			},
			wantOut: "cicada.txt\n",
		},
		{
			name: "multiple file, list files",
			opts: &ViewOptions{
				Selector:  "1234",
				Raw:       false,
				ListFiles: true,
			},
			gist: &shared.Gist{
				Description: "some files",
				Files: map[string]*shared.GistFile{
					"cicada.txt": {
						Content: "bwhiizzzbwhuiiizzzz",
						Type:    "text/plain",
					},
					"foo.md": {
						Content: "- foo",
						Type:    "application/markdown",
					},
				},
			},
			wantOut: "cicada.txt\nfoo.md\n",
		},
	}

	for _, tt := range tests {
		reg := &httpmock.Registry{}
		if tt.gist == nil {
			reg.Register(httpmock.REST("GET", "gists/1234"),
				httpmock.StatusStringResponse(404, "Not Found"))
		} else {
			reg.Register(httpmock.REST("GET", "gists/1234"),
				httpmock.JSONResponse(tt.gist))
		}

		if tt.mockGistList {
			sixHours, _ := time.ParseDuration("6h")
			sixHoursAgo := time.Now().Add(-sixHours)
			reg.Register(
				httpmock.GraphQL(`query GistList\b`),
				httpmock.StringResponse(fmt.Sprintf(
					`{ "data": { "viewer": { "gists": { "nodes": [
							{
								"name": "1234",
								"files": [{ "name": "cool.txt" }],
								"description": "",
								"updatedAt": "%s",
								"isPublic": true
							}
						] } } } }`,
					sixHoursAgo.Format(time.RFC3339),
				)),
			)

			//nolint:staticcheck // SA1019: prompt.NewAskStubber is deprecated: use PrompterMock
			as := prompt.NewAskStubber(t)
			as.StubPrompt("Select a gist").AnswerDefault()
		}

		if tt.opts == nil {
			tt.opts = &ViewOptions{}
		}

		tt.opts.HttpClient = func() (*http.Client, error) {
			return &http.Client{Transport: reg}, nil
		}

		tt.opts.Config = func() (config.Config, error) {
			return config.NewBlankConfig(), nil
		}

		ios, _, stdout, _ := iostreams.Test()
		ios.SetStdoutTTY(true)
		tt.opts.IO = ios

		t.Run(tt.name, func(t *testing.T) {
			err := viewRun(tt.opts)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)

			assert.Equal(t, tt.wantOut, stdout.String())
			reg.Verify(t)
		})
	}
}

func Test_promptGists(t *testing.T) {
	tests := []struct {
		name     string
		askStubs func(as *prompt.AskStubber)
		response string
		wantOut  string
		gist     *shared.Gist
		wantErr  bool
	}{
		{
			name: "multiple files, select first gist",
			askStubs: func(as *prompt.AskStubber) {
				as.StubPrompt("Select a gist").AnswerWith("cool.txt  about 6 hours ago")
			},
			response: `{ "data": { "viewer": { "gists": { "nodes": [
							{
								"name": "gistid1",
								"files": [{ "name": "cool.txt" }],
								"description": "",
								"updatedAt": "%[1]v",
								"isPublic": true
							},
							{
								"name": "gistid2",
								"files": [{ "name": "gistfile0.txt" }],
								"description": "",
								"updatedAt": "%[1]v",
								"isPublic": true
							}
						] } } } }`,
			wantOut: "gistid1",
		},
		{
			name: "multiple files, select second gist",
			askStubs: func(as *prompt.AskStubber) {
				as.StubPrompt("Select a gist").AnswerWith("gistfile0.txt  about 6 hours ago")
			},
			response: `{ "data": { "viewer": { "gists": { "nodes": [
							{
								"name": "gistid1",
								"files": [{ "name": "cool.txt" }],
								"description": "",
								"updatedAt": "%[1]v",
								"isPublic": true
							},
							{
								"name": "gistid2",
								"files": [{ "name": "gistfile0.txt" }],
								"description": "",
								"updatedAt": "%[1]v",
								"isPublic": true
							}
						] } } } }`,
			wantOut: "gistid2",
		},
		{
			name:     "no files",
			response: `{ "data": { "viewer": { "gists": { "nodes": [] } } } }`,
			wantOut:  "",
		},
	}

	ios, _, _, _ := iostreams.Test()

	for _, tt := range tests {
		reg := &httpmock.Registry{}

		const query = `query GistList\b`
		sixHours, _ := time.ParseDuration("6h")
		sixHoursAgo := time.Now().Add(-sixHours)
		reg.Register(
			httpmock.GraphQL(query),
			httpmock.StringResponse(fmt.Sprintf(
				tt.response,
				sixHoursAgo.Format(time.RFC3339),
			)),
		)
		client := &http.Client{Transport: reg}

		t.Run(tt.name, func(t *testing.T) {
			//nolint:staticcheck // SA1019: prompt.NewAskStubber is deprecated: use PrompterMock
			as := prompt.NewAskStubber(t)
			if tt.askStubs != nil {
				tt.askStubs(as)
			}

			gistID, err := promptGists(client, "github.com", ios.ColorScheme())
			assert.NoError(t, err)
			assert.Equal(t, tt.wantOut, gistID)
			reg.Verify(t)
		})
	}
}
