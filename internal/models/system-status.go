package models

type SystemStatus struct {
	CPUPercent    float64
	MemoryPercent float64
	MemoryUsed    uint64
	MemoryTotal   uint64
	Temperatures  []SystemTemperature
	SerialData    string
}

type SystemTemperature struct {
	Sensor  string
	Celsius float64
}
