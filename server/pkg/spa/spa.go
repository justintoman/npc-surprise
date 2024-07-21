package spa

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
)

func Middleware(urlPrefix, spaDirectory string) gin.HandlerFunc {
	directory := static.LocalFile(spaDirectory, true)
	fileserver := http.FileServer(directory)
	if urlPrefix != "" {
		fileserver = http.StripPrefix(urlPrefix, fileserver)
	}
	return func(c *gin.Context) {
		if directory.Exists(urlPrefix, c.Request.URL.Path) {
			slog.Info("servering SPA", "path", c.Request.URL.Path)
			fileserver.ServeHTTP(c.Writer, c.Request)
			c.Abort()
		} else if len(c.HandlerNames()) == 3 {
			c.Request.URL.Path = "/"
			fileserver.ServeHTTP(c.Writer, c.Request)
			c.Abort()
		} else {
			slog.Info("handler names", "names", c.HandlerNames())
			slog.Info("should be an API route", "path", c.Request.URL.Path)
			c.Next()
		}
	}
}
