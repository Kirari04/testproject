package t

type ProxyStatusResponse struct {
	Proxy  string  `json:"proxy"`
	Server string  `json:"server"`
	State  string  `json:"state"`
	Value  float64 `json:"value"`
}

type ProxyStatus struct {
	FrontendId uint `json:"frontend_id"`
	Servers    []ProxyStatusServer

	BytesInTotal  *float64 `json:"bytes_in_total"`
	BytesOutTotal *float64 `json:"bytes_out_total"`

	RequestsTotal       *float64 `json:"requests_total"`
	ResponsesTotal1xx   *float64 `json:"responses_total_1xx"`
	ResponsesTotal2xx   *float64 `json:"responses_total_2xx"`
	ResponsesTotal3xx   *float64 `json:"responses_total_3xx"`
	ResponsesTotal4xx   *float64 `json:"responses_total_4xx"`
	ResponsesTotal5xx   *float64 `json:"responses_total_5xx"`
	ResponsesTotalOther *float64 `json:"responses_total_other"`
}

type ProxyStatusServer struct {
	ServerId uint   `json:"server_id"`
	Address  string `json:"address"`

	HANA     float64 `json:"hana"`
	SOCKERR  float64 `json:"sockerr"`
	L4OK     float64 `json:"l4ok"`
	L4TOUT   float64 `json:"l4tout"`
	L4CON    float64 `json:"l4con"`
	L6OK     float64 `json:"l6ok"`
	L6TOUT   float64 `json:"l6tout"`
	L6RSP    float64 `json:"l6rsp"`
	L7TOUT   float64 `json:"l7tout"`
	L7RSP    float64 `json:"l7rsp"`
	L7OK     float64 `json:"l7ok"`
	L7OKC    float64 `json:"l7okc"`
	L7STS    float64 `json:"l7sts"`
	PROCERR  float64 `json:"procerr"`
	PROCTOUT float64 `json:"proctout"`
	PROCOK   float64 `json:"procok"`
}

type ProxyStatusStats struct {
	FrontendId          uint     `json:"frontend_id"`
	BytesInTotal        *float64 `json:"bytes_in_total"`
	BytesOutTotal       *float64 `json:"bytes_out_total"`
	RequestsTotal       *float64 `json:"requests_total"`
	ResponsesTotal1xx   *float64 `json:"responses_total_1xx"`
	ResponsesTotal2xx   *float64 `json:"responses_total_2xx"`
	ResponsesTotal3xx   *float64 `json:"responses_total_3xx"`
	ResponsesTotal4xx   *float64 `json:"responses_total_4xx"`
	ResponsesTotal5xx   *float64 `json:"responses_total_5xx"`
	ResponsesTotalOther *float64 `json:"responses_total_other"`
}
