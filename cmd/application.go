package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rivo/tview"
)

func (app *appEmulator) start() {

	app.createDevices()

	app.initMqttBrokerClient()
	app.initMqttReaderClient()

	app.logger.Info("Emutale started")

	//for graceful shutdown
	exit := make(chan os.Signal)
	done := make(chan bool)
	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)

	//goroutine for graceful shutdown
	go func() {
		sig := <-exit
		app.logger.Infof("caught sig: %+v", sig)
		app.logger.Info("Disconnected from mqttserver and exit...")
		app.brockerClient.Disconnect(10)
		app.readerClient.Disconnect(10)
		done <- true
	}()

	//Goroutine for connect to read coltrol chanel
	go app.runCommandReader()

	cntStop := 0

	if app.cfg.Ir.Enabled {
		cntStop++
		go app.runWrapEmulator(app.devs.devIrArray)
	}
	if app.cfg.Temperature.Enabled {
		cntStop++
		go app.runWrapEmulator(app.devs.devTemperatureArray)
	}
	if app.cfg.Light.Enabled {
		cntStop++
		go app.runWrapEmulator(app.devs.devLightArray)
	}
	if app.cfg.Co2.Enabled {
		cntStop++
		go app.runWrapEmulator(app.devs.devCo2Array)
	}

	//app.aggDataMetro = make(chan metroSensor, cntStop)
	app.Lock()
	app.stop = make(chan string, cntStop)
	app.Unlock()

	if cntStop == 0 {
		exit <- syscall.SIGINT
	}
	if *loops != 0 && cntStop > 0 {
		for i := 0; i < cntStop; i++ {
			app.logger.Info(<-app.stop)
		}
		exit <- syscall.SIGINT
	}

	<-done
}

func (app *appEmulator) startGui() {

	var (
		guiapp         *tview.Application
		vDropDown      *tview.DropDown
		vForm1, vForm2 *tview.Form
		vLog           *tview.TextView
	)

	guiapp = tview.NewApplication()
	vLog = tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWordWrap(true).
		SetChangedFunc(func() {
			guiapp.Draw()
		})
	vLog.SetBorder(true).SetTitle("Log")
	vForm1 = tview.NewForm().
		AddDropDown("Title", []string{"Mr.", "Ms.", "Mrs.", "Dr.", "Prof."}, 0, nil).
		AddInputField("First name", "", 20, nil, nil).
		AddInputField("Last name", "", 20, nil, nil).
		AddCheckbox("Age 18+", false, nil).
		AddPasswordField("Password", "", 10, '*', nil).
		AddButton("Save", func() {
			guiapp.SetFocus(vDropDown)
		}).
		AddButton("Quit", func() {
			guiapp.Stop()
		})
	vForm2 = tview.NewForm().
		AddDropDown("Title", []string{"1.", "2.", "1.", "2.", "1."}, 0, nil).
		AddInputField("First name", "", 20, nil, nil).
		AddInputField("Last name", "", 20, nil, nil).
		AddCheckbox("Age 18+", false, nil).
		AddPasswordField("Password", "", 10, '*', nil).
		AddButton("Save", func() {
			guiapp.SetFocus(vDropDown)
		}).
		AddButton("Quit", func() {
			guiapp.Stop()
		})

	vDropDown = tview.NewDropDown()

	grid := tview.NewGrid().
		SetRows(1, 0, 1).
		SetColumns(0, 1).
		AddItem(vLog, 1, 1, 1, 2, 0, 0, false).
		AddItem(vDropDown, 1, 0, 1, 1, 0, 0, false)

	vDropDown.AddOption("1", func() {
		grid.AddItem(vForm1, 1, 0, 1, 1, 0, 0, false)
		guiapp.SetFocus(vForm1)
	}).AddOption("2", func() {
		grid.AddItem(vForm2, 1, 0, 1, 1, 0, 0, false)
		guiapp.SetFocus(vForm2)
	}).SetLabel("dropdrop").SetTitle("drop1")

	go func() {
		for {
			time.Sleep(200 * time.Millisecond)
			//fmt.Fprintf(wLog, "PUSH %s\n", "asd11")
			fmt.Fprintf(vLog, "PUSH %s\n", <-app.guiLog)
		}
	}()

	if err := guiapp.SetRoot(grid, true).SetFocus(vDropDown).Run(); err != nil {
		panic(err)
	}

}
