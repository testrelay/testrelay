package assignmentuser

import (
	"github.com/testrelay/testrelay/backend/internal/core/assignment"
	"github.com/testrelay/testrelay/backend/internal/core/user"
)

type RawReviewer struct {
	ID           int
	UserID       int
	AssignmentID int
}

type ReviewerDetail struct {
	User       user.Short
	Assignment assignment.Short
}
