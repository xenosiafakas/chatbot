package utils

import (
	"context"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/shomali11/slacker"
	"github.com/slack-go/slack"
	"github.com/xenosiafakas/review-chatbot/pkg/models"
)

var (
	customer models.Customer
	product  models.Product
	Rating float64
	NewReview string
)

type BotState int

const (
	StateInitial BotState = iota
	StateRating
	StateReview
	LastState
)

var currentState = make(map[string]BotState)

func GetData(newCustomer models.Customer, newProduct models.Product) {
	customer = newCustomer
	product = newProduct
	currentState = make(map[string]BotState)
	InitiateConversation()
}

func containsAnyWord(input string, words ...string) bool {
	input = strings.ToLower(input)
	for _, word := range words {
		if strings.Contains(input, word) {
			return true
		}
	}

	return false
}

func findRating(input string) (bool, float64) {

	match := regexp.MustCompile(`(\d+(\.\d+)?)`).FindString(input)

	if match == "" {
		return false, 0
	}
	rating, _ := strconv.ParseFloat(match, 64)

	return true, rating
}

func InitiateConversation() {
	os.Setenv("SLACK_BOT_TOKEN")
	os.Setenv("SLACK_APP_TOKEN")

	bot := slacker.NewClient(os.Getenv("SLACK_BOT_TOKEN"), os.Getenv("SLACK_APP_TOKEN"))

	bot.Init(func() {
		message := fmt.Sprintf("Hello again %s! We noticed you've recently received your %s. We'd love to hear about your experience. Can you spare a few minutes to share your thoughts?", customer.Name, product.Name)
		bot.SocketModeClient().Client.PostMessage("C069G7EMN74", slack.MsgOptionText(message, false))
	})

	Conversation(bot)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := bot.Listen(ctx)
	if err != nil {
		log.Fatal(err)
	}
}

func Conversation(bot *slacker.Slacker) {
	bot.Command("<response>", &slacker.CommandDefinition{
		Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
			userResponse := strings.ToLower(request.Param("response"))

			state := currentState[botCtx.Event().UserID]

			switch state {
			case StateInitial:
				handleInitialResponse(botCtx, request, response, userResponse)
			case StateRating:
				handleRatingResponse(botCtx, request, response, userResponse)
			case StateReview:
				handleReviewResponse(botCtx, request, response, userResponse)
			default:

			}
		},
	})
}


func handleInitialResponse(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter, userResponse string) {
	positive := containsAnyWord(userResponse, "yes", "sure", "of course", "i would love to", "ok")
	negative := containsAnyWord(userResponse, "no", "not", "dont", "nope")
	defaultMessage := "I'm sorry, I didn't understand your response."

	if positive {
		currentState[botCtx.Event().UserID] = StateRating
		message := fmt.Sprintf("Fantastic! On a scale of 1-5, how would you rate the %s", product.Name)
		response.Reply(message)
	} else if negative {
		currentState[botCtx.Event().UserID] = LastState
		message := "We are sorry to hear that... Have a great day"
		response.Reply(message)
	} else {
		response.Reply(defaultMessage)
	}
}

func handleRatingResponse(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter, userResponse string) {
	foundRating, number := findRating(userResponse)
	defaultMessage := "I'm sorry, I didn't understand your response."

	if foundRating {
		if number>=1.0 && number<=5.0{
		Rating = number
		message := fmt.Sprintf("Fantastic! On a scale of 1-5, you rated the product %.1f stars. Since you have been so helpful please leave a product review as well.", number)
		response.Reply(message)
		currentState[botCtx.Event().UserID] = StateReview}else{
		response.Reply("Please provide a valid rating")
		}
	} else {
		response.Reply(defaultMessage)
	}
}

func handleReviewResponse(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter, userResponse string) {
		
		currentState[botCtx.Event().UserID] = LastState

		NewReview = userResponse
	    newReview := models.Review{
			Rating:  int(Rating),
			Comment: NewReview,
			CustomerName: customer.Name,
			ProductName: product.Name,
	
		}
		newReview.CreateReview()

		message := "Thank you for the review"
		response.Reply(message)
		fmt.Print(customer)
}
