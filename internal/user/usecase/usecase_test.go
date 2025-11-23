package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/silentmol/avito-backend-trainee/internal/testutils/mocks"
	"github.com/silentmol/avito-backend-trainee/internal/user/domain"
	"github.com/silentmol/avito-backend-trainee/internal/user/dto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserUsecase_GetUser(t *testing.T) {
	t.Parallel()

	type tc struct {
		name      string
		req       *dto.GetUserRequest
		stubUser  *domain.User
		stubErr   error
		wantErr   bool
	}

	tests := []tc{
		{
			name: "success",
			req: &dto.GetUserRequest{
				UserID: "u1",
			},
			stubUser: &domain.User{
				ID:       "u1",
				Name:     "Alice",
				TeamName: "team",
				IsActive: true,
			},
			stubErr: nil,
			wantErr: false,
		},
		{
			name: "provider_error",
			req: &dto.GetUserRequest{
				UserID: "u2",
			},
			stubUser: nil,
			stubErr:  errors.New("db error"),
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userProvider := mocks.NewMockUserProvider(ctrl)

			userProvider.EXPECT().
				GetUser(gomock.Any(), tt.req.UserID).
				Return(tt.stubUser, tt.stubErr)

			uc := &UserUsecase{userProvider: userProvider}

			resp, err := uc.GetUser(context.Background(), tt.req)
			if tt.wantErr {
				require.Error(t, err)
				require.Nil(t, resp)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, resp)
			assert.Equal(t, tt.stubUser.ID, resp.User.ID)
			assert.Equal(t, tt.stubUser.Name, resp.User.Name)
			assert.Equal(t, tt.stubUser.TeamName, resp.User.TeamName)
			assert.Equal(t, tt.stubUser.IsActive, resp.User.IsActive)
		})
	}
}

func TestUserUsecase_SetIsActive(t *testing.T) {
	t.Parallel()

	type tc struct {
		name     string
		req      *dto.SetIsActiveRequest
		stubUser *domain.User
		stubErr  error
		wantErr  bool
	}

	tests := []tc{
		{
			name: "success",
			req: &dto.SetIsActiveRequest{
				UserID:   "u1",
				IsActive: true,
			},
			stubUser: &domain.User{
				ID:       "u1",
				Name:     "Alice",
				TeamName: "team",
				IsActive: true,
			},
			stubErr: nil,
			wantErr: false,
		},
		{
			name: "provider_error",
			req: &dto.SetIsActiveRequest{
				UserID:   "u2",
				IsActive: false,
			},
			stubUser: nil,
			stubErr:  errors.New("db error"),
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userProvider := mocks.NewMockUserProvider(ctrl)

			userProvider.EXPECT().
				SetIsActive(gomock.Any(), tt.req.UserID, tt.req.IsActive).
				Return(tt.stubUser, tt.stubErr)

			uc := &UserUsecase{userProvider: userProvider}

			resp, err := uc.SetIsActive(context.Background(), tt.req)
			if tt.wantErr {
				require.Error(t, err)
				require.Nil(t, resp)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, resp)
			assert.Equal(t, tt.stubUser.ID, resp.User.ID)
			assert.Equal(t, tt.stubUser.IsActive, resp.User.IsActive)
		})
	}
}
