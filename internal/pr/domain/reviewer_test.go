package domain

import (
	"testing"

	"github.com/silentmol/avito-backend-trainee/internal/apperr"
	teamdomain "github.com/silentmol/avito-backend-trainee/internal/team/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSelectReviewersForTeam(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		team     *teamdomain.Team
		authorID string
		wantLen  int
	}{
		{
			name:     "nil_team_returns_nil",
			team:     nil,
			authorID: "u1",
			wantLen:  -1, // special: expect nil
		},
		{
			name: "only_author_in_team",
			team: &teamdomain.Team{
				Members: []teamdomain.TeamMember{
					{ID: "u1", IsActive: true},
				},
			},
			authorID: "u1",
			wantLen:  0,
		},
		{
			name: "one_active_non_author",
			team: &teamdomain.Team{
				Members: []teamdomain.TeamMember{
					{ID: "u1", IsActive: true},
					{ID: "u2", IsActive: true},
				},
			},
			authorID: "u1",
			wantLen:  1,
		},
		{
			name: "two_active_non_author",
			team: &teamdomain.Team{
				Members: []teamdomain.TeamMember{
					{ID: "u1", IsActive: true},
					{ID: "u2", IsActive: true},
					{ID: "u3", IsActive: true},
				},
			},
			authorID: "u1",
			wantLen:  2,
		},
		{
			name: "more_than_two_candidates",
			team: &teamdomain.Team{
				Members: []teamdomain.TeamMember{
					{ID: "u1", IsActive: true},
					{ID: "u2", IsActive: true},
					{ID: "u3", IsActive: true},
					{ID: "u4", IsActive: true},
				},
			},
			authorID: "u1",
			wantLen:  2,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := SelectReviewersForTeam(tt.team, tt.authorID)

			if tt.wantLen == -1 {
				assert.Nil(t, got)
				return
			}

			require.Len(t, got, tt.wantLen)

			// все ревьюеры должны быть не автором и активными членами команды
			if tt.team != nil {
				active := tt.team.ActiveMembersExcept(tt.authorID)
				activeSet := make(map[string]struct{}, len(active))
				for _, m := range active {
					activeSet[m.ID] = struct{}{}
				}

				for _, id := range got {
					assert.NotEqual(t, tt.authorID, id)
					_, ok := activeSet[id]
					assert.True(t, ok, "reviewer must be active team member")
				}
			}
		})
	}
}

func TestReassignReviewer(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		pr          *PullRequest
		team        *teamdomain.Team
		oldReviewer string
		wantErr     error
		wantNewID   string
	}{
		{
			name:        "nil_pr_or_team_return_no_candidate",
			pr:          nil,
			team:        nil,
			oldReviewer: "u2",
			wantErr:     apperr.ErrNoCandidate,
		},
		{
			name: "merged_pr_cannot_be_reassigned",
			pr: &PullRequest{
				Status:            StatusMerged,
				AssignedReviewers: []string{"u2"},
			},
			team: &teamdomain.Team{
				Members: []teamdomain.TeamMember{
					{ID: "u2", IsActive: true},
					{ID: "u3", IsActive: true},
				},
			},
			oldReviewer: "u2",
			wantErr:     apperr.ErrPRMerged,
		},
		{
			name: "old_reviewer_not_assigned",
			pr: &PullRequest{
				Status:            StatusOpen,
				AssignedReviewers: []string{"u1"},
			},
			team: &teamdomain.Team{
				Members: []teamdomain.TeamMember{
					{ID: "u2", IsActive: true},
				},
			},
			oldReviewer: "u2",
			wantErr:     apperr.ErrNotAssigned,
		},
		{
			name: "no_candidates_in_team",
			pr: &PullRequest{
				Status:            StatusOpen,
				AssignedReviewers: []string{"u2", "u3"},
			},
			team: &teamdomain.Team{
				Members: []teamdomain.TeamMember{
					{ID: "u2", IsActive: true},
					{ID: "u3", IsActive: true},
				},
			},
			oldReviewer: "u2",
			wantErr:     apperr.ErrNoCandidate,
		},
		{
			name: "success_with_single_candidate",
			pr: &PullRequest{
				Status:            StatusOpen,
				AssignedReviewers: []string{"u2", "u1"},
			},
			team: &teamdomain.Team{
				Members: []teamdomain.TeamMember{
					{ID: "u2", IsActive: true}, // old reviewer
					{ID: "u3", IsActive: true}, // candidate
					{ID: "u1", IsActive: true}, // already reviewer
				},
			},
			oldReviewer: "u2",
			wantErr:     nil,
			wantNewID:   "u3",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			newID, err := ReassignReviewer(tt.pr, tt.team, tt.oldReviewer)

			if tt.wantErr != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErr)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.wantNewID, newID)

			// старый ревьюер должен быть заменён
			require.NotNil(t, tt.pr)
			require.Contains(t, tt.pr.AssignedReviewers, tt.wantNewID)
			assert.NotContains(t, tt.pr.AssignedReviewers, tt.oldReviewer)
		})
	}
}
