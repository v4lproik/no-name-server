package main

import (
	"log"
	"net/http"
	"github.com/emicklei/go-restful"
	"github.com/emicklei/go-restful/swagger"
	"github.com/yinkozi/no-name-domain"
	"strconv"
	"fmt"
)

type ReportResource struct {
	// normally one would use DAO (data access object)
	reports map[string]domain.Report
}

func (r ReportResource) Register(container *restful.Container) {
	ws := new(restful.WebService)
	ws.
	Path("/no-name").
		Doc("Manage Report").
		Consumes(restful.MIME_XML, restful.MIME_JSON).
		Produces(restful.MIME_JSON, restful.MIME_XML) // you can specify this per route as well

	ws.Route(ws.GET("/{report-id}").To(r.findReport).
	// docs
		Doc("get a report").
		Operation("findReport").
		Param(ws.PathParameter("report-id", "identifier of the report").DataType("string")).
		Writes(domain.Report{})) // on the response

	ws.Route(ws.PUT("/{report-id}").To(r.updateReport).
	// docs
		Doc("update a user").
		Operation("updateUser").
		Param(ws.PathParameter("report-id", "identifier of the report").DataType("string")).
		ReturnsError(409, "duplicate report-id", nil).
		Reads(domain.Report{})) // from the request

	ws.Route(ws.POST("/{report-id}").To(r.createReport).
	// docs
		Doc("create a report").
		Operation("createReport").
		Reads(domain.Report{})) // from the request

	ws.Route(ws.DELETE("/{report-id}").To(r.removeReport).
	// docs
		Doc("delete a report").
		Operation("removeReport").
		Param(ws.PathParameter("report-id", "identifier of the report").DataType("string")))

	container.Add(ws)
}

func (u ReportResource) findReport(request *restful.Request, response *restful.Response) {
	id := request.PathParameter("report-id")
	report := u.reports[id]
	if len(report.Id) == 0 {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusNotFound, "404: User could not be found.")
		return
	}
	response.WriteEntity(report)
}

func (r *ReportResource) createReport(request *restful.Request, response *restful.Response) {
	report := new(domain.Report)
	err := request.ReadEntity(report)
	if err != nil {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusInternalServerError, err.Error())
		return
	}
	report.Id = strconv.Itoa(len(r.reports) + 1) // simple id generation
	r.reports[report.Id] = *report
	response.WriteHeaderAndEntity(http.StatusCreated, report)

	fmt.Printf("%s", report)
}

func (r *ReportResource) updateReport(request *restful.Request, response *restful.Response) {
	report := new(domain.Report)
	err := request.ReadEntity(&report)
	if err != nil {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusInternalServerError, err.Error())
		return
	}
	r.reports[report.Id] = *report
	response.WriteEntity(report)
}

func (r *ReportResource) removeReport(request *restful.Request, response *restful.Response) {
	id := request.PathParameter("report-id")
	delete(r.reports, id)
}

func main() {
	// to see what happens in the package, uncomment the following
	//restful.TraceLogger(log.New(os.Stdout, "[restful] ", log.LstdFlags|log.Lshortfile))

	wsContainer := restful.NewContainer()
	u := ReportResource{map[string]domain.Report{}}
	u.Register(wsContainer)


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