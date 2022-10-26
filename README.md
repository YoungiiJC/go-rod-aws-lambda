Basic [Rod](https://github.com/go-rod/rod) implementation in AWS Lambda. Spins up a [function url](https://docs.aws.amazon.com/lambda/latest/dg/lambda-urls.html).

Deploy it via the AWS SAM cli:

```bash
sam build && sam deploy --guided
```

then test it out in the console with this Event JSON:

```json
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
curl <your_function_url>?url=https://www.uber.com
```

Feel free to open an issue with questions or a PR if we could be doing something better!