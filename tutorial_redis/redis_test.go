package tutorial_redis

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/sky1309/log"
)

func connect() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "foobared",
		DB:       0,
		OnConnect: func(ctx context.Context, cn *redis.Conn) error {
			fmt.Println("redis on connect!")
			return nil
		},
	})
}

func TestRedis_String(t *testing.T) {
	rdb := connect()
	defer rdb.Close()

	// set foo 1
	expiration := time.Second
	rdb.Set(context.Background(), "foo", 1, expiration)

	foo, err := rdb.Get(context.Background(), "foo").Result()
	fmt.Printf("foo=%s, err=%v\n", foo, err)

	tick := time.NewTicker(expiration)
	<-tick.C

	foo, err = rdb.Get(context.Background(), "foo").Result()
	fmt.Printf("foo=%s, err=%v\n", foo, err)
}

func TestRedis_Hash(t *testing.T) {
	rdb := connect()
	defer rdb.Close()

	ctx := context.Background()
	rdb.HSet(ctx, "online_nums", 1, 100, 2, 66)

	// get all
	onlineNumsMap := rdb.HGetAll(ctx, "online_nums").Val()
	log.Info("online_nums_map=%v", onlineNumsMap)

	rdb.HIncrBy(ctx, "online_nums", "1", 1)
	onlineNumsMap = rdb.HGetAll(ctx, "online_nums").Val()
	log.Info("online_nums_map1=%v", onlineNumsMap)

	// not exist
	num, err := rdb.HGet(ctx, "online_nums", "10").Result()
	log.Info("num=%s, err=%v", num, err)
}

func TestRedis_List(t *testing.T) {
	rdb := connect()
	defer rdb.Close()
	ctx := context.Background()

	logs := []string{
		"balabala",
		"lol",
		"fff",
	}

	for _, s := range logs {
		// add left
		rdb.LPush(ctx, "logs", s)
		// add right
		rdb.RPush(ctx, "logs", s)
	}

	query, err := rdb.LRange(ctx, "logs", 0, 100).Result()
	if err != nil {
		t.Fatal(err)
	}

	for i, s := range query {
		log.Info("query i=%d, s=%s", i, s)
	}

	rdb.Del(ctx, "logs")
}

func TestRedis_Set(t *testing.T) {
	rdb := connect()
	defer rdb.Close()
	ctx := context.Background()

	rdb.SAdd(ctx, "role_ids", 1, 2, 3, 4, 5)

	// get all members
	roleIds, err := rdb.SMembers(ctx, "role_ids").Result()
	if err != nil {
		t.Fatal(err)
	}
	log.Info("role_ids111=%v", roleIds)

	// remove member
	removeIds := []string{"1"}
	rdb.SRem(ctx, "role_ids", removeIds)
	roleIds, err = rdb.SMembers(ctx, "role_ids").Result()
	if err != nil {
		t.Fatal(err)
	}
	log.Info("role_ids222=%v", roleIds)
}

func TestRedis_ZSet(t *testing.T) {
	rdb := connect()
	defer rdb.Close()
	ctx := context.Background()

	members := []redis.Z{
		{Score: 100, Member: "1"},
		{Score: 102, Member: "2"},
		{Score: 66, Member: "3"},
	}

	// add rank
	rdb.ZAdd(ctx, "rank_level", members...)

	rank := rdb.ZRevRank(ctx, "rank_level", "2").Val()
	score := rdb.ZScore(ctx, "rank_level", "2").Val()
	log.Info("no1 rank=%d, score=%v", rank, score)
}
