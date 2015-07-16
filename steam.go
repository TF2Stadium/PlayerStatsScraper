package scraper

import (
	"encoding/json"
	"strconv"
)

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
