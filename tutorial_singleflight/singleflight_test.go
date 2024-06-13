package tutorial_singleflight

import (
	"fmt"
	"testing"
	"time"

	"github.com/sky1309/log"
	"golang.org/x/sync/singleflight"
)

func TestSingleFlight(t *testing.T) {
	group := singleflight.Group{}

	for i := 0; i < 2; i++ {
		go func(reqId int) {
			value, err, shared := group.Do("query", func() (any, error) {
				fmt.Println("===do query")
				time.Sleep(time.Second)
				fmt.Println("===finish query")
				return 100, nil
			})
			log.Info("i=%d, value=%+v, err=%v, shared=%v", reqId, value, err, shared)
		}(i)
	}

	time.Sleep(time.Second * 2)
}

func TestSingleFightChan(t *testing.T) {
	group := singleflight.Group{}

	queryFn := func() (interface{}, error) {
		fmt.Println("===do query")
		time.Sleep(time.Second)
		fmt.Println("===finish query")
		return 100, nil
	}

	for i := 0; i < 5; i++ {
		go func(reqId int) {
			ch := group.DoChan("query", queryFn)
			log.Info("idx=%d, result=%+v", reqId, (<-ch).Val)
		}(i)
	}

	time.Sleep(time.Second * 2)
}
