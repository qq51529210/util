package util

import (
	"io"
	"io/fs"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// InitGinStaticDir 初始化静态文件
func InitGinStaticDir(r gin.IRouter, dir fs.FS) (err error) {
	index := "/index.html"
	return fs.WalkDir(dir, ".", func(path string, d fs.DirEntry, err error) error {
		if d == nil || d.IsDir() {
			return nil
		}
		r.StaticFileFS(path, path, http.FS(dir))
		if strings.HasSuffix(path, index) {
			initGinStaticIndex(r, dir, path, index)
		}
		return nil
	})
}

// 以免 gin 内部对 index.html 一直重定向
func initGinStaticIndex(r gin.IRouter, statics fs.FS, path, index string) {
	r.GET(path[:len(path)-len(index)], func(ctx *gin.Context) {
		f, err := statics.Open(path)
		if err != nil {
			ctx.Status(http.StatusNotFound)
			return
		}
		//
		ctx.Writer.Header().Set("Content-Type", gin.MIMEHTML)
		//
		io.Copy(ctx.Writer, f)
	})
}
