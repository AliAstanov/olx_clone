package redis

import (
	"context"
	"encoding/json"
	"log"

	"github.com/AliAstanov/olx_clone/models"
	"github.com/go-redis/redis/v8"
)

func GetData(ctx context.Context, cli *redis.Client, key string)(*models.CreateUserReq, error){
	
	result, err := cli.Get(ctx,key).Result()
	if err != nil {
		if err == redis.Nil{
			log.Printf("data not found in Redis for key: %s:",err)
		return nil, err
		}
		log.Printf("Failed to get data from Redis: %s:",err)
		return nil, err
	}

	var userData models.CreateUserReq
	err = json.Unmarshal([]byte(result),&userData)
	if err != nil {
		log.Printf("Failed to unmarshel data: %s",err)
		return nil, err
	}
	return &userData, nil

}