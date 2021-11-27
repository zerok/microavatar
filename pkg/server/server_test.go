package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFallback(t *testing.T) {
	// If no default is requested and also no "nobody" default is defined then
	// a 404 is returned for a missing avatar.
	t.Run("no-fallback", func(t *testing.T) {
		s := New()
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/avatar/something", nil)
		s.ServeHTTP(w, r)
		require.Equal(t, 404, w.Code)
	})

	// If gravatar support is enabled, an unknown avatar request results in a
	// redirection to that.
	t.Run("gravatar", func(t *testing.T) {
		s := New(func(c *Configuration) {
			c.FallbackToGravatar = true
		})
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/avatar/something", nil)
		s.ServeHTTP(w, r)
		require.Equal(t, 307, w.Code)
		require.Equal(t, "https://secure.gravatar.com/avatar/something?s=80", w.Header().Get("Location"))
	})

	// If gravatar support is enabled, an unknown avatar request results in a
	// redirection to that no matter if a default is configured.
	t.Run("gravatar-before-default", func(t *testing.T) {
		s := New(func(c *Configuration) {
			c.FallbackToGravatar = true
			c.DefaultToImage = map[string]string{"nobody": "image.jpg"}
		})
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/avatar/something", nil)
		s.ServeHTTP(w, r)
		require.Equal(t, 307, w.Code)
		require.Equal(t, "https://secure.gravatar.com/avatar/something?s=80", w.Header().Get("Location"))
	})

	// If you pick the 404 fallback, that's what you get.
	t.Run("explicit-fallback", func(t *testing.T) {
		s := New(func(c *Configuration) {
			c.FallbackToGravatar = false
		})
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/avatar/something?default=404", nil)
		s.ServeHTTP(w, r)
		require.Equal(t, 404, w.Code)
	})

	// If a nobody default is set and a non-existing avatar is requested,
	// nobody is returned as redirect.
	t.Run("nobody-fallback", func(t *testing.T) {
		s := New(func(c *Configuration) {
			c.DefaultToImage = map[string]string{"nobody": "somewhere"}
			c.FallbackToGravatar = false
		})
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/avatar/something", nil)
		s.ServeHTTP(w, r)
		require.Equal(t, 307, w.Code)
		require.Equal(t, "/static/defaults/nobody.80.jpg", w.Header().Get("Location"))
	})

	// If an unknown default is requested by a nobody default is available,
	// this should be used instead.
	t.Run("unknown-default-with-nobody", func(t *testing.T) {
		s := New(func(c *Configuration) {
			c.DefaultToImage = map[string]string{"nobody": "somewhere"}
			c.FallbackToGravatar = false
		})
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/avatar/something?d=unknown", nil)
		s.ServeHTTP(w, r)
		require.Equal(t, 307, w.Code)
		require.Equal(t, "/static/defaults/nobody.80.jpg", w.Header().Get("Location"))
	})

	// Without a nobody-default an unknown default results in a 404.
	t.Run("unknown-default-without-nobody", func(t *testing.T) {
		s := New(func(c *Configuration) {
			c.FallbackToGravatar = false
		})
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/avatar/something?d=unknown", nil)
		s.ServeHTTP(w, r)
		require.Equal(t, 404, w.Code)
	})

	t.Run("force-default", func(t *testing.T) {
		s := New(func(c *Configuration) {
			c.FallbackToGravatar = false
			c.EmailToImage = map[string]string{"test@example.org": "exists.jpg"}
			c.DefaultToImage = map[string]string{"nobody": "somewhere"}
		})
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/avatar/0c17bf66e649070167701d2d3cd71711?d=nobody&f=y", nil)
		s.ServeHTTP(w, r)
		require.Equal(t, 307, w.Code)
		require.Equal(t, "/static/defaults/nobody.80.jpg", w.Header().Get("Location"))
	})
}
