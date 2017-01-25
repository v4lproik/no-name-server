package main

import (
	"log"
	"net/http"
	"github.com/emicklei/go-restful"
	"github.com/emicklei/go-restful/swagger"
	"github.com/yinkozi/no-name-domain"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"

	"encoding/json"
)
const dbPath = "foo.db"

type ReportResource struct {
	db *sql.DB
}

type ReportDAO struct {
	Id string
	TypeEnum string

	Form string
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

func (u ReportResource) findReport(request *restful.Request, response *restful.Response) {
	id := request.PathParameter("report-id")

	report := ReadItem(u.db, id)

	if len(report) == 0 {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusNotFound, "404: Report could not be found.")
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

	StoreItem(r.db, *report)

	response.WriteHeaderAndEntity(http.StatusCreated, report)
}

func (r *ReportResource) updateReport(request *restful.Request, response *restful.Response) {
	report := new(domain.Report)
	err := request.ReadEntity(&report)
	if err != nil {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusInternalServerError, err.Error())
		return
	}
	//r.reports[report.Id] = *report
	response.WriteEntity(report)
}

func (r *ReportResource) removeReport(request *restful.Request, response *restful.Response) {
	DeleteItem(r.db, request.PathParameter("report-id"))
}

func main() {
	// initi db
	db := InitDB(dbPath)
	defer db.Close()

	CreateTable(db)

	// to see what happens in the package, uncomment the following
	//restful.TraceLogger(log.New(os.Stdout, "[restful] ", log.LstdFlags|log.Lshortfile))

	wsContainer := restful.NewContainer()
	u := ReportResource{db}
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


func InitDB(filepath string) *sql.DB {
	db, err := sql.Open("sqlite3", filepath)
	if err != nil { panic(err) }
	if db == nil { panic("db nil") }
	return db
}

func StoreItem(db *sql.DB, report domain.Report) {
	sql_additem := `
	INSERT OR REPLACE INTO reports(
		Id,
		TypeEnum,
		InsertedDatetime,
		Form
	) values(?, ?, CURRENT_TIMESTAMP, ?)
	`

	stmt, err := db.Prepare(sql_additem)
	if err != nil { panic(err) }
	defer stmt.Close()

	formToString, err := json.Marshal(report.Form)
	if err != nil { panic(err) }
	defer stmt.Close()

	_, err2 := stmt.Exec(report.Id, report.TypeEnum, formToString)
	if err2 != nil { panic(err2) }
}

func ReadItem(db *sql.DB, id string) []domain.Report {
	sql_readall := `
	SELECT Id, TypeEnum, Form FROM reports Where Id = ?
	`

	rows, err := db.Query(sql_readall, id)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var resultsDAO []ReportDAO
	var results []domain.Report
	for rows.Next() {
		item := ReportDAO{}

		err2 := rows.Scan(&item.Id, &item.TypeEnum, &item.Form)
		if err2 != nil { panic(err2) }

		resultsDAO = append(resultsDAO, item)
	}

	for _,result := range resultsDAO {
		var form domain.Form
		err3 := json.Unmarshal([]byte(result.Form), &form)
		if err3 != nil { panic(err3) }

		results = append(results, *domain.NewReport(result.Id, result.TypeEnum, &form))
	}

	return results
}

func DeleteItem(db *sql.DB, id string) {
	sql := `
	DELETE FROM reports Where Id = ?
	`

	stmt, err := db.Prepare(sql)
	if err != nil { panic(err) }
	defer stmt.Close()

	_, err2 := stmt.Exec(id)
	if err2 != nil { panic(err2) }
}

func UpdateItem(db *sql.DB, report domain.Report) {
	sql := `
	UPDATE reports SET TypeEnum = ?, Form = ? WHERE Id = ?
	`

	stmt, err := db.Prepare(sql)
	if err != nil { panic(err) }
	defer stmt.Close()

	formToString, err := json.Marshal(report.Form)
	if err != nil { panic(err) }
	defer stmt.Close()

	_, err2 := stmt.Exec(report.TypeEnum, formToString, report.Id)
	if err2 != nil { panic(err2) }
}

func ReadAllItems(db *sql.DB) []domain.Report {
	sql_readall := `
	SELECT Id, TypeEnum, Form FROM reports
	ORDER BY datetime(InsertedDatetime) DESC
	`

	rows, err := db.Query(sql_readall)
	if err != nil { panic(err) }
	defer rows.Close()

	var resultsDAO []ReportDAO
	var results []domain.Report
	for rows.Next() {
		item := ReportDAO{}

		err2 := rows.Scan(&item.Id, &item.TypeEnum, &item.Form)
		if err2 != nil { panic(err2) }

		resultsDAO = append(resultsDAO, item)
	}

	for _,result := range resultsDAO {
		var form domain.Form
		err3 := json.Unmarshal([]byte(result.Form), &form)
		if err3 != nil { panic(err3) }

		results = append(results, *domain.NewReport(result.Id, result.TypeEnum, &form))
	}
	return results
}

func CreateTable(db *sql.DB) {
	// create table if not exists
	sql_table := `
	CREATE TABLE IF NOT EXISTS reports(
		Id TEXT NOT NULL PRIMARY KEY,
		TypeEnum TEXT,
		InsertedDatetime DATETIME,
		Form TEXT
	);
	`

	_, err := db.Exec(sql_table)
	if err != nil { panic(err) }
}