package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTeam_ActiveMembersExcept(t *testing.T) {
	t.Parallel()

	type fields struct {
		members []TeamMember
	}

	type args struct {
		excluded []string
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []TeamMember
		teamNil bool
	}{
		{
			name:    "nil_team_returns_nil",
			teamNil: true,
		},
		{
			name: "empty_members_returns_empty",
			fields: fields{
				members: []TeamMember{},
			},
			args: args{
				excluded: nil,
			},
			want: []TeamMember{},
		},
		{
			name: "filters_inactive_and_excluded",
			fields: fields{
				members: []TeamMember{
					{ID: "u1", Name: "Alice", IsActive: true},
					{ID: "u2", Name: "Bob", IsActive: false},
					{ID: "u3", Name: "Charlie", IsActive: true},
				},
			},
			args: args{
				excluded: []string{"u3"},
			},
			want: []TeamMember{
				{ID: "u1", Name: "Alice", IsActive: true},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var team *Team
			if !tt.teamNil {
				team = &Team{
					Members: tt.fields.members,
				}
			}

			got := team.ActiveMembersExcept(tt.args.excluded...)

			if tt.teamNil {
				assert.Nil(t, got)
				return
			}

			require.Len(t, got, len(tt.want))
			assert.ElementsMatch(t, tt.want, got)
		})
	}
}
