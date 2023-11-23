package holdero

import (
	"strconv"
	"time"

	dreams "github.com/dReam-dApps/dReams"
	"github.com/dReam-dApps/dReams/bundle"
	"github.com/dReam-dApps/dReams/menu"
	"github.com/dReam-dApps/dReams/rpc"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

var tables_menu bool

// Sets bet amount and current bet readout
func ifBet(w, r uint64) {
	if w > 0 && r > 0 && !signals.placedBet {
		float := float64(w) / 100000
		wager := strconv.FormatFloat(float, 'f', 1, 64)
		table.betEntry.SetText(wager)
		round.display.results = round.raiser + " Raised, " + wager + " to Call "
	} else if w > 0 && !signals.placedBet {
		float := float64(w) / 100000
		wager := strconv.FormatFloat(float, 'f', 1, 64)
		table.betEntry.SetText(wager)
		round.display.results = round.bettor + " Bet " + wager
	} else if r > 0 && signals.placedBet {
		float := float64(r) / 100000
		raised := strconv.FormatFloat(float, 'f', 1, 64)
		table.betEntry.SetText(raised)
		round.display.results = round.raiser + " Raised, " + raised + " to Call"
	} else if w == 0 && !signals.bet {
		var float float64
		if round.Ante == 0 {
			float = float64(round.BB) / 100000
		} else {
			float = float64(round.Ante) / 100000
		}
		this := strconv.FormatFloat(float, 'f', 1, 64)
		table.betEntry.SetText(this)
		if !signals.reveal {
			round.display.results = "Check or Bet"
			table.betEntry.Enable()
		}
	} else if !signals.deal {
		round.display.results = "Deal Hand"
	}

	table.betEntry.Refresh()
}

// Single shot triggering ifBet() on players turn
func singleShot(turn, trigger bool) bool {
	if turn && !trigger {
		ifBet(round.Wager, round.Raised)
		return true
	}

	if !turn {
		return false
	} else {
		return turn
	}
}

// Main Holdero process
func fetch(d *dreams.AppObject) {
	initValues()
	time.Sleep(3 * time.Second)
	var autoCF, autoD, autoB, trigger bool
	var skip, delay, offset int
	for {
		select {
		case <-d.Receive():
			if !rpc.Wallet.IsConnected() || !rpc.Daemon.IsConnected() {
				signals.contract = false
				disableActions()
				Settings.synced = false
				setHolderoLabel()
				d.WorkDone()
				continue
			}

			if !Settings.synced && menu.GnomonScan(d.IsConfiguring()) {
				logger.Println("[Holdero] Syncing")
				createTableList()
				Settings.synced = true
				H.Actions.Show()
			}

			if signals.contract {
				Settings.check.SetChecked(true)
			} else {
				Settings.check.SetChecked(false)
				disableOwnerControls(true)
				signals.sit = true
			}

			FetchHolderoSC()

			table.stats.Container = *container.NewVBox(
				container.NewStack(tableIcon(bundle.ResourceAvatarFramePng)),
				table.stats.name,
				table.stats.desc,
				table.stats.owner,
				table.stats.chips,
				table.stats.blinds,
				table.stats.version,
				table.stats.last,
				table.stats.seats)
			table.stats.Container.Refresh()

			if (round.Turn == round.ID && rpc.Wallet.Height > signals.height+4) ||
				(round.Turn != round.ID && round.ID >= 1) || (!signals.myTurn && round.ID >= 1) {
				if signals.clicked {
					trigger = false
					autoCF = false
					autoD = false
					autoB = false
					signals.reveal = false
				}
				signals.clicked = false
			}

			if !signals.clicked {
				if round.first {
					round.first = false
					delay = 0
					round.delay = false
				}

				if round.delay {
					now := time.Now().Unix()
					delay++
					if delay >= 17 || now > round.Last+60 {
						delay = 0
						round.delay = false
					}
				} else {
					setHolderoLabel()
					GetUrls(round.cards.Faces.Url, round.cards.Backs.Url)
					Called(round.flop, round.Wager)
					trigger = singleShot(signals.myTurn, trigger)
					holderoRefresh(d, offset)
					// Auto check
					if Settings.auto.check && signals.myTurn && !autoCF {
						if !signals.reveal && !signals.end && !round.localEnd {
							if round.cards.Local1 != "" {
								ActionBuffer()
								Check()
								H.TopLabel.Text = "Auto Check/Fold Tx Sent"
								H.TopLabel.Refresh()
								autoCF = true

								go func() {
									if !d.IsWindows() {
										time.Sleep(500 * time.Millisecond)
										round.notified = d.Notification("dReams - Holdero", "Auto Check/Fold TX Sent")
									}
								}()
							}
						}
					}

					// Auto deal
					if Settings.auto.deal && signals.myTurn && !autoD && GameIsActive() {
						if !signals.reveal && !signals.end && !round.localEnd {
							if round.cards.Local1 == "" {
								autoD = true
								go func() {
									time.Sleep(2100 * time.Millisecond)
									ActionBuffer()
									DealHand()
									H.TopLabel.Text = "Auto Deal Tx Sent"
									H.TopLabel.Refresh()

									if !d.IsWindows() {
										time.Sleep(300 * time.Millisecond)
										round.notified = d.Notification("dReams - Holdero", "Auto Deal TX Sent")
									}
								}()
							}
						}
					}

					// Auto bet
					if Odds.IsRunning() && signals.myTurn && !autoB && GameIsActive() {
						if !signals.reveal && !signals.end && !round.localEnd {
							if round.cards.Local1 != "" {
								autoB = true
								go func() {
									time.Sleep(2100 * time.Millisecond)
									ActionBuffer()
									odds, future := MakeOdds()
									BetLogic(odds, future, true)
									H.TopLabel.Text = "Auto Bet Tx Sent"
									H.TopLabel.Refresh()

									if !d.IsWindows() {
										time.Sleep(300 * time.Millisecond)
										round.notified = d.Notification("dReams - Holdero", "Auto Bet TX Sent")
									}
								}()
							}
						}
					}

					if round.ID > 1 && signals.myTurn && !signals.end && !round.localEnd {
						now := time.Now().Unix()
						if now > round.Last+100 {
							table.warning.Show()
						} else {
							table.warning.Hide()
						}
					} else {
						table.warning.Hide()
					}

					skip = 0
				}
			} else {
				waitLabel()
				revealingKey(d)
				skip++
				if skip >= 25 {
					signals.clicked = false
					skip = 0
					trigger = false
					autoCF = false
					autoD = false
					autoB = false
					signals.reveal = false
				}
			}

			offset++
			if offset >= 21 {
				offset = 0
			}

			d.WorkDone()
		case <-d.CloseDapp():
			logger.Println("[Holdero] Done")
			return
		}
	}
}

// Do when disconnected
func Disconnected(b bool) {
	if b {
		round.ID = 0
		round.display.playerId = ""
		Odds.Stop()
		Settings.faces.Select.Options = []string{"Light", "Dark"}
		Settings.backs.Select.Options = []string{"Light", "Dark"}
		Settings.avatars.Select.Options = []string{"None"}
		Settings.faces.URL = ""
		Settings.backs.URL = ""
		Settings.avatar.url = ""
		Settings.faces.Select.SetSelectedIndex(0)
		Settings.backs.Select.SetSelectedIndex(0)
		Settings.avatars.Select.SetSelectedIndex(0)
		Settings.faces.Select.Refresh()
		Settings.backs.Select.Refresh()
		Settings.avatars.Select.Refresh()
		DisableHolderoTools()
		Settings.synced = false
		table.owner.valid = false
		table.Public.List.UnselectAll()
	}
}

func disableActions() {
	H.Actions.Hide()
	H.DApp.Objects[4].(*fyne.Container).Objects[0].(*fyne.Container).Objects[1].(*fyne.Container).Objects[1].(*widget.Check).SetChecked(false)
	H.DApp.Objects[4].(*fyne.Container).Objects[0].(*fyne.Container).Objects[1].(*fyne.Container).Objects[2].(*widget.Check).SetChecked(false)
	Settings.check.SetChecked(false)
	table.entry.SetText("")
	clearShared()
	disableOwnerControls(true)
	table.Public.SCIDs = []string{}
	table.Owned.SCIDs = []string{}
	table.unlock.Hide()
	table.new.Hide()
	table.tournament.Hide()
	table.unlock.Refresh()
	table.new.Refresh()
	table.tournament.Refresh()
}

// Disable Holdero owner actions
func disableOwnerControls(d bool) {
	if d {
		table.owner.owners_left.Hide()
		table.owner.owners_mid.Hide()
	} else {
		table.owner.owners_left.Show()
		table.owner.owners_mid.Show()
	}

	table.owner.owners_left.Refresh()
	table.owner.owners_mid.Refresh()
}

// Sets Holdero table info labels
func setHolderoLabel() {
	H.TopLabel.Text = round.display.results
	H.LeftLabel.SetText("Seats: " + round.display.seats + "      Pot: " + round.display.pot + "      Blinds: " + round.display.blinds + "      Ante: " + round.display.ante + "      Dealer: " + round.display.dealer)
	if round.asset {
		if round.tourney {
			H.RightLabel.SetText(round.display.readout + "      Player ID: " + round.display.playerId + "      Chip Balance: " + rpc.DisplayBalance("Tournament") + "      Dero Balance: " + rpc.DisplayBalance("Dero") + "      Height: " + rpc.Wallet.Display.Height)
		} else {
			asset_name := rpc.GetAssetSCIDName(round.assetID)
			H.RightLabel.SetText(round.display.readout + "      Player ID: " + round.display.playerId + "      " + asset_name + " Balance: " + rpc.DisplayBalance(asset_name) + "      Dero Balance: " + rpc.DisplayBalance("Dero") + "      Height: " + rpc.Wallet.Display.Height)
		}
	} else {
		H.RightLabel.SetText(round.display.readout + "      Player ID: " + round.display.playerId + "      Dero Balance: " + rpc.DisplayBalance("Dero") + "      Height: " + rpc.Wallet.Display.Height)
	}

	if signals.contract {
		Settings.shared.Enable()
	} else {
		Settings.shared.Disable()
	}

	H.TopLabel.Refresh()
	H.LeftLabel.Refresh()
	H.RightLabel.Refresh()
}

// Holdero label for waiting for block
func waitLabel() {
	H.TopLabel.Text = ""
	if round.asset {
		if round.tourney {
			H.RightLabel.SetText("Wait for Block" + "      Player ID: " + round.display.playerId + "      Chip Balance: " + rpc.DisplayBalance("Tournament") + "      Dero Balance: " + rpc.DisplayBalance("Dero") + "      Height: " + rpc.Wallet.Display.Height)
		} else {
			asset_name := rpc.GetAssetSCIDName(round.assetID)
			H.RightLabel.SetText("Wait for Block" + "      Player ID: " + round.display.playerId + "      " + asset_name + " Balance: " + rpc.DisplayBalance(asset_name) + "      Dero Balance: " + rpc.DisplayBalance("Dero") + "      Height: " + rpc.Wallet.Display.Height)
		}

	} else {
		H.RightLabel.SetText("Wait for Block" + "      Player ID: " + round.display.playerId + "      Dero Balance: " + rpc.DisplayBalance("Dero") + "      Height: " + rpc.Wallet.Display.Height)
	}
	H.TopLabel.Refresh()
	H.RightLabel.Refresh()
}

// Refresh all Holdero gui objects
func holderoRefresh(d *dreams.AppObject, offset int) {
	ShowAvatar(d.OnTab("Holdero"))
	refreshHolderoCards(round.cards.Local1, round.cards.Local2, d)
	refreshHolderoPlayers()
	if !signals.clicked {
		if round.ID == 0 && rpc.Wallet.IsConnected() {
			if signals.sit {
				table.sit.Hide()
			} else {
				table.sit.Show()
			}
			table.leave.Hide()
			table.deal.Hide()
			table.check.Hide()
			table.bet.Hide()
			table.betEntry.Hide()
		} else if !signals.end && !signals.reveal && signals.myTurn && rpc.Wallet.IsConnected() {
			if signals.sit {
				table.sit.Hide()
			} else {
				table.sit.Show()
			}

			if signals.leave {
				table.leave.Hide()
			} else {
				table.leave.Show()
			}

			if signals.deal {
				table.deal.Hide()
			} else {
				table.deal.Show()
			}

			table.check.SetText(round.display.checkButton)
			table.bet.SetText(round.display.betButton)
			if signals.bet {
				table.check.Hide()
				table.bet.Hide()
				table.betEntry.Hide()
			} else {
				table.check.Show()
				table.bet.Show()
				table.betEntry.Show()
			}

			if !round.notified {
				if !d.IsWindows() {
					round.notified = d.Notification("dReams - Holdero", "Your Turn")
				}
			}
		} else {
			if signals.sit {
				table.sit.Hide()
			} else if !signals.sit && rpc.Wallet.IsConnected() {
				table.sit.Show()
			}
			table.leave.Hide()
			table.deal.Hide()
			table.check.Hide()
			table.bet.Hide()
			table.betEntry.Hide()

			if !signals.myTurn && !signals.end && !round.localEnd {
				round.display.results = ""
				round.notified = false
			}
		}
	}

	if tables_menu {
		if offset%3 == 0 {
			go getTableStats(round.Contract, false)
		}
	}
}

// Refresh Holdero player names and avatars
func refreshHolderoPlayers() {
	H.Back.Objects[0] = HolderoTable(ResourcePokerTablePng)
	// H.Back.Objects[0].Refresh()

	H.Back.Objects[1] = Player1_label(ResourceUnknownAvatarPng, bundle.ResourceAvatarFramePng, ResourceTurnFramePng)
	// H.Back.Objects[1].Refresh()

	H.Back.Objects[2] = Player2_label(ResourceUnknownAvatarPng, bundle.ResourceAvatarFramePng, ResourceTurnFramePng)
	// H.Back.Objects[2].Refresh()

	H.Back.Objects[3] = Player3_label(ResourceUnknownAvatarPng, bundle.ResourceAvatarFramePng, ResourceTurnFramePng)
	// H.Back.Objects[3].Refresh()

	H.Back.Objects[4] = Player4_label(ResourceUnknownAvatarPng, bundle.ResourceAvatarFramePng, ResourceTurnFramePng)
	// H.Back.Objects[4].Refresh()

	H.Back.Objects[5] = Player5_label(ResourceUnknownAvatarPng, bundle.ResourceAvatarFramePng, ResourceTurnFramePng)
	// H.Back.Objects[5].Refresh()

	H.Back.Objects[6] = Player6_label(ResourceUnknownAvatarPng, bundle.ResourceAvatarFramePng, ResourceTurnFramePng)
	// H.Back.Objects[6].Refresh()
}

// Reveal key notification and display
func revealingKey(d *dreams.AppObject) {
	if signals.reveal && signals.myTurn && !signals.end {
		if !round.notified {
			round.display.results = "Revealing Key"
			H.TopLabel.Text = round.display.results
			H.TopLabel.Refresh()

			if !d.IsWindows() {
				round.notified = d.Notification("dReams - Holdero", "Revealing Key")
			}
		}
	}
}
