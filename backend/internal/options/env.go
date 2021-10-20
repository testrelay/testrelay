package options

import (
	"errors"
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	AppURL     string
	BackendURL string
	APPEnv     string

	SMTPHost       string
	SMTPPort       int64
	SMTPUsername   string
	SMTPPassword   string
	MailFromDomain string

	HasuraURL   string
	HasuraToken string

	GithubAccessToken        string
	GithubPrivateKeyLocation string
	GithubAppID              int64

	GoogleServiceAccountLocation string
}

func ConfigFromEnv() (Config, error) {
	var e errs

	c := Config{
		AppURL:                       envOrDefaultString("APP_URL", "localhost"),
		BackendURL:                   envOrDefaultString("BACKEND_URL", "backend"),
		APPEnv:                       envOrDefaultString("APP_ENV", "development"),
		SMTPHost:                     e.envOrError("SMTP_HOST"),
		SMTPPort:                     e.envOrErrorInt("SMTP_PORT"),
		SMTPUsername:                 os.Getenv("SMTP_USER"),
		SMTPPassword:                 os.Getenv("SMTP_PASS"),
		MailFromDomain:               envOrDefaultString("MAIL_FROM_DOMAIN", "@testrelay.io"),
		HasuraURL:                    envOrDefaultString("HASURA_URL", "hasura"),
		HasuraToken:                  e.envOrError("HASURA_TOKEN"),
		GithubAccessToken:            e.envOrError("GITHUB_ACCESS_TOKEN"),
		GithubPrivateKeyLocation:     e.envOrError("GITHUB_PRIVATE_KEY"),
		GithubAppID:                  e.envOrErrorInt("GITHUB_APP_ID"),
		GoogleServiceAccountLocation: e.envOrError("GOOGLE_SERVICE_ACC"),
	}

	return c, e.Error()
}

type errs []error

func (e errs) Error() error {
	if len(e) == 0 {
		return nil
	}

	var msg string
	for _, err := range e {
		msg += err.Error() + "\n"
	}

	return errors.New(msg)
}

func (e *errs) envOrErrorInt(key string) int64 {
	v := e.envOrError(key)
	if v == "" {
		return 0
	}

	i, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		*e = append(*e, fmt.Errorf("%s is not a valid int", key))
	}

	return i
}

func (e *errs) envOrError(key string) string {
	v := os.Getenv(key)
	if v == "" {
		*e = append(*e, fmt.Errorf("%s must be set to boot application", key))
	}

	return v
}

func envOrDefaultString(key, def string) string {
	v := os.Getenv(key)
	if v != "" {
		return v
	}

	return def
}

func envOrDefaultInt(key string, def int64) int64 {
	v := os.Getenv(key)
	if v == "" {
		return def
	}

	i, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return def
	}

	return i
}
