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
		swaggerui.WithSpecURLs("AdminCouch",
			[]swaggerui.SpecURL{
				{
					Name: "AdminCouch",
					URL:  path.Join(swaggerPath, "admin_couch.swagger.json"),
				},
				{
					Name: "AdminSubscriptions",
					URL:  path.Join(swaggerPath, "admin_subscriptions.swagger.json"),
				},
				{
					Name: "AdminUsers",
					URL:  path.Join(swaggerPath, "admin_users.swagger.json"),
				},
				{
					Name: "ArtelApi",
					URL:  path.Join(swaggerPath, "artel_api.swagger.json"),
				},
				{
					Name: "Auth",
					URL:  path.Join(swaggerPath, "auth.swagger.json"),
				},
				{
					Name: "CouchInstances",
					URL:  path.Join(swaggerPath, "couch_instances.swagger.json"),
				},
				{
					Name: "DockerHosts",
					URL:  path.Join(swaggerPath, "docker_hosts.swagger.json"),
				},
				{
					Name: "ExternalConnections",
					URL:  path.Join(swaggerPath, "external_connections.swagger.json"),
				},
				{
					Name: "McpKeys",
					URL:  path.Join(swaggerPath, "mcp_keys.swagger.json"),
				},
				{
					Name: "Notes",
					URL:  path.Join(swaggerPath, "notes.swagger.json"),
				},
				{
					Name: "Prompts",
					URL:  path.Join(swaggerPath, "prompts.swagger.json"),
				},
				{
					Name: "S3Instances",
					URL:  path.Join(swaggerPath, "s3_instances.swagger.json"),
				},
				{
					Name: "TaskTrackers",
					URL:  path.Join(swaggerPath, "task_trackers.swagger.json"),
				},
				{
					Name: "Tracts",
					URL:  path.Join(swaggerPath, "tracts.swagger.json"),
				},
				{
					Name: "UserErrors",
					URL:  path.Join(swaggerPath, "user_errors.swagger.json"),
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
