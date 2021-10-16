package scheduler

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sfn"
)

type StartInput struct {
	ID           int64       `json:"id"`
	TestStart    string      `json:"testStart"`
	TestDuration int         `json:"testDuration"`
	Data         interface{} `json:"data"`
}

type AssignmentScheduler interface {
	Stop(id string) error
	Start(input StartInput) (string, error)
}

type SFNClient interface {
	StartExecution(input *sfn.StartExecutionInput) (*sfn.StartExecutionOutput, error)
}

type StepFunctionAssignmentScheduler struct {
	StateMachineArn string
	SFNClient       *sfn.SFN
}

func (s StepFunctionAssignmentScheduler) Stop(id string) error {
	_, err := s.SFNClient.StopExecution(&sfn.StopExecutionInput{
		ExecutionArn: aws.String(id),
	})
	if err != nil {
		return fmt.Errorf("could not start aws state machine with arn %s err %w", id, err)
	}

	return nil
}

func (s StepFunctionAssignmentScheduler) Start(input StartInput) (string, error) {
	b, err := json.Marshal(input)
	if err != nil {
		return "", fmt.Errorf("could not marshal input %w", err)
	}

	stateName := fmt.Sprintf("assignment-%d-%d", input.ID, time.Now().Unix())
	out, err := s.SFNClient.StartExecution(&sfn.StartExecutionInput{
		Input:           aws.String(string(b)),
		Name:            aws.String(stateName),
		StateMachineArn: aws.String(s.StateMachineArn),
	})
	if err != nil {
		return "", fmt.Errorf("could not start step func execution arn %s %w", s.StateMachineArn, err)
	}

	return fmt.Sprintf("%p", out.ExecutionArn), nil
}
