// wizard
package main

import (
	"fmt"

	mbapi "github.com/MarinX/go-mercedes-dealer"
)

type Wizard struct {
	stepId int
	params *mbapi.DealerParam
	steps  []Step
}

type Step interface {
	IsQuestionSubmited() bool
	GetQuestion([]string) string
	SubmitAnswer([]string)
	GetAnswer() string
	GetStepName() string
}

func NewWizard() *Wizard {
	return &Wizard{
		params: new(mbapi.DealerParam),
		steps: []Step{
			&CityStep{},
			&ActivityStep{},
		},
	}
}

func (w *Wizard) GetStep(next bool) Step {
	if next {
		w.stepId += 1
	}
	if w.isEnd() {
		return nil
	}
	return w.steps[w.stepId]

}

// append valid answers to query api
func (w *Wizard) SubmitStep(args []string) {
	s := w.steps[w.stepId]
	s.SubmitAnswer(args)
	switch s.GetStepName() {
	case "city":
		w.params.City = s.GetAnswer()
		break
	case "activity":
		w.params.Activity = s.GetAnswer()
		break
	}
}

// dealer params
func (w *Wizard) GetParams() *mbapi.DealerParam {
	return w.params
}

func (w *Wizard) GetWelcomeText(name string) string {
	return fmt.Sprintf("Beep beep, I'am Mercedes Benz Dealer Bot.\n%s, let's find a dealer by filling some questions", name)
}

func (w *Wizard) isEnd() bool {
	if w.stepId >= len(w.steps) {
		return true
	}
	return false
}
