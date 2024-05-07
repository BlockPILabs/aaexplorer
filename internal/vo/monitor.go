package vo

type AddMonitorRequest struct {
	Network        string `json:"network"`
	MonitorAddress string `json:"monitorAddress"`
	Sign           string `json:"sign"`
	Message        string `json:"message"`
}

type AddMonitorResponse struct {
}

type RemoveMonitorRequest struct {
	Network        string `json:"network"`
	MonitorAddress string `json:"monitorAddress"`
	Sign           string `json:"sign"`
	Message        string `json:"message"`
}

type RemoveMonitorResponse struct {
}
