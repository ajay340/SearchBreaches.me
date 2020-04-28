package echofw

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"runtime"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	"github.com/ajay340/SearchBreaches.me/echofw/handler"
)

type TemplateRegistry struct {
	templates *template.Template
}

func (t *TemplateRegistry) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

var errorPath string = ""

func customHTTPErrorHandler(err error, c echo.Context) {
	code := http.StatusInternalServerError
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
	}
	errorPage := fmt.Sprintf(errorPath+"%d.html", code)
	if err := c.File(errorPage); err != nil {
		c.Logger().Error(err)
	}
	c.Logger().Error(err)
}

func StartServer(path string) {
	e := echo.New()
	e.Use(middleware.Secure())
	if runtime.GOOS == "windows" {
		e.Static("/static", path+"\\echofw\\view\\assets")
		errorPath = path + "\\echofw\\view\\html\\ErrorPages\\"
		e.Renderer = &TemplateRegistry{
			templates: template.Must(template.ParseGlob(path + "\\echofw\\view\\html\\*.html")),
		}
	} else {
		e.Static("/static", path+"/echofw/view/assets")
		errorPath = path + "/echofw/view/html/ErrorPages/"
		e.Renderer = &TemplateRegistry{
			templates: template.Must(template.ParseGlob(path + "/echofw/view/html/*.html")),
		}
	}
	e.HTTPErrorHandler = customHTTPErrorHandler

	//Paths
	e.GET("/", handler.IndexHandler)
	e.GET("/search=:search", handler.ListingHandler)
	e.GET("/product/:id", handler.ProductHandler)
	e.GET("/login", handler.LoginGetHandler)
	e.GET("/register", handler.RegisterGetHandler)
	e.GET("/logout", handler.LogoutGetHandler)
	e.GET("/pdf/:id", handler.PDFDownloaderHandler)
	e.GET("/about", handler.AboutHandler)
	e.GET("/searchq=:search", handler.SearchQueryHandler)

	e.GET("/user", handler.UserGetHandler)
	e.POST("/user", handler.UserPostHandler)

	e.POST("/login", handler.LoginPostHandler)
	e.POST("/register", handler.RegisterPostHandler)

	e.Logger.Fatal(e.Start(":80"))
}
