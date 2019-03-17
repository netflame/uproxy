package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
)

type handler struct {
	c *client // client for redis-server
}

func defaultHandler() *handler {
	return &handler{c: defaultClient()}
}

func (h *handler) Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	resp := `
	# >> Welcome to Proxy Pool <<< #
	# ---------------------------- #
	# >>>>>> Avalibale Path <<<<<< #
	# >> /all - /random - /size << #
	`
	fmt.Fprint(w, resp)
}

func (h *handler) All(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if resp, err := h.c.all(config.PPool.RedisZKey); err != nil {
		fmt.Fprintf(w, "Got error: %s\n", err)
	} else {
		fmt.Fprintf(w, strings.Join(resp, "\n"))
	}
}

func (h *handler) Random(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if resp, err := h.c.random(config.PPool.RedisZKey); err != nil {
		fmt.Fprintf(w, "Got error: %s\n", err)
	} else {
		fmt.Fprintf(w, resp)
	}
}

func (h *handler) Size(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if resp, err := h.c.size(config.PPool.RedisZKey); err != nil {
		fmt.Fprintf(w, "Got error: %s\n", err)
	} else {
		fmt.Fprintf(w, strconv.Itoa(resp))
	}
}
