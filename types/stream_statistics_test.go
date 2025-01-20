package types

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewStreamStatistics(t *testing.T) {
	stats, err := NewStreamStatistics("", "")

	assert.NoError(t, err)
	assert.NotNil(t, stats.clientsBlue)
	assert.NotNil(t, stats.clientsGreen)
	assert.True(t, stats.useBlue)
}

func TestAddClient(t *testing.T) {
	stats, err := NewStreamStatistics("", "")

	assert.NoError(t, err)

	clientID := "test-client"
	stats.Add(clientID)

	assert.Contains(t, stats.clientsBlue, clientID)
	assert.Contains(t, stats.clientsGreen, clientID)
}

func TestLen(t *testing.T) {
	stats, err := NewStreamStatistics("", "")
	assert.NoError(t, err)

	stats.Add("client1")
	stats.Add("client2")

	assert.Len(t, stats.clientsBlue, 2)
	assert.Len(t, stats.clientsGreen, 2)
	assert.Equal(t, 2, stats.Len())
}

func TestRotate(t *testing.T) {
	stats, err := NewStreamStatistics("", "")
	assert.NoError(t, err)

	clientID := "test-client"
	stats.Add(clientID)
	err = stats.Rotate()
	assert.NoError(t, err)

	assert.NotNil(t, stats.clientsBlue)
	assert.NotNil(t, stats.clientsGreen)
	assert.Len(t, stats.clientsBlue, 0)
	assert.Len(t, stats.clientsGreen, 1)
	assert.False(t, stats.useBlue)
	assert.Contains(t, stats.clientsGreen, clientID)
	assert.Equal(t, 1, stats.Len())
}

func TestDuplicateClients(t *testing.T) {
	stats, err := NewStreamStatistics("", "")
	assert.NoError(t, err)

	clientID := "test-client"
	stats.Add(clientID)
	stats.Add(clientID)
	assert.Len(t, stats.clientsBlue, 1)
}

func TestConsecutiveRotations(t *testing.T) {
	stats, err := NewStreamStatistics("", "")
	assert.NoError(t, err)

	for i := 0; i < 5; i++ {
		clientID := fmt.Sprintf("client-%d", i)
		stats.Add(clientID)
		err = stats.Rotate()
		assert.NoError(t, err)
	}
	assert.Equal(t, 1, stats.Len())
	assert.Len(t, stats.Timeline(), 5)
}

func TestNilPointer(t *testing.T) {
	var stats *StreamStatistics
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic when calling method on nil pointer")
		}
	}()
	stats.Add("test-client")
}
