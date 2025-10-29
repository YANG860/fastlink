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

func GetLinkSourceFromCache(shortUrl string) (string, error) {
	source, err := RedisClient.Get(Ctx, "short_url:"+shortUrl+":source").Result()
	if err != nil {
		return "", err
	}

	return source, nil
}

func SetLinkToCache(link *Link) error {

	err := RedisClient.Set(Ctx, "short_url:"+link.ShortUrl+":source", link.SourceUrl, 30*time.Minute).Err()
	if err != nil {
		return err
	}

	err = RedisClient.Set(Ctx, "short_url:"+link.ShortUrl+":link_id", link.ID, 30*time.Minute).Err()
	if err != nil {
		return err
	}
	err = RedisClient.Set(Ctx, "short_url:"+link.ShortUrl+":expire_at", link.ExpireAt.Unix(), 30*time.Minute).Err()
	if err != nil {
		return err
	}
	err = RedisClient.Set(Ctx, "short_url:"+link.ShortUrl+":click_count", link.ClickCount, 30*time.Minute).Err()
	if err != nil {
		return err
	}
	return nil
}
