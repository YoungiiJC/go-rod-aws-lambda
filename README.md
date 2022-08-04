Basic [Rod](https://github.com/go-rod/rod) implementation in AWS Lambda. Spins up an API Gateway endpoint.

Deploy it via the AWS SAM cli:
```bash
sam build && sam deploy --guided
```

then test it out in the console with this Event JSON:
```js
{
  "path": "/",
  "httpMethod": "GET",
  "queryStringParameters": {
    "url": "https://www.uber.com"
  }
}
```

or hit the endpoint directly:
```bash
curl <your_lambda_api_gateway_uri>?url=https://www.uber.com
```

Feel free to open an issue with questions or a PR if we could be doing something better!