package main

import (
	"log"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/youtube/vitess/go/pools"

	"golang.org/x/net/context"
	"sync"

	"github.com/rainycape/memcache"
	"os"
	"strconv"
)

func main() {

	var err error
	time.Sleep(time.Second * 5)
	var parallel uint64
	parallel, err = strconv.ParseUint(os.Getenv("PARALLEL"), 10, 64)
	if err != nil {
		panic(err)
	}
	var start, end time.Time
	var wg *sync.WaitGroup

	var iter uint64
	iter, err = strconv.ParseUint(os.Getenv("ITERATIONS"), 10, 64)
	if err != nil {
		panic(err)
	}

	calcRps:= func(start time.Time , end time.Time){
			var x float64
		x = float64((float64(end.UnixNano()-start.UnixNano())/float64(1000000000)))
		log.Println("requests per second", float64(parallel*iter)/x)
	}

	//// tarantool
	//tntclient, err := getTntClient()
	//if err != nil {
	//	panic(err)
	//}
	//
	//log.Println("Start tarantool test bench")
	//wg = &sync.WaitGroup{}
	//start = time.Now()
	//for i := uint64(0); i < parallel; i++ {
	//	wg.Add(1)
	//	go tntWorker(tntclient, wg, iter)
	//}
	//wg.Wait()
	//end = time.Now()
	//calcRps(start, end)


	// redis
	p := pools.NewResourcePool(func() (pools.Resource, error) {
		c, err := redis.Dial("tcp", "redis:6379")
		return ResourceConn{c}, err
	}, 100, 1000000, time.Minute)
	defer p.Close()

	redisSetHelper(p)



	log.Println("Start redis test bench")
	wg = &sync.WaitGroup{}
	start = time.Now()
	for i := uint64(0); i < parallel; i++ {
		wg.Add(1)
		go redisWorker(p, wg, iter)
	}
	wg.Wait()
	end = time.Now()
	calcRps(start, end)


	//memcached
	memcachedClient, err := memcache.New("memcached:11211")
	memcachedClient.SetTimeout(time.Second * 10)
	memcachedClient.SetMaxIdleConnsPerAddr(100)
	memcachedClient.Set(&memcache.Item{Key: "key", Value: []byte("data")})
	if err != nil {
		panic(err)
	}

	log.Println("Start memcached test bench")
	wg = &sync.WaitGroup{}
	start = time.Now()
	for i := uint64(0); i < parallel; i++ {
		wg.Add(1)
		go memcachedWorker(memcachedClient, wg, iter)
	}
	wg.Wait()
	end = time.Now()
	calcRps(start, end)


}

type ResourceConn struct {
	redis.Conn
}

func (r ResourceConn) Close() {
	r.Conn.Close()
}

func memcachedWorker(m *memcache.Client, wg *sync.WaitGroup, n uint64) {
	for i := uint64(0); i < n; i++ {

		_, err := m.Get("key")

		if err != nil {
			panic(err)
		}
	}
	wg.Done()
}

func redisSetHelper(p *pools.ResourcePool) {
	ctx := context.TODO()

	r, err := p.Get(ctx)
	if err != nil {
		panic(err)
	}
	c := r.(ResourceConn)

	_, err = c.Do("SET", "key", "data")
	if err != nil {
		panic(err)
	}
	p.Put(r)
}

func redisWorker(p *pools.ResourcePool, wg *sync.WaitGroup, n uint64) {
	for i := uint64(0); i < n; i++ {
		ctx := context.TODO()

		r, err := p.Get(ctx)
		if err != nil {
			panic(err)
		}
		c := r.(ResourceConn)

		_, err = c.Do("GET", "key")
		if err != nil {
			panic(err)
		}
		p.Put(r)

	}
	wg.Done()
}
