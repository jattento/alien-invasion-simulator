package client

import (
	"fmt"
	"math/rand"
	"os"
	"sync/atomic"
	"time"

	"github.com/jattento/alien-invasion-simulator/internal/earth"
	"github.com/jattento/alien-invasion-simulator/internal/interface/terminal"
	"github.com/jattento/alien-invasion-simulator/internal/simulation"
)

type Simulation interface {
	Tick() (bool, simulation.TickReport)
	Cities() map[string]map[earth.Direction]string
}

// Run blocks until the program is ended or an error happen
func Run(invSimulation Simulation, aliens int) error {
	logsCh := make(chan string)
	DaysCh := make(chan string)
	citiesCh := make(chan []string)

	worldMatrix := worldMap{cities: make([]city, 0), citiesIndex: make(map[string]int), alive: aliens}
	for cityInfo := range invSimulation.Cities() {
		worldMatrix.save(city{name: cityInfo})
	}

	t := terminal.New(os.Stdout, logsCh, DaysCh, citiesCh)

	go func() {
		for keepTicking := true; keepTicking && worldMatrix.alive > 0; {
			now := time.Now()

			var report simulation.TickReport
			keepTicking, report = invSimulation.Tick()

			worldMatrix.clear()

			for _, battleReport := range report.Battles {
				worldMatrix.save(city{name: battleReport.City, destroyed: true, aliens: battleReport.InvolvedAliens})
				logsCh <- killLog(battleReport)
			}
			for cityName := range report.AlienPositions {
				worldMatrix.save(city{name: cityName, aliens: report.AlienPositions[cityName]})
			}

			citiesCh <- worldMatrix.prettySlice()
			DaysCh <- fmt.Sprintf("🕒  :  %v   |   👽  :  %v   |   💀  :  %v   |   🏡  :  %v   |   🔥  :  %v",
				report.Tick, worldMatrix.alive, worldMatrix.dead, worldMatrix.notDestroyed, worldMatrix.destroyed)

			time.Sleep(time.Duration(atomic.LoadInt64(&t.WaitTime)) - (time.Now().Sub(now)))
		}

		logsCh <- "Congratulations on completing the alien simulation! "
		logsCh <- "Just remember, if any actual aliens come to visit, don't blame me if this isn't accurucate."

	}()

	if err := t.Run(); err != nil {
		return err
	}

	return nil
}

func killLog(report earth.BattleReport) string {
	skull := "💀"
	knife := "🗡️"
	gun := "🔫"
	bomb := "💣"
	wrench := "🔧"
	poison := "🧪"
	syringe := "💉"
	fire := "🔥"
	paperClip := "📎"

	weapons := []string{skull, knife, gun, bomb, wrench, poison, syringe, fire, paperClip}
	weapon := weapons[randomInt(0, len(weapons)-1)]

	fmtText := fmt.Sprintf("👽 %q", report.InvolvedAliens[0])
	for i := 1; i < len(report.InvolvedAliens); i++ {
		if i == len(report.InvolvedAliens)-1 {
			fmtText += fmt.Sprintf(" and 👽 %q", report.InvolvedAliens[i])
			continue
		}
		fmtText += fmt.Sprintf(", 👽 %q", report.InvolvedAliens[i])
	}

	return fmt.Sprintf("%s killed each other in %q in a %s  duel %s",
		fmtText, report.City, weapon, skull)
}

func randomInt(min int, max int) int {
	return min + rand.Intn(max-min+1)
}
