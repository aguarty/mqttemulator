package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

//PayloadChan format payload chanel
type PayloadChan struct {
	Payload string
	ID      int
	Type    string
}

//WithTimeCorrection - correct value with time rate
func WithTimeCorrection() float64 {
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
func Emulate(DevArray interface{}) <-chan PayloadChan {

	c := make(chan PayloadChan)

	switch DevArray.(type) {

	case []*DevIrModel:
		go func() {
			for {
				rnd := rand.New(crysrc)
				var chance bool
				if cfg.Ir.All == "" {
					chance = rnd.Intn(100) > 100-cfg.Ir.Chance
				} else {
					chance = true
				}
				if chance {

					for k, v := range DevIrArray {
						v.Lock.Lock()
						payload := PayloadChan{}
						payload.ID = k

						if v.Enabled {
							v.Data.Time = time.Now()

							if cfg.Ir.All != "" {
								flo, err := strconv.ParseFloat(cfg.Ir.All, 64)
								if err == nil {
									v.Data.Value = flo
								}
								payload.Payload = fmt.Sprintf(`{"ID":"%d", "time":"%s", "type":"%s", "value":"%.0f"}`, v.Data.ID, time.Now().Format(timeformat), "ir", flo)
								payload.Type = "ir"
								c <- payload
							} else {
								if rnd.Float64()*100 < ChanceIr {
									oldval := int(v.Data.Value)
									v.Data.Value = float64(oldval ^ 1)
									payload.Payload = fmt.Sprintf(`{"ID":"%d", "time":"%s", "type":"%s", "value":"%.0f"}`, v.Data.ID, v.Data.Time.Format(timeformat), "ir", v.Data.Value)
									payload.Type = "ir"
									c <- payload
								}
							}
						}
						v.Lock.Unlock()
					}
				}
				time.Sleep(time.Millisecond * time.Duration(cfg.Ir.Interval))
			}
		}()

	case []*DevTemperatureModel:

		go func() {
			for {
				for k, v := range DevTemperatureArray {
					v.Lock.Lock()
					if v.Enabled {
						payload := PayloadChan{}
						payload.ID = k
						rnd := rand.New(crysrc)
						oldval := v.Data.Value
						var newvalue float64
						if rnd.Intn(100) > 100-cfg.Temperature.Chance {
							if rnd.Intn(100) > v.Balance {
								newvalue = (oldval + (rnd.Float64()+WithTimeCorrection())*v.Correction)
								if newvalue >= float64(v.Range.High) {
									newvalue = oldval
									v.Balance = v.Balance + 5
								}
							} else {
								newvalue = (oldval - (rnd.Float64()+WithTimeCorrection())*v.Correction)
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
						payload.Payload = fmt.Sprintf(`{"ID":"%d", "time":"%s", "type":"%s", "value":"%.2f"}`, v.Data.ID, time.Now().Format(timeformat), "temperature", newvalue)
						c <- payload
					}
					v.Lock.Unlock()
				}
				time.Sleep(time.Millisecond * time.Duration(cfg.Temperature.Interval))
			}
		}()

	case []*DevLightModel:
		go func() {
			for {
				for k, v := range DevLightArray {
					v.Lock.Lock()
					if v.Enabled {
						payload := PayloadChan{}
						payload.ID = k
						rnd := rand.New(crysrc)
						oldval := v.Data.Value
						var newvalue float64
						if rnd.Intn(100) > 100-cfg.Light.Chance {
							if rnd.Intn(100) > v.Balance {
								newvalue = (oldval + (rnd.Float64()+WithTimeCorrection())*v.Correction)
								if newvalue > float64(v.Range.High) {
									newvalue = oldval
									v.Balance = v.Balance + 5
								}
							} else {
								newvalue = (oldval - (rnd.Float64()+WithTimeCorrection())*v.Correction)
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

						payload.Payload = fmt.Sprintf(`{"ID":"%d", "time":"%s", "type":"%s", "value":"%.2f"}`, v.Data.ID, time.Now().Format(timeformat), "light", newvalue)
						payload.Type = "light"
						c <- payload
					}
					v.Lock.Unlock()
				}
				time.Sleep(time.Millisecond * time.Duration(cfg.Light.Interval))
			}
		}()

	case []*DevCo2Model:
		go func() {
			for {
				for k, v := range DevCo2Array {
					v.Lock.Lock()
					if v.Enabled {
						payload := PayloadChan{}
						payload.ID = k
						rnd := rand.New(crysrc)
						oldval := v.Data.Value
						var newvalue float64
						if rnd.Intn(100) > 100-cfg.Co2.Chance {
							if rnd.Intn(100) > v.Balance {
								newvalue = (oldval + (rnd.Float64()+WithTimeCorrection())*v.Correction)
								if newvalue > float64(v.Range.High) {
									newvalue = oldval
									v.Balance = v.Balance + 5
								}
							} else {
								newvalue = (oldval - (rnd.Float64()+WithTimeCorrection())*v.Correction)
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
						payload.Payload = fmt.Sprintf(`{"ID":"%d", "time":"%s", "type":"%s", "value":"%.2f"}`, v.Data.ID, time.Now().Format(timeformat), "co2", newvalue)
						payload.Type = "co2"
						c <- payload
					}
					v.Lock.Unlock()
				}
				time.Sleep(time.Millisecond * time.Duration(cfg.Co2.Interval))
			}
		}()

	}

	return c
}
