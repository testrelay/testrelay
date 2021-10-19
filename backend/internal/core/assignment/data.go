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
	SchedulerID        *string   `json:"step_arn"`
	CandidateEmail     string    `json:"candidate_email"`
	TestTimezoneChosen *string   `json:"test_timezone_chosen"`
}

type WithTestDetails struct {
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
	SchedulerID        string    `json:"step_arn"`
	Candidate          Candidate `json:"candidate"`
	Recruiter          Recruiter `json:"recruiter"`
	Test               Test      `json:"test"`
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
	Business   Business `json:"business"`
	Name       string   `json:"name"`
	GithubRepo string   `json:"github_repo"`
}

type Business struct {
	Name string `json:"name"`
}

type SentDetails struct {
	ID           int64
	RecruiterID  int64
	CandidateUID string
	BusinessID   int64
}
