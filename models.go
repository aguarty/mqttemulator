package main

import (
	"sync"
	"time"
)

const (
	timeformat = "2006-01-02 15:04:05"
)

var (
	crysrc cryptoSource
	cfg    *Config

	DevIrArray          []*DevIrModel
	DevTemperatureArray []*DevTemperatureModel
	DevLightArray       []*DevLightModel
	DevCo2Array         []*DevCo2Model

	ChanceIr float64
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
type DevIrModel struct {
	Data       *devData `json:"data"`
	Lock       sync.Mutex
	Correction float64 `json:"correction"`
	Enabled    bool    `json:"enabled"`
}

//DevTemperatureModel device model
type DevTemperatureModel struct {
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
type DevLightModel struct {
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
type DevCo2Model struct {
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

type ControlCommand struct {
	Type  string `json:"type"`
	ID    int    `json:"id"`
	Cmd   string `json:"cmd"`
	Value string `json:"value"`
}
