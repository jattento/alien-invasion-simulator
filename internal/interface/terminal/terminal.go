package terminal

import (
	"io"
	"strings"
	"sync/atomic"
	"time"

	"github.com/liamg/gobless"
)

type Manager struct {
	output       io.Writer
	LogsCh       <-chan string
	DayCounterCh <-chan string
	CitiesCh     <-chan []string

	logs string

	WaitTime int64
}

func New(output io.Writer, logsCh <-chan string, dayCounterCh <-chan string, citiesCh <-chan []string) *Manager {
	return &Manager{
		output:       output,
		LogsCh:       logsCh,
		DayCounterCh: dayCounterCh,
		CitiesCh:     citiesCh,
		WaitTime:     int64(time.Second),
	}
}

func (manager *Manager) Run() error {
	gui := gobless.NewGUI()
	if err := gui.Init(); err != nil {
		panic(err)
	}
	defer gui.Close()

	logsBox := gobless.NewTextBox()
	logsBox.SetTitle("LOGS")
	logsBox.SetTextWrap(true)
	logsColumn := gobless.NewColumn(
		gobless.GridSizeHalf,
		logsBox,
	)

	citiesBox1 := gobless.NewTextBox()
	citiesBox1.SetTextWrap(true)
	citiesBox1.SetTitle("CITIES")
	citiesColumn1 := gobless.NewColumn(
		gobless.GridSizeOneQuarter,
		citiesBox1,
	)

	citiesBox2 := gobless.NewTextBox()
	citiesBox2.SetTextWrap(true)
	citiesColumn2 := gobless.NewColumn(
		gobless.GridSizeOneQuarter,
		citiesBox2,
	)

	dayCounterBox := gobless.NewTextBox()
	dayCounterBox.SetText("Day: 0")

	ControllerBox := gobless.NewTextBox()
	ControllerBox.SetText("Control + Q: Close | Control + A: Time speed down | Control + S: Time speed up")

	rows := []gobless.Component{
		gobless.NewRow(
			gobless.GridSizeThreeQuarters,
			logsColumn,
			citiesColumn1,
			citiesColumn2,
		),
		gobless.NewRow(
			gobless.GridSizeOneQuarter,
			gobless.NewColumn(
				gobless.GridSizeHalf,
				dayCounterBox,
			),
			gobless.NewColumn(
				gobless.GridSizeHalf,
				ControllerBox,
			),
		),
	}

	gui.HandleKeyPress(gobless.KeyCtrlQ, func(event gobless.KeyPressEvent) {
		gui.Close()
	})
	gui.HandleKeyPress(gobless.KeyCtrlS, func(event gobless.KeyPressEvent) {
		atomic.AddInt64(&manager.WaitTime, -atomic.LoadInt64(&manager.WaitTime)/3)
	})
	gui.HandleKeyPress(gobless.KeyCtrlA, func(event gobless.KeyPressEvent) {
		atomic.AddInt64(&manager.WaitTime, atomic.LoadInt64(&manager.WaitTime)/3)
	})

	gui.Render(rows...)

	go func() {
		for {
			select {
			case log := <-manager.LogsCh:
				manager.logs = log + "\n" + manager.logs
				logsBox.SetText(manager.logs)
			case info := <-manager.DayCounterCh:
				dayCounterBox.SetText(info)
			case cities := <-manager.CitiesCh:
				cities1 := cities[:len(cities)/2]
				citiesBox1text := strings.Join(cities1, "\n")
				citiesBox1.SetText(citiesBox1text)

				cities2 := cities[len(cities)/2:]
				citiesBox2text := strings.Join(cities2, "\n")
				citiesBox2.SetText(citiesBox2text)
			}
			gui.Render(rows...)
		}
	}()
	gui.Loop()

	return nil
}
