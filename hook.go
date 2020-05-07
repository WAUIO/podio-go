package podio

import (
	"fmt"
)

type Hook struct {
	Id     uint   `json:"hook_id"`
	Status string `json:"status"`
	Type   string `json:"type"`
	Url    string `json:"url"`
}

// https://developers.podio.com/doc/hooks/create-hook-215056
func (client *Client) HookCreate(refType string, refId uint, params map[string]interface{}) (hook *Hook, err error) {
	path := fmt.Sprintf("/hook/%s/%d", refType, refId)
	err = client.RequestWithParams("POST", path, nil, params, &hook)
	if err != nil {
		hook.Url = params["url"].(string)
		hook.Type = params["type"].(string)
	}

	return
}

// https://developers.podio.com/doc/hooks/get-hooks-215285
func (client *Client) HookGetFor(refType string, refId uint) (hooks []*Hook, err error) {
	path := fmt.Sprintf("/hook/%s/%d", refType, refId)
	err = client.Request("GET", path, nil, nil, &hooks)
	return
}

// https://developers.podio.com/doc/hooks/delete-hook-215291
func (client *Client) HookDelete(hookId uint) (response interface{}, err error) {
	path := fmt.Sprintf("/hook/%d", hookId)
	err = client.RequestWithParams("DELETE", path, nil, nil, &response)
	return
}

// https://developers.podio.com/doc/hooks/request-hook-verification-215232
func (client *Client) HookVerify(hookId uint, params map[string]interface{}) (response interface{}, err error) {
	path := fmt.Sprintf("/hook/%d/verify/request", hookId)
	err = client.RequestWithParams("POST", path, nil, params, &response)
	return
}

// https://developers.podio.com/doc/hooks/validate-hook-verification-215241
func (client *Client) HookValidate(hookId uint, code string) (response interface{}, err error) {
	params := make(map[string]interface{})
	params["code"] = code

	path := fmt.Sprintf("/hook/%d/verify/validate", hookId)
	err = client.RequestWithParams("POST", path, nil, params, &response)
	return
}
