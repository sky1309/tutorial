package tutorial_redis

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/sky1309/log"
)

func TestRedis_ListMockMessageQueue(t *testing.T) {
	rdb := connect()
	defer rdb.Close()

	productor := func(ctx context.Context) {
		tick := time.NewTicker(time.Second)
		i := 0
		for {
			select {
			case <-tick.C:
				rdb.LPush(ctx, "logs", fmt.Sprintf("idx%d", i))
				i++
			case <-ctx.Done():
				log.Info("quit productor!")
				return
			}
		}
	}

	consumer := func(ctx context.Context) {
		tick := time.NewTicker(time.Microsecond * 200)
		for {
			select {
			case <-tick.C:
				s, err := rdb.RPop(ctx, "logs").Result()
				if err != nil {
					if err == redis.Nil {
						continue
					}
					log.Error("consumer rpop err %v", err)
					return
				}
				log.Info("consumer get log %s", s)
			case <-ctx.Done():
				log.Info("quit consumer!")
				return
			}
		}
	}

	ctx, _ := context.WithTimeout(context.Background(), time.Second*10)

	go productor(ctx)
	go consumer(ctx)

	<-ctx.Done()
	log.Info("finish")
}

func TestRedis_SimpleDistributedLock(t *testing.T) {
	rdb := connect()
	defer rdb.Close()

	ctx := context.Background()
	lock := func(key, value string) (bool, error) {
		// retry 10 times
		for i := 1; i <= 10; i++ {
			success, err := rdb.SetNX(ctx, key, value, time.Second*3).Result()
			if err != nil {
				return false, err
			}

			if success {
				log.Info("lock key=%s, value=%s success!", key, value)
				return true, nil
			}

			time.Sleep(time.Millisecond * time.Duration(100*i))
			log.Warn("lock key=%s retry=%d", key, i)
		}

		log.Info("lock key=%s fail", key)
		return false, nil
	}

	// use redis script
	unlockScript := redis.NewScript(`
		local key = KEYS[1]
		local compareValue = ARGV[1]

		local value = redis.call("GET", key)
		if value ~= compareValue then
			return false
		end

		redis.call("DEL", key)
		return true
	`)

	unlock := func(key, value string) (bool, error) {
		result, err := unlockScript.Run(ctx, rdb, []string{key}, []string{value}).Bool()
		log.Info("unlock key=%s, value=%s, result=%v, err=%v", key, value, result, err)
		return result, err
	}

	lock("update_name", "1")

	time.AfterFunc(time.Second, func() {
		// fail
		unlock("update_name", "ddddddd")

		// success
		unlock("update_name", "1")
	})

	lock("update_name", "1")
}
