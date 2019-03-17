package main

func init() {
	initConfig()
	initRedisPool()
	initMiniLog()
	// initServer()
}

func main() {
	toScrape, err := defaultSites()
	check(err)
	spider := &Spider{
		Name:     "Go spider",
		ToScrape: toScrape,
	}

	s := defaultServer()
	s.receive(spider)
	s.run()
}
