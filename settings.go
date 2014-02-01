package main

import "os"

func ReadConsumerKey() string {
	return os.Getenv("ConsumerKey")
}

func ReadConsumerSecret() string {
	return os.Getenv("ConsumerSecret")
}

func ReadAccessToken() string {
	return os.Getenv("AccessToken")
}

func ReadAccessTokenSecret() string {
	return os.Getenv("AccessTokenSecret")
}
