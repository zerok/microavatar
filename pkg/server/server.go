package server

import (
	"context"
	"crypto/md5"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/rs/zerolog"
	"github.com/zerok/microavatar/pkg/resizer"
)

type Configuration struct {
	EmailToImage       map[string]string
	DefaultToImage     map[string]string
	CacheFolder        string
	Logger             zerolog.Logger
	FallbackToGravatar bool
	Resizer            resizer.Resizer
}

type Configurator func(c *Configuration)

func New(configs ...Configurator) *server {
	c := Configuration{
		CacheFolder: "cache",
	}
	for _, cfg := range configs {
		cfg(&c)
	}
	h2i := make(map[string]string)
	for e, i := range c.EmailToImage {
		h2i[generateMD5(e)] = i
	}
	s := &server{
		hashToImage: h2i,
		cfg:         c,
	}
	r := chi.NewRouter()
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := c.Logger.WithContext(r.Context())
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	})
	r.Get("/avatar/{emailhash}", s.handleAvatar)
	s.r = r
	return s
}

type server struct {
	r           chi.Router
	hashToImage map[string]string
	cfg         Configuration
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.r.ServeHTTP(w, r)
}

func (s *server) handleAvatar(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := zerolog.Ctx(ctx)
	emailhash := chi.URLParam(r, "emailhash")
	if emailhash == "" {
		http.Error(w, "No hash specified", http.StatusBadRequest)
		return
	}
	size := intFromQuery(r, []string{"size", "s"}, 80)
	if size > 512 || size < 1 {
		http.Error(w, "Invalid size", http.StatusBadRequest)
		return
	}
	defaultImage := strFromQuery(r, []string{"default", "d"}, "")
	forceDefault := boolFromQuery(r, []string{"forcedefault", "f"}, false)
	image, ok := s.hashToImage[emailhash]
	var cachePath string
	if !ok || forceDefault {
		if !forceDefault && s.cfg.FallbackToGravatar {
			http.Redirect(w, r, fmt.Sprintf("https://secure.gravatar.com/avatar/%s?s=%d", emailhash, size), http.StatusTemporaryRedirect)
			return
		}
		if defaultImage != "" {
			switch defaultImage {
			case "404":
				http.Error(w, "Not found", http.StatusNotFound)
				return
			default:
				image, defaultImage = s.getDefaultImagePath(defaultImage)
			}
		} else {
			image, defaultImage = s.getDefaultImagePath("nobody")
		}
		if image == "" {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}
		cachePath = filepath.Join(s.cfg.CacheFolder, fmt.Sprintf("_default_%s.s%d.jpg", defaultImage, size))
	} else {
		cachePath = filepath.Join(s.cfg.CacheFolder, fmt.Sprintf("%s.s%d.jpg", emailhash, size))
	}
	if err := s.createIfMissing(ctx, cachePath, image, size); err != nil {
		logger.Error().Err(err).Msg("Failed to generate image")
		http.Error(w, "Failed to generate image", http.StatusInternalServerError)
		return
	}

	http.ServeFile(w, r, cachePath)
}

func (s *server) getDefaultImagePath(defaultImageName string) (string, string) {
	nobody := s.cfg.DefaultToImage["nobody"]
	requested, hasRequested := s.cfg.DefaultToImage[defaultImageName]
	if hasRequested {
		return requested, defaultImageName
	}
	return nobody, "nobody"
}

func (s *server) createIfMissing(ctx context.Context, output string, input string, size int64) error {
	inst, err := os.Stat(input)
	if err != nil {
		return err
	}
	outst, err := os.Stat(output)
	if !os.IsNotExist(err) {
		return err
	}
	if outst != nil && outst.ModTime().After(inst.ModTime()) {
		return nil
	}
	if s.cfg.Resizer == nil {
		return fmt.Errorf("no resizer available")
	}
	return s.cfg.Resizer.Resize(ctx, output, input, size, size)
}

func intFromQuery(r *http.Request, paramNames []string, defaultValue int64) int64 {
	val := strFromQuery(r, paramNames, "")
	if val == "" {
		return defaultValue
	}
	i, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return defaultValue
	}
	return i
}

func strFromQuery(r *http.Request, paramNames []string, defaultValue string) string {
	for _, p := range paramNames {
		val := r.URL.Query().Get(p)
		if val != "" {
			return val
		}
	}
	return defaultValue
}

func boolFromQuery(r *http.Request, paramNames []string, defaultValue bool) bool {
	for _, p := range paramNames {
		val := r.URL.Query().Get(p)
		if val != "" {
			return val == "y"
		}
	}
	return defaultValue
}

func generateMD5(s string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))
}
