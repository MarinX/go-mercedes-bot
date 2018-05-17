// main
package main

import (
	"fmt"
	"os"
	"regexp"

	mbapi "github.com/MarinX/go-mercedes-dealer"
	"github.com/go-chat-bot/bot"
	"github.com/go-chat-bot/bot/slack"
)

const (
	DEALER_TEMPLATE = "Name: %s\nAddress: %s\nWebsite: %s\nEmail Address: %s\nTelephone: %s\nGoogle maps link: \nhttps://www.google.com/maps/?q=%f,%f"
)

var api *mbapi.API
var wizards map[string]*Wizard
var rexp = regexp.MustCompile(`(?m)GS\d\d\d\d\d\d\d`)

func init() {
	bot.RegisterPassiveCommand("mbdealerbot", mbdealerbot)
	wizards = make(map[string]*Wizard)
}

func main() {
	slackKey := os.Getenv("SLACK_KEY")
	mbKey := os.Getenv("MB_KEY")
	if len(slackKey) == 0 {
		fmt.Println("Missing slack token in env variable SLACK_KEY")
		return
	}
	if len(mbKey) == 0 {
		fmt.Println("Missing mercedes benz dev key in env variable MB_KEY")
		return
	}
	api = mbapi.New(mbKey)
	slack.Run(slackKey)
}

func mbdealerbot(command *bot.PassiveCmd) (string, error) {
	userId := command.User.ID
	args := []string{userId, command.Raw}

	if rexp.MatchString(command.Raw) {
		if wizards[userId] != nil {
			delete(wizards, userId)
		}

		return GetDealer(rexp.FindString(command.Raw), command.User.RealName), nil
	}

	// new user, create a flow
	if wizards[userId] == nil {
		wizards[userId] = NewWizard()
		welcomeText := wizards[userId].GetWelcomeText(command.User.RealName)
		return fmt.Sprintf("%s\n\n%s", welcomeText, wizards[userId].GetStep(false).GetQuestion(args)), nil
	}

	// get step
	currStep := wizards[userId].GetStep(false)

	// if the user already got a question,
	// next is to submit answer and get new step
	if currStep.IsQuestionSubmited() {
		wizards[userId].SubmitStep(args)
		currStep = wizards[userId].GetStep(true)
	}

	// no steps left, clear all
	if currStep == nil {
		params := wizards[userId].GetParams()
		delete(wizards, userId)
		return GetDealers(params), nil
	}

	return currStep.GetQuestion(args), nil
}

func GetDealers(params *mbapi.DealerParam) string {
	response, err := api.GetDealers(params)
	if err != nil {
		return fmt.Sprintf("Uh no, error occured :( %s ", err.Error())
	}

	if len(response.Dealers) <= 0 {
		return "Cannot find dealer that match your criteria! Let's try again."
	}

	return generateResponse(response.Dealers)
}

func GetDealer(dealerId string, user string) string {
	msg := fmt.Sprintf("%s here are the details for %s\n\n", user, dealerId)
	dealer, err := api.GetDealer(dealerId)
	if err != nil {
		return fmt.Sprintf("Uh no, error occured :( %s", err.Error())
	}

	msg += fmt.Sprintf(DEALER_TEMPLATE,
		dealer.LegalName,
		dealer.Address.Street,
		dealer.Communication.Website,
		dealer.Communication.Email,
		dealer.Communication.Phone,
		dealer.GeoCoordinates.Latitude,
		dealer.GeoCoordinates.Longitude,
	)

	return msg
}

func generateResponse(dealers []mbapi.Dealer) string {
	msg := fmt.Sprintf("I found some dealers near you! To get detail information about dealer, type dealer id eg %s\n\n", dealers[0].DealerId)
	for _, d := range dealers {
		msg += fmt.Sprintf("ID: %s\nName: %s\n\n", d.DealerId, d.LegalName)
	}

	return msg + "\nHave a nice day! :)"
}
