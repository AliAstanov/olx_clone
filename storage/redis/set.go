package redis

import (
	"encoding/json"
	"log"
	"time"

	"github.com/AliAstanov/olx_clone/models"
	"github.com/go-redis/redis/v8"
	"golang.org/x/net/context"
)

func SetData(ctx context.Context, cli *redis.Client,userData *models.CreateUserReq) error {
	
	key := userData.Email

	bData, err := json.Marshal(userData)
	if err != nil {
		log.Printf("Failed to marshal user data: %v",err)
		return err
	}

	sData := string(bData)

	expiration := time.Minute * 5
	
	_, err = cli.Set(ctx,key,sData,expiration).Result()
	if err != nil {
		log.Printf("Failed to set data in Redis: %v",err)
		return err
	}
	return nil 
}
