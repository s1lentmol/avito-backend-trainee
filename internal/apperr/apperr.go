package apperr

import "errors"

var (
	ErrNotFound    = errors.New("not found")
	ErrTeamExists  = errors.New("team exists")
	ErrPRExists    = errors.New("pr exists")
	ErrPRMerged    = errors.New("pr merged")
	ErrNotAssigned = errors.New("not assigned to pr")
	ErrNoCandidate = errors.New("no candidate in team")
)
