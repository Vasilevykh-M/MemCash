package database

import (
	"encoding/json"
	"io"
	"os"
	"reflect"
	"sync"
)

type SelectQuery struct {
	Columns []int
	Table   string
	Values  map[string]string
}

type InsertQuery struct {
	Table   string
	Columns []string
}

type Database struct {
	Tables map[string]Table
	mx     sync.Mutex
}

type Table struct {
	Name    string
	CurPath string
	OldPath string
	Columns map[string]Column
}

type Column struct {
	Name string
	Num  int
	Type reflect.Type
}

func buildRow(str []string, selectColumns []int) []string {
	result := make([]string, 0, 1)
	for _, column := range selectColumns {
		result = append(result, str[column])
	}
	return result
}

type TableConn struct {
	PathSchema string `json:"pathSchema"`
	PathData   string `json:"pathData"`
	Name       string `json:"name"`
}

type ColumnConn struct {
	Name string `json:"name"`
	Id   int    `json:"id"`
	Type string `json:"type"`
}

func NewDatabase(pathToDB string) *Database {
	file, err := os.OpenFile(pathToDB, os.O_APPEND|os.O_RDWR, os.ModeAppend)
	defer file.Close()

	if err != nil {
		return nil
	}

	db := Database{}

	tables := []TableConn{}
	byteValue, err := io.ReadAll(file)

	if err != nil {
		return nil
	}

	err = json.Unmarshal(byteValue, &tables)

	if err != nil {
		return nil
	}

	tableMap := map[string]Table{}

	for _, table := range tables {
		columns := []ColumnConn{}
		file, err := os.OpenFile(table.PathSchema, os.O_APPEND|os.O_RDWR, os.ModeAppend)
		defer file.Close()

		if err != nil {
			return nil
		}

		byteValue, err := io.ReadAll(file)

		if err != nil {
			return nil
		}

		err = json.Unmarshal(byteValue, &columns)

		if err != nil {
			return nil
		}

		columnMap := map[string]Column{}

		for _, column := range columns {
			var columnType reflect.Type
			switch column.Type {
			case "string":
				columnType = reflect.TypeOf("string")
			}
			columnMap[column.Name] = Column{column.Name, column.Id, columnType}
		}

		tableMap[table.Name] = Table{table.Name, table.PathData, table.PathData + "_reserve", columnMap}
	}
	db.Tables = tableMap
	return &db
}
