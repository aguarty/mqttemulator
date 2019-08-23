package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

type metroSensor struct {
	Name  string  `json:"name"`
	Value float64 `json:"value"`
	Range struct {
		Min   int `json:"min"`
		Max   int `json:"max"`
		Delta int `json:"delta"`
		Time  int `json:"time"`
	} `json:"range"`
	Alarms struct {
		Min   bool `json:"min"`
		Max   bool `json:"max"`
		Delta bool `json:"delta"`
	} `json:"alarms"`
}

type payloadMetro struct {
	ID   int           `json:"id"`
	Time time.Time     `json:"time"`
	Type []metroSensor `json:"type"`
}

//payloadChan format payload chanel
type payloadChan struct {
	Payload string
	ID      int
	Type    string
}

//withTimeCorrection - correct value with time rate
func withTimeCorrection() float64 {
	t := time.Now()
	var floatCorrection float64
	if t.Hour() < 12 {
		strCorrection := strconv.Itoa(t.Hour()) + "." + strconv.Itoa(t.Minute())
		floatCorrection, _ = strconv.ParseFloat(strCorrection, 64)
		floatCorrection = floatCorrection / 100

	} else {
		strCorrection := "-" + strconv.Itoa(t.Hour()-10) + "." + strconv.Itoa(t.Minute())
		floatCorrection, _ = strconv.ParseFloat(strCorrection, 64)
		floatCorrection = floatCorrection / 100
	}
	return floatCorrection
}

//Emulate - emulating data
func (a *appEmulator) emulate(devArray interface{}) <-chan payloadChan {

	c := make(chan payloadChan)

	switch devArray.(type) {

	case []*devIrModel:
		go func() {
			var cnt int
			for {
				rnd := rand.New(a.devs.crysrc)
				var chance bool
				if a.cfg.Ir.All == "" {
					chance = rnd.Intn(100) > 100-a.cfg.Ir.Chance
				} else {
					chance = true
				}
				if chance {

					for k, v := range a.devs.devIrArray {
						v.Lock.Lock()
						payload := payloadChan{}
						payload.ID = k

						if v.Enabled {
							v.Data.Time = time.Now()

							if a.cfg.Ir.All != "" {
								flo, err := strconv.ParseFloat(a.cfg.Ir.All, 64)
								if err == nil {
									v.Data.Value = flo
								}
								payload.Payload = fmt.Sprintf(dataFormat, v.Data.ID, time.Now().Format(timeformat), "ir", flo)
								//payload.Payload = fmt.Sprintf(dataFormatMetro, v.Data.ID, time.Now().Format(timeformatMetro), "ir", flo, 0, 1, a.cfg.Ir.Interval)
								payload.Type = "ir"
								c <- payload
							} else {
								if rnd.Float64()*100 < a.devs.chanceIr {
									oldval := int(v.Data.Value)
									v.Data.Value = float64(oldval ^ 1)
									payload.Payload = fmt.Sprintf(dataFormat, v.Data.ID, v.Data.Time.Format(timeformat), "ir", v.Data.Value)
									//payload.Payload = fmt.Sprintf(dataFormatMetro, v.Data.ID, time.Now().Format(timeformatMetro), "ir", v.Data.Value, 0, 1, a.cfg.Ir.Interval)
									payload.Type = "ir"
									c <- payload
								}
							}
						}
						v.Lock.Unlock()
					}
				}
				if *loops != 0 {
					cnt++
					if cnt >= *loops {
						a.Lock()
						a.stop <- "IR stop"
						a.Unlock()
						return
					}
				}
				time.Sleep(time.Millisecond * time.Duration(a.cfg.Ir.Interval))
			}
		}()

	case []*devTemperatureModel:

		go func() {
			var cnt int
			for {
				for k, v := range a.devs.devTemperatureArray {
					v.Lock.Lock()
					if v.Enabled {
						payload := payloadChan{}
						payload.ID = k
						rnd := rand.New(a.devs.crysrc)
						oldval := v.Data.Value
						var newvalue float64
						if rnd.Intn(100) > 100-a.cfg.Temperature.Chance {
							if rnd.Intn(100) > v.Balance {
								newvalue = (oldval + (rnd.Float64()+withTimeCorrection())*v.Correction)
								if newvalue >= float64(v.Range.High) {
									newvalue = oldval
									v.Balance = v.Balance + 5
								}
							} else {
								newvalue = (oldval - (rnd.Float64()+withTimeCorrection())*v.Correction)
								if newvalue <= float64(v.Range.Low) {
									newvalue = oldval
									v.Balance = v.Balance - 5
								}
							}
						} else {
							newvalue = oldval
						}
						v.Data.Value = newvalue
						v.Data.Time = time.Now()

						payload.Type = "temperature"
						payload.Payload = fmt.Sprintf(dataFormat, v.Data.ID, time.Now().Format(timeformat), "temperature", newvalue)
						//payload.Payload = fmt.Sprintf(dataFormatMetro, v.Data.ID, time.Now().Format(timeformatMetro), "temperature", v.Data.Value, a.cfg.Temperature.Range.Low, a.cfg.Temperature.Range.High, a.cfg.Temperature.Interval)
						c <- payload
					}
					v.Lock.Unlock()
				}
				if *loops != 0 {
					cnt++
					if cnt >= *loops {
						a.Lock()
						a.stop <- "Temperature stop"
						a.Unlock()
						return
					}
				}
				time.Sleep(time.Millisecond * time.Duration(a.cfg.Temperature.Interval))
			}
		}()

	case []*devLightModel:
		go func() {
			var cnt int
			for {
				for k, v := range a.devs.devLightArray {
					v.Lock.Lock()
					if v.Enabled {
						payload := payloadChan{}
						payload.ID = k
						rnd := rand.New(a.devs.crysrc)
						oldval := v.Data.Value
						var newvalue float64
						if rnd.Intn(100) > 100-a.cfg.Light.Chance {
							if rnd.Intn(100) > v.Balance {
								newvalue = (oldval + (rnd.Float64()+withTimeCorrection())*v.Correction)
								if newvalue > float64(v.Range.High) {
									newvalue = oldval
									v.Balance = v.Balance + 5
								}
							} else {
								newvalue = (oldval - (rnd.Float64()+withTimeCorrection())*v.Correction)
								if newvalue < float64(v.Range.Low) {
									newvalue = oldval
									v.Balance = v.Balance - 5
								}
							}
						} else {
							newvalue = oldval
						}
						v.Data.Value = newvalue
						v.Data.Time = time.Now()

						payload.Payload = fmt.Sprintf(dataFormat, v.Data.ID, time.Now().Format(timeformat), "light", newvalue)
						//payload.Payload = fmt.Sprintf(dataFormatMetro, v.Data.ID, time.Now().Format(timeformatMetro), "light", v.Data.Value, a.cfg.Light.Range.Low, a.cfg.Light.Range.High, a.cfg.Light.Interval)
						payload.Type = "light"
						c <- payload
					}
					v.Lock.Unlock()
				}
				if *loops != 0 {
					cnt++
					if cnt >= *loops {
						a.Lock()
						a.stop <- "Light stop"
						a.Unlock()
						return
					}
				}
				time.Sleep(time.Millisecond * time.Duration(a.cfg.Light.Interval))
			}
		}()

	case []*devCo2Model:
		go func() {
			var cnt int
			for {
				for k, v := range a.devs.devCo2Array {
					v.Lock.Lock()
					if v.Enabled {
						payload := payloadChan{}
						payload.ID = k
						rnd := rand.New(a.devs.crysrc)
						oldval := v.Data.Value
						var newvalue float64
						if rnd.Intn(100) > 100-a.cfg.Co2.Chance {
							if rnd.Intn(100) > v.Balance {
								newvalue = (oldval + (rnd.Float64()+withTimeCorrection())*v.Correction)
								if newvalue > float64(v.Range.High) {
									newvalue = oldval
									v.Balance = v.Balance + 5
								}
							} else {
								newvalue = (oldval - (rnd.Float64()+withTimeCorrection())*v.Correction)
								if newvalue < float64(v.Range.Low) {
									newvalue = oldval
									v.Balance = v.Balance - 5
								}
							}
						} else {
							newvalue = oldval
						}
						v.Data.Time = time.Now()
						v.Data.Value = newvalue
						payload.Payload = fmt.Sprintf(dataFormat, v.Data.ID, time.Now().Format(timeformat), "co2", newvalue)
						//payload.Payload = fmt.Sprintf(dataFormatMetro, v.Data.ID, time.Now().Format(timeformatMetro), "co2", v.Data.Value, a.cfg.Co2.Range.Low, a.cfg.Co2.Range.High, a.cfg.Co2.Interval)
						payload.Type = "co2"
						c <- payload
					}
					v.Lock.Unlock()
				}
				if *loops != 0 {
					cnt++
					if cnt >= *loops {
						a.Lock()
						a.stop <- "CO2 stop"
						a.Unlock()
						return
					}
				}
				time.Sleep(time.Millisecond * time.Duration(a.cfg.Co2.Interval))
			}
		}()

	}
	return c
}
