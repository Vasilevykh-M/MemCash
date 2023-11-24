package database

import (
	"bufio"
	"context"
	"errors"
	"os"
	"reflect"
	"strings"
)

func (d *Database) Exec(ctx context.Context, sql string, arguments ...any) (err error) {
	query := NewInsertQuery(sql, d)
	if query == nil {
		return errors.New("error")
	}

	path := d.Tables[query.Table].CurPath
	file, err := os.OpenFile(path, os.O_APPEND|os.O_RDWR, os.ModeAppend)
	defer file.Close()

	if err != nil {
		return err
	}

	s := bufio.NewWriter(file)

	i := 0

	for _, val := range query.Columns {
		if d.Tables[query.Table].Columns[val].Type != reflect.TypeOf(arguments[i]) {
			return errors.New("error")
		}
		i++
	}

	arg := make([]string, 0, 1)

	for _, argument := range arguments {
		arg = append(arg, argument.(string))
	}

	strRes := strings.Join(arg, " | ")
	_, err = s.WriteString(strRes + "\n")

	if err != nil {
		return err
	}

	s.Flush()

	return nil
}
