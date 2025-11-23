package domain

type Team struct {
	Name    string       `json:"team_name" validate:"required"`
	Members []TeamMember `json:"members" validate:"required"`
}

type TeamMember struct {
	ID       string `json:"user_id" validate:"required"`
	Name     string `json:"username" validate:"required"`
	IsActive bool   `json:"is_active"`
}

//  возвращает активных участников, исключая переданные ID.
func (t *Team) ActiveMembersExcept(excludedIDs ...string) []TeamMember {
	if t == nil {
		return nil
	}

	exclude := make(map[string]struct{}, len(excludedIDs))
	for _, id := range excludedIDs {
		exclude[id] = struct{}{}
	}

	members := make([]TeamMember, 0, len(t.Members))
	for _, m := range t.Members {
		if !m.IsActive {
			continue
		}
		if _, blocked := exclude[m.ID]; blocked {
			continue
		}
		members = append(members, m)
	}

	return members
}
