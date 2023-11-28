package psg

import (
	"addressBookServer/models/dto"
	"context"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"text/template"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Psg представляет гейт к базе данных PostgreSQL.
type Psg struct {
	Conn *pgxpool.Pool
}

// NewPsg создает новый экземпляр Psg.
func NewPsg(url string) *Psg {
    dbpool, err := pgxpool.New(context.Background(), url)
    if err != nil {
        panic(fmt.Sprintf("%s connection error", url))
    }

    return &Psg{dbpool}
}

// RecordAdd добавляет новую запись в базу данных.
func (p *Psg) RecordAdd(record dto.Record) (int64, error) {
    query := `INSERT INTO address_book (name, last_name, middle_name, address, phone) VALUES ($1, $2, $3, $4, $5) RETURNING id`

    var id int64 
    err := p.Conn.QueryRow(context.Background(), query, 
        record.Name,
        record.LastName,
        record.MiddleName,
        record.Address,
        record.Phone,
    ).Scan(&id)

	return id, err 
}

// RecordsGet возвращает записи из базы данных на основе предоставленных полей Record.
func (p *Psg) RecordsGet(record dto.Record) ([]dto.Record, error) {

	return nil, nil
}

// RecordUpdate обновляет существующую запись в базе данных по номеру телефона.
func (p *Psg) RecordUpdate(record dto.Record) error {
       
	return nil
}

// RecordDeleteByPhone удаляет запись из базы данных по номеру телефона.
func (p *Psg) RecordDeleteByPhone(phone string) error {
    query := `DELETE FROM note WHERE phone = $1` 
    _, err := p.Conn.Exec(context.Background(), query, phone)

	return err 
}

type Cond struct {
	Lop    string
	PgxInd string
	Field  string
	Value  any
}

func SelectRecord(r dto.Record) (query string, err error) {
	sqlFields, values, err := GetTagsAndFieldsValues(r, "sql.field")
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

	queryTmpl := `
        SELECT 
            id, name, last_name, middle_name, address, phone
        FROM
            address_book
        WHERE
            {{range .}} {{.Lop}} {{.Field}} = {{.PgxInd}}{{end}}
        ;
    `

	tmpl, err := template.New("").Parse(queryTmpl)
	if err != nil {
		return
	}

	var sb strings.Builder
	err = tmpl.Execute(&sb, conds)
	if err != nil {
		return
	}

    query = sb.String()

	return query, nil
}

func GetTagsAndFieldsValues(inst any, tag string) (fieldsTagValues []string, fieldsValues []any, err error) {
	rv := reflect.ValueOf(inst)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}

	if rv.Kind() != reflect.Struct {
		return nil, nil, errors.New("s must be a struct")
	}

	for i := 0; i < rv.NumField(); i++ {
		field := rv.Type().Field(i)
		fieldTagValue := strings.TrimSpace(field.Tag.Get(tag))
		if fieldTagValue == "" || fieldTagValue == "-" {
			continue
		}
		fieldTagValue = strings.Split(fieldTagValue, ",")[0]

		fv := rv.Field(i)
		if reflectValIsZero(fv) {
			continue
		}

		fieldsTagValues = append(fieldsTagValues, fieldTagValue)
		fieldsValues = append(fieldsValues, fv.Interface())
	}

	return
}

func reflectValIsZero(v reflect.Value) bool {
    isZero := false

    switch v.Kind() {
    case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
        isZero = v.Int() == 0
    case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
        isZero = v.Uint() == 0
    case reflect.Float32, reflect.Float64:
        isZero = v.Float() == 0
    case reflect.Complex64, reflect.Complex128:
        isZero = v.Complex() == complex(0, 0)
    case reflect.Bool:
        isZero = !v.Bool()
    case reflect.String:
        isZero = v.String() == ""
    case reflect.Array, reflect.Slice:
        isZero = v.Len() == 0
    }

    return isZero
}
