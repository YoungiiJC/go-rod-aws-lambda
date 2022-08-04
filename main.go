package main

import (
	"context"
	"fmt"

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
	const requiredQueryParameter = "url"

	// parse the request
	url, ok := request.QueryStringParameters[requiredQueryParameter]
	if !ok {
		err := fmt.Sprintf(`Query parameter "%s" is missing or incorrect`, requiredQueryParameter)
		return resp(err, 400)
	}

	// get the page's html
	html, err := getPageHTML(url)
	if nil != err {
		err := fmt.Sprintf("Failed to fetch page HTML. Error: %s", err.Error())
		return resp(err, 500)
	}

	return resp(html, 200)
}

func main() {
	lambda.Start(handler)
}
