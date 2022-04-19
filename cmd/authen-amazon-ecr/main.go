package main

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/ecr"
)

func main() {
	os.Exit(run())
}

var envCreds = credentials.NewStaticCredentialsProvider(os.Getenv("AWS_ACCESS_KEY_ID"), os.Getenv("AWS_SECRET_ACCESS_KEY"), os.Getenv("AWS_SESSION_TOKEN"))

func run() int {
	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx, config.WithCredentialsProvider(envCreds))
	if err != nil {
		fmt.Println(err)
		return 1
	}
	client := ecr.NewFromConfig(cfg)
	out, err := client.GetAuthorizationToken(ctx, nil)
	if err != nil {
		fmt.Println(err)
		return 1
	}
	creds, err := getCredentials(out)
	if err != nil {
		fmt.Println(err)
		return 1
	}
	output("username", creds.user)
	output("password", creds.password)
	output("server", creds.server)
	return 0
}

type ecrCredentials struct {
	user     string
	password string
	server   string
}

var protocol = "https://"

func getCredentials(out *ecr.GetAuthorizationTokenOutput) (*ecrCredentials, error) {
	if len(out.AuthorizationData) != 1 {
		return nil, errors.New("malformed output")
	}
	authData := out.AuthorizationData[0]
	if authData.AuthorizationToken == nil {
		return nil, errors.New("AuthorizationToken is empty")
	}
	if authData.ProxyEndpoint == nil {
		return nil, errors.New("ProxyEndpoint is empty")
	}
	proxyEndpoint := *authData.ProxyEndpoint
	if strings.HasPrefix(proxyEndpoint, protocol) {
		proxyEndpoint = proxyEndpoint[len(protocol):]
	}
	decoded, err := base64.StdEncoding.DecodeString(*authData.AuthorizationToken)
	if err != nil {
		return nil, err
	}
	decodedStr := string(decoded)
	idx := strings.Index(decodedStr, ":")
	if idx == -1 {
		return nil, errors.New("malformed AuthorizationToken")
	}
	return &ecrCredentials{user: decodedStr[0:idx], password: decodedStr[idx+1:], server: proxyEndpoint}, nil
}

func output(name, value string) {
	fmt.Fprintf(os.Stdout, "::set-output name=%s::%s\n", name, value)
}
