package t

type ProxyStatus struct {
	Proxy  string  `json:"proxy"`
	Server string  `json:"server"`
	State  string  `json:"state"`
	Value  float64 `json:"value"`
}
