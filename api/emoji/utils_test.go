package main

import (
	"reflect"
	"testing"

	"github.com/pkg/errors"
)

func TestGetEmojis(t *testing.T) {
	tests := []struct {
		name      string
		qs        map[string]string
		wantRes   []string
		wantError error
	}{
		{
			name:      "InvalidSearchKey",
			qs:        map[string]string{"invalid_query_param": "Test"},
			wantRes:   nil,
			wantError: errors.Errorf("%q is a required query string parameter", searchKey),
		},
	}

	for _, tc := range tests {
		gotRes, gotError := getEmojis(tc.qs)
		if !reflect.DeepEqual(tc.wantRes, gotRes) {
			t.Fatalf("%v - expected response: %v, got: %v", tc.name, tc.wantRes, gotRes)
		}

		if tc.wantError.Error() != gotError.Error() {
			t.Fatalf("%v - expected error: %v, got: %v", tc.name, tc.wantError, gotError)
		}
	}
}
