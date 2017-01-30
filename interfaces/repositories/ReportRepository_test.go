package repositories_test

import (
	"github.com/yinkozi/no-name-server/infrastructure"
	"testing"
	"github.com/yinkozi/no-name-server/interfaces/repositories"
	"github.com/yinkozi/no-name-domain"
	"os"
)

var repo *repositories.DbReportRepo

func TestMain(m *testing.M) {
	dbHandler := infrastructure.NewSqliteHandler("foo_test.db")

	handlers := make(map[string]repositories.DbHandler)
	handlers["DbReportRepo"] = dbHandler
	repo = repositories.NewDbReportRepo(handlers)

	repo.CreateTable();

	os.Exit(m.Run())
}

func Test_create_a_report(t *testing.T) {
	//given
	clean()
	r := domain.Report{}
	r.Id = "555"
	r.TypeEnum = "TEST"

	//when
	repo.Store(&r)

	//then
	report, _ := repo.Find("555")
	if report.Id != "555" {
		t.Error()
	}
	clean()
}

func Test_delete_a_report(t *testing.T) {
	//given
	clean()
	r := domain.Report{}
	r.Id = "556"
	repo.Store(&r)

	//when
	repo.Delete(r.Id)

	//then
	report, _ := repo.Find("556")
	if report.Id == "556" {
		t.Error()
	}
	clean()
}

func Test_list_all_reports(t *testing.T) {
	//given
	clean()
	r := domain.Report{}
	r.Id = "556"

	r2 := domain.Report{}
	r2.Id = "557"

	repo.Store(&r)
	repo.Store(&r2)

	//when
	reports, _ := repo.List()

	//then
	if reports[0].Id != "556" {
		t.Error()
	}
	if reports[1].Id != "557" {
		t.Error()
	}
	clean()
}

func clean(){
	reports, _ := repo.List()
	for _, report := range (reports) {
		repo.Delete(report.Id)
	}
}