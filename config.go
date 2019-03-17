package main

import (
	"sync"
	"time"

	"github.com/BurntSushi/toml"
)

var (
	config *Config
	once   sync.Once
)

// try to use this
func initConfig() {
	if _, err := toml.DecodeFile("./config.toml", &config); err != nil {
		panic(err)
	}
}

// Config reading from `config.toml`
type Config struct {
	Redis redisConfig
	RPool rPoolConfig `toml:"RedisPool"`
	PPool pPoolConfig `toml:"ProxyPool"`
}

type redisConfig struct {
	// Address's format is `host:port`
	// `host` can be ip or host(`service`) name of `redis server`
	// `port` exposed by `redis server`
	Address        string
	Password       string   // `Password` in redis conf
	DB             int      `toml:"Database"` // `DB` the index of redis database
	ConnectTimeout duration // connect timeout
	ReadTimeout    duration // read timeout
	WriteTimeout   duration // write timeout
}

type rPoolConfig struct {
	// Maximum number of idle connections in the pool
	MaxIdle int

	// Maximum number of connections allocated by the pool at a given time
	MaxActive int

	// Close connections after remaining idle for this duration
	IdleTimeout duration

	// If Wait is true and the pool is at the MaxActive limit, then Get() waits
	// for a connection to be returned to the pool before returning.
	Wait bool

	// Close connections older than this duration
	MaxConnLifetime duration
}

type pPoolConfig struct {
	Address        string   // ip:port
	MinScore       int      // used by sorted set
	MaxScore       int      // used ...
	DefaultScore   int      // used ...
	RedisZKey      string   // key for sorted set
	ScanInterval   duration // used by scheduler
	ScrapeInterval duration // used...
	IncrAfterScan  int
	DecrAfterScan  int
}

type duration struct {
	time.Duration
}

func (d *duration) UnmarshalText(text []byte) (err error) {
	d.Duration, err = time.ParseDuration(string(text))
	return
}

func getConfig() *Config {
	// var c Config
	once.Do(func() {
		if _, err := toml.DecodeFile("./config.toml", &config); err != nil {
			panic(err)
		}
	})
	return config
}
