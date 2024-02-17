package main

import (
	"log"
	"math"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/gdamore/tcell/v2"
)

// Debug Menu Enabler (for the cheaters)
var DEBUG bool = true

// Styles
var (
	aimingstyle   tcell.Style = tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorRed)
	commentstyle  tcell.Style = tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorGray).Italic(true)
	grassstyle    tcell.Style = tcell.StyleDefault.Background(tcell.ColorGreen).Foreground(tcell.ColorBlack)
	lightbluetext tcell.Style = tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorLightBlue)
	yellowtext    tcell.Style = tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorYellow)
)

// choose = choosing action, waitforkeypress = Wait for Key Press, idle = not player turn, move = moving action, moved = just moved a step (for checking barriers), wantmove = y/n for move, attack = attack action, youcannotreach = show cannot reach msg, enemy# = attack enemy, moveattack = aim attack, movedattack = choose aim attack, noenemy = show no enemy msg, inventory = view inventory
var playerstate string = "choose"

// Keegan's Stats
var (
	hp                 int         = 10
	maxhp              int         = 10
	armor              int         = 0
	armorname          string      = "Military Clothes" // Display Armor Name
	strength           int         = 10
	weaponname         string      = "Pistol" // Weapon Name and Damage Checker
	items                          = []string{"Stick", "Pistol"}
	firstname          string      = "Keegan"
	lastname           string      = "Miller"
	keeganstyle        tcell.Style = tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)
	defaultkeeganstyle tcell.Style = tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)
)

type maphandle struct {
	givengrounds []string
	givenx       []int
	giveny       []int
	blockx       []int
	blocky       []int
	pathx        []int
	pathy        []int
	pathdir      []string
	chestitems   []string
}

// https://stackoverflow.com/questions/15323767/does-go-have-if-x-in-construct-similar-to-python
func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// get distance between two points
func getDistance(x int, y int, x2 int, y2 int) (distance float64) {
	newx, newy := float64(x), float64(y)
	newx2, newy2 := float64(x2), float64(y2)

	return math.Sqrt(math.Pow(newx2-newx, 2) + math.Pow(newy2-newy, 2))
}

// cool function to see if your weapon can hit the enemy lol
//
// depending on the weapon, it might make it not be able to hit, or set it to disadvantage
//
// if good = false, and atdis = true, then the weapon will hit, but be at disadvantage
func WeaponDistanceGood(distance float64, weapon string) (good bool, atdis bool) {
	atdis = false

	if weapon == "Stick" {
		if distance <= 1.6 {
			good = true
		} else {
			good = false
		}
	} else if weapon == "Pistol" {
		if distance <= 5.0 {
			good = true
		} else if distance <= 8.0 {
			good = false
			atdis = true
		} else {
			good = false
		}
	}

	return
}

// AddObj Ground Type to Map
func (m *maphandle) AddObj(objtype string, x int, y int) {
	checkforPath := []string{"pathup", "pathdown", "pathleft", "pathright"}

	if !stringInSlice(objtype, checkforPath) {
		m.givengrounds = append(m.givengrounds, objtype)
		m.givenx = append(m.givenx, x)
		m.giveny = append(m.giveny, y)
	}

	if objtype == "horiwall" || objtype == "vertwall" || objtype == "enemy" || objtype == "chest" {
		m.blockx = append(m.blockx, x)
		m.blocky = append(m.blocky, y)
	}

	if stringInSlice(objtype, checkforPath) {
		m.pathx = append(m.pathx, x)
		m.pathy = append(m.pathy, y)
		m.pathdir = append(m.pathdir, objtype)
	}
}

func (m maphandle) Show(s tcell.Screen) {
	blocki := 0
	i := 0

	for i < len(m.givengrounds) {
		if m.givengrounds[i] == "horiwall" {
			drawText(s, m.blockx[blocki], m.blocky[blocki], "-")
			blocki++
		} else if m.givengrounds[i] == "vertwall" {
			drawText(s, m.blockx[blocki], m.blocky[blocki], "|")
			blocki++
		} else if m.givengrounds[i] == "g" || m.givengrounds[i] == "grass" {
			drawTextStyle(s, m.givenx[i], m.giveny[i], grassstyle, " ")
		} else if m.givengrounds[i] == "enemy" {
			drawTextStyle(s, m.blockx[blocki], m.blocky[blocki], aimingstyle, "B")
			blocki++
		} else if m.givengrounds[i] == "chest" {
			drawTextStyle(s, m.blockx[blocki], m.blocky[blocki], yellowtext, "󰜦")
			blocki++
		}
		i++
	}
}

func (m *maphandle) AddChestItem(item string) {
	m.chestitems = append(m.chestitems, item)
}

func (m maphandle) CoordIsCollide(x int, y int) (collide bool) {
	i := 0
	collide = false

	for i < len(m.blockx) {
		if m.blockx[i] == x {
			if m.blocky[i] == y {
				collide = true
				return
			}
		}

		i++
	}

	return
}

func (m maphandle) FindObjectAtCoord(x int, y int, objtype string) (objatcoord bool, objx int, objy int) {
	i := 0
	objatcoord = false
	objx = 0
	objy = 0

	for i < len(m.givengrounds) {
		if m.givengrounds[i] == objtype {
			if m.givenx[i] == x && m.giveny[i] == y {
				objatcoord = true
				objx = m.givenx[i]
				objy = m.giveny[i]
				return
			}
		}

		i++
	}

	return
}

func (m maphandle) GetObjectAtCoord(x int, y int) (objatcoord bool, objtype string, objx int, objy int) {
	i := 0
	objatcoord = false
	objtype = ""
	objx = 0
	objy = 0

	for i < len(m.givengrounds) {
		if m.givenx[i] == x && m.giveny[i] == y {
			objatcoord = true
			objx = m.givenx[i]
			objy = m.giveny[i]
			objtype = m.givengrounds[i]
			return
		}

		i++
	}

	return
}

func (m maphandle) GroundType(x int, y int) (groundtype string) {
	i := 0

	for i < len(m.givenx) {
		if m.givenx[i] == x {
			if m.giveny[i] == y {
				break
			}
		}
		i++
	}

	if i != len(m.givenx) {
		if m.givengrounds[i] == "grass" || m.givengrounds[i] == "g" {
			groundtype = "grass"
		} else {
			groundtype = "none"
		}
	}

	return
}

func (m maphandle) getPathDirAtPoint(x int, y int) (dir string) {
	i := 0
	dir = "none"

	for i < len(m.pathx) {
		if x == m.pathx[i] {
			if y == m.pathy[i] {
				dir = m.pathdir[i]
				return
			}
		}
	}

	return
}

// mode: chase, away, range
//
// chase = chase the player
//
// away = run away from the player
//
// range = get in range with player
func (m maphandle) EnemyMove(ex int, ey int, x int, y int, erange int, mode string) (nex int, ney int) {
	// 2 means unused
	leftright := 2
	updown := 2

	steps := 0

	if x > ex {
		leftright = 1
	} else if x < ex {
		leftright = 0
	}

	if y > ey {
		updown = 1
	} else if y < ey {
		updown = 0
	}

	if mode == "chase" {
		for steps < 6 {
			exb := ex
			eyb := ey

			if ex == x-1 || ex == x+1 || ex == x {
				if ey == y+1 || ey == y-1 || ey == y {
					break
				}
			}

			if ex == x {
				leftright = 2
			}
			if ey == y {
				updown = 2
			}

			if leftright == 0 {
				ex--
				steps++
			} else if leftright == 1 {
				ex++
				steps++
			}

			if m.CoordIsCollide(ex, ey) {
				ex = exb
			}

			if updown == 0 {
				ey--
				steps++
			} else if updown == 1 {
				ey++
				steps++
			}

			if m.CoordIsCollide(ex, ey) {
				ey = eyb
			}

			if ex == x && ey == y {
				ex = exb
				ey = eyb
			}
		}
	} else if mode == "away" {
		invertupdown := false
		invertleftright := false

		if ey > y {
			invertupdown = true
			invertleftright = true
		}

		for steps < 6 {
			exb := ex
			eyb := ey

			dir := m.getPathDirAtPoint(ex, ey)

			if dir == "pathup" {
				if !invertupdown {
					updown = 0
				} else {
					updown = 1
				}
				leftright = 2
			} else if dir == "pathdown" {
				if !invertupdown {
					updown = 1
				} else {
					updown = 0
				}
				leftright = 2
			} else if dir == "pathleft" {
				if !invertleftright {
					leftright = 0
				} else {
					leftright = 1
				}
				updown = 2
			} else if dir == "pathright" {
				if !invertleftright {
					leftright = 1
				} else {
					leftright = 0
				}
				updown = 2
			}

			if leftright == 0 {
				ex--
				steps++
			} else if leftright == 1 {
				ex++
				steps++
			}

			if m.CoordIsCollide(ex, ey) {
				ex = exb
			}

			if updown == 0 {
				ey--
				steps++
			} else if updown == 1 {
				ey++
				steps++
			}

			if m.CoordIsCollide(ex, ey) {
				ey = eyb
			}

			if ex == x && ey == y {
				ex = exb
				ey = eyb
			}
		}
	}

	nex = ex
	ney = ey

	return
}

// probably going to be unused, but could be helpful for other projects
func (m *maphandle) RemoveObj(x int, y int) {
	i := 0

	for i < len(m.givenx) {
		if m.givenx[i] == x {
			if m.giveny[i] == y {
				break
			}
		}

		i++
	}

	bi := 0

	for bi < len(m.blockx) {
		if m.blockx[bi] == x {
			if m.blocky[bi] == y {
				break
			}
		}

		bi++
	}

	if i != len(m.givenx) {
		// givenx
		newArr := make([]int, len(m.givenx)-1)
		copy(newArr[:i], m.givenx[:i])
		copy(newArr[i:], m.givenx[i+1:])

		m.givenx = newArr

		// giveny
		newArr = make([]int, len(m.giveny)-1)
		copy(newArr[:i], m.giveny[:i])
		copy(newArr[i:], m.giveny[i+1:])

		m.giveny = newArr

		// givengrounds
		newstrArr := make([]string, len(m.givengrounds)-1)
		copy(newstrArr[:i], m.givengrounds[:i])
		copy(newstrArr[i:], m.givengrounds[i+1:])
	}

	if bi != len(m.givenx) {
		// blockx
		newArr := make([]int, len(m.blockx)-1)
		copy(newArr[:bi], m.blockx[:bi])
		copy(newArr[bi:], m.blockx[bi+1:])

		m.blockx = newArr

		// blocky
		newArr = make([]int, len(m.blocky)-1)
		copy(newArr[:bi], m.blocky[:bi])
		copy(newArr[bi:], m.blocky[bi+1:])
		m.blocky = newArr
	}
}

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
	drawText(s, 0, 6, "Weapon (Damage, Distance): Pistol (󰇏 💥 10, 󱡀 5-8)")
	drawText(s, 0, 8, "Press any key to continue...")

	s.Show()

	// Poll event
	s.PollEvent()

	s.Clear()
	drawText(s, 0, 0, "Symbols:")
	drawText(s, 0, 2, "💥  : Damage Number Indicator")
	drawText(s, 0, 3, "󱡀  : Weapon Distance (if distance has a -, then past the first number it is disadvantage, then miss)")
	drawText(s, 0, 4, "󰇏  : Weapon rolls dice (based on damage amount, ex. 10 damage would be a d10)")
	drawText(s, 0, 5, "#󰇏 : Weapon rolls a # amount of dice (based on damage amount)")
	drawText(s, 0, 6, "  : Armor Number Indicator")
	drawText(s, 0, 8, "Press any key to quit help...")

	// Update screen
	s.Show()

	// Poll event
	s.PollEvent()
}

func yourstats(s tcell.Screen) {
	dicesymbol := ""
	distance := ""

	if weaponname == "Pistol" {
		dicesymbol = "󰇏 "
		distance = "5-8"
	} else if weaponname == "Stick" {
		distance = "1"
	}

	statsdisplay := func(dice string) {
		s.Clear()
		drawText(s, 0, 0, "Name: "+firstname+" "+lastname)
		drawText(s, 0, 2, "Health: "+strconv.Itoa(hp)+"/"+strconv.Itoa(maxhp))
		drawText(s, 0, 3, "Weapon (Damage, Distance): "+weaponname+" ("+dice+"💥 "+strconv.Itoa(strength)+", 󱡀 "+distance+")")
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

// pimp named slickback
func debugmenu(s tcell.Screen) {
	s.Clear()
	drawText(s, 0, 0, "DEBUG MENU")
	drawTextStyle(s, 0, 2, lightbluetext, "[1] +Keegan")
	drawText(s, 0, 3, "[2] YES")
}

func testmap(s tcell.Screen) {
	// Keegan
	y := 3
	x := 3
	bx := x
	by := y
	canattack := true

	// Keegan Aim
	ax := x
	ay := y

	// Keegan Move
	kx := x
	ky := y
	movestyle := keeganstyle

	// Enemy
	ex := 1
	ey := 1
	ehp := 10

	// Enemy Move
	nex := ex
	ney := ey

	// ehp := 10
	steps := 0
	bsteps := 0
	controltxt := ""
	hudtxt := ""
	playerstate = "choose"
	beingattacked := false
	enemymoving := false
	disadvantage := false

	// r := rand.New(rand.NewSource(time.Now().UnixMicro()))

	println(bx, by)
	for {
		// r = rand.New(rand.NewSource(time.Now().UnixMicro()))
		gamemap := maphandle{}

		// objects created before state checks, and player controls (or else enemy moves through walls)
		for i := 0; i <= 5; i++ {
			gamemap.AddObj("horiwall", i, 0)
			gamemap.AddObj("horiwall", i, 10)
		}
		for i := 1; i <= 9; i++ {
			gamemap.AddObj("vertwall", 0, i)
			gamemap.AddObj("vertwall", 5, i)
		}

		if ehp > 0 {
			gamemap.AddObj("enemy", ex, ey)
		}

		gamemap.AddObj("g", 1, 3)

		if playerstate == "attack" && !canattack {
			hudtxt = "You already used your attack!"
			controltxt = "Press any key to continue..."
			playerstate = "waitforkeypress"
			beingattacked = true
		}

		if playerstate == "idle" {
			steps = 0
			enemyhit := false

			if ehp > 0 {
				if ex == x-1 || ex == x+1 || ex == x {
					if ey == y+1 || ey == y-1 || ey == y {
						enemyhit = true
					} else {
						enemymoving = true
					}
				} else {
					enemymoving = true
				}
			}

			if enemymoving {
				nex, ney = gamemap.EnemyMove(ex, ey, x, y, 0, "chase")
				hudtxt = "The enemy cutout is sliding towards you."
			}

			if enemyhit {
				hudtxt = "The enemy cutout falls over, and cuts you. You lost 1 HP."
				hp--
			} else {
				if ehp > 0 && !enemymoving {
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

			bsteps = 0
			canattack = true
		}

		if playerstate == "enemy1" {
			damage, hit, crit := 0, false, false

			if disadvantage {
				damage, hit, crit, _ = startattackplayer(10 + 4)
			} else {
				damage, hit, crit, _ = startattackplayer(10)
			}

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

			beingattacked = true
			canattack = false
			playerstate = "waitforkeypress"
		}

		if playerstate != "move" && playerstate != "wantmove" {
			x = kx
			y = ky
		}

		if playerstate == "choose" {
			// controltxt = "[m]ove [a]ttack [w]orld [s]tats [i]nventory [e]nd turn"
			hudtxt = "HP: " + strconv.Itoa(hp) + "/" + strconv.Itoa(maxhp) + ", Armor: " + strconv.Itoa(armor) + ", Weapon: " + weaponname + ", Status: Choosing Action"
		} else if playerstate == "move" {
			hudtxt = "HP: " + strconv.Itoa(hp) + "/" + strconv.Itoa(maxhp) + ", Armor: " + strconv.Itoa(armor) + ", Status: Moving"
			controltxt = "Steps: " + strconv.Itoa(steps) + "/6"
		} else if playerstate == "attack" {
			hudtxt = "HP: " + strconv.Itoa(hp) + "/" + strconv.Itoa(maxhp) + ", Armor: " + strconv.Itoa(armor) + ", Weapon: " + weaponname + ", Status: Attacking"
			controltxt = "Attack Here? (y/n)"
		} else if playerstate == "wantmove" {
			hudtxt = "HP: " + strconv.Itoa(hp) + "/" + strconv.Itoa(maxhp) + ", Armor: " + strconv.Itoa(armor) + ", Status: Moving"
			controltxt = "Move Here? (y/n)"
		}

		// if playerstate == "choose" && !canattack {
		// 	controltxt = "[m]ove [w]orld [s]tats [i]nventory [e]nd turn"
		// }

		// if playerstate == "choose" && steps == 6 && !canattack {
		// 	controltxt = "[w]orld [s]tats [i]nventory [e]nd turn"
		// } else if playerstate == "choose" && steps == 6 {
		// 	controltxt = "[a]ttack [w]orld [s]tats [i]nventory [e]nd turn"
		// }

		if playerstate == "choose" {
			controltxt = ""

			if steps != 6 {
				controltxt += "[m]ove "
			}

			if canattack {
				controltxt += "[a]ttack "
			}

			controltxt += "[w]orld [s]tats [i]nventory [e]nd"
		}

		if playerstate == "youcannotreach" {
			hudtxt = "You cannot reach that far!"
			controltxt = "Press any key to continue..."
			playerstate = "waitforkeypress"
			beingattacked = true
		}

		if playerstate == "noenemy" {
			hudtxt = "There is no enemy to hit!"
			controltxt = "Press any key to continue..."
			beingattacked = true
			playerstate = "waitforkeypress"
		}

		s.Clear()
		gamemap.Show(s)
		// drawText(s, 0, 0, "------")
		// drawText(s, 0, 1, "|    |")
		// drawText(s, 0, 2, "|    |")
		// drawText(s, 0, 3, "|    |")
		// drawText(s, 0, 4, "------")
		drawTextStyle(s, x, y, keeganstyle, "K")

		if playerstate == "attack" {
			drawTextStyle(s, ax, ay, aimingstyle, "+")
		}

		if playerstate == "move" || playerstate == "wantmove" {
			drawTextStyle(s, kx, ky, movestyle, "󰖃")
		}

		if enemymoving {
			drawTextStyle(s, nex, ney, aimingstyle, "󰖃")
		}

		drawText(s, 0, 16, "Distance: "+strconv.FormatFloat(getDistance(x, y, ex, ey), 'f', -1, 64))

		// Draw HUD
		drawText(s, 0, 12, hudtxt)

		// Draw Controls
		drawText(s, 0, 14, controltxt)

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
							kx = x
							ky = y
							playerstate = "move"
						} else if playerstate == "move" {
							kx = x
							ky = y
							steps = bsteps
							playerstate = "choose"
						}
					} else if ev.Rune() == 'a' {
						// attack or move in moving state
						if playerstate == "move" {
							bx = kx
							by = ky
							kx -= 1
							steps += 1
							playerstate = "moved"
						} else if playerstate != "attack" && playerstate != "wantmove" {
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
							bx = kx
							by = ky
							ky += 1
							steps += 1
							playerstate = "moved"
						} else if playerstate != "attack" && playerstate != "wantmove" {
							yourstats(s)
						}
					} else if ev.Rune() == 'd' {
						// for moving
						if playerstate == "move" {
							bx = kx
							by = ky
							kx += 1
							steps += 1
							playerstate = "moved"
						} else if playerstate == "choose" {
							if DEBUG {
								debugmenu(s)
							}
						}
					} else if ev.Rune() == 'w' {
						// for moving
						if playerstate == "move" {
							bx = kx
							by = ky
							ky -= 1
							steps += 1
							playerstate = "moved"
						}
					} else if ev.Rune() == 'n' {
						if playerstate == "attack" {
							playerstate = "choose"
						} else if playerstate == "wantmove" {
							kx = x
							ky = y
							playerstate = "choose"
							steps = bsteps
						}
					} else if ev.Rune() == 'y' {
						if playerstate == "attack" {
							good, atdis := false, false
							if ax == ex && ay == ey {
								good, atdis = WeaponDistanceGood(getDistance(x, y, ex, ey), weaponname)
							}
							println(good, atdis)
							// cantreach := false
							// if ex == 1 && ey == 1 {
							// 	if x == 1 || x == 2 {
							// 		if y == 1 || y == 2 {
							// 			playerstate = "enemy1"
							// 		} else {
							// 			if weaponname != "Pistol" {
							// 				cantreach = true
							// 			} else {
							// 				playerstate = "enemy1"
							// 			}
							// 		}
							// 	} else {
							// 		if weaponname != "Pistol" {
							// 			cantreach = true
							// 		} else {
							// 			playerstate = "enemy1"
							// 		}
							// 	}
							// }

							if good {
								playerstate = "enemy1"
							} else if !good && atdis {
								playerstate = "enemy1"
								disadvantage = true
							} else if !good && !atdis {
								playerstate = "youcannotreach"
							}
						} else if playerstate == "wantmove" || playerstate == "move" {
							ground := gamemap.GroundType(kx, ky)

							if ground == "grass" {
								keeganstyle = grassstyle
							} else {
								keeganstyle = defaultkeeganstyle
							}

							playerstate = "choose"
							bsteps = steps
						}
					} else if ev.Rune() == 'e' {
						playerstate = "idle"
					}
				}
			}

			if playerstate == "moved" {
				// Barriers Checks here

				if gamemap.CoordIsCollide(kx, ky) {
					kx = bx
					ky = by
					steps--
				}

				// if x == 0 {
				// 	x += 1
				// 	steps -= 1
				// }
				// if y == 0 {
				// 	y += 1
				// 	steps -= 1
				// }
				// if x == 5 {
				// 	x -= 1
				// 	steps -= 1
				// }
				// if y == 4 {
				// 	y -= 1
				// 	steps -= 1
				// }
				// if x == 1 && y == 1 {
				// 	if ehp <= 0 {
				// 		x = bx
				// 		y = by
				// 		steps -= 1
				// 	}
				// }

				ground := gamemap.GroundType(kx, ky)

				if ground == "grass" {
					movestyle = grassstyle
				} else {
					movestyle = defaultkeeganstyle
				}

				if steps != 6 {
					playerstate = "move"
				} else {
					playerstate = "wantmove"
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

			if enemymoving {
				ex, ey = nex, ney
				enemymoving = false

				if ex == x-1 || ex == x+1 || ex == x {
					if ey == y+1 || ey == y-1 || ey == y {
						playerstate = "idle"
					}
				}
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

	needwidth, needheight := 60, 27

	width, height := s.Size()

	term_small_err_msg := "ERROR: Your terminal is too small for this game.\nPlease change your terminal size and try again.\n\nYour Width: " + strconv.Itoa(width) + "\nYour Height: " + strconv.Itoa(height) + "\nNeeds Width: " + strconv.Itoa(needwidth) + "\nNeeds Height: " + strconv.Itoa(needheight)

	if width < needwidth || height < needheight {
		s.Fini()
		log.Fatalf(term_small_err_msg)
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
