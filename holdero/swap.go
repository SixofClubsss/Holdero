package holdero

import (
	"fmt"
	"image/color"
	"net/url"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	dreams "github.com/dReam-dApps/dReams"
	"github.com/dReam-dApps/dReams/bundle"
	"github.com/dReam-dApps/dReams/dwidget"
	"github.com/dReam-dApps/dReams/menu"
	"github.com/dReam-dApps/dReams/rpc"
)

// Create a swap pair container
func CreateSwapContainer(pair string) (*dwidget.AmountEntry, *fyne.Container) {
	split := strings.Split(pair, "-")
	if len(split) != 2 {
		return dwidget.NewAmountEntry("", 0, 0), container.NewStack(widget.NewLabel("Invalid Pair"))
	}

	incr := 0.1
	switch split[0] {
	case "dReams":
		incr = 1
	}

	color1 := color.RGBA{0, 0, 0, 0}
	color2 := color.RGBA{0, 0, 0, 0}
	image1 := canvas.NewImageFromResource(ResourceSwapFrame1Png)
	image2 := canvas.NewImageFromResource(ResourceSwapFrame2Png)

	rect2 := canvas.NewRectangle(color2)
	rect2.SetMinSize(fyne.NewSize(200, 100))
	swap2_label := canvas.NewText(split[1], bundle.TextColor)
	swap2_label.Alignment = fyne.TextAlignCenter
	swap2_label.TextSize = 18
	swap2_entry := dwidget.NewAmountEntry("", incr, uint(menu.CoinDecimal(split[0])))
	swap2_entry.SetText("0")
	swap2_entry.Disable()

	pad2 := container.NewBorder(layout.NewSpacer(), layout.NewSpacer(), layout.NewSpacer(), layout.NewSpacer(), swap2_entry)

	swap2 := container.NewBorder(nil, pad2, nil, nil, container.NewCenter(swap2_label))
	cont2 := container.NewStack(rect2, image2, swap2)

	rect1 := canvas.NewRectangle(color1)
	rect1.SetMinSize(fyne.NewSize(200, 100))
	swap1_label := canvas.NewText(split[0], bundle.TextColor)
	swap1_label.Alignment = fyne.TextAlignCenter
	swap1_label.TextSize = 18
	swap1_entry := dwidget.NewAmountEntry("", incr, uint(menu.CoinDecimal(split[0])))
	swap1_entry.SetText("0")
	swap1_entry.Validator = validation.NewRegexp(`^\d{1,}\.\d{1,5}$|^[^0.]\d{0,}$`, "Int or float required")
	swap1_entry.OnChanged = func(s string) {
		switch pair {
		case "DERO-dReams", "dReams-DERO":
			if f, err := strconv.ParseFloat(s, 64); err == nil {
				ex := float64(333)
				if split[0] == "dReams" {
					new := f / ex
					swap2_entry.SetText(fmt.Sprintf("%.5f", new))
					return
				}

				new := f * ex
				swap2_entry.SetText(fmt.Sprintf("%.5f", new))

			}
		default:
			// other pairs
		}
	}

	pad1 := container.NewBorder(layout.NewSpacer(), layout.NewSpacer(), layout.NewSpacer(), layout.NewSpacer(), swap1_entry)

	swap1 := container.NewBorder(nil, pad1, nil, nil, container.NewCenter(swap1_label))
	cont1 := container.NewStack(rect1, image1, swap1)

	return swap1_entry, container.NewAdaptiveGrid(2, cont1, cont2)
}

// Balance and swap container
func PlaceSwap(d *dreams.AppObject) *container.Split {
	pair_opts := []string{"DERO-dReams", "dReams-DERO"}
	select_pair := widget.NewSelect(pair_opts, nil)
	select_pair.PlaceHolder = "Pairs"
	select_pair.SetSelectedIndex(0)

	var selectedAsset string
	_, menu.Assets.Balances.SCIDs = rpc.Wallet.Balances()

	menu.Assets.Balances.List = widget.NewList(
		func() int {
			return len(menu.Assets.Balances.SCIDs)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(fmt.Sprintf("%s: %s", menu.Assets.Balances.SCIDs[i], rpc.Wallet.BalanceF(menu.Assets.Balances.SCIDs[i])))
		})

	balance_tabs := container.NewAppTabs(
		container.NewTabItem("Balances", container.NewStack(menu.Assets.Balances.List)),
		container.NewTabItem("Explorer", layout.NewSpacer()),
		container.NewTabItem("Send Message", layout.NewSpacer()))

	balance_tabs.OnSelected = func(ti *container.TabItem) {
		switch ti.Text {
		case "Explorer":
			link, _ := url.Parse("https://explorer.dero.io")
			fyne.CurrentApp().OpenURL(link)
			balance_tabs.SelectIndex(0)
		case "Send Message":
			if rpc.Wallet.IsConnected() {
				go menu.SendMessageMenu("", bundle.ResourceDReamsIconAltPng)
			} else {
				dialog.NewInformation("Not Connected", "Connect a wallet to send a message", d.Window).Show()
			}
			balance_tabs.SelectIndex(0)
		}
	}

	var swap_entry *dwidget.AmountEntry
	var swap_boxes *fyne.Container

	max := container.NewStack()

	swap_button := widget.NewButton("Swap", nil)
	swap_button.Importance = widget.HighImportance
	swap_button.OnTapped = func() {
		switch select_pair.Selected {
		case "DERO-dReams":
			f, err := strconv.ParseFloat(swap_entry.Text, 64)
			if err == nil && swap_entry.Validate() == nil {
				if amt := (f * 333) * 100000; amt > 0 {
					DreamsConfirm(1, amt, d)
				}
			} else {
				dialog.NewInformation("Swap", "Amount error", d.Window).Show()
			}
		case "dReams-DERO":
			f, err := strconv.ParseFloat(swap_entry.Text, 64)
			if err == nil && swap_entry.Validate() == nil {
				if amt := f * 100000; amt > 0 {
					DreamsConfirm(2, amt, d)
				}
			} else {
				dialog.NewInformation("Swap", "Amount error", d.Window).Show()
			}
		}
	}

	swap_entry, swap_boxes = CreateSwapContainer(select_pair.Selected)
	menu.Assets.Swap = container.NewBorder(nil, container.NewBorder(nil, nil, nil, swap_button, select_pair), nil, nil, swap_boxes)
	menu.Assets.Swap.Hide()

	select_pair.OnChanged = func(s string) {
		split := strings.Split(s, "-")
		if len(split) != 2 {
			return
		}

		swap_entry, swap_boxes = CreateSwapContainer(s)

		menu.Assets.Swap.Objects[0] = swap_boxes
		menu.Assets.Swap.Refresh()
	}

	btnTokenRemove := widget.NewButtonWithIcon("", dreams.FyneIcon("contentRemove"), nil)
	btnTokenRemove.Importance = widget.LowImportance
	btnTokenRemove.Hide()
	btnTokenRemove.OnTapped = func() {
		if selectedAsset == "" {
			dialog.NewError(fmt.Errorf("select a token to remove"), d.Window).Show()
			return
		}

		if selectedAsset == "DERO" || selectedAsset == "dReams" || selectedAsset == "HGC" {
			dialog.NewInformation("Default", "Can't remove default assets", d.Window).Show()
			return
		}

		dialog.NewConfirm("Remove", fmt.Sprintf("Would you like to remove %s from your balance list?", selectedAsset), func(b bool) {
			if b {
				menu.Assets.Balances.List.UnselectAll()
				btnTokenRemove.Hide()
				rpc.Wallet.TokenRemove(selectedAsset)
				tkn, names := rpc.Wallet.Balances()
				menu.Assets.Balances.SCIDs = names
				err := dreams.StoreAccount(dreams.AddAccountData(tkn, "tokens"))
				if err != nil {
					err = fmt.Errorf("storing account %s", err)
					dialog.NewError(err, d.Window).Show()
					logger.Errorf("[%s] %s\n", d.Name(), err)
				}
			}
		}, d.Window).Show()
	}

	btnTokenDefault := widget.NewButtonWithIcon("", dreams.FyneIcon("mediaReplay"), nil)
	btnTokenDefault.Importance = widget.LowImportance
	btnTokenDefault.OnTapped = func() {
		dialog.NewConfirm("Default Balances", "Set balances to default assets", func(b bool) {
			if b {
				menu.Assets.Balances.List.UnselectAll()
				btnTokenRemove.Hide()
				rpc.Wallet.SetDefaultTokens()
				tkn, names := rpc.Wallet.Balances()
				menu.Assets.Balances.SCIDs = names
				err := dreams.StoreAccount(dreams.AddAccountData(tkn, "tokens"))
				if err != nil {
					err = fmt.Errorf("storing account %s", err)
					dialog.NewError(err, d.Window).Show()
					logger.Errorf("[%s] %s\n", d.Name(), err)
				}
			}
		}, d.Window).Show()
	}

	btnTokenAdd := widget.NewButtonWithIcon("", dreams.FyneIcon("contentAdd"), nil)
	btnTokenAdd.Importance = widget.LowImportance
	btnTokenAdd.OnTapped = func() {
		entryName := widget.NewEntry()
		entryName.SetPlaceHolder("Name:")
		entryName.Validator = func(s string) error {
			if s == "" {
				return fmt.Errorf("enter a name")
			}

			return nil
		}

		entrySCID := widget.NewEntry()
		entrySCID.SetPlaceHolder("SCID:")
		entrySCID.Validator = func(s string) error {
			if len(s) == 64 {
				return nil
			}

			return fmt.Errorf("not a valid scid")
		}

		entryDeci := dwidget.NewAmountEntry("", 1, 0)
		entryDeci.SetPlaceHolder("Decimal:")
		entryDeci.Validator = func(s string) error {
			u, err := entryDeci.Uint64()
			if err != nil {
				return fmt.Errorf("enter a number 0-5")
			}

			if u > 5 {
				return fmt.Errorf("less than 6")
			}

			return nil
		}

		var add *dialog.CustomDialog
		btnAdd := widget.NewButton("Add", nil)
		btnAdd.Importance = widget.HighImportance
		btnAdd.OnTapped = func() {
			err := entryName.Validate()
			if err != nil {
				dialog.NewError(err, d.Window).Show()
				return
			}

			err = entrySCID.Validate()
			if err != nil {
				dialog.NewError(err, d.Window).Show()
				return
			}

			err = entryDeci.Validate()
			if err != nil {
				dialog.NewError(err, d.Window).Show()
				return
			}

			u, err := entryDeci.Uint64()
			if err != nil {
				dialog.NewError(err, d.Window).Show()
				return
			}

			err = rpc.Wallet.TokenAdd(entryName.Text, entrySCID.Text, int(u))
			if err != nil {
				dialog.NewError(err, d.Window).Show()
				return
			}

			tkn, names := rpc.Wallet.Balances()
			menu.Assets.Balances.SCIDs = names

			err = dreams.StoreAccount(dreams.AddAccountData(tkn, "tokens"))
			if err != nil {
				err = fmt.Errorf("storing account %s", err)
				dialog.NewError(err, d.Window).Show()
				logger.Errorf("[%s] %s\n", d.Name(), err)
			}

			add.Hide()
			add = nil
		}

		btnCancel := widget.NewButton("Cancel", func() {
			add.Hide()
			add = nil
		})

		var form []*widget.FormItem
		form = append(form, widget.NewFormItem("Name", entryName))
		form = append(form, widget.NewFormItem("", container.NewVBox(dwidget.NewLine(20, 1, bundle.TextColor))))
		form = append(form, widget.NewFormItem("SCID", entrySCID))
		form = append(form, widget.NewFormItem("Decimal", entryDeci))

		add = dialog.NewCustom("Add Token", "", container.NewStack(dwidget.NewSpacer(400, 0), widget.NewForm(form...)), d.Window)

		add.SetButtons([]fyne.CanvasObject{btnAdd, btnCancel})

		add.Show()
	}

	menu.Assets.Balances.List.OnSelected = func(id widget.ListItemID) {
		str := menu.Assets.Balances.SCIDs[id]
		selectedAsset = str
		if str == "DERO" || str == "dReams" || str == "HGC" {
			btnTokenRemove.Hide()
		} else {
			btnTokenRemove.Show()
		}
	}

	swap_tabs := container.NewAppTabs(container.NewTabItem("Swap", container.NewCenter(menu.Assets.Swap)))
	max.Add(swap_tabs)

	menu.Assets.AddRmv = container.NewVBox(layout.NewSpacer(), container.NewHBox(layout.NewSpacer(), btnTokenDefault, btnTokenRemove, btnTokenAdd))
	menu.Assets.AddRmv.Hide()

	full := container.NewHSplit(container.NewStack(bundle.NewAlpha120(), balance_tabs, menu.Assets.AddRmv), max)
	full.SetOffset(0.66)

	return full
}
