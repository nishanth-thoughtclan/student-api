package utils

import (
	"context"
	"log"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"google.golang.org/api/option"
)

func InitializeFirebase() *firebase.App {
	opt := option.WithCredentialsFile("google-services.json")
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("error initializing firebase app: %v", err)
	}
	return app
}

func GenerateFirebaseToken(userID string) (string, error) {
	app := InitializeFirebase()
	client, err := app.Auth(context.Background())
	if err != nil {
		return "", err
	}

	token, err := client.CustomToken(context.Background(), userID)
	if err != nil {
		return "", err
	}

	return token, nil
}

func ValidateFirebaseToken(idToken string) (*auth.Token, error) {
	app := InitializeFirebase()
	client, err := app.Auth(context.Background())
	if err != nil {
		return nil, err
	}

	token, err := client.VerifyIDToken(context.Background(), idToken)
	if err != nil {
		return nil, err
	}

	return token, nil
}
