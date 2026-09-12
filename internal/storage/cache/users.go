package cache

import (
	"context"
	"encoding/json"
	"fmt"
	store "go-social/internal/storage"
	"time"

	"github.com/go-redis/redis/v8"
)

type UserStore struct {
	rdb *redis.Client
}

const UserExpTime = time.Minute

func (u UserStore) Get(ctx context.Context, userId string) (*store.User, error) {
	cacheKey := fmt.Sprintf("user-%v", userId)

	data, err := u.rdb.Get(ctx, cacheKey).Result()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	var user store.User
	if data != "" {
		err := json.Unmarshal([]byte(data), &user)
		if err != nil {
			return nil, err
		}
	}

	return &user, nil
}

func (u UserStore) Set(ctx context.Context, user *store.User) error {
	if user.Id == "" {
		return fmt.Errorf("user id not found!")
	}

	cacheKey := fmt.Sprintf("user-%v", user.Id)

	json, err := json.Marshal(user)
	if err != nil {
		return err
	}

	return u.rdb.SetEX(ctx, cacheKey, json, UserExpTime).Err()
}
