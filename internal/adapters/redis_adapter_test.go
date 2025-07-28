package adapters_test

import (
	"context"
	"testing"
	"time"

	"github.com/go-redis/redismock/v9"
	"github.com/slilp/go-wallet/internal/adapters"
	"github.com/stretchr/testify/assert"
)

type walletInfo struct {
	Id      string
	Balance string
}

func TestRedisAdapter_SetJSON_Array(t *testing.T) {
	client, mock := redismock.NewClientMock()
	adapter := adapters.NewRedisAdapter(client)

	key := "wallets"
	value := []walletInfo{
		{Id: "1", Balance: "20"},
		{Id: "2", Balance: "50"},
	}
	jsonBytes := []byte(`[{"Id":"1","Balance":"20"},{"Id":"2","Balance":"50"}]`)

	testCases := []struct {
		name    string
		mock    func()
		wantErr bool
	}{
		{
			name: "GivenValidArray_WhenSetJSON_ThenSuccess",
			mock: func() {
				mock.ExpectSet(key, jsonBytes, time.Minute).SetVal("OK")
			},
			wantErr: false,
		},
		{
			name: "GivenRedisError_WhenSetJSON_ThenError",
			mock: func() {
				mock.ExpectSet(key, jsonBytes, time.Minute).SetErr(assert.AnError)
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mock()
			err := adapter.SetJSON(context.Background(), key, value, time.Minute)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestRedisAdapter_GetJSON_Array(t *testing.T) {
	client, mock := redismock.NewClientMock()
	adapter := adapters.NewRedisAdapter(client)

	key := "wallets"
	jsonBytes := []byte(`[{"Id":"1","Balance":"20"},{"Id":"2","Balance":"50"}]`)

	testCases := []struct {
		name     string
		mock     func()
		wantErr  bool
		expected []walletInfo
	}{
		{
			name: "GivenValidArray_WhenGetJSON_ThenSuccess",
			mock: func() {
				mock.ExpectGet(key).SetVal(string(jsonBytes))
			},
			wantErr:  false,
			expected: []walletInfo{{Id: "1", Balance: "20"}, {Id: "2", Balance: "50"}},
		},
		{
			name: "GivenRedisError_WhenGetJSON_ThenError",
			mock: func() {
				mock.ExpectGet(key).SetErr(assert.AnError)
			},
			wantErr:  true,
			expected: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mock()
			var result []walletInfo
			err := adapter.GetJSON(context.Background(), key, &result)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, result)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
