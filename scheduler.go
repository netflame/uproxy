package main

import "time"

func getScheduler(spider *Spider) *scheduler {
	// for creating spider in next scrape
	s := defaultScheduler()
	s.spider = spider
	return s
}

type scheduler struct {
	scanInterval   time.Duration
	scrapeInterval time.Duration
	spider         *Spider
}

func defaultScheduler() *scheduler {
	return &scheduler{
		scanInterval:   config.PPool.ScanInterval.Duration,
		scrapeInterval: config.PPool.ScrapeInterval.Duration,
		spider:         nil,
	}
}

func (s *scheduler) Start() {
	go s.Scrape()
	go s.Scan()
}

func (s *scheduler) Scrape() {
	scrapeTicker := time.NewTicker(s.scrapeInterval)

	// scrape first
	s.spider.Scrape()

	for range scrapeTicker.C {
		s.spider.Scrape()
	}
}

func (s *scheduler) Scan() {
	scanTicker := time.NewTicker(s.scanInterval)

	for range scanTicker.C {
		// get a client first
		c := defaultClient()
		c.scan()
		c.close()
	}
}
