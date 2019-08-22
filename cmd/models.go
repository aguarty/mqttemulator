package main

import (
	"sync"
	"time"
)

const (
	timeformat      = "2006-01-02 15:04:05"
	timeformatMetro = time.RFC3339

	dataFormat      = `{"ID":"%d", "time":"%s", "type":"%s", "value":"%.2f"}`
	dataFormatMetro = `{"id":%d,"time":"%s","type":[{"name":"%s","value":%.2f,"range":{"min":%d,"max":%d,"delta":1,"time":%d},"alarms":{"min":false,"max":false,"delta":false}}]}`
)

var (
	crysrc cryptoSource
	cfg    *Config

	devIrArray          []*devIrModel
	devTemperatureArray []*devTemperatureModel
	devLightArray       []*devLightModel
	devCo2Array         []*devCo2Model

	chanceIr float64
)

type mqttServer struct {
	Host         string `json:"host"`
	Port         string `json:"port"`
	Login        string `json:"login"`
	Password     string `json:"password"`
	Topic        string `json:"topic"`
	CommandTopic string `json:"commandtopic"`
}

type irConf struct {
	Enabled  bool   `json:"enabled"`
	Count    int    `json:"count"`
	Interval int    `json:"interval"`
	All      string `json:"all"`
	Chance   int    `json:"chance"`
}

type devConf struct {
	Enabled  bool `json:"enabled"`
	Count    int  `json:"count"`
	Interval int  `json:"interval"`
	Range    struct {
		Low  int `json:"low"`
		High int `json:"high"`
	} `json:"range"`
	Overflowcount int
	Chance        int `json:"chance"`
}

//Config configuration emulator
type Config struct {
	Mqttserver  *mqttServer `json:"mqttserver"`
	Ir          *irConf     `json:"ir"`
	Temperature *devConf    `json:"temperature"`
	Light       *devConf    `json:"light"`
	Co2         *devConf    `json:"co2"`
}

type devData struct {
	Time  time.Time `json:"time"`
	ID    int       `json:"id"`
	Value float64   `json:"value"`
}

//DevIrModel device model
type devIrModel struct {
	Data       *devData `json:"data"`
	Lock       sync.Mutex
	Correction float64 `json:"correction"`
	Enabled    bool    `json:"enabled"`
}

//DevTemperatureModel device model
type devTemperatureModel struct {
	Data       *devData `json:"data"`
	Lock       sync.Mutex
	Correction float64 `json:"correction"`
	Enabled    bool    `json:"enabled"`
	GoodDevice bool    `json:"gooddevice"`
	Balance    int     `json:"balance"`
	Chance     int     `json:"chance"`
	Range      struct {
		Low  int `json:"low"`
		High int `json:"high"`
	} `json:"range"`
}

//DevLightModel device model
type devLightModel struct {
	Data       *devData `json:"data"`
	Lock       sync.Mutex
	Correction float64 `json:"correction"`
	Enabled    bool    `json:"enabled"`
	GoodDevice bool    `json:"gooddevice"`
	Balance    int     `json:"balance"`
	Chance     int     `json:"chance"`
	Range      struct {
		Low  int `json:"low"`
		High int `json:"high"`
	} `json:"range"`
}

//DevCo2Model device model
type devCo2Model struct {
	Data       *devData `json:"data"`
	Lock       sync.Mutex
	Correction float64 `json:"correction"`
	Enabled    bool    `json:"enabled"`
	GoodDevice bool    `json:"gooddevice"`
	Balance    int     `json:"balance"`
	Chance     int     `json:"chance"`
	Range      struct {
		Low  int `json:"low"`
		High int `json:"high"`
	} `json:"range"`
}

type controlCommand struct {
	Type  string `json:"type"`
	ID    int    `json:"id"`
	Cmd   string `json:"cmd"`
	Value string `json:"value"`
}
