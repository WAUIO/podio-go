package podio

import (
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"regexp"
	"time"
)

var (
	useStubApi = false
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func random(min int, max int) int {
	return rand.Intn(max-min) + min
}

func podioError(r *http.Request, err ...string) string {
	var errType, description string
	description = err[0]
	errType = err[0]
	if len(err) > 1 {
		errType = err[1]
	}

	return fmt.Sprintf(`
{
  "request": {
    "url": "%s"
  },
  "error_description": "%s",
  "error" : "%s"
}`, r.RequestURI, description, errType)
}

func UseStub() {
	useStubApi = true
}

var randomErrCode = func() int {
	return random(399, 520)
}

func SetErrCodeFunc(errCodeFunc func() int) {
	randomErrCode = errCodeFunc
}

// only appliable for usual request, authentication are not affected
func (client *Client) useStub() {
	var api = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var resp string

		time.Sleep(time.Duration(random(500, 1500)) * time.Millisecond)

		w.Header().Set("X-Podio-Auth-Ref", "app_19162664")
		w.Header().Set("X-Rate-Limit-Limit", "1000")
		w.Header().Set("X-Rate-Limit-Remaining", fmt.Sprintf("%d", random(0, 1000)))

		// simulate 4xx and some 5xx errors
		errCode := randomErrCode()
		switch errCode {
		case 400:
			http.Error(w, podioError(r, "bad_request"), http.StatusBadRequest)
			return
		case 401:
			http.Error(w, podioError(r, "unauthorized"), http.StatusUnauthorized)
			return
		case 403:
			http.Error(w, podioError(r, "forbidden access", "forbidden"), http.StatusForbidden)
			return
		case 404:
			http.Error(w, podioError(r, "item not found", "gone"), http.StatusNotFound)
			return
		case 405:
			http.Error(w, podioError(r, fmt.Sprintf("method '%s' not allowed", r.Method)), http.StatusMethodNotAllowed)
			return
		case 408:
			http.Error(w, podioError(r, "request timeout", "timeout"), http.StatusRequestTimeout)
			return
		case 410:
			http.Error(w, podioError(r, "gone"), http.StatusGone)
			return
		case 420:
			w.Header().Set("X-Rate-Limit-Remaining", "0")
			http.Error(w, podioError(r, "hit rate limit", "rate_limit"), 420)
			return
		case 500:
			http.Error(w, podioError(r, "server error", "server_error"), http.StatusInternalServerError)
			return
		case 502:
			http.Error(w, podioError(r, "bad gateway", "bad_gateway"), http.StatusBadGateway)
			return
		case 503:
			http.Error(w, podioError(r, "service unavailable"), http.StatusServiceUnavailable)
			return
		}

		itemReq, _ := regexp.Compile("^/item/([0-9]+)$")

		switch {
		case itemReq.MatchString(r.RequestURI):
			m := itemReq.FindStringSubmatch(r.RequestURI)
			resp = fmt.Sprintf(
				`
{
  "app": {
    "app_id": 19162664
  },
  "push": {
    "expires_in": 300,
    "channel": "/item/%s",
    "signature": "11a8bd86ad2c442023d5be1509c1d85e4e0404d0"
  },
  "item_id": %s,
  "title": "Mock item #%s"
}
`, m[1], m[1], m[1])
		default:
			http.Error(w, fmt.Sprintf("%s not handled", r.RequestURI), http.StatusNotFound)
			return
		}

		buff := []byte(resp)

		if len(buff) > 0 {
			buff = buff[:len(buff)-1]
		}

		if _, err := w.Write(buff); err != nil {
			w.WriteHeader(500)
		}
	}))

	client.URL = api.URL
}
