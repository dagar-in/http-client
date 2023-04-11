# http-client

http-client is a custom HTTP client package that provides a fluent interface for making HTTP requests and handling HTTP responses.

## Installation

To install http-client, use the go get command:

```bash
go get github.com/dagar-in/http-client
```

## Usage

To use http-client, import the package in your Go file:
```go
import "github.com/dagar-in/http-client"
```
Then create a new client instance using the New function:
```go
client := request.New()
```
You can use the client instance to make HTTP requests using the Get, Post, Put, Patch, or Delete methods. These methods take a URL string as an argument and return a Response instance or an error. For example:
```go
resp, err := client.Get("https://example.com")
if err != nil {
    // handle error
}
```
You can also chain methods to configure the request before sending it. For example, you can use the WithHeaders method to set the request headers, the WithQuery method to set the query parameters, and the WithBody method to set the request body. For example:
```go
resp, err := client.WithHeaders(map[string]string{
    "Content-Type": "application/json",
}).WithQuery(map[string]string{
    "q": "golang",
}).WithBody([]byte(`{"name":"Alice"}`)).Post("https://example.com")
if err != nil {
    // handle error
}
```
You can use the Response instance to access the response status code, headers, and body. For example, you can use the StatusCode method to get the status code as an int, the Header method to get the headers as an http.Header map, and the BodyMap method to get the body as a map[string]interface{}. For example:
```go
fmt.Println(resp.StatusCode())
fmt.Println(resp.Header())
bodyMap, err := resp.BodyMap()
if err != nil {
    // handle error
}
fmt.Println(bodyMap)
```
You can also use the DoAll method to make multiple requests with or without concurrency. This method takes a method string, a slice of URLs, and a boolean flag as arguments and returns a slice of Response instances or an error. For example:
```go
urls := []string{
    "https://example.com/foo",
    "https://example.com/bar",
    "https://example.com/baz",
}
responses, err := client.DoAll("GET", urls, true)
if err != nil {
    // handle error
}
for _, resp := range responses {
    fmt.Println(resp.StatusCode())
}
```
## Error Handling
The httpclient package returns errors that may occur during parsing, creating, executing, or reading requests or responses. The errors are wrapped with more context using fmt.Errorf and %w verb. You can use the errors.Is or errors.As functions to check or unwrap the original errors if needed. For example:
```go
resp, err := client.Post("invalid_url")
if err != nil {
    fmt.Println(err)
    var urlErr *url.Error
    if errors.As(err, &urlErr) {
        fmt.Println(urlErr)
    }
}
```