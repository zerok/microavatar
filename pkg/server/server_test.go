package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFallback(t *testing.T) {
	t.Run("no-fallback", func(t *testing.T) {
		s := New()
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/avatar/something", nil)
		s.ServeHTTP(w, r)
		require.Equal(t, 404, w.Code)
	})
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
	t.Run("skip-fallback", func(t *testing.T) {
		s := New(func(c *Configuration) {
			c.FallbackToGravatar = true
		})
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/avatar/something?default=404", nil)
		s.ServeHTTP(w, r)
		require.Equal(t, 404, w.Code)
	})
}
