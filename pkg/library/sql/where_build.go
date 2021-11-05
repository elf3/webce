package sql

import (
	"fmt"
	"strings"
)

type NullType byte

const (
	_ NullType = iota
	// IsNull the same as `is null`
	IsNull
	// IsNotNull the same as `is not null`
	IsNotNull
)

// WhereBuild sql build where
func WhereBuild(where map[string]interface{}) (sql string, args []interface{}, err error) {
	args = make([]interface{}, 0)
	for k, v := range where {
		ks := strings.Split(k, " ")
		if len(ks) > 2 {
			return "", nil, fmt.Errorf("Error in query condition: %s. ", k)
		}
		if sql != "" {
			sql += " AND "
		}
		switch len(ks) {
		case 1:
			switch v := v.(type) {
			case NullType:
				if v == IsNotNull {
					sql += fmt.Sprint(k, " IS NOT NULL")
				} else {
					sql += fmt.Sprint(k, " IS NULL")
				}
			default:
				sql += fmt.Sprint(k, "=?")
				args = append(args, v)
			}

		case 2:
			k = ks[0]
			switch ks[1] {
			case "=":
				sql += fmt.Sprint(k, "=?")
				args = append(args, v)
			case ">":
				sql += fmt.Sprint(k, ">?")
				args = append(args, v)

			case ">=":
				sql += fmt.Sprint(k, ">=?")
				args = append(args, v)

			case "<":
				sql += fmt.Sprint(k, "<?")
				args = append(args, v)

			case "<=":
				sql += fmt.Sprint(k, "<=?")
				args = append(args, v)

			case "!=":
				sql += fmt.Sprint(k, "!=?")
				args = append(args, v)

			case "<>":
				sql += fmt.Sprint(k, "!=?")
				args = append(args, v)

			case "in":
				sql += fmt.Sprint(k, " in (?) ")
				args = append(args, v)

			case "like":
				sql += fmt.Sprint(k, " like ? ")
				args = append(args, v)
			}

		}
	}
	return
}
