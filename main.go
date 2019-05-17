package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/rivo/tview"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

var (
	GuiLog chan string
	gui    *bool
	clog   *bool
	loops  *int
	stop   chan string

	aggDataMetro chan MetroSensor
)

func main() {

	gui = flag.Bool("gui", false, "run emulator with GUI")
	clog = flag.Bool("clog", false, "run emulator with log to console")
	loops = flag.Int("loops", 0, "run emulator with fixed loops")

	flag.Parse()

	//** Console GUI :TODO**
	if *gui {
		var (
			app            *tview.Application
			vDropDown      *tview.DropDown
			vForm1, vForm2 *tview.Form
			vLog           *tview.TextView
		)
		GuiLog = make(chan string)
		go func() {
			app = tview.NewApplication()
			vLog = tview.NewTextView().
				SetDynamicColors(true).
				SetRegions(true).
				SetWordWrap(true).
				SetChangedFunc(func() {
					app.Draw()
				})
			vLog.SetBorder(true).SetTitle("Log")
			vForm1 = tview.NewForm().
				AddDropDown("Title", []string{"Mr.", "Ms.", "Mrs.", "Dr.", "Prof."}, 0, nil).
				AddInputField("First name", "", 20, nil, nil).
				AddInputField("Last name", "", 20, nil, nil).
				AddCheckbox("Age 18+", false, nil).
				AddPasswordField("Password", "", 10, '*', nil).
				AddButton("Save", func() {
					app.SetFocus(vDropDown)
				}).
				AddButton("Quit", func() {
					app.Stop()
				})
			vForm2 = tview.NewForm().
				AddDropDown("Title", []string{"1.", "2.", "1.", "2.", "1."}, 0, nil).
				AddInputField("First name", "", 20, nil, nil).
				AddInputField("Last name", "", 20, nil, nil).
				AddCheckbox("Age 18+", false, nil).
				AddPasswordField("Password", "", 10, '*', nil).
				AddButton("Save", func() {
					app.SetFocus(vDropDown)
				}).
				AddButton("Quit", func() {
					app.Stop()
				})

			vDropDown = tview.NewDropDown()

			grid := tview.NewGrid().
				SetRows(1, 0, 1).
				SetColumns(0, 1).
				AddItem(vLog, 1, 1, 1, 2, 0, 0, false).
				AddItem(vDropDown, 1, 0, 1, 1, 0, 0, false)

			vDropDown.AddOption("1", func() {
				grid.AddItem(vForm1, 1, 0, 1, 1, 0, 0, false)
				app.SetFocus(vForm1)
			}).AddOption("2", func() {
				grid.AddItem(vForm2, 1, 0, 1, 1, 0, 0, false)
				app.SetFocus(vForm2)
			}).SetLabel("dropdrop").SetTitle("drop1")

			go func() {
				for {
					time.Sleep(200 * time.Millisecond)
					//fmt.Fprintf(wLog, "PUSH %s\n", "asd11")
					fmt.Fprintf(vLog, "PUSH %s\n", <-GuiLog)
				}
			}()

			if err := app.SetRoot(grid, true).SetFocus(vDropDown).Run(); err != nil {
				panic(err)
			}
		}()

	}

	loadConfig()

	//THIS IS MATHEMATICAAAAAAAAA (chance to send data for IR devices)
	ChanceIr = (float64(cfg.Ir.Chance) * math.Log(float64(cfg.Ir.Chance)) / (math.Log1p(float64(cfg.Ir.Count))))

	CreateDevices()

	brockerClient := InitMqttBrokerClient()
	readerClient := InitMqttReaderClient()

	log.Println("Emutale started")

	//for graceful shutdown
	exit := make(chan os.Signal)
	done := make(chan bool)
	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)
	//goroutine for graceful shutdown
	go func(brockerClient MQTT.Client, done chan bool) {
		sig := <-exit
		log.Printf("caught sig: %+v", sig)
		log.Println("Disconnected from mqttserver and exit...")
		brockerClient.Disconnect(10)
		readerClient.Disconnect(10)
		done <- true
		//time.Sleep(time.Second)
		//os.Exit(0)
	}(brockerClient, done)

	//Goroutine for connect to read coltrol chanel
	go RunCommandReader(readerClient)

	cntStop := 0

	if cfg.Ir.Enabled {
		cntStop++
		go RunWrapEmulator(brockerClient, DevIrArray)
	}
	if cfg.Temperature.Enabled {
		cntStop++
		go RunWrapEmulator(brockerClient, DevTemperatureArray)
	}
	if cfg.Light.Enabled {
		cntStop++
		go RunWrapEmulator(brockerClient, DevLightArray)
	}
	if cfg.Co2.Enabled {
		cntStop++
		go RunWrapEmulator(brockerClient, DevCo2Array)
	}

	aggDataMetro = make(chan MetroSensor, cntStop)
	stop = make(chan string, cntStop)

	if cntStop == 0 {
		exit <- syscall.SIGINT
	}
	if *loops != 0 && cntStop > 0 {
		for i := 0; i < cntStop; i++ {
			log.Println(<-stop)
		}
		exit <- syscall.SIGINT
	}

	<-done
}
