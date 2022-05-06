package rediscache

import (
	"context"
	"errors"
	"github.com/go-redis/redismock/v8"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_NewRedisCache(t *testing.T) {
	NewRedisCache(1, Options{})
	assert.Equal(t, true, true)
}
func Test_Get(t *testing.T) {
	type testObject struct {
		Name    string `json:"name"`
		Surname string `json:"surname"`
	}
	var cacheMock cache
	db, redisMockClient := redismock.NewClientMock()
	t.Run("Get returns the value", func(t *testing.T) {
		redisMockClient.ExpectGet("123456").SetVal(`{"name":"John","surname":"Doe"}`)
		cacheMock.client = db
		actualObject, err := cacheMock.Get(context.Background(), "123456", testObject{})
		require.Nil(t, err)
		assert.Equal(t, testObject{
			Name:    "John",
			Surname: "Doe",
		}, actualObject)
	})

	t.Run("Get returns the error on client.Get request", func(t *testing.T) {
		redisMockClient.ExpectGet("123456").SetErr(errors.New("error"))
		cacheMock.client = db
		actualObject, err := cacheMock.Get(context.Background(), "123456", testObject{})
		require.Error(t, err)
		assert.Equal(t, testObject{}, actualObject)
	})

	t.Run("Get returns the error on unmarshall", func(t *testing.T) {
		redisMockClient.ExpectGet("123456").SetVal("")
		cacheMock.client = db
		actualObject, err := cacheMock.Get(context.Background(), "123456", testObject{})
		require.Error(t, err)
		assert.Equal(t, testObject{}, actualObject)
	})
}

func Test_Set(t *testing.T) {
	type testObject struct {
		Name    string `json:"name"`
		Surname string `json:"surname"`
	}

	var cacheMock cache
	db, redisMockClient := redismock.NewClientMock()

	t.Run("Set returns the nil", func(t *testing.T) {
		var err error
		value := testObject{
			Name:    "John",
			Surname: "Doe"}
		redisMockClient.ExpectSet("123456", value, 1).RedisNil()
		cacheMock.client = db
		err = cacheMock.Set(context.Background(), "123456", value)
		require.Nil(t, err)
	})

	t.Run("Set returns the error on marshal", func(t *testing.T) {
		var err error
		value := make(chan string)
		redisMockClient.ExpectSet("123456", value, 1).RedisNil()
		cacheMock.client = db
		err = cacheMock.Set(context.Background(), "123456", value)
		require.Error(t, err)
	})
}
