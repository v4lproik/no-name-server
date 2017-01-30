package main

import (
	"log"
	"net/http"
	"github.com/emicklei/go-restful"
	"github.com/emicklei/go-restful/swagger"
	_ "github.com/mattn/go-sqlite3"

	"github.com/yinkozi/no-name-server/infrastructure"
	"github.com/yinkozi/no-name-server/interfaces/repositories"
	"github.com/yinkozi/no-name-server/usecases"
	"github.com/yinkozi/no-name-server/interfaces/webcontrollers"
)
const dbPath = "foo.db"

func main() {
	// init db
	dbHandler := infrastructure.NewSqliteHandler(dbPath)

	// init repositories
	handlers := make(map[string] repositories.DbHandler)
	handlers["DbReportRepo"] = dbHandler

	// init interactors
	appInteractor := new(usecases.AppInteractor)
	appInteractor.ReportRepository = repositories.NewDbReportRepo(handlers)

	// init web resources
	webServiceHandler := webcontrollers.WebServiceHandler{}
	webServiceHandler.AppInteractor = appInteractor

	//init web services containers (specific to the library)
	// to see what happens in the package, uncomment the following
	//restful.TraceLogger(log.New(os.Stdout, "[restful] ", log.LstdFlags|log.Lshortfile))
	wsContainer := restful.NewContainer()

	// inject containers into web resources so they can be accessible
	webServiceHandler.Register(wsContainer)


	// Optionally, you can install the Swagger Service which provides a nice Web UI on your REST API
	// You need to download the Swagger HTML5 assets and change the FilePath location in the config below.
	// Open http://localhost:8080/apidocs and enter http://localhost:8080/apidocs.json in the api input field.
	config := swagger.Config{
		WebServices:    wsContainer.RegisteredWebServices(), // you control what services are visible
		WebServicesUrl: "http://localhost:8080",
		ApiPath:        "/apidocs.json",

		// Optionally, specifiy where the UI is located
		SwaggerPath:     "/apidocs/",
		SwaggerFilePath: "/Users/emicklei/xProjects/swagger-ui/dist"}
	swagger.RegisterSwaggerService(config, wsContainer)

	log.Printf("start listening on localhost:8080")
	server := &http.Server{Addr: ":8080", Handler: wsContainer}
	log.Fatal(server.ListenAndServe())
}