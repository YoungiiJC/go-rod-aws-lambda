
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func resp(body string, status int) (*events.APIGatewayProxyResponse, error) {
	return &events.APIGatewayProxyResponse{
		StatusCode: status,
		Headers:    map[string]string{"Content-Type": "text/html"},
		Body:       body,
	}, nil
}

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	const (
		urlParam      = "url"
		emailParam    = "email"
		passwordParam = "password"
		proxyParam    = "proxy"
	)
	var requiredQueryParameters = []string{urlParam, emailParam, passwordParam, proxyParam}

	// Log the received query parameters
	log.Printf("Received query parameters: %v", request.QueryStringParameters)

	// Parse the request
	for _, requiredQueryParameter := range requiredQueryParameters {
		value, ok := request.QueryStringParameters[requiredQueryParameter]
		if !ok || value == "" {
			err := fmt.Sprintf(`Query parameter "%s" is missing or incorrect`, requiredQueryParameter)
			log.Println(err)
			return resp(err, 400)
		}
	}

	urle := request.QueryStringParameters[urlParam]
	email := request.QueryStringParameters[emailParam]
	password := request.QueryStringParameters[passwordParam]
	proxy := request.QueryStringParameters[proxyParam]

	// Log the parsed parameters
	log.Printf("Parsed parameters - URL: %s, Email: %s, Password: %s, Proxy: %s", urle, email, password, proxy)

	// Get the page's HTML
	html, err := getPageHTML(urle, email, password, proxy)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to fetch page HTML. Error: %s", err.Error())
		log.Println(errMsg)
		return resp(errMsg, 500)
	}

	return resp(html, 200)
}

func main() {
	lambda.Start(handler)
}
