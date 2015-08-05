package scraper

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"time"
)

type PersonaState int

const (
	PersonaStateOffline PersonaState = 0
	PersonaStateOnline  PersonaState = 1
	PersonaStateBusy    PersonaState = 2
	PersonaStateAway    PersonaState = 3
	PersonaStateSnooze  PersonaState = 4
	PersonaStateLTTrade PersonaState = 5
	PersonaStateLTPlay  PersonaState = 6
)

var PersonaStatesStr = [...]string{
	"Offline",
	"Online",
	"Busy",
	"Away",
	"Snooze",
	"Looking to Trade",
	"Looking to Play",
}

func (ps *PersonaState) String() string {
	return PersonaStatesStr[int(*ps)]
}

type PlayerInfo struct {
	// name
	Name        string
	Personaname string
	Realname    string

	// profile
	Profileurl        string
	Personastate      PersonaState
	Profilestate      int
	Commentpermission string

	// steam
	Steamid                  string
	Communityvisibilitystate int
	Visibility               string

	// time creaTed
	Timecreated int

	// location
	Loccountrycode string
	Locstatecode   string
	Loccityid      int

	// logoff
	Lastlogoff int

	// avatar
	Avatar       string
	Avatarfull   string
	Avatarmedium string
}

// get player "timecreated" as Time
func (p *PlayerInfo) GetTimeCreated() time.Time {
	return time.Unix(int64(p.Timecreated), 0)
}

// get player "lastlogoff" as Time
func (p *PlayerInfo) GetLogoffTime() time.Time {
	return time.Unix(int64(p.Lastlogoff), 0)
}

func tryParsingString(arg interface{}) string {
	res, ok := arg.(string)
	if !ok {
		return ""
	}
	return res
}

func tryParsingNumber(arg interface{}) (int, error) {
	json, ok := arg.(json.Number)
	if !ok {
		return 0, errors.New("Not a json number")
	}

	res, err := json.Int64()
	if err != nil {
		return 0, err
	}
	return int(res), nil
}

// parse player json
func (p *PlayerInfo) Parse(elem map[string]interface{}) error {
	p.Steamid = elem["steamid"].(string)
	p.Profileurl = elem["profileurl"].(string)

	// profile visibility
	visibility := elem["communityvisibilitystate"].(json.Number).String()
	var vState string

	switch {
	case visibility == "1":
		vState = "private"
	case visibility == "3":
		vState = "public"
	}

	// same as communityvisibilitystate
	// but as string (public or private)
	p.Visibility = vState

	playerV, vErr := strconv.Atoi(visibility)
	if vErr != nil {
		return vErr
	}

	// int as str
	p.Communityvisibilitystate = playerV

	// Logoff
	playerL, lErr := tryParsingNumber(elem["lastlogoff"])
	if lErr != nil {
		return lErr
	}

	p.Lastlogoff = playerL

	// if the account has a steam community profile set then this should be 1
	profileState, stErr := tryParsingNumber(elem["profilestate"])
	if stErr != nil {
		return stErr
	}

	p.Profilestate = profileState

	// variables that are only available when
	// the profile visibility is set to public
	if vState == "public" {
		p.Realname = tryParsingString(elem["realname"])

		// user is online, offline...
		playerState, pstErr := tryParsingNumber(elem["personastate"])
		if pstErr != nil {
			return pstErr
		}

		p.Personastate = PersonaState(playerState)

		// timecreated
		playerTC, tcErr := tryParsingNumber(elem["timecreated"])
		if tcErr != nil {
			return tcErr
		}

		p.Timecreated = playerTC

		// location
		p.Loccountrycode = tryParsingString(elem["loccountrycode"])
		p.Locstatecode = tryParsingString(elem["locstatecode"])

		cityId, cErr := tryParsingNumber(elem["loccityid"])
		if cErr != nil {
			cityId = 0
		}

		p.Loccityid = cityId
	}

	// avatar
	p.Avatar = elem["avatar"].(string)
	p.Avatarfull = elem["avatarfull"].(string)
	p.Avatarmedium = elem["avatarmedium"].(string)

	// name
	p.Personaname = elem["personaname"].(string)
	p.Name = p.Personaname // alias

	return nil
}

var steamApiKey string

func SetSteamApiKey(key string) {
	steamApiKey = key
}

func GetTF2Hours(steamid string) (int, error) {
	games, err := GetSteamGamesOwned(steamid)
	if err != nil {
		return 0, err
	}

	return GetTF2HoursFromGamesOwned(games), nil
}

func GetTF2HoursFromGamesOwned(res *map[string]string) int {
	hours, ok := (*res)["440"]

	if !ok {
		return 0
	}

	num, err := strconv.ParseInt(hours, 10, 32)
	if err != nil {
		return 0
	}

	return int(num) / 60
}

func GetSteamGamesOwned(steamid string) (*map[string]string, error) {
	url := "http://api.steampowered.com/IPlayerService/GetOwnedGames/v0001/?key=" +
		steamApiKey + "&steamid=" + steamid + "&include_played_free_games=1&format=json"

	response, err := getJsonFromUrl(url)

	if err != nil {
		return nil, err
	}

	res := make(map[string]string)

	ps, err := response.Get("response").Get("games").Array()
	if err != nil {
		return nil, err
	}

	for _, _elem := range ps {
		elem := _elem.(map[string]interface{})
		name := elem["appid"].(json.Number).String()
		value := elem["playtime_forever"].(json.Number).String()
		res[name] = value
	}

	return &res, nil
}

func GetTF2Stats(steamid string) (*map[string]string, error) {
	url := "http://api.steampowered.com/ISteamUserStats/GetUserStatsForGame/v0002/?appid=440&key=" +
		steamApiKey + "&steamid=" + steamid

	res := make(map[string]string)

	response, err := getJsonFromUrl(url)
	if err != nil {
		return nil, err
	}

	ps, err := response.Get("playerstats").Get("stats").Array()
	if err != nil {
		return nil, err
	}

	for _, _elem := range ps {
		elem := _elem.(map[string]interface{})
		name, _ := elem["name"].(string)
		value := elem["value"].(json.Number).String()
		res[name] = value
	}

	if err != nil {
		return nil, err
	}

	return &res, nil
}

// https://developer.valvesoftware.com/wiki/Steam_Web_API#GetPlayerSummaries_.28v0002.29
func GetPlayersInfo(steamids []string) (map[string]*PlayerInfo, error) {
	url := "http://api.steampowered.com/ISteamUser/GetPlayerSummaries/v0002/?key=" +
		steamApiKey + "&steamids=" + strings.Join(steamids, ",")

	profiles := make(map[string]*PlayerInfo)

	response, err := getJsonFromUrl(url)
	if err != nil {
		return nil, err
	}

	ps, err := response.Get("response").Get("players").Array()
	if err != nil {
		return nil, err
	}

	for _, _elem := range ps {
		elem := _elem.(map[string]interface{})
		player := new(PlayerInfo)

		pErr := player.Parse(elem)

		if pErr != nil {
			return nil, pErr
		}

		profiles[player.Steamid] = player
	}

	return profiles, nil
}

func GetPlayerInfo(steamid string) (*PlayerInfo, error) {
	url := "http://api.steampowered.com/ISteamUser/GetPlayerSummaries/v0002/?key=" +
		steamApiKey + "&steamids=" + steamid

	player := new(PlayerInfo)

	response, err := getJsonFromUrl(url)
	if err != nil {
		return nil, err
	}

	ps, err := response.Get("response").Get("players").Array()
	if err != nil {
		return nil, err
	}

	if len(ps) == 1 {
		elem := ps[0].(map[string]interface{})
		pErr := player.Parse(elem)

		if pErr != nil {
			return nil, pErr
		}
	}

	return player, nil
}
