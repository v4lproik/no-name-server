package repositories

import (
	"github.com/yinkozi/no-name-domain"
	"encoding/json"
)

type DbHandler interface {
	Execute(statement string) error
	ExecuteWithParam(statement string, args ...interface{}) error
	Query(statement string) (Row, error)
	QueryWithParam(statement string, args ...interface{}) (Row, error)
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

	return repo.dbHandler.ExecuteWithParam(sql_additem, params...)
}

func (repo *DbReportRepo) Update(report *domain.Report) error {
	sql := `
	UPDATE reports SET TypeEnum = ?, Form = ? WHERE Id = ?
	`

	formToString, err := json.Marshal(report.Form)
	if err != nil { return err }

	params := []interface{}{report.Id, report.TypeEnum, formToString}

	return repo.dbHandler.ExecuteWithParam(sql, params...)
}

func (repo *DbReportRepo) Find(id string) (domain.Report, error) {
	newReport := domain.Report{}
	sql_readall := `
	SELECT Id, TypeEnum, Form FROM reports Where Id = ?
	`

	params := []interface{}{id}

	rows, err := repo.dbHandler.QueryWithParam(sql_readall, params...)
	if err != nil { return newReport, err }
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

	return repo.dbHandler.ExecuteWithParam(sql, params...)
}

func (repo *DbReportRepo) List() ([]domain.Report, error) {
	sql_readall := `
	SELECT Id, TypeEnum, Form FROM reports
	ORDER BY datetime(InsertedDatetime) DESC
	`

	rows, err := repo.dbHandler.Query(sql_readall)
	if err != nil { return nil, err }

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
		if err3 != nil { return nil, err3 }

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
