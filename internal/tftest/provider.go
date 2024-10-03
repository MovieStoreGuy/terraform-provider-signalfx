package tftest

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/signalfx/signalfx-go"

	"github.com/splunk-terraform/terraform-provider-signalfx/internal/definition/provider/providerhelper"
)

// NewTestHTTPMockState allows for a state object to be configured as if it would be set by the provider.
// The underlying mux uses the [http.ServerMux] using `HandlerFunc` method which allows setting the expected
// method on routes, for example:
//
//	routes := map[string]http.HandlerFunc{
//		"GET /v2/<resource>": func(...){ ... } // Only accepts GET method request matching the route
//		"/v2/login": func(...){ ... }          // Handles all request for `/v2/login` regardless of method
//	}
func NewTestHTTPMockState(routes map[string]http.HandlerFunc) func(testing.TB) any {
	mux := http.NewServeMux()
	for pattern, handler := range routes {
		mux.HandleFunc(pattern, handler)
	}

	return func(t testing.TB) any {
		secret := t.Name()
		s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if v := r.Header.Get("X-SF-Token"); v != secret {
				http.Error(w, "Incorrectly authenticated", http.StatusForbidden)
				return
			}
			mux.ServeHTTP(w, r)
		}))
		t.Cleanup(s.Close)

		sfx, _ := signalfx.NewClient(
			secret,
			signalfx.HTTPClient(s.Client()),
			signalfx.APIUrl(s.URL),
		)

		return &providerhelper.State{
			APIURL:       s.URL,
			CustomAppURL: s.URL,
			Client:       sfx,
		}
	}
}
