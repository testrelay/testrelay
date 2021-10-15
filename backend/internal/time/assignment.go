package time

import (
	"fmt"
	"time"
)

type AssignmentChoices struct {
	DayChosen  string
	TimeChosen string
	Timezone   string
}

type ScheduleOutput struct {
	StartAssignmentAt  string
	SendNotificationAt string
}

func Parse(input AssignmentChoices) (*ScheduleOutput, error) {
	t, err := time.Parse("2006-01-02 15:04:05",
		fmt.Sprintf("%s %s", input.DayChosen, input.TimeChosen),
	)
	if err != nil {
		return nil, fmt.Errorf(
			"could not parse time input given, day: %s time: %s err: %w",
			input.DayChosen,
			input.TimeChosen,
			err,
		)
	}

	loc, err := time.LoadLocation(string(input.Timezone))
	if err != nil {
		return nil, fmt.Errorf("could not load location %s err: %w", input.Timezone, err)
	}
	t = t.In(loc)

	return &ScheduleOutput{
		StartAssignmentAt:  t.Format(time.RFC3339),
		SendNotificationAt: t.Add(-(time.Minute * 5)).Format(time.RFC3339),
	}, nil
}
