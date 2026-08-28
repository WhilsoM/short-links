package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type LinksCache struct {
	redis *redis.Client
	ttl   time.Duration
}

func NewLinksCache(redis *redis.Client) *LinksCache {
	return &LinksCache{
		redis: redis,
		ttl:   15 * time.Minute,
	}
}

func (l *LinksCache) GetLinkByCode(ctx context.Context, code string) (string, error) {
	val, err := l.redis.Get(ctx, "link:"+code).Result()
	if err != nil {
		return "", err
	}

	return val, nil
}

func (l *LinksCache) SetLink(ctx context.Context, code, original_link string) error {
	err := l.redis.Set(ctx, "link:"+code, original_link, l.ttl).Err()
	if err != nil {
		return err
	}

	return nil
}
