package steamid

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCommunityIdToSteamId(t *testing.T) {
	res, err := CommIdToSteamId("76561197999073985")
	assert.Nil(t, err)
	assert.Equal(t, res, "[U:1:38808257:1]")

	res, err = CommIdToSteamId("99999999999999999999999999999999999999999999999999999999999999999999999999999999999999999")
	assert.NotNil(t, err)

	res, err = CommIdToLegacySteamId("76561197999073985")
	assert.Nil(t, err)
	assert.Equal(t, res, "STEAM_0:1:19404128")

	_, err = CommIdToLegacySteamId("999999999999999999999999999999999999999999999999999")
	assert.NotNil(t, err)
}

func TestNewSteamIdToCommunityId(t *testing.T) {
	res, err := SteamIdToCommId("[U:1:38808257]")
	assert.Nil(t, err)
	assert.Equal(t, res, "76561197999073985")

	res, err = SteamIdToCommId("U:1:38808257:1")
	assert.Nil(t, err)
	assert.Equal(t, res, "76561197999073985")

	_, err = CommIdToSteamId("abcd")
	assert.NotNil(t, err)
}

func TestLegacySteamIdToCommunityId(t *testing.T) {
	res, err := SteamIdToCommId("STEAM_0:1:19404128")
	assert.Nil(t, err)
	assert.Equal(t, res, "76561197999073985")

	_, err = CommIdToSteamId("STEAM_1:1:19404128")
	assert.NotNil(t, err)

	_, err = CommIdToSteamId("")
	assert.NotNil(t, err)
}
