package main

import (
	"net/http"

	"github.com/spf13/cobra"
	"github.com/zerok/microavatar/pkg/resizer"
	"github.com/zerok/microavatar/pkg/server"
)

var addr string
var cacheFolder string
var gravatarFallback bool
var emails map[string]string
var defaults map[string]string

var serveCmd = &cobra.Command{
	Use: "serve",
	RunE: func(cmd *cobra.Command, args []string) error {
		s := server.New(func(c *server.Configuration) {
			c.EmailToImage = emails
			c.DefaultToImage = defaults
			c.Logger = logger
			c.FallbackToGravatar = gravatarFallback
			c.Resizer = resizer.NewImageMagick()
			c.CacheFolder = cacheFolder
		})
		hs := http.Server{}
		hs.Handler = s
		hs.Addr = addr
		logger.Info().Msgf("Listening on %s", addr)
		return hs.ListenAndServe()
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
	serveCmd.Flags().StringVar(&addr, "addr", "localhost:8888", "Address to listen on")
	serveCmd.Flags().StringVar(&cacheFolder, "cache-folder", "cache", "Path to the cache folder")
	serveCmd.Flags().StringToStringVar(&emails, "email", make(map[string]string), "email=image mapping(s)")
	serveCmd.Flags().StringToStringVar(&defaults, "default", make(map[string]string), "nobody=./static.png for instance")
	serveCmd.Flags().BoolVar(&gravatarFallback, "gravatar", false, "Fall back to gravatar")
}
