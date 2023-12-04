package holdero

import (
	"fmt"
	"image/color"
	"strconv"
	"strings"
	"time"

	dreams "github.com/dReam-dApps/dReams"
	"github.com/dReam-dApps/dReams/dwidget"
	"github.com/dReam-dApps/dReams/menu"
	"github.com/dReam-dApps/dReams/rpc"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func HolderoIndicator() (ind menu.DreamsIndicator) {
	purple := color.RGBA{105, 90, 205, 210}
	blue := color.RGBA{31, 150, 200, 210}
	alpha := &color.RGBA{0, 0, 0, 0}

	ind.Img = canvas.NewImageFromResource(ResourceHolderoCirclePng)
	ind.Img.SetMinSize(fyne.NewSize(30, 30))
	ind.Rect = canvas.NewRectangle(alpha)
	ind.Rect.SetMinSize(fyne.NewSize(36, 36))

	ind.Animation = canvas.NewColorRGBAAnimation(purple, blue,
		time.Second*3, func(c color.Color) {
			if Odds.IsRunning() {
				ind.Rect.FillColor = c
				ind.Img.Show()
				canvas.Refresh(ind.Rect)
			} else {
				ind.Rect.FillColor = alpha
				ind.Img.Hide()
				canvas.Refresh(ind.Rect)
			}
		})

	ind.Animation.RepeatCount = fyne.AnimationRepeatForever
	ind.Animation.AutoReverse = true

	return
}

// Holdero owner control objects inside owners tab
func ownersBox(d *dreams.AppObject) fyne.CanvasObject {
	players := []string{"2 Players", "3 Players", "4 Players", "5 Players", "6 Players"}
	player_select := widget.NewSelect(players, nil)
	player_select.SetSelectedIndex(0)

	blinds_entry := dwidget.NewDeroEntry("Big Blind: ", 0.1, 1)
	blinds_entry.SetPlaceHolder("Dero:")
	blinds_entry.SetText("Big Blind: 0.0")
	blinds_entry.Validator = validation.NewRegexp(`^(Big Blind: )\d{1,}\.\d{0,1}$|^(Big Blind: )\d{1,}$`, "Int or float required")
	blinds_entry.OnChanged = func(s string) {
		if blinds_entry.Validate() != nil {
			blinds_entry.SetText("Big Blind: 0.0")
			table.owner.blinds = 0
		} else {
			trimmed := strings.Trim(s, "Biglnd: ")
			if f, err := strconv.ParseFloat(trimmed, 64); err == nil {
				if uint64(f*100000)%10000 == 0 {
					blinds_entry.SetText(blinds_entry.Prefix + strconv.FormatFloat(roundFloat(f, 1), 'f', int(blinds_entry.Decimal), 64))
					table.owner.blinds = uint64(roundFloat(f*100000, 1))
				} else {
					blinds_entry.SetText(blinds_entry.Prefix + strconv.FormatFloat(roundFloat(f, 1), 'f', int(blinds_entry.Decimal), 64))
				}
			}
		}
	}

	ante_entry := dwidget.NewDeroEntry("Ante: ", 0.1, 1)
	ante_entry.SetPlaceHolder("Ante:")
	ante_entry.SetText("Ante: 0.0")
	ante_entry.Validator = validation.NewRegexp(`^(Ante: )\d{1,}\.\d{0,1}$|^(Ante: )\d{1,}$`, "Int or float required")
	ante_entry.OnChanged = func(s string) {
		if ante_entry.Validate() != nil {
			ante_entry.SetText("Ante: 0.0")
			table.owner.ante = 0
		} else {
			trimmed := strings.Trim(s, ante_entry.Prefix)
			if f, err := strconv.ParseFloat(trimmed, 64); err == nil {
				if uint64(f*100000)%10000 == 0 {
					ante_entry.SetText(ante_entry.Prefix + strconv.FormatFloat(roundFloat(f, 1), 'f', int(ante_entry.Decimal), 64))
					table.owner.ante = uint64(roundFloat(f*100000, 1))
				} else {
					ante_entry.SetText(ante_entry.Prefix + strconv.FormatFloat(roundFloat(f, 1), 'f', int(ante_entry.Decimal), 64))
				}
			}
		}
	}

	options := []string{"DERO", "ASSET"}
	table.owner.chips = widget.NewRadioGroup(options, nil)
	table.owner.chips.SetSelected("DERO")
	table.owner.chips.Horizontal = true
	table.owner.chips.Required = true
	table.owner.chips.OnChanged = func(s string) {
		if s == "ASSET" {
			blinds_entry.Increment = 1
			blinds_entry.Decimal = 0
			blinds_entry.SetText("0")
			blinds_entry.Refresh()

			ante_entry.Increment = 1
			ante_entry.Decimal = 0
			ante_entry.SetText("0")
			ante_entry.Refresh()
		} else {
			blinds_entry.Increment = 0.1
			blinds_entry.Decimal = 1
			blinds_entry.Refresh()

			ante_entry.Increment = 0.1
			ante_entry.Decimal = 1
			ante_entry.Refresh()
		}
	}

	set_button := widget.NewButton("Set Table", func() {
		if round.display.seats != "" {
			info := fmt.Sprintf("This table is already open with %s seats", round.display.seats)
			dialog.NewInformation("Set Table", info, d.Window).Show()
			return
		}

		bb := table.owner.blinds
		sb := table.owner.blinds / 2
		ante := table.owner.ante
		chips := table.owner.chips.Selected
		if menu.Username != "" {
			trim := strings.TrimSuffix(player_select.Selected, " Players")
			if players, err := strconv.ParseInt(trim, 10, 64); err == nil {
				info := fmt.Sprintf("Setting table for,\n\nPlayers: (%d)\n\nChips: (%s)\n\nBlinds: (%s/%s)\n\nAnte: (%s)", players, chips, rpc.FromAtomic(bb, 5), rpc.FromAtomic(sb, 5), rpc.FromAtomic(ante, 5))
				dialog.NewConfirm("Set Table", info, func(b bool) {
					if b {
						SetTable(int(players), bb, sb, ante, chips, menu.Username, Settings.avatar.url)
					}
				}, d.Window).Show()
			}
		} else {
			dialog.NewInformation("Set Table", "Choose a name before setting table", d.Window).Show()
		}
	})

	clean_entry := dwidget.NewDeroEntry("Clean: ", 1, 0)
	clean_entry.AllowFloat = false
	clean_entry.SetPlaceHolder("Atomic:")
	clean_entry.SetText("Clean: 0")
	clean_entry.Validator = validation.NewRegexp(`^(Clean: )\d{1,}`, "Int required")
	clean_entry.OnChanged = func(s string) {
		if clean_entry.Validate() != nil {
			clean_entry.SetText("Clean: 0")
		}
	}

	clean_button := widget.NewButton("Clean Table", func() {
		if round.display.seats == "" {
			dialog.NewInformation("Clean Table", "Table needs to be opened to clean", d.Window).Show()
			return
		}

		trimmed := strings.Trim(clean_entry.Text, "Clean: ")
		c, err := strconv.Atoi(trimmed)
		if err != nil {
			dialog.NewInformation("Clean Table", "Invalid clean amount", d.Window).Show()
			logger.Errorln("[Holdero] Invalid Clean Amount")
			return
		}

		if c > int(round.Pot) {
			if round.Pot == 0 {
				dialog.NewInformation("Clean Table", "This pot is empty", d.Window).Show()
				return
			}

			dialog.NewInformation("Clean Table", fmt.Sprintf("There is only %s %s in this pot", rpc.FromAtomic(round.Pot, 5), table.owner.chips.Selected), d.Window).Show()
			return
		}

		if c == 0 {
			dialog.NewConfirm("Clean Table", "Would you like to reset this table?", func(b bool) {
				if b {
					CleanTable(0)
				}
			}, d.Window).Show()

			return
		}

		dialog.NewConfirm("Clean Table", fmt.Sprintf("Would you like to withdraw %s %s from this table and reset it? ", rpc.FromAtomic(trimmed, 5), table.owner.chips.Selected), func(b bool) {
			if b {
				CleanTable(uint64(c))
			}
		}, d.Window).Show()
	})

	table.owner.timeout = widget.NewButton("Timeout", func() {
		if round.display.seats == "" {
			dialog.NewInformation("Timeout", "This table is closed", d.Window).Show()
			return
		}

		dialog.NewConfirm("Timeout", "Would you like to timeout the current player at this table?", func(b bool) {
			if b {
				TimeOut()
			}
		}, d.Window).Show()
	})

	force := widget.NewButton("Force Start", func() {
		if round.display.seats == "" {
			dialog.NewInformation("Force Start", "This table is closed", d.Window).Show()
			return
		}

		if round.Pot != 0 {
			dialog.NewInformation("Force Start", "This table is already started", d.Window).Show()
			return
		}

		dialog.NewConfirm("Force Start", "Would you like to start this table before all seats are filled?", func(b bool) {
			if b {
				ForceStat()
			}
		}, d.Window).Show()
	})

	close := widget.NewButton("Close Table", func() {
		if round.Pot != 0 {
			dialog.NewInformation("Close Table", "There is still funds to be paid out at this table", d.Window).Show()
			return
		}

		if round.display.seats == "" {
			dialog.NewInformation("Close Table", "This table is already closed", d.Window).Show()
			return
		}

		dialog.NewConfirm("Close Table", "Would you like to close this table?", func(b bool) {
			if b {
				SetTable(1, 0, 0, 0, "", "", "")
			}
		}, d.Window).Show()
	})

	spacer := canvas.NewRectangle(color.RGBA{0, 0, 0, 0})
	spacer.SetMinSize(table.owner.chips.Size())

	first_form := []*widget.FormItem{}
	first_form = append(first_form, widget.NewFormItem("Seats", player_select))
	first_form = append(first_form, widget.NewFormItem("Chips", table.owner.chips))
	first_form = append(first_form, widget.NewFormItem("Blinds", blinds_entry))
	first_form = append(first_form, widget.NewFormItem("Ante", ante_entry))
	first_form = append(first_form, widget.NewFormItem("", set_button))
	first_form = append(first_form, widget.NewFormItem("", force))
	first_form = append(first_form, widget.NewFormItem("", close))
	first_form = append(first_form, widget.NewFormItem("", spacer))

	first_form = append(first_form, widget.NewFormItem("Clean Amount", clean_entry))
	first_form = append(first_form, widget.NewFormItem("", clean_button))
	first_form = append(first_form, widget.NewFormItem("", spacer))

	k_times := []string{"Off", "2m", "5m"}
	auto_remove := widget.NewSelect(k_times, func(s string) {
		switch s {
		case "Off":
			signals.times.kick = 0
		case "2m":
			signals.times.kick = 120
		case "5m":
			signals.times.kick = 300
		default:
			signals.times.kick = 0
		}
	})
	auto_remove.PlaceHolder = "Kick after:"

	p_times := []string{"30s", "60s"}
	delay := widget.NewSelect(p_times, func(s string) {
		switch s {
		case "30s":
			signals.times.delay = 30
		case "60s":
			signals.times.delay = 60
		default:
			signals.times.delay = 30
		}
	})
	delay.PlaceHolder = "Payout delay:"

	table.tournament = widget.NewButton("Tournament", func() {
		bal := rpc.TokenBalance(TourneySCID)
		balance := float64(bal) / 100000
		if balance == 0 {
			dialog.NewInformation("Tournament Deposit", "You have no Tournament chips to deposit", d.Window).Show()
			return
		}

		info := fmt.Sprintf("Would you like to deposit %s Tournament Chips into leader board contract?", strconv.FormatFloat(balance, 'f', 5, 64))
		dialog.NewConfirm("Tournament Deposit", info, func(b bool) {
			if b {
				TourneyDeposit(bal, menu.Username)
			}
		}, d.Window).Show()
	})

	table.tournament.Hide()

	second_form := []*widget.FormItem{}
	second_form = append(second_form, widget.NewFormItem("Auto Kick", auto_remove))
	second_form = append(second_form, widget.NewFormItem("", table.owner.timeout))
	second_form = append(second_form, widget.NewFormItem("", spacer))
	second_form = append(second_form, widget.NewFormItem("Payout Delay  ", delay))

	table.owner.owners_mid = container.NewVBox(widget.NewForm(second_form...))

	instructions := "To start a game on a table you own:\n---\nSelect number of seats at the table (6 max)\n\nSelect DERO or ASSET as chips\n\nSelect blinds and any required ante (can be 0)\n\nClick 'Set Table' to open your table for others to join\n\nClick 'Force Start' if you'd like to start the table before all the seats are filled\n\nWhen done playing, click 'Close Table' to close it\n\n'Clean Table' is your reset button, it shuffles the deck and move the turn to the next player,\nif clean amount is above 0 it will withdraw that amount (in atomic units) from the table\n\nAuto kick time default is off, and payout default is 30 seconds\n\nVisit dreamdapps.io for more docs"
	help_button := widget.NewButton("Help", func() {
		dialog.NewInformation("Owners Manual", instructions, d.Window).Show()
	})
	first_form = append(first_form, widget.NewFormItem("", table.tournament))
	first_form = append(first_form, widget.NewFormItem("", help_button))
	first_form = append(first_form, widget.NewFormItem("", layout.NewSpacer()))

	table.unlock = widget.NewButton("Unlock Holdero Contract", nil)
	table.unlock.Hide()

	table.new = widget.NewButton("New Holdero Table", nil)
	table.new.Hide()

	third_form := []*widget.FormItem{}
	third_form = append(third_form, widget.NewFormItem("", spacer))
	third_form = append(third_form, widget.NewFormItem("", container.NewVBox(table.unlock, table.new)))

	table.owner.owners_left = container.NewVBox(widget.NewForm(first_form...))
	table.owner.owners_left.Hide()

	help_spacer := canvas.NewRectangle(color.RGBA{0, 0, 0, 0})
	help_spacer.SetMinSize(fyne.NewSize(180, 0))

	return container.NewVScroll(container.NewVBox(table.owner.owners_left, table.owner.owners_mid, layout.NewSpacer(), container.NewVBox(widget.NewForm(third_form...))))
}

func holderoMenuConfirm(c int, d *dreams.AppObject) {
	gas_fee := 0.3
	unlock_fee := float64(rpc.UnlockFee) / 100000
	var text, title string
	switch c {
	case 1:
		table.unlock.Hide()
		title = "Holdero Unlock"
		text = `You are about to unlock and install your first Holdero Table
		
To help support the project, there is a ` + fmt.Sprintf("%.5f", unlock_fee) + ` DERO donation attached to preform this action

Unlocking a Holdero table gives you unlimited access to table uploads and all base level owner features

Including gas fee, transaction total will be ` + fmt.Sprintf("%0.5f", unlock_fee+gas_fee) + ` DERO


Select a public or private table

Public will show up in indexed list of tables

Private will not show up in the list

All standard tables can use dReams or DERO


HGC holders can choose to install a HGC table

Public table that uses HGC or DERO`
	case 2:
		table.new.Hide()
		title = "Holdero Install"
		text = `You are about to install a new Holdero table

Gas fee to install new table is 0.30000 DERO


Select a public or private table

Public will show up in indexed list of tables

Private will not show up in the list

All standard tables can use dReams or DERO


HGC holders can choose to install a HGC table

Public table that uses HGC or DERO`
	}

	label := widget.NewLabel(text)
	label.Wrapping = fyne.TextWrapWord
	label.Alignment = fyne.TextAlignCenter

	var choice *widget.Select
	var confirm *dialog.CustomDialog

	confirm_button := widget.NewButton("Confirm", func() {
		if choice.SelectedIndex() < 3 && choice.SelectedIndex() >= 0 {
			uploadHolderoContract(choice.SelectedIndex())
		}

		if c == 2 {
			table.new.Show()
		}

		confirm.Hide()
		confirm = nil
	})

	options := []string{"Public", "Private"}
	if hgc := rpc.TokenBalance(rpc.HgcSCID); hgc > 0 {
		options = append(options, "HGC")
	}

	choice = widget.NewSelect(options, func(s string) {
		if s == "Public" || s == "Private" || s == "HGC" {
			confirm_button.Show()
		} else {
			confirm_button.Hide()
		}
	})

	cancel_button := widget.NewButton("Cancel", func() {
		switch c {
		case 1:
			table.unlock.Show()
		case 2:
			table.new.Show()
		default:

		}

		confirm.Hide()
		confirm = nil
	})

	confirm_button.Hide()

	left := container.NewVBox(confirm_button)
	right := container.NewVBox(cancel_button)
	buttons := container.NewAdaptiveGrid(2, left, right)
	actions := container.NewVBox(choice, buttons)

	spacer := canvas.NewRectangle(color.Transparent)
	spacer.SetMinSize(fyne.NewSize(360, 100))

	confirm = dialog.NewCustom(title, "", container.NewStack(spacer, label), d.Window)
	confirm.SetButtons([]fyne.CanvasObject{actions})
	confirm.Show()

	go func() {
		for rpc.IsReady() {
			time.Sleep(time.Second)
		}

		if confirm != nil {
			confirm.Hide()
			confirm = nil
		}
	}()
}
