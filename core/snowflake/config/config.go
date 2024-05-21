package config

type SnowFlakeSetting struct {
	MachineId int64 `json:"machineId" yaml:"machineId"`
}

var Setting *SnowFlakeSetting = &SnowFlakeSetting{
	MachineId: 13,
}
