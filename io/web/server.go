package web

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gnames/gnfmt"
	"github.com/gnames/gnparser"
	"github.com/gnames/gnparser/ent/parsed"
	"github.com/gnames/gnparser/io/fs"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const withLogs = false

type inputPOST struct {
	Names       []string `json:"names"`
	WithDetails bool     `json:"withDetails,omitempty"`
	CSV         bool     `json:"csv,omitempty"`
}

// Run starts the GNparser web service and servies both RESTful API and
// a website.
func Run(gnps GNparserService) {
	e := echo.New()
	e.Renderer = templates()
	e.Use(middleware.Gzip())
	e.Use(middleware.CORS())
	if withLogs {
		e.Use(middleware.Logger())
	}
	e.GET("/", home(gnps))
	e.GET("/doc/api", docAPI())
	e.GET("/api", info())
	e.GET("/api/v1", info())
	e.GET("/api/v1/ping", ping(gnps))
	e.GET("/api/v1/version", ver(gnps))
	e.GET("/api/v1/:names", parseNamesGET(gnps))
	e.GET("/api/:names", parseNamesGET(gnps))
	e.POST("/api/v1", parseNamesPOST(gnps))
	e.POST("/api", parseNamesPOST(gnps))

	assetHandler := http.FileServer(fs.Files)
	e.GET("/static/*", echo.WrapHandler(http.StripPrefix("/static/", assetHandler)))

	addr := fmt.Sprintf(":%d", gnps.Port())
	s := &http.Server{
		Addr:         addr,
		ReadTimeout:  5 * time.Minute,
		WriteTimeout: 5 * time.Minute,
	}
	e.Logger.Fatal(e.StartServer(s))
}

func info() func(c echo.Context) error {
	return func(c echo.Context) error {
		return c.String(
			http.StatusOK,
			`OpenAPI for gnparser is described at

https://app.swaggerhub.com/apis-docs/dimus/gnparser/1.0.0`,
		)
	}
}

func ping(gnps GNparserService) func(echo.Context) error {
	return func(c echo.Context) error {
		result := gnps.Ping()
		return c.String(http.StatusOK, result)
	}
}

func ver(gnps GNparserService) func(echo.Context) error {
	return func(c echo.Context) error {
		version := gnps.GetVersion()
		return c.JSON(http.StatusOK, version)
	}
}

func parseNamesGET(gnps GNparserService) func(echo.Context) error {
	return func(c echo.Context) error {
		nameStr, _ := url.QueryUnescape(c.Param("names"))
		csv := c.QueryParam("csv") == "true"
		det := c.QueryParam("with_details") == "true"
		gnp := gnps.ChangeConfig(opts(c, csv, det)...)
		names := strings.Split(nameStr, "|")
		res := gnp.ParseNames(names)
		return formatNames(c, res, gnp.Format())
	}
}

func parseNamesPOST(gnps GNparserService) func(echo.Context) error {
	return func(c echo.Context) error {
		var input inputPOST
		if err := c.Bind(&input); err != nil {
			return err
		}
		gnp := gnps.ChangeConfig(opts(c, input.CSV, input.WithDetails)...)
		res := gnp.ParseNames(input.Names)
		return formatNames(c, res, gnp.Format())
	}
}

func formatNames(
	c echo.Context,
	res []parsed.Parsed,
	f gnfmt.Format,
) error {
	switch f {
	case gnfmt.CSV:
		resCSV := make([]string, 0, len(res)+1)
		resCSV = append(resCSV, parsed.HeaderCSV())
		for i := range res {
			resCSV = append(resCSV, res[i].Output(f))
		}
		return c.String(http.StatusOK, strings.Join(resCSV, "\n"))
	default:
		return c.JSON(http.StatusOK, res)
	}
}

func opts(c echo.Context, csv, details bool) []gnparser.Option {
	if csv {
		return []gnparser.Option{gnparser.OptFormat("csv")}
	}
	var res []gnparser.Option
	if details {
		res = []gnparser.Option{gnparser.OptWithDetails(true)}
	}
	return res
}
