package main

import (
	"database/sql"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	// randSource источник псевдо случайных чисел.
	// Для повышения уникальности в качестве seed
	// используется текущее время в unix формате (в виде числа)
	randSource = rand.NewSource(time.Now().UnixNano())
	// randRange использует randSource для генерации случайных чисел
	randRange = rand.New(randSource)
)

// getTestParcel возвращает тестовую посылку
func getTestParcel() Parcel {
	return Parcel{
		Client:    1000,
		Status:    ParcelStatusRegistered,
		Address:   "test",
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}
}

// TestAddGetDelete проверяет добавление, получение и удаление посылки
func TestAddGetDelete(t *testing.T) {
	// prepare
	db, err := sql.Open("sqlite", "tracker.db")
	require.NoError(t, err)
	defer db.Close()

	_, err = db.Exec("DELETE FROM parcel")
	require.NoError(t, err)

	store := NewParcelStore(db)
	parcel := getTestParcel()

	num, err := store.Add(parcel)
	require.NoError(t, err)
	assert.NotEmpty(t, num)

	newParcel, err := store.Get(num)
	require.NoError(t, err)
	parcel.Number = newParcel.Number
	assert.Equal(t, parcel, newParcel)

	err = store.Delete(num)
	require.NoError(t, err)

	_, err = store.Get(num)
	require.ErrorIs(t, err, sql.ErrNoRows)
}

// TestSetAddress проверяет обновление адреса
func TestSetAddress(t *testing.T) {
	// prepare
	db, err := sql.Open("sqlite", "tracker.db")
	require.NoError(t, err)

	defer db.Close()

	store := NewParcelStore(db)
	parcel := getTestParcel()

	number, err := store.Add(parcel)
	require.NoError(t, err)
	assert.NotEmpty(t, number)

	newAddress := "new test address"
	store.SetAddress(number, newAddress)
	require.NoError(t, err)

	uppdateParcel, err := store.Get(number)
	require.NoError(t, err)

	assert.Equal(t, newAddress, uppdateParcel.Address)
}

// TestSetStatus проверяет обновление статуса
func TestSetStatus(t *testing.T) {
	// prepare
	db, err := sql.Open("sqlite", "tracker.db")
	require.NoError(t, err)
	defer db.Close()

	store := NewParcelStore(db)
	parcel := getTestParcel()

	number, err := store.Add(parcel)
	require.NoError(t, err)

	assert.NotEmpty(t, number)

	store.SetStatus(number, ParcelStatusDelivered)
	require.NoError(t, err)

	uppdateParcel, err := store.Get(number)
	require.NoError(t, err)
	assert.Equal(t, ParcelStatusDelivered, uppdateParcel.Status)
}

// TestGetByClient проверяет получение посылок по идентификатору клиента
func TestGetByClient(t *testing.T) {
	// prepare
	db, err := sql.Open("sqlite", "tracker.db")
	require.NoError(t, err)

	defer db.Close()

	store := NewParcelStore(db)

	parcels := []Parcel{
		getTestParcel(),
		getTestParcel(),
		getTestParcel(),
	}
	parcelMap := map[int]Parcel{}

	// задаём всем посылкам один и тот же идентификатор клиента
	client := randRange.Intn(10_000_000)
	parcels[0].Client = client
	parcels[1].Client = client
	parcels[2].Client = client

	// add
	for i := 0; i < len(parcels); i++ {
		id, err := store.Add(parcels[i])
		require.NoError(t, err)
		assert.NotEmpty(t, id)

		// обновляем идентификатор добавленной у посылки
		parcels[i].Number = id

		// сохраняем добавленную посылку в структуру map, чтобы её можно было легко достать по идентификатору посылки
		parcelMap[id] = parcels[i]
	}

	// get by client
	storedParcels, err := store.GetByClient(client)
	require.NoError(t, err)

	assert.Len(t, parcels, len(storedParcels))

	// check
	for _, parcel := range storedParcels {
		assert.NotEmpty(t, parcel.Number, parcelMap[parcel.Number])
		assert.Equal(t, parcel, parcelMap[parcel.Number])
		// в parcelMap лежат добавленные посылки, ключ - идентификатор посылки, значение - сама посылка
		// убедитесь, что все посылки из storedParcels есть в parcelMap
		// убедитесь, что значения полей полученных посылок заполнены верно
	}
}
