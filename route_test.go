package deer

import (
	"testing"
)

func Test_trPathPatternToRegexpPattern(t *testing.T) {
	tests := []struct {
		name string
		path string
		want string
	}{
		{name: "/users", path: "/users", want: "^/users$"},
		{name: "/users/:uid", path: "/users/:uid", want: "^/users/(?P<uid>[^/]+)$"},
		{name: "/orgs/:oid/users/:uid", path: "/orgs/:oid/users/:uid", want: "^/orgs/(?P<oid>[^/]+)/users/(?P<uid>[^/]+)$"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := trPathPatternToRegexpPattern(tt.path); got != tt.want {
				t.Errorf("trPathPatternToRegexpPattern() = %v, want %v", got, tt.want)
			}
		})
	}
}
