package main

import (
	"flag"
	"runtime"
	"sync"

	logger "github.com/aguarty/litelogger"
	MQTT "github.com/eclipse/paho.mqtt.golang"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

type appEmulator struct {
	sync.Mutex
	logger        *logger.Logger
	cfg           *Config
	guiLog        chan string
	stop          chan string
	aggDataMetro  chan metroSensor
	devs          *storage
	brockerClient MQTT.Client
	readerClient  MQTT.Client
}

// flags variables
var (
	gui   *bool
	clog  *bool
	loops *int
)

func main() {
	gui = flag.Bool("gui", false, "run emulator with GUI")
	clog = flag.Bool("clog", false, "run emulator with log to console")
	loops = flag.Int("loops", 0, "run emulator with fixed loops")
	flag.Parse()

	app := appEmulator{}
	app.logger = logger.Init("info", "")
	app.guiLog = make(chan string)
	app.cfg = app.loadConfig()
	if err := app.validateConfig(); err != nil {
		app.logger.Fatal(err.Error())
	}

	//** Console GUI :TODO**
	if *gui {
		go app.startGui()
	}

	app.start()
}
