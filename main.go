package main

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"text/template"
)

func main() {
	r := Record{
		Name:       "Иван",
		LastName:   "Иванов",
		MiddleName: "Иванович",
		Address:    "Москва",
		Phone:      "1234567890",
	}

	err := SelectRecord(r)
	if err != nil {
		fmt.Println(err)
	}
}

type Cond struct {
	Lop    string
	PgxInd string
	Field  string
	Value  any
}

func SelectRecord(r Record) (err error) {
	sqlFields, values, err := StructToFieldsValues(r, "sql.field")
	if err != nil {
		return
	}

	var conds []Cond

	for i := range sqlFields {
		if i == 0 {
			conds = append(conds, Cond{
				Lop:    "",
				PgxInd: "$" + strconv.Itoa(i+1),
				Field:  sqlFields[i],
				Value:  values[i],
			})
			continue
		}
		conds = append(conds, Cond{
			Lop:    "AND",
			PgxInd: "$" + strconv.Itoa(i+1),
			Field:  sqlFields[i],
			Value:  values[i],
		})
	}

	query := `
        SELECT 
            id, name, last_name, middle_name, address, phone
        FROM
            address_book
        WHERE
            {{range .}} {{.Lop}} {{.Field}} = {{.PgxInd}}{{end}}
        ;
    `

	tmpl, err := template.New("").Parse(query)
	if err != nil {
		return
	}

	var sb strings.Builder
	err = tmpl.Execute(&sb, conds)
	if err != nil {
		return
	}
	fmt.Println(sb.String())
	return
}

func StructToFieldsValues(s any, tag string) (sqlFields []string, values []any, err error) {
	rv := reflect.ValueOf(s)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}

	if rv.Kind() != reflect.Struct {
		return nil, nil, errors.New("s must be a struct")
	}

	for i := 0; i < rv.NumField(); i++ {
		field := rv.Type().Field(i)
		tg := strings.TrimSpace(field.Tag.Get(tag))
		if tg == "" || tg == "-" {
			continue
		}
		tgs := strings.Split(tg, ",")
		tg = tgs[0]

		fv := rv.Field(i)
		isZero := false
		switch fv.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			isZero = fv.Int() == 0
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			isZero = fv.Uint() == 0
		case reflect.Float32, reflect.Float64:
			isZero = fv.Float() == 0
		case reflect.Complex64, reflect.Complex128:
			isZero = fv.Complex() == complex(0, 0)
		case reflect.Bool:
			isZero = !fv.Bool()
		case reflect.String:
			isZero = fv.String() == ""
		case reflect.Array, reflect.Slice:
			isZero = fv.Len() == 0
		}

		if isZero {
			continue
		}

		sqlFields = append(sqlFields, tg)
		values = append(values, fv.Interface())
	}

	return
}
