package psg

import (
	"github.com/jackc/pgx/v5"
)

// Psg представляет гейт к базе данных PostgreSQL.
type Psg struct {
	Conn *pgxpool.Pool
}

// NewPsg создает новый экземпляр Psg.
func NewPsg(psgAddr string, login, password string) *Psg {
	return
}

// RecordAdd добавляет новую запись в базу данных.
func (p *Psg) RecordAdd(record Record) (int64, error) {
	// TODO: Реализовать добавление записи
	return 0, nil
}

// RecordsGet возвращает записи из базы данных на основе предоставленных полей Record.
func (p *Psg) RecordsGet(record Record) ([]Record, error) {
	// TODO: Реализовать поиск записей
	return nil, nil
}

// RecordUpdate обновляет существующую запись в базе данных по номеру телефона.
func (p *Psg) RecordUpdate(record Record) error {
	// TODO: Реализовать обновление записи
	return nil
}

// RecordDeleteByPhone удаляет запись из базы данных по номеру телефона.
func (p *Psg) RecordDeleteByPhone(phone string) error {
	// TODO: Реализовать удаление записи по номеру телефона
	return nil
}
