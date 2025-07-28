package adapters_test

import (
	"testing"

	"github.com/go-redis/redismock/v9"
	"github.com/slilp/go-wallet/internal/adapters"
	"github.com/stretchr/testify/suite"
)

type AdaptersTestSuite struct {
	suite.Suite
	mockRedis    redismock.ClientMock
	redisAdapter adapters.RedisAdapter
}

func (suite *AdaptersTestSuite) SetupTest() {
	client, mockRedis := redismock.NewClientMock()
	suite.mockRedis = mockRedis
	suite.redisAdapter = adapters.NewRedisAdapter(client)
}

func TestAdaptersTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(AdaptersTestSuite))
}
