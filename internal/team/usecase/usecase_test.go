package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/silentmol/avito-backend-trainee/internal/team/domain"
	"github.com/silentmol/avito-backend-trainee/internal/team/dto"
	"github.com/silentmol/avito-backend-trainee/internal/testutils/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTeamUsecase_CreateTeam(t *testing.T) {
	t.Parallel()

	type tc struct {
		name     string
		req      *dto.AddTeamRequest
		stubTeam *domain.Team
		stubErr  error
		wantErr  bool
	}

	tests := []tc{
		{
			name: "success",
			req: &dto.AddTeamRequest{
				Team: domain.Team{
					Name: "test",
					Members: []domain.TeamMember{
						{ID: "u1", Name: "Alice", IsActive: true},
					},
				},
			},
			stubTeam: &domain.Team{
				Name: "test",
				Members: []domain.TeamMember{
					{ID: "u1", Name: "Alice", IsActive: true},
				},
			},
			stubErr: nil,
			wantErr: false,
		},
		{
			name: "provider_error",
			req: &dto.AddTeamRequest{
				Team: domain.Team{Name: "bad"},
			},
			stubTeam: nil,
			stubErr:  errors.New("db error"),
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			teamProvider := mocks.NewMockTeamProvider(ctrl)

			teamProvider.EXPECT().
				CreateTeam(gomock.Any(), gomock.Any()).
				DoAndReturn(func(_ context.Context, team *domain.Team) (*domain.Team, error) {
					assert.Equal(t, tt.req.Name, team.Name)
					assert.Equal(t, tt.req.Members, team.Members)
					return tt.stubTeam, tt.stubErr
				})

			uc := &TeamUsecase{teamProvider: teamProvider}

			resp, err := uc.CreateTeam(context.Background(), tt.req)
			if tt.wantErr {
				require.Error(t, err)
				require.Nil(t, resp)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, resp)
			assert.Equal(t, tt.stubTeam.Name, resp.Team.Name)
			assert.Equal(t, tt.stubTeam.Members, resp.Team.Members)
		})
	}
}

func TestTeamUsecase_GetTeam(t *testing.T) {
	t.Parallel()

	type tc struct {
		name     string
		req      *dto.GetTeamRequest
		stubTeam *domain.Team
		stubErr  error
		wantErr  bool
	}

	tests := []tc{
		{
			name: "success",
			req: &dto.GetTeamRequest{
				TeamName: "team",
			},
			stubTeam: &domain.Team{
				Name: "team",
				Members: []domain.TeamMember{
					{ID: "u1", Name: "Alice", IsActive: true},
				},
			},
			stubErr: nil,
			wantErr: false,
		},
		{
			name: "provider_error",
			req: &dto.GetTeamRequest{
				TeamName: "missing",
			},
			stubTeam: nil,
			stubErr:  errors.New("db error"),
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			teamProvider := mocks.NewMockTeamProvider(ctrl)

			teamProvider.EXPECT().
				GetTeam(gomock.Any(), tt.req.TeamName).
				Return(tt.stubTeam, tt.stubErr)

			uc := &TeamUsecase{teamProvider: teamProvider}

			resp, err := uc.GetTeam(context.Background(), tt.req)
			if tt.wantErr {
				require.Error(t, err)
				require.Nil(t, resp)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, resp)
			assert.Equal(t, tt.stubTeam.Name, resp.Team.Name)
			assert.Equal(t, tt.stubTeam.Members, resp.Team.Members)
		})
	}
}
