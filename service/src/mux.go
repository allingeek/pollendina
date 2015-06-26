package main

import (
	"net/http"
	"regexp"
	"github.com/pollendina/logger"
)

type MuxRoute struct {
        pattern *regexp.Regexp
        handler http.Handler
}

type RegexHandler struct {
        rs []*MuxRoute
}

func (rh *RegexHandler) Handler(p *regexp.Regexp, h http.Handler) {
        rh.rs = append(rh.rs, &MuxRoute{p, h})
}

func (rh *RegexHandler) HandleFunc(p *regexp.Regexp, h func(http.ResponseWriter, *http.Request)) {
        rh.rs = append(rh.rs, &MuxRoute{p, http.HandlerFunc(h)})
}

func (rh *RegexHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
        for _, route := range rh.rs {
                if route.pattern.MatchString(r.URL.Path) {
                        route.handler.ServeHTTP(w, r)
                        return
                }
        }
        logger.Warning.Printf("Route not found: %s", r.URL.Path)
        http.NotFound(w, r)
}
