package dbhandle

import (
	"database/sql"
	"fmt"
	"strings"
)

type DbHandler struct {
	db *sql.DB
}

func GetDbHandler(dbUser, dbPassword, dbHost, dbPort, dbName string) (*DbHandler, error) {

	db, err := sql.Open("mysql",
		fmt.Sprintf(
			"%s:%s@tcp(%s:%s)/%s",
			dbUser,
			dbPassword,
			dbHost,
			dbPort,
			dbName,
		))

	return &DbHandler{
		db: db,
	}, err
}
func (dh *DbHandler) Close() error {
	return dh.db.Close()
}

func (dh *DbHandler) ReadRowByID(id int, object Table) (Table, error) {

	pointers, err := GetPointers(object)

	if err != nil {
		return object, fmt.Errorf("error getting pointers (%e)", err)
	}

	err = dh.db.QueryRow(
		fmt.Sprintf("SELECT * FROM %s WHERE %s = ?",
			object.GetTableName(),
			object.GetFieldsNames()[0],
		),
		id,
	).Scan(pointers...)

	return object, err
}
func (dh *DbHandler) SendRow(object Table) (id int, err error) {

	format := strings.Join(object.GetFieldsNames()[1:], ", ")
	values := GetValues(object)[1:]
	placeholders := "?"
	for i, l := 0, len(values); i < l-1; i++ {
		placeholders += ", ?"
	}

	result, err := dh.db.Exec(
		fmt.Sprintf(
			"INSERT INTO %s (%s) VALUES (%s)",
			object.GetTableName(),
			format,
			placeholders,
		), values...)

	if err != nil {
		return -1, err
	}

	id64, err := result.LastInsertId()
	id = int(id64)

	return
}

func (dh *DbHandler) GetRowCountOf(table Table) (count int, err error) {

	query := fmt.Sprintf("SELECT COUNT(*) AS count FROM %s", table.GetTableName())

	err = dh.db.QueryRow(query).Scan(&count)
	return
}
