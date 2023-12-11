package holdero

import (
	"fmt"
	"image/color"
	"sort"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	dreams "github.com/dReam-dApps/dReams"
	"github.com/dReam-dApps/dReams/bundle"
	"github.com/dReam-dApps/dReams/dwidget"
	"github.com/dReam-dApps/dReams/menu"
	"github.com/dReam-dApps/dReams/rpc"
)

// Create a swap pair container
func CreateSwapContainer(pair string) (*dwidget.DeroAmts, *fyne.Container) {
	split := strings.Split(pair, "-")
	if len(split) != 2 {
		return dwidget.NewDeroEntry("", 0, 0), container.NewStack(widget.NewLabel("Invalid Pair"))
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
	swap2_entry := dwidget.NewDeroEntry("", incr, uint(menu.CoinDecimal(split[0])))
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
	swap1_entry := dwidget.NewDeroEntry("", incr, uint(menu.CoinDecimal(split[0])))
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

	assets := []string{}
	for asset := range rpc.Wallet.Display.Balance {
		assets = append(assets, asset)
	}

	sort.Strings(assets)

	menu.Assets.Balances = widget.NewList(
		func() int {
			return len(assets)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(assets[i] + fmt.Sprintf(": %s", rpc.DisplayBalance(assets[i])))
		})

	balance_tabs := container.NewAppTabs(
		container.NewTabItem("Balances", container.NewStack(menu.Assets.Balances)))

	var swap_entry *dwidget.DeroAmts
	var swap_boxes *fyne.Container

	max := container.NewStack()

	swap_button := widget.NewButton("Swap", nil)
	swap_button.OnTapped = func() {
		switch select_pair.Selected {
		case "DERO-dReams":
			f, err := strconv.ParseFloat(swap_entry.Text, 64)
			if err == nil && swap_entry.Validate() == nil {
				if amt := (f * 333) * 100000; amt > 0 {
					DreamsConfirm(1, amt, d)
				}
			}
		case "dReams-DERO":
			f, err := strconv.ParseFloat(swap_entry.Text, 64)
			if err == nil && swap_entry.Validate() == nil {
				if amt := f * 100000; amt > 0 {
					DreamsConfirm(2, amt, d)
				}
			}
		}
	}

	swap_entry, swap_boxes = CreateSwapContainer(select_pair.Selected)
	menu.Assets.Swap = container.NewBorder(select_pair, swap_button, nil, nil, swap_boxes)
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

	swap_tabs := container.NewAppTabs(container.NewTabItem("Swap", container.NewCenter(menu.Assets.Swap)))
	max.Add(swap_tabs)

	full := container.NewHSplit(container.NewStack(bundle.NewAlpha120(), balance_tabs), max)
	full.SetOffset(0.66)

	return full
}
