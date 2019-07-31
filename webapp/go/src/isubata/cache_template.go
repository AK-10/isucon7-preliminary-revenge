package main

import (
	"encoding/json"
	"fmt"

	"github.com/gomodule/redigo/redis"
)

// example
func setNum(key string, value int64) error {
	conn, err := redis.Dial("tcp", "localhost:6379")
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = conn.Do("SET", key, value)
	return err
}

func getNum(key string) (int64, error) {
	conn, err := redis.Dial("tcp", "localhost:6379")
	if err != nil {
		return -1, err
	}
	defer conn.Close()

	return redis.Int64(conn.Do("GET", key))
}

// 取得時のキャスト
// redis.Int
// redis.Ints []int
// redis.Int64
// redis.Int64s []int64
// redis.Float64
// redis.Float64s []float64
// redis.String
// redis.Strings []string
// redis.Bool
// redis.Bytes []byte

// 取得時の呼び出し側のエラー処理
func callGetFromRedis() {
	v, err := getNum("num")
	if err == redis.ErrNil {
		// 値がない
	}
	if err != nil {
		// それ以外のエラー
	}
	fmt.Println(v)
}

func setStruct(key string, value interface{}) error {
	conn, err := redis.Dial("tcp", "localhost:6379")
	if err != nil {
		return err
	}
	defer conn.Close()

	serialized, _ := json.Marshal(value)

	_, err = conn.Do("SET", key, serialized)
	return err
}

func getStruct(key string) (interface{}, error) {
	conn, err := redis.Dial("tcp", "localhost:6379")
	if err != nil {
		return -1, err
	}
	defer conn.Close()
	data, _ := redis.Bytes(conn.Do("GET", key))

	// JSON to struct
	if data != nil {
		deserialized := new(User)
		json.Unmarshal(data, deserialized)
		return deserialized, nil
	}
	return nil, err

}
