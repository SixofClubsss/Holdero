package holdero

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/civilware/Gnomon/structures"
	dreams "github.com/dReam-dApps/dReams"
	"github.com/dReam-dApps/dReams/bundle"
	"github.com/dReam-dApps/dReams/dwidget"
	"github.com/dReam-dApps/dReams/gnomes"
	"github.com/dReam-dApps/dReams/menu"
	"github.com/dReam-dApps/dReams/rpc"
	"github.com/sirupsen/logrus"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type holderoObjects struct {
	Public     dwidget.Lists
	Favorites  dwidget.Lists
	Owned      dwidget.Lists
	entry      *widget.SelectEntry
	unlock     *widget.Button
	new        *widget.Button
	sit        *widget.Button
	leave      *widget.Button
	deal       *widget.Button
	bet        *widget.Button
	check      *widget.Button
	tournament *widget.Button
	betEntry   *dwidget.DeroAmts
	warning    *fyne.Container
	owner      struct {
		valid    bool
		blinds   uint64
		ante     uint64
		chips    *widget.RadioGroup
		settings *fyne.Container
		times    *fyne.Container
	}
}

type tableInfo struct {
	name    string
	desc    string
	scid    string
	owner   string
	chips   string
	blinds  string
	version string
	last    string
	seats   string
	rating  uint64
	image   *canvas.Image
}

type settings struct {
	avatar struct {
		name string
		url  string
	}
	synced  bool
	sharing bool
	auto    struct {
		check bool
		deal  bool
	}
	check   *widget.Check
	avatars dreams.AssetSelect
	faces   dreams.AssetSelect
	backs   dreams.AssetSelect
	tools   *widget.Button
	shared  *widget.RadioGroup
}

var publicTables []tableInfo
var ownedTables []tableInfo
var favoriteTables []tableInfo
var table holderoObjects
var Settings settings
var logger = structures.Logger.WithFields(logrus.Fields{})

func DreamsMenuIntro() (entries map[string][]string) {
	entries = map[string][]string{
		"Holdero": {
			"Multiplayer Texas Hold'em style on chain poker",
			"No limit, single raise game. Table owners choose game params",
			"Six players max at a table",
			"No side pots, must call or fold",
			"Standard tables can be public or private, and can use Dero or dReam Tokens",
			"dReam Tools", "Tournament tables can be set up to use any Token",
			"View table listings or launch your own Holdero contract in the owned tab"},

		"dReam Tools": {
			"A suite of tools for Holdero, unlocked with ownership of a AZY or SIX playing card assets",
			"Odds calculator",
			"Bot player with 12 customizable parameters",
			"Track playing stats for user and bot players"},
	}

	return
}

func OnConnected() {
	table.entry.CursorColumn = 1
	table.entry.Refresh()
	if len(rpc.Wallet.Address) == 66 {
		CheckExistingKey()
		menu.Assets.Names.ClearSelected()
		menu.Assets.Names.Options = []string{}
		menu.Assets.Names.Refresh()
		menu.Assets.Names.Options = append(menu.Assets.Names.Options, rpc.Wallet.Address[0:12])
		if menu.Assets.Names.Options != nil {
			menu.Assets.Names.SetSelectedIndex(0)
		}
	}
}

func (s *settings) EnableCardSelects() {
	if round.ID == 1 {
		s.faces.Select.Enable()
		s.backs.Select.Enable()
	}
}

func (s *settings) ClearAssets() {
	s.faces.Select.Options = []string{}
	s.backs.Select.Options = []string{}
	s.avatars.Select.Options = []string{}
}

func (s *settings) SortCardAssets() {
	sort.Strings(s.faces.Select.Options)
	sort.Strings(s.backs.Select.Options)

	ld := []string{"Light", "Dark"}
	s.faces.Select.Options = append(ld, s.faces.Select.Options...)
	s.backs.Select.Options = append(ld, s.backs.Select.Options...)
}

func (s *settings) SortAvatarAsset() {
	sort.Strings(s.avatars.Select.Options)
	s.avatars.Select.Options = append([]string{"None"}, s.avatars.Select.Options...)
}

func (s *settings) AddAvatar(add, check string) {
	if check == rpc.Wallet.Address {
		avatars := s.avatars.Select.Options
		new_avatar := append(avatars, add)
		s.avatars.Select.Options = new_avatar
		s.avatars.Select.Refresh()
	}
}

func (s *settings) AddFaces(add, check string) {
	if check == rpc.Wallet.Address {
		current := s.faces.Select.Options
		new := append(current, add)
		s.faces.Select.Options = new
		s.faces.Select.Refresh()
	}
}

func (s *settings) CurrentFaces() []string {
	return s.faces.Select.Options
}

func (s *settings) CurrentBacks() []string {
	return s.backs.Select.Options
}

func (s *settings) AddBacks(add, check string) {
	if check == rpc.Wallet.Address {
		current := s.backs.Select.Options
		new := append(current, add)
		s.backs.Select.Options = new
		s.backs.Select.Refresh()
	}
}

func initValues() {
	signals.times.delay = 30
	signals.times.kick = 0
	Odds.Stop()
	Settings.faces.Name = "light/"
	Settings.backs.Name = "back1.png"
	Settings.avatar.name = "None"
	Settings.faces.URL = ""
	Settings.backs.URL = ""
	Settings.avatar.url = ""
	Settings.auto.deal = false
	Settings.auto.check = false
	signals.sit = true
	autoBetDefault()
}

// Sets contract entry text when list item pressed, returns item for use in list funcs
func setHolderoEntryText(t tableInfo) (item tableInfo) {
	if len(t.scid) == 64 {
		potIsEmpty(0)
		item = t
		table.entry.SetText(t.scid)
		signals.times.block = rpc.Wallet.Height
	}

	return
}

// List item with table info and icon
func listItem() fyne.CanvasObject {
	spacer := canvas.NewImageFromImage(nil)
	spacer.SetMinSize(fyne.NewSize(70, 70))
	return container.NewStack(container.NewBorder(
		nil,
		nil,
		container.NewStack(container.NewCenter(spacer), canvas.NewImageFromImage(nil)),
		nil,
		container.NewBorder(
			nil,
			container.NewHBox(canvas.NewImageFromImage(nil), widget.NewLabel("")),
			container.NewStack(canvas.NewImageFromImage(nil), canvas.NewImageFromImage(nil)),
			nil,
			widget.NewLabel(""))))
}

// Update func for listItem
func updateListItem(i widget.ListItemID, o fyne.CanvasObject, t []tableInfo) {
	if i > len(t)-1 {
		return
	}

	header := fmt.Sprintf("%s   %s   %s", t[i].name, t[i].desc, t[i].scid)
	if o.(*fyne.Container).Objects[0].(*fyne.Container).Objects[0].(*fyne.Container).Objects[0].(*widget.Label).Text != header {

		o.(*fyne.Container).Objects[0].(*fyne.Container).Objects[0].(*fyne.Container).Objects[0].(*widget.Label).SetText(header)
		if t[i].chips != "" {
			o.(*fyne.Container).Objects[0].(*fyne.Container).Objects[0].(*fyne.Container).Objects[1].(*fyne.Container).Objects[1].(*widget.Label).SetText(fmt.Sprintf("%s   %s   %s   Last move: %s   v%s", t[i].seats, t[i].chips, t[i].blinds, t[i].last, t[i].version))
		}

		badge := canvas.NewImageFromResource(menu.DisplayRating(menu.Control.Ratings[t[i].scid]))
		badge.SetMinSize(fyne.NewSize(35, 35))
		o.(*fyne.Container).Objects[0].(*fyne.Container).Objects[0].(*fyne.Container).Objects[1].(*fyne.Container).Objects[0] = badge

		if t[i].image != nil {
			o.(*fyne.Container).Objects[0].(*fyne.Container).Objects[1].(*fyne.Container).Objects[0].(*fyne.Container).Objects[0] = t[i].image
		} else {
			unknown := canvas.NewImageFromResource(ResourceUnknownAvatarPng)
			unknown.SetMinSize(fyne.NewSize(66, 66))
			o.(*fyne.Container).Objects[0].(*fyne.Container).Objects[1].(*fyne.Container).Objects[0].(*fyne.Container).Objects[0] = unknown
		}

		frame := canvas.NewImageFromResource(bundle.ResourceFramePng)
		frame.SetMinSize(fyne.NewSize(70, 70))
		o.(*fyne.Container).Objects[0].(*fyne.Container).Objects[1].(*fyne.Container).Objects[1] = frame
		o.Refresh()
	}
}

// Public Holdero table list object
func publicList(d *dreams.AppObject) fyne.CanvasObject {
	table.Public.List = widget.NewList(
		func() int {
			return len(publicTables)
		},
		listItem,
		func(i widget.ListItemID, o fyne.CanvasObject) {
			updateListItem(i, o, publicTables)
		})

	var item tableInfo

	table.Public.List.OnSelected = func(id widget.ListItemID) {
		if gnomes.Connected() {
			item = setHolderoEntryText(publicTables[id])
			table.Favorites.List.UnselectAll()
			table.Owned.List.UnselectAll()
		}
	}

	save_favorite := widget.NewButton("Favorite", func() {
		if item.name == "" {
			return
		}

		for _, sc := range table.Favorites.SCIDs {
			if item.scid == sc {
				return
			}
		}
		favoriteTables = append(favoriteTables, item)
		table.Favorites.SCIDs = append(table.Favorites.SCIDs, item.scid)
	})
	save_favorite.Importance = widget.LowImportance

	rate_contract := widget.NewButton("Rate", func() {
		if len(round.Contract) == 64 {
			if !checkTableOwner(round.Contract) {
				menu.RateConfirm(round.Contract, d)
			} else {
				dialog.NewInformation("Can't rate", "You are the owner of this SCID", d.Window).Show()
				logger.Warnln("[Holdero] Can't rate, you own this contract")
			}
		}
	})
	rate_contract.Importance = widget.LowImportance

	return container.NewBorder(
		nil,
		container.NewBorder(nil, nil, save_favorite, rate_contract, layout.NewSpacer()),
		nil,
		nil,
		table.Public.List)
}

// Favorite Holdero tables list object
func favoritesList() fyne.CanvasObject {
	table.Favorites.List = widget.NewList(
		func() int {
			return len(favoriteTables)
		},
		listItem,
		func(i widget.ListItemID, o fyne.CanvasObject) {
			updateListItem(i, o, favoriteTables)
		})

	var item tableInfo

	table.Favorites.List.OnSelected = func(id widget.ListItemID) {
		if gnomes.Connected() {
			item = setHolderoEntryText(favoriteTables[id])
			table.Public.List.UnselectAll()
			table.Owned.List.UnselectAll()
		}
	}

	remove := widget.NewButton("Remove", func() {
		if len(table.Favorites.SCIDs) > 0 {
			table.Favorites.List.UnselectAll()
			table.Favorites.RemoveSCID(item.scid)
			for i := range favoriteTables {
				if favoriteTables[i].scid == item.scid {
					copy(favoriteTables[i:], favoriteTables[i+1:])
					favoriteTables[len(favoriteTables)-1] = tableInfo{}
					favoriteTables = favoriteTables[:len(favoriteTables)-1]
					break
				}
			}
		}
		table.Favorites.List.Refresh()
	})
	remove.Importance = widget.LowImportance

	return container.NewBorder(
		nil,
		container.NewBorder(nil, nil, nil, remove, layout.NewSpacer()),
		nil,
		nil,
		table.Favorites.List)
}

// Returns table.Favorites.SCIDs
func GetFavoriteTables() []string {
	return table.Favorites.SCIDs
}

// Set table.Favorites.SCIDs
func SetFavoriteTables(fav []string) {
	table.Favorites.SCIDs = fav
}

// Owned Holdero tables list object
func ownedList(d *dreams.AppObject) fyne.CanvasObject {
	table.Owned.List = widget.NewList(
		func() int {
			return len(ownedTables)
		},
		listItem,
		func(i widget.ListItemID, o fyne.CanvasObject) {
			updateListItem(i, o, ownedTables)
		})

	table.Owned.List.OnSelected = func(id widget.ListItemID) {
		if gnomes.Connected() {
			setHolderoEntryText(ownedTables[id])
			table.Public.List.UnselectAll()
			table.Favorites.List.UnselectAll()
		}
	}

	return container.NewBorder(
		nil,
		nil,
		nil,
		ownersBox(d),
		table.Owned.List)
}

// Table owner name and avatar objects
//   - Pass a and f as avatar and its frame resource, shared avatar is set here if image exists
//   - Pass t for player's turn frame resource
func Player1_label(a, f, t fyne.Resource) fyne.CanvasObject {
	var name fyne.CanvasObject
	var frame fyne.CanvasObject
	var out fyne.CanvasObject
	if signals.In1 {
		if round.Turn == 1 {
			name = canvas.NewText(round.p1.name, color.RGBA{105, 90, 205, 210})
		} else {
			name = canvas.NewText(round.p1.name, color.White)
		}
	} else {
		name = canvas.NewText("", color.Transparent)
	}

	if a != nil && signals.In1 {
		if round.p1.url != "" {
			if !shared.loaded.p1 {
				shared.loaded.p1 = true
				if img, err := dreams.DownloadCanvas(round.p1.url, "P1"); err == nil {
					shared.avatar.p1 = &img
				} else {
					logger.Errorln("[Holdero] Player 1 avatar:", err)
					shared.avatar.p1 = canvas.NewImageFromResource(a)
				}
			}
		} else {
			shared.loaded.p1 = false
			shared.avatar.p1 = canvas.NewImageFromResource(a)
		}

		if round.Turn == 1 {
			frame = canvas.NewImageFromResource(t)
		} else {
			frame = canvas.NewImageFromResource(f)
		}
	} else {
		shared.avatar.p1 = canvas.NewImageFromImage(nil)
		frame = canvas.NewImageFromImage(nil)
	}

	if signals.Out1 {
		out = canvas.NewText("Sitting out", color.White)
		out.Resize(fyne.NewSize(100, 25))
		out.Move(fyne.NewPos(253, 45))
	} else {
		out = canvas.NewText("", color.RGBA{0, 0, 0, 0})
	}

	name.Resize(fyne.NewSize(100, 25))
	name.Move(fyne.NewPos(242, 20))

	shared.avatar.p1.Resize(fyne.NewSize(74, 74))
	shared.avatar.p1.Move(fyne.NewPos(359, 50))

	frame.Resize(fyne.NewSize(78, 78))
	frame.Move(fyne.NewPos(357, 48))

	return container.NewWithoutLayout(name, out, shared.avatar.p1, frame)
}

// Player 2 name and avatar objects
//   - Pass a and f as avatar and its frame resource, shared avatar is set here if image exists
//   - Pass t for player's turn frame resource
func Player2_label(a, f, t fyne.Resource) fyne.CanvasObject {
	var name fyne.CanvasObject
	var frame fyne.CanvasObject
	if signals.In2 {
		if round.Turn == 2 {
			name = canvas.NewText(round.p2.name, color.RGBA{105, 90, 205, 210})
		} else {
			name = canvas.NewText(round.p2.name, color.White)
		}
	} else {
		name = canvas.NewText("", color.Transparent)
	}

	if a != nil && signals.In2 {
		if round.p2.url != "" {
			if !shared.loaded.p2 {
				shared.loaded.p2 = true
				if img, err := dreams.DownloadCanvas(round.p2.url, "P2"); err == nil {
					shared.avatar.p2 = &img
				} else {
					logger.Errorln("[Holdero] Player 2 avatar:", err)
					shared.avatar.p2 = canvas.NewImageFromResource(a)
				}
			}
		} else {
			shared.loaded.p2 = false
			shared.avatar.p2 = canvas.NewImageFromResource(a)
		}

		if round.Turn == 2 {
			frame = canvas.NewImageFromResource(t)
		} else {
			frame = canvas.NewImageFromResource(f)
		}
	} else {
		shared.avatar.p2 = canvas.NewImageFromImage(nil)
		frame = canvas.NewImageFromImage(nil)
	}

	name.Resize(fyne.NewSize(100, 25))
	name.Move(fyne.NewPos(667, 20))

	shared.avatar.p2.Resize(fyne.NewSize(74, 74))
	shared.avatar.p2.Move(fyne.NewPos(782, 50))

	frame.Resize(fyne.NewSize(78, 78))
	frame.Move(fyne.NewPos(780, 48))

	return container.NewWithoutLayout(name, shared.avatar.p2, frame)
}

// Player 3 name and avatar objects
//   - Pass a and f as avatar and its frame resource, shared avatar is set here if image exists
//   - Pass t for player's turn frame resource
func Player3_label(a, f, t fyne.Resource) fyne.CanvasObject {
	var name fyne.CanvasObject
	var frame fyne.CanvasObject
	if signals.In3 {
		if round.Turn == 3 {
			name = canvas.NewText(round.p3.name, color.RGBA{105, 90, 205, 210})
		} else {
			name = canvas.NewText(round.p3.name, color.White)
		}
	} else {
		name = canvas.NewText("", color.Transparent)
	}

	if a != nil && signals.In3 {
		if round.p3.url != "" {
			if !shared.loaded.p3 {
				shared.loaded.p3 = true
				if img, err := dreams.DownloadCanvas(round.p3.url, "P3"); err == nil {
					shared.avatar.p3 = &img
				} else {
					logger.Errorln("[Holdero] Player 3 avatar:", err)
					shared.avatar.p3 = canvas.NewImageFromResource(a)
				}
			}
		} else {
			shared.loaded.p3 = false
			shared.avatar.p3 = canvas.NewImageFromResource(a)
		}

		if round.Turn == 3 {
			frame = canvas.NewImageFromResource(t)
		} else {
			frame = canvas.NewImageFromResource(f)
		}
	} else {
		shared.avatar.p3 = canvas.NewImageFromImage(nil)
		frame = canvas.NewImageFromImage(nil)
	}

	name.Resize(fyne.NewSize(100, 25))
	name.Move(fyne.NewPos(889, 300))

	shared.avatar.p3.Resize(fyne.NewSize(74, 74))
	shared.avatar.p3.Move(fyne.NewPos(987, 327))

	frame.Resize(fyne.NewSize(78, 78))
	frame.Move(fyne.NewPos(985, 325))

	return container.NewWithoutLayout(name, shared.avatar.p3, frame)
}

// Player 4 name and avatar objects
//   - Pass a and f as avatar and its frame resource, shared avatar is set here if image exists
//   - Pass t for player's turn frame resource
func Player4_label(a, f, t fyne.Resource) fyne.CanvasObject {
	var name fyne.CanvasObject
	var frame fyne.CanvasObject
	if signals.In4 {
		if round.Turn == 4 {
			name = canvas.NewText(round.p4.name, color.RGBA{105, 90, 205, 210})
		} else {
			name = canvas.NewText(round.p4.name, color.White)
		}
	} else {
		name = canvas.NewText("", color.Transparent)
	}

	if a != nil && signals.In4 {
		if round.p4.url != "" {
			if !shared.loaded.p4 {
				shared.loaded.p4 = true
				if img, err := dreams.DownloadCanvas(round.p4.url, "P4"); err == nil {
					shared.avatar.p4 = &img
				} else {
					logger.Errorln("[Holdero] Player 4 avatar:", err)
					shared.avatar.p4 = canvas.NewImageFromResource(a)
				}
			}
		} else {
			shared.loaded.p4 = false
			shared.avatar.p4 = canvas.NewImageFromResource(a)
		}

		if round.Turn == 4 {
			frame = canvas.NewImageFromResource(t)
		} else {
			frame = canvas.NewImageFromResource(f)
		}
	} else {
		shared.avatar.p4 = canvas.NewImageFromImage(nil)
		frame = canvas.NewImageFromImage(nil)
	}

	name.Resize(fyne.NewSize(100, 25))
	name.Move(fyne.NewPos(688, 555))

	shared.avatar.p4.Resize(fyne.NewSize(74, 74))
	shared.avatar.p4.Move(fyne.NewPos(686, 480))

	frame.Resize(fyne.NewSize(78, 78))
	frame.Move(fyne.NewPos(684, 478))

	return container.NewWithoutLayout(name, shared.avatar.p4, frame)
}

// Player 5 name and avatar objects
//   - Pass a and f as avatar and its frame resource, shared avatar is set here if image exists
//   - Pass t for player's turn frame resource
func Player5_label(a, f, t fyne.Resource) fyne.CanvasObject {
	var name fyne.CanvasObject
	var frame fyne.CanvasObject
	if signals.In5 {
		if round.Turn == 5 {
			name = canvas.NewText(round.p5.name, color.RGBA{105, 90, 205, 210})
		} else {
			name = canvas.NewText(round.p5.name, color.White)
		}
	} else {
		name = canvas.NewText("", color.Transparent)
	}

	if a != nil && signals.In5 {
		if round.p5.url != "" {
			if !shared.loaded.p5 {
				shared.loaded.p5 = true
				if img, err := dreams.DownloadCanvas(round.p5.url, "P5"); err == nil {
					shared.avatar.p5 = &img
				} else {
					logger.Errorln("[Holdero] Player 5 avatar:", err)
					shared.avatar.p5 = canvas.NewImageFromResource(a)
				}
			}
		} else {
			shared.loaded.p5 = false
			shared.avatar.p5 = canvas.NewImageFromResource(a)
		}

		if round.Turn == 5 {
			frame = canvas.NewImageFromResource(t)
		} else {
			frame = canvas.NewImageFromResource(f)
		}
	} else {
		shared.avatar.p5 = canvas.NewImageFromResource(nil)
		frame = canvas.NewImageFromResource(nil)
	}

	name.Resize(fyne.NewSize(100, 25))
	name.Move(fyne.NewPos(258, 555))

	shared.avatar.p5.Resize(fyne.NewSize(74, 74))
	shared.avatar.p5.Move(fyne.NewPos(257, 480))

	frame.Resize(fyne.NewSize(78, 78))
	frame.Move(fyne.NewPos(255, 478))

	return container.NewWithoutLayout(name, shared.avatar.p5, frame)
}

// Player 6 name and avatar objects
//   - Pass a and f as avatar and its frame resource, shared avatar is set here if image exists
//   - Pass t for player's turn frame resource
func Player6_label(a, f, t fyne.Resource) fyne.CanvasObject {
	var name fyne.CanvasObject
	var frame fyne.CanvasObject
	if signals.In6 {
		if round.Turn == 6 {
			name = canvas.NewText(round.p6.name, color.RGBA{105, 90, 205, 210})
		} else {
			name = canvas.NewText(round.p6.name, color.White)
		}
	} else {
		name = canvas.NewText("", color.Transparent)
	}

	if a != nil && signals.In6 {
		if round.p6.url != "" {
			if !shared.loaded.p6 {
				shared.loaded.p6 = true
				if img, err := dreams.DownloadCanvas(round.p6.url, "P6"); err == nil {
					shared.avatar.p6 = &img
				} else {
					logger.Errorln("[Holdero] Player 6 avatar:", err)
					shared.avatar.p6 = canvas.NewImageFromResource(a)
				}
			}
		} else {
			shared.loaded.p6 = false
			shared.avatar.p6 = canvas.NewImageFromResource(a)
		}

		if round.Turn == 6 {
			frame = canvas.NewImageFromResource(t)
		} else {
			frame = canvas.NewImageFromResource(f)
		}
	} else {
		shared.avatar.p6 = canvas.NewImageFromResource(nil)
		frame = canvas.NewImageFromResource(nil)
	}

	name.Resize(fyne.NewSize(100, 27))
	name.Move(fyne.NewPos(56, 267))

	shared.avatar.p6.Resize(fyne.NewSize(74, 74))
	shared.avatar.p6.Move(fyne.NewPos(55, 193))

	frame.Resize(fyne.NewSize(78, 78))
	frame.Move(fyne.NewPos(53, 191))

	return container.NewWithoutLayout(name, shared.avatar.p6, frame)
}

// Set Holdero table image resource
func HolderoTable(img fyne.Resource) fyne.CanvasObject {
	table_image := canvas.NewImageFromResource(img)
	table_image.Resize(fyne.NewSize(1100, 600))
	table_image.Move(fyne.NewPos(5, 0))

	return table_image
}

// Holdero object buffer when action triggered
func ActionBuffer() {
	table.sit.Hide()
	table.leave.Hide()
	table.deal.Hide()
	table.bet.Hide()
	table.check.Hide()
	table.betEntry.Hide()
	table.warning.Hide()
	round.display.results = ""
	signals.clicked = true
	signals.height = rpc.Wallet.Height
}

// Checking for current player names at connected Holdero table
//   - If name exists, prompt user to select new name
func checkNames(seats string) bool {
	if round.ID == 1 {
		return true
	}

	err := "[Holdero] Name already used"

	switch seats {
	case "2":
		if menu.Username == round.p1.name {
			logger.Warnln(err)
			return false
		}
		return true
	case "3":
		if menu.Username == round.p1.name || menu.Username == round.p2.name || menu.Username == round.p3.name {
			logger.Warnln(err)
			return false
		}
		return true
	case "4":
		if menu.Username == round.p1.name || menu.Username == round.p2.name || menu.Username == round.p3.name || menu.Username == round.p4.name {
			logger.Warnln(err)
			return false
		}
		return true
	case "5":
		if menu.Username == round.p1.name || menu.Username == round.p2.name || menu.Username == round.p3.name || menu.Username == round.p4.name || menu.Username == round.p5.name {
			logger.Warnln(err)
			return false
		}
		return true
	case "6":
		if menu.Username == round.p1.name || menu.Username == round.p2.name || menu.Username == round.p3.name || menu.Username == round.p4.name || menu.Username == round.p5.name || menu.Username == round.p6.name {
			logger.Warnln(err)
			return false
		}
		return true
	default:
		return false
	}
}

// Holdero player sit down button to join current table
func SitButton() fyne.Widget {
	table.sit = widget.NewButton("Sit Down", func() {
		if menu.Username != "" {
			if checkNames(round.display.seats) {
				SitDown(menu.Username, Settings.avatar.url)
				ActionBuffer()
			}
		} else {
			logger.Warnln("[Holdero] Pick a name")
		}
	})

	table.sit.Hide()

	return table.sit
}

// Holdero player leave button to leave current table seat
func LeaveButton() fyne.Widget {
	table.leave = widget.NewButton("Leave", func() {
		Leave()
		ActionBuffer()
	})

	table.leave.Hide()

	return table.leave
}

// Holdero player deal hand button
func DealHandButton() fyne.Widget {
	table.deal = widget.NewButton("Deal Hand", func() {
		if tx := DealHand(); tx != "" {
			ActionBuffer()
		}
	})

	table.deal.Hide()

	return table.deal
}

// Holdero bet entry amount
//   - Setting the initial value based on if PlacedBet, Wager and Ante
//   - If entry invalid, set to min bet value
func BetAmount() fyne.CanvasObject {
	table.betEntry = dwidget.NewDeroEntry("", 0.1, 1)
	table.betEntry.Enable()
	if table.betEntry.Text == "" {
		table.betEntry.SetText("0.0")
	}
	table.betEntry.Validator = validation.NewRegexp(`^\d{1,}\.\d{1,5}$|^[^0.]\d{0,}$`, "Int or float required")
	table.betEntry.OnChanged = func(s string) {
		if f, err := strconv.ParseFloat(s, 64); err == nil {
			if signals.placedBet {
				table.betEntry.SetText(strconv.FormatFloat(float64(round.Raised)/100000, 'f', int(table.betEntry.Decimal), 64))
				if table.betEntry.Validate() != nil {
					table.betEntry.SetText(strconv.FormatFloat(float64(round.Raised)/100000, 'f', int(table.betEntry.Decimal), 64))
				}
			} else {

				if round.Wager > 0 {
					table.check.SetText("Fold")
					if round.Raised > 0 {
						table.bet.SetText("Call")
						if signals.placedBet {
							table.betEntry.SetText(strconv.FormatFloat(float64(round.Raised)/100000, 'f', int(table.betEntry.Decimal), 64))
						} else {
							table.betEntry.SetText(strconv.FormatFloat(float64(round.Wager)/100000, 'f', int(table.betEntry.Decimal), 64))
						}
						if table.betEntry.Validate() != nil {
							if signals.placedBet {
								table.betEntry.SetText(strconv.FormatFloat(float64(round.Raised)/100000, 'f', int(table.betEntry.Decimal), 64))
							} else {
								table.betEntry.SetText(strconv.FormatFloat(float64(round.Wager)/100000, 'f', int(table.betEntry.Decimal), 64))
							}
						}
					} else {
						if f > float64(round.Wager)/100000 {
							table.bet.SetText("Raise")
						} else {
							table.bet.SetText("Call")
						}

						if f < float64(round.Wager)/100000 {
							table.betEntry.SetText(strconv.FormatFloat(float64(round.Wager)/100000, 'f', int(table.betEntry.Decimal), 64))
						}

						if table.betEntry.Validate() != nil {
							float := f * 100000
							if uint64(float)%10000 == 0 {
								table.betEntry.SetText(strconv.FormatFloat(roundFloat(f, 1), 'f', int(table.betEntry.Decimal), 64))
							} else if table.betEntry.Validate() != nil {
								table.betEntry.SetText(strconv.FormatFloat(roundFloat(f, 1), 'f', int(table.betEntry.Decimal), 64))
							}
						}
					}
				} else {
					table.bet.SetText("Bet")
					table.check.SetText("Check")
					if rpc.Daemon.IsConnected() {
						float := f * 100000
						if uint64(float)%10000 == 0 {
							table.betEntry.SetText(strconv.FormatFloat(roundFloat(f, 1), 'f', int(table.betEntry.Decimal), 64))
						} else if table.betEntry.Validate() != nil {
							table.betEntry.SetText(strconv.FormatFloat(roundFloat(f, 1), 'f', int(table.betEntry.Decimal), 64))
						}

						if round.Ante > 0 {
							if f < float64(round.Ante)/100000 {
								table.betEntry.SetText(strconv.FormatFloat(float64(round.Ante)/100000, 'f', int(table.betEntry.Decimal), 64))
							}

							if table.betEntry.Validate() != nil {
								table.betEntry.SetText(strconv.FormatFloat(float64(round.Ante)/100000, 'f', int(table.betEntry.Decimal), 64))
							}

						} else {
							if f < float64(round.BB)/100000 {
								table.betEntry.SetText(strconv.FormatFloat(float64(round.BB)/100000, 'f', int(table.betEntry.Decimal), 64))
							}

							if table.betEntry.Validate() != nil {
								table.betEntry.SetText(strconv.FormatFloat(float64(round.BB)/100000, 'f', int(table.betEntry.Decimal), 64))
							}
						}
					}
				}
			}
		} else {
			logger.Errorln("[BetAmount]", err)
			if round.Ante == 0 {
				table.betEntry.SetText(strconv.FormatFloat(float64(round.BB)/100000, 'f', int(table.betEntry.Decimal), 64))
			} else {
				table.betEntry.SetText(strconv.FormatFloat(float64(round.Ante)/100000, 'f', int(table.betEntry.Decimal), 64))
			}
		}
	}

	amt_box := container.NewHScroll(table.betEntry)
	amt_box.SetMinSize(fyne.NewSize(100, 40))
	table.betEntry.Hide()

	return amt_box

}

// Round float val to precision
func roundFloat(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}

// Holdero place bet button
//   - Input from Table.BetEntry
func BetButton() fyne.Widget {
	table.bet = widget.NewButton("Bet", func() {
		if table.betEntry.Validate() == nil {
			if tx := Bet(table.betEntry.Text); tx != "" {
				signals.bet = true
				ActionBuffer()
			}
		}
	})

	table.bet.Hide()

	return table.bet
}

// Holdero check and fold button
func CheckButton() fyne.Widget {
	table.check = widget.NewButton("Check", func() {
		if tx := Check(); tx != "" {
			signals.bet = true
			ActionBuffer()
		}
	})

	table.check.Hide()

	return table.check
}

// Automated options object for Holdero
func AutoOptions(d *dreams.AppObject) fyne.CanvasObject {
	refresh := widget.NewButtonWithIcon("", fyne.Theme.Icon(fyne.CurrentApp().Settings().Theme(), "viewRefresh"), func() {
		if !rpc.Daemon.IsConnected() || !rpc.Wallet.IsConnected() {
			dialog.NewInformation("Not connected", "You are not connected to daemon or wallet", d.Window).Show()
			return
		}

		if !signals.contract {
			dialog.NewInformation("Not connected", "You are not connected to a Holdero SC", d.Window).Show()
			return
		}
		fetchHolderoSC()
	})

	cf := widget.NewCheck("Auto Check/Fold", func(b bool) {
		if b {
			Settings.auto.check = true
		} else {
			Settings.auto.check = false
		}
	})

	deal := widget.NewCheck("Auto Deal", func(b bool) {
		if b {
			Settings.auto.deal = true
		} else {
			Settings.auto.deal = false
		}
	})

	Settings.tools = widget.NewButton("Tools", func() {
		go holderoTools(deal, cf, Settings.tools)
	})

	DisableHolderoTools()

	return container.NewVBox(container.NewHBox(refresh, layout.NewSpacer()), deal, cf, Settings.tools)
}

// Holdero warning label displayed when player is risking being timed out
func TimeOutWarning() *fyne.Container {
	rect := canvas.NewRectangle(color.RGBA{0, 0, 0, 210})
	msg := canvas.NewText("Make your move, or you will be Timed Out", color.RGBA{240, 0, 0, 240})
	msg.TextSize = 15

	table.warning = container.NewStack(rect, msg)

	table.warning.Hide()

	return container.NewVBox(layout.NewSpacer(), table.warning)
}

// Set default params for auto bet functions
func autoBetDefault() {
	Odds.Bot.Risk[2] = 21
	Odds.Bot.Risk[1] = 9
	Odds.Bot.Risk[0] = 3
	Odds.Bot.Bet[2] = 6
	Odds.Bot.Bet[1] = 3
	Odds.Bot.Bet[0] = 1
	Odds.Bot.Luck = 0
	Odds.Bot.Slow = 4
	Odds.Bot.Aggr = 1
	Odds.Bot.Max = 10
	Odds.Bot.Random[0] = 0
	Odds.Bot.Random[1] = 0
}

// Setting current auto bet random option when menu opened
func setRandomOpts(opts *widget.RadioGroup) {
	if Odds.Bot.Random[0] == 0 {
		opts.Disable()
	} else {
		switch Odds.Bot.Random[1] {
		case 1:
			opts.SetSelected("Risk")
		case 2:
			opts.SetSelected("Bet")
		case 3:
			opts.SetSelected("Both")
		default:
			opts.SetSelected("")
		}
	}
}

// dReam Tools menu for Holdero
//   - deal check and button widgets are passed when setting auto objects for control
func holderoTools(deal, check *widget.Check, button *widget.Button) {
	button.Hide()
	bm := fyne.CurrentApp().NewWindow("Holdero Tools")
	bm.Resize(fyne.NewSize(330, 700))
	bm.SetFixedSize(true)
	bm.SetIcon(bundle.ResourceDReamsIconAltPng)
	bm.SetCloseIntercept(func() {
		button.Show()
		bm.Close()
	})

	stats = readSavedStats()
	config_opts := []string{}
	for i := range stats.Bots {
		config_opts = append(config_opts, stats.Bots[i].Name)
	}

	entry := widget.NewSelectEntry(config_opts)
	entry.SetPlaceHolder("Default")
	entry.SetText(Odds.Bot.Name)

	curr := " Dero"
	max_bet := float64(100)
	if round.asset {
		curr = " Tokens"
		max_bet = 2500
	}

	mb_label := widget.NewLabel("Max Bet: " + fmt.Sprintf("%.0f", Odds.Bot.Max) + curr)
	mb_slider := widget.NewSlider(1, max_bet)
	mb_slider.SetValue(Odds.Bot.Max)
	mb_slider.OnChanged = func(f float64) {
		go func() {
			min := float64(MinBet()) / 100000
			if min == 0 {
				min = 0.1
			}

			if f < (min*Odds.Bot.Bet[2])*Odds.Bot.Aggr {
				Odds.Bot.Max = (min*Odds.Bot.Bet[2])*Odds.Bot.Aggr + 3
				mb_slider.SetValue(Odds.Bot.Max)
				mb_label.SetText("Max Bet: " + fmt.Sprintf("%.0f", Odds.Bot.Max) + curr)
			} else {
				Odds.Bot.Max = f
				mb_label.SetText("Max Bet: " + fmt.Sprintf("%.0f", f) + curr)
			}
		}()
	}

	rh_label := widget.NewLabel("Risk High: " + fmt.Sprintf("%.0f", Odds.Bot.Risk[2]) + "%")
	rh_slider := widget.NewSlider(1, 90)
	rh_slider.SetValue(Odds.Bot.Risk[2])
	rh_slider.OnChanged = func(f float64) {
		go func() {
			if f < Odds.Bot.Risk[1] {
				Odds.Bot.Risk[2] = Odds.Bot.Risk[1] + 1
				rh_slider.SetValue(Odds.Bot.Risk[2])
			} else {
				Odds.Bot.Risk[2] = f
			}

			rh_label.SetText("Risk High: " + fmt.Sprintf("%.0f", Odds.Bot.Risk[2]) + "%")
		}()
	}

	rm_label := widget.NewLabel("Risk Medium: " + fmt.Sprintf("%.0f", Odds.Bot.Risk[1]) + "%")
	rm_slider := widget.NewSlider(1, 89)
	rm_slider.SetValue(Odds.Bot.Risk[1])
	rm_slider.OnChanged = func(f float64) {
		go func() {
			Odds.Bot.Risk[1] = f
			if f <= Odds.Bot.Risk[0] {
				Odds.Bot.Risk[1] = Odds.Bot.Risk[0] + 1
				rm_slider.SetValue(Odds.Bot.Risk[1])
			}

			if f >= Odds.Bot.Risk[2] {
				Odds.Bot.Risk[2] = f + 1
				rh_slider.SetValue(Odds.Bot.Risk[2])
			}

			rm_label.SetText("Risk Medium: " + fmt.Sprintf("%.0f", Odds.Bot.Risk[1]) + "%")
		}()
	}

	rl_label := widget.NewLabel("Risk Low: " + fmt.Sprintf("%.0f", Odds.Bot.Risk[0]) + "%")
	rl_slider := widget.NewSlider(1, 88)
	rl_slider.SetValue(Odds.Bot.Risk[0])
	rl_slider.OnChanged = func(f float64) {
		go func() {
			if Odds.Bot.Risk[1] <= f {
				rm_slider.SetValue(f + 1)
			}

			Odds.Bot.Risk[0] = f
			rl_label.SetText("Risk Low: " + fmt.Sprintf("%.0f", Odds.Bot.Risk[0]) + "%")
		}()
	}

	bh_label := widget.NewLabel("Bet High: " + fmt.Sprintf("%.0f", Odds.Bot.Bet[2]) + "x")
	bh_slider := widget.NewSlider(1, 100)
	bh_slider.SetValue(Odds.Bot.Bet[2])
	bh_slider.OnChanged = func(f float64) {
		go func() {
			if f < Odds.Bot.Bet[1] {
				Odds.Bot.Bet[2] = Odds.Bot.Bet[1] + 1
				bh_slider.SetValue(Odds.Bot.Bet[2])
			} else {
				Odds.Bot.Bet[2] = f
			}

			min := float64(MinBet()) / 100000
			if min == 0 {
				min = 0.1
			}

			if Odds.Bot.Max < (min*Odds.Bot.Bet[2])*Odds.Bot.Aggr {
				Odds.Bot.Max = (min * Odds.Bot.Bet[2]) * Odds.Bot.Aggr
				mb_slider.SetValue(Odds.Bot.Max)
			}

			bh_label.SetText("Bet High: " + fmt.Sprintf("%.0f", Odds.Bot.Bet[2]) + "x")
		}()
	}

	bm_label := widget.NewLabel("Bet Medium: " + fmt.Sprintf("%.0f", Odds.Bot.Bet[1]) + "x")
	bm_slider := widget.NewSlider(1, 99)
	bm_slider.SetValue(Odds.Bot.Bet[1])
	bm_slider.OnChanged = func(f float64) {
		go func() {
			Odds.Bot.Bet[1] = f
			if f <= Odds.Bot.Bet[0] {
				Odds.Bot.Bet[1] = Odds.Bot.Bet[0] + 1
				bm_slider.SetValue(Odds.Bot.Bet[1])
			}

			if f >= Odds.Bot.Bet[2] {
				Odds.Bot.Bet[2] = f + 1
				bh_slider.SetValue(Odds.Bot.Bet[2])
			}

			bm_label.SetText("Bet Medium: " + fmt.Sprintf("%.0f", Odds.Bot.Bet[1]) + "x")
		}()
	}

	bl_label := widget.NewLabel("Bet Low: " + fmt.Sprintf("%.0f", Odds.Bot.Bet[0]) + "x")
	bl_slider := widget.NewSlider(1, 98)
	bl_slider.SetValue(Odds.Bot.Bet[0])
	bl_slider.OnChanged = func(f float64) {
		go func() {
			if Odds.Bot.Bet[1] <= f {
				bm_slider.SetValue(f + 1)
			}

			Odds.Bot.Bet[0] = f
			bl_label.SetText("Bet Low: " + fmt.Sprintf("%.0f", Odds.Bot.Bet[0]) + "x")
		}()
	}

	luck_label := widget.NewLabel("Luck: " + fmt.Sprintf("%.2f", Odds.Bot.Luck))
	luck_slider := widget.NewSlider(0, 10)
	luck_slider.Step = 0.25
	luck_slider.SetValue(Odds.Bot.Luck)
	luck_slider.OnChanged = func(f float64) {
		go func() {
			Odds.Bot.Luck = f
			luck_label.SetText("Luck: " + fmt.Sprintf("%.2f", f))
		}()
	}

	random_label := widget.NewLabel("Randomize: Off")
	if Odds.Bot.Random[0] == 0 {
		random_label.SetText("Randomize: Off")
	} else {
		random_label.SetText("Randomize: " + fmt.Sprintf("%.2f", Odds.Bot.Random[0]))
	}

	random_opts := widget.NewRadioGroup([]string{"Risk", "Bet", "Both"}, func(s string) {
		switch s {
		case "Risk":
			Odds.Bot.Random[1] = 1
		case "Bet":
			Odds.Bot.Random[1] = 2
		case "Both":
			Odds.Bot.Random[1] = 3
		default:
			Odds.Bot.Random[1] = 0
		}
	})

	setRandomOpts(random_opts)

	random_slider := widget.NewSlider(0, 10)
	random_slider.Step = 0.25
	random_slider.SetValue(Odds.Bot.Random[0])
	random_slider.OnChanged = func(f float64) {
		go func() {
			Odds.Bot.Random[0] = f
			if f >= 0.5 {
				random_label.SetText("Randomize: " + fmt.Sprintf("%.2f", f))
				random_opts.Enable()
			} else {
				Odds.Bot.Random[0] = 0
				Odds.Bot.Random[1] = 0
				random_label.SetText("Randomize: Off")
				random_opts.SetSelected("")
				random_opts.Disable()
			}
		}()
	}

	slow_label := widget.NewLabel("Slowplay: " + fmt.Sprintf("%.0f", Odds.Bot.Slow))
	slow_slider := widget.NewSlider(1, 5)
	slow_slider.SetValue(Odds.Bot.Slow)
	slow_slider.OnChanged = func(f float64) {
		go func() {
			Odds.Bot.Slow = f
			slow_label.SetText("Slowplay: " + fmt.Sprintf("%.0f", f))
		}()
	}

	aggr_label := widget.NewLabel("Aggression: " + fmt.Sprintf("%.0f", Odds.Bot.Aggr))
	aggr_slider := widget.NewSlider(1, 5)
	aggr_slider.SetValue(Odds.Bot.Aggr)
	aggr_slider.OnChanged = func(f float64) {
		go func() {
			Odds.Bot.Aggr = f
			min := float64(MinBet()) / 100000
			if min == 0 {
				min = 0.1
			}

			if Odds.Bot.Max < (min*Odds.Bot.Bet[2])*Odds.Bot.Aggr {
				Odds.Bot.Max = (min * Odds.Bot.Bet[2]) * Odds.Bot.Aggr
				mb_slider.SetValue(Odds.Bot.Max)
			}

			aggr_label.SetText("Aggression: " + fmt.Sprintf("%.0f", f))
		}()
	}

	rem := widget.NewButtonWithIcon("", fyne.Theme.Icon(fyne.CurrentApp().Settings().Theme(), "delete"), func() {
		if entry.Text != "" {
			var new []Bot_config
			for i := range stats.Bots {
				if stats.Bots[i].Name == entry.Text {
					logger.Println("[Holdero] Deleting bot config")
					if i > 0 {
						new = append(stats.Bots[0:i], stats.Bots[i+1:]...)
						config_opts = append(config_opts[0:i], config_opts[i+1:]...)
						break
					} else {
						if len(config_opts) < 2 {
							new = nil
							config_opts = []string{}
						} else {
							new = stats.Bots[1:]
							config_opts = append(config_opts[1:2], config_opts[2:]...)
						}
						break
					}
				}
			}

			stats.Bots = new
			WriteHolderoStats(stats)
			entry.SetOptions(config_opts)
			entry.SetText("")
		}
	})

	reset := widget.NewButtonWithIcon("", fyne.Theme.Icon(fyne.CurrentApp().Settings().Theme(), "viewRefresh"), func() {
		autoBetDefault()
		mb_slider.SetValue(Odds.Bot.Max)
		rh_slider.SetValue(Odds.Bot.Risk[2])
		rm_slider.SetValue(Odds.Bot.Risk[1])
		rl_slider.SetValue(Odds.Bot.Risk[0])
		bh_slider.SetValue(Odds.Bot.Bet[2])
		bm_slider.SetValue(Odds.Bot.Bet[1])
		bl_slider.SetValue(Odds.Bot.Bet[0])
		luck_slider.SetValue(Odds.Bot.Luck)
		random_slider.SetValue(Odds.Bot.Random[0])
		slow_slider.SetValue(Odds.Bot.Slow)
		aggr_slider.SetValue(Odds.Bot.Aggr)
		random_opts.SetSelected("")
		entry.SetText("")
		Odds.Bot.Name = ""
	})

	save := widget.NewButton("Save", func() {
		if entry.Text != "" {
			var ex bool
			for i := range stats.Bots {
				if entry.Text == stats.Bots[i].Name {
					ex = true
					logger.Warnln("[Holdero] Bot config name exists")
				}
			}

			if !ex {
				stats.Bots = append(stats.Bots, Odds.Bot)
				if WriteHolderoStats(stats) {
					config_opts = append(config_opts, entry.Text)
					entry.SetOptions(config_opts)
					logger.Println("[Holdero] Saved bot config")
				}
			}
		}
	})

	entry.OnChanged = func(s string) {
		if s != "" {
			Odds.Bot.Name = entry.Text
			for i := range config_opts {
				if s == config_opts[i] {
					for r := range stats.Bots {
						if stats.Bots[r].Name == config_opts[i] {
							SetBotConfig(stats.Bots[r])
							mb_slider.SetValue(Odds.Bot.Max)
							rh_slider.SetValue(Odds.Bot.Risk[2])
							rm_slider.SetValue(Odds.Bot.Risk[1])
							rl_slider.SetValue(Odds.Bot.Risk[0])
							bh_slider.SetValue(Odds.Bot.Bet[2])
							bm_slider.SetValue(Odds.Bot.Bet[1])
							bl_slider.SetValue(Odds.Bot.Bet[0])
							luck_slider.SetValue(Odds.Bot.Luck)
							random_slider.SetValue(Odds.Bot.Random[0])
							slow_slider.SetValue(Odds.Bot.Slow)
							aggr_slider.SetValue(Odds.Bot.Aggr)
							setRandomOpts(random_opts)
						}
					}
				}
			}
		}
	}

	enable := widget.NewCheck("Auto Bet Enabled", func(b bool) {
		if b {
			Odds.Start()
			if check.Checked {
				check.SetChecked(false)
			}
			check.Disable()
			deal.SetChecked(true)
		} else {
			Odds.Stop()
			check.Enable()
			if deal.Checked {
				deal.SetChecked(false)
			}
		}
	})

	if Odds.IsRunning() {
		enable.SetChecked(true)
	}

	config_buttons := container.NewBorder(nil, nil, nil, container.NewHBox(reset, rem), save)

	top_box := container.NewAdaptiveGrid(2,
		container.NewVBox(
			luck_label,
			luck_slider,
			slow_label,
			slow_slider,
			aggr_label,
			aggr_slider,
			mb_label,
			mb_slider,
			layout.NewSpacer(),
			enable),

		container.NewVBox(
			random_label,
			random_slider,
			random_opts,
			layout.NewSpacer(),
			entry,
			config_buttons))

	Odds.Label = widget.NewLabel("")
	Odds.Label.Wrapping = fyne.TextWrapWord
	scroll := container.NewVScroll(Odds.Label)
	odds_button := widget.NewButton("Odds", func() {
		odds, future := MakeOdds()
		BetLogic(odds, future, false)
	})

	r_box := container.NewVBox(
		rh_label,
		rh_slider,
		rm_label,
		rm_slider,
		rl_label,
		rl_slider)

	b_box := container.NewVBox(
		bh_label,
		bh_slider,
		bm_label,
		bm_slider,
		bl_label,
		bl_slider)

	bet_bot := container.NewVBox(
		r_box,
		layout.NewSpacer(),
		b_box,
		layout.NewSpacer(),
		top_box)

	odds_box := container.NewVBox(layout.NewSpacer(), odds_button)
	max := container.NewStack(scroll, odds_box)

	stats_label := widget.NewLabel("")

	tabs := container.NewAppTabs(
		container.NewTabItem("Bot", container.NewBorder(nil, nil, nil, nil, bet_bot)),
		container.NewTabItem("Odds", max),
		container.NewTabItem("Stats", stats_label),
	)

	tabs.OnSelected = func(ti *container.TabItem) {
		switch ti.Text {
		case "Stats":
			stats_label.SetText("Total Player Stats\n\nWins: " + strconv.Itoa(stats.Player.Win) + "\n\nLost: " + strconv.Itoa(stats.Player.Lost) +
				"\n\nFolded: " + strconv.Itoa(stats.Player.Fold) + "\n\nPush: " + strconv.Itoa(stats.Player.Push) +
				"\n\nWagered: " + fmt.Sprintf("%.5f", stats.Player.Wagered) + "\n\nEarnings: " + fmt.Sprintf("%.5f", stats.Player.Earnings))

			if Odds.Bot.Name != "" {
				for i := range stats.Bots {
					if Odds.Bot.Name == stats.Bots[i].Name {
						stats_label.SetText(stats_label.Text + "\n\n\nBot Stats\n\nBot: " + Odds.Bot.Name + "\n\nWins: " + strconv.Itoa(stats.Bots[i].Stats.Win) +
							"\n\nLost: " + strconv.Itoa(stats.Bots[i].Stats.Lost) + "\n\nFolded: " + strconv.Itoa(stats.Bots[i].Stats.Fold) + "\n\nPush: " + strconv.Itoa(stats.Bots[i].Stats.Push) +
							"\n\nWagered: " + fmt.Sprintf("%.5f", stats.Bots[i].Stats.Wagered) + "\n\nEarnings: " + fmt.Sprintf("%.5f", stats.Bots[i].Stats.Earnings))
					}
				}
			}
		}
	}

	go func() {
		for rpc.Wallet.IsConnected() {
			time.Sleep(1 * time.Second)
		}

		button.Show()
		bm.Close()
	}()

	var err error
	var img image.Image
	var rast *canvas.Raster
	if img, _, err = image.Decode(bytes.NewReader(menu.Theme.Img.Resource.Content())); err != nil {
		if img, _, err = image.Decode(bytes.NewReader(menu.DefaultThemeResource().StaticContent)); err != nil {
			logger.Warnf("[holderoTools] Fallback %s\n", err)
			source := image.Rect(2, 2, 4, 4)

			rast = canvas.NewRasterFromImage(source)
		} else {
			rast = canvas.NewRasterFromImage(img)
		}
	} else {
		rast = canvas.NewRasterFromImage(img)
	}

	bm.SetContent(
		container.NewStack(
			rast,
			bundle.Alpha180,
			tabs))
	bm.Show()
}

func DisableHolderoTools() {
	Odds.Enabled = false
	Settings.tools.Hide()
	if len(Settings.backs.Select.Options) > 2 || len(Settings.faces.Select.Options) > 2 {
		cards := false
		for _, f := range Settings.faces.Select.Options {
			asset := strings.Trim(f, "0123456789")
			switch asset {
			case "AZYPC":
				cards = true
			case "SIXPC":
				cards = true
			default:

			}

			if cards {
				break
			}
		}

		if !cards {
			for _, b := range Settings.backs.Select.Options {
				asset := strings.Trim(b, "0123456789")
				switch asset {
				case "AZYPCB":
					cards = true
				case "SIXPCB":
					cards = true
				default:

				}

				if cards {
					break
				}
			}
		}

		if cards {
			Odds.Enabled = true
			Settings.tools.Show()
			if !dreams.FileExists("config/stats.json", "Holdero") {
				WriteHolderoStats(stats)
				logger.Println("[Holdero] Created stats.json")
			} else {
				stats = readSavedStats()
			}
		}
	}
}

// Reading saved Holdero stats from config file
func readSavedStats() (saved Player_stats) {
	file, err := os.ReadFile("config/stats.json")

	if err != nil {
		logger.Errorln("[readSavedStats]", err)
		return
	}

	err = json.Unmarshal(file, &saved)
	if err != nil {
		logger.Errorln("[readSavedStats]", err)
		return
	}

	return
}
