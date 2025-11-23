package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/silentmol/avito-backend-trainee/internal/apperr"
	"github.com/silentmol/avito-backend-trainee/internal/pr/domain"
	"github.com/silentmol/avito-backend-trainee/internal/pr/dto"
	teamdomain "github.com/silentmol/avito-backend-trainee/internal/team/domain"
	"github.com/silentmol/avito-backend-trainee/internal/testutils/mocks"
	userdomain "github.com/silentmol/avito-backend-trainee/internal/user/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPRUsecase_CreatePR(t *testing.T) {
	t.Parallel()

	type tc struct {
		name          string
		req           *dto.CreatePRRequest
		stubUser      *userdomain.User
		stubUserErr   error
		stubTeam      *teamdomain.Team
		stubTeamErr   error
		stubCreated   *domain.PullRequest
		stubCreateErr error
		wantErr       error
	}

	tests := []tc{
		{
			name: "success",
			req: &dto.CreatePRRequest{
				PrID:     "pr-1",
				Name:     "Add feature",
				AuthorId: "u1",
			},
			stubUser: &userdomain.User{
				ID:       "u1",
				Name:     "Alice",
				TeamName: "team",
				IsActive: true,
			},
			stubTeam: &teamdomain.Team{
				Name: "team",
				Members: []teamdomain.TeamMember{
					{ID: "u1", Name: "Alice", IsActive: true},
					{ID: "u2", Name: "Bob", IsActive: true},
				},
			},
			stubCreated: &domain.PullRequest{
				ID:       "pr-1",
				Name:     "Add feature",
				AuthorId: "u1",
				Status:   domain.StatusOpen,
			},
			wantErr: nil,
		},
		{
			name: "author_not_found",
			req: &dto.CreatePRRequest{
				PrID:     "pr-2",
				Name:     "Add search",
				AuthorId: "u2",
			},
			stubUser:    nil,
			stubUserErr: apperr.ErrNotFound,
			wantErr:     apperr.ErrNotFound,
		},
		{
			name: "team_not_found",
			req: &dto.CreatePRRequest{
				PrID:     "pr-3",
				Name:     "Add search",
				AuthorId: "u3",
			},
			stubUser: &userdomain.User{
				ID:       "u3",
				Name:     "Charlie",
				TeamName: "missing",
				IsActive: true,
			},
			stubTeam:    nil,
			stubTeamErr: apperr.ErrNotFound,
			wantErr:     apperr.ErrNotFound,
		},
		{
			name: "provider_error",
			req: &dto.CreatePRRequest{
				PrID:     "pr-4",
				Name:     "Fail in provider",
				AuthorId: "u4",
			},
			stubUser: &userdomain.User{
				ID:       "u4",
				Name:     "Dave",
				TeamName: "team",
				IsActive: true,
			},
			stubTeam: &teamdomain.Team{
				Name: "team",
				Members: []teamdomain.TeamMember{
					{ID: "u4", Name: "Dave", IsActive: true},
				},
			},
			stubCreated:   nil,
			stubCreateErr: apperr.ErrPRExists,
			wantErr:       apperr.ErrPRExists,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userReader := mocks.NewMockUserReader(ctrl)
			teamReader := mocks.NewMockTeamReader(ctrl)
			prProvider := mocks.NewMockPRProvider(ctrl)

			if tt.stubUser != nil || tt.stubUserErr != nil {
				userReader.EXPECT().
					GetUser(gomock.Any(), tt.req.AuthorId).
					DoAndReturn(func(_ context.Context, id string) (*userdomain.User, error) {
						if tt.stubUserErr != nil {
							return nil, tt.stubUserErr
						}
						assert.Equal(t, tt.req.AuthorId, id)
						return tt.stubUser, nil
					})
			}

			if tt.stubUser != nil || tt.stubTeamErr != nil {
				teamReader.EXPECT().
					GetTeam(gomock.Any(), tt.stubUser.TeamName).
					DoAndReturn(func(_ context.Context, teamName string) (*teamdomain.Team, error) {
						if tt.stubTeamErr != nil {
							return nil, tt.stubTeamErr
						}
						if tt.stubUser != nil {
							assert.Equal(t, tt.stubUser.TeamName, teamName)
						}
						return tt.stubTeam, nil
					})
			}

			if tt.stubUserErr == nil && tt.stubTeamErr == nil {
				prProvider.EXPECT().
					CreatePR(gomock.Any(), gomock.Any()).
					DoAndReturn(func(_ context.Context, pr *domain.PullRequest) (*domain.PullRequest, error) {
						if tt.stubCreated == nil {
							return nil, tt.stubCreateErr
						}
						assert.Equal(t, tt.req.PrID, pr.ID)
						assert.Equal(t, tt.req.Name, pr.Name)
						assert.Equal(t, tt.req.AuthorId, pr.AuthorId)
						return tt.stubCreated, nil
					})
			}

			uc := &PRUsecase{
				prProvider: prProvider,
				userReader: userReader,
				teamReader: teamReader,
			}

			resp, err := uc.CreatePR(context.Background(), tt.req)
			if tt.wantErr != nil {
				require.Error(t, err)
				assert.True(t, errors.Is(err, tt.wantErr))
				require.Nil(t, resp)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, resp)
			assert.Equal(t, tt.stubCreated.ID, resp.PullRequest.ID)
			assert.Equal(t, tt.stubCreated.Name, resp.PullRequest.Name)
		})
	}
}

func TestPRUsecase_MergePR(t *testing.T) {
	t.Parallel()

	type tc struct {
		name    string
		req     *dto.MergePRRequest
		stubPR  *domain.PullRequest
		stubErr error
		wantErr bool
	}

	tests := []tc{
		{
			name: "success",
			req: &dto.MergePRRequest{
				PrID: "pr-1",
			},
			stubPR: &domain.PullRequest{
				ID:     "pr-1",
				Status: domain.StatusMerged,
			},
		},
		{
			name: "provider_error",
			req: &dto.MergePRRequest{
				PrID: "pr-2",
			},
			stubErr: errors.New("db error"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			prProvider := mocks.NewMockPRProvider(ctrl)

			prProvider.EXPECT().
				MergePR(gomock.Any(), tt.req.PrID).
				Return(tt.stubPR, tt.stubErr)

			uc := &PRUsecase{prProvider: prProvider}

			resp, err := uc.MergePR(context.Background(), tt.req)
			if tt.wantErr {
				require.Error(t, err)
				require.Nil(t, resp)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, resp)
			assert.Equal(t, tt.stubPR.ID, resp.PullRequest.ID)
			assert.Equal(t, tt.stubPR.Status, resp.PullRequest.Status)
		})
	}
}

func TestPRUsecase_ReassignPR_Success(t *testing.T) {
	t.Parallel()

	req := &dto.ReassignPRRequest{
		PrID:          "pr-1",
		OldReviewerId: "u2",
	}

	stubPR := &domain.PullRequest{
		ID:                "pr-1",
		Status:            domain.StatusOpen,
		AssignedReviewers: []string{"u2", "u1"},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	prProvider := mocks.NewMockPRProvider(ctrl)
	userReader := mocks.NewMockUserReader(ctrl)
	teamReader := mocks.NewMockTeamReader(ctrl)

	prProvider.EXPECT().
		GetPR(gomock.Any(), req.PrID).
		Return(stubPR, nil)

	userReader.EXPECT().
		GetUser(gomock.Any(), req.OldReviewerId).
		Return(&userdomain.User{
			ID:       "u2",
			TeamName: "team-1",
			IsActive: true,
		}, nil)

	teamReader.EXPECT().
		GetTeam(gomock.Any(), "team-1").
		Return(&teamdomain.Team{
			Name: "team-1",
			Members: []teamdomain.TeamMember{
				{ID: "u2", Name: "R1", IsActive: true}, // old reviewer
				{ID: "u3", Name: "R2", IsActive: true}, // candidate
				{ID: "u1", Name: "Author", IsActive: true},
			},
		}, nil)

	prProvider.EXPECT().
		UpdatePR(gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, pr *domain.PullRequest) (*domain.PullRequest, error) {
			assert.Equal(t, "pr-1", pr.ID)
			assert.Equal(t, domain.StatusOpen, pr.Status)

			require.Len(t, pr.AssignedReviewers, 2)
			assert.Contains(t, pr.AssignedReviewers, "u1")
			assert.Contains(t, pr.AssignedReviewers, "u3")

			updated := *pr
			return &updated, nil
		})

	uc := &PRUsecase{
		prProvider: prProvider,
		userReader: userReader,
		teamReader: teamReader,
	}

	resp, err := uc.ReassignPR(context.Background(), req)
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, "pr-1", resp.PullRequest.ID)
	assert.Equal(t, "u3", resp.ReplacedBy)
	assert.Contains(t, resp.PullRequest.AssignedReviewers, "u3")
	assert.Contains(t, resp.PullRequest.AssignedReviewers, "u1")
	assert.NotContains(t, resp.PullRequest.AssignedReviewers, "u2")
}

func TestPRUsecase_ReassignPR_PRNotFound(t *testing.T) {
	t.Parallel()

	req := &dto.ReassignPRRequest{
		PrID:          "pr-404",
		OldReviewerId: "u2",
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	prProvider := mocks.NewMockPRProvider(ctrl)
	userReader := mocks.NewMockUserReader(ctrl)
	teamReader := mocks.NewMockTeamReader(ctrl)

	prProvider.EXPECT().
		GetPR(gomock.Any(), req.PrID).
		Return(nil, apperr.ErrNotFound)

	uc := &PRUsecase{
		prProvider: prProvider,
		userReader: userReader,
		teamReader: teamReader,
	}

	resp, err := uc.ReassignPR(context.Background(), req)
	require.Error(t, err)
	require.True(t, errors.Is(err, apperr.ErrNotFound))
	require.Nil(t, resp)
}

func TestPRUsecase_ReassignPR_PRProviderError(t *testing.T) {
	t.Parallel()

	req := &dto.ReassignPRRequest{
		PrID:          "pr-err",
		OldReviewerId: "u2",
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	prProvider := mocks.NewMockPRProvider(ctrl)
	userReader := mocks.NewMockUserReader(ctrl)
	teamReader := mocks.NewMockTeamReader(ctrl)

	prProvider.EXPECT().
		GetPR(gomock.Any(), req.PrID).
		Return(nil, errors.New("db error"))

	uc := &PRUsecase{
		prProvider: prProvider,
		userReader: userReader,
		teamReader: teamReader,
	}

	resp, err := uc.ReassignPR(context.Background(), req)
	require.Error(t, err)
	require.Nil(t, resp)
}

func TestPRUsecase_ReassignPR_PRIsMerged(t *testing.T) {
	t.Parallel()

	req := &dto.ReassignPRRequest{
		PrID:          "pr-merged",
		OldReviewerId: "u2",
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	prProvider := mocks.NewMockPRProvider(ctrl)
	userReader := mocks.NewMockUserReader(ctrl)
	teamReader := mocks.NewMockTeamReader(ctrl)

	prProvider.EXPECT().
		GetPR(gomock.Any(), req.PrID).
		Return(&domain.PullRequest{
			ID:                "pr-merged",
			Status:            domain.StatusMerged,
			AssignedReviewers: []string{"u2"},
		}, nil)

	uc := &PRUsecase{
		prProvider: prProvider,
		userReader: userReader,
		teamReader: teamReader,
	}

	resp, err := uc.ReassignPR(context.Background(), req)
	require.Error(t, err)
	require.True(t, errors.Is(err, apperr.ErrPRMerged))
	require.Nil(t, resp)
}

func TestPRUsecase_ReassignPR_OldReviewerNotFound(t *testing.T) {
	t.Parallel()

	req := &dto.ReassignPRRequest{
		PrID:          "pr-1",
		OldReviewerId: "u2",
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	prProvider := mocks.NewMockPRProvider(ctrl)
	userReader := mocks.NewMockUserReader(ctrl)
	teamReader := mocks.NewMockTeamReader(ctrl)

	prProvider.EXPECT().
		GetPR(gomock.Any(), req.PrID).
		Return(&domain.PullRequest{
			ID:                "pr-1",
			Status:            domain.StatusOpen,
			AssignedReviewers: []string{"u2"},
		}, nil)

	userReader.EXPECT().
		GetUser(gomock.Any(), req.OldReviewerId).
		Return(nil, apperr.ErrNotFound)

	uc := &PRUsecase{
		prProvider: prProvider,
		userReader: userReader,
		teamReader: teamReader,
	}

	resp, err := uc.ReassignPR(context.Background(), req)
	require.Error(t, err)
	require.True(t, errors.Is(err, apperr.ErrNotFound))
	require.Nil(t, resp)
}

func TestPRUsecase_ReassignPR_TeamNotFound(t *testing.T) {
	t.Parallel()

	req := &dto.ReassignPRRequest{
		PrID:          "pr-1",
		OldReviewerId: "u2",
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	prProvider := mocks.NewMockPRProvider(ctrl)
	userReader := mocks.NewMockUserReader(ctrl)
	teamReader := mocks.NewMockTeamReader(ctrl)

	prProvider.EXPECT().
		GetPR(gomock.Any(), req.PrID).
		Return(&domain.PullRequest{
			ID:                "pr-1",
			Status:            domain.StatusOpen,
			AssignedReviewers: []string{"u2"},
		}, nil)

	userReader.EXPECT().
		GetUser(gomock.Any(), req.OldReviewerId).
		Return(&userdomain.User{
			ID:       "u2",
			TeamName: "team-missing",
			IsActive: true,
		}, nil)

	teamReader.EXPECT().
		GetTeam(gomock.Any(), "team-missing").
		Return(nil, apperr.ErrNotFound)

	uc := &PRUsecase{
		prProvider: prProvider,
		userReader: userReader,
		teamReader: teamReader,
	}

	resp, err := uc.ReassignPR(context.Background(), req)
	require.Error(t, err)
	require.True(t, errors.Is(err, apperr.ErrNotFound))
	require.Nil(t, resp)
}

func TestPRUsecase_ReassignPR_NoCandidates(t *testing.T) {
	t.Parallel()

	req := &dto.ReassignPRRequest{
		PrID:          "pr-1",
		OldReviewerId: "u2",
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	prProvider := mocks.NewMockPRProvider(ctrl)
	userReader := mocks.NewMockUserReader(ctrl)
	teamReader := mocks.NewMockTeamReader(ctrl)

	prProvider.EXPECT().
		GetPR(gomock.Any(), req.PrID).
		Return(&domain.PullRequest{
			ID:                "pr-1",
			Status:            domain.StatusOpen,
			AssignedReviewers: []string{"u2", "u1"},
		}, nil)

	userReader.EXPECT().
		GetUser(gomock.Any(), req.OldReviewerId).
		Return(&userdomain.User{
			ID:       "u2",
			TeamName: "team-1",
			IsActive: true,
		}, nil)

	teamReader.EXPECT().
		GetTeam(gomock.Any(), "team-1").
		Return(&teamdomain.Team{
			Name: "team-1",
			Members: []teamdomain.TeamMember{
				{ID: "u2", Name: "R1", IsActive: true}, // old reviewer
				{ID: "u1", Name: "Author", IsActive: true},
			},
		}, nil)

	uc := &PRUsecase{
		prProvider: prProvider,
		userReader: userReader,
		teamReader: teamReader,
	}

	resp, err := uc.ReassignPR(context.Background(), req)
	require.Error(t, err)
	require.True(t, errors.Is(err, apperr.ErrNoCandidate))
	require.Nil(t, resp)
}

func TestPRUsecase_ReassignPR_OldReviewerNotAssigned(t *testing.T) {
	t.Parallel()

	req := &dto.ReassignPRRequest{
		PrID:          "pr-1",
		OldReviewerId: "u2",
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	prProvider := mocks.NewMockPRProvider(ctrl)
	userReader := mocks.NewMockUserReader(ctrl)
	teamReader := mocks.NewMockTeamReader(ctrl)

	prProvider.EXPECT().
		GetPR(gomock.Any(), req.PrID).
		Return(&domain.PullRequest{
			ID:                "pr-1",
			Status:            domain.StatusOpen,
			AssignedReviewers: []string{"u1"}, // old reviewer not in list
		}, nil)

	userReader.EXPECT().
		GetUser(gomock.Any(), req.OldReviewerId).
		Return(&userdomain.User{
			ID:       "u2",
			TeamName: "team-1",
			IsActive: true,
		}, nil)

	teamReader.EXPECT().
		GetTeam(gomock.Any(), "team-1").
		Return(&teamdomain.Team{
			Name: "team-1",
			Members: []teamdomain.TeamMember{
				{ID: "u2", Name: "R1", IsActive: true}, // candidate
				{ID: "u1", Name: "Author", IsActive: true},
			},
		}, nil)

	uc := &PRUsecase{
		prProvider: prProvider,
		userReader: userReader,
		teamReader: teamReader,
	}

	resp, err := uc.ReassignPR(context.Background(), req)
	require.Error(t, err)
	require.True(t, errors.Is(err, apperr.ErrNotAssigned))
	require.Nil(t, resp)
}

func TestPRUsecase_ReassignPR_UpdateError(t *testing.T) {
	t.Parallel()

	req := &dto.ReassignPRRequest{
		PrID:          "pr-1",
		OldReviewerId: "u2",
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	prProvider := mocks.NewMockPRProvider(ctrl)
	userReader := mocks.NewMockUserReader(ctrl)
	teamReader := mocks.NewMockTeamReader(ctrl)

	prProvider.EXPECT().
		GetPR(gomock.Any(), req.PrID).
		Return(&domain.PullRequest{
			ID:                "pr-1",
			Status:            domain.StatusOpen,
			AssignedReviewers: []string{"u2", "u1"},
		}, nil)

	userReader.EXPECT().
		GetUser(gomock.Any(), req.OldReviewerId).
		Return(&userdomain.User{
			ID:       "u2",
			TeamName: "team-1",
			IsActive: true,
		}, nil)

	teamReader.EXPECT().
		GetTeam(gomock.Any(), "team-1").
		Return(&teamdomain.Team{
			Name: "team-1",
			Members: []teamdomain.TeamMember{
				{ID: "u2", Name: "R1", IsActive: true},
				{ID: "u3", Name: "R2", IsActive: true},
				{ID: "u1", Name: "Author", IsActive: true},
			},
		}, nil)

	prProvider.EXPECT().
		UpdatePR(gomock.Any(), gomock.Any()).
		Return(nil, errors.New("update error"))

	uc := &PRUsecase{
		prProvider: prProvider,
		userReader: userReader,
		teamReader: teamReader,
	}

	resp, err := uc.ReassignPR(context.Background(), req)
	require.Error(t, err)
	require.Nil(t, resp)
}

func TestPRUsecase_GetReview(t *testing.T) {
	t.Parallel()

	type tc struct {
		name    string
		req     *dto.GetReviewRequest
		stubPRs []domain.PullRequest
		stubErr error
		wantErr bool
	}

	tests := []tc{
		{
			name: "success_no_prs",
			req: &dto.GetReviewRequest{
				UserId: "u1",
			},
			stubPRs: nil,
			stubErr: nil,
		},
		{
			name: "success_with_prs",
			req: &dto.GetReviewRequest{
				UserId: "u2",
			},
			stubPRs: []domain.PullRequest{
				{
					ID:       "pr-1",
					Name:     "Add search",
					AuthorId: "u1",
					Status:   domain.StatusOpen,
				},
			},
			stubErr: nil,
		},
		{
			name: "provider_error",
			req: &dto.GetReviewRequest{
				UserId: "u3",
			},
			stubErr: errors.New("db error"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			prProvider := mocks.NewMockPRProvider(ctrl)

			prProvider.EXPECT().
				GetReview(gomock.Any(), tt.req.UserId).
				DoAndReturn(func(_ context.Context, userID string) (*[]domain.PullRequest, error) {
					assert.Equal(t, tt.req.UserId, userID)
					if tt.stubErr != nil {
						return nil, tt.stubErr
					}
					if tt.stubPRs == nil {
						empty := make([]domain.PullRequest, 0)
						return &empty, nil
					}
					cp := make([]domain.PullRequest, len(tt.stubPRs))
					copy(cp, tt.stubPRs)
					return &cp, nil
				})

			uc := &PRUsecase{prProvider: prProvider}

			resp, err := uc.GetReview(context.Background(), tt.req)
			if tt.wantErr {
				require.Error(t, err)
				require.Nil(t, resp)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, resp)
			assert.Equal(t, tt.req.UserId, resp.UserId)
			assert.Len(t, resp.PullRequests, len(tt.stubPRs))
			for i, pr := range tt.stubPRs {
				assert.Equal(t, pr.ID, resp.PullRequests[i].ID)
				assert.Equal(t, pr.Name, resp.PullRequests[i].Name)
				assert.Equal(t, pr.AuthorId, resp.PullRequests[i].AuthorId)
				assert.Equal(t, pr.Status, resp.PullRequests[i].Status)
			}
		})
	}
}
