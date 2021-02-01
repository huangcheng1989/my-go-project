package dbredis

import (
	goredislib "github.com/go-redis/redis"
	"github.com/go-redsync/redsync"
	"github.com/go-redsync/redsync/redis/goredis/v8"
	"strings"
	"time"
)

var rediSync *redsync.Redsync

var redisCluster *RedisCluster

type RedisCluster struct {
	Client *goredislib.ClusterClient
}

func InitRedisCluster(env string) {
	var addressString string
	if env == "prod" {
		addressString = "vskit-redis-major.1rhy6q.clustercfg.euw1.cache.amazonaws.com:6379"
	} else if env == "test" {
		addressString = "db.mylichking.com:7000,db.mylichking.com:7001,db.mylichking.com:7002,db.mylichking.com:7003,db.mylichking.com:7004,db.mylichking.com:7005"
	} else if env == "dev" {
		addressString = "10.200.50.49:7001,10.200.50.49:7002,10.200.50.49:7003,10.200.50.49:7004,10.200.50.49:7005,10.200.50.49:7006"
	}

	redisAddress := strings.Split(addressString, ",")
	opt := &goredislib.ClusterOptions{
		Addrs:           redisAddress,
		DialTimeout:     60 * time.Second,
		ReadTimeout:     60 * time.Second,
		WriteTimeout:    60 * time.Second,
		PoolSize:        2000,
		Password:        "",
		MaxRetries:      3,
		MinRetryBackoff: -1,
		MaxRetryBackoff: -1,
	}
	client := goredislib.NewClusterClient(opt)
	redisCluster = &RedisCluster{Client: client}
	pool := goredis.NewPool(client)
	rediSync = redsync.New(pool)
}

func GetRedisCluster() *goredislib.ClusterClient {
	return redisCluster.Client
}
