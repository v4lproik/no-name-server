package webcontrollers

import (
	"github.com/emicklei/go-restful"
	"github.com/yinkozi/no-name-domain"
	"net/http"
)


type AppInteractor interface {
	FindReport(id string) (domain.Report, error)
	CreateReport(report *domain.Report) error
	UpdateReport(report *domain.Report) error
	DeleteReport(id string) error
}

type WebServiceHandler struct {
	AppInteractor AppInteractor
}

func (r WebServiceHandler) Register(container *restful.Container) {
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
		Doc("update a report").
		Operation("updateReport").
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

func (handler WebServiceHandler) findReport(request *restful.Request, response *restful.Response) {
	id := request.PathParameter("report-id")

	report, _ := handler.AppInteractor.FindReport(id)

	if report.Id == "" {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusNotFound, "404: Report could not be found.")
		return
	}
	response.WriteEntity(report)
}

func (handler WebServiceHandler) createReport(request *restful.Request, response *restful.Response) {
	report := new(domain.Report)
	report.Id = request.PathParameter("report-id")
	err := request.ReadEntity(report)
	if err != nil {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusInternalServerError, err.Error())
		return
	}

	err = handler.AppInteractor.CreateReport(report)
	if err != nil {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusInternalServerError, err.Error())
		return
	}

	response.WriteHeaderAndEntity(http.StatusCreated, report)
}

func (handler WebServiceHandler) updateReport(request *restful.Request, response *restful.Response) {
	report := new(domain.Report)
	err := request.ReadEntity(&report)
	if err != nil {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusInternalServerError, err.Error())
		return
	}

	err = handler.AppInteractor.UpdateReport(report)
	if err != nil {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusInternalServerError, err.Error())
		return
	}

	response.WriteEntity(report)
}

func (handler WebServiceHandler) removeReport(request *restful.Request, response *restful.Response) {
	err := handler.AppInteractor.DeleteReport(request.PathParameter("report-id"))
	if err != nil {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusInternalServerError, err.Error())
		return
	}
}
