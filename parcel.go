package main

import (
	"database/sql"
	"fmt"
)

type ParcelStore struct {
	db *sql.DB
}

func NewParcelStore(db *sql.DB) ParcelStore {
	return ParcelStore{db: db}
}

func (s ParcelStore) Add(p Parcel) (int, error) {
	// реализуйте добавление строки в таблицу parcel, используйте данные из переменной p
	query := "INSERT INTO parcel (client, status, address, created_at) VALUES (?, ?, ?, ?)"
	res, err := s.db.Exec(query, p.Client, p.Status, p.Address, p.CreatedAt)
	if err != nil {
		return 0, err
	}
	// верните идентификатор последней добавленной записи
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func (s ParcelStore) Get(number int) (Parcel, error) {

	p := Parcel{}

	// реализуйте чтение строки по заданному number
	query := "SELECT number, client, status, address, created_at FROM parcel WHERE number = ?"

	// выполняем запрос
	row := s.db.QueryRow(query, number)
	// здесь из таблицы должна вернуться только одна строка
	err := row.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return p, fmt.Errorf("посылка с номером %d не найдена", number)
		}
		return p, err
	}
	// заполните объект Parcel данными из таблицы
	return p, nil
}

func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {
	// заполните срез Parcel данными из таблицы
	var res []Parcel

	// реализуйте чтение строк из таблицы parcel по заданному client
	rows, err := s.db.Query("SELECT number, client, status, address, created_at FROM parcel WHERE client = ?", client)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// здесь из таблицы может вернуться несколько строк
	for rows.Next() {
		var p Parcel
		err := rows.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
		if err != nil {
			return nil, err
		}
		res = append(res, p)
	}
	// исправлено ошибка при чтении строк
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return res, nil
}

func (s ParcelStore) SetStatus(number int, status string) error {
	// реализуйте обновление статуса в таблице parcel
	result, err := s.db.Exec("UPDATE parcel SET status = ? WHERE number = ?", status, number)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no rows were ipdated")
	}
	return nil
}

func (s ParcelStore) SetAddress(number int, address string) error {
	// реализуйте обновление адреса в таблице parcel
	// менять адрес можно только если значение статуса registered
	result, err := s.db.Exec("UPDATE parcel SET address = ? WHERE number = ? AND status = 'registered'", address, number)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no rows were updated")
	}
	return nil
}

func (s ParcelStore) Delete(number int) error {
	// реализуйте удаление строки из таблицы parcel
	// удалять строку можно только если значение статуса registered
	result, err := s.db.Exec("DELETE FROM parcel WHERE number = ? and status = 'registered'", number)
	if err != nil {
		return err
	}
	// проверяем количество затронутых строк, но не возвращаем ошибку
	_, _ = result.RowsAffected()
	return nil
}
