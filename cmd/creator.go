package main

import (
	"math"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

//calcCorrection calc correction
func calcCorrection(dev *devConf) float64 {

	dif := dev.Range.High - dev.Range.Low
	//a := strconv.FormatFloat(float64(dif), 'f', 1, 64)
	//i := strings.Index(a, ".")
	corr := float64(dif) / 100 //math.Pow(10, float64(i))
	return corr
}

//CreateDevices - first stap create devices
func (a *appEmulator) createDevices() {

	a.devs = &storage{}
	//THIS IS MATHEMATICAAAAAAAAA (chance to send data for IR devices)
	a.devs.chanceIr = (float64(a.cfg.Ir.Chance) * math.Log(float64(a.cfg.Ir.Chance)) / (math.Log1p(float64(a.cfg.Ir.Count))))

	var wg sync.WaitGroup
	wg.Add(4)

	if a.cfg.Ir.Enabled {
		go func(wg *sync.WaitGroup) {
			a.devs.devIrArray = make([]*devIrModel, a.cfg.Ir.Count)
			for i := 0; i < a.cfg.Ir.Count; i++ {
				irSencor := &devData{}
				irSencor.ID = i
				irSencor.Time = time.Now()

				//generate current value
				rnd := rand.New(a.devs.crysrc)
				if a.cfg.Ir.All == "" {
					if rnd.Intn(100) > 80 {
						irSencor.Value = 0.0
					} else {
						irSencor.Value = 1.0
					}
				} else {
					flo, err := strconv.ParseFloat(a.cfg.Ir.All, 64)
					if err == nil {
						irSencor.Value = flo
					}
				}
				a.devs.devIrArray[i] = &devIrModel{
					Data:    irSencor,
					Enabled: true,
				}
			}
			a.logger.Info("IR devices was generated")
			wg.Done()
		}(&wg)
	} else {
		wg.Done()
	}
	if a.cfg.Temperature.Enabled {
		correction := calcCorrection(a.cfg.Temperature)
		go func(wg *sync.WaitGroup) {

			a.devs.devTemperatureArray = make([]*devTemperatureModel, a.cfg.Temperature.Count)
			//rand.Perm: from [0, 1, 2, 3, 4, 5, ...] -> to [3, 0, 1, 5, 2, 4, ...]
			p := rand.Perm(a.cfg.Temperature.Count)
			//ganerate values in normal range
			for _, r := range p[a.cfg.Temperature.Overflowcount:a.cfg.Temperature.Count] {
				TemperSencor := &devData{}
				TemperSencor.ID = r
				TemperSencor.Time = time.Now()

				rnd := rand.New(a.devs.crysrc)
				TemperSencor.Value = float64(rnd.Intn(a.cfg.Temperature.Range.High-a.cfg.Temperature.Range.Low) + a.cfg.Temperature.Range.Low)

				a.devs.devTemperatureArray[r] = &devTemperatureModel{
					Data:       TemperSencor,
					Correction: correction,
					Enabled:    true,
					GoodDevice: true,
					Balance:    50,
					Chance:     a.cfg.Temperature.Chance,
				}
				a.devs.devTemperatureArray[r].Range.Low = a.cfg.Temperature.Range.Low
				a.devs.devTemperatureArray[r].Range.High = a.cfg.Temperature.Range.High
			}
			//ganerate values out of normal range
			if a.cfg.Temperature.Overflowcount > 0 {

				checkzero := a.cfg.Temperature.Range.Low - (a.cfg.Temperature.Range.High - a.cfg.Temperature.Range.Low)

				for _, r := range p[0:a.cfg.Temperature.Overflowcount] {
					TemperSencor := &devData{}
					TemperSencor.ID = r
					TemperSencor.Time = time.Now()
					rnd := rand.New(a.devs.crysrc)
					var (
						LowRange  int
						HighRange int
					)
					if rnd.Intn(2) == 1 {
						val := float64(a.cfg.Temperature.Range.Low - rnd.Intn(a.cfg.Temperature.Range.High-a.cfg.Temperature.Range.Low)/2)
						if val < 0 {
							val = 0
						}
						TemperSencor.Value = val
						if checkzero < 0 {
							LowRange = 0
						} else {
							LowRange = checkzero
						}
						HighRange = a.cfg.Temperature.Range.Low
					} else {
						TemperSencor.Value = float64(rnd.Intn(a.cfg.Temperature.Range.High-a.cfg.Temperature.Range.Low)/2 + a.cfg.Temperature.Range.High)
						LowRange = a.cfg.Temperature.Range.High
						HighRange = a.cfg.Temperature.Range.High + (a.cfg.Temperature.Range.High-a.cfg.Temperature.Range.Low)/2
					}
					//fmt.Info(TemperSencor.Value)
					a.devs.devTemperatureArray[r] = &devTemperatureModel{
						Data:       TemperSencor,
						Correction: correction,
						Enabled:    true,
						GoodDevice: false,
						Balance:    50,
						Chance:     a.cfg.Temperature.Chance,
					}
					a.devs.devTemperatureArray[r].Range.Low = LowRange
					a.devs.devTemperatureArray[r].Range.High = HighRange
				}
			}
			a.logger.Info("Temperature devices was generated")
			wg.Done()
		}(&wg)
	} else {
		wg.Done()
	}
	if a.cfg.Light.Enabled {
		correction := calcCorrection(a.cfg.Light)
		go func(wg *sync.WaitGroup) {

			a.devs.devLightArray = make([]*devLightModel, a.cfg.Light.Count)
			//rand.Perm: from [0, 1, 2, 3, 4, 5, ...] -> to [3, 0, 1, 5, 2, 4, ...]
			p := rand.Perm(a.cfg.Light.Count)
			//ganerate values in normal range
			for _, r := range p[a.cfg.Light.Overflowcount:a.cfg.Light.Count] {
				LightSencor := &devData{}
				LightSencor.ID = r
				LightSencor.Time = time.Now()
				rnd := rand.New(a.devs.crysrc)

				LightSencor.Value = float64(rnd.Intn(a.cfg.Light.Range.High-a.cfg.Light.Range.Low) + a.cfg.Light.Range.Low)

				a.devs.devLightArray[r] = &devLightModel{
					Data:       LightSencor,
					Correction: correction,
					Enabled:    true,
					GoodDevice: true,
					Balance:    50,
					Chance:     a.cfg.Light.Chance,
				}
				a.devs.devLightArray[r].Range.Low = a.cfg.Light.Range.Low
				a.devs.devLightArray[r].Range.High = a.cfg.Light.Range.High
			}
			//ganerate values out of normal range
			if a.cfg.Light.Overflowcount > 0 {

				checkzero := a.cfg.Light.Range.Low - (a.cfg.Light.Range.High - a.cfg.Light.Range.Low)

				for _, r := range p[0:a.cfg.Light.Overflowcount] {
					LightSencor := &devData{}
					LightSencor.ID = r
					LightSencor.Time = time.Now()

					rnd := rand.New(a.devs.crysrc)
					var (
						LowRange  int
						HighRange int
					)
					if rnd.Intn(2) == 1 {
						val := float64(a.cfg.Light.Range.Low - rnd.Intn(a.cfg.Light.Range.High-a.cfg.Light.Range.Low)/2)
						if val < 0 {
							val = 0
						}
						LightSencor.Value = val
						if checkzero < 0 {
							LowRange = 0
						} else {
							LowRange = checkzero
						}
						HighRange = a.cfg.Light.Range.Low
					} else {
						LightSencor.Value = float64(rnd.Intn(a.cfg.Light.Range.High-a.cfg.Light.Range.Low)/2 + a.cfg.Light.Range.High)
						LowRange = a.cfg.Light.Range.High
						HighRange = a.cfg.Light.Range.High + (a.cfg.Light.Range.High-a.cfg.Light.Range.Low)/2
					}
					a.devs.devLightArray[r] = &devLightModel{
						Data:       LightSencor,
						Correction: correction,
						Enabled:    true,
						GoodDevice: false,
						Balance:    50,
						Chance:     a.cfg.Light.Chance,
					}
					a.devs.devLightArray[r].Range.Low = LowRange
					a.devs.devLightArray[r].Range.High = HighRange
				}
			}
			a.logger.Info("Light devices was generated")
			wg.Done()
		}(&wg)
	} else {
		wg.Done()
	}
	if a.cfg.Co2.Enabled {
		correction := calcCorrection(a.cfg.Co2)
		go func(wg *sync.WaitGroup) {
			a.devs.devCo2Array = make([]*devCo2Model, a.cfg.Co2.Count)
			//rand.Perm: from [0, 1, 2, 3, 4, 5, ...] -> to [3, 0, 1, 5, 2, 4, ...]
			p := rand.Perm(a.cfg.Co2.Count)
			//ganerate values in normal range
			for _, r := range p[a.cfg.Co2.Overflowcount:a.cfg.Co2.Count] {
				Co2Sencor := &devData{}
				Co2Sencor.ID = r
				Co2Sencor.Time = time.Now()

				rnd := rand.New(a.devs.crysrc)
				Co2Sencor.Value = float64(rnd.Intn(a.cfg.Co2.Range.High-a.cfg.Co2.Range.Low) + a.cfg.Co2.Range.Low)

				a.devs.devCo2Array[r] = &devCo2Model{
					Data:       Co2Sencor,
					Correction: correction,
					Enabled:    true,
					GoodDevice: true,
					Balance:    50,
					Chance:     a.cfg.Co2.Chance,
				}
				a.devs.devCo2Array[r].Range.Low = a.cfg.Co2.Range.Low
				a.devs.devCo2Array[r].Range.High = a.cfg.Co2.Range.High
			}
			//ganerate values out of normal range
			if a.cfg.Co2.Overflowcount > 0 {

				checkzero := a.cfg.Co2.Range.Low - (a.cfg.Co2.Range.High - a.cfg.Co2.Range.Low)

				for _, r := range p[0:a.cfg.Co2.Overflowcount] {
					Co2Sencor := &devData{}
					Co2Sencor.ID = r
					Co2Sencor.Time = time.Now()

					rnd := rand.New(a.devs.crysrc)
					var (
						LowRange  int
						HighRange int
					)
					if rnd.Intn(2) == 1 {
						val := float64(a.cfg.Co2.Range.Low - rnd.Intn(a.cfg.Co2.Range.High-a.cfg.Co2.Range.Low)/2)
						if val < 0 {
							val = 0
						}
						Co2Sencor.Value = val
						if checkzero < 0 {
							LowRange = 0
						} else {
							LowRange = checkzero
						}
						HighRange = a.cfg.Co2.Range.Low
					} else {
						Co2Sencor.Value = float64(rnd.Intn(a.cfg.Co2.Range.High-a.cfg.Co2.Range.Low)/2 + a.cfg.Co2.Range.High)
						LowRange = a.cfg.Co2.Range.High
						HighRange = a.cfg.Co2.Range.High + (a.cfg.Co2.Range.High-a.cfg.Co2.Range.Low)/2
					}
					a.devs.devCo2Array[r] = &devCo2Model{
						Data:       Co2Sencor,
						Correction: correction,
						Enabled:    true,
						GoodDevice: false,
						Balance:    50,
						Chance:     a.cfg.Co2.Chance,
					}
					a.devs.devCo2Array[r].Range.Low = LowRange
					a.devs.devCo2Array[r].Range.High = HighRange
				}
			}
			a.logger.Info("Co2 devices was generated")
			wg.Done()
		}(&wg)
	} else {
		wg.Done()
	}
	wg.Wait()
}
