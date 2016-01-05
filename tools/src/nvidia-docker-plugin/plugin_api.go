// Copyright (c) 2015-2016, NVIDIA CORPORATION. All rights reserved.

package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path"

	"graceful"
)

const (
	socketName   = "nvidia.sock"
	acceptHeader = "application/vnd.docker.plugins.v1.1+json"
)

type plugin interface {
	implement() string
	register(*PluginAPI)
}

type PluginAPI struct {
	*graceful.HTTPServer

	plugins []plugin
}

func accept(handler http.Handler) http.Handler {
	f := func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Accept") != acceptHeader {
			w.WriteHeader(http.StatusNotAcceptable)
			return
		}
		handler.ServeHTTP(w, r)
	}
	return http.HandlerFunc(f)
}

func NewPluginAPI(prefix string) *PluginAPI {
	os.MkdirAll(prefix, 0700)

	a := &PluginAPI{
		HTTPServer: graceful.NewHTTPServer("unix", path.Join(prefix, socketName), accept),
	}
	a.Handle("POST", "/Plugin.Activate", a.activate)

	a.register(
		new(pluginVolume),
	)
	return a
}

func (a *PluginAPI) register(plugins ...plugin) {
	for _, p := range plugins {
		p.register(a)
		a.plugins = append(a.plugins, p)
	}
}

func (a *PluginAPI) activate(resp http.ResponseWriter, req *http.Request) {
	r := struct{ Implements []string }{}

	log.Println("Received handshake request")
	r.Implements = make([]string, len(a.plugins))
	for i, p := range a.plugins {
		r.Implements[i] = p.implement()
	}
	assert(json.NewEncoder(resp).Encode(r))
	log.Println("Plugins activated", r.Implements)
}
