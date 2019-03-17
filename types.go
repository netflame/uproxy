package main

type site struct {
	Name string
	URL  string
}

type sites struct {
	Site []site
}

type httpbinIP struct {
	Origin string
}

type miniChan struct{}

func getMCElement() miniChan {
	return miniChan{}
}
