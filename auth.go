package main

import (
	"context"
	"os"

	"golang.org/x/oauth2"
)

func auth2() (ctx context.Context) {
	ctx = context.Background()

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv(Contexts[currentContext].Github.Token)},
	)
	client = oauth2.NewClient(ctx, ts)
	return
}
