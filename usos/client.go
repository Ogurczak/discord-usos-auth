package usos

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/dghubble/oauth1"
)

// ErrUnableToCall represents an error which took place during calling an usos-api method
type ErrUnableToCall struct {
	cause error
}

// newErrUnableToCall returns a pointer to a new ErrUnableToCall
func newErrUnableToCall(cause error) *ErrUnableToCall {
	return &ErrUnableToCall{
		cause: cause,
	}
}

func (e *ErrUnableToCall) Error() string {
	return "Error during calling an usos-api method"
}

// ErrHTTP represtents a HTTP error
type ErrHTTP struct {
	Code    int
	Message string
}

func newErrHTTP(Code int, Message string) *ErrHTTP {
	return &ErrHTTP{
		Code:    Code,
		Message: Message,
	}
}
func (e *ErrHTTP) Error() string {
	return e.Message
}

func usosURL(key string) string {
	var baseURL = "https://apps.usos.pw.edu.pl/"

	var urls = map[string]string{
		"":              "",
		"requestToken":  "services/oauth/request_token",
		"authorize":     "services/oauth/authorize",
		"accessToken":   "services/oauth/access_token",
		"user":          "services/users/user?fields=%s",
		"groups":        "services/groups/user?fields=%s&active_terms=%v",
		"registrations": "services/registrations/user_registrations?fields=%s",
		"term":          "services/terms/term?term_id=%s",
		"courses":       "services/courses/user?fields=%s",
	}
	return baseURL + urls[key]
}

var config = oauth1.Config{
	ConsumerKey:    "774c544Rjd7R3hevEzkg",
	ConsumerSecret: "hFH6hFfEqJmbvHn7VcrPqfchKn357U6mErGN7F2F",
	CallbackURL:    "oob",
	Endpoint: oauth1.Endpoint{
		RequestTokenURL: usosURL("requestToken"),
		AuthorizeURL:    usosURL("authorize"),
		AccessTokenURL:  usosURL("accessToken"),
	},
}

// RequestToken represents oauth1 request token
type RequestToken struct {
	Token            string
	Secret           string
	AuthorizationURL *url.URL
}

// GetAccessToken returns an access token from the request token and verifier
func (rt *RequestToken) GetAccessToken(verifier string) (*oauth1.Token, error) {
	token, secret, err := config.AccessToken(rt.Token, rt.Secret, verifier)
	if err != nil {
		return nil, err
	}
	return oauth1.NewToken(token, secret), nil
}

// NewRequestToken returns an usos unauthorized
func NewRequestToken() (*RequestToken, error) {
	token, secret, err := config.RequestToken()
	if err != nil {
		return nil, err
	}
	authorizationURL, err := config.AuthorizationURL(token)
	if err != nil {
		return nil, err
	}
	return &RequestToken{token, secret, authorizationURL}, nil
}

func makeCall(client *http.Client, key string, a ...interface{}) (io.ReadCloser, error) {
	url := fmt.Sprintf(usosURL(key), a...)
	resp, err := client.Get(url)
	if err != nil {
		return nil, newErrUnableToCall(err)
	}
	if resp.StatusCode != 200 {
		return nil, newErrHTTP(resp.StatusCode, resp.Status)
	}

	return resp.Body, err
}

func printResponse(r *io.Reader) error {
	dat, err := ioutil.ReadAll(*r)
	if err != nil {
		return err
	}
	*r = bytes.NewReader(dat)
	fmt.Println("Response: ", string(dat))
	return nil
}
