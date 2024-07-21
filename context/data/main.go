package main

import (
	"context"
	"fmt"
	"net/http"
)

type ctxKey string

// add value to context
var ctxRequestId ctxKey
var ctxRequestToken ctxKey

func main() {
	token := "Bearer Hello"

	// define context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ctxRequestId = "requestId"
	ctxRequestToken = "authToken"

	ctx = context.WithValue(ctx, ctxRequestId, 1)
	ctx = context.WithValue(ctx, ctxRequestToken, token)

	resp, err := processRequest(ctx)
	if err != nil {
		fmt.Println("request failed", err)
		return
	}

	handleResponse(ctx, resp)
}

func processRequest(ctx context.Context) (interface{}, error) {
	requstURL := fmt.Sprintf("https://jsonplaceholder.typicode.com/posts/%d", ctx.Value(ctxKey(ctxRequestId)))
	resp, err := http.Get(requstURL)
	if err != nil {
		return nil, err
	}

	return resp, err
}

func handleResponse(ctx context.Context, res interface{}) {
	fmt.Println("requestId", ctx.Value(ctxKey(ctxRequestId)))
	fmt.Println("authToken", ctx.Value(ctxKey(ctxRequestToken)))
	fmt.Println(res)
}
