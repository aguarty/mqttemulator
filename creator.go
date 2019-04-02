package main

import (
	"log"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

//CalcCorrection calc correction
func CalcCorrection(dev *devConf) float64 {

	dif := dev.Range.High - dev.Range.Low
	//a := strconv.FormatFloat(float64(dif), 'f', 1, 64)
	//i := strings.Index(a, ".")
	corr := float64(dif) / 100 //math.Pow(10, float64(i))
	return corr
}

//CreateDevices - first stap create devices
func CreateDevices() {
	var wg sync.WaitGroup
	wg.Add(4)

	if cfg.Ir.Enabled {
		go func(wg *sync.WaitGroup) {
			DevIrArray = make([]*DevIrModel, cfg.Ir.Count)
			for i := 0; i < cfg.Ir.Count; i++ {
				irSencor := &devData{}
				irSencor.ID = i
				irSencor.Time = time.Now()

				//generate current value
				rnd := rand.New(crysrc)
				if cfg.Ir.All == "" {
					if rnd.Intn(100) > 80 {
						irSencor.Value = 0.0
					} else {
						irSencor.Value = 1.0
					}
				} else {
					flo, err := strconv.ParseFloat(cfg.Ir.All, 64)
					if err == nil {
						irSencor.Value = flo
					}
				}
				DevIrArray[i] = &DevIrModel{
					Data:    irSencor,
					Enabled: true,
				}
			}
			log.Println("IR devices was generated")
			wg.Done()
		}(&wg)
	} else {
		wg.Done()
	}
	if cfg.Temperature.Enabled {
		correction := CalcCorrection(cfg.Temperature)
		go func(wg *sync.WaitGroup) {

			DevTemperatureArray = make([]*DevTemperatureModel, cfg.Temperature.Count)
			//rand.Perm: from [0, 1, 2, 3, 4, 5, ...] -> to [3, 0, 1, 5, 2, 4, ...]
			p := rand.Perm(cfg.Temperature.Count)
			//ganerate values in normal range
			for _, r := range p[cfg.Temperature.Overflowcount:cfg.Temperature.Count] {
				TemperSencor := &devData{}
				TemperSencor.ID = r
				TemperSencor.Time = time.Now()

				rnd := rand.New(crysrc)
				TemperSencor.Value = float64(rnd.Intn(cfg.Temperature.Range.High-cfg.Temperature.Range.Low) + cfg.Temperature.Range.Low)

				DevTemperatureArray[r] = &DevTemperatureModel{
					Data:       TemperSencor,
					Correction: correction,
					Enabled:    true,
					GoodDevice: true,
					Balance:    50,
					Chance:     cfg.Temperature.Chance,
				}
				DevTemperatureArray[r].Range.Low = cfg.Temperature.Range.Low
				DevTemperatureArray[r].Range.High = cfg.Temperature.Range.High
			}
			//ganerate values out of normal range
			if cfg.Temperature.Overflowcount > 0 {

				checkzero := cfg.Temperature.Range.Low - (cfg.Temperature.Range.High - cfg.Temperature.Range.Low)

				for _, r := range p[0:cfg.Temperature.Overflowcount] {
					TemperSencor := &devData{}
					TemperSencor.ID = r
					TemperSencor.Time = time.Now()
					rnd := rand.New(crysrc)
					var (
						LowRange  int
						HighRange int
					)
					if rnd.Intn(2) == 1 {
						val := float64(cfg.Temperature.Range.Low - rnd.Intn(cfg.Temperature.Range.High-cfg.Temperature.Range.Low)/2)
						if val < 0 {
							val = 0
						}
						TemperSencor.Value = val
						if checkzero < 0 {
							LowRange = 0
						} else {
							LowRange = checkzero
						}
						HighRange = cfg.Temperature.Range.Low
					} else {
						TemperSencor.Value = float64(rnd.Intn(cfg.Temperature.Range.High-cfg.Temperature.Range.Low)/2 + cfg.Temperature.Range.High)
						LowRange = cfg.Temperature.Range.High
						HighRange = cfg.Temperature.Range.High + (cfg.Temperature.Range.High-cfg.Temperature.Range.Low)/2
					}
					//fmt.Println(TemperSencor.Value)
					DevTemperatureArray[r] = &DevTemperatureModel{
						Data:       TemperSencor,
						Correction: correction,
						Enabled:    true,
						GoodDevice: false,
						Balance:    50,
						Chance:     cfg.Temperature.Chance,
					}
					DevTemperatureArray[r].Range.Low = LowRange
					DevTemperatureArray[r].Range.High = HighRange
				}
			}
			log.Println("Temperature devices was generated")
			wg.Done()
		}(&wg)
	} else {
		wg.Done()
	}
	if cfg.Light.Enabled {
		correction := CalcCorrection(cfg.Light)
		go func(wg *sync.WaitGroup) {

			DevLightArray = make([]*DevLightModel, cfg.Light.Count)
			//rand.Perm: from [0, 1, 2, 3, 4, 5, ...] -> to [3, 0, 1, 5, 2, 4, ...]
			p := rand.Perm(cfg.Light.Count)
			//ganerate values in normal range
			for _, r := range p[cfg.Light.Overflowcount:cfg.Light.Count] {
				LightSencor := &devData{}
				LightSencor.ID = r
				LightSencor.Time = time.Now()
				rnd := rand.New(crysrc)

				LightSencor.Value = float64(rnd.Intn(cfg.Light.Range.High-cfg.Light.Range.Low) + cfg.Light.Range.Low)

				DevLightArray[r] = &DevLightModel{
					Data:       LightSencor,
					Correction: correction,
					Enabled:    true,
					GoodDevice: true,
					Balance:    50,
					Chance:     cfg.Light.Chance,
				}
				DevLightArray[r].Range.Low = cfg.Light.Range.Low
				DevLightArray[r].Range.High = cfg.Light.Range.High
			}
			//ganerate values out of normal range
			if cfg.Light.Overflowcount > 0 {

				checkzero := cfg.Light.Range.Low - (cfg.Light.Range.High - cfg.Light.Range.Low)

				for _, r := range p[0:cfg.Light.Overflowcount] {
					LightSencor := &devData{}
					LightSencor.ID = r
					LightSencor.Time = time.Now()

					rnd := rand.New(crysrc)
					var (
						LowRange  int
						HighRange int
					)
					if rnd.Intn(2) == 1 {
						val := float64(cfg.Light.Range.Low - rnd.Intn(cfg.Light.Range.High-cfg.Light.Range.Low)/2)
						if val < 0 {
							val = 0
						}
						LightSencor.Value = val
						if checkzero < 0 {
							LowRange = 0
						} else {
							LowRange = checkzero
						}
						HighRange = cfg.Light.Range.Low
					} else {
						LightSencor.Value = float64(rnd.Intn(cfg.Light.Range.High-cfg.Light.Range.Low)/2 + cfg.Light.Range.High)
						LowRange = cfg.Light.Range.High
						HighRange = cfg.Light.Range.High + (cfg.Light.Range.High-cfg.Light.Range.Low)/2
					}
					DevLightArray[r] = &DevLightModel{
						Data:       LightSencor,
						Correction: correction,
						Enabled:    true,
						GoodDevice: false,
						Balance:    50,
						Chance:     cfg.Light.Chance,
					}
					DevLightArray[r].Range.Low = LowRange
					DevLightArray[r].Range.High = HighRange
				}
			}
			log.Println("Light devices was generated")
			wg.Done()
		}(&wg)
	} else {
		wg.Done()
	}
	if cfg.Co2.Enabled {
		correction := CalcCorrection(cfg.Co2)
		go func(wg *sync.WaitGroup) {
			DevCo2Array = make([]*DevCo2Model, cfg.Co2.Count)
			//rand.Perm: from [0, 1, 2, 3, 4, 5, ...] -> to [3, 0, 1, 5, 2, 4, ...]
			p := rand.Perm(cfg.Co2.Count)
			//ganerate values in normal range
			for _, r := range p[cfg.Co2.Overflowcount:cfg.Co2.Count] {
				Co2Sencor := &devData{}
				Co2Sencor.ID = r
				Co2Sencor.Time = time.Now()

				rnd := rand.New(crysrc)
				Co2Sencor.Value = float64(rnd.Intn(cfg.Co2.Range.High-cfg.Co2.Range.Low) + cfg.Co2.Range.Low)

				DevCo2Array[r] = &DevCo2Model{
					Data:       Co2Sencor,
					Correction: correction,
					Enabled:    true,
					GoodDevice: true,
					Balance:    50,
					Chance:     cfg.Co2.Chance,
				}
				DevCo2Array[r].Range.Low = cfg.Co2.Range.Low
				DevCo2Array[r].Range.High = cfg.Co2.Range.High
			}
			//ganerate values out of normal range
			if cfg.Co2.Overflowcount > 0 {

				checkzero := cfg.Co2.Range.Low - (cfg.Co2.Range.High - cfg.Co2.Range.Low)

				for _, r := range p[0:cfg.Co2.Overflowcount] {
					Co2Sencor := &devData{}
					Co2Sencor.ID = r
					Co2Sencor.Time = time.Now()

					rnd := rand.New(crysrc)
					var (
						LowRange  int
						HighRange int
					)
					if rnd.Intn(2) == 1 {
						val := float64(cfg.Co2.Range.Low - rnd.Intn(cfg.Co2.Range.High-cfg.Co2.Range.Low)/2)
						if val < 0 {
							val = 0
						}
						Co2Sencor.Value = val
						if checkzero < 0 {
							LowRange = 0
						} else {
							LowRange = checkzero
						}
						HighRange = cfg.Co2.Range.Low
					} else {
						Co2Sencor.Value = float64(rnd.Intn(cfg.Co2.Range.High-cfg.Co2.Range.Low)/2 + cfg.Co2.Range.High)
						LowRange = cfg.Co2.Range.High
						HighRange = cfg.Co2.Range.High + (cfg.Co2.Range.High-cfg.Co2.Range.Low)/2
					}
					DevCo2Array[r] = &DevCo2Model{
						Data:       Co2Sencor,
						Correction: correction,
						Enabled:    true,
						GoodDevice: false,
						Balance:    50,
						Chance:     cfg.Co2.Chance,
					}
					DevCo2Array[r].Range.Low = LowRange
					DevCo2Array[r].Range.High = HighRange
				}
			}
			log.Println("Co2 devices was generated")
			wg.Done()
		}(&wg)
	} else {
		wg.Done()
	}
	wg.Wait()
}
