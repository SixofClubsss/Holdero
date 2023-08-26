package holdero

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	dreams "github.com/dReam-dApps/dReams"
	"github.com/dReam-dApps/dReams/menu"
	"github.com/dReam-dApps/dReams/rpc"
	dero "github.com/deroproject/derohe/rpc"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

// Holdero card face selection object
//   - Sets shared face url on selected
//   - If deck is not present locally, it is downloaded
func FaceSelect() fyne.Widget {
	options := []string{"Light", "Dark"}
	Settings.faces.Select = widget.NewSelect(options, func(s string) {
		switch Settings.faces.Select.SelectedIndex() {
		case -1:
			Settings.faces.Name = "light/"
		case 0:
			Settings.faces.Name = "light/"
		case 1:
			Settings.faces.Name = "dark/"
		default:
			Settings.faces.Name = s
		}

		check := strings.Trim(s, "0123456789")
		if check == "AZYPC" {
			Settings.faces.URL = "https://raw.githubusercontent.com/Azylem/" + s + "/main/" + s + ".zip?raw=true"
			dir := dreams.GetDir()
			face := dir + "/cards/" + Settings.faces.Name + "/card1.png"
			if !dreams.FileExists(face, "Holdero") {
				logger.Println("[Holdero] Downloading " + Settings.faces.URL)
				go GetZipDeck(Settings.faces.Name, Settings.faces.URL)
			}
		} else if check == "SIXPC" {
			Settings.faces.URL = "https://raw.githubusercontent.com/SixofClubsss/" + s + "/main/" + s + ".zip?raw=true"
			dir := dreams.GetDir()
			face := dir + "/cards/" + Settings.faces.Name + "/card1.png"
			if !dreams.FileExists(face, "Holdero") {
				logger.Println("[Holdero] Downloading " + Settings.faces.URL)
				go GetZipDeck(Settings.faces.Name, Settings.faces.URL)
			}
		} else if check == "High-Strangeness" {
			Settings.faces.URL = "https://raw.githubusercontent.com/High-Strangeness/High-Strangeness/main/HS_Deck/HS_Deck.zip?raw=true"
			dir := dreams.GetDir()
			face := dir + "/cards/" + Settings.faces.Name + "/card1.png"
			if !dreams.FileExists(face, "Holdero") {
				logger.Println("[Holdero] Downloading " + Settings.faces.URL)
				go GetZipDeck(Settings.faces.Name, Settings.faces.URL)
			}
		} else {
			Settings.faces.URL = ""
		}
	})

	Settings.faces.Select.SetSelectedIndex(0)
	Settings.faces.Select.PlaceHolder = "Faces"

	return Settings.faces.Select
}

// Holdero card back selection object for all games
//   - Sets shared back url on selected
//   - If back is not present locally, it is downloaded
func BackSelect() fyne.Widget {
	options := []string{"Light", "Dark"}
	Settings.backs.Select = widget.NewSelect(options, func(s string) {
		switch Settings.backs.Select.SelectedIndex() {
		case -1:
			Settings.backs.Name = "back1.png"
		case 0:
			Settings.backs.Name = "back1.png"
		case 1:
			Settings.backs.Name = "back2.png"
		default:
			Settings.backs.Name = s
		}

		go func() {
			check := strings.Trim(s, "0123456789")
			if check == "AZYPCB" {
				Settings.backs.URL = "https://raw.githubusercontent.com/Azylem/" + s + "/main/" + s + ".png"
				dir := dreams.GetDir()
				file := dir + "/cards/backs/" + s + ".png"
				if !dreams.FileExists(file, "Holdero") {
					logger.Println("[Holdero] Downloading " + Settings.backs.URL)
					downloadFileLocal("cards/backs/"+Settings.backs.Name+".png", Settings.backs.URL)
				}
			} else if check == "SIXPCB" {
				Settings.backs.URL = "https://raw.githubusercontent.com/SixofClubsss/" + s + "/main/" + s + ".png"
				dir := dreams.GetDir()
				back := dir + "/cards/backs/" + s + ".png"
				if !dreams.FileExists(back, "Holdero") {
					logger.Println("[Holdero] Downloading " + Settings.backs.URL)
					downloadFileLocal("cards/backs/"+Settings.backs.Name+".png", Settings.backs.URL)
				}
			} else if check == "High-Strangeness" {
				Settings.backs.URL = "https://raw.githubusercontent.com/High-Strangeness/" + s + "/main/HS_Back/HS_Back.png"
				dir := dreams.GetDir()
				back := dir + "/cards/backs/" + s + ".png"
				if !dreams.FileExists(back, "Holdero") {
					logger.Println("[Holdero] Downloading " + Settings.backs.URL)
					downloadFileLocal("cards/backs/"+Settings.backs.Name+".png", Settings.backs.URL)
				}
			} else {
				Settings.backs.URL = ""
			}
		}()
	})

	Settings.backs.Select.SetSelectedIndex(0)
	Settings.backs.Select.PlaceHolder = "Backs"

	return Settings.backs.Select
}

// dReams app avatar selection object
//   - Sets shared avatar url on selected
func AvatarSelect(asset_map map[string]string) fyne.Widget {
	options := []string{"None"}
	Settings.avatars.Select = widget.NewSelect(options, func(s string) {
		switch Settings.avatars.Select.SelectedIndex() {
		case -1:
			Settings.avatar.name = "None"
		case 0:
			Settings.avatar.name = "None"
		default:
			Settings.avatar.name = s
		}

		check := strings.Trim(s, " #0123456789")
		if check == "DBC" {
			Settings.avatar.url = "https://raw.githubusercontent.com/Azylem/" + s + "/main/" + s + ".PNG"
		} else if check == "HighStrangeness" {
			Settings.avatar.url = "https://raw.githubusercontent.com/High-Strangeness/High-Strangeness/main/" + s + "/" + s + ".jpg"
		} else if check == "AZYDS" {
			Settings.avatar.url = "https://raw.githubusercontent.com/Azylem/" + s + "/main/" + s + "-IC.png"
		} else if check == "SIXART" {
			Settings.avatar.url = "https://raw.githubusercontent.com/SixofClubsss/SIXART/main/" + s + "/" + s + "-IC.png"
		} else if check == "Dero Seals" {
			seal := strings.Trim(s, "Dero Sals#")
			Settings.avatar.url = "https://ipfs.io/ipfs/QmP3HnzWpiaBA6ZE8c3dy5ExeG7hnYjSqkNfVbeVW5iEp6/low/" + seal + ".jpg"
		} else if check == "Dero Degen" {
			degen := strings.Trim(s, "Dero gn#")
			Settings.avatar.url = "https://ipfs.io/ipfs/QmZM6onfiS8yUHFwfVypYnc6t9ZrvmpT43F9HFTou6LJyg/" + degen + ".png"
		} else if ValidAsset(asset_map[s]) {
			if url := menu.GetAssetUrl(1, asset_map[s]); url != "" {
				Settings.avatar.url = url
				return
			}

			agent := getAgentNumber(asset_map[s])
			if agent >= 0 && agent < 172 {
				Settings.avatar.url = "https://ipfs.io/ipfs/QmaRHXcQwbFdUAvwbjgpDtr5kwGiNpkCM2eDBzAbvhD7wh/low/" + strconv.Itoa(agent) + ".jpg"
			} else if agent < 1200 {
				Settings.avatar.url = "https://ipfs.io/ipfs/QmQQyKoE9qDnzybeDCXhyMhwQcPmLaVy3AyYAzzC2zMauW/low/" + strconv.Itoa(agent) + ".jpg"
			}
		} else if s == "None" {
			Settings.avatar.url = ""
		}
	})

	Settings.avatars.Select.PlaceHolder = "Avatar"

	return Settings.avatars.Select
}

// Confirm if asset map is valid
func ValidAsset(s string) bool {
	if s != "" && len(s) == 64 {
		return true
	}
	return false
}

// Rpc call to get A-Team agent number
func getAgentNumber(scid string) int {
	if rpc.Daemon.IsConnected() {
		rpcClientD, ctx, cancel := rpc.SetDaemonClient(rpc.Daemon.Rpc)
		defer cancel()

		var result *dero.GetSC_Result
		params := dero.GetSC_Params{
			SCID:      scid,
			Code:      false,
			Variables: true,
		}

		err := rpcClientD.CallFor(ctx, &result, "DERO.GetSC", params)
		if err != nil {
			logger.Errorln("[getAgentNumber]", err)
			return 1200
		}

		data := result.VariableStringKeys["metadata"]
		var agent menu.Agent

		hx, _ := hex.DecodeString(data.(string))
		if err := json.Unmarshal(hx, &agent); err == nil {
			return agent.ID
		}

	}
	return 1200
}

// Holdero shared cards toggle object
//   - Do not send a blank url
//   - If cards are not present locally, it is downloaded
func SharedDecks() fyne.Widget {
	options := []string{"Shared Decks"}
	Settings.shared = widget.NewRadioGroup(options, func(string) {
		if Settings.sharing || ((len(round.cards.Faces.Name) < 3 || len(round.cards.Backs.Name) < 3) && round.ID != 1) {
			logger.Println("[Holdero] Shared Decks Off")
			Settings.sharing = false
			Settings.faces.Select.Enable()
			Settings.backs.Select.Enable()
		} else {
			logger.Println("[Holdero] Shared Decks On")
			Settings.sharing = true
			if round.ID == 1 {
				if Settings.faces.Name != "" && Settings.faces.URL != "" && Settings.backs.Name != "" && Settings.backs.URL != "" {
					SharedDeckUrl(Settings.faces.Name, Settings.faces.URL, Settings.backs.Name, Settings.backs.URL)
					dir := dreams.GetDir()
					back := "/cards/backs/" + Settings.backs.Name + ".png"
					face := "/cards/" + Settings.faces.Name + "/card1.png"

					if !dreams.FileExists(dir+face, "Holdero") {
						go GetZipDeck(Settings.faces.Name, Settings.faces.URL)
					}

					if !dreams.FileExists(dir+back, "Holdero") {
						downloadFileLocal("cards/backs/"+Settings.backs.Name+".png", Settings.backs.URL)
					}
				}
			} else {
				Settings.faces.Select.Disable()
				Settings.backs.Select.Disable()
				dir := dreams.GetDir()
				back := "/cards/backs/" + round.cards.Backs.Name + ".png"
				face := "/cards/" + round.cards.Faces.Name + "/card1.png"

				if !dreams.FileExists(dir+face, "Holdero") {
					go GetZipDeck(round.cards.Faces.Name, round.cards.Faces.Url)
				}

				if !dreams.FileExists(dir+back, "Holdero") {
					downloadFileLocal("cards/backs/"+round.cards.Backs.Name+".png", round.cards.Backs.Url)
				}
			}
		}
	})

	Settings.shared.Disable()

	return Settings.shared
}

// Confirmation for dReams-Dero swap pairs
//   - c defines swap for Dero or dReams
//   - amt of Dero in atomic units
//   - Pass main window obj to reset to
func DreamsConfirm(c, amt float64, obj *fyne.Container, reset fyne.CanvasObject) fyne.CanvasObject {
	var text string
	dero := (amt / 100000) / 333
	ratio := math.Pow(10, float64(5))
	x := math.Round(dero*ratio) / ratio
	a := fmt.Sprint(strconv.FormatFloat(dero, 'f', 5, 64))
	switch c {
	case 1:
		text = fmt.Sprintf("You are about to swap %s DERO for %.5f dReams", a, amt/100000)
	case 2:
		text = fmt.Sprintf("You are about to swap %.5f dReams for %s Dero", amt/100000, a)
	}

	done := false
	label := widget.NewLabel(text)
	label.Wrapping = fyne.TextWrapWord
	label.Alignment = fyne.TextAlignCenter

	confirm_button := widget.NewButton("Confirm", func() {
		switch c {
		case 1:
			rpc.GetdReams(uint64(x * 100000))
		case 2:
			rpc.TradedReams(uint64(amt))
		default:

		}

		done = true
		obj.Objects[0] = reset
		obj.Objects[0].Refresh()
	})

	cancel_button := widget.NewButton("Cancel", func() {
		done = true
		obj.Objects[0] = reset
		obj.Objects[0].Refresh()
	})

	buttons := container.NewAdaptiveGrid(2, confirm_button, cancel_button)
	content := container.NewVBox(layout.NewSpacer(), label, layout.NewSpacer(), buttons)

	go func() {
		for rpc.IsReady() {
			time.Sleep(time.Second)
			if done {
				return
			}
		}

		obj.Objects[0] = reset
		obj.Objects[0].Refresh()
	}()

	return container.NewMax(content)
}
