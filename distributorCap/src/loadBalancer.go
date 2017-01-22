package main



import(
       "log"
        "math/rand"
        "net/http"
        "net/http/httputil"
        "net/url"
)

func LoadBalancer(sites[] *url.URL)*httputil.ReverseProxy{
	balancer := func(req *http.Request){
		target := sites[rand.Int()%len(sites)]
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.URL.Path = target.Path
	}
	return &httputil.ReverseProxy{Director: balancer}
}


func main() {
        proxy := LoadBalancer([]*url.URL{
                {
                        Scheme: "http",
                        Host:   "localhost:9091",
                },
                {
                        Scheme: "http",
                        Host:   "localhost:9092",
                },
        })
        log.Fatal(http.ListenAndServe(":9090", proxy))
}
