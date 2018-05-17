// step_activity
package main

import (
	"strings"
)

type ActivityStep struct {
	sent   bool
	answer string
}

func (a *ActivityStep) IsQuestionSubmited() bool {
	return a.sent
}

func (a *ActivityStep) GetQuestion(args []string) string {
	a.sent = true
	return "Are you looking for specific activity ? Parts, sales or service ?"
}

func (a *ActivityStep) SubmitAnswer(args []string) {
	res := strings.ToLower(args[1])
	if strings.Contains(res, "part") {
		a.answer = "PARTS"
	}

	if strings.Contains(res, "sale") {
		a.answer = "SALES"
	}

	if strings.Contains(res, "service") {
		a.answer = "SERVICE"
	}

}

func (a *ActivityStep) GetAnswer() string {
	return a.answer
}

func (a *ActivityStep) GetStepName() string {
	return "activity"
}
