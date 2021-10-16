package assignment

import "time"

type Short struct {
	CandidateName string
	GithubRepoUrl string
}

type Full struct {
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

type SentDetails struct {
	ID           int64
	RecruiterID  int64
	CandidateUID string
	BusinessID   int64
}
