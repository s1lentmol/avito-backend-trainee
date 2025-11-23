package domain

import (
	"math/rand"
	"time"

	"github.com/silentmol/avito-backend-trainee/internal/apperr"
	teamdomain "github.com/silentmol/avito-backend-trainee/internal/team/domain"
)

func SelectReviewersForTeam(team *teamdomain.Team, authorID string) []string {
	if team == nil {
		return nil
	}

	// берём активных членов команды, кроме автора
	active := team.ActiveMembersExcept(authorID)
	if len(active) == 0 {
		return nil
	}

	candidates := make([]string, 0, len(active))
	for _, m := range active {
		candidates = append(candidates, m.ID)
	}

	if len(candidates) <= 2 {
		return candidates
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	r.Shuffle(len(candidates), func(i, j int) {
		candidates[i], candidates[j] = candidates[j], candidates[i]
	})

	return candidates[:2]
}

func ReassignReviewer(pr *PullRequest, team *teamdomain.Team, oldReviewerID string) (string, error) {
	if pr == nil || team == nil {
		return "", apperr.ErrNoCandidate
	}

	// на уже слитых PR переставлять ревьюера нельзя
	if err := pr.CanReassign(); err != nil {
		return "", err
	}

	reviewerSet := make(map[string]struct{}, len(pr.AssignedReviewers))
	for _, id := range pr.AssignedReviewers {
		reviewerSet[id] = struct{}{}
	}

	if _, ok := reviewerSet[oldReviewerID]; !ok {
		return "", apperr.ErrNotAssigned
	}

	// ищем активных кандидатов вместо старого ревьюера
	activeMembers := team.ActiveMembersExcept(oldReviewerID)

	candidates := make([]string, 0, len(activeMembers))
	for _, member := range activeMembers {
		if _, alreadyAssigned := reviewerSet[member.ID]; alreadyAssigned {
			continue
		}
		candidates = append(candidates, member.ID)
	}

	if len(candidates) == 0 {
		return "", apperr.ErrNoCandidate
	}

	newReviewerID := pickRandomFrom(candidates)

	if err := pr.ReplaceReviewer(oldReviewerID, newReviewerID); err != nil {
		return "", err
	}

	return newReviewerID, nil
}

func pickRandomFrom(candidates []string) string {
	if len(candidates) == 0 {
		return ""
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return candidates[r.Intn(len(candidates))]
}
