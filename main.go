package main

import (
	"log"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/gopxl/beep"
	"github.com/gopxl/beep/mp3"
	"github.com/gopxl/beep/speaker"
)

// Styles
var (
	aimingstyle  tcell.Style = tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorRed)
	commentstyle tcell.Style = tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorGray).Italic(true)
)

// choose = choosing action, waitforkeypress = Wait for Key Press, idle = not player turn, move = moving action, moved = just moved a step (for checking barriers), attack = attack action, youcannotreach = show cannot reach msg, enemy# = attack enemy, moveattack = aim attack, movedattack = choose aim attack, noenemy = show no enemy msg, inventory = view inventory
var playerstate string = "choose"

// Keegan's Stats
var (
	hp         int    = 10
	maxhp      int    = 10
	armor      int    = 0
	armorname  string = "Military Clothes" // Display Armor Name
	strength   int    = 10
	weaponname string = "Pistol" // Weapon Name and Damage Checker
	items             = []string{"Stick", "Pistol"}
	firstname  string = "Keegan"
	lastname   string = "Miller"
)

// Weapons on map
var pistolonmap = false

// Draw Text with Tcell
func drawText(s tcell.Screen, x, y int, text string) {
	row := y
	col := x
	defStyle := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)

	for _, r := range string(text) {
		s.SetContent(col, row, r, nil, defStyle)
		col++
	}
}

// Draw Text with Tcell, and a custom style for that text
func drawTextStyle(s tcell.Screen, x, y int, style tcell.Style, text string) {
	row := y
	col := x

	for _, r := range string(text) {
		s.SetContent(col, row, r, nil, style)
		col++
	}
}

func stats_help(s tcell.Screen) {
	s.Clear()
	drawText(s, 0, 0, "Stats:")
	drawText(s, 0, 2, "Stats show many things that have to do with your character.")
	drawText(s, 0, 4, "Example of a stat:")
	drawText(s, 0, 6, "Weapon (Damage): Pistol (󰇏 💥10)")
	drawText(s, 0, 8, "Press any key to continue...")

	s.Show()

	// Poll event
	s.PollEvent()

	s.Clear()
	drawText(s, 0, 0, "Symbols:")
	drawText(s, 0, 2, "💥  : Damage Number Indicator")
	drawText(s, 0, 3, "󰇏  : Weapon rolls dice (based on damage amount, ex. 10 damage would be a d10)")
	drawText(s, 0, 4, "#󰇏 : Weapon rolls a # amount of dice (based on damage amount)")
	drawText(s, 0, 5, "  : Armor Number Indicator")
	drawText(s, 0, 7, "Press any key to quit help...")

	// Update screen
	s.Show()

	// Poll event
	s.PollEvent()
}

func yourstats(s tcell.Screen) {
	dicesymbol := ""

	if weaponname == "Pistol" {
		dicesymbol = "󰇏 "
	}

	statsdisplay := func(dice string) {
		s.Clear()
		drawText(s, 0, 0, "Name: "+firstname+" "+lastname)
		drawText(s, 0, 2, "Health: "+strconv.Itoa(hp)+"/"+strconv.Itoa(maxhp))
		drawText(s, 0, 3, "Weapon (Damage): "+weaponname+" ("+dice+"💥 "+strconv.Itoa(strength)+")")
		drawText(s, 0, 4, "Armor (Defense): "+armorname+" ( "+strconv.Itoa(armor)+")")
		drawText(s, 0, 6, "Press ? for help, or any other key to go back...")
	}

	statsdisplay(dicesymbol)

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
			if ev.Rune() == '?' {
				stats_help(s)
				statsdisplay(dicesymbol)
				s.Show()
			} else {
				return
			}
		}
	}
}

func rolld4() (roll int) {
	r := rand.New(rand.NewSource(time.Now().UnixMicro()))
	roll = r.Intn(3) + 1

	return
}

func rolld6() (roll int) {
	r := rand.New(rand.NewSource(time.Now().UnixMicro()))
	roll = r.Intn(5) + 1

	return
}

func rolld8() (roll int) {
	r := rand.New(rand.NewSource(time.Now().UnixMicro()))
	roll = r.Intn(7) + 1

	return
}

func rolld10() (roll int) {
	r := rand.New(rand.NewSource(time.Now().UnixMicro()))
	roll = r.Intn(9) + 1

	return
}

func rolld12() (roll int) {
	r := rand.New(rand.NewSource(time.Now().UnixMicro()))
	roll = r.Intn(11) + 1

	return
}

func rolld20() (roll int) {
	r := rand.New(rand.NewSource(time.Now().UnixMicro()))
	roll = r.Intn(19) + 1

	return
}

func startattackplayer(hitchance int) (damage int, hit bool, crit bool, roll int) {
	// r := rand.New(rand.NewSource(time.Now().UnixMicro()))
	roll = rolld20()

	if roll < hitchance {
		hit = false
		damage = 0
	} else if roll != 20 {
		hit = true
	} else {
		hit = true
		crit = true
	}

	if hit {
		// Check weapon
		if weaponname == "Stick" {
			damage = 1
			if crit {
				damage += 1
			}
		} else if weaponname == "Pistol" {
			damage = rolld10()
			if crit {
				damage += rolld10()
			}
		}
	}

	return
}

func playMP3(filename string) {
	var wg sync.WaitGroup
	wg.Add(1)

	go helperPlayMP3(filename, &wg)
}

func helperPlayMP3(filename string, wg *sync.WaitGroup) {
	defer wg.Done()
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}

	streamer, format, err := mp3.Decode(f)
	if err != nil {
		log.Fatal(err)
	}
	defer streamer.Close()

	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))

	done := make(chan bool)
	speaker.Play(beep.Seq(streamer, beep.Callback(func() {
		done <- true
	})))

	<-done
}

func testmap(s tcell.Screen) {
	// Keegan
	x := 3
	y := 3
	bx := x
	by := y

	// Keegan Aim
	ax := x
	ay := y

	// Enemy
	ex := 1
	ey := 1
	ehp := 10

	// ehp := 10
	steps := 0
	controltxt := ""
	hudtxt := ""
	playerstate = "choose"
	beingattacked := false

	// r := rand.New(rand.NewSource(time.Now().UnixMicro()))

	println(bx, by)
	for {
		// r = rand.New(rand.NewSource(time.Now().UnixMicro()))

		if playerstate == "idle" {
			steps = 0
			enemyhit := false

			if x == 1 || x == 2 {
				if y == 1 || y == 2 {
					if ehp > 0 {
						enemyhit = true
					}
				}
			}

			if enemyhit {
				hudtxt = "The enemy cutout falls over, and cuts you. You lost 1 HP (but you have infinity health)."
			} else {
				if ehp > 0 {
					hudtxt = "The enemy cutout cannot do anything to you."
				}
			}

			if ehp > 0 {
				controltxt = "Press any key to continue..."
				playerstate = "waitforkeypress"
				beingattacked = true
			} else {
				playerstate = "choose"
			}
		}

		if playerstate == "enemy1" {
			damage, hit, crit, _ := startattackplayer(10)

			if hit {
				if crit {
					hudtxt = "You got critical hit with a damage of " + strconv.Itoa(damage) + "!"
				} else {
					hudtxt = "You got a hit with a damage of " + strconv.Itoa(damage) + "!"
				}

				ehp -= damage
			} else {
				hudtxt = "You missed."
			}

			// if randnum < 10 {
			// 	hudtxt = "You missed with a roll of " + strconv.Itoa(randnum) + "."
			// 	playerstate = "waitforkeypress"
			// } else if randnum == 20 {
			// 	hudtxt = "You got critical hit with a damage of " + strconv.Itoa(strength+4) + "!"
			// 	playerstate = "waitforkeypress"
			// 	ehp -= strength + 4
			// } else {
			// 	hudtxt = "You got a hit with a damage of " + strconv.Itoa(strength) + "!"
			// 	playerstate = "waitforkeypress"
			// 	ehp -= strength
			// }

			controltxt = "Press any key to continue..."

			playerstate = "waitforkeypress"
		}

		if playerstate == "choose" {
			controltxt = "[m]ove [a]ttack/action [s]tats [i]nventory [e]nd turn"
			hudtxt = "HP: " + strconv.Itoa(hp) + "/" + strconv.Itoa(maxhp) + ", Armor: " + strconv.Itoa(armor) + ", Weapon: " + weaponname + ", Status: Choosing Action"
		} else if playerstate == "move" {
			hudtxt = "HP: " + strconv.Itoa(hp) + "/" + strconv.Itoa(maxhp) + ", Armor: " + strconv.Itoa(armor) + ", Status: Moving"
			controltxt = "Steps: " + strconv.Itoa(steps) + "/6"
		} else if playerstate == "attack" {
			hudtxt = "HP: " + strconv.Itoa(hp) + "/" + strconv.Itoa(maxhp) + ", Armor: " + strconv.Itoa(armor) + ", Weapon: " + weaponname + ", Status: Attacking"
			controltxt = "Attack Here? (y/n)"
		}

		if playerstate == "choose" && steps == 6 {
			controltxt = "[a]ttack/action [s]tats [i]nventory [e]nd turn"
		}

		if playerstate == "youcannotreach" {
			hudtxt = "You cannot reach that far!"
			controltxt = "Press any key to continue..."
		}

		if playerstate == "noenemy" {
			hudtxt = "There is no enemy to hit!"
			controltxt = "Press any key to continue..."
			beingattacked = true
			playerstate = "waitforkeypress"
		}

		s.Clear()
		drawText(s, 0, 0, "------")
		drawText(s, 0, 1, "|    |")
		drawText(s, 0, 2, "|    |")
		drawText(s, 0, 3, "|    |")
		drawText(s, 0, 4, "------")
		drawText(s, x, y, "K")

		if ehp > 0 {
			drawTextStyle(s, ex, ey, aimingstyle, "B")
		}

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
				// Quiting for debug
				if ev.Rune() == 'q' {
					s.Fini()
					os.Exit(0)
				}

				if playerstate != "waitforkeypress" {
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
							playMP3("move.mp3")
						} else if playerstate != "attack" {
							ax = ex
							ay = ey
							if ehp > 0 {
								playerstate = "attack"
							} else {
								playerstate = "noenemy"
							}
						}
					} else if ev.Rune() == 's' {
						// check stats or move in moving state
						if playerstate == "move" {
							bx = x
							by = y
							y += 1
							steps += 1
							playerstate = "moved"
							playMP3("move.mp3")
						} else {
							yourstats(s)
						}
					} else if ev.Rune() == 'd' {
						// for moving
						if playerstate == "move" {
							bx = x
							by = y
							x += 1
							steps += 1
							playerstate = "moved"
							playMP3("move.mp3")
						}
					} else if ev.Rune() == 'w' {
						// for moving
						if playerstate == "move" {
							bx = x
							by = y
							y -= 1
							steps += 1
							playerstate = "moved"
							playMP3("move.mp3")
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
										if weaponname != "Pistol" {
											cantreach = true
										} else {
											playerstate = "enemy1"
										}
									}
								} else {
									if weaponname != "Pistol" {
										cantreach = true
									} else {
										playerstate = "enemy1"
									}
								}
							}

							if cantreach {
								playerstate = "youcannotreach"
							}
						}
					} else if ev.Rune() == 'e' {
						playerstate = "idle"
					}
				}
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
					if ehp <= 0 {
						x = bx
						y = by
						steps -= 1
					}
				}

				if steps != 6 {
					playerstate = "move"
				} else {
					playerstate = "choose"
				}
			}

			break
		}

		if playerstate == "waitforkeypress" {
			if beingattacked {
				playerstate = "choose"
				beingattacked = false
			} else {
				playerstate = "idle"
			}
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
