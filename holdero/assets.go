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
	"github.com/dReam-dApps/dReams/bundle"
	"github.com/dReam-dApps/dReams/gnomes"
	"github.com/dReam-dApps/dReams/menu"
	"github.com/dReam-dApps/dReams/rpc"
	dero "github.com/deroproject/derohe/rpc"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

// Downloads card deck if it does not exists locally
func getCardDeck(url string) {
	Settings.faces.URL = url
	dir := dreams.GetDir()
	face := dir + "/cards/" + Settings.faces.Name + "/card1.png"
	if !dreams.FileExists(face, "Holdero") {
		logger.Println("[Holdero] Downloading " + Settings.faces.URL)
		go GetZipDeck(Settings.faces.Name, Settings.faces.URL)
	}
}

// Holdero card face selection object
//   - Sets shared face url on selected
//   - If deck is not present locally, it is downloaded
func FaceSelect(assets map[string]string) fyne.CanvasObject {
	var max *fyne.Container
	options := []string{"Light", "Dark"}
	icon := menu.AssetIcon(ResourceCardsCirclePng.StaticContent, "", 60)
	Settings.faces.Select = widget.NewSelect(options, nil)
	Settings.faces.Select.SetSelectedIndex(0)
	Settings.faces.Select.OnChanged = func(s string) {
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
			url := "https://raw.githubusercontent.com/Azylem/" + s + "/main/" + s + ".zip?raw=true"
			getCardDeck(url)
			max.Objects[1].(*fyne.Container).Objects[0].(*fyne.Container).Objects[0] = menu.SwitchProfileIcon("AZY-Playing card backs", s, gnomes.GetAssetUrl(1, assets[s]), 60)
		} else if check == "SIXPC" {
			url := "https://raw.githubusercontent.com/SixofClubsss/" + s + "/main/" + s + ".zip?raw=true"
			getCardDeck(url)
			max.Objects[1].(*fyne.Container).Objects[0].(*fyne.Container).Objects[0] = menu.SwitchProfileIcon("AZY-Playing card backs", s, gnomes.GetAssetUrl(1, assets[s]), 60)
		} else if check == "HS_Deck" {
			url := "https://raw.githubusercontent.com/High-Strangeness/High-Strangeness/main/" + s + "/" + s + ".zip?raw=true"
			getCardDeck(url)
			hs_icon := "https://raw.githubusercontent.com/High-Strangeness/High-Strangeness/main/HighStrangeness-IC.jpg"
			max.Objects[1].(*fyne.Container).Objects[0].(*fyne.Container).Objects[0] = menu.SwitchProfileIcon("High Strangeness", "HighStrangeness1", hs_icon, 60)
		} else {
			Settings.faces.URL = ""
			img := canvas.NewImageFromResource(ResourceCardsCirclePng)
			img.SetMinSize(fyne.NewSize(60, 60))
			max.Objects[1].(*fyne.Container).Objects[0].(*fyne.Container).Objects[0] = img
		}
	}

	Settings.faces.Select.PlaceHolder = "Faces:"
	max = container.NewBorder(nil, nil, icon, nil, container.NewVBox(Settings.faces.Select))

	return max
}

// Downloads card back if it does not exists locally
func getCardBack(s, url string) {
	Settings.backs.URL = url
	dir := dreams.GetDir()
	back := dir + "/cards/backs/" + s + ".png"
	if !dreams.FileExists(back, "Holdero") {
		logger.Println("[Holdero] Downloading " + Settings.backs.URL)
		downloadFileLocal("cards/backs/"+Settings.backs.Name+".png", Settings.backs.URL)
	}
}

// Holdero card back selection object for all games
//   - Sets shared back url on selected
//   - If back is not present locally, it is downloaded
func BackSelect(assets map[string]string) *fyne.Container {
	var max *fyne.Container
	options := []string{"Light", "Dark"}
	icon := menu.AssetIcon(ResourceCardsCirclePng.StaticContent, "", 60)
	Settings.backs.Select = widget.NewSelect(options, nil)
	Settings.backs.Select.SetSelectedIndex(0)
	Settings.backs.Select.OnChanged = func(s string) {
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
				url := "https://raw.githubusercontent.com/Azylem/" + s + "/main/" + s + ".png"
				getCardBack(s, url)
				max.Objects[1].(*fyne.Container).Objects[0].(*fyne.Container).Objects[0] = menu.SwitchProfileIcon("AZY-Playing card backs", s, gnomes.GetAssetUrl(1, assets[s]), 60)
			} else if check == "SIXPCB" {
				url := "https://raw.githubusercontent.com/SixofClubsss/" + s + "/main/" + s + ".png"
				getCardBack(s, url)
				max.Objects[1].(*fyne.Container).Objects[0].(*fyne.Container).Objects[0] = menu.SwitchProfileIcon("SIXPCB", s, gnomes.GetAssetUrl(1, assets[s]), 60)
			} else if check == "HS_Back" {
				url := "https://raw.githubusercontent.com/High-Strangeness/High-Strangeness/main/" + s + "/" + s + ".png"
				getCardBack(s, url)
				hs_icon := "https://raw.githubusercontent.com/High-Strangeness/High-Strangeness/main/HighStrangeness-IC.jpg"
				max.Objects[1].(*fyne.Container).Objects[0].(*fyne.Container).Objects[0] = menu.SwitchProfileIcon("High Strangeness", "HighStrangeness1", hs_icon, 60)
			} else {
				Settings.backs.URL = ""
				img := canvas.NewImageFromResource(ResourceCardsCirclePng)
				img.SetMinSize(fyne.NewSize(60, 60))
				max.Objects[1].(*fyne.Container).Objects[0].(*fyne.Container).Objects[0] = img
			}
		}()
	}

	Settings.backs.Select.PlaceHolder = "Backs:"
	max = container.NewBorder(nil, nil, icon, nil, container.NewVBox(Settings.backs.Select))

	return max
}

// Avatar selection object
//   - Sets shared avatar url on selected
func AvatarSelect(assets map[string]string) fyne.CanvasObject {
	var max *fyne.Container
	options := []string{"None"}
	icon := menu.AssetIcon(bundle.ResourceFigure1CirclePng.StaticContent, "", 60)
	Settings.avatars.Select = widget.NewSelect(options, nil)
	Settings.avatars.Select.SetSelectedIndex(0)
	Settings.avatars.Select.OnChanged = func(s string) {
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
			Settings.avatar.url = gnomes.GetAssetUrl(1, assets[s]) //"https://raw.githubusercontent.com/Azylem/" + s + "/main/" + s + ".PNG"
			max.Objects[1].(*fyne.Container).Objects[0].(*fyne.Container).Objects[0] = menu.SwitchProfileIcon("Death By Cupcake", s, Settings.avatar.url, 60)
		} else if check == "HighStrangeness" {
			Settings.avatar.url = gnomes.GetAssetUrl(1, assets[s]) //"https://raw.githubusercontent.com/High-Strangeness/High-Strangeness/main/" + s + "/" + s + ".jpg"
			max.Objects[1].(*fyne.Container).Objects[0].(*fyne.Container).Objects[0] = menu.SwitchProfileIcon("High Strangeness", s, Settings.avatar.url, 60)
		} else if check == "AZYDS" {
			Settings.avatar.url = gnomes.GetAssetUrl(1, assets[s]) //"https://raw.githubusercontent.com/Azylem/" + s + "/main/" + s + "-IC.png"
			max.Objects[1].(*fyne.Container).Objects[0].(*fyne.Container).Objects[0] = menu.SwitchProfileIcon("AZY-Deroscapes", s, Settings.avatar.url, 60)
		} else if check == "SIXART" {
			Settings.avatar.url = gnomes.GetAssetUrl(1, assets[s]) // "https://raw.githubusercontent.com/SixofClubsss/SIXART/main/" + s + "/" + s + "-IC.png"
			max.Objects[1].(*fyne.Container).Objects[0].(*fyne.Container).Objects[0] = menu.SwitchProfileIcon("SIXART", s, Settings.avatar.url, 60)
		} else if check == "Desperado" {
			Settings.avatar.url = gnomes.GetAssetUrl(1, assets[s])
			max.Objects[1].(*fyne.Container).Objects[0].(*fyne.Container).Objects[0] = menu.SwitchProfileIcon("Dero Desperados", s, Settings.avatar.url, 60)
		} else if check == "Gun" {
			Settings.avatar.url = gnomes.GetAssetUrl(1, assets[s])
			max.Objects[1].(*fyne.Container).Objects[0].(*fyne.Container).Objects[0] = menu.SwitchProfileIcon("Desperado Guns", s, Settings.avatar.url, 60)
		} else if check == "Dero Seals" {
			seal := strings.Trim(s, "Dero Sals#")
			Settings.avatar.url = "https://ipfs.io/ipfs/QmP3HnzWpiaBA6ZE8c3dy5ExeG7hnYjSqkNfVbeVW5iEp6/low/" + seal + ".jpg"
			max.Objects[1].(*fyne.Container).Objects[0].(*fyne.Container).Objects[0] = menu.SwitchProfileIcon("Dero Seals", s, Settings.avatar.url, 60)
		} else if check == "Dero Degen" {
			degen := strings.Trim(s, "Dero gn#")
			Settings.avatar.url = "https://ipfs.io/ipfs/QmZM6onfiS8yUHFwfVypYnc6t9ZrvmpT43F9HFTou6LJyg/" + degen + ".png"
			max.Objects[1].(*fyne.Container).Objects[0].(*fyne.Container).Objects[0] = menu.SwitchProfileIcon("Dero Degens", s, Settings.avatar.url, 60)
		} else if ValidAsset(assets[s]) {
			if url := gnomes.GetAssetUrl(1, assets[s]); url != "" {
				Settings.avatar.url = url
				max.Objects[1].(*fyne.Container).Objects[0].(*fyne.Container).Objects[0] = menu.SwitchProfileIcon("", s, url, 60)
				return
			}

			agent := getAgentNumber(assets[s])
			if agent >= 0 && agent < 172 {
				Settings.avatar.url = "https://ipfs.io/ipfs/QmaRHXcQwbFdUAvwbjgpDtr5kwGiNpkCM2eDBzAbvhD7wh/low/" + strconv.Itoa(agent) + ".jpg"
			} else if agent < 1200 {
				Settings.avatar.url = "https://ipfs.io/ipfs/QmQQyKoE9qDnzybeDCXhyMhwQcPmLaVy3AyYAzzC2zMauW/low/" + strconv.Itoa(agent) + ".jpg"
			} else {
				return
			}
			max.Objects[1].(*fyne.Container).Objects[0].(*fyne.Container).Objects[0] = menu.SwitchProfileIcon("Dero A-Team", s, Settings.avatar.url, 60)
		} else if s == "None" {
			Settings.avatar.url = ""
			img := canvas.NewImageFromResource(bundle.ResourceFigure1CirclePng)
			img.SetMinSize(fyne.NewSize(60, 60))
			max.Objects[1].(*fyne.Container).Objects[0].(*fyne.Container).Objects[0] = img
		}
	}

	Settings.avatars.Select.PlaceHolder = "Avatar:"
	max = container.NewBorder(nil, nil, icon, nil, container.NewVBox(Settings.avatars.Select))

	return max
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
func DreamsConfirm(c, amt float64, d *dreams.AppObject) {
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

	done := make(chan struct{})
	confirm := dialog.NewConfirm("Swap", text, func(b bool) {
		if b {
			switch c {
			case 1:
				if tx := rpc.GetdReams(uint64(x * 100000)); tx != "" {
					go menu.ShowTxDialog("Swap", fmt.Sprintf("TXID: %s", tx), tx, 3*time.Second, d.Window)
				} else {
					go menu.ShowTxDialog("Swap", "TX error, check logs", tx, 3*time.Second, d.Window)
				}
			case 2:
				if tx := rpc.TradedReams(uint64(amt)); tx != "" {
					go menu.ShowTxDialog("Swap", fmt.Sprintf("TXID: %s", tx), tx, 3*time.Second, d.Window)
				} else {
					go menu.ShowTxDialog("Swap", "TX error, check logs", tx, 3*time.Second, d.Window)
				}
			default:

			}
		}
		done <- struct{}{}
	}, d.Window)

	go menu.ShowConfirmDialog(done, confirm)
}
