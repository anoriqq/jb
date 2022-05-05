package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/anoriqq/jb/internal/config"
	"github.com/pkg/errors"
)

// client

type client struct {
	client *http.Client
	cfg    config.Config
}

type PostChatCommandResp struct {
	OK               bool   `json:"ok"`
	Error            string `json:"error"`
	ResponseMetadata struct {
		Messages string `json:"messages"`
	}
}

func (c *client) PostChatCommand(touchChannel, text string) (*PostChatCommandResp, error) {
	u, err := url.Parse("https://clustervr.slack.com/api/chat.command")
	if err != nil {
		return nil, errors.WithStack(err)
	}

	v := url.Values{}
	v.Set("command", "/jobcan_touch")
	v.Set("channel", touchChannel)

	if len(text) > 0 {
		blocks := fmt.Sprintf(`[{"type":"rich_text","elements":[{"type":"rich_text_section","elements":[{"type":"text","text":"%v"}]}]}]`, text)
		v.Set("blocks", blocks)
	}

	req, err := http.NewRequest(http.MethodPost, u.String(), bytes.NewReader([]byte(v.Encode())))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer resp.Body.Close()

	var postChatCommandResp PostChatCommandResp
	err = json.NewDecoder(resp.Body).Decode(&postChatCommandResp)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(postChatCommandResp.Error)
	}

	return &postChatCommandResp, nil
}

type PostAuthTestResp struct {
	OK    bool   `json:"ok"`
	Error string `json:"error"`
	Team  string `json:"team"`
	User  string `json:"user"`
}

func (c *client) PostAuthTest() (*PostAuthTestResp, error) {
	u, err := url.Parse("https://clustervr.slack.com/api/auth.test")
	if err != nil {
		return nil, errors.WithStack(err)
	}

	req, err := http.NewRequest(http.MethodPost, u.String(), nil)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer resp.Body.Close()

	var postAuthTestResp PostAuthTestResp
	err = json.NewDecoder(resp.Body).Decode(&postAuthTestResp)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(postAuthTestResp.Error)
	}

	return &postAuthTestResp, nil
}

func New(cfg config.Config) *client {
	c := http.DefaultClient
	c.Transport = newBearerTokenRoundTripper(cfg.D, cfg.Token)

	return &client{
		client: c,
	}
}

// RoundTripper

type bearerTokenRoundTripper struct {
	d     string
	token string
}

func (r *bearerTokenRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	cookie := &http.Cookie{
		Name:  "d",
		Value: r.d,
	}
	req.AddCookie(cookie)

	v := url.Values{}

	if req.Body != nil {
		b, err := io.ReadAll(req.Body)
		if err != nil {
			return nil, err
		}

		v, err = url.ParseQuery(string(b))
		if err != nil {
			return nil, err
		}
	}

	v.Set("token", r.token)

	req.Body = io.NopCloser(strings.NewReader(v.Encode()))

	contentLength := len(v.Encode())
	req.ContentLength = int64(contentLength)
	req.Header.Set("content-length", strconv.Itoa(contentLength))
	req.Header.Set("content-type", "application/x-www-form-urlencoded")

	return http.DefaultTransport.RoundTrip(req)
}

func newBearerTokenRoundTripper(d, token string) http.RoundTripper {
	return &bearerTokenRoundTripper{
		d:     d,
		token: token,
	}
}
