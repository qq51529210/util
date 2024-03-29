package util

import (
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"path"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zht "github.com/go-playground/validator/v10/translations/zh"
)

const indexHTML = "/index.html"

// GinStaticDir 初始化静态文件
func GinStaticDir(r gin.IRouter, relativeRoot, staticsRoot string, staticsDir fs.FS) (err error) {
	return fs.WalkDir(staticsDir, ".", func(p string, d fs.DirEntry, err error) error {
		if d == nil || d.IsDir() {
			return nil
		}
		rp := path.Join(relativeRoot, strings.TrimPrefix(p, staticsRoot))
		r.StaticFileFS(rp, p, http.FS(staticsDir))
		if strings.HasSuffix(p, indexHTML) {
			ginStaticIndex(r, staticsDir, rp, p)
		}
		return nil
	})
}

// 以免 gin 内部对 index.html 一直重定向
func ginStaticIndex(r gin.IRouter, statics fs.FS, path, filePath string) {
	r.GET(path[:len(path)-len(indexHTML)], func(ctx *gin.Context) {
		f, err := statics.Open(filePath)
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

// GinValidateZH 验证器设置为中文
func GinValidateZH(errs map[string]string) ut.Translator {
	if va, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// 设置
		lt := zh.New()
		_ut := ut.New(lt, lt)
		t, _ := _ut.GetTranslator("zh")
		zht.RegisterDefaultTranslations(va, t)
		//
		for k, v := range errs {
			va.RegisterTranslation(k, t, func(ut ut.Translator) error {
				return ut.Add(k, fmt.Sprintf("[{0}]%s", v), true)
			}, func(ut ut.Translator, fe validator.FieldError) string {
				t, _ := ut.T(k, fe.Field())
				return t
			})
		}
		return t
	}
	return nil
}

// GinTranslate 返回翻译的错误数组
func GinTranslate(t ut.Translator, err error) []string {
	var ss []string
	if _err, ok := err.(validator.ValidationErrors); ok {
		for _, v := range _err.Translate(t) {
			ss = append(ss, v)
		}
	} else {
		ss = append(ss, err.Error())
	}
	return ss
}
