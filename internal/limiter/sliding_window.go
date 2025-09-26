package limiter

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	LIMITER_CACHE_KEY = "limitr_window:%v"
)

type SlidingWindow struct {
	rdb *redis.Client
}

func New(rdb *redis.Client) *SlidingWindow {
	return &SlidingWindow{rdb: rdb}
}

// Allow reports whether a given key is allowed to make a request
func (l *SlidingWindow) Allow(
	ctx context.Context, hashedKey string, window time.Duration, maxReqs int,
) (bool, error) {
	now := time.Now().UnixMilli()
	rdbKey := fmt.Sprintf(LIMITER_CACHE_KEY, hashedKey)

	pipe := l.rdb.TxPipeline()

	pipe.ZAdd(ctx, rdbKey, redis.Z{
		Score:  float64(now),
		Member: now,
	})

	cutoff := now - window.Milliseconds()
	pipe.ZRemRangeByScore(ctx, rdbKey, "0", fmt.Sprint(cutoff))

	count := pipe.ZCard(ctx, rdbKey)

	pipe.Expire(ctx, rdbKey, window)

	if _, err := pipe.Exec(ctx); err != nil {
		return false, err
	}

	return count.Val() <= int64(maxReqs), nil
}
