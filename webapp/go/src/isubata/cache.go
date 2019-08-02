package main

import (
	"errors"
	"strconv"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/patrickmn/go-cache"
)

const (
	haveReadPrefix   = "HAVE-READ-"
	messageNumPrefix = "MSG-NUM-"
	userNumKey       = "USER-NUM"
)

var (
	c = cache.New(5*time.Minute, 10*time.Minute)
)

// var (
// 	redisPool = &redis.Pool{
// 		MaxIdle:     3,
// 		MaxActive:   18,
// 		IdleTimeout: 240 * time.Second,
// 		Dial:        func() (redis.Conn, error) { return redis.Dial("tcp", "127.0.0.1:6379") },
// 	}
// )

func makeMSGNumKey(chanID int64) string {
	return messageNumPrefix + strconv.Itoa(int(chanID))
}

func makeHaveReadKey(chanID, userID int64) string {
	return haveReadPrefix + strconv.Itoa(int(chanID)) + "-" + strconv.Itoa(int(userID))
}

func setLastIDtoGC(chanID, userID, lastID int64) {
	c.Set(makeHaveReadKey(chanID, userID), lastID, cache.DefaultExpiration)
}

func getLastIDFromGC(chanID, userID int64) (int64, error) {
	if lastID, found := c.Get(makeHaveReadKey(chanID, userID)); found {
		return lastID.(int64), nil
	}
	return -1, errors.New("go cache: key not found")
}

func setLastIDtoRedis(chanID, userID, lastID int64) error {
	// conn := redisPool.Get()
	// defer conn.Close()
	conn, err := redis.Dial("tcp", "localhost:6379")
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	key := makeHaveReadKey(chanID, userID)
	_, err = conn.Do("SET", key, strconv.Itoa(int(lastID)))
	return err
}

func getLastIDFromRedis(chanID, userID int64) (int64, error) {
	// conn := redisPool.Get()
	// defer conn.Close()
	conn, err := redis.Dial("tcp", "localhost:6379")
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	key := makeHaveReadKey(chanID, userID)
	lastID, err := redis.Int64(conn.Do("GET", key))
	return lastID, err
}

func setMessageNumToGC(chanID, num int64) {
	c.Set(makeMSGNumKey(chanID), num, cache.DefaultExpiration)
}

func getMessageNumFromGC(chanID int64) (int64, error) {
	if lastID, found := c.Get(makeMSGNumKey(chanID)); found {
		num := lastID.(int64)
		return num, nil
	}
	return -1, errors.New("go-cache: key not found")
}

func incrementMessageNumAtGC(chanID int64) error {
	oldNum, err := getMessageNumFromGC(chanID)
	if err != nil {
		return err
	}
	setMessageNumToGC(chanID, oldNum+1)
	return nil
}

func decrementMessageNumAtGC(chanID int64) error {
	oldNum, err := getMessageNumFromGC(chanID)
	if err != nil {
		return err
	}
	setMessageNumToGC(chanID, oldNum-1)
	return nil
}

func setUserNumToGC(num int64) {
	c.Set(userNumKey, num, cache.DefaultExpiration)
}

func getUserNumFromGC() (int64, error) {
	if num, found := c.Get(userNumKey); found {
		return num.(int64), nil
	}
	return -1, errors.New("go-cache: key not found")
}

func incrementUserNumtoGC() error {
	oldNum, err := getUserNumFromGC()
	if err != nil {
		return err
	}
	setUserNumToGC(oldNum + 1)
	return nil
}

func incrementMessageNumtoRedis(chanID int64) error {
	oldNum, err := getMessageNumFromRedis(chanID)
	if err != nil {
		return err
	}
	if err = setMessageNumFromRedis(chanID, oldNum+1); err != nil {
		return err
	}
	return nil
}

func decrementMessageNumtoRedis(chanID int64) error {
	oldNum, err := getMessageNumFromRedis(chanID)
	if err != nil {
		return err
	}
	if err = setMessageNumFromRedis(chanID, oldNum-1); err != nil {
		return err
	}
	return nil
}

func setMessageNumFromRedis(chanID, num int64) error {
	// conn := redisPool.Get()
	// defer conn.Close()
	conn, err := redis.Dial("tcp", "localhost:6379")
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	key := makeMSGNumKey(chanID)
	_, err = conn.Do("SET", key, strconv.Itoa(int(num)))
	return err
}

func getMessageNumFromRedis(chanID int64) (int64, error) {
	// conn := redisPool.Get()
	// defer conn.Close()
	conn, err := redis.Dial("tcp", "localhost:6379")
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	key := makeMSGNumKey(chanID)
	num, err := redis.Int64(conn.Do("GET", key))
	return num, err
}

func incrementUserNumtoRedis() error {
	oldNum, err := getUserNumFromRedis()
	if err != nil {
		return err
	}
	if err = setUserNumFromRedis(oldNum + 1); err != nil {
		return err
	}
	return nil
}

func setUserNumFromRedis(num int64) error {
	// conn := redisPool.Get()
	// defer conn.Close()
	conn, err := redis.Dial("tcp", "localhost:6379")
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	key := userNumKey
	_, err = conn.Do("SET", key, strconv.Itoa(int(num)))
	return err
}

func getUserNumFromRedis() (int64, error) {
	// conn := redisPool.Get()
	// defer conn.Close()
	conn, err := redis.Dial("tcp", "localhost:6379")
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	key := userNumKey
	num, err := redis.Int64(conn.Do("GET", key))
	return num, err
}
