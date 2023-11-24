package database

import (
	"bufio"
	"context"
	"os"
	"strings"
)

type Row struct {
	Values []string
}

func (d *Database) QueryRow(ctx context.Context, sql string, args ...any) Row {
	query := NewSelectQuery(sql, d)
	if query == nil {
		return Row{}
	}

	if len(args) != len(query.Values) {
		return Row{}
	}

	path := d.Tables[query.Table].CurPath
	file, err := os.OpenFile(path, os.O_APPEND|os.O_RDWR, os.ModeAppend)
	defer file.Close()

	if err != nil {
		return Row{}
	}
	s := bufio.NewScanner(file)

	for s.Scan() {
		str := strings.Split(s.Text(), " | ")

		if query.Values == nil {
			return Row{buildRow(str, query.Columns)}
		}

		ok, err := calcExpr(query, d, str, args)

		if err != nil {
			return Row{}
		}

		if ok {
			return Row{buildRow(str, query.Columns)}
		}
	}
	return Row{}
}
