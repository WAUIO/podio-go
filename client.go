package podio

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"golang.org/x/net/context"
)

type Client struct {
	httpClient *http.Client
	authToken  *AuthToken
	Emitter    EventEmitterPodioWrapper
	Context    context.Context
	ServiceEndpoint
}

type Error struct {
	Parameters interface{} `json:"error_parameters"`
	Detail     interface{} `json:"error_detail"`
	Propagate  bool        `json:"error_propagate"`
	Request struct {
		URL   string `json:"url"`
		Query string `json:"query_string"`
	} `json:"request"`
	Description string `json:"error_description"`
	Type        string `json:"error"`
}

func (p *Error) Error() string {
	return fmt.Sprintf("%s: %s", p.Type, p.Description)
}

func NewClient(authToken *AuthToken, emiterConf func(e EventEmitterPodioWrapper)) *Client {
	Emitter := GetPodioEmitter()
	emiterConf(Emitter)

	client :=  &Client{
		httpClient: &http.Client{},
		authToken:  authToken,
		Emitter:    Emitter,
	}

	client.URL = "https://api.podio.com"

	if useStubApi {
		client.useStub()
	}

	return client
}

func (client *Client) WithContext(ctx context.Context) *Client{
	client.Context = ctx
	return client
}

func (client *Client) Request(method string, path string, headers map[string]string, body io.Reader, out interface{}) error {
	req, err := http.NewRequest(method, client.URL + path, body)

	if err != nil {
		return err
	}

	client.Emitter.FireBackground("podio.request", method, req.URL.Path, struct {
		Form interface{} `json:"form"`
	}{
		Form:req.PostForm,
	})

	for k, v := range headers {
		req.Header.Add(k, v)
	}

	req.Header.Add("Authorization", "OAuth2 "+client.authToken.AccessToken)
	resp, err := client.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if !(200 <= resp.StatusCode && resp.StatusCode < 300) {
		podioErr := &Error{}
		err := json.Unmarshal(respBody, podioErr)
		if err != nil {
			return errors.New(string(respBody))
		}

		client.Emitter.FireBackground("podio.error", podioErr, resp.StatusCode)
		client.Emitter.FireBackground("podio.context.error", client.Context, podioErr, resp)

		return podioErr
	}

	client.Emitter.FireBackground("podio.response", client.Context, respBody, resp.Header)

	if out != nil {
		return json.Unmarshal(respBody, out)
	}

	return nil
}

func (client *Client) RequestWithParams(method string, path string, headers map[string]string, params map[string]interface{}, out interface{}) error {
	buf, err := json.Marshal(params)
	if err != nil {
		return err
	}

	return client.Request(method, path, headers, bytes.NewReader(buf), out)
}
