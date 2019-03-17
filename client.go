package main

import (
	"math/rand"
	"time"

	"github.com/gomodule/redigo/redis"
)

// connection of redis server
type client struct {
	conn redis.Conn
}

func defaultClient() *client {
	return &client{conn: getConn()}
}

// ------ begin: useful func ------
// for handler use
func (c *client) all(k interface{}) (members []string, err error) {
	members, err = c.zrevrange(k, 0, -1)
	if err != nil {
		return nil, err
	}
	return
}

func (c *client) random(k interface{}) (member string, err error) {
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)
	members, err := c.all(k)
	if err != nil {
		return "", nil
	}
	size := len(members)
	if !(size > 0) {
		return "", nil
	}
	return members[r.Intn(size)], nil
}

func (c *client) size(k interface{}) (size int, err error) {
	size, err = c.zcard(k)
	if err != nil {
		return 0, nil
	}
	return
}

// used by redis pipeline to store
func (c *client) add(member interface{}) error {
	rzk, ds := config.PPool.RedisZKey, config.PPool.DefaultScore
	return c.zadd(rzk, ds, member)
}

// for scheduler scan
func (c *client) scan() {
	proxyURLs, err := c.all(config.PPool.RedisZKey)
	if err != nil {
		miniLog.error("(client.go:scan) Scan Failed")
		return
	}
	puch := make(chan miniChan, len(proxyURLs))
	for _, pURL := range proxyURLs {
		go validate(pURL, c, puch)
	}
	for i := 0; i < len(proxyURLs); i++ {
		<-puch
	}
}

func validate(pURL string, c *client, ch chan<- miniChan) {
	rzk := config.PPool.RedisZKey
	mins, maxs := config.PPool.MinScore, config.PPool.MaxScore
	incr, decr := config.PPool.IncrAfterScan, config.PPool.DecrAfterScan
	score, _ := c.zscore(rzk, pURL)

	err := checkProxyURL(pURL)
	// proxy not works
	if err != nil {
		c.zincrby(rzk, decr, pURL)
		// just remove it
		if score < mins {
			c.zrem(rzk, pURL)
		}
		miniLog.error("(client.go:validate) ", err)
	} else { // good proxy
		c.zincrby(rzk, incr, pURL)
		if score >= maxs {
			c.zadd(rzk, maxs, pURL)
		}
	}

	ch <- getMCElement()
}

// ------        end         ------

func (c *client) get(k interface{}) (v string, err error) {
	v, err = redis.String(c.conn.Do("GET", k))
	if err != nil {
		return "", err
	}
	return
}

func (c *client) set(k, v interface{}) error {
	_, err := c.conn.Do("SET", k, v)
	return err
}

func (c *client) del(k interface{}) error {
	_, err := c.conn.Do("DEL", k)
	return err
}

func (c *client) exists(k interface{}) (y bool, err error) {
	y, err = redis.Bool(c.conn.Do("EXISTS", k))
	if err != nil {
		return false, nil
	}
	return
}

// ksm: k->key, sm->score-member pairs
func (c *client) zadd(ksm ...interface{}) error {
	k := ksm[0]
	smPairs, smCount := ksm[1:], len(ksm[1:])
	if smCount < 2 || smCount%2 != 0 {
		miniLog.fatal("Unexpected `score member` pairs")
	}

	for i := 0; i < smCount; i += 2 {
		s, m := smPairs[i], smPairs[i+1]
		c.conn.Send("ZADD", k, s, m)
	}
	c.conn.Flush()
	_, err := c.conn.Receive()
	return err
}

func (c *client) zcard(k interface{}) (size int, err error) {
	size, err = redis.Int(c.conn.Do("ZCARD", k))
	if err != nil {
		return 0, err
	}
	return
}

func (c *client) zexists(key, member interface{}) (y bool, err error) {
	y = true
	if _, err = c.zscore(key, member); err != nil {
		y = false
		return
	}
	return
}

func (c *client) zincrby(key, increment, member interface{}) error {
	_, err := c.conn.Do("ZINCRBY", key, increment, member)
	return err
}

// km: k->key, m->members
func (c *client) zrem(km ...interface{}) error {
	if len(km) < 2 {
		miniLog.fatal("Unexpected `km` format, expect k, member, [member...]")
	}
	k, members := km[0], km[1:]
	for _, m := range members {
		c.conn.Send("ZREM", k, m)
	}
	c.conn.Flush()
	_, err := c.conn.Receive()
	return err
}

func (c *client) zrevrange(key, start, stop interface{}) (members []string, err error) {
	reply, err := redis.Values(c.conn.Do("ZREVRANGE", key, start, stop))
	if err != nil {
		return nil, err
	}
	if err = redis.ScanSlice(reply, &members); err != nil {
		return nil, err
	}
	return
}

func (c *client) zrevrangebyscore(key, max, min interface{}) (members []string, err error) {
	reply, err := redis.Values(c.conn.Do("ZREVRANGEBYSCORE", key, max, min))
	if err != nil {
		return nil, err
	}
	if err = redis.ScanSlice(reply, &members); err != nil {
		return nil, err
	}
	return
}

func (c *client) zscore(key, member interface{}) (score int, err error) {
	score, err = redis.Int(c.conn.Do("ZSCORE", key, member))
	if err != nil {
		return -1, nil
	}
	return
}

func (c *client) close() {
	c.conn.Close()
}
