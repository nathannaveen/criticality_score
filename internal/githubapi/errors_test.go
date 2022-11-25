package githubapi

import (
	"errors"
	"net/http"
	"testing"

	"github.com/google/go-github/v47/github"
)

func TestErrorResponseStatusCode(t *testing.T) {
	tests := []struct { //nolint:govet
		name string
		err  error
		want int
	}{
		{
			name: "nil error",
			want: 0,
		},
		{
			name: "non-nil error that is not in ErrorResponse",
			err:  errors.New("some error"),
			want: 0,
		},
		{
			name: "error that is in ErrorResponse",
			err: &github.ErrorResponse{
				Response: &http.Response{
					StatusCode: 404,
				},
			},
			want: 404,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := ErrorResponseStatusCode(test.err); got != test.want {
				t.Errorf("ErrorResponseStatusCode() = %v, want %v", got, test.want)
			}
		})
	}
}

func TestGraphQLErrors_Error(t *testing.T) {
	tests := []struct { //nolint:govet
		name      string
		errors    []GraphQLError
		want      string
		wantPanic bool
	}{
		{
			name:      "zero errors",
			errors:    []GraphQLError{},
			wantPanic: true,
		},
		{
			name: "one error",
			errors: []GraphQLError{
				{Message: "one"},
			},
			want: "one",
		},
		{
			name: "more than one error",
			errors: []GraphQLError{
				{Message: "one"},
				{Message: "two"},
				{Message: "three"},
			},
			want: "3 GraphQL errors",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			defer func() {
				r := recover()
				if r != nil && !test.wantPanic {
					t.Fatalf("Error() panic: %v, %s", r, test.name)
				}
			}()

			e := &GraphQLErrors{test.errors}

			if got := e.Error(); got != test.want {
				t.Errorf("Error() = %v, want %v", got, test.want)
			}
		})
	}
}

func TestGraphQLErrors_HasType(t *testing.T) {
	tests := []struct { //nolint:govet
		name   string
		errors []GraphQLError
		t      string
		want   bool
	}{
		{
			name: "type is equal to t",
			errors: []GraphQLError{
				{Type: "NOT_FOUND"},
			},
			t:    "NOT_FOUND",
			want: true,
		},
		{
			name: "type is not equal to t",
			errors: []GraphQLError{
				{Type: "NOT_FOUND"},
			},
			t:    "invalid type",
			want: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			e := &GraphQLErrors{
				errors: test.errors,
			}
			if got := e.HasType(test.t); got != test.want {
				t.Errorf("HasType() = %v, want %v", got, test.want)
			}
		})
	}
}

func TestGraphQLErrors_Is(t *testing.T) {
	tests := []struct { //nolint:govet
		name   string
		errors []GraphQLError
		target error
		want   bool
	}{
		{
			name: "target equals ErrGraphQLNotFound and returns true",
			errors: []GraphQLError{
				{Message: "one", Type: "NOT_FOUND"},
			},
			target: ErrGraphQLNotFound,
			want:   true,
		},
		{
			name: "target equals ErrGraphQLNotFound and returns false",
			errors: []GraphQLError{
				{Message: "one"},
			},
			target: ErrGraphQLNotFound,
			want:   false,
		},
		{
			name: "regular testcase",
			errors: []GraphQLError{
				{Message: "one"},
			},
			target: errors.New("some error"),
			want:   false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			e := &GraphQLErrors{
				errors: test.errors,
			}
			if got := e.Is(test.target); got != test.want {
				t.Errorf("Is() = %v, want %v", got, test.want)
			}
		})
	}
}
