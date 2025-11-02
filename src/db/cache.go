package db

import (
	"strconv"
	"time"
)

func GetVersionFromCache(userID int) (int, error) {
	var tokenID_str = strconv.Itoa(userID)
	version_str, err := RedisClient.Get(Ctx, "user:"+tokenID_str+":token_version").Result()
	if err != nil {
		return 0, err
	}
	version, err := strconv.Atoi(version_str)
	if err != nil {
		return 0, err
	}
	return version, nil
}

func SetVersionToCache(userID int, version int) error {
	var tokenID_str = strconv.Itoa(userID)
	err := RedisClient.Set(Ctx, "user:"+tokenID_str+":token_version", version, 30*time.Minute).Err()
	return err
}

func GetLinkFromCache(shortUrl string) (*Link, error) {
	var link Link
	err := RedisClient.HGetAll(Ctx, "link:"+shortUrl).Scan(&link)
	if err != nil {
		return nil, err
	}
	return &link, nil
}


func SetLinkToCache(link *Link) error {

	err := RedisClient.HSet(Ctx, "link:"+link.ShortUrl, link).Err()
	if err != nil {
		return err
	}

	return nil
}
