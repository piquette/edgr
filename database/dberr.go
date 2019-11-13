package database

import pg "github.com/go-pg/pg"

func isDuplicate(err error) bool {
	var code string

	if err, ok := err.(pg.Error); ok {
		code = err.Field('C')
	}

	switch code {
	case "23000":
		return true
	case "23505":
		return true
	}

	return false
}
