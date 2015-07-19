package scraper

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const steamKey = "your steam dev key"
const robinAvatar = "https://steamcdn-a.akamaihd.net/steamcommunity/public/images/avatars/f1/f1dd60a188883caf82d0cbfccfe6aba0af1732d4.jpg"
const robinPUrl = "http://steamcommunity.com/id/robinwalker/"

func TestSteamPlayers(t *testing.T) {
	SetSteamApiKey(steamKey)

	playerInfo, playErr := GetPlayersInfo([]string{"76561197960435530", "76561198067132047"})
	assert.Nil(t, playErr, "There is an error!")

	assert.NotNil(t, playerInfo["76561197960435530"], "The id should be available!")
	assert.NotNil(t, playerInfo["76561198067132047"], "The id should be available!")

	// vars
	visibility := playerInfo["76561197960435530"].Visibility
	realName := playerInfo["76561197960435530"].Realname
	avatar := playerInfo["76561197960435530"].Avatar

	pState := playerInfo["76561197960435530"].Personastate
	name := playerInfo["76561197960435530"].Name
	pUrl := playerInfo["76561197960435530"].Profileurl

	// tests
	assert.Equal(t, PersonaStateOffline, pState, "wtf is he doing not offline?")
	assert.Equal(t, "Robin Walker", realName, "Name should be Robin Walker")
	assert.Equal(t, robinAvatar, avatar, "Avatar should be "+robinAvatar)

	assert.Equal(t, "Offline", pState.String(), "yo, hes online!")
	assert.Equal(t, robinPUrl, pUrl, "no comments")
	assert.Equal(t, "public", visibility, "Profile visibility should be public")

	assert.Equal(t, "Robin", name, "Name should be Robin")
}

func TestSteamPlayer(t *testing.T) {
	SetSteamApiKey(steamKey)

	playerInfo, playErr := GetPlayerInfo("76561197960435530")
	assert.Nil(t, playErr, "There is an error!")

	assert.Equal(t, PersonaStateOffline, playerInfo.Personastate, "wtf is he doing not offline?")
	assert.Equal(t, "Robin Walker", playerInfo.Realname, "Name should be Robin Walker")
	assert.Equal(t, robinAvatar, playerInfo.Avatar, "Avatar should be "+robinAvatar)

	assert.Equal(t, "Offline", playerInfo.Personastate.String(), "yo, hes online!")
	assert.Equal(t, robinPUrl, playerInfo.Profileurl, "no comments")
	assert.Equal(t, "public", playerInfo.Visibility, "Profile visibility should be public")

	assert.Equal(t, "Robin", playerInfo.Name, "Name should be Robin")
}
