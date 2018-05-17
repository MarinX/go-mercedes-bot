// step_city
package main

type CityStep struct {
	sent   bool
	answer string
}

func (c *CityStep) IsQuestionSubmited() bool {
	return c.sent
}

func (c *CityStep) GetQuestion(args []string) string {
	c.sent = true
	return "In which city do you need a dealer ?"
}

func (c *CityStep) SubmitAnswer(args []string) {
	c.answer = args[1]
}

func (c *CityStep) GetAnswer() string {
	return c.answer
}

func (c *CityStep) GetStepName() string {
	return "city"
}
