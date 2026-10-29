package usecase

import (
	"strings"
	"time"

	"github.com/Impisigmatus/alice-connect/internal/models"
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/mem"
	"github.com/shirou/gopsutil/v4/sensors"
)

type UseCase struct{}

func New() *UseCase {
	return &UseCase{}
}

func (*UseCase) SystemStatus() (*models.SystemStatus, error) {
	percent, err := cpu.Percent(time.Second, false)
	if err != nil {
		return nil, err
	}

	vm, err := mem.VirtualMemory()
	if err != nil {
		return nil, err
	}

	status := &models.SystemStatus{
		CPUPercent:    percent[0],
		MemoryPercent: vm.UsedPercent,
		MemoryUsed:    vm.Used,
		MemoryTotal:   vm.Total,
	}

	// Температура (ошибка здесь часто «частичная», поэтому смотрим на сами данные)
	temps, _ := sensors.SensorsTemperatures()
	for _, t := range temps {
		key := strings.ToLower(t.SensorKey)
		if strings.Contains(key, "core") || strings.Contains(key, "cpu") ||
			strings.Contains(key, "tctl") || strings.Contains(key, "package") {
			status.Temperatures = append(status.Temperatures, models.SystemTemperature{Sensor: t.SensorKey, Celsius: t.Temperature})
		}
	}

	return status, nil
}
