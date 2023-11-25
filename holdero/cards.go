package holdero

import (
	dreams "github.com/dReam-dApps/dReams"
	"github.com/dReam-dApps/dReams/bundle"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
)

type playerCards struct {
	a fyne.CanvasObject
	b fyne.CanvasObject
}

type cards struct {
	hole1  fyne.CanvasObject
	hole2  fyne.CanvasObject
	flop1  fyne.CanvasObject
	flop2  fyne.CanvasObject
	flop3  fyne.CanvasObject
	turn   fyne.CanvasObject
	river  fyne.CanvasObject
	p1     playerCards
	p2     playerCards
	p3     playerCards
	p4     playerCards
	p5     playerCards
	p6     playerCards
	layout *fyne.Container
}

var card cards

// Set player hole card one image
//   - w and h of main window for resize
func Hole_1(c int, w, h float32) fyne.CanvasObject {
	card.hole1 = DisplayCard(c)
	card.hole1.Resize(fyne.NewSize(165, 225))
	card.hole1.Move(fyne.NewPos(w-335, h-335))

	return card.hole1
}

// Set player hole card two image
//   - w and h of main window for resize
func Hole_2(c int, w, h float32) fyne.CanvasObject {
	card.hole2 = DisplayCard(c)
	card.hole2.Resize(fyne.NewSize(165, 225))
	card.hole2.Move(fyne.NewPos(w-275, h-335))

	return card.hole2
}

// Set first flop card image
func Flop_1(c int) fyne.CanvasObject {
	size := fyne.NewSize(110, 150)
	pos := fyne.NewPos(260, 203)
	card.flop1 = DisplayCard(c)
	card.flop1.Resize(size)
	card.flop1.Move(pos)

	for _, i := range round.winningHand {
		if c == i {
			highlight := canvas.NewRectangle(bundle.Highlight)
			highlight.Resize(size)
			highlight.Move(pos)

			return container.NewWithoutLayout(card.flop1, highlight)
		}
	}

	return card.flop1
}

// Set second flop card image
func Flop_2(c int) fyne.CanvasObject {
	size := fyne.NewSize(110, 150)
	pos := fyne.NewPos(380, 203)
	card.flop2 = DisplayCard(c)
	card.flop2.Resize(size)
	card.flop2.Move(pos)

	for _, i := range round.winningHand {
		if c == i {
			highlight := canvas.NewRectangle(bundle.Highlight)
			highlight.Resize(size)
			highlight.Move(pos)

			return container.NewWithoutLayout(card.flop2, highlight)
		}
	}

	return card.flop2
}

// Set third flop card image
func Flop_3(c int) fyne.CanvasObject {
	size := fyne.NewSize(110, 150)
	pos := fyne.NewPos(500, 203)
	card.flop3 = DisplayCard(c)
	card.flop3.Resize(size)
	card.flop3.Move(pos)

	for _, i := range round.winningHand {
		if c == i {
			highlight := canvas.NewRectangle(bundle.Highlight)
			highlight.Resize(size)
			highlight.Move(pos)

			return container.NewWithoutLayout(card.flop3, highlight)
		}
	}

	return card.flop3
}

// Set turn card image
func Turn(c int) fyne.CanvasObject {
	size := fyne.NewSize(110, 150)
	pos := fyne.NewPos(620, 203)
	card.turn = DisplayCard(c)
	card.turn.Resize(size)
	card.turn.Move(pos)

	for _, i := range round.winningHand {
		if c == i {
			highlight := canvas.NewRectangle(bundle.Highlight)
			highlight.Resize(size)
			highlight.Move(pos)

			return container.NewWithoutLayout(card.turn, highlight)
		}
	}

	return card.turn
}

// Set river card image
func River(c int) fyne.CanvasObject {
	size := fyne.NewSize(110, 150)
	pos := fyne.NewPos(740, 203)
	card.river = DisplayCard(c)
	card.river.Resize(size)
	card.river.Move(pos)

	for _, i := range round.winningHand {
		if c == i {
			highlight := canvas.NewRectangle(bundle.Highlight)
			highlight.Resize(size)
			highlight.Move(pos)

			return container.NewWithoutLayout(card.river, highlight)
		}
	}

	return card.river
}

// Set first players table card one image
func P1_a(c int) fyne.CanvasObject {
	size := fyne.NewSize(110, 150)
	pos := fyne.NewPos(190, 25)
	card.p1.a = DisplayCard(c)
	card.p1.a.Resize(size)
	card.p1.a.Move(pos)

	for _, i := range round.winningHand {
		if c == i {
			highlight := canvas.NewRectangle(bundle.Highlight)
			highlight.Resize(size)
			highlight.Move(pos)

			return container.NewWithoutLayout(card.p1.a, highlight)
		}
	}

	return card.p1.a
}

// Set first players table card two image
func P1_b(c int) fyne.CanvasObject {
	size := fyne.NewSize(110, 150)
	pos := fyne.NewPos(242, 25)
	card.p1.b = DisplayCard(c)
	card.p1.b.Resize(size)
	card.p1.b.Move(pos)

	for _, i := range round.winningHand {
		if c == i {
			highlight := canvas.NewRectangle(bundle.Highlight)
			highlight.Resize(size)
			highlight.Move(pos)

			return container.NewWithoutLayout(card.p1.b, highlight)
		}
	}

	return card.p1.b
}

// Set second players table card one image
func P2_a(c int) fyne.CanvasObject {
	size := fyne.NewSize(110, 150)
	pos := fyne.NewPos(614, 25)
	card.p2.a = DisplayCard(c)
	card.p2.a.Resize(size)
	card.p2.a.Move(pos)

	for _, i := range round.winningHand {
		if c == i {
			highlight := canvas.NewRectangle(bundle.Highlight)
			highlight.Resize(size)
			highlight.Move(pos)

			return container.NewWithoutLayout(card.p2.a, highlight)
		}
	}

	return card.p2.a
}

// Set second players table card two image
func P2_b(c int) fyne.CanvasObject {
	size := fyne.NewSize(110, 150)
	pos := fyne.NewPos(666, 25)
	card.p2.b = DisplayCard(c)
	card.p2.b.Resize(size)
	card.p2.b.Move(pos)

	for _, i := range round.winningHand {
		if c == i {
			highlight := canvas.NewRectangle(bundle.Highlight)
			highlight.Resize(size)
			highlight.Move(pos)

			return container.NewWithoutLayout(card.p2.b, highlight)
		}
	}

	return card.p2.b
}

// Set third players table card one image
func P3_a(c int) fyne.CanvasObject {
	size := fyne.NewSize(110, 150)
	pos := fyne.NewPos(883, 129)
	card.p3.a = DisplayCard(c)
	card.p3.a.Resize(size)
	card.p3.a.Move(pos)

	for _, i := range round.winningHand {
		if c == i {
			highlight := canvas.NewRectangle(bundle.Highlight)
			highlight.Resize(size)
			highlight.Move(pos)

			return container.NewWithoutLayout(card.p3.a, highlight)
		}
	}

	return card.p3.a

}

// Set third players table card two image
func P3_b(c int) fyne.CanvasObject {
	size := fyne.NewSize(110, 150)
	pos := fyne.NewPos(935, 129)
	card.p3.b = DisplayCard(c)
	card.p3.b.Resize(size)
	card.p3.b.Move(pos)

	for _, i := range round.winningHand {
		if c == i {
			highlight := canvas.NewRectangle(bundle.Highlight)
			highlight.Resize(size)
			highlight.Move(pos)

			return container.NewWithoutLayout(card.p3.b, highlight)
		}
	}

	return card.p3.b
}

// Set fourth players table card one image
func P4_a(c int) fyne.CanvasObject {
	size := fyne.NewSize(110, 150)
	pos := fyne.NewPos(766, 383)
	card.p4.a = DisplayCard(c)
	card.p4.a.Resize(size)
	card.p4.a.Move(pos)

	for _, i := range round.winningHand {
		if c == i {
			highlight := canvas.NewRectangle(bundle.Highlight)
			highlight.Resize(size)
			highlight.Move(pos)

			return container.NewWithoutLayout(card.p4.a, highlight)
		}
	}

	return card.p4.a
}

// Set fourth players table card two image
func P4_b(c int) fyne.CanvasObject {
	size := fyne.NewSize(110, 150)
	pos := fyne.NewPos(818, 383)
	card.p4.b = DisplayCard(c)
	card.p4.b.Resize(size)
	card.p4.b.Move(pos)

	for _, i := range round.winningHand {
		if c == i {
			highlight := canvas.NewRectangle(bundle.Highlight)
			highlight.Resize(size)
			highlight.Move(pos)

			return container.NewWithoutLayout(card.p4.b, highlight)
		}
	}

	return card.p4.b
}

// Set fifth players table card one image
func P5_a(c int) fyne.CanvasObject {
	size := fyne.NewSize(110, 150)
	pos := fyne.NewPos(336, 383)
	card.p5.a = DisplayCard(c)
	card.p5.a.Resize(size)
	card.p5.a.Move(pos)

	for _, i := range round.winningHand {
		if c == i {
			highlight := canvas.NewRectangle(bundle.Highlight)
			highlight.Resize(size)
			highlight.Move(pos)

			return container.NewWithoutLayout(card.p5.a, highlight)
		}
	}

	return card.p5.a
}

// Set fifth players table card two image
func P5_b(c int) fyne.CanvasObject {
	size := fyne.NewSize(110, 150)
	pos := fyne.NewPos(388, 383)
	card.p5.b = DisplayCard(c)
	card.p5.b.Resize(size)
	card.p5.b.Move(pos)

	for _, i := range round.winningHand {
		if c == i {
			highlight := canvas.NewRectangle(bundle.Highlight)
			highlight.Resize(size)
			highlight.Move(pos)

			return container.NewWithoutLayout(card.p5.b, highlight)
		}
	}

	return card.p5.b
}

// Set sixth players table card one image
func P6_a(c int) fyne.CanvasObject {
	size := fyne.NewSize(110, 150)
	pos := fyne.NewPos(65, 269)
	card.p6.a = DisplayCard(c)
	card.p6.a.Resize(size)
	card.p6.a.Move(pos)

	for _, i := range round.winningHand {
		if c == i {
			highlight := canvas.NewRectangle(bundle.Highlight)
			highlight.Resize(size)
			highlight.Move(pos)

			return container.NewWithoutLayout(card.p6.a, highlight)
		}
	}

	return card.p6.a
}

// Set sixth players table card two image
func P6_b(c int) fyne.CanvasObject {
	size := fyne.NewSize(110, 150)
	pos := fyne.NewPos(117, 269)
	card.p6.b = DisplayCard(c)
	card.p6.b.Resize(size)
	card.p6.b.Move(pos)

	for _, i := range round.winningHand {
		if c == i {
			highlight := canvas.NewRectangle(bundle.Highlight)
			highlight.Resize(size)
			highlight.Move(pos)

			return container.NewWithoutLayout(card.p6.b, highlight)
		}
	}

	return card.p6.b
}

// Returns int value for player table cards display.
// If player has no card hash values, no cards will be shown show
func Is_In(hash string, who int, end bool) int {
	if hash != "" {
		if end {
			return keyCard(hash, who)
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

// Place Holdero card images
func placeHolderoCards(w fyne.Window) *fyne.Container {
	size := w.Content().Size()
	card.layout = container.NewWithoutLayout(
		Hole_1(0, size.Width, size.Height),
		Hole_2(0, size.Width, size.Height),
		P1_a(Is_In(round.cards.P1C1, 1, signals.end)),
		P1_b(Is_In(round.cards.P1C2, 1, signals.end)),
		P2_a(Is_In(round.cards.P2C1, 2, signals.end)),
		P2_b(Is_In(round.cards.P2C2, 2, signals.end)),
		P3_a(Is_In(round.cards.P3C1, 3, signals.end)),
		P3_b(Is_In(round.cards.P3C2, 3, signals.end)),
		P4_a(Is_In(round.cards.P4C1, 4, signals.end)),
		P4_b(Is_In(round.cards.P4C2, 4, signals.end)),
		P5_a(Is_In(round.cards.P5C1, 5, signals.end)),
		P5_b(Is_In(round.cards.P5C2, 5, signals.end)),
		P6_a(Is_In(round.cards.P6C1, 6, signals.end)),
		P6_b(Is_In(round.cards.P6C2, 6, signals.end)),
		Flop_1(round.cards.flop1),
		Flop_2(round.cards.flop2),
		Flop_3(round.cards.flop3),
		Turn(round.cards.turn),
		River(round.cards.river))

	return card.layout
}

// Refresh Holdero card images
func refreshHolderoCards(l1, l2 string, d *dreams.AppObject) {
	size := d.Window.Content().Size()
	align := float32(0)
	if d.OS() == "darwin" {
		align = 10
	}
	card.layout.Objects[0] = Hole_1(findCard(l1), size.Width+align, size.Height)
	card.layout.Objects[0].Refresh()

	card.layout.Objects[1] = Hole_2(findCard(l2), size.Width+align, size.Height)
	card.layout.Objects[1].Refresh()

	card.layout.Objects[2] = P1_a(Is_In(round.cards.P1C1, 1, signals.end))
	card.layout.Objects[2].Refresh()

	card.layout.Objects[3] = P1_b(Is_In(round.cards.P1C2, 1, signals.end))
	card.layout.Objects[3].Refresh()

	card.layout.Objects[4] = P2_a(Is_In(round.cards.P2C1, 2, signals.end))
	card.layout.Objects[4].Refresh()

	card.layout.Objects[5] = P2_b(Is_In(round.cards.P2C2, 2, signals.end))
	card.layout.Objects[5].Refresh()

	card.layout.Objects[6] = P3_a(Is_In(round.cards.P3C1, 3, signals.end))
	card.layout.Objects[6].Refresh()

	card.layout.Objects[7] = P3_b(Is_In(round.cards.P3C2, 3, signals.end))
	card.layout.Objects[7].Refresh()

	card.layout.Objects[8] = P4_a(Is_In(round.cards.P4C1, 4, signals.end))
	card.layout.Objects[8].Refresh()

	card.layout.Objects[9] = P4_b(Is_In(round.cards.P4C2, 4, signals.end))
	card.layout.Objects[9].Refresh()

	card.layout.Objects[10] = P5_a(Is_In(round.cards.P5C1, 5, signals.end))
	card.layout.Objects[10].Refresh()

	card.layout.Objects[11] = P5_b(Is_In(round.cards.P5C2, 5, signals.end))
	card.layout.Objects[11].Refresh()

	card.layout.Objects[12] = P6_a(Is_In(round.cards.P6C1, 6, signals.end))
	card.layout.Objects[12].Refresh()

	card.layout.Objects[13] = P6_b(Is_In(round.cards.P6C2, 6, signals.end))
	card.layout.Objects[13].Refresh()

	card.layout.Objects[14] = Flop_1(round.cards.flop1)
	card.layout.Objects[14].Refresh()

	card.layout.Objects[15] = Flop_2(round.cards.flop2)
	card.layout.Objects[15].Refresh()

	card.layout.Objects[16] = Flop_3(round.cards.flop3)
	card.layout.Objects[16].Refresh()

	card.layout.Objects[17] = Turn(round.cards.turn)
	card.layout.Objects[17].Refresh()

	card.layout.Objects[18] = River(round.cards.river)
	card.layout.Objects[18].Refresh()

	card.layout.Refresh()
}

// Main switch used to display playing card images
func DisplayCard(card int) *canvas.Image {
	if !Settings.sharing || round.ID == 1 {
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
			return CustomCard(card, round.cards.Faces.Name)
		} else {
			return CustomBack(round.cards.Backs.Name)
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
