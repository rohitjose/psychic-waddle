package main

import (
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

func TestHandler(t *testing.T) {
	tests := []struct {
		name       string
		emojiRes   string
		searchKey  string
		wantRes    string
		wantError  error
		wantStatus int
	}{
		{
			name:       "Success",
			emojiRes:   "[ { \"emoji\": \"👌\", \"aliases\": [\"dummy-emoji\"] } ]",
			searchKey:  "dummy-emoji",
			wantRes:    "👌",
			wantError:  nil,
			wantStatus: http.StatusOK,
		},
		{
			name:       "FailEmojiDoesNotExist",
			emojiRes:   "[ { \"emoji\": \"👌\", \"aliases\": [\"dummy-emoji\"] } ]",
			searchKey:  "emoji-that-does-not-exist",
			wantRes:    "no results for \"emoji-that-does-not-exist\"",
			wantError:  nil,
			wantStatus: http.StatusBadRequest,
		},
	}

	for cn, tc := range tests {
		defer teardown(setup(t, tc.emojiRes, fmt.Sprintf("127.0.0.1:929%v", cn)))

		gotRes, gotError := handler(events.APIGatewayProxyRequest{QueryStringParameters: map[string]string{searchKey: tc.searchKey}})

		// Verify response content
		if v := gotRes.Body; v != tc.wantRes {
			t.Fatalf("%v - expected response: %v, got: %v", tc.name, tc.wantRes, v)
		}

		// Verify response status code
		if gotRes.StatusCode != tc.wantStatus {
			t.Fatalf("%v - expected response status: %v, got: %v", tc.name, tc.wantStatus, gotRes.StatusCode)
		}

		// Verify error
		if tc.wantError != gotError {
			t.Fatalf("%v - expected error: %v, got: %v", tc.name, tc.wantError, gotError)
		}

	}
}

func TestHandlerFailNoRouteToHost(t *testing.T) {
	res, _ := handler(events.APIGatewayProxyRequest{QueryStringParameters: map[string]string{searchKey: "dummy-emoji"}})
	if res.StatusCode == http.StatusOK {
		t.Fatalf("TestHandlerFailNoRouteToHost failed: expected an error")
	}

	x := "could not retrieve emoji data: could not make the request: Get \"http://127.0.0.1:9291\": dial tcp 127.0.0.1:9291: connect: connection refused"
	if v := res.Body; v != x {
		t.Fatalf("TestHandlerFailNoRouteToHost failed: have %q, want %q", v, x)
	}
}

func setup(t *testing.T, body string, mockServerURL string) *httptest.Server {
	os.Setenv("SOURCE_URL", "http://"+mockServerURL)

	ts := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { fmt.Fprintln(w, body) }))

	l, err := net.Listen("tcp", mockServerURL)
	if err != nil {
		t.Fatalf("could not start the mock server: %v", err)
	}

	ts.Listener = l
	ts.Start()

	return ts
}

func teardown(ts *httptest.Server) {
	ts.Close()
}
