package callbacks

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"gorm.io/gorm"
)

type Callbacks struct {
	Summary *prometheus.SummaryVec
}

func NewCallbacks(service string) *Callbacks {

	s := prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Namespace: "LinInk",
		Subsystem: service,
		Name:      "grom_spend_time",

		Objectives: map[float64]float64{
			0.50: 0.050, // 50% 的情况，误差为 5.0%
			0.75: 0.025, // 75% 的情况，误差为 2.5%
			0.90: 0.010, // 90% 的情况，误差为 1.0%
			0.95: 0.005, // 95% 的情况，误差为 0.5%
			0.99: 0.001, // 99% 的情况，误差为 0.1%
		},
	}, []string{"type", "table"})

	return &Callbacks{
		Summary: s,
	}
}

func (c *Callbacks) Before(typ string) func(*gorm.DB) {
	return func(d *gorm.DB) {
		start_time := time.Now()
		d.Set("start_time", start_time)
	}
}

func (c *Callbacks) After(typ string) func(*gorm.DB) {
	return func(d *gorm.DB) {
		value, _ := d.Get("start_time")
		start_time, ok := value.(time.Time)
		if !ok {
			return
		}
		table := d.Statement.Table
		if table == "" {
			table = "unknow"
		}
		c.Summary.WithLabelValues(typ, table).Observe(float64(time.Since(start_time).Milliseconds()))
	}
}
