package main

import "net/url"

type roundRobbin struct {
	servers []url.URL
}

func (r roundRobbin) GetNext() url.URL {
	//not sure if this is the most efficent way to do it
	//I think this copys the slice, let me know how it should be done
	tmp := r.servers[0]
	r.servers = append(r.servers[:0], r.servers[0+1:]...)
	r.servers = append(r.servers, tmp)
	return tmp
}
