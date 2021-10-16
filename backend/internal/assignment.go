package internal

import (
	"encoding/json"
	"time"
)

type AssignmentUser struct {
	ID           int
	UserID       int
	AssignmentID int
}

type AssignmentEvent struct {
	ID           int             `json:"id"`
	UserID       int             `json:"user_id"`
	AssignmentID int             `json:"assignment_id"`
	Meta         json.RawMessage `json:"meta"`
	EventType    string          `json:"event_type"`
	CreatedAt    time.Time       `json:"created_at"`
}

type Assignment struct {
	Status             string    `json:"status"`
	TestTimeChosen     *string   `json:"test_time_chosen"`
	ChooseUntil        string    `json:"choose_until"`
	TestDayChosen      *string   `json:"test_day_chosen"`
	TestId             int       `json:"test_id"`
	TimeLimit          int       `json:"time_limit"`
	UpdatedAt          time.Time `json:"updated_at"`
	CandidateId        *int      `json:"candidate_id"`
	CreatedAt          time.Time `json:"created_at"`
	Id                 int       `json:"id"`
	CandidateName      string    `json:"candidate_name"`
	RecruiterId        int       `json:"recruiter_id"`
	InviteCode         string    `json:"invite_code"`
	GithubRepoURL      *string   `json:"github_repo_url"`
	CandidateEmail     string    `json:"candidate_email"`
	TestTimezoneChosen *string   `json:"test_timezone_chosen"`
}

type FullAssignment struct {
	Status             string    `json:"status"`
	TestTimeChosen     string    `json:"test_time_chosen"`
	ChooseUntil        string    `json:"choose_until"`
	TestDayChosen      string    `json:"test_day_chosen"`
	TestID             int       `json:"test_id"`
	TimeLimit          int       `json:"time_limit"`
	CandidateID        int       `json:"candidate_id"`
	ID                 int       `json:"id"`
	CandidateName      string    `json:"candidate_name"`
	RecruiterID        int       `json:"recruiter_id"`
	InviteCode         string    `json:"invite_code"`
	GithubRepoURL      string    `json:"github_repo_url"`
	CandidateEmail     string    `json:"candidate_email"`
	TestTimezoneChosen string    `json:"test_timezone_chosen"`
	StepArn   string    `json:"step_arn"`
	Candidate Candidate `json:"candidate"`
	Recruiter Recruiter `json:"recruiter"`
	Test      Test      `json:"test"`
}

type Candidate struct {
	Email             string `json:"email"`
	GithubUsername    string `json:"github_username"`
	GithubAccessToken string `json:"github_access_token"`
}

type Recruiter struct {
	Email string `json:"email"`
}

type Test struct {
	Business Business `json:"business"`
	Name     string   `json:"name"`
	GithubRepo string   `json:"github_repo"`
}

type Business struct {
	Name string `json:"name"`
}

type Reviewer struct {
	GithubUsername string `graphql:"github_username" json:"github_username"`
}
