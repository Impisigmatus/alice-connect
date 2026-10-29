package application

import (
	"fmt"
	"strings"
	"time"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/mem"
	"github.com/shirou/gopsutil/v4/sensors"
)

type Application struct{}

func NewApplication() *Application {
	return &Application{}
}

func (app *Application) Hello() (string, error) {
	percent, err := cpu.Percent(time.Second, false)
	if err != nil {
		return "", err
	}

	// Оперативная память
	vm, err := mem.VirtualMemory()
	if err != nil {
		return "", err
	}

	data := "" +
		fmt.Sprintf("ЦПУ: %.1f%%; ", percent[0]) +
		fmt.Sprintf("ОЗУ: %.1f%% (%d из %d МБ).",
			vm.UsedPercent,
			vm.Used/1024/1024,
			vm.Total/1024/1024,
		)
	// Температура (ошибка здесь часто «частичная», поэтому смотрим на сами данные)
	temps, _ := sensors.SensorsTemperatures()
	for _, t := range temps {
		key := strings.ToLower(t.SensorKey)
		if strings.Contains(key, "core") || strings.Contains(key, "cpu") ||
			strings.Contains(key, "tctl") || strings.Contains(key, "package") {
			data += fmt.Sprintf(" Температура %s: %.0f°C;", t.SensorKey, t.Temperature)
		}
	}

	return "Привет! Go сервис на Orange Pi работает. Используется " + data, nil
}
