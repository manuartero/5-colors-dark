package core

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/martero/herocademy/core/time"
	"github.com/martero/herocademy/models"
)

type eventListener interface {
	OnA()
	OnC()
	OnP()
	OnS()
	OnT()
	OnUp()
	OnDown()
	OnEnter()
	OneSecond()
}

func (game gameImpl) Start() {
	game.PrintDate(game.Now())
	game.FullDraw()

	for i := 0; i < 5; i++ {
		hero := game.NewHero()
		game.AddHero(hero)
		game.PrintHeroInTable(hero, i)
		game.Log(fmt.Sprintf("Hero: %s ready to rock [%s]", hero.HeroName(), hero.ID()))
	}
}

func (game gameImpl) OnA() {
	done := make(chan bool, 1)
	go game.PrintAventureOptionsInTheMenu()
	go game.seekForAdventure(done)
	<-done
}

func (game gameImpl) seekForAdventure(done chan bool) {
	adventure := game.NewAdventure()
	game.AdventureInProgress()
	game.PrintInDetailBox(adventure)
	game.Log(fmt.Sprintf("Starting adventure \"%s\"", adventure))
	done <- true
}

func (game gameImpl) OnC() {
	go game.CleanMenuOptions()
	if game.IsAdventureInProgress() {
		game.Log("Scaping from current adventure")
	}
	game.InitialStatus()
}

func (game gameImpl) OnP() {
	if game.IsHeroTableFocused() {
		hero := game.HeroStartsPatrolling()
		game.PrintHeroStatusInTable(game.GetHeroStatus(models.Patrolling))
		game.Log(fmt.Sprintf("%s starts patrolling", hero.HeroName()))
	}
}

func (game gameImpl) OnS() {
	if game.IsHeroTableFocused() {
		hero, prevAction := game.HeroStopsCurrentAction()
		game.PrintHeroStatusInTable(game.GetHeroStatus(models.Available))
		game.Log(fmt.Sprintf("%s stops %d", hero.HeroName(), prevAction))
	}
}

func (game gameImpl) OnT() {
	if game.IsHeroTableFocused() {
		hero := game.StartTraining()
		game.PrintHeroStatusInTable(game.GetHeroStatus(models.Training))
		game.Log(fmt.Sprintf("%s starts training", hero.HeroName()))
	} else {
		game.FocusHeroTable()
		hero, relations := game.NextHero()
		game.PrintFocusInHeroTable(foo(relations))
		game.PrintInDetailBox(hero)
		game.PrintTableOptionsInTheMenu()
	}
}

func (game gameImpl) OnUp() {
	if game.IsHeroTableFocused() {
		hero, relations := game.PrevHero()
		game.Log(fmt.Sprintf("%v", relations))
		game.FocusPrevHeroInTable(foo(relations))
		game.PrintInDetailBox(hero)
	}
}

func (game gameImpl) OnDown() {
	if game.IsHeroTableFocused() {
		hero, relations := game.NextHero()
		game.Log(fmt.Sprintf("%s => %v | %v", hero.ID(), relations, foo(relations)))
		game.FocusNextHeroInTable(foo(relations))
		game.PrintInDetailBox(hero)
	}
}

func (game gameImpl) OnEnter() {
	if !game.IsAdventureInProgress() {
		if game.AllHeroesBusy() {
			// TODO calculate loop
			game.PassTime(time.FewMonths)
			game.PrintDate(game.Now())
		} else {
			game.PrintAlert("Some heroes are idle! specify an action for all of them")
		}
	}
}

func (game gameImpl) OneSecond() {
	if game.IsAdventureInProgress() {
		// TODO
		game.Log("one second in adventure")
	}
}

// FIXME temporal solution, when removing heroes this will break
func foo(m map[string]int) map[int]int {
	response := make(map[int]int)
	for id, v := range m {
		tmp := strings.Split(id, "#")[1]
		pos, _ := strconv.Atoi(tmp)
		response[pos] = v
	}
	return response
}
