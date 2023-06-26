package holdero

import (
	"math/rand"
	"time"

	dreams "github.com/dReam-dApps/dReams"
	"github.com/dReam-dApps/dReams/bundle"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
)

type cards struct {
	Hole1 fyne.CanvasObject
	Hole2 fyne.CanvasObject
	Flop1 fyne.CanvasObject
	Flop2 fyne.CanvasObject
	Flop3 fyne.CanvasObject
	Turn  fyne.CanvasObject
	River fyne.CanvasObject

	P1a fyne.CanvasObject
	P1b fyne.CanvasObject

	P2a fyne.CanvasObject
	P2b fyne.CanvasObject

	P3a fyne.CanvasObject
	P3b fyne.CanvasObject

	P4a fyne.CanvasObject
	P4b fyne.CanvasObject

	P5a fyne.CanvasObject
	P5b fyne.CanvasObject

	P6a fyne.CanvasObject
	P6b fyne.CanvasObject

	Layout *fyne.Container
}

var Cards cards

// Set player hole card one image
//   - w and h of main window for resize
func Hole_1(c int, w, h float32) fyne.CanvasObject {
	Cards.Hole1 = DisplayCard(c)
	Cards.Hole1.Resize(fyne.NewSize(165, 225))
	Cards.Hole1.Move(fyne.NewPos(w-335, h-335))

	return Cards.Hole1
}

// Set player hole card two image
//   - w and h of main window for resize
func Hole_2(c int, w, h float32) fyne.CanvasObject {
	Cards.Hole2 = DisplayCard(c)
	Cards.Hole2.Resize(fyne.NewSize(165, 225))
	Cards.Hole2.Move(fyne.NewPos(w-275, h-335))

	return Cards.Hole2
}

// Set first flop card image
func Flop_1(c int) fyne.CanvasObject {
	size := fyne.NewSize(110, 150)
	pos := fyne.NewPos(260, 203)
	Cards.Flop1 = DisplayCard(c)
	Cards.Flop1.Resize(size)
	Cards.Flop1.Move(pos)

	for _, i := range Round.Winning_hand {
		if c == i {
			highlight := canvas.NewRectangle(bundle.Highlight)
			highlight.Resize(size)
			highlight.Move(pos)

			return container.NewWithoutLayout(Cards.Flop1, highlight)
		}
	}

	return Cards.Flop1
}

// Set second flop card image
func Flop_2(c int) fyne.CanvasObject {
	size := fyne.NewSize(110, 150)
	pos := fyne.NewPos(380, 203)
	Cards.Flop2 = DisplayCard(c)
	Cards.Flop2.Resize(size)
	Cards.Flop2.Move(pos)

	for _, i := range Round.Winning_hand {
		if c == i {
			highlight := canvas.NewRectangle(bundle.Highlight)
			highlight.Resize(size)
			highlight.Move(pos)

			return container.NewWithoutLayout(Cards.Flop2, highlight)
		}
	}

	return Cards.Flop2
}

// Set third flop card image
func Flop_3(c int) fyne.CanvasObject {
	size := fyne.NewSize(110, 150)
	pos := fyne.NewPos(500, 203)
	Cards.Flop3 = DisplayCard(c)
	Cards.Flop3.Resize(size)
	Cards.Flop3.Move(pos)

	for _, i := range Round.Winning_hand {
		if c == i {
			highlight := canvas.NewRectangle(bundle.Highlight)
			highlight.Resize(size)
			highlight.Move(pos)

			return container.NewWithoutLayout(Cards.Flop3, highlight)
		}
	}

	return Cards.Flop3
}

// Set turn card image
func Turn(c int) fyne.CanvasObject {
	size := fyne.NewSize(110, 150)
	pos := fyne.NewPos(620, 203)
	Cards.Turn = DisplayCard(c)
	Cards.Turn.Resize(size)
	Cards.Turn.Move(pos)

	for _, i := range Round.Winning_hand {
		if c == i {
			highlight := canvas.NewRectangle(bundle.Highlight)
			highlight.Resize(size)
			highlight.Move(pos)

			return container.NewWithoutLayout(Cards.Turn, highlight)
		}
	}

	return Cards.Turn
}

// Set river card image
func River(c int) fyne.CanvasObject {
	size := fyne.NewSize(110, 150)
	pos := fyne.NewPos(740, 203)
	Cards.River = DisplayCard(c)
	Cards.River.Resize(size)
	Cards.River.Move(pos)

	for _, i := range Round.Winning_hand {
		if c == i {
			highlight := canvas.NewRectangle(bundle.Highlight)
			highlight.Resize(size)
			highlight.Move(pos)

			return container.NewWithoutLayout(Cards.River, highlight)
		}
	}

	return Cards.River
}

// Set first players table card one image
func P1_a(c int) fyne.CanvasObject {
	size := fyne.NewSize(110, 150)
	pos := fyne.NewPos(190, 25)
	Cards.P1a = DisplayCard(c)
	Cards.P1a.Resize(size)
	Cards.P1a.Move(pos)

	for _, i := range Round.Winning_hand {
		if c == i {
			highlight := canvas.NewRectangle(bundle.Highlight)
			highlight.Resize(size)
			highlight.Move(pos)

			return container.NewWithoutLayout(Cards.P1a, highlight)
		}
	}

	return Cards.P1a
}

// Set first players table card two image
func P1_b(c int) fyne.CanvasObject {
	size := fyne.NewSize(110, 150)
	pos := fyne.NewPos(242, 25)
	Cards.P1b = DisplayCard(c)
	Cards.P1b.Resize(size)
	Cards.P1b.Move(pos)

	for _, i := range Round.Winning_hand {
		if c == i {
			highlight := canvas.NewRectangle(bundle.Highlight)
			highlight.Resize(size)
			highlight.Move(pos)

			return container.NewWithoutLayout(Cards.P1b, highlight)
		}
	}

	return Cards.P1b
}

// Set second players table card one image
func P2_a(c int) fyne.CanvasObject {
	size := fyne.NewSize(110, 150)
	pos := fyne.NewPos(614, 25)
	Cards.P2a = DisplayCard(c)
	Cards.P2a.Resize(size)
	Cards.P2a.Move(pos)

	for _, i := range Round.Winning_hand {
		if c == i {
			highlight := canvas.NewRectangle(bundle.Highlight)
			highlight.Resize(size)
			highlight.Move(pos)

			return container.NewWithoutLayout(Cards.P2a, highlight)
		}
	}

	return Cards.P2a
}

// Set second players table card two image
func P2_b(c int) fyne.CanvasObject {
	size := fyne.NewSize(110, 150)
	pos := fyne.NewPos(666, 25)
	Cards.P2b = DisplayCard(c)
	Cards.P2b.Resize(size)
	Cards.P2b.Move(pos)

	for _, i := range Round.Winning_hand {
		if c == i {
			highlight := canvas.NewRectangle(bundle.Highlight)
			highlight.Resize(size)
			highlight.Move(pos)

			return container.NewWithoutLayout(Cards.P2b, highlight)
		}
	}

	return Cards.P2b
}

// Set third players table card one image
func P3_a(c int) fyne.CanvasObject {
	size := fyne.NewSize(110, 150)
	pos := fyne.NewPos(883, 129)
	Cards.P3a = DisplayCard(c)
	Cards.P3a.Resize(size)
	Cards.P3a.Move(pos)

	for _, i := range Round.Winning_hand {
		if c == i {
			highlight := canvas.NewRectangle(bundle.Highlight)
			highlight.Resize(size)
			highlight.Move(pos)

			return container.NewWithoutLayout(Cards.P3a, highlight)
		}
	}

	return Cards.P3a

}

// Set third players table card two image
func P3_b(c int) fyne.CanvasObject {
	size := fyne.NewSize(110, 150)
	pos := fyne.NewPos(935, 129)
	Cards.P3b = DisplayCard(c)
	Cards.P3b.Resize(size)
	Cards.P3b.Move(pos)

	for _, i := range Round.Winning_hand {
		if c == i {
			highlight := canvas.NewRectangle(bundle.Highlight)
			highlight.Resize(size)
			highlight.Move(pos)

			return container.NewWithoutLayout(Cards.P3b, highlight)
		}
	}

	return Cards.P3b
}

// Set fourth players table card one image
func P4_a(c int) fyne.CanvasObject {
	size := fyne.NewSize(110, 150)
	pos := fyne.NewPos(766, 383)
	Cards.P4a = DisplayCard(c)
	Cards.P4a.Resize(size)
	Cards.P4a.Move(pos)

	for _, i := range Round.Winning_hand {
		if c == i {
			highlight := canvas.NewRectangle(bundle.Highlight)
			highlight.Resize(size)
			highlight.Move(pos)

			return container.NewWithoutLayout(Cards.P4a, highlight)
		}
	}

	return Cards.P4a
}

// Set fourth players table card two image
func P4_b(c int) fyne.CanvasObject {
	size := fyne.NewSize(110, 150)
	pos := fyne.NewPos(818, 383)
	Cards.P4b = DisplayCard(c)
	Cards.P4b.Resize(size)
	Cards.P4b.Move(pos)

	for _, i := range Round.Winning_hand {
		if c == i {
			highlight := canvas.NewRectangle(bundle.Highlight)
			highlight.Resize(size)
			highlight.Move(pos)

			return container.NewWithoutLayout(Cards.P4b, highlight)
		}
	}

	return Cards.P4b
}

// Set fifth players table card one image
func P5_a(c int) fyne.CanvasObject {
	size := fyne.NewSize(110, 150)
	pos := fyne.NewPos(336, 383)
	Cards.P5a = DisplayCard(c)
	Cards.P5a.Resize(size)
	Cards.P5a.Move(pos)

	for _, i := range Round.Winning_hand {
		if c == i {
			highlight := canvas.NewRectangle(bundle.Highlight)
			highlight.Resize(size)
			highlight.Move(pos)

			return container.NewWithoutLayout(Cards.P5a, highlight)
		}
	}

	return Cards.P5a
}

// Set fifth players table card two image
func P5_b(c int) fyne.CanvasObject {
	size := fyne.NewSize(110, 150)
	pos := fyne.NewPos(388, 383)
	Cards.P5b = DisplayCard(c)
	Cards.P5b.Resize(size)
	Cards.P5b.Move(pos)

	for _, i := range Round.Winning_hand {
		if c == i {
			highlight := canvas.NewRectangle(bundle.Highlight)
			highlight.Resize(size)
			highlight.Move(pos)

			return container.NewWithoutLayout(Cards.P5b, highlight)
		}
	}

	return Cards.P5b
}

// Set sixth players table card one image
func P6_a(c int) fyne.CanvasObject {
	size := fyne.NewSize(110, 150)
	pos := fyne.NewPos(65, 269)
	Cards.P6a = DisplayCard(c)
	Cards.P6a.Resize(size)
	Cards.P6a.Move(pos)

	for _, i := range Round.Winning_hand {
		if c == i {
			highlight := canvas.NewRectangle(bundle.Highlight)
			highlight.Resize(size)
			highlight.Move(pos)

			return container.NewWithoutLayout(Cards.P6a, highlight)
		}
	}

	return Cards.P6a
}

// Set sixth players table card two image
func P6_b(c int) fyne.CanvasObject {
	size := fyne.NewSize(110, 150)
	pos := fyne.NewPos(117, 269)
	Cards.P6b = DisplayCard(c)
	Cards.P6b.Resize(size)
	Cards.P6b.Move(pos)

	for _, i := range Round.Winning_hand {
		if c == i {
			highlight := canvas.NewRectangle(bundle.Highlight)
			highlight.Resize(size)
			highlight.Move(pos)

			return container.NewWithoutLayout(Cards.P6b, highlight)
		}
	}

	return Cards.P6b
}

// Returns int value for player table cards display.
// If player has no card hash values, no cards will be shown show
func Is_In(hash string, who int, end bool) int {
	if hash != "" {
		if end {
			return KeyCard(hash, who)
		} else {
			return 0
		}
	} else {
		return 99
	}
}

// Returns a custom card face image
//   - face defines which deck to look for
func CustomCard(c int, face string) *canvas.Image {
	dir := dreams.GetDir()
	mid := "/cards/" + face + "/"
	path := dir + mid + cardEnd(c)

	if dreams.FileExists(path, "Holdero") {
		return canvas.NewImageFromFile(path)
	}

	return canvas.NewImageFromImage(nil)
}

// Returns a custom card back image
//   - back defines which back to look for
func CustomBack(back string) *canvas.Image {
	dir := dreams.GetDir()
	post := "/cards/backs/" + back + ".png"
	path := dir + post

	if dreams.FileExists(path, "Holdero") {
		return canvas.NewImageFromFile(path)
	}

	return canvas.NewImageFromImage(nil)
}

// Used in CustomCard() to build image path
func cardEnd(card int) (suffix string) {
	if card > 0 && card < 53 {
		switch card {
		case 1:
			suffix = "card1.png"
		case 2:
			suffix = "card2.png"
		case 3:
			suffix = "card3.png"
		case 4:
			suffix = "card4.png"
		case 5:
			suffix = "card5.png"
		case 6:
			suffix = "card6.png"
		case 7:
			suffix = "card7.png"
		case 8:
			suffix = "card8.png"
		case 9:
			suffix = "card9.png"
		case 10:
			suffix = "card10.png"
		case 11:
			suffix = "card11.png"
		case 12:
			suffix = "card12.png"
		case 13:
			suffix = "card13.png"
		case 14:
			suffix = "card14.png"
		case 15:
			suffix = "card15.png"
		case 16:
			suffix = "card16.png"
		case 17:
			suffix = "card17.png"
		case 18:
			suffix = "card18.png"
		case 19:
			suffix = "card19.png"
		case 20:
			suffix = "card20.png"
		case 21:
			suffix = "card21.png"
		case 22:
			suffix = "card22.png"
		case 23:
			suffix = "card23.png"
		case 24:
			suffix = "card24.png"
		case 25:
			suffix = "card25.png"
		case 26:
			suffix = "card26.png"
		case 27:
			suffix = "card27.png"
		case 28:
			suffix = "card28.png"
		case 29:
			suffix = "card29.png"
		case 30:
			suffix = "card30.png"
		case 31:
			suffix = "card31.png"
		case 32:
			suffix = "card32.png"
		case 33:
			suffix = "card33.png"
		case 34:
			suffix = "card34.png"
		case 35:
			suffix = "card35.png"
		case 36:
			suffix = "card36.png"
		case 37:
			suffix = "card37.png"
		case 38:
			suffix = "card38.png"
		case 39:
			suffix = "card39.png"
		case 40:
			suffix = "card40.png"
		case 41:
			suffix = "card41.png"
		case 42:
			suffix = "card42.png"
		case 43:
			suffix = "card43.png"
		case 44:
			suffix = "card44.png"
		case 45:
			suffix = "card45.png"
		case 46:
			suffix = "card46.png"
		case 47:
			suffix = "card47.png"
		case 48:
			suffix = "card48.png"
		case 49:
			suffix = "card49.png"
		case 50:
			suffix = "card50.png"
		case 51:
			suffix = "card51.png"
		case 52:
			suffix = "card52.png"
		}
	} else {
		suffix = "card1.png"
	}
	return suffix

}

// Set Baccarat player card images
func PlayerCards(c1, c2, c3 int) fyne.CanvasObject {
	card1 := DisplayCard(c1)
	card2 := DisplayCard(c2)
	card3 := DisplayCard(c3)

	card1.Resize(fyne.NewSize(110, 150))
	card1.Move(fyne.NewPos(180, 39))

	card2.Resize(fyne.NewSize(110, 150))
	card2.Move(fyne.NewPos(290, 39))

	card3.Resize(fyne.NewSize(110, 150))
	card3.Move(fyne.NewPos(400, 39))

	return container.NewWithoutLayout(card1, card2, card3)
}

// Set Baccarat banker card images
func BankerCards(c1, c2, c3 int) fyne.CanvasObject {
	card1 := DisplayCard(c1)
	card2 := DisplayCard(c2)
	card3 := DisplayCard(c3)

	card1.Resize(fyne.NewSize(110, 150))
	card1.Move(fyne.NewPos(600, 39))

	card2.Resize(fyne.NewSize(110, 150))
	card2.Move(fyne.NewPos(710, 39))

	card3.Resize(fyne.NewSize(110, 150))
	card3.Move(fyne.NewPos(820, 39))

	return container.NewWithoutLayout(card1, card2, card3)
}

// Create a random suit for baccarat card
func BaccSuit(card int) int {
	if card == 99 {
		return 99
	}

	var suited int
	seed := rand.NewSource(time.Now().UnixNano())
	y := rand.New(seed)
	x := y.Intn(4) + 1

	if card == 0 {
		seed := rand.NewSource(time.Now().UnixNano())
		y := rand.New(seed)
		x := y.Intn(16) + 1

		switch x {
		case 1:
			suited = 10
		case 2:
			suited = 11
		case 3:
			suited = 12
		case 4:
			suited = 13
		case 5:
			suited = 23
		case 6:
			suited = 24
		case 7:
			suited = 25
		case 8:
			suited = 26
		case 9:
			suited = 36
		case 10:
			suited = 37
		case 11:
			suited = 38
		case 12:
			suited = 39
		case 13:
			suited = 49
		case 14:
			suited = 50
		case 15:
			suited = 51
		case 16:
			suited = 52
		}

		return suited
	}

	switch card {
	case 1:
		switch x {
		case 1:
			suited = 1
		case 2:
			suited = 14
		case 3:
			suited = 27
		case 4:
			suited = 40
		}
	case 2:
		switch x {
		case 1:
			suited = 2
		case 2:
			suited = 15
		case 3:
			suited = 28
		case 4:
			suited = 41
		}
	case 3:
		switch x {
		case 1:
			suited = 3
		case 2:
			suited = 16
		case 3:
			suited = 29
		case 4:
			suited = 42
		}
	case 4:
		switch x {
		case 1:
			suited = 4
		case 2:
			suited = 17
		case 3:
			suited = 30
		case 4:
			suited = 43
		}
	case 5:
		switch x {
		case 1:
			suited = 5
		case 2:
			suited = 18
		case 3:
			suited = 31
		case 4:
			suited = 44
		}
	case 6:
		switch x {
		case 1:
			suited = 6
		case 2:
			suited = 19
		case 3:
			suited = 32
		case 4:
			suited = 45
		}
	case 7:
		switch x {
		case 1:
			suited = 7
		case 2:
			suited = 20
		case 3:
			suited = 33
		case 4:
			suited = 46
		}
	case 8:
		switch x {
		case 1:
			suited = 8
		case 2:
			suited = 21
		case 3:
			suited = 34
		case 4:
			suited = 47
		}
	case 9:
		switch x {
		case 1:
			suited = 9
		case 2:
			suited = 22
		case 3:
			suited = 35
		case 4:
			suited = 48
		}
	case 10:
		switch x {
		case 1:
			suited = 10
		case 2:
			suited = 23
		case 3:
			suited = 36
		case 4:
			suited = 49
		}
	case 11:
		switch x {
		case 1:
			suited = 11
		case 2:
			suited = 24
		case 3:
			suited = 37
		case 4:
			suited = 50
		}
	case 12:
		switch x {
		case 1:
			suited = 12
		case 2:
			suited = 25
		case 3:
			suited = 38
		case 4:
			suited = 51
		}
	case 13:
		switch x {
		case 1:
			suited = 13
		case 2:
			suited = 26
		case 3:
			suited = 39
		case 4:
			suited = 52
		}
	}

	return suited
}

// Place Holdero card images
func placeHolderoCards(w fyne.Window) *fyne.Container {
	size := w.Content().Size()
	Cards.Layout = container.NewWithoutLayout(
		Hole_1(0, size.Width, size.Height),
		Hole_2(0, size.Width, size.Height),
		P1_a(Is_In(Round.Cards.P1C1, 1, Signal.End)),
		P1_b(Is_In(Round.Cards.P1C2, 1, Signal.End)),
		P2_a(Is_In(Round.Cards.P2C1, 2, Signal.End)),
		P2_b(Is_In(Round.Cards.P2C2, 2, Signal.End)),
		P3_a(Is_In(Round.Cards.P3C1, 3, Signal.End)),
		P3_b(Is_In(Round.Cards.P3C2, 3, Signal.End)),
		P4_a(Is_In(Round.Cards.P4C1, 4, Signal.End)),
		P4_b(Is_In(Round.Cards.P4C2, 4, Signal.End)),
		P5_a(Is_In(Round.Cards.P5C1, 5, Signal.End)),
		P5_b(Is_In(Round.Cards.P5C2, 5, Signal.End)),
		P6_a(Is_In(Round.Cards.P6C1, 6, Signal.End)),
		P6_b(Is_In(Round.Cards.P6C2, 6, Signal.End)),
		Flop_1(Round.Flop1),
		Flop_2(Round.Flop2),
		Flop_3(Round.Flop3),
		Turn(Round.TurnCard),
		River(Round.RiverCard))

	return Cards.Layout
}

// Refresh Holdero card images
func refreshHolderoCards(l1, l2 string, w fyne.Window) {
	size := w.Content().Size()
	Cards.Layout.Objects[0] = Hole_1(Card(l1), size.Width, size.Height)
	Cards.Layout.Objects[0].Refresh()

	Cards.Layout.Objects[1] = Hole_2(Card(l2), size.Width, size.Height)
	Cards.Layout.Objects[1].Refresh()

	Cards.Layout.Objects[2] = P1_a(Is_In(Round.Cards.P1C1, 1, Signal.End))
	Cards.Layout.Objects[2].Refresh()

	Cards.Layout.Objects[3] = P1_b(Is_In(Round.Cards.P1C2, 1, Signal.End))
	Cards.Layout.Objects[3].Refresh()

	Cards.Layout.Objects[4] = P2_a(Is_In(Round.Cards.P2C1, 2, Signal.End))
	Cards.Layout.Objects[4].Refresh()

	Cards.Layout.Objects[5] = P2_b(Is_In(Round.Cards.P2C2, 2, Signal.End))
	Cards.Layout.Objects[5].Refresh()

	Cards.Layout.Objects[6] = P3_a(Is_In(Round.Cards.P3C1, 3, Signal.End))
	Cards.Layout.Objects[6].Refresh()

	Cards.Layout.Objects[7] = P3_b(Is_In(Round.Cards.P3C2, 3, Signal.End))
	Cards.Layout.Objects[7].Refresh()

	Cards.Layout.Objects[8] = P4_a(Is_In(Round.Cards.P4C1, 4, Signal.End))
	Cards.Layout.Objects[8].Refresh()

	Cards.Layout.Objects[9] = P4_b(Is_In(Round.Cards.P4C2, 4, Signal.End))
	Cards.Layout.Objects[9].Refresh()

	Cards.Layout.Objects[10] = P5_a(Is_In(Round.Cards.P5C1, 5, Signal.End))
	Cards.Layout.Objects[10].Refresh()

	Cards.Layout.Objects[11] = P5_b(Is_In(Round.Cards.P5C2, 5, Signal.End))
	Cards.Layout.Objects[11].Refresh()

	Cards.Layout.Objects[12] = P6_a(Is_In(Round.Cards.P6C1, 6, Signal.End))
	Cards.Layout.Objects[12].Refresh()

	Cards.Layout.Objects[13] = P6_b(Is_In(Round.Cards.P6C2, 6, Signal.End))
	Cards.Layout.Objects[13].Refresh()

	Cards.Layout.Objects[14] = Flop_1(Round.Flop1)
	Cards.Layout.Objects[14].Refresh()

	Cards.Layout.Objects[15] = Flop_2(Round.Flop2)
	Cards.Layout.Objects[15].Refresh()

	Cards.Layout.Objects[16] = Flop_3(Round.Flop3)
	Cards.Layout.Objects[16].Refresh()

	Cards.Layout.Objects[17] = Turn(Round.TurnCard)
	Cards.Layout.Objects[17].Refresh()

	Cards.Layout.Objects[18] = River(Round.RiverCard)
	Cards.Layout.Objects[18].Refresh()

	Cards.Layout.Refresh()
}

// Main switch used to display playing card images
func DisplayCard(card int) *canvas.Image {
	if !Settings.Shared || Round.ID == 1 {
		if card == 99 {
			return canvas.NewImageFromImage(nil)
		}

		if card > 0 {
			i := Settings.faces.Select.SelectedIndex()
			switch i {
			case -1:
				return canvas.NewImageFromResource(DisplayLightCard(card))
			case 0:
				return canvas.NewImageFromResource(DisplayLightCard(card))
			case 1:
				return canvas.NewImageFromResource(DisplayDarkCard(card))
			default:
				return CustomCard(card, Settings.faces.Name)
			}
		}

		i := Settings.backs.Select.SelectedIndex()
		switch i {
		case -1:
			return canvas.NewImageFromResource(ResourceBack1Png)
		case 0:
			return canvas.NewImageFromResource(ResourceBack1Png)
		case 1:
			return canvas.NewImageFromResource(ResourceBack2Png)
		default:
			return CustomBack(Settings.backs.Name)
		}

	} else {
		if card == 99 {
			return canvas.NewImageFromImage(nil)
		} else if card > 0 {
			return CustomCard(card, Round.Face)
		} else {
			return CustomBack(Round.Back)
		}
	}
}

// Switch for standard light deck image
func DisplayLightCard(card int) fyne.Resource {
	if card > 0 && card < 53 {
		switch card {
		case 1:
			return ResourceLightcard1Png
		case 2:
			return ResourceLightcard2Png
		case 3:
			return ResourceLightcard3Png
		case 4:
			return ResourceLightcard4Png
		case 5:
			return ResourceLightcard5Png
		case 6:
			return ResourceLightcard6Png
		case 7:
			return ResourceLightcard7Png
		case 8:
			return ResourceLightcard8Png
		case 9:
			return ResourceLightcard9Png
		case 10:
			return ResourceLightcard10Png
		case 11:
			return ResourceLightcard11Png
		case 12:
			return ResourceLightcard12Png
		case 13:
			return ResourceLightcard13Png
		case 14:
			return ResourceLightcard14Png
		case 15:
			return ResourceLightcard15Png
		case 16:
			return ResourceLightcard16Png
		case 17:
			return ResourceLightcard17Png
		case 18:
			return ResourceLightcard18Png
		case 19:
			return ResourceLightcard19Png
		case 20:
			return ResourceLightcard20Png
		case 21:
			return ResourceLightcard21Png
		case 22:
			return ResourceLightcard22Png
		case 23:
			return ResourceLightcard23Png
		case 24:
			return ResourceLightcard24Png
		case 25:
			return ResourceLightcard25Png
		case 26:
			return ResourceLightcard26Png
		case 27:
			return ResourceLightcard27Png
		case 28:
			return ResourceLightcard28Png
		case 29:
			return ResourceLightcard29Png
		case 30:
			return ResourceLightcard30Png
		case 31:
			return ResourceLightcard31Png
		case 32:
			return ResourceLightcard32Png
		case 33:
			return ResourceLightcard33Png
		case 34:
			return ResourceLightcard34Png
		case 35:
			return ResourceLightcard35Png
		case 36:
			return ResourceLightcard36Png
		case 37:
			return ResourceLightcard37Png
		case 38:
			return ResourceLightcard38Png
		case 39:
			return ResourceLightcard39Png
		case 40:
			return ResourceLightcard40Png
		case 41:
			return ResourceLightcard41Png
		case 42:
			return ResourceLightcard42Png
		case 43:
			return ResourceLightcard43Png
		case 44:
			return ResourceLightcard44Png
		case 45:
			return ResourceLightcard45Png
		case 46:
			return ResourceLightcard46Png
		case 47:
			return ResourceLightcard47Png
		case 48:
			return ResourceLightcard48Png
		case 49:
			return ResourceLightcard49Png
		case 50:
			return ResourceLightcard50Png
		case 51:
			return ResourceLightcard51Png
		case 52:
			return ResourceLightcard52Png
		}
	}
	return nil
}

// Switch for standard dark deck image
func DisplayDarkCard(card int) fyne.Resource {
	if card > 0 && card < 53 {
		switch card {
		case 1:
			return ResourceDarkcard1Png
		case 2:
			return ResourceDarkcard2Png
		case 3:
			return ResourceDarkcard3Png
		case 4:
			return ResourceDarkcard4Png
		case 5:
			return ResourceDarkcard5Png
		case 6:
			return ResourceDarkcard6Png
		case 7:
			return ResourceDarkcard7Png
		case 8:
			return ResourceDarkcard8Png
		case 9:
			return ResourceDarkcard9Png
		case 10:
			return ResourceDarkcard10Png
		case 11:
			return ResourceDarkcard11Png
		case 12:
			return ResourceDarkcard12Png
		case 13:
			return ResourceDarkcard13Png
		case 14:
			return ResourceDarkcard14Png
		case 15:
			return ResourceDarkcard15Png
		case 16:
			return ResourceDarkcard16Png
		case 17:
			return ResourceDarkcard17Png
		case 18:
			return ResourceDarkcard18Png
		case 19:
			return ResourceDarkcard19Png
		case 20:
			return ResourceDarkcard20Png
		case 21:
			return ResourceDarkcard21Png
		case 22:
			return ResourceDarkcard22Png
		case 23:
			return ResourceDarkcard23Png
		case 24:
			return ResourceDarkcard24Png
		case 25:
			return ResourceDarkcard25Png
		case 26:
			return ResourceDarkcard26Png
		case 27:
			return ResourceDarkcard27Png
		case 28:
			return ResourceDarkcard28Png
		case 29:
			return ResourceDarkcard29Png
		case 30:
			return ResourceDarkcard30Png
		case 31:
			return ResourceDarkcard31Png
		case 32:
			return ResourceDarkcard32Png
		case 33:
			return ResourceDarkcard33Png
		case 34:
			return ResourceDarkcard34Png
		case 35:
			return ResourceDarkcard35Png
		case 36:
			return ResourceDarkcard36Png
		case 37:
			return ResourceDarkcard37Png
		case 38:
			return ResourceDarkcard38Png
		case 39:
			return ResourceDarkcard39Png
		case 40:
			return ResourceDarkcard40Png
		case 41:
			return ResourceDarkcard41Png
		case 42:
			return ResourceDarkcard42Png
		case 43:
			return ResourceDarkcard43Png
		case 44:
			return ResourceDarkcard44Png
		case 45:
			return ResourceDarkcard45Png
		case 46:
			return ResourceDarkcard46Png
		case 47:
			return ResourceDarkcard47Png
		case 48:
			return ResourceDarkcard48Png
		case 49:
			return ResourceDarkcard49Png
		case 50:
			return ResourceDarkcard50Png
		case 51:
			return ResourceDarkcard51Png
		case 52:
			return ResourceDarkcard52Png
		}
	}
	return nil
}
