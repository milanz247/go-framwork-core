package database

import (
	"database/sql"

	_ "modernc.org/sqlite"

	"ticket-sync-mvc/framework/support"
)

var DB *sql.DB

func Rows(q string, args ...any) []support.M {
	r, e := DB.Query(q, args...)
	if e != nil {
		panic(e)
	}
	defer r.Close()
	cols, e := r.Columns()
	if e != nil {
		panic(e)
	}
	out := []support.M{}
	for r.Next() {
		v := make([]any, len(cols))
		p := make([]any, len(cols))
		for i := range v {
			p[i] = &v[i]
		}
		if e = r.Scan(p...); e != nil {
			panic(e)
		}
		m := support.M{}
		for i, c := range cols {
			m[c] = v[i]
		}
		out = append(out, m)
	}
	if e = r.Err(); e != nil {
		panic(e)
	}
	return out
}

func Exec(q string, args ...any) sql.Result {
	r, e := DB.Exec(q, args...)
	if e != nil {
		panic(e)
	}
	return r
}
