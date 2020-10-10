package main

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws/credentials"

	"github.com/stretchr/testify/assert"
)

func TestSigninUrl(t *testing.T) {
	cmd := Console{}

	tests := map[string]struct {
		request struct {
			credentials credentials.Value
			timeout     string
		}
		response string
		err      error
	}{
		"timeout_is_ignored_for_sts": {
			response: "https://signin.aws.amazon.com/federation?Action=getSigninToken&Session=%7B%22sessionId%22%3A%22key%22%2C%22sessionKey%22%3A%22secret%22%2C%22sessionToken%22%3A%22token%22%7D",
			err:      nil,
			request: struct {
				credentials credentials.Value
				timeout     string
			}{
				credentials: credentials.Value{
					AccessKeyID:     "key",
					SecretAccessKey: "secret",
					SessionToken:    "token",
				},
				timeout: "1",
			},
		},
		"timeout_is_passed_for_iam_creds": {
			response: "https://signin.aws.amazon.com/federation?Action=getSigninToken&Session=%7B%22sessionId%22%3A%22key%22%2C%22sessionKey%22%3A%22secret%22%2C%22sessionToken%22%3A%22%22%7D&SessionDuration=1",
			err:      nil,
			request: struct {
				credentials credentials.Value
				timeout     string
			}{
				credentials: credentials.Value{
					AccessKeyID:     "key",
					SecretAccessKey: "secret",
				},
				timeout: "1",
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			r, err := cmd.signinURL(test.request.credentials, test.request.timeout)
			assert.Equal(t, test.response, r)
			assert.Equal(t, test.err, err)
		})
	}
}

func TestConsoleUrl(t *testing.T) {
	cmd := Console{}

	tests := map[string]struct {
		request struct {
			signinToken    string
			destinationUrl string
		}
		response string
		err      error
	}{
		"siginurl_is_returned_as_expected": {
			response: "https://signin.aws.amazon.com/federation?Action=login&Destination=http%3A%2F%2Fgoogle.com.au&Issuer=awscli-console-plugin&SigninToken=token",
			request: struct {
				signinToken    string
				destinationUrl string
			}{
				signinToken:    "token",
				destinationUrl: "http://google.com.au",
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			r := cmd.consoleURL(test.request.signinToken, test.request.destinationUrl)
			assert.Equal(t, test.response, r)
		})
	}
}
