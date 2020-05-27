package utils

import (
    "github.com/go-redis/redis"
)

var redisClientMap map[string]*redis.Client

const (
    DefaultRedis = "default"
)

func Redis(name string) *redis.Client{
    if redisClientMap == nil{
       return nil
    }
    //Log().Debug("return redis client", zap.Any("name", name))
    return redisClientMap[name]
}

func SetRedis(name string, client *redis.Client){
    //Log().Debug("set redis client", zap.Any("name", name))
    if redisClientMap == nil{
        redisClientMap = make(map[string]*redis.Client)
    }

    redisClientMap[name] = client
}
