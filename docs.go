// +build doc

package main

import (
    "fmt"
    "log/slog"
    "strings"
    "text/template"
    "bytes"
    "os"
    "io/ioutil"
    "path/filepath"
    "net/http"
    "github.com/gin-gonic/gin"
    ginSwagger "github.com/swaggo/gin-swagger"
    "github.com/swaggo/files"
    _ "suglider-auth/docs"
    "suglider-auth/configs"
)

type specUrl struct {
    Url string
}

func redocHandler(c *gin.Context) {
    if c.Request.Method != http.MethodGet {
        c.AbortWithStatus(http.StatusMethodNotAllowed)
        return
    }
    // matcher := regexp.MustCompile(`(.*)(index\.html|doc\.json|redoc\.standalone\.js)[?|.]*`)
    // matches := matcher.FindStringSubmatch(c.Request.RequestURI)
    // if len(matches) != 3 {
    //     c.String(http.StatusNotFound, http.StatusText(http.StatusNotFound))
    //     return
    // }
    // filename := matches[2]
    idx := strings.LastIndex(c.Request.RequestURI, "/")
    if idx <= 1 {
        c.String(http.StatusNotFound, http.StatusText(http.StatusNotFound))
        return
    }
    filename := c.Request.RequestURI[idx + 1:]
    switch filepath.Ext(filename) {
    case ".html":
        c.Header("Content-Type", "text/html; charset=utf-8")
    case ".js":
        c.Header("Content-Type", "application/javascript")
    case ".json":
        c.Header("Content-Type", "application/json; charset=utf-8")
    }
    switch filename {
    case "index.html":
        indexHtml, err := ioutil.ReadFile("configs/redoc.tmpl")
        if err != nil {
            errorMessage := fmt.Sprintf("Read HTML Template File Error: %v", err)
            slog.Error(errorMessage)
            
            c.AbortWithStatus(http.StatusInternalServerError)
            return
        }
        bodyTemplate, err := template.New("redoc").Parse(string(indexHtml))
        if err != nil {
            slog.Error(err.Error())
            c.AbortWithStatus(http.StatusInternalServerError)
            return
        }
        specUrl := specUrl {
            Url: fmt.Sprintf("http://127.0.0.1:%d%sswagger/doc.json", configs.Args.Port, configs.Args.Subpath),
        }
        buf := new(bytes.Buffer)
        if err = bodyTemplate.Execute(buf, specUrl); err != nil {
            slog.Error(err.Error())
             c.AbortWithStatus(http.StatusInternalServerError)
             return
        }
        _, _ = c.Writer.Write(buf.Bytes())
    case "doc.json":
        doc, err := ioutil.ReadFile("docs/swagger.json")
        if err != nil {
            slog.Error(err.Error())
            c.AbortWithStatus(http.StatusInternalServerError)
            return
        }
        c.String(http.StatusOK, string(doc))
    default:
        c.AbortWithStatus(http.StatusNotFound)
    }
}

func DisableCustomHandker(hdlf gin.HandlerFunc, env string) gin.HandlerFunc {
  if os.Getenv(env) != "" {
    return func(c *gin.Context) {
        c.String(http.StatusNotFound, "")
    }
  }
  return hdlf
}

func init() {
    // set SWAGGER_OFF env to disable swagger.
    if configs.ApplicationConfig.Swagger.Theme == "redoc" {
        swag = DisableCustomHandker(redocHandler, "SWAGGER_OFF")
    } else {
        swag = ginSwagger.DisablingWrapHandler(swaggerFiles.Handler, "SWAGGER_OFF")
    }
}
