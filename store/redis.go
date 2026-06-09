package store

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisStore struct{


	client *redis.Client

}

func NewRedisStore(addr, password string , db int)(*RedisStore, error){

	client := redis.NewClient(&redis.Options{
		Addr: addr,
		Password: password,
		DB: db,
	})

	ctx, cancel:= context.WithTimeout(context.Background(),5*time.Second)

	defer cancel()

	if err:=client.Ping(ctx).Err(); err!=nil{
		return nil,err
	}

	return &RedisStore{client: client},nil

}

