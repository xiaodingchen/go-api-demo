package xredis

import (
   "errors"
   "fmt"
   "github.com/go-redis/redis"
   "github.com/spf13/viper"
   "test.local/pkg/utils"
   "time"
)

const errConfigKeyNil = "redis config nil: %s"

func NewClient(name string) (client *redis.Client, err error) {
   key := "redis." + name
   if !viper.IsSet(key) {
      err = errors.New(fmt.Sprintf(errConfigKeyNil, key))
      return
   }
   client = redis.NewClient(
      &redis.Options{
         Addr:         viper.GetString(key + ".addr"),
         Password:     viper.GetString(key + ".password"),
         DB:           viper.GetInt(key + ".db"),
         DialTimeout:  viper.GetDuration(key+".dialTimeout") * time.Millisecond,
         ReadTimeout:  viper.GetDuration(key+".readTimeout") * time.Millisecond,
         WriteTimeout: viper.GetDuration(key+".writeTimeout") * time.Millisecond,
         MaxRetries:   viper.GetInt(key + ".maxRetries"),
         PoolSize:     viper.GetInt(key + ".poolSize"),
         MinIdleConns: viper.GetInt(key + ".minIdleConns"),
      })

   utils.SetRedis(name, client)
   return
}

func loop(client *redis.Client) {
   go func() {
      for {
         client.Ping()
         time.Sleep(3 * time.Second)
      }
   }()
}