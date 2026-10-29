package usecase

import (
	"strings"
	"time"

	"github.com/Impisigmatus/alice-connect/internal/apis"
	"github.com/Impisigmatus/alice-connect/internal/config"
	"github.com/Impisigmatus/alice-connect/internal/models"
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/mem"
	"github.com/shirou/gopsutil/v4/sensors"
)

type UseCase struct {
	serialDriverAPI *apis.SerialDriverAPI
}

func New(cfg config.SerialDriverAPIConfig) *UseCase {
	return &UseCase{
		serialDriverAPI: apis.NewSerialDriverAPI(cfg),
	}
}

func (uc *UseCase) SystemStatus() (*models.SystemStatus, error) {
	percents, err := cpu.Percent(time.Second, false)
	if err != nil {
		return nil, err
	}

	vm, err := mem.VirtualMemory()
	if err != nil {
		return nil, err
	}

	status := &models.SystemStatus{
		CPUPercent:    average(percents),
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

	status.SerialData, err = uc.serialDriverAPI.GetSerialLatest()
	if err != nil {
		return nil, err
	}

	return status, nil
}

func average(values []float64) float64 {
	if len(values) == 0 {
		return -1
	}

	var sum float64
	for _, value := range values {
		sum += value
	}

	return sum / float64(len(values))
}
