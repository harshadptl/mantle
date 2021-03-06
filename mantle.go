package mantle

import (
	"github.com/goibibo/mantle/backends"
)

//only strings are supported
type Mantle interface {
	Get(key string) string
	Set(key string, value interface{}) bool
	Delete(keys ...interface{}) int
	Setex(key string, duration int, value interface{}) bool
	MGet(keys ...interface{}) []string
	MSet(keyValMap map[string]interface{}) bool
	Expire(key string, duration int) bool
	Execute(cmd string, args ...interface{}) (interface{}, error)
}

//helper func
func redisConns(settings mantle.PoolSettings) *mantle.Redis {
	redis := &mantle.Redis{}
	redis.Configure(settings)
	return redis
}

func memcacheConns(settings mantle.PoolSettings) *mantle.Memcache {
	redis := &mantle.Memcache{}
	redis.Configure(settings)
	return redis
}

//generic pool settings
func getSettings(o *Orm) mantle.PoolSettings {
	return mantle.PoolSettings{
		HostAndPorts: o.HostAndPorts,
		Capacity:     o.Capacity,
		MaxCapacity:  o.Capacity,
		Options:      o.Options}

}

//this struct is exported
type Orm struct {
	//redis|memcache|cassandra
	Driver string
	//arrays of ip:port,ip:port
	HostAndPorts []string
	//pool size
	Capacity int
	//any other options thats needed for creating a connection
	Options map[string]string
}

//mantle is a wrapper for many nosql dbs
func (o *Orm) New() Mantle {
	settings := getSettings(o)
	if o.Driver == "memcache" {
		return memcacheConns(settings)
	} else {
		return redisConns(settings)
	}
}

//override mantle and get a redis client
func (o *Orm) GetRedisConn() (*mantle.RedisConn, error) {
	settings := getSettings(o)
	redisPool := redisConns(settings)
	return redisPool.GetClient()
}
