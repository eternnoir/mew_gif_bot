package redis

import (
	log "github.com/Sirupsen/logrus"
	"github.com/eternnoir/mew_gif_bot/mew/utils"
	"github.com/garyburd/redigo/redis"
	"math/rand"
	"strings"
)

type RedisStore struct {
	Addr   string
	Passwd string
}

func NewRedisStore(addr, passwd string) (*RedisStore, error) {
	ret := &RedisStore{Addr: addr, Passwd: passwd}
	c, err := ret.getConn(0)
	if err != nil {
		return nil, err
	}
	defer c.Close()
	_, err = redis.String(c.Do("PING"))
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func (rs *RedisStore) getConn(db int) (redis.Conn, error) {
	c, err := redis.Dial("tcp", rs.Addr)

	if err != nil {
		return nil, err
	}

	return c, err
}

func (rs *RedisStore) Put(key, fileid string) error {
	fileid = key + "::" + fileid
	log.Debugf("Redis put %s  %s", key, fileid)
	c, err := rs.getConn(0)
	if err != nil {
		return err
	}
	defer c.Close()
	_, err = c.Do("ZADD", "ALLGIF", 0, fileid)
	if err != nil {
		return err
	}
	return nil
}

func (rs *RedisStore) Get(key string) ([]string, error) {
	allkeys, err := rs.GetAll()
	if err != nil {
		return nil, err
	}
	ret := []string{}
	for _, k := range allkeys {
		if strings.HasPrefix(k, key) {
			ret = append(ret, k)
		}
	}
	return ret, nil
}

func (rs *RedisStore) GetAll() ([]string, error) {
	c, err := rs.getConn(0)
	if err != nil {
		return nil, err
	}
	defer c.Close()
	c.Do("INCR", "COUNT:TOTALQUERY")
	return redis.Strings(c.Do("ZREVRANGE", "ALLGIF", 0, -1))
}

func (rs *RedisStore) Hint(key string) {
	c, err := rs.getConn(0)
	if err != nil {
		return
	}
	defer c.Close()
	c.Do("INCR", "COUNT:TOTALSEND")
	_, err = c.Do("ZINCRBY", "ALLGIF", 1, key)
	if err != nil {
		log.Error(err)
	}
}

func (rs *RedisStore) Reset() error {
	c, err := rs.getConn(0)
	if err != nil {
		return err
	}
	defer c.Close()
	allkeys, err := rs.GetAll()
	if err != nil {
		return err
	}
	for _, k := range allkeys {
		c.Do("ZADD", "ALLGIF", rand.Intn(10), k)
	}
	return nil
}

func (rs *RedisStore) GetStatus() (*utils.Status, error) {
	c, err := rs.getConn(0)
	if err != nil {
		return nil, err
	}
	ret := &utils.Status{}
	allkeys, err := rs.GetAll()
	if err != nil {
		return nil, err
	}
	ret.TotalGifs = int64(len(allkeys))
	totalQuerys, err := redis.Int64(c.Do("GET", "COUNT:TOTALQUERY"))
	if err != nil {
		return nil, err
	}
	totalSends, err := redis.Int64(c.Do("GET", "COUNT:TOTALSEND"))
	if err != nil {
		return nil, err
	}
	ret.TotalQuerys = totalQuerys
	ret.TotalSend = totalSends

	return ret, nil
}
