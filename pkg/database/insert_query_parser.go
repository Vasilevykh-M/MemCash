package database

import (
	"encoding/json"
	sqlParse "github.com/krasun/gosqlparser"
	"strings"
)

func columnIsTable(column, table string, d *Database) bool {
	if _, ok := d.Tables[table].Columns[column]; !ok {
		return false
	} else {
		return true
	}
}

func NewInsertQuery(sql string, d *Database) *InsertQuery {

	i := strings.Index(sql, "INSERT")
	if i == -1 {
		return nil
	}

	query, err := sqlParse.Parse(sql)
	if err != nil {
		return nil
	}
	jsonStr, err := json.Marshal(query)
	if err != nil {
		return nil
	}

	x := map[string]interface{}{}
	err = json.Unmarshal([]byte(string(jsonStr)), &x)
	if err != nil {
		return nil
	}

	result := InsertQuery{}
	result.Table = x["Table"].(string)

	columnsI := x["Columns"].([]interface{})

	validColumn := make([]string, 0, 1)

	columns := make([]string, len(columnsI))

	for i, column := range columnsI {
		columns[i] = column.(string)
	}

	for _, column := range columns {

		ok := columnIsTable(column, result.Table, d)

		if !ok {
			return nil
		} else {
			validColumn = append(validColumn, column)
		}
	}

	if len(validColumn) != len(d.Tables[result.Table].Columns) {
		return nil
	}
	if len(x["Values"].([]interface{})) != len(d.Tables[result.Table].Columns) {
		return nil
	}

	result.Columns = validColumn

	return &result
}
