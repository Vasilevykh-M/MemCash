package database

import (
	"encoding/json"
	sqlParse "github.com/krasun/gosqlparser"
	"strings"
)

const notFoundSelect = -1

func NewSelectQuery(sql string, d *Database) *SelectQuery {

	i := strings.Index(sql, "SELECT")
	if i == notFoundSelect {
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

	err = json.Unmarshal([]byte(jsonStr), &x)
	if err != nil {
		return nil
	}

	result := SelectQuery{}
	result.Table = x["Table"].(string)
	columnsI := x["Columns"].([]interface{})

	columns := make([]string, len(columnsI))

	for i, column := range columnsI {
		columns[i] = column.(string)
	}

	validColumn := make([]int, 0, 1)

	for _, column := range columns {
		if columnDb, ok := d.Tables[result.Table].Columns[column]; !ok {
			return nil
		} else {
			validColumn = append(validColumn, columnDb.Num)
		}
	}

	result.Columns = validColumn

	i = strings.Index(sql, "WHERE")

	if i == -1 {
		return &result
	}

	result.Values = make(map[string]string)

	qwery := sql[i+5 : len(sql)]
	qwery = strings.ReplaceAll(qwery, "AND", "")
	qwery = strings.ReplaceAll(qwery, "==", "")
	sliceParam := strings.Split(qwery, " ")
	idColumn := ""
	for i, param := range sliceParam {

		if i%2 != 1 {
			continue
		}
		if (i+1)%4 == 0 {
			result.Values[idColumn] = param
			continue
		}

		columnDb, ok := d.Tables[result.Table].Columns[param]
		if !ok {
			return nil
		}

		idColumn = columnDb.Name
	}

	return &result
}
