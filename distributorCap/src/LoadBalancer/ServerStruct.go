package LoadBalancer

type ServerStruct struct {
	HostName string `json:"hostName"`
	Port     string `json:"port"`
	Scheme   string `json:"Scheme"`
}
