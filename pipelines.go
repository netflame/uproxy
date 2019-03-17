package main

type redisPL struct {
	c *client
}

func defaultRedisPL() *redisPL {
	return &redisPL{c: defaultClient()}
}

// get a redis conn every time when process proxies
// can auto close redis conn when done
func processProxies(proxies []*proxy, c chan<- miniChan) {
	rpl := defaultRedisPL()
	defer rpl.close()

	ch := make(chan miniChan, len(proxies))
	for _, i := range proxies {
		go rpl.processItem(i, ch)
	}
	for i := 0; i < len(proxies); i++ {
		<-ch
	}
	c <- getMCElement()
}

// store item into redis zset
func (rpl *redisPL) processItem(i item, ch chan<- miniChan) (err error) {
	defer func() {
		ch <- getMCElement()
	}()

	if err = checkItem(i); err != nil {
		miniLog.error("(pipelines.go:processItem) Invalid: ", i.String())
		return
	}

	// get valid item
	switch v := i.(type) {
	case *proxy:
		if err = rpl.c.add(v.String()); err != nil {
			miniLog.error("(pipelines.go:processItem) ", err)
		} else {
			miniLog.info("Valid: ", v.String())
		}
	default:
	}
	return
}

func (rpl *redisPL) close() {
	rpl.c.close()
}
