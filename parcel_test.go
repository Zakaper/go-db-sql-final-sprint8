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
	db, err := sql.Open("sqlite", "trcker.db")
	if err != nil {
		require.NoError(t, err)
	}
	defer db.Close()
	store := NewParcelStore(db)
	parcel := getTestParcel()

	number, err := store.Add(parcel)
	if err != nil {
		require.NoError(t, err)
	}
	assert.NotEmpty(t, number)

	newParcel, err := store.Get(number)
	if err != nil {
		require.NoError(t, err)
	}
	assert.Equal(t, parcel.Number, newParcel.Number)
	assert.Equal(t, parcel.Address, newParcel.Address)
	assert.Equal(t, parcel.Client, newParcel.Client)
	assert.Equal(t, parcel.CreatedAt, newParcel.CreatedAt)

	err = store.Delete(number)
	if err != nil {
		require.NoError(t, err)
	}
	_, err = store.Get(number)
	if err != nil {
		require.Equal(t, sql.ErrNoRows, err)
	}
}

// TestSetAddress проверяет обновление адреса
func TestSetAddress(t *testing.T) {
	// prepare
	db, err := sql.Open("sqlite", "tracker.db")
	if err != nil {
		require.NoError(t, err)
	}
	defer db.Close()

	store := NewParcelStore(db)
	parcel := getTestParcel()

	number, err := store.Add(parcel)
	if err != nil {
		require.NoError(t, err)
	}
	assert.Empty(t, number)

	newAddress := "new test address"
	store.SetAddress(number, newAddress)
	require.NoError(t, err)

	uppdateParcel, err := store.Get(number)
	if err != nil {
		require.NoError(t, err)
	}
	assert.Equal(t, newAddress, uppdateParcel.Address)
}

// TestSetStatus проверяет обновление статуса
func TestSetStatus(t *testing.T) {
	// prepare
	db, err := sql.Open("sqlite", "tracker.db")
	if err != nil {
		require.NoError(t, err)
	}
	defer db.Close()

	store := NewParcelStore(db)
	parcel := getTestParcel()

	number, err := store.Add(parcel)
	if err != nil {
		require.NoError(t, err)
	}
	assert.Empty(t, number)

	store.SetStatus(number, ParcelStatusSent)
	require.NoError(t, err)

	uppdateParcel, err := store.Get(number)
	if err != nil {
		require.NoError(t, err)
	}
	assert.Equal(t, ParcelStatusSent, uppdateParcel.Status)
}

// TestGetByClient проверяет получение посылок по идентификатору клиента
func TestGetByClient(t *testing.T) {
	// prepare
	db, err := sql.Open("sqlite", "tracker.db")
	if err != nil {
		require.NoError(t, err)
	}
	defer db.Close()
	store := NewParcelStore(db)
	parcel := getTestParcel()

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
		id, err := store.Add(parcel)
		if err != nil {
			require.NoError(t, err)
		}
		// обновляем идентификатор добавленной у посылки
		parcels[i].Number = id

		// сохраняем добавленную посылку в структуру map, чтобы её можно было легко достать по идентификатору посылки
		parcelMap[id] = parcels[i]
	}

	// get by client
	storedParcels, err := store.GetByClient(client)
	if err != nil {
		require.NoError(t, err)
	}
	assert.Len(t, parcels, len(storedParcels))

	// check
	for _, parcel := range storedParcels {
		assert.NotEmpty(t, parcelMap[parcel.Number])
		assert.Equal(t, parcel, parcelMap[parcel.Number])
		// в parcelMap лежат добавленные посылки, ключ - идентификатор посылки, значение - сама посылка
		// убедитесь, что все посылки из storedParcels есть в parcelMap
		// убедитесь, что значения полей полученных посылок заполнены верно
	}
}
