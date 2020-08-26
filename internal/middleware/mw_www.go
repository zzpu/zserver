package middleware

import (
	"os"
	"path/filepath"

	bm "github.com/zzpu/kratos/pkg/net/http/gin"
)

// WWWMiddleware 静态站点中间件
func WWWMiddleware(root string, skippers ...SkipperFunc) bm.HandlerFunc {
	return func(c *bm.Context) {
		if SkipHandler(c, skippers...) {
			c.Next()
			return
		}

		p := c.Request.URL.Path
		fpath := filepath.Join(root, filepath.FromSlash(p))
		_, err := os.Stat(fpath)
		if err != nil && os.IsNotExist(err) {
			fpath = filepath.Join(root, "index.html")
		}

		c.File(fpath)
		c.Abort()
	}
}
