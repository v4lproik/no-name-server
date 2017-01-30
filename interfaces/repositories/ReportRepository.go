package repositories

import (
	"github.com/yinkozi/no-name-domain"
	"encoding/json"
)

type DbHandler interface {
	Execute(statement string)
	ExecuteWithParam(statement string, args ...interface{})
	Query(statement string) Row
	QueryWithParam(statement string, args ...interface{}) Row
}

type Row interface {
	Scan(dest ...interface{})
	Next() bool
	Close()
}

type DbRepo struct {
	dbHandlers map[string]DbHandler
	dbHandler  DbHandler
}

type ReportDAO struct {
	Id string
	TypeEnum string

	Form string
}


type DbReportRepo DbRepo

func NewDbReportRepo(dbHandlers map[string]DbHandler) *DbReportRepo {
	dbUserRepo := new(DbReportRepo)
	dbUserRepo.dbHandlers = dbHandlers
	dbUserRepo.dbHandler = dbHandlers["DbReportRepo"]
	return dbUserRepo
}


func (repo *DbReportRepo) Store(report *domain.Report) error{
	sql_additem := `
	INSERT OR REPLACE INTO reports(
		Id,
		TypeEnum,
		InsertedDatetime,
		Form
	) values(?, ?, CURRENT_TIMESTAMP, ?)
	`

	formToString, err := json.Marshal(report.Form)
	if err != nil { return err }

	params := []interface{}{report.Id, report.TypeEnum, formToString}

	repo.dbHandler.ExecuteWithParam(sql_additem, params...)

	//TODO
	return nil
}

func (repo *DbReportRepo) Update(report *domain.Report) error {
	sql := `
	UPDATE reports SET TypeEnum = ?, Form = ? WHERE Id = ?
	`

	formToString, err := json.Marshal(report.Form)
	if err != nil { panic(err) }

	params := []interface{}{report.Id, report.TypeEnum, formToString}

	repo.dbHandler.ExecuteWithParam(sql, params...)

	//TODO
	return nil
}

func (repo *DbReportRepo) Find(id string) (domain.Report, error) {
	newReport := domain.Report{}
	sql_readall := `
	SELECT Id, TypeEnum, Form FROM reports Where Id = ?
	`

	params := []interface{}{id}

	rows := repo.dbHandler.QueryWithParam(sql_readall, params...)
	rows.Next()

	item := ReportDAO{}
	rows.Scan([]interface{}{&item.Id, &item.TypeEnum, &item.Form}...)
	rows.Close()

	var form domain.Form
	err3 := json.Unmarshal([]byte(item.Form), &form)
	if err3 != nil { return newReport, err3 }

	newReport.Id = item.Id
	newReport.Form = &form
	newReport.TypeEnum = item.TypeEnum

	return newReport, nil
}

func (repo *DbReportRepo) Delete(id string) error {
	sql := `
	DELETE FROM reports Where Id = ?
	`

	params := []interface{}{id}

	repo.dbHandler.ExecuteWithParam(sql, params...)

	//TODO
	return nil
}

func (repo *DbReportRepo) List() ([]domain.Report, error) {
	sql_readall := `
	SELECT Id, TypeEnum, Form FROM reports
	ORDER BY datetime(InsertedDatetime) DESC
	`

	rows := repo.dbHandler.Query(sql_readall)

	var resultsDAO []ReportDAO
	var results []domain.Report
	for rows.Next() {
		item := ReportDAO{}
		rows.Scan(&item.Id, &item.TypeEnum, &item.Form)

		resultsDAO = append(resultsDAO, item)
	}
	rows.Close()

	for _,result := range resultsDAO {
		var form domain.Form
		err3 := json.Unmarshal([]byte(result.Form), &form)
		if err3 != nil { panic(err3) }

		results = append(results, *domain.NewReport(result.Id, result.TypeEnum, &form))
	}
	return results, nil
}

func (repo *DbReportRepo) CreateTable() {
	sql_table := `
	CREATE TABLE IF NOT EXISTS reports(
		Id TEXT NOT NULL PRIMARY KEY,
		TypeEnum TEXT,
		InsertedDatetime DATETIME,
		Form TEXT
	);
	`

	repo.dbHandler.Execute(sql_table)
}
