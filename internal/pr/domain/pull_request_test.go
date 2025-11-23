package domain

import (
	"testing"
	"time"

	"github.com/silentmol/avito-backend-trainee/internal/apperr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPullRequest_IsMerged(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		status PrStatus
		want   bool
	}{
		{
			name:   "open_is_not_merged",
			status: StatusOpen,
			want:   false,
		},
		{
			name:   "merged_is_merged",
			status: StatusMerged,
			want:   true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			pr := &PullRequest{Status: tt.status}
			assert.Equal(t, tt.want, pr.IsMerged())
		})
	}
}

func TestPullRequest_CanReassign(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		status  PrStatus
		wantErr error
	}{
		{
			name:    "open_pr_can_be_reassigned",
			status:  StatusOpen,
			wantErr: nil,
		},
		{
			name:    "merged_pr_cannot_be_reassigned",
			status:  StatusMerged,
			wantErr: apperr.ErrPRMerged,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			pr := &PullRequest{Status: tt.status}
			err := pr.CanReassign()
			if tt.wantErr == nil {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErr)
			}
		})
	}
}

func TestPullRequest_ReplaceReviewer(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		reviewers []string
		oldID     string
		newID     string
		want      []string
		wantErr   error
	}{
		{
			name:      "replace_existing_reviewer",
			reviewers: []string{"u1", "u2"},
			oldID:     "u1",
			newID:     "u3",
			want:      []string{"u3", "u2"},
			wantErr:   nil,
		},
		{
			name:      "same_reviewer_no_change",
			reviewers: []string{"u1", "u2"},
			oldID:     "u1",
			newID:     "u1",
			want:      []string{"u1", "u2"},
			wantErr:   nil,
		},
		{
			name:      "old_reviewer_not_assigned",
			reviewers: []string{"u1", "u2"},
			oldID:     "u3",
			newID:     "u4",
			want:      []string{"u1", "u2"},
			wantErr:   apperr.ErrNotAssigned,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			pr := &PullRequest{
				AssignedReviewers: append([]string(nil), tt.reviewers...),
			}

			err := pr.ReplaceReviewer(tt.oldID, tt.newID)
			if tt.wantErr == nil {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErr)
			}

			assert.Equal(t, tt.want, pr.AssignedReviewers)
		})
	}
}

func TestPullRequest_Merge(t *testing.T) {
	t.Parallel()

	pr := &PullRequest{
		Status:    StatusOpen,
		CreatedAt: time.Now().Add(-time.Hour),
	}

	require.False(t, pr.IsMerged())
	require.Nil(t, pr.MergedAt)

	pr.Merge()

	require.True(t, pr.IsMerged())
	require.NotNil(t, pr.MergedAt)

	mergedAt := pr.MergedAt

	// повторный вызов не должен ничего менять
	pr.Merge()
	assert.Same(t, mergedAt, pr.MergedAt)
	assert.True(t, pr.IsMerged())
}
