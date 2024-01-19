package holdero

import (
	"image/color"
	"time"

	dreams "github.com/dReam-dApps/dReams"
	"github.com/dReam-dApps/dReams/bundle"
	"github.com/dReam-dApps/dReams/dwidget"
	"github.com/dReam-dApps/dReams/menu"
	"github.com/dReam-dApps/dReams/rpc"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

var H dreams.ContainerStack

// Holdero tables menu tab layout
func placeContract(change_screen *fyne.Container, d *dreams.AppObject) *fyne.Container {
	Settings.check = widget.NewCheck("", func(b bool) {
		if !b {
			disableOwnerControls(true)
		}
	})
	Settings.check.Disable()

	tabs := container.NewAppTabs(
		container.NewTabItemWithIcon("", ResourceHolderoCirclePng, layout.NewSpacer()),
		container.NewTabItem("Tables", publicList(d)),
		container.NewTabItem("Favorites", favoritesList()),
		container.NewTabItem("Owned", ownedList(d)),
		container.NewTabItem("View Table", layout.NewSpacer()),
		container.NewTabItem("How to Play", layout.NewSpacer()))

	tabs.DisableIndex(0)
	tabs.SelectIndex(1)

	tabs.OnSelected = func(ti *container.TabItem) {
		switch ti.Text {
		case "Tables":
			if rpc.Daemon.IsConnected() {
				go createTableList(nil)
			}
		case "How to Play":
			instructions := "Connect to your wallet and daemon and wait for tables to sync\n\nClick on a table in the list to connect to it\n\nClick on 'View Table' to view it\n\nIf there is a open seat you can click 'Sit Down' to join the game\n\nWhen it is your turn you can click 'Deal Hand' to get your cards\n\nHoldero is a no limit single raise version of Hold'em\n\nThere is no all in, players must call the bet or fold\n\nAt the start of each deal players can leave the table\n\nIf you are inactive during the hand you will be timed out and removed from the table\n\nAssets that unlock dReam Tools give access to bot players and odds calculators\n\nYou can create and view your tables in the 'Owned' tab\n\nVisit dreamdapps.io for more docs"
			dialog.NewInformation("How to play", instructions, d.Window).Show()
			tabs.SelectIndex(1)

		default:
		}

		if ti.Text == "View Table" {
			go func() {
				if len(round.Contract) == 64 {
					fetchHolderoSC()
					if round.ID == 0 {
						if menu.Assets.Names.SelectedIndex() == 0 && len(menu.Assets.Names.Options) > 1 {
							dialog.NewInformation("Anon", "You are playing as a anon player, select name in assets/profile tab", d.Window).Show()
						}
					}
					d.Window.Content().(*fyne.Container).Objects[1].(*fyne.Container).Objects[1].(*container.AppTabs).Selected().Content = change_screen
					tabs.SelectIndex(1)
					now := time.Now().Unix()
					if now > round.Last+33 {
						holderoRefresh(d, 0)
					}
				} else {
					if rpc.IsReady() {
						dialog.NewInformation("Select Table", "Select a table SCID to play at", d.Window).Show()
					} else {
						dialog.NewInformation("Not Connected", "Connect to daemon and wallet", d.Window).Show()
					}
					tabs.SelectIndex(1)
				}
			}()
		}
	}

	// Holdero SCID entry bound to round.Contract
	// entry text set on all List selections
	table.entry = widget.NewSelectEntry(nil)
	options := []string{""}
	table.entry.SetOptions(options)
	table.entry.PlaceHolder = "Holdero Contract Address: "

	this := binding.BindString(&round.Contract)
	table.entry.Bind(this)

	tabs.SetTabLocation(container.TabLocationLeading)

	contract_cont := container.NewBorder(nil, nil, nil, Settings.check, table.entry)

	table.unlock.OnTapped = func() {
		holderoMenuConfirm(1, d)
	}

	table.new.OnTapped = func() {
		holderoMenuConfirm(2, d)
	}

	// Changes to SCID entry clear table and check if current entry is valid table
	var wait bool
	table.entry.OnCursorChanged = func() {
		if rpc.Daemon.IsConnected() && !wait {
			wait = true
			text := table.entry.Text
			go clearShared()
			if len(text) == 64 {
				if checkTableOwner(text) {
					disableOwnerControls(false)
					if checkTableVersion(text) >= 110 {
						table.owner.chips.Show()
						table.owner.times.Show()
					} else {
						table.owner.chips.Hide()
						table.owner.chips.SetSelected("DERO")
						table.owner.times.Hide()
					}
				} else {
					disableOwnerControls(true)
				}

				if rpc.Wallet.IsConnected() && checkHolderoContract(text) {
					table.tournament.Show()
				} else {
					table.tournament.Hide()
				}
			} else {
				signals.contract = false
				Settings.check.SetChecked(false)
				table.tournament.Hide()
			}
			fetchHolderoSC()
			wait = false
		}
	}

	return container.NewStack(bundle.Alpha120, container.NewBorder(contract_cont, nil, nil, nil, tabs))
}

// Holdero tab layout
func placeHoldero(change_screen *widget.Button, d *dreams.AppObject) *fyne.Container {
	H.Back = *container.NewWithoutLayout(
		HolderoTable(ResourcePokerTablePng),
		Player1_label(nil, nil, nil),
		Player2_label(nil, nil, nil),
		Player3_label(nil, nil, nil),
		Player4_label(nil, nil, nil),
		Player5_label(nil, nil, nil),
		Player6_label(nil, nil, nil),
		H.TopLabel)

	holdero_label := container.NewHBox(H.LeftLabel, layout.NewSpacer(), H.RightLabel)

	H.Front = *placeHolderoCards(d.Window)

	H.Actions = *container.NewVBox(
		layout.NewSpacer(),
		SitButton(d),
		LeaveButton(d),
		DealHandButton(d),
		CheckButton(d),
		BetButton(d),
		BetAmount())

	options := container.NewVBox(layout.NewSpacer(), AutoOptions(d), change_screen)

	holdero_actions := container.NewHBox(options, layout.NewSpacer(), TimeOutWarning(), layout.NewSpacer(), layout.NewSpacer(), &H.Actions)

	H.DApp = container.NewVBox(
		dwidget.LabelColor(holdero_label),
		&H.Back,
		&H.Front,
		layout.NewSpacer(),
		holdero_actions)

	return H.DApp
}

// Layout all objects for Holdero dApp
func LayoutAllItems(d *dreams.AppObject) *fyne.Container {
	H.LeftLabel = widget.NewLabel("")
	H.RightLabel = widget.NewLabel("")
	H.TopLabel = canvas.NewText(round.display.results, color.White)
	H.TopLabel.Move(fyne.NewPos(387, 204))
	H.LeftLabel.SetText("Seats: " + round.display.seats + "      Pot: " + round.display.pot + "      Blinds: " + round.display.blinds + "      Ante: " + round.display.ante + "      Dealer: " + round.display.dealer)
	H.RightLabel.SetText(round.display.readout + "      Player ID: " + round.display.playerId + "      Dero Balance: " + rpc.DisplayBalance("Dero") + "      Height: " + rpc.Wallet.Display.Height)

	var holdero_objs *fyne.Container
	var contract_objs *fyne.Container
	contract_change_screen := widget.NewButton("Tables", nil)
	contract_change_screen.Importance = widget.HighImportance
	contract_change_screen.OnTapped = func() {
		go func() {
			d.Window.Content().(*fyne.Container).Objects[1].(*fyne.Container).Objects[1].(*container.AppTabs).Selected().Content = contract_objs
		}()
	}

	holdero_objs = placeHoldero(contract_change_screen, d)
	contract_objs = placeContract(holdero_objs, d)

	// Main process
	go fetch(d, contract_objs)

	return contract_objs
}
