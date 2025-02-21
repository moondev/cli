package shared

import (
	"testing"

	"github.com/moondev/cli/v2/api"
	"github.com/moondev/cli/v2/internal/ghrepo"
	"github.com/moondev/cli/v2/pkg/iostreams"
	"github.com/moondev/cli/v2/pkg/prompt"
	"github.com/stretchr/testify/assert"
)

type metadataFetcher struct {
	metadataResult *api.RepoMetadataResult
}

func (mf *metadataFetcher) RepoMetadataFetch(input api.RepoMetadataInput) (*api.RepoMetadataResult, error) {
	return mf.metadataResult, nil
}

func TestMetadataSurvey_selectAll(t *testing.T) {
	ios, _, stdout, stderr := iostreams.Test()

	repo := ghrepo.New("OWNER", "REPO")

	fetcher := &metadataFetcher{
		metadataResult: &api.RepoMetadataResult{
			AssignableUsers: []api.RepoAssignee{
				{Login: "hubot"},
				{Login: "monalisa"},
			},
			Labels: []api.RepoLabel{
				{Name: "help wanted"},
				{Name: "good first issue"},
			},
			Projects: []api.RepoProject{
				{Name: "Huge Refactoring"},
				{Name: "The road to 1.0"},
			},
			Milestones: []api.RepoMilestone{
				{Title: "1.2 patch release"},
			},
		},
	}

	//nolint:staticcheck // SA1019: prompt.InitAskStubber is deprecated: use NewAskStubber
	as, restoreAsk := prompt.InitAskStubber()
	defer restoreAsk()

	//nolint:staticcheck // SA1019: as.Stub is deprecated: use StubPrompt
	as.Stub([]*prompt.QuestionStub{
		{
			Name:  "metadata",
			Value: []string{"Labels", "Projects", "Assignees", "Reviewers", "Milestone"},
		},
	})
	//nolint:staticcheck // SA1019: as.Stub is deprecated: use StubPrompt
	as.Stub([]*prompt.QuestionStub{
		{
			Name:  "reviewers",
			Value: []string{"monalisa"},
		},
		{
			Name:  "assignees",
			Value: []string{"hubot"},
		},
		{
			Name:  "labels",
			Value: []string{"good first issue"},
		},
		{
			Name:  "projects",
			Value: []string{"The road to 1.0"},
		},
		{
			Name:  "milestone",
			Value: "(none)",
		},
	})

	state := &IssueMetadataState{
		Assignees: []string{"hubot"},
	}
	err := MetadataSurvey(ios, repo, fetcher, state)
	assert.NoError(t, err)

	assert.Equal(t, "", stdout.String())
	assert.Equal(t, "", stderr.String())

	assert.Equal(t, []string{"hubot"}, state.Assignees)
	assert.Equal(t, []string{"monalisa"}, state.Reviewers)
	assert.Equal(t, []string{"good first issue"}, state.Labels)
	assert.Equal(t, []string{"The road to 1.0"}, state.Projects)
	assert.Equal(t, []string{}, state.Milestones)
}

func TestMetadataSurvey_keepExisting(t *testing.T) {
	ios, _, stdout, stderr := iostreams.Test()

	repo := ghrepo.New("OWNER", "REPO")

	fetcher := &metadataFetcher{
		metadataResult: &api.RepoMetadataResult{
			Labels: []api.RepoLabel{
				{Name: "help wanted"},
				{Name: "good first issue"},
			},
			Projects: []api.RepoProject{
				{Name: "Huge Refactoring"},
				{Name: "The road to 1.0"},
			},
		},
	}

	//nolint:staticcheck // SA1019: prompt.InitAskStubber is deprecated: use NewAskStubber
	as, restoreAsk := prompt.InitAskStubber()
	defer restoreAsk()

	//nolint:staticcheck // SA1019: as.Stub is deprecated: use StubPrompt
	as.Stub([]*prompt.QuestionStub{
		{
			Name:  "metadata",
			Value: []string{"Labels", "Projects"},
		},
	})
	//nolint:staticcheck // SA1019: as.Stub is deprecated: use StubPrompt
	as.Stub([]*prompt.QuestionStub{
		{
			Name:  "labels",
			Value: []string{"good first issue"},
		},
		{
			Name:  "projects",
			Value: []string{"The road to 1.0"},
		},
	})

	state := &IssueMetadataState{
		Assignees: []string{"hubot"},
	}
	err := MetadataSurvey(ios, repo, fetcher, state)
	assert.NoError(t, err)

	assert.Equal(t, "", stdout.String())
	assert.Equal(t, "", stderr.String())

	assert.Equal(t, []string{"hubot"}, state.Assignees)
	assert.Equal(t, []string{"good first issue"}, state.Labels)
	assert.Equal(t, []string{"The road to 1.0"}, state.Projects)
}
