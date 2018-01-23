package funcs

import (
	"errors"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/gaochao1/swcollector/g"

	"github.com/open-falcon/common/model"
	"strconv"
)

func ready() error {
	log.SetFlags(log.Lshortfile + log.LstdFlags)
	// init influxdb
	if g.Config() == nil {
		g.ParseConfig("../cfg.json")
		if g.Config() == nil {
			return errors.New("cfg parse failed.")
		}
		/*
			if g.Config().SwitchHosts.Enabled {
				hostcfg := g.Config().SwitchHosts.Hosts
				g.ParseHostConfig(hostcfg)
			}
			if g.Config().CustomMetrics.Enabled {
				custMetrics := g.Config().CustomMetrics.Template
				g.ParseCustConfig(custMetrics)
			}
			g.InitRootDir()
			g.InitLocalIps()
			g.InitRpcClients()
		*/
	}
	NewLastifMap()
	return nil
}

func getval(val interface{}) (float64, bool) {
	switch val.(type) {
	case float64:
		return val.(float64), true
	case float32:
		f32 := val.(float32)
		return float64(f32), true
	case int:
		it := val.(int)
		return float64(it), true
	case int32:
		it := val.(int32)
		return float64(it), true
	case int64:
		it := val.(int64)
		return float64(it), true
	case uint:
		it := val.(uint)
		return float64(it), true
	case uint8:
		it := val.(uint8)
		return float64(it), true
	case uint16:
		it := val.(uint16)
		return float64(it), true
	case uint32:
		it := val.(uint32)
		return float64(it), true
	case uint64:
		it := val.(uint64)
		return float64(it), true
	case string:
		str := val.(string)
		it, err := strconv.Atoi(str)
		if err != nil {
			return float64(0), false
		}
		return float64(it), true
	default:
		return float64(0), false
	}
}

func metricsValueString(this *model.MetricValue) string {
	var unit string = "b"
	val, ok := getval(this.Value)
	if !ok {
		log.Fatalf("Can't convert %v to float64", this.Value)
	}
	switch {
	case val > 1000:
		val = val / 1000
		unit = "Kb"
	case val > (1000 * 1000):
		val = val / (1000 * 10000)
		unit = "Mb"
	case val > (1000 * 1000 * 1000):
		val = val / (1000 * 1000 * 1000)
		unit = "Gb"
	case val > (1000 * 1000 * 1000 * 1000):
		val = val / (1000 * 1000 * 1000 * 1000)
		unit = "Tb"
	}
	return fmt.Sprintf(
		"<Endpoint:%s, Metric:%s, Type:%s, Tags:%s, Step:%d, Time:%d, Value:%.4f %s>",
		this.Endpoint,
		this.Metric,
		this.Type,
		this.Tags,
		this.Step,
		this.Timestamp,
		val,
		unit,
	)
}

func TestSwIfMetrics(t *testing.T) {
	ready()

	c := 0
	for {
		L := swIfMetrics()
		log.Printf("Collectioon [%d]th swIfMetrics, metrics list %d", c, len(L))

		for i, l := range L {
			log.Printf("%d => %d %s", c, i, metricsValueString(l))
		}

		c++
		time.Sleep(time.Second * 10)
	}

}
