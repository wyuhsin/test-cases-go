package tests

import (
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func TestMysqlInsert(t *testing.T) {

	const (
		DATABASE_DRIVER = "mysql"
	)
	info := os.Getenv("MYSQL_CONNECT_INFO")

	db, err := sqlx.Open(DATABASE_DRIVER, info)
	if err != nil {
		t.Fatalf("Failed to connect database: %s\n", err.Error())
	}

	// a := ``
	// b := ``
	query := ``

	_, err = db.NamedExec(query, map[string]any{})
	if err != nil {
		t.Fatalf("Failed to execute command: %s\n", err.Error())
	}

	t.Logf("data inserted successfully\n")

}
