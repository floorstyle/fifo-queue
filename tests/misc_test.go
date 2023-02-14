package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"net/url"

	"github.com/floorstyle/fifo-queue/server"
	"github.com/floorstyle/fifo-queue/util"
	"github.com/pkg/errors"
)

const LOWERCASE_LETTERS = "abcdefghijklmnopqrstuvwxyz"

func tJsonRequest(t TB, ctx context.Context, server server.Server, method string, path string, body interface{}, out interface{}) error {
	recorder := httptest.NewRecorder()
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, method, path, bytes.NewReader(bodyBytes))
	if err != nil {
		return err
	}

	server.Router.ServeHTTP(recorder, req)
	if !util.IsHttpStatusOk(recorder.Code) {
		return errors.Errorf(`self-request returned a non-OK status code; code: %v; body: %s`,
			recorder.Code, recorder.Body.Bytes())
	}

	if out == nil {
		return nil
	}

	return json.Unmarshal(recorder.Body.Bytes(), out)
}

func RandomCharSample(str string, count int) string {
	chars := []rune(str)
	buf := make([]rune, count)
	for i := range buf {
		buf[i] = chars[rand.Intn(len(chars))]
	}
	return string(buf)
}

func RandomShortString() string {
	return RandomCharSample(LOWERCASE_LETTERS, 12)
}

func tRequestGetByName(t TB, ctx context.Context, server server.Server, query url.Values, name string, out interface{}) error {
	path := fmt.Sprintf("/%v?%v", name, query.Encode())
	return tJsonRequest(t, ctx, server, http.MethodGet, path, nil, out)
}

func tRequestPostByName(t TB, ctx context.Context, server server.Server, query url.Values, name string, body interface{}, out interface{}) error {
	path := fmt.Sprintf("/%v?%v", name, query.Encode())
	err := tJsonRequest(t, ctx, server, http.MethodPost, path, body, out)
	return err
}

func tRequestDeleteByName(t TB, ctx context.Context, server server.Server, query url.Values, name string, body interface{}, out interface{}) error {
	path := fmt.Sprintf("/%v?%v", name, query.Encode())
	err := tJsonRequest(t, ctx, server, http.MethodDelete, path, body, out)
	return err
}
