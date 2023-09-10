package holdero

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/dReam-dApps/dReams/rpc"
)

type playerId struct {
	Id     string `json:"id"`
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
}
type CardSpecs struct {
	Faces struct {
		Name string `json:"Name"`
		Url  string `json:"Url"`
	} `json:"Faces"`
	Backs struct {
		Name string `json:"Name"`
		Url  string `json:"Url"`
	} `json:"Backs"`
}

type TableSpecs struct {
	MaxBet float64 `json:"Maxbet"`
	MinBuy float64 `json:"Minbuy"`
	MaxBuy float64 `json:"Maxbuy"`
	Time   int     `json:"Time"`
}

type tableSignals struct {
	contract  bool
	deal      bool
	bet       bool
	called    bool
	reveal    bool
	end       bool
	sit       bool
	leave     bool
	In1       bool
	In2       bool
	In3       bool
	In4       bool
	In5       bool
	In6       bool
	Out1      bool
	myTurn    bool
	placedBet bool
	paid      bool
	log       bool
	odds      bool
	clicked   bool
	height    int
	times     struct {
		kick  int
		delay int
		block int
	}
}

var signals tableSignals

// Make blinds display string
func blindString(b, s interface{}) string {
	if bb, ok := b.(float64); ok {
		if sb, ok := s.(float64); ok {
			return fmt.Sprintf("%.5f / %.5f", bb/100000, sb/100000)
		}
	}

	return "? / ?"
}

// If Holdero table is closed set vars accordingly
func closedTable() {
	round.winner = ""
	round.Players = 0
	round.Pot = 0
	round.ID = 0
	round.tourney = false
	round.p1.url = ""
	round.p2.url = ""
	round.p3.url = ""
	round.p4.url = ""
	round.p5.name = ""
	round.p6.url = ""
	round.p1.name = ""
	round.p2.name = ""
	round.p3.name = ""
	round.p4.name = ""
	round.p5.name = ""
	round.p6.name = ""
	round.bettor = ""
	round.raiser = ""
	round.Turn = 0
	round.Last = 0
	round.trigger.local = false
	round.trigger.flop = false
	round.trigger.turn = false
	round.trigger.river = false
	signals.Out1 = false
	signals.sit = true
	signals.In1 = false
	signals.In2 = false
	signals.In3 = false
	signals.In4 = false
	signals.In5 = false
	signals.In6 = false
	round.display.seats = ""
	round.display.pot = ""
	round.display.blinds = ""
	round.display.ante = ""
	round.display.dealer = ""
	round.display.playerId = ""
}

// Clear a single players name and avatar values
func singleNameClear(p int) {
	switch p {
	case 1:
		round.p1.name = ""
		round.p1.url = ""
	case 2:
		round.p2.name = ""
		round.p2.url = ""
	case 3:
		round.p3.name = ""
		round.p3.url = ""
	case 4:
		round.p4.name = ""
		round.p4.url = ""
	case 5:
		round.p5.name = ""
		round.p5.name = ""
	case 6:
		round.p6.name = ""
		round.p6.url = ""
	default:

	}
}

// Returns name of Holdero player who bet
func findBettor(p interface{}) string {
	if p != nil {
		switch rpc.Float64Type(p) {
		case 0:
			if round.p6.name != "" && !round.p6.folded {
				return round.p6.name
			} else if round.p5.name != "" && !round.p5.folded {
				return round.p5.name
			} else if round.p4.name != "" && !round.p4.folded {
				return round.p4.name
			} else if round.p3.name != "" && !round.p3.folded {
				return round.p3.name
			} else if round.p2.name != "" && !round.p2.folded {
				return round.p2.name
			}
		case 1:
			if round.p1.name != "" && !round.p1.folded {
				return round.p1.name
			} else if round.p6.name != "" && !round.p6.folded {
				return round.p6.name
			} else if round.p5.name != "" && !round.p5.folded {
				return round.p5.name
			} else if round.p4.name != "" && !round.p4.folded {
				return round.p4.name
			} else if round.p3.name != "" && !round.p3.folded {
				return round.p3.name
			}
		case 2:
			if round.p2.name != "" && !round.p2.folded {
				return round.p2.name
			} else if round.p1.name != "" && !round.p1.folded {
				return round.p1.name
			} else if round.p6.name != "" && !round.p6.folded {
				return round.p6.name
			} else if round.p5.name != "" && !round.p5.folded {
				return round.p5.name
			} else if round.p4.name != "" && !round.p4.folded {
				return round.p4.name
			}
		case 3:
			if round.p3.name != "" && !round.p3.folded {
				return round.p3.name
			} else if round.p2.name != "" && !round.p2.folded {
				return round.p2.name
			} else if round.p1.name != "" && !round.p1.folded {
				return round.p1.name
			} else if round.p6.name != "" && !round.p6.folded {
				return round.p6.name
			} else if round.p5.name != "" && !round.p5.folded {
				return round.p5.name
			}
		case 4:
			if round.p4.name != "" && !round.p4.folded {
				return round.p4.name
			} else if round.p3.name != "" && !round.p3.folded {
				return round.p3.name
			} else if round.p2.name != "" && !round.p2.folded {
				return round.p2.name
			} else if round.p1.name != "" && !round.p1.folded {
				return round.p1.name
			} else if round.p6.name != "" && !round.p6.folded {
				return round.p6.name
			}
		case 5:
			if round.p5.name != "" && !round.p5.folded {
				return round.p5.name
			} else if round.p4.name != "" && !round.p4.folded {
				return round.p4.name
			} else if round.p3.name != "" && !round.p3.folded {
				return round.p3.name
			} else if round.p2.name != "" && !round.p2.folded {
				return round.p2.name
			} else if round.p1.name != "" && !round.p1.folded {
				return round.p1.name
			}
		default:
			return ""
		}
	}

	return ""
}

// Gets Holdero player name and avatar, returns player Id string
func getAvatar(p int, id interface{}) string {
	if id == nil {
		singleNameClear(p)
		return "nil"
	}

	str := fmt.Sprint(id)

	if len(str) == 64 {
		return str
	}

	av := rpc.HexToString(str)

	var player playerId

	if err := json.Unmarshal([]byte(av), &player); err != nil {
		logger.Errorln("[getAvatar]", err)
		return ""
	}

	switch p {
	case 1:
		round.p1.name = player.Name
		round.p1.url = player.Avatar
	case 2:
		round.p2.name = player.Name
		round.p2.url = player.Avatar
	case 3:
		round.p3.name = player.Name
		round.p3.url = player.Avatar
	case 4:
		round.p4.name = player.Name
		round.p4.url = player.Avatar
	case 5:
		round.p5.name = player.Name
		round.p5.url = player.Avatar
	case 6:
		round.p6.name = player.Name
		round.p6.url = player.Avatar
	}

	return player.Id
}

// Check if player Id matches rpc.Wallet.IdHash
func checkPlayerId(one, two, three, four, five, six string) string {
	var id string
	if rpc.Wallet.IdHash == one {
		id = "1"
		round.ID = 1
	} else if rpc.Wallet.IdHash == two {
		id = "2"
		round.ID = 2
	} else if rpc.Wallet.IdHash == three {
		id = "3"
		round.ID = 3
	} else if rpc.Wallet.IdHash == four {
		id = "4"
		round.ID = 4
	} else if rpc.Wallet.IdHash == five {
		id = "5"
		round.ID = 5
	} else if rpc.Wallet.IdHash == six {
		id = "6"
		round.ID = 6
	} else {
		id = ""
		round.ID = 0
	}

	return id
}

// Set Holdero name signals for when player is at table
func setHolderoName(one, two, three, four, five, six interface{}) {
	if one != nil {
		signals.In1 = true
	} else {
		signals.In1 = false
	}

	if two != nil {
		signals.In2 = true
	} else {
		signals.In2 = false
	}

	if three != nil {
		signals.In3 = true
	} else {
		signals.In3 = false
	}

	if four != nil {
		signals.In4 = true
	} else {
		signals.In4 = false
	}

	if five != nil {
		signals.In5 = true
	} else {
		signals.In5 = false
	}

	if six != nil {
		signals.In6 = true
	} else {
		signals.In6 = false
	}
}

// When Holdero pot is empty set vars accordingly
func potIsEmpty(pot uint64) {
	if pot == 0 {
		if !signals.myTurn {
			rpc.Wallet.KeyLock = false
		}
		round.winningHand = []int{}
		round.cards.flop1 = 0
		round.cards.flop2 = 0
		round.cards.flop3 = 0
		round.cards.turn = 0
		round.cards.river = 0
		round.localEnd = false
		round.Wager = 0
		round.Raised = 0
		round.winner = ""
		round.printed = false
		round.cards.Local1 = ""
		round.cards.Local2 = ""
		round.cards.P1C1 = ""
		round.cards.P1C2 = ""
		round.cards.P2C1 = ""
		round.cards.P2C2 = ""
		round.cards.P3C1 = ""
		round.cards.P3C2 = ""
		round.cards.P4C1 = ""
		round.cards.P4C2 = ""
		round.cards.P5C1 = ""
		round.cards.P5C2 = ""
		round.cards.P6C1 = ""
		round.cards.P6C2 = ""
		round.cards.Key1 = ""
		round.cards.Key2 = ""
		round.cards.Key3 = ""
		round.cards.Key4 = ""
		round.cards.Key5 = ""
		round.cards.Key6 = ""
		signals.called = false
		signals.placedBet = false
		signals.reveal = false
		signals.end = false
		signals.paid = false
		signals.log = false
		signals.odds = false
		round.display.results = ""
		round.bettor = ""
		round.raiser = ""
		round.trigger.local = false
		round.trigger.flop = false
		round.trigger.turn = false
		round.trigger.river = false
	}
}

// Sets Holdero sit signal if table has open seats
func tableOpen(seats, full, two, three, four, five, six interface{}) {
	players := 1
	if two != nil {
		players++
	}

	if three != nil {
		players++
	}

	if four != nil {
		players++
	}

	if five != nil {
		players++
	}

	if six != nil {
		players++
	}

	if signals.Out1 {
		players--
	}

	round.Players = players

	if round.ID > 1 {
		signals.sit = true
		return
	}
	s := rpc.IntType(seats)
	if s >= 2 && two == nil && round.ID != 1 {
		signals.sit = false
	}

	if s >= 3 && three == nil && round.ID != 1 {
		signals.sit = false
	}

	if s >= 4 && four == nil && round.ID != 1 {
		signals.sit = false
	}

	if s >= 5 && five == nil && round.ID != 1 {
		signals.sit = false
	}

	if s == 6 && six == nil && round.ID != 1 {
		signals.sit = false
	}

	if full != nil {
		signals.sit = true
	}
}

// Gets Holdero community card values
func getCommCardValues(f1, f2, f3, t, r interface{}) {
	if f1 != nil {
		round.cards.flop1 = rpc.IntType(f1)
		if !round.trigger.flop {
			round.delay = true
		}
		round.trigger.flop = true
	} else {
		round.cards.flop1 = 0
		round.trigger.flop = false
	}

	if f2 != nil {
		round.cards.flop2 = rpc.IntType(f2)
	} else {
		round.cards.flop2 = 0
	}

	if f3 != nil {
		round.cards.flop3 = rpc.IntType(f3)
	} else {
		round.cards.flop3 = 0
	}

	if t != nil {
		round.cards.turn = rpc.IntType(t)
		if !round.trigger.turn {
			round.delay = true
		}
		round.trigger.turn = true
	} else {
		round.cards.turn = 0
		round.trigger.turn = false
	}

	if r != nil {
		round.cards.river = rpc.IntType(r)
		if !round.trigger.river {
			round.delay = true
		}
		round.trigger.river = true
	} else {
		round.cards.river = 0
		round.trigger.river = false
	}
}

// Gets Holdero player card hash values
func getPlayerCardValues(a1, a2, b1, b2, c1, c2, d1, d2, e1, e2, f1, f2 interface{}) {
	if round.ID == 1 {
		if a1 != nil {
			round.cards.Local1 = fmt.Sprint(a1)
			round.cards.Local2 = fmt.Sprint(a2)
			if !round.trigger.local {
				round.delay = true
			}
			round.trigger.local = true
		} else {
			round.cards.Local1 = ""
			round.cards.Local2 = ""
			round.trigger.local = false
		}
	}

	if a1 != nil {
		round.cards.P1C1 = fmt.Sprint(a1)
		round.cards.P1C2 = fmt.Sprint(a2)
	} else {
		round.cards.P1C1 = ""
		round.cards.P1C2 = ""
	}

	if round.ID == 2 {
		if b1 != nil {
			round.cards.Local1 = fmt.Sprint(b1)
			round.cards.Local2 = fmt.Sprint(b2)
			if !round.trigger.local {
				round.delay = true
			}
			round.trigger.local = true
		} else {
			round.cards.Local1 = ""
			round.cards.Local2 = ""
			round.trigger.local = false
		}
	}

	if b1 != nil {
		round.cards.P2C1 = fmt.Sprint(b1)
		round.cards.P2C2 = fmt.Sprint(b2)
	} else {
		round.cards.P2C1 = ""
		round.cards.P2C2 = ""
	}

	if round.ID == 3 {
		if c1 != nil {
			round.cards.Local1 = fmt.Sprint(c1)
			round.cards.Local2 = fmt.Sprint(c2)
			if !round.trigger.local {
				round.delay = true
			}
			round.trigger.local = true
		} else {
			round.cards.Local1 = ""
			round.cards.Local2 = ""
			round.trigger.local = false
		}
	}

	if c1 != nil {
		round.cards.P3C1 = fmt.Sprint(c1)
		round.cards.P3C2 = fmt.Sprint(c2)
	} else {
		round.cards.P3C1 = ""
		round.cards.P3C2 = ""
	}

	if round.ID == 4 {
		if d1 != nil {
			round.cards.Local1 = fmt.Sprint(d1)
			round.cards.Local2 = fmt.Sprint(d2)
			if !round.trigger.local {
				round.delay = true
			}
			round.trigger.local = true
		} else {
			round.cards.Local1 = ""
			round.cards.Local2 = ""
			round.trigger.local = false
		}
	}

	if d1 != nil {
		round.cards.P4C1 = fmt.Sprint(d1)
		round.cards.P4C2 = fmt.Sprint(d2)
	} else {
		round.cards.P4C1 = ""
		round.cards.P4C2 = ""
	}

	if round.ID == 5 {
		if e1 != nil {
			round.cards.Local1 = fmt.Sprint(e1)
			round.cards.Local2 = fmt.Sprint(e2)
			if !round.trigger.local {
				round.delay = true
			}
			round.trigger.local = true
		} else {
			round.cards.Local1 = ""
			round.cards.Local2 = ""
			round.trigger.local = false
		}
	}

	if e1 != nil {
		round.cards.P5C1 = fmt.Sprint(e1)
		round.cards.P5C2 = fmt.Sprint(e2)
	} else {
		round.cards.P5C1 = ""
		round.cards.P5C2 = ""
	}

	if round.ID == 6 {
		if f1 != nil {
			round.cards.Local1 = fmt.Sprint(f1)
			round.cards.Local2 = fmt.Sprint(f2)
			if !round.trigger.local {
				round.delay = true
			}
			round.trigger.local = true
		} else {
			round.cards.Local1 = ""
			round.cards.Local2 = ""
			round.trigger.local = false
		}
	}

	if f1 != nil {
		round.cards.P6C1 = fmt.Sprint(f1)
		round.cards.P6C2 = fmt.Sprint(f2)
	} else {
		round.cards.P6C1 = ""
		round.cards.P6C2 = ""
	}

	if round.ID == 0 {
		round.cards.Local1 = ""
		round.cards.Local2 = ""
	}
}

// If Holdero player has called set Signal.Called, and reset Signal.PlacedBet when no wager
func Called(fb bool, w uint64) {
	if w == 0 {
		if fb {
			signals.called = true
		} else {
			signals.called = false
		}

		if signals.called {
			round.Raised = 0
			round.Wager = 0
			signals.placedBet = false
			signals.called = false
		}

		round.display.betButton = "Bet"
		round.display.checkButton = "Check"
	}
}

// Holdero players turn display string
func turnReadout(t interface{}) (turn string) {
	if t != nil {
		switch rpc.AddOne(t) {
		case round.display.playerId:
			turn = "Your Turn"
		case "1":
			turn = "Player 1's Turn"
		case "2":
			turn = "Player 2's Turn"
		case "3":
			turn = "Player 3's Turn"
		case "4":
			turn = "Player 4's Turn"
		case "5":
			turn = "Player 5's Turn"
		case "6":
			turn = "Player 6's Turn"
		}
	}

	return
}

// Sets Holdero action signals
func setSignals(pot uint64, one interface{}) {
	if !round.localEnd {
		if len(round.cards.Local1) != 64 {
			signals.deal = false
			signals.leave = false
			signals.bet = true
		} else {
			signals.deal = true
			signals.leave = true
			if pot != 0 {
				signals.bet = false
			} else {
				signals.bet = true
			}
		}
	} else {
		signals.deal = true
		signals.leave = true
		signals.bet = true
	}

	if round.ID > 1 {
		signals.sit = true
	}

	if round.ID == 1 {
		if one != nil {
			signals.sit = false
		} else {
			signals.sit = true
		}
	}
}

// If Holdero player has folded, set Round folded bools and clear cards
func hasFolded(one, two, three, four, five, six interface{}) {
	if one != nil {
		round.p1.folded = true
		round.cards.P1C1 = ""
		round.cards.P1C2 = ""
	} else {
		round.p1.folded = false
	}

	if two != nil {
		round.p2.folded = true
		round.cards.P2C1 = ""
		round.cards.P2C2 = ""
	} else {
		round.p2.folded = false
	}

	if three != nil {
		round.p3.folded = true
		round.cards.P3C1 = ""
		round.cards.P3C2 = ""
	} else {
		round.p3.folded = false
	}

	if four != nil {
		round.p4.folded = true
		round.cards.P4C1 = ""
		round.cards.P4C2 = ""
	} else {
		round.p4.folded = false
	}

	if five != nil {
		round.p5.folded = true
		round.cards.P5C1 = ""
		round.cards.P5C2 = ""
	} else {
		round.p5.folded = false
	}

	if six != nil {
		round.p6.folded = true
		round.cards.P6C1 = ""
		round.cards.P6C2 = ""
	} else {
		round.p6.folded = false
	}
}

// Determine if all players have folded and trigger payout
func allFolded(p1, p2, p3, p4, p5, p6, s interface{}) {
	var a, b, c, d, e, f int
	var who string
	var display string
	seats := rpc.IntType(s)
	if seats >= 2 {
		if p1 != nil {
			a = rpc.IntType(p1)
		} else {
			who = "Player1"
			display = round.p1.name
		}
		if p2 != nil {
			b = rpc.IntType(p2)
		} else {
			who = "Player2"
			display = round.p2.name
		}
	}
	if seats >= 3 {
		if p3 != nil {
			c = rpc.IntType(p3)
		} else {
			who = "Player3"
			display = round.p3.name
		}
	}

	if seats >= 4 {
		if p4 != nil {
			d = rpc.IntType(p4)
		} else {
			who = "Player4"
			display = round.p4.name
		}
	}

	if seats >= 5 {
		if p5 != nil {
			e = rpc.IntType(p5)
		} else {
			who = "Player5"
			display = round.p5.name
		}
	}

	if seats >= 6 {
		if p6 != nil {
			f = rpc.IntType(p6)
		} else {
			who = "Player6"
			display = round.p6.name
		}
	}

	i := a + b + c + d + e + f

	if 1+i-seats == 0 {
		round.localEnd = true
		round.winner = who
		round.display.results = display + " Wins, All Players Have Folded"
		if GameIsActive() && round.Pot > 0 {
			if !signals.log {
				signals.log = true
				rpc.AddLog(round.display.results)
			}

			updateStatsWins(round.Pot, who, true)
		}
	}
}

// Payout routine when all Holdero players have folded
func allFoldedWinner() {
	if round.ID == 1 {
		if round.localEnd && !rpc.Startup {
			if !signals.paid {
				signals.paid = true
				go func() {
					time.Sleep(time.Duration(signals.times.delay) * time.Second)
					retry := 0
					for retry < 4 {
						tx := PayOut(round.winner)
						time.Sleep(time.Second)
						retry += rpc.ConfirmTxRetry(tx, "Holdero", 36)
					}
				}()
			}
		}
	}
}

// If Holdero showdown, trigger the hand ranker routine
func winningHand(e interface{}) {
	if e != nil && !rpc.Startup && !round.localEnd {
		go func() {
			getHands(rpc.StringToInt(round.display.seats))
		}()
	}
}
