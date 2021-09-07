package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/mainak90/lambda-okta/src/cache"
	"github.com/mainak90/lambda-okta/src/cachingclient"
	"github.com/mainak90/lambda-okta/src/client"
	okta_actions "github.com/mainak90/lambda-okta/src/okta-actions"
	"os"
)

var (
	SecretParameterName = os.Getenv("SECRET_PARAM")
	secret = cachingclient.GetSecretCached(SecretParameterName)
	OktaUrl = os.Getenv("OKTA_URL")
	url = cachingclient.GetSecretCached(OktaUrl)
	ApiToken = os.Getenv("API_TOKEN")
	token = cachingclient.GetSecretCached(ApiToken)
)

func handle(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	cfg := client.DefaultConfig()
	var err error
	//x-auth-header from okta-actions
	if secret == "" {
		secret, err = cache.GenerateSecretCache(cfg, SecretParameterName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error %s\n", err)
			return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 503}, nil
		}
	}
	if url == "" {
		secret, err = cache.GenerateSecretCache(cfg, OktaUrl)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error %s\n", err)
			return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 503}, nil
		}
	}
	if token == "" {
		secret, err = cache.GenerateSecretCache(cfg, ApiToken)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error %s\n", err)
			return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 503}, nil
		}
	}
	//Handle the verification first
	method := req.HTTPMethod
	switch method {
	case "GET":
		if val, ok := req.Headers["x-okta-actions-verification-challenge"]; ok {
			resp := fmt.Sprintf("{\"verification\": %s }", val)
			return events.APIGatewayProxyResponse{Body: resp, StatusCode: 201}, nil
		} else {
			return events.APIGatewayProxyResponse{Body: "{ \"Error\": \"x-okta-actions-verification-challenge header missing\" }", StatusCode: 503}, nil
		}
	case "POST":
		var userid string
		var groupid string
		var respB map[string]interface{}
		json.Unmarshal([]byte(req.Body), &respB)
		action1:=respB["data"].(map[string]interface{})["data"].(map[string]interface{})["events"].([]interface{})[0].(map[string]interface{})["target"].([]interface{})[0].(map[string]interface{})
		action2:=respB["data"].(map[string]interface{})["data"].(map[string]interface{})["events"].([]interface{})[0].(map[string]interface{})["target"].([]interface{})[1].(map[string]interface{})
		if action1["type"] == "User" {
			userid = fmt.Sprintf("%v", action1["id"])
			groupid = fmt.Sprintf("%v", action2["id"])
		} else {
			userid = fmt.Sprintf("%v", action2["id"])
			groupid = fmt.Sprintf("%v", action1["id"])
		}
		err := okta_actions.RemoveUserFromGroup(context.Background(), url, token, userid, groupid)
		if err != nil {
			fmt.Println(err.Error())
			return events.APIGatewayProxyResponse{Body:err.Error(), StatusCode: 500}, nil
		}
	}
	return events.APIGatewayProxyResponse{Body: "{\"Error\": \"bad request\"}", StatusCode: 400}, nil
}

func main() {
	lambda.Start(handle)
}