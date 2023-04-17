package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"

	"github.com/aws/aws-sdk-go/aws/credentials"

	"net/http"
	"net/url"

	"github.com/alecthomas/kong"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/skratchdot/open-golang/open"
)

type Context struct {
	Timeout int
	Profile string
	Output  bool
}

type Console struct {
}

var cli struct {
	Command Console `cmd help:"login to aws console" default:"1"`
	Timeout int     `help:"Console session timeout in seconds (for non STS credentials only)" default:"360"`
	Profile string  `help:"Assume credentials in the given profile before opening the console."`
	Output  bool    `help:"Output the signin url instead of opening it in the browser."`
}

func (cli *Console) signinURL(credentials credentials.Value, timeout string) (string, error) {
	credentialsJSON, err := json.Marshal(map[string]string{
		"sessionId":    credentials.AccessKeyID,
		"sessionKey":   credentials.SecretAccessKey,
		"sessionToken": credentials.SessionToken,
	})
	if err != nil {
		return "", fmt.Errorf("error while marshalling credentials, %s", err)
	}

	signinQuery := url.Values{}
	signinQuery.Add("Action", "getSigninToken")
	// signinQuery.Add("SessionType", "json")

	if credentials.SessionToken == "" {
		signinQuery.Add("SessionDuration", timeout)
	}

	signinQuery.Add("Session", string(credentialsJSON))
	signinURL := url.URL{}
	signinURL.Host = "signin.aws.amazon.com"
	signinURL.Scheme = "https"
	signinURL.Path = "federation"
	signinURL.RawQuery = signinQuery.Encode()

	return signinURL.String(), nil
}

func (cli *Console) consoleURL(signinToken string, destinationURL string) string {
	consoleQuery := url.Values{}

	consoleQuery.Add("Action", "login")
	consoleQuery.Add("Issuer", "awscli-console-plugin")
	consoleQuery.Add("Destination", destinationURL)
	consoleQuery.Add("SigninToken", signinToken)

	consoleURL := url.URL{}
	consoleURL.Host = "signin.aws.amazon.com"
	consoleURL.Scheme = "https"
	consoleURL.Path = "federation"
	consoleURL.RawQuery = consoleQuery.Encode()

	return consoleURL.String()
}

func (cli *Console) Credentials(profile string) (credentials.Value, string, error) {
	var sess *session.Session

	if profile != "" {
		sess = session.Must(session.NewSessionWithOptions(session.Options{
			Profile: profile,
			// according to some otherwise unrelated Github Issue
			SharedConfigState: session.SharedConfigEnable,
		}))
	} else {
		sess = session.Must(session.NewSessionWithOptions(session.Options{
			SharedConfigState: session.SharedConfigEnable,
		}))
	}

	credentials, err := sess.Config.Credentials.Get()

	return credentials, aws.StringValue(sess.Config.Region), err
}

func (cli *Console) Run(ctx *Context) error {
	credentials, region, err := cli.Credentials(ctx.Profile)
	if err != nil {
		return fmt.Errorf("error while retrieving credentials, %s", err)
	}

	signinURL, err := cli.signinURL(credentials, strconv.Itoa(ctx.Timeout))
	if err != nil {
		return err
	}

	resp, err := http.Get(signinURL)
	if err != nil {
		return fmt.Errorf("error while making signin request to %s, %s", signinURL, err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unsuccessful response from the signin request, %s", resp.Status)
	}

	responseBody := map[string]string{}
	err = json.NewDecoder(resp.Body).Decode(&responseBody)
	if err != nil {
		return fmt.Errorf("error while unmarshalling response body, %s", err)
	}

	_, ok := responseBody["SigninToken"]
	if !ok {
		return errors.New("missing 'SigninToken' field in the response body")
	}

	destinationURL := fmt.Sprintf("https://%s.console.aws.amazon.com/console/home?region=%s", region, region)

	consoleURL := cli.consoleURL(responseBody["SigninToken"], destinationURL)

	if ctx.Output {
		fmt.Fprint(os.Stdout, consoleURL)
		return nil
	}

	err = open.Run(consoleURL)
	if err != nil {
		return fmt.Errorf("error while opening browser, %s", err)
	}

	return nil
}

func main() {
	ctx := kong.Parse(&cli)

	err := ctx.Run(&Context{
		Timeout: cli.Timeout,
		Profile: cli.Profile,
		Output:  cli.Output})

	ctx.FatalIfErrorf(err)
}
