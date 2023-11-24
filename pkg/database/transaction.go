package database

import (
	"context"
	"io"
	"os"
)

func copyFileContents(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		return
	}
	err = out.Sync()
	return
}

func (d *Database) Commit(ctx context.Context) error {
	defer d.mx.Unlock()
	for _, table := range d.Tables {
		err := os.Remove(table.OldPath)
		if err != nil {
			return err
		}
	}
	return nil
}

func (d *Database) Begin(ctx context.Context) error {
	for _, table := range d.Tables {
		err := copyFileContents(table.CurPath, table.OldPath)
		if err != nil {
			d.mx.Unlock()
			return err
		}
	}
	d.mx.Lock()
	return nil
}

func (d *Database) Rollback(ctx context.Context) error {
	defer d.mx.Unlock()

	for _, table := range d.Tables {
		err := os.Remove(table.CurPath)
		if err != nil {
			return err
		}
		err = copyFileContents(table.OldPath, table.CurPath)

		if err != nil {
			return err
		}

		err = os.Remove(table.OldPath)
		if err != nil {
			return err
		}
	}
	return nil
}
