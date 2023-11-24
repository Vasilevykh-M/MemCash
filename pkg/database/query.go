package database

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"reflect"
	"strings"
)

type Rows struct {
	Values []Row
}

func calcExpr(query *SelectQuery, d *Database, str []string, args []any) (bool, error) {
	i := 0
	flag := true
	for key, _ := range query.Values {
		if d.Tables[query.Table].Columns[key].Type != reflect.TypeOf(args[i]) {
			return false, fmt.Errorf("error: not equal types, expected: %t, received: %t", d.Tables[query.Table].Columns[key].Type, reflect.TypeOf(args[i]))
		}
		if str[d.Tables[query.Table].Columns[key].Num] != args[i] {
			flag = false
		}
		i++
	}
	return flag, nil
}

func (d *Database) Query(ctx context.Context, sql string, args ...any) (Rows, error) {

	query := NewSelectQuery(sql, d)
	if query == nil {
		return Rows{}, fmt.Errorf("error: %s", "not valid query")
	}

	if len(args) != len(query.Values) {
		return Rows{}, fmt.Errorf("error: to many args, expected: %d, received: %d", len(query.Values), len(args))
	}

	path := d.Tables[query.Table].CurPath
	file, err := os.OpenFile(path, os.O_APPEND|os.O_RDWR, os.ModeAppend)
	defer file.Close()

	if err != nil {
		return Rows{}, fmt.Errorf("error: %s, %s", err, "to opend table")
	}
	s := bufio.NewScanner(file)

	result := Rows{make([]Row, 0)}

	for s.Scan() {
		str := strings.Split(s.Text(), " | ")
		if query.Values == nil {
			result.Values = append(result.Values, Row{buildRow(str, query.Columns)})
			continue
		}

		ok, err := calcExpr(query, d, str, args)

		if err != nil {
			return Rows{}, err
		}
		if ok {
			result.Values = append(result.Values, Row{buildRow(str, query.Columns)})
		}
	}

	if len(result.Values) == 0 {
		return Rows{}, fmt.Errorf("error: not found")
	}

	return result, nil
}
