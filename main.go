package main

import (
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/gdamore/tcell/v2"
)

var aimingstyle tcell.Style = tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorRed)

// choose = choosing action, waitforkeypress = Wait for Key Press, idle = not player turn, move = moving action, moved = just moved a step (for checking barriers), attack = attack action, youcannotreach = show cannot reach msg, enemy# = attack enemy, moveattack = aim attack, movedattack = choose aim attack, inventory = view inventory
var playerstate string = "choose"

// Keegan's Health
var (
	hp    int = 10
	maxhp int = 10
	armor int = 0
)

func drawText(s tcell.Screen, x, y int, text string) {
	row := y
	col := x
	defStyle := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)

	for _, r := range []rune(text) {
		s.SetContent(col, row, r, nil, defStyle)
		col++
	}
}

func drawTextStyle(s tcell.Screen, x, y int, style tcell.Style, text string) {
	row := y
	col := x

	for _, r := range []rune(text) {
		s.SetContent(col, row, r, nil, style)
		col++
	}
}

func testmap(s tcell.Screen) {
	// Keegan
	x := 3
	y := 3
	bx := x
	by := y
	equiped := "Stick"
	strength := 0

	// Keegan Aim
	ax := x
	ay := y

	// Enemy
	ex := 1
	ey := 1

	// ehp := 10
	steps := 0
	controltxt := ""
	hudtxt := ""
	playerstate = "choose"

	r := rand.New(rand.NewSource(time.Now().UnixMicro()))

	println(bx, by)
	for {
		r = rand.New(rand.NewSource(time.Now().UnixMicro()))

		if playerstate == "idle" {
			steps = 0
			enemyhit := false
			if x == 1 || x == 2 {
				if y == 1 || y == 2 {
					enemyhit = true
				}
			}
			if enemyhit {
			}

			playerstate = "waitforkeypress"
		}

		if equiped == "Stick" {
			strength = 1
		}

		if playerstate == "enemy1" {
			randnum := r.Intn(19) + 1

			if randnum < 10 {
				hudtxt = "You missed with a roll of " + strconv.Itoa(randnum) + "."
				playerstate = "waitforkeypress"
			} else if randnum == 20 {
				hudtxt = "You got critical hit with a damage of " + strconv.Itoa(strength+4) + "!"
				playerstate = "waitforkeypress"
			} else {
				hudtxt = "You got a hit with a damage of " + strconv.Itoa(strength) + "!"
				playerstate = "waitforkeypress"
			}

			controltxt = "Press any key to continue..."
		}

		if playerstate == "choose" {
			controltxt = "[m]ove [a]ttack [s]tats [i]nventory"
			hudtxt = "HP: " + strconv.Itoa(hp) + "/" + strconv.Itoa(maxhp) + ", Armor: " + strconv.Itoa(armor) + ", Weapon: " + equiped + ", Status: Choosing Action"
		} else if playerstate == "move" {
			hudtxt = "HP: " + strconv.Itoa(hp) + "/" + strconv.Itoa(maxhp) + ", Armor: " + strconv.Itoa(armor) + ", Status: Moving"
			controltxt = "Steps: " + strconv.Itoa(steps) + "/6"
		} else if playerstate == "attack" {
			hudtxt = "HP: " + strconv.Itoa(hp) + "/" + strconv.Itoa(maxhp) + ", Armor: " + strconv.Itoa(armor) + ", Weapon: " + equiped + ", Status: Attacking"
			controltxt = "Attack Here? (y/n)"
		}

		if playerstate == "choose" && steps == 6 {
			controltxt = "[a]ttack [s]tats [i]nventory"
		}

		if playerstate == "youcannotreach" {
			hudtxt = "You cannot reach that far!"
			controltxt = "Press any key to continue..."
		}

		s.Clear()
		drawText(s, 0, 0, "------")
		drawText(s, 0, 1, "|    |")
		drawText(s, 0, 2, "|    |")
		drawText(s, 0, 3, "|    |")
		drawText(s, 0, 4, "------")
		drawText(s, x, y, "K")
		drawTextStyle(s, ex, ey, aimingstyle, "B")

		if playerstate == "attack" {
			drawTextStyle(s, ax, ay, aimingstyle, "+")
		}

		// Draw HUD
		drawText(s, 0, 6, hudtxt)

		// Draw Controls
		drawText(s, 0, 8, controltxt)

		s.Sync()

		for {
			// Update screen
			s.Show()

			// Poll event
			ev := s.PollEvent()

			// Process event
			switch ev := ev.(type) {
			case *tcell.EventResize:
				s.Sync()
			case *tcell.EventKey:
				if ev.Rune() == 'm' {
					if playerstate != "move" && steps != 6 {
						// moving
						playerstate = "move"
					} else if playerstate == "move" {
						playerstate = "choose"
					}
				} else if ev.Rune() == 'a' {
					// attack or move in moving state
					if playerstate == "move" {
						bx = x
						by = y
						x -= 1
						steps += 1
						playerstate = "moved"
					} else if playerstate != "attack" {
						ax = ex
						ay = ey
						playerstate = "attack"
					}
				} else if ev.Rune() == 's' {
					// check stats or move in moving state
					if playerstate == "move" {
						bx = x
						by = y
						y += 1
						steps += 1
						playerstate = "moved"
					}
				} else if ev.Rune() == 'd' {
					// for moving
					if playerstate == "move" {
						bx = x
						by = y
						x += 1
						steps += 1
						playerstate = "moved"
					}
				} else if ev.Rune() == 'w' {
					// for moving
					if playerstate == "move" {
						bx = x
						by = y
						y -= 1
						steps += 1
						playerstate = "moved"
					}
				} else if ev.Rune() == 'n' {
					if playerstate == "attack" {
						playerstate = "choose"
					}
				} else if ev.Rune() == 'y' {
					if playerstate == "attack" {
						cantreach := false
						if ex == 1 && ey == 1 {
							if x == 1 || x == 2 {
								if y == 1 || y == 2 {
									playerstate = "enemy1"
								} else {
									cantreach = true
								}
							} else {
								cantreach = true
							}
						}

						if cantreach {
							playerstate = "youcannotreach"
						}
					}
				}
			}

			if playerstate == "waitforkeypress" {
				playerstate = "idle"
			}

			if playerstate == "youcannotreach" {
				playerstate = "choose"
			}

			if playerstate == "moved" {
				// Barriers Checks here
				if x == 0 {
					x += 1
					steps -= 1
				}
				if y == 0 {
					y += 1
					steps -= 1
				}
				if x == 5 {
					x -= 1
					steps -= 1
				}
				if y == 4 {
					y -= 1
					steps -= 1
				}
				if x == 1 && y == 1 {
					x = bx
					y = by
					steps -= 1
				}

				if steps != 6 {
					playerstate = "move"
				} else {
					playerstate = "choose"
				}
			}

			break
		}
	}
}

func main() {
	s, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("%+v", err)
	}
	if err := s.Init(); err != nil {
		log.Fatalf("%+v", err)
	}

	// Set default text style
	defStyle := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)
	s.SetStyle(defStyle)

	// Clear screen
	s.Clear()
	drawText(s, 0, 0, "Ghost Team: A Unfriendly Meeting")
	drawText(s, 0, 2, "Want to play? (y/n)")
	s.Sync()

	quit := func() {
		s.Fini()
		os.Exit(0)
	}
	for {
		// Update screen
		s.Show()

		// Poll event
		ev := s.PollEvent()

		// Process event
		switch ev := ev.(type) {
		case *tcell.EventResize:
			s.Sync()
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlC {
				quit()
			} else if ev.Rune() == 'y' {
				// start game here
				testmap(s)
			} else if ev.Rune() == 'n' {
				s.Fini()
				println("Bye!")
				os.Exit(0)
			}
		}
	}
}
