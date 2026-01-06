
# Concurrency Limiter Middleware

This middleware functions to limit the number of requests that can be processed globally at the same time. This feature is very useful to prevent the server from being overwhelmed by sudden request spikes and to maintain service stability.


- **Custom Default Response:**
    Developers can set a custom default response when declaring a new limiter by using the `DefResError{}` option. This allows you to control the response sent when a request is rejected due to concurrency limits.

    **Example:**
    ```go
    limiter := concurrencylimiter.NewConcurrencyLimit(10, concurrencylimiter.DefResError{
            Status: 429, // HTTP status code
            Message: "oops! please try again later.",
    })
    ```

    The above example will return a 429 status code and a custom message instead of the default 503 response.

## Usage Global Limit

1. **Initialize the Concurrency Limiter**

Create a limiter instance with the maximum number of concurrent requests allowed:

```go
import "github.com/he-end/simproute/concurrency_limiter"

// Example: limit to a maximum of 10 concurrent requests
limiter := concurrencylimiter.NewConcurrencyLimit(10)
```

// Example with custom default response
limiter := concurrencylimiter.NewConcurrencyLimit(10, concurrencylimiter.DefResError{
        Status: 429,
        Message: "oops! please try again later.",
})

2. **Attach the Middleware**

Use `MwCCLimit()` as middleware on your router/handler:

```go
mux := http.NewServeMux()
// ...
handler := limiter.MwCCLimit()(yourHandler)
mux.Handle("/endpoint", handler)
```

## Explanation
- If the number of requests exceeds the specified capacity, new requests will be immediately rejected with a `503 Service Unavailable` error response.
- This middleware is suitable for use at the application's entry point (global middleware) to limit server load.

## Notes
- Make sure the capacity is adjusted to your server's capability.
- This middleware does not perform queueing; requests exceeding the limit will be immediately rejected.

# Usage Per-Handler Limit

This feature allows you to limit the number of requests specifically for each handler or route. It is suitable if you want to restrict access to certain endpoints without affecting others.

## How to Use

Suppose you already have a handler function named `User`:

```go
func User(w http.ResponseWriter, r *http.Request) {
    // ...handler implementation...
}
```

To limit a maximum of 10 concurrent requests on the `/user` endpoint, use the following:

```go
import "github.com/he-end/simproute/concurrency_limiter"

r.Get("/user", concurrencylimiter.PerHandler(10, User))
```

## Explanation
- Each handler wrapped with `PerHandler` will have its own concurrency limit.
- If the number of requests exceeds the specified capacity, new requests will be immediately rejected with a `503 Service Unavailable` response.
- There is no queueing; requests exceeding the limit are immediately rejected.

## Notes
- Make sure the limit value is adjusted to the needs and capabilities of the server for each endpoint.
- This feature is very useful for endpoints that require large resources or are prone to overload.
