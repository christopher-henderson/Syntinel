package LoadBalancer

import "net/url"

//todo: Read load balancing algo from config file
type loadBalancer interface {
	getNext() url.URL
}
