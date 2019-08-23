package main

import (
	crand "crypto/rand"
	"crypto/tls"
	"encoding/binary"
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"strconv"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

func (a *appEmulator) loadConfig() *Config {
	file, err := ioutil.ReadFile("config.cfg")
	if err != nil {
		a.logger.Info("open config error: ", err)
	}
	cfg := &Config{}
	if err = json.Unmarshal(file, cfg); err != nil {
		a.logger.Fatal("parse config error: ", err)
	}
	return cfg
}

//block for unique random
type cryptoSource struct{}

func (s cryptoSource) Seed(seed int64) {}

func (s cryptoSource) Int63() int64 {
	return int64(s.Uint64() & ^uint64(1<<63))
}

func (s cryptoSource) Uint64() (v uint64) {
	err := binary.Read(crand.Reader, binary.BigEndian, &v)
	if err != nil {
		panic(err)
	}
	return v
}

//validateConfig check config
func (a *appEmulator) validateConfig() error {
	if a.cfg.Co2.Range.Low >= a.cfg.Co2.Range.High || a.cfg.Co2.Count < a.cfg.Co2.Overflowcount {
		return errors.New("CO2 range failed")
	}
	if a.cfg.Temperature.Range.Low >= a.cfg.Temperature.Range.High || a.cfg.Temperature.Count < a.cfg.Temperature.Overflowcount {
		return errors.New("Temperature range failed")
	}
	if a.cfg.Light.Range.Low >= a.cfg.Light.Range.High || a.cfg.Light.Count < a.cfg.Light.Overflowcount {
		return errors.New("Light range failed")
	}
	if a.cfg.Ir.Chance > 100 {
		return errors.New("Chance over 100%")
	}
	return nil
}

//validateCommand validate command
func (a *appEmulator) validateCommand(c *controlCommand) error {

	switch c.Type {
	case "co2":
		if c.ID > a.cfg.Co2.Count-1 || len(a.devs.devCo2Array) == 0 {
			return errors.New("CO2 device ID does not exist")
		}
		switch c.Cmd {
		case "enabled":
			return nil
		case "balance":
			tmpval, err := strconv.Atoi(c.Value)
			if err != nil {
				return errors.New("CO2 Balance parse command failed")
			}
			if tmpval > 100 || tmpval < 0 {
				return errors.New("CO2 Balance must be in range [0-100]")
			}
		case "rangeLow":
			tmpval, err := strconv.Atoi(c.Value)
			if err != nil {
				return errors.New("CO2 Balance parse command failed")
			}
			if tmpval < 0 || tmpval >= a.devs.devCo2Array[c.ID].Range.High {
				return errors.New("CO2 Balance in out of range")
			}
		case "rangeHigh":
			tmpval, err := strconv.Atoi(c.Value)
			if err != nil {
				return errors.New("CO2 Balance parse command failed")
			}
			if tmpval <= a.devs.devCo2Array[c.ID].Range.Low {
				return errors.New("CO2 Balance in out of range")
			}
		case "chance":
			tmpval, err := strconv.Atoi(c.Value)
			if err != nil {
				return errors.New("CO2 Chance parse command failed")
			}
			if tmpval > 100 || tmpval < 0 {
				return errors.New("CO2 Chance must be in range [0-100]")
			}
		default:
			return errors.New("CO2 command failed")
		}
	case "light":
		if c.ID > a.cfg.Light.Count-1 || len(a.devs.devLightArray) == 0 {
			return errors.New("Light device ID does not exist")
		}
		switch c.Cmd {
		case "enabled":
			return nil
		case "balance":
			tmpval, err := strconv.Atoi(c.Value)
			if err != nil {
				return errors.New("Light Balance parse command failed")
			}
			if tmpval > 100 || tmpval < 0 {
				return errors.New("Light Balance must be in range [0-100]")
			}
		case "rangeLow":
			tmpval, err := strconv.Atoi(c.Value)
			if err != nil {
				return errors.New("Light Balance parse command failed")
			}
			if tmpval < 0 || tmpval >= a.devs.devLightArray[c.ID].Range.High {
				return errors.New("Light Balance in out of range")
			}
		case "rangeHigh":
			tmpval, err := strconv.Atoi(c.Value)
			if err != nil {
				return errors.New("Light Balance parse command failed")
			}
			a.logger.Info(a.devs.devLightArray[c.ID].Range.Low)
			if tmpval <= a.devs.devLightArray[c.ID].Range.Low {
				return errors.New("Light Balance in out of range")
			}
		case "chance":
			tmpval, err := strconv.Atoi(c.Value)
			if err != nil {
				return errors.New("Light Chance parse command failed")
			}
			if tmpval > 100 || tmpval < 0 {
				return errors.New("Light Chance must be in range [0-100]")
			}
		default:
			return errors.New("Light command failed")
		}
	case "temperature":
		if c.ID > a.cfg.Temperature.Count-1 || len(a.devs.devTemperatureArray) == 0 {
			return errors.New("Temperature device ID does not exist")
		}
		switch c.Cmd {
		case "enabled":
			return nil
		case "balance":
			tmpval, err := strconv.Atoi(c.Value)
			if err != nil {
				return errors.New("Temperature Balance parse command failed")
			}
			if tmpval > 100 || tmpval < 0 {
				return errors.New("Temperature Balance must be in range [0-100]")
			}
		case "rangeLow":
			tmpval, err := strconv.Atoi(c.Value)
			if err != nil {
				return errors.New("Temperature Balance parse command failed")
			}
			if tmpval < 0 || tmpval >= a.devs.devTemperatureArray[c.ID].Range.High {
				return errors.New("Temperature Balance in out of range")
			}
		case "rangeHigh":
			tmpval, err := strconv.Atoi(c.Value)
			if err != nil {
				return errors.New("Temperature Balance parse command failed")
			}
			if tmpval <= a.devs.devTemperatureArray[c.ID].Range.Low {
				return errors.New("Temperature Balance in out of range")
			}
		case "chance":
			tmpval, err := strconv.Atoi(c.Value)
			if err != nil {
				return errors.New("Temperature Chance parse command failed")
			}
			if tmpval > 100 || tmpval < 0 {
				return errors.New("Temperature Chance must be in range [0-100]")
			}
		default:
			return errors.New("Temperature command failed")
		}
	case "ir":
		if c.ID > a.cfg.Ir.Count-1 || len(a.devs.devIrArray) == 0 {
			return errors.New("IR device ID does not exist")
		}
	}
	return nil
}

//InitMqttReaderClient initialize connection to mqtt-server
func (a *appEmulator) initMqttReaderClient() {
	server := "tcp://" + a.cfg.Mqttserver.Host + ":" + a.cfg.Mqttserver.Port

	connOpts := MQTT.NewClientOptions()
	connOpts.AddBroker(server)
	connOpts.SetClientID("Bw8TNW1isvXHecvuadZj")
	connOpts.SetCleanSession(true)
	connOpts.SetUsername("")
	connOpts.SetPassword("")
	connOpts.SetTLSConfig(&tls.Config{InsecureSkipVerify: true})

	a.readerClient = MQTT.NewClient(connOpts)
	if token := a.readerClient.Connect(); token.Wait() && token.Error() != nil {
		a.logger.Fatal(token.Error())
		os.Exit(0)
	}

	a.logger.Info("Reader connected to server:", server)
}

//InitMqttBrokerClient initialize connection to mqtt-server
func (a *appEmulator) initMqttBrokerClient() {
	server := "tcp://" + a.cfg.Mqttserver.Host + ":" + a.cfg.Mqttserver.Port

	connOpts := MQTT.NewClientOptions()
	connOpts.AddBroker(server)
	connOpts.SetClientID("gxCXCxKujmU26zgi8UTa")
	connOpts.SetCleanSession(true)
	connOpts.SetUsername("")
	connOpts.SetPassword("")
	connOpts.SetTLSConfig(&tls.Config{InsecureSkipVerify: true})

	a.brockerClient = MQTT.NewClient(connOpts)
	if token := a.brockerClient.Connect(); token.Wait() && token.Error() != nil {
		a.logger.Fatal(token.Error())
		os.Exit(0)
	}
	a.logger.Info("Brocker connected to server:", server)
}

//RunCommandReader read commands from topic
// {"type":"light", "id":1, "cmd":"enabled", "value":"true"}
func (a *appEmulator) runCommandReader() {
	for {
		a.readerClient.Subscribe(a.cfg.Mqttserver.CommandTopic, 0, func(client MQTT.Client, msg MQTT.Message) {
			cmd := &controlCommand{}
			if err := json.Unmarshal(msg.Payload(), cmd); err == nil {
				if err := a.validateCommand(cmd); err == nil {
					switch cmd.Type {
					case "ir":
						switch cmd.Cmd {
						case "enabled":
							tmpval, err := strconv.ParseBool(cmd.Value)
							if err == nil {
								a.devs.devIrArray[cmd.ID].Lock.Lock()
								a.devs.devIrArray[cmd.ID].Enabled = tmpval
								a.devs.devIrArray[cmd.ID].Lock.Unlock()
							}
						}
					case "co2":
						switch cmd.Cmd {
						case "enabled":
							tmpval, err := strconv.ParseBool(cmd.Value)
							if err == nil {
								a.devs.devCo2Array[cmd.ID].Lock.Lock()
								a.devs.devCo2Array[cmd.ID].Enabled = tmpval
								a.devs.devCo2Array[cmd.ID].Lock.Unlock()
							}
						case "balance":
							tmpval, err := strconv.Atoi(cmd.Value)
							if err == nil {
								a.devs.devCo2Array[cmd.ID].Lock.Lock()
								a.devs.devCo2Array[cmd.ID].Balance = tmpval
								a.devs.devCo2Array[cmd.ID].Lock.Unlock()
							}
						case "rangeLow":
							tmpval, err := strconv.Atoi(cmd.Value)
							if err == nil {
								a.devs.devCo2Array[cmd.ID].Lock.Lock()
								a.devs.devCo2Array[cmd.ID].Range.Low = tmpval
								a.devs.devCo2Array[cmd.ID].Lock.Unlock()
							}
						case "rangeHigh":
							tmpval, err := strconv.Atoi(cmd.Value)
							if err == nil {
								a.devs.devCo2Array[cmd.ID].Lock.Lock()
								a.devs.devCo2Array[cmd.ID].Range.High = tmpval
								a.devs.devCo2Array[cmd.ID].Lock.Unlock()
							}
						case "chance":
							tmpval, err := strconv.Atoi(cmd.Value)
							if err == nil {
								a.devs.devCo2Array[cmd.ID].Lock.Lock()
								a.devs.devCo2Array[cmd.ID].Chance = tmpval
								a.devs.devCo2Array[cmd.ID].Lock.Unlock()
							}
						}
					case "light":
						switch cmd.Cmd {
						case "enabled":
							tmpval, err := strconv.ParseBool(cmd.Value)
							if err == nil {
								a.devs.devLightArray[cmd.ID].Lock.Lock()
								a.devs.devLightArray[cmd.ID].Enabled = tmpval
								a.devs.devLightArray[cmd.ID].Lock.Unlock()
							}
						case "balance":
							tmpval, err := strconv.Atoi(cmd.Value)
							if err == nil {
								a.devs.devLightArray[cmd.ID].Lock.Lock()
								a.devs.devLightArray[cmd.ID].Balance = tmpval
								a.devs.devLightArray[cmd.ID].Lock.Unlock()
							}
						case "rangeLow":
							tmpval, err := strconv.Atoi(cmd.Value)
							if err == nil {
								a.devs.devLightArray[cmd.ID].Lock.Lock()
								a.devs.devLightArray[cmd.ID].Range.Low = tmpval
								a.devs.devLightArray[cmd.ID].Lock.Unlock()
								a.logger.Info(a.devs.devLightArray[cmd.ID].Range.Low)
							}
						case "rangeHigh":
							tmpval, err := strconv.Atoi(cmd.Value)
							if err == nil {
								a.devs.devLightArray[cmd.ID].Lock.Lock()
								a.devs.devLightArray[cmd.ID].Range.High = tmpval
								a.devs.devLightArray[cmd.ID].Lock.Unlock()
							}
						case "chance":
							tmpval, err := strconv.Atoi(cmd.Value)
							if err == nil {
								a.devs.devLightArray[cmd.ID].Lock.Lock()
								a.devs.devLightArray[cmd.ID].Chance = tmpval
								a.devs.devLightArray[cmd.ID].Lock.Unlock()
							}
						}
					case "temperature":
						switch cmd.Cmd {
						case "enabled":
							tmpval, err := strconv.ParseBool(cmd.Value)
							if err == nil {
								a.devs.devTemperatureArray[cmd.ID].Lock.Lock()
								a.devs.devTemperatureArray[cmd.ID].Enabled = tmpval
								a.devs.devTemperatureArray[cmd.ID].Lock.Unlock()
							}
						case "balance":
							tmpval, err := strconv.Atoi(cmd.Value)
							if err == nil {
								a.devs.devTemperatureArray[cmd.ID].Lock.Lock()
								a.devs.devTemperatureArray[cmd.ID].Balance = tmpval
								a.devs.devTemperatureArray[cmd.ID].Lock.Unlock()
							}
						case "rangeLow":
							tmpval, err := strconv.Atoi(cmd.Value)
							if err == nil {
								a.devs.devTemperatureArray[cmd.ID].Lock.Lock()
								a.devs.devTemperatureArray[cmd.ID].Range.Low = tmpval
								a.devs.devTemperatureArray[cmd.ID].Lock.Unlock()
							}
						case "rangeHigh":
							tmpval, err := strconv.Atoi(cmd.Value)
							if err == nil {
								a.devs.devTemperatureArray[cmd.ID].Lock.Lock()
								a.devs.devTemperatureArray[cmd.ID].Range.High = tmpval
								a.devs.devTemperatureArray[cmd.ID].Lock.Unlock()
							}
						case "chance":
							tmpval, err := strconv.Atoi(cmd.Value)
							if err == nil {
								a.devs.devTemperatureArray[cmd.ID].Lock.Lock()
								a.devs.devTemperatureArray[cmd.ID].Chance = tmpval
								a.devs.devTemperatureArray[cmd.ID].Lock.Unlock()
							}
						}
					}
					a.logger.Infof("GET: %s \n", msg.Payload())
				} else {
					a.logger.Info("error parse command: ", err)
				}
			} else {
				a.logger.Info("error parse command: ", err)
			}
		})

	}
}

//RunWrapEmulator wrapper around emulator
func (a *appEmulator) runWrapEmulator(DevArray interface{}) {
	for payload := range a.emulate(DevArray) {
		a.brockerClient.Publish(a.cfg.Mqttserver.Topic+"/"+payload.Type+"/"+strconv.Itoa(payload.ID), byte(0), false, payload.Payload)
		if *gui && *clog {
			a.guiLog <- payload.Payload
		} else if *clog {
			a.logger.Infof("PUSH %s\n", payload.Payload)
		}
	}
}
