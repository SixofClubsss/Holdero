package holdero

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"fyne.io/fyne/v2/canvas"
)

type sharedCards struct {
	avatar struct {
		p1 *canvas.Image
		p2 *canvas.Image
		p3 *canvas.Image
		p4 *canvas.Image
		p5 *canvas.Image
		p6 *canvas.Image
	}
	loaded struct {
		p1 bool
		p2 bool
		p3 bool
		p4 bool
		p5 bool
		p6 bool
	}
}

var shared sharedCards

// Clear Holdero card values when player changes table
func clearShared() {
	round.winningHand = []int{}
	round.display.results = ""
	round.first = true
	round.assetID = ""
	round.p1.url = ""
	round.p2.url = ""
	round.p3.url = ""
	round.p4.url = ""
	round.p5.name = ""
	round.p6.url = ""
	round.p1.name = ""
	round.p2.name = ""
	round.p3.name = ""
	round.p4.name = ""
	round.p5.name = ""
	round.p6.name = ""
	round.bettor = ""
	round.raiser = ""
	round.Last = 0
	signals.reveal = false
	signals.Out1 = false
	signals.odds = false
	Odds.Bot.Name = ""
	autoBetDefault()
	shared.loaded.p1 = false
	shared.loaded.p2 = false
	shared.loaded.p3 = false
	shared.loaded.p4 = false
	shared.loaded.p5 = false
	shared.loaded.p6 = false
	shared.avatar.p1 = canvas.NewImageFromImage(nil)
	shared.avatar.p2 = canvas.NewImageFromImage(nil)
	shared.avatar.p3 = canvas.NewImageFromImage(nil)
	shared.avatar.p4 = canvas.NewImageFromImage(nil)
	shared.avatar.p5 = canvas.NewImageFromImage(nil)
	shared.avatar.p6 = canvas.NewImageFromImage(nil)
}

// Code for storing card deck in memory
/*
func downloadMemoryDeck(url string) {
	var prog float64
	if url != "" {
		go func() {
			for i := 0; i < 53; i++ {
				float := float64(downloadSharedImages(url, i))
				prog = float / 53
				downloadPopUp(prog, i)
			}
			Settings.Shared = true
			Shared.window.Close()
		}()
	}
}

func downloadSharedImages(Url string, i int) int {
	fileName := "card" + strconv.Itoa(i) + ".png"
	logger.Println("[Holdero] Downloading ", Url+fileName)

	switch i {
	case 0:
		Shared.Back, _ = DownloadFile(Settings.BackUrl, fileName)
	case 1:
		Shared.Card1, _ = DownloadFile(Url+fileName, fileName)
	case 2:
		Shared.Card2, _ = DownloadFile(Url+fileName, fileName)
	case 3:
		Shared.Card3, _ = DownloadFile(Url+fileName, fileName)
	case 4:
		Shared.Card4, _ = DownloadFile(Url+fileName, fileName)
	case 5:
		Shared.Card5, _ = DownloadFile(Url+fileName, fileName)
	case 6:
		Shared.Card6, _ = DownloadFile(Url+fileName, fileName)
	case 7:
		Shared.Card7, _ = DownloadFile(Url+fileName, fileName)
	case 8:
		Shared.Card8, _ = DownloadFile(Url+fileName, fileName)
	case 9:
		Shared.Card9, _ = DownloadFile(Url+fileName, fileName)
	case 10:
		Shared.Card10, _ = DownloadFile(Url+fileName, fileName)
	case 11:
		Shared.Card11, _ = DownloadFile(Url+fileName, fileName)
	case 12:
		Shared.Card12, _ = DownloadFile(Url+fileName, fileName)
	case 13:
		Shared.Card13, _ = DownloadFile(Url+fileName, fileName)
	case 14:
		Shared.Card14, _ = DownloadFile(Url+fileName, fileName)
	case 15:
		Shared.Card15, _ = DownloadFile(Url+fileName, fileName)
	case 16:
		Shared.Card16, _ = DownloadFile(Url+fileName, fileName)
	case 17:
		Shared.Card17, _ = DownloadFile(Url+fileName, fileName)
	case 18:
		Shared.Card18, _ = DownloadFile(Url+fileName, fileName)
	case 19:
		Shared.Card19, _ = DownloadFile(Url+fileName, fileName)
	case 20:
		Shared.Card20, _ = DownloadFile(Url+fileName, fileName)
	case 21:
		Shared.Card21, _ = DownloadFile(Url+fileName, fileName)
	case 22:
		Shared.Card22, _ = DownloadFile(Url+fileName, fileName)
	case 23:
		Shared.Card23, _ = DownloadFile(Url+fileName, fileName)
	case 24:
		Shared.Card24, _ = DownloadFile(Url+fileName, fileName)
	case 25:
		Shared.Card25, _ = DownloadFile(Url+fileName, fileName)
	case 26:
		Shared.Card26, _ = DownloadFile(Url+fileName, fileName)
	case 27:
		Shared.Card27, _ = DownloadFile(Url+fileName, fileName)
	case 28:
		Shared.Card28, _ = DownloadFile(Url+fileName, fileName)
	case 29:
		Shared.Card29, _ = DownloadFile(Url+fileName, fileName)
	case 30:
		Shared.Card30, _ = DownloadFile(Url+fileName, fileName)
	case 31:
		Shared.Card31, _ = DownloadFile(Url+fileName, fileName)
	case 32:
		Shared.Card32, _ = DownloadFile(Url+fileName, fileName)
	case 33:
		Shared.Card33, _ = DownloadFile(Url+fileName, fileName)
	case 34:
		Shared.Card34, _ = DownloadFile(Url+fileName, fileName)
	case 35:
		Shared.Card35, _ = DownloadFile(Url+fileName, fileName)
	case 36:
		Shared.Card36, _ = DownloadFile(Url+fileName, fileName)
	case 37:
		Shared.Card37, _ = DownloadFile(Url+fileName, fileName)
	case 38:
		Shared.Card38, _ = DownloadFile(Url+fileName, fileName)
	case 39:
		Shared.Card39, _ = DownloadFile(Url+fileName, fileName)
	case 40:
		Shared.Card40, _ = DownloadFile(Url+fileName, fileName)
	case 41:
		Shared.Card41, _ = DownloadFile(Url+fileName, fileName)
	case 42:
		Shared.Card42, _ = DownloadFile(Url+fileName, fileName)
	case 43:
		Shared.Card43, _ = DownloadFile(Url+fileName, fileName)
	case 44:
		Shared.Card44, _ = DownloadFile(Url+fileName, fileName)
	case 45:
		Shared.Card45, _ = DownloadFile(Url+fileName, fileName)
	case 46:
		Shared.Card46, _ = DownloadFile(Url+fileName, fileName)
	case 47:
		Shared.Card47, _ = DownloadFile(Url+fileName, fileName)
	case 48:
		Shared.Card48, _ = DownloadFile(Url+fileName, fileName)
	case 49:
		Shared.Card49, _ = DownloadFile(Url+fileName, fileName)
	case 50:
		Shared.Card50, _ = DownloadFile(Url+fileName, fileName)
	case 51:
		Shared.Card51, _ = DownloadFile(Url+fileName, fileName)
	case 52:
		Shared.Card52, _ = DownloadFile(Url+fileName, fileName)
	}

	return i
}

// func SharedMemoryImage(c int) *canvas.Image {
// 	var card canvas.Image
// 	switch c {
// 	case 0:
// 		card = Shared.Back
// 	case 1:
// 		card = Shared.Card1
// 	case 2:
// 		card = Shared.Card2
// 	case 3:
// 		card = Shared.Card3
// 	case 4:
// 		card = Shared.Card4
// 	case 5:
// 		card = Shared.Card5
// 	case 6:
// 		card = Shared.Card6
// 	case 7:
// 		card = Shared.Card7
// 	case 8:
// 		card = Shared.Card8
// 	case 9:
// 		card = Shared.Card9
// 	case 10:
// 		card = Shared.Card10
// 	case 11:
// 		card = Shared.Card11
// 	case 12:
// 		card = Shared.Card12
// 	case 13:
// 		card = Shared.Card13
// 	case 14:
// 		card = Shared.Card14
// 	case 15:
// 		card = Shared.Card15
// 	case 16:
// 		card = Shared.Card16
// 	case 17:
// 		card = Shared.Card17
// 	case 18:
// 		card = Shared.Card18
// 	case 19:
// 		card = Shared.Card19
// 	case 20:
// 		card = Shared.Card20
// 	case 21:
// 		card = Shared.Card21
// 	case 22:
// 		card = Shared.Card22
// 	case 23:
// 		card = Shared.Card23
// 	case 24:
// 		card = Shared.Card24
// 	case 25:
// 		card = Shared.Card25
// 	case 26:
// 		card = Shared.Card26
// 	case 27:
// 		card = Shared.Card27
// 	case 28:
// 		card = Shared.Card28
// 	case 29:
// 		card = Shared.Card29
// 	case 30:
// 		card = Shared.Card30
// 	case 31:
// 		card = Shared.Card31
// 	case 32:
// 		card = Shared.Card32
// 	case 33:
// 		card = Shared.Card33
// 	case 34:
// 		card = Shared.Card34
// 	case 35:
// 		card = Shared.Card35
// 	case 36:
// 		card = Shared.Card36
// 	case 37:
// 		card = Shared.Card37
// 	case 38:
// 		card = Shared.Card38
// 	case 39:
// 		card = Shared.Card39
// 	case 40:
// 		card = Shared.Card40
// 	case 41:
// 		card = Shared.Card41
// 	case 42:
// 		card = Shared.Card42
// 	case 43:
// 		card = Shared.Card43
// 	case 44:
// 		card = Shared.Card44
// 	case 45:
// 		card = Shared.Card45
// 	case 46:
// 		card = Shared.Card46
// 	case 47:
// 		card = Shared.Card47
// 	case 48:
// 		card = Shared.Card48
// 	case 49:
// 		card = Shared.Card49
// 	case 50:
// 		card = Shared.Card50
// 	case 51:
// 		card = Shared.Card51
// 	case 52:
// 		card = Shared.Card52
// 	default:
// 		card = *canvas.NewImageFromFile("")
// 	}

// 	return &card
// }

func downloadProgress(p float64) fyne.Widget {
	Shared.progress = widget.NewProgressBar()
	this := binding.BindFloat(&p)
	Shared.progress.Bind(this)

	return Shared.progress
}

func downloadPopUp(p float64, i int) { /// pop up for loading progress
	if i == 0 {
		Shared.window = fyne.CurrentApp().NewWindow("Loading Custom Deck")
		Shared.window.Resize(fyne.NewSize(300, 30))
		Shared.window.SetFixedSize(true)
		Shared.window.SetIcon(nil)
		content := container.NewStack(downloadProgress(p))
		Shared.window.SetContent(content)
		Shared.window.Show()
	} else {
		content := container.NewStack(downloadProgress(p))
		Shared.window.SetContent(content)
	}
}
*/
/*
// for on demand

func SharedImage(c string) *canvas.Image {
	var card canvas.Image
	switch c {
	case "card0.png":
		card, _ = DownloadFile(table.Settings.BackUrl, c)
	case "card1.png":
		card, _ = DownloadFile(table.Settings.FaceUrl+c, c)
	case "card2.png":
		card, _ = DownloadFile(table.Settings.FaceUrl+c, c)
	case "card3.png":
		card, _ = DownloadFile(table.Settings.FaceUrl+c, c)
	case "card4.png":
		card, _ = DownloadFile(table.Settings.FaceUrl+c, c)
	case "card5.png":
		card, _ = DownloadFile(table.Settings.FaceUrl+c, c)
	case "card6.png":
		card, _ = DownloadFile(table.Settings.FaceUrl+c, c)
	case "card7.png":
		card, _ = DownloadFile(table.Settings.FaceUrl+c, c)
	case "card8.png":
		card, _ = DownloadFile(table.Settings.FaceUrl+c, c)
	case "card9.png":
		card, _ = DownloadFile(table.Settings.FaceUrl+c, c)
	case "card10.png":
		card, _ = DownloadFile(table.Settings.FaceUrl+c, c)
	case "card11.png":
		card, _ = DownloadFile(table.Settings.FaceUrl+c, c)
	case "card12.png":
		card, _ = DownloadFile(table.Settings.FaceUrl+c, c)
	case "card13.png":
		card, _ = DownloadFile(table.Settings.FaceUrl+c, c)
	case "card14.png":
		card, _ = DownloadFile(table.Settings.FaceUrl+c, c)
	case "card15.png":
		card, _ = DownloadFile(table.Settings.FaceUrl+c, c)
	case "card16.png":
		card, _ = DownloadFile(table.Settings.FaceUrl+c, c)
	case "card17.png":
		card, _ = DownloadFile(table.Settings.FaceUrl+c, c)
	case "card18.png":
		card, _ = DownloadFile(table.Settings.FaceUrl+c, c)
	case "card19.png":
		card, _ = DownloadFile(table.Settings.FaceUrl+c, c)
	case "card20.png":
		card, _ = DownloadFile(table.Settings.FaceUrl+c, c)
	case "card21.png":
		card, _ = DownloadFile(table.Settings.FaceUrl+c, c)
	case "card22.png":
		card, _ = DownloadFile(table.Settings.FaceUrl+c, c)
	case "card23.png":
		card, _ = DownloadFile(table.Settings.FaceUrl+c, c)
	case "card24.png":
		card, _ = DownloadFile(table.Settings.FaceUrl+c, c)
	case "card25.png":
		card, _ = DownloadFile(table.Settings.FaceUrl+c, c)
	case "card26.png":
		card, _ = DownloadFile(table.Settings.FaceUrl+c, c)
	case "card27.png":
		card, _ = DownloadFile(table.Settings.FaceUrl+c, c)
	case "card28.png":
		card, _ = DownloadFile(table.Settings.FaceUrl+c, c)
	case "card29.png":
		card, _ = DownloadFile(table.Settings.FaceUrl+c, c)
	case "card30.png":
		card, _ = DownloadFile(table.Settings.FaceUrl+c, c)
	case "card31.png":
		card, _ = DownloadFile(table.Settings.FaceUrl+c, c)
	case "card32.png":
		card, _ = DownloadFile(table.Settings.FaceUrl+c, c)
	case "card33.png":
		card, _ = DownloadFile(table.Settings.FaceUrl+c, c)
	case "card34.png":
		card, _ = DownloadFile(table.Settings.FaceUrl+c, c)
	case "card35.png":
		card, _ = DownloadFile(table.Settings.FaceUrl+c, c)
	case "card36.png":
		card, _ = DownloadFile(table.Settings.FaceUrl+c, c)
	case "card37.png":
		card, _ = DownloadFile(table.Settings.FaceUrl+c, c)
	case "card38.png":
		card, _ = DownloadFile(table.Settings.FaceUrl+c, c)
	case "card39.png":
		card, _ = DownloadFile(table.Settings.FaceUrl+c, c)
	case "card40.png":
		card, _ = DownloadFile(table.Settings.FaceUrl+c, c)
	case "card41.png":
		card, _ = DownloadFile(table.Settings.FaceUrl+c, c)
	case "card42.png":
		card, _ = DownloadFile(table.Settings.FaceUrl+c, c)
	case "card43.png":
		card, _ = DownloadFile(table.Settings.FaceUrl+c, c)
	case "card44.png":
		card, _ = DownloadFile(table.Settings.FaceUrl+c, c)
	case "card45.png":
		card, _ = DownloadFile(table.Settings.FaceUrl+c, c)
	case "card46.png":
		card, _ = DownloadFile(table.Settings.FaceUrl+c, c)
	case "card47.png":
		card, _ = DownloadFile(table.Settings.FaceUrl+c, c)
	case "card48.png":
		card, _ = DownloadFile(table.Settings.FaceUrl+c, c)
	case "card49.png":
		card, _ = DownloadFile(table.Settings.FaceUrl+c, c)
	case "card50.png":
		card, _ = DownloadFile(table.Settings.FaceUrl+c, c)
	case "card51.png":
		card, _ = DownloadFile(table.Settings.FaceUrl+c, c)
	case "card52.png":
		card, _ = DownloadFile(table.Settings.FaceUrl+c, c)
	default:
		card = *canvas.NewImageFromFile("")
	}

	fmt.Println(card)

	return &card
}
*/

var downloading bool

// Download a single uncompressed image image file to filepath
func downloadFileLocal(filepath string, url string) (err error) {
	downloading = true
	defer func() {
		downloading = false
	}()

	_, dir := os.Stat("cards")
	if os.IsNotExist(dir) {
		logger.Println("[Holdero] Creating Cards Dir")
		mkdir := os.Mkdir("cards", 0755)
		if mkdir != nil {
			logger.Errorln("[Holdero]", mkdir)
		} else {
			mksub := os.Mkdir("cards/backs", 0755)
			if mksub != nil {
				logger.Errorln("[Holdero]", mksub)
			}
		}
	}

	_, subdir := os.Stat("cards/backs")
	if os.IsNotExist(subdir) {
		logger.Println("[Holdero] Creating Backs Dir")
		mkdir := os.Mkdir("cards/backs", 0755)
		if mkdir != nil {
			logger.Errorln("[Holdero]", mkdir)
		}
	}

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

// Function to get and prepare deck assets for use in dReams
//   - face will be download path
func GetZipDeck(face, url string) {
	downloading = true
	downloadFileLocal("cards/"+face+".zip", url)
	files, err := Unzip("cards/"+face+".zip", "cards/"+face)

	if err != nil {
		logger.Errorln("[GetZipDeck]", err)
	}

	logger.Println("[Holdero] Unzipped files:\n" + strings.Join(files, "\n"))
	downloading = false
}

// Unzip a src file into destination
func Unzip(src string, destination string) ([]string, error) {
	var filenames []string

	r, err := zip.OpenReader(src)
	if err != nil {
		return filenames, err
	}

	defer r.Close()

	for _, f := range r.File {
		fpath := filepath.Join(destination, f.Name)

		if !strings.HasPrefix(fpath, filepath.Clean(destination)+string(os.PathSeparator)) {
			return filenames, fmt.Errorf("%s is an illegal filepath", fpath)
		}

		filenames = append(filenames, fpath)

		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return filenames, err
		}

		outFile, err := os.OpenFile(fpath,
			os.O_WRONLY|os.O_CREATE|os.O_TRUNC,
			f.Mode())

		if err != nil {
			return filenames, err
		}

		rc, err := f.Open()
		if err != nil {
			return filenames, err
		}

		_, err = io.Copy(outFile, rc)
		outFile.Close()
		rc.Close()

		if err != nil {
			return filenames, err
		}
	}
	return filenames, nil
}
