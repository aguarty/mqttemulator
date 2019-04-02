package main

import (
	crand "crypto/rand"
	"crypto/tls"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

func loadConfig() {
	file, err := ioutil.ReadFile("config.cfg")
	if err != nil {
		log.Println("open config error: ", err)
	}
	cfg = &Config{}
	if err = json.Unmarshal(file, cfg); err != nil {
		log.Fatal("parse config error: ", err)
	}
	if err := ValidateConfig(); err != nil {
		log.Fatal(err.Error())
	}
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
		log.Fatal(err)
	}
	return v
}

//ValidateConfig check config
func ValidateConfig() error {
	if cfg.Co2.Range.Low >= cfg.Co2.Range.High || cfg.Co2.Count < cfg.Co2.Overflowcount {
		return errors.New("CO2 range failed")
	}
	if cfg.Temperature.Range.Low >= cfg.Temperature.Range.High || cfg.Temperature.Count < cfg.Temperature.Overflowcount {
		return errors.New("Temperature range failed")
	}
	if cfg.Light.Range.Low >= cfg.Light.Range.High || cfg.Light.Count < cfg.Light.Overflowcount {
		return errors.New("Light range failed")
	}
	if cfg.Ir.Chance > 100 {
		return errors.New("Chance over 100%")
	}
	return nil
}

//ValidateCommand validate command
func ValidateCommand(c *ControlCommand) error {

	switch c.Type {
	case "co2":
		if c.ID > cfg.Co2.Count-1 || len(DevCo2Array) == 0 {
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
			if tmpval < 0 || tmpval >= DevCo2Array[c.ID].Range.High {
				return errors.New("CO2 Balance in out of range")
			}
		case "rangeHigh":
			tmpval, err := strconv.Atoi(c.Value)
			if err != nil {
				return errors.New("CO2 Balance parse command failed")
			}
			if tmpval <= DevCo2Array[c.ID].Range.Low {
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
		if c.ID > cfg.Light.Count-1 || len(DevLightArray) == 0 {
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
			if tmpval < 0 || tmpval >= DevLightArray[c.ID].Range.High {
				return errors.New("Light Balance in out of range")
			}
		case "rangeHigh":
			tmpval, err := strconv.Atoi(c.Value)
			if err != nil {
				return errors.New("Light Balance parse command failed")
			}
			fmt.Println(DevLightArray[c.ID].Range.Low)
			if tmpval <= DevLightArray[c.ID].Range.Low {
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
		if c.ID > cfg.Temperature.Count-1 || len(DevTemperatureArray) == 0 {
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
			if tmpval < 0 || tmpval >= DevTemperatureArray[c.ID].Range.High {
				return errors.New("Temperature Balance in out of range")
			}
		case "rangeHigh":
			tmpval, err := strconv.Atoi(c.Value)
			if err != nil {
				return errors.New("Temperature Balance parse command failed")
			}
			if tmpval <= DevTemperatureArray[c.ID].Range.Low {
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
		if c.ID > cfg.Ir.Count-1 || len(DevIrArray) == 0 {
			return errors.New("IR device ID does not exist")
		}
	}
	return nil
}

//InitMqttReaderClient initialize connection to mqtt-server
func InitMqttReaderClient() MQTT.Client {
	server := "tcp://" + cfg.Mqttserver.Host + ":" + cfg.Mqttserver.Port

	connOpts := MQTT.NewClientOptions()
	connOpts.AddBroker(server)
	connOpts.SetClientID("Bw8TNW1isvXHecvuadZj")
	connOpts.SetCleanSession(true)
	connOpts.SetUsername("")
	connOpts.SetPassword("")
	connOpts.SetTLSConfig(&tls.Config{InsecureSkipVerify: true})

	client := MQTT.NewClient(connOpts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
		os.Exit(0)
	}

	log.Println("Reader connected to server:", server)
	return client
}

//InitMqttBrokerClient initialize connection to mqtt-server
func InitMqttBrokerClient() MQTT.Client {
	server := "tcp://" + cfg.Mqttserver.Host + ":" + cfg.Mqttserver.Port

	connOpts := MQTT.NewClientOptions()
	connOpts.AddBroker(server)
	connOpts.SetClientID("gxCXCxKujmU26zgi8UTa")
	connOpts.SetCleanSession(true)
	connOpts.SetUsername("")
	connOpts.SetPassword("")
	connOpts.SetTLSConfig(&tls.Config{InsecureSkipVerify: true})

	client := MQTT.NewClient(connOpts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
		os.Exit(0)
	}

	log.Println("Brocker connected to server:", server)
	return client
}

//RunCommandReader read commands from topic
// {"type":"light", "id":1, "cmd":"enabled", "value":"true"}
func RunCommandReader(readerClient MQTT.Client) {
	for {
		readerClient.Subscribe(cfg.Mqttserver.CommandTopic, 0, func(readerClient MQTT.Client, msg MQTT.Message) {
			cmd := &ControlCommand{}
			if err := json.Unmarshal(msg.Payload(), cmd); err == nil {
				if err := ValidateCommand(cmd); err == nil {
					switch cmd.Type {
					case "ir":
						switch cmd.Cmd {
						case "enabled":
							tmpval, err := strconv.ParseBool(cmd.Value)
							if err == nil {
								DevIrArray[cmd.ID].Lock.Lock()
								DevIrArray[cmd.ID].Enabled = tmpval
								DevIrArray[cmd.ID].Lock.Unlock()
							}
						}
					case "co2":
						switch cmd.Cmd {
						case "enabled":
							tmpval, err := strconv.ParseBool(cmd.Value)
							if err == nil {
								DevCo2Array[cmd.ID].Lock.Lock()
								DevCo2Array[cmd.ID].Enabled = tmpval
								DevCo2Array[cmd.ID].Lock.Unlock()
							}
						case "balance":
							tmpval, err := strconv.Atoi(cmd.Value)
							if err == nil {
								DevCo2Array[cmd.ID].Lock.Lock()
								DevCo2Array[cmd.ID].Balance = tmpval
								DevCo2Array[cmd.ID].Lock.Unlock()
							}
						case "rangeLow":
							tmpval, err := strconv.Atoi(cmd.Value)
							if err == nil {
								DevCo2Array[cmd.ID].Lock.Lock()
								DevCo2Array[cmd.ID].Range.Low = tmpval
								DevCo2Array[cmd.ID].Lock.Unlock()
							}
						case "rangeHigh":
							tmpval, err := strconv.Atoi(cmd.Value)
							if err == nil {
								DevCo2Array[cmd.ID].Lock.Lock()
								DevCo2Array[cmd.ID].Range.High = tmpval
								DevCo2Array[cmd.ID].Lock.Unlock()
							}
						case "chance":
							tmpval, err := strconv.Atoi(cmd.Value)
							if err == nil {
								DevCo2Array[cmd.ID].Lock.Lock()
								DevCo2Array[cmd.ID].Chance = tmpval
								DevCo2Array[cmd.ID].Lock.Unlock()
							}
						}
					case "light":
						switch cmd.Cmd {
						case "enabled":
							tmpval, err := strconv.ParseBool(cmd.Value)
							if err == nil {
								DevLightArray[cmd.ID].Lock.Lock()
								DevLightArray[cmd.ID].Enabled = tmpval
								DevLightArray[cmd.ID].Lock.Unlock()
							}
						case "balance":
							tmpval, err := strconv.Atoi(cmd.Value)
							if err == nil {
								DevLightArray[cmd.ID].Lock.Lock()
								DevLightArray[cmd.ID].Balance = tmpval
								DevLightArray[cmd.ID].Lock.Unlock()
							}
						case "rangeLow":
							tmpval, err := strconv.Atoi(cmd.Value)
							if err == nil {
								DevLightArray[cmd.ID].Lock.Lock()
								DevLightArray[cmd.ID].Range.Low = tmpval
								DevLightArray[cmd.ID].Lock.Unlock()
								fmt.Println(DevLightArray[cmd.ID].Range.Low)
							}
						case "rangeHigh":
							tmpval, err := strconv.Atoi(cmd.Value)
							if err == nil {
								DevLightArray[cmd.ID].Lock.Lock()
								DevLightArray[cmd.ID].Range.High = tmpval
								DevLightArray[cmd.ID].Lock.Unlock()
							}
						case "chance":
							tmpval, err := strconv.Atoi(cmd.Value)
							if err == nil {
								DevLightArray[cmd.ID].Lock.Lock()
								DevLightArray[cmd.ID].Chance = tmpval
								DevLightArray[cmd.ID].Lock.Unlock()
							}
						}
					case "temperature":
						switch cmd.Cmd {
						case "enabled":
							tmpval, err := strconv.ParseBool(cmd.Value)
							if err == nil {
								DevTemperatureArray[cmd.ID].Lock.Lock()
								DevTemperatureArray[cmd.ID].Enabled = tmpval
								DevTemperatureArray[cmd.ID].Lock.Unlock()
							}
						case "balance":
							tmpval, err := strconv.Atoi(cmd.Value)
							if err == nil {
								DevTemperatureArray[cmd.ID].Lock.Lock()
								DevTemperatureArray[cmd.ID].Balance = tmpval
								DevTemperatureArray[cmd.ID].Lock.Unlock()
							}
						case "rangeLow":
							tmpval, err := strconv.Atoi(cmd.Value)
							if err == nil {
								DevTemperatureArray[cmd.ID].Lock.Lock()
								DevTemperatureArray[cmd.ID].Range.Low = tmpval
								DevTemperatureArray[cmd.ID].Lock.Unlock()
							}
						case "rangeHigh":
							tmpval, err := strconv.Atoi(cmd.Value)
							if err == nil {
								DevTemperatureArray[cmd.ID].Lock.Lock()
								DevTemperatureArray[cmd.ID].Range.High = tmpval
								DevTemperatureArray[cmd.ID].Lock.Unlock()
							}
						case "chance":
							tmpval, err := strconv.Atoi(cmd.Value)
							if err == nil {
								DevTemperatureArray[cmd.ID].Lock.Lock()
								DevTemperatureArray[cmd.ID].Chance = tmpval
								DevTemperatureArray[cmd.ID].Lock.Unlock()
							}
						}
					}
					log.Printf("GET: %s \n", msg.Payload())
				} else {
					log.Println("error parse command: ", err)
				}
			} else {
				log.Println("error parse command: ", err)
			}
		})

	}
}

//RunWrapEmulator wrapper around emulator
func RunWrapEmulator(brockerClient MQTT.Client, DevArray interface{}) {
	for payload := range Emulate(DevArray) {
		brockerClient.Publish(cfg.Mqttserver.Topic+"/"+payload.Type+"/"+strconv.Itoa(payload.ID), byte(0), false, payload.Payload)
		if *gui && *clog {
			GuiLog <- payload.Payload
		} else if *clog {
			log.Printf("PUSH %s\n", payload.Payload)
		}
	}
}
