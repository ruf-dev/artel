package docs

import (
    "embed"
    "io/fs"
    "log"
    "net/http"
    "path"

    swaggerui "github.com/Red-Sock/go-swagger-ui"
)

//go:embed all:swaggers
var swaggers embed.FS

const (
    BasePath    = "/docs/"
    swaggerPath = BasePath + "swaggers/"
)

func Swagger() (p string, handler http.HandlerFunc) {
    mux := http.NewServeMux()

    mux.Handle(BasePath, swaggerui.NewHandler(
        swaggerui.WithBasePath(BasePath),
        swaggerui.WithHTMLTitle("Swagger"),
        swaggerui.WithSpecURLs("ArtelApi",
            []swaggerui.SpecURL{ 
            {
                Name: "ArtelApi",
                URL:  path.Join(swaggerPath, "artel_api.swagger.json"),
            },
            {
                Name: "Users",
                URL:  path.Join(swaggerPath, "users.swagger.json"),
            },
            {
                Name: "Vaults",
                URL:  path.Join(swaggerPath, "vaults.swagger.json"),
            },
        }),
        swaggerui.WithShowExtensions(true),
    ))

    stripped, err := fs.Sub(swaggers, "swaggers")
    if err != nil {
        log.Fatal(err)
    }

    ffs := http.StripPrefix(swaggerPath, http.FileServer(http.FS(stripped)))
    mux.Handle(swaggerPath, ffs)

    return BasePath, mux.ServeHTTP
}