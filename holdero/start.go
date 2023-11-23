package holdero

import (
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"syscall"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	dreams "github.com/dReam-dApps/dReams"
	"github.com/dReam-dApps/dReams/bundle"
	"github.com/dReam-dApps/dReams/dwidget"
	"github.com/dReam-dApps/dReams/menu"
	"github.com/dReam-dApps/dReams/rpc"
	"github.com/sirupsen/logrus"
)

const app_tag = "Holdero"

// Start Holdero dApp
func StartApp() {
	n := runtime.NumCPU()
	runtime.GOMAXPROCS(n)
	menu.InitLogrusLog(logrus.InfoLevel)
	config := menu.ReadDreamsConfig(app_tag)

	// Initialize Fyne app and window
	a := app.NewWithID(fmt.Sprintf("%s Desktop Client", app_tag))
	a.Settings().SetTheme(bundle.DeroTheme(config.Skin))
	w := a.NewWindow(app_tag)
	w.SetIcon(ResourcePokerBotIconPng)
	w.Resize(fyne.NewSize(1400, 800))
	w.SetMaster()
	done := make(chan struct{})

	// Initialize dReams AppObject and close func
	dreams.Theme.Img = *canvas.NewImageFromResource(nil)
	d := dreams.AppObject{
		App:        a,
		Window:     w,
		Background: container.NewStack(&dreams.Theme.Img),
	}
	d.SetChannels(1)
	d.SetTab("Holdero")

	closeFunc := func() {
		save := dreams.SaveData{
			Skin:   config.Skin,
			DBtype: menu.Gnomes.DBType,
		}

		if rpc.Daemon.Rpc == "" {
			save.Daemon = config.Daemon
		} else {
			save.Daemon = []string{rpc.Daemon.Rpc}
		}

		menu.WriteDreamsConfig(save)
		menu.CloseAppSignal(true)
		menu.Gnomes.Stop(app_tag)
		d.StopProcess()
		w.Close()
	}

	w.SetCloseIntercept(closeFunc)

	// Handle ctrl-c close
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println()
		closeFunc()
	}()

	// Initialize vars
	rpc.InitBalances()
	menu.Control.Contract_rating = make(map[string]uint64)
	menu.Gnomes.DBType = "boltdb"
	menu.Gnomes.Fast = true

	// Initialize asset widgets
	names := menu.NameEntry()
	asset_selects := []fyne.Widget{
		names.(*fyne.Container).Objects[1].(*widget.Select),
		FaceSelect(),
		BackSelect(),
		dreams.ThemeSelect(),
		AvatarSelect(menu.Assets.Asset_map),
		SharedDecks(),
	}

	// Create dwidget connection box with controls
	connect_box := dwidget.NewHorizontalEntries(app_tag, 1)
	connect_box.Button.OnTapped = func() {
		rpc.GetAddress(app_tag)
		rpc.Ping()
		if len(rpc.Wallet.Address) > 13 {
			menu.Control.Names.Options = []string{rpc.Wallet.Address[0:12]}
			menu.Control.Names.Refresh()
		}
		if rpc.Daemon.IsConnected() && !menu.Gnomes.IsInitialized() && !menu.Gnomes.Start {
			filter := []string{
				GetHolderoCode(0),
				GetHolderoCode(2),
				menu.NFA_SEARCH_FILTER,
				rpc.GetSCCode(rpc.GnomonSCID),
				rpc.GetSCCode(rpc.RatingSCID),
				rpc.GetSCCode(rpc.NameSCID)}

			go menu.StartGnomon(app_tag, menu.Gnomes.DBType, filter, 0, 0, nil)
		}
	}

	connect_box.Disconnect.OnChanged = func(b bool) {
		if !b {
			menu.Gnomes.Stop(app_tag)
		}
	}

	connect_box.AddDaemonOptions(config.Daemon)

	connect_box.Container.Objects[0].(*fyne.Container).Add(menu.StartIndicators())

	// Layout tabs
	tabs := container.NewAppTabs(
		container.NewTabItem(app_tag, LayoutAllItems(&d)),
		container.NewTabItem("Assets", menu.PlaceAssets(app_tag, asset_selects, ResourcePokerBotIconPng, d.Window)),
		container.NewTabItem("Swap", PlaceSwap()),
		container.NewTabItem("Log", rpc.SessionLog()))

	tabs.SetTabLocation(container.TabLocationBottom)

	// Stand alone process
	go func() {
		logger.Printf("[%s] %s %s %s", app_tag, rpc.DREAMSv, runtime.GOOS, runtime.GOARCH)
		time.Sleep(6 * time.Second)
		ticker := time.NewTicker(3 * time.Second)

		for {
			select {
			case <-ticker.C: // do on interval
				rpc.Ping()
				rpc.EchoWallet(app_tag)
				rpc.GetDreamsBalances(rpc.SCIDs)

				connect_box.RefreshBalance()
				if !rpc.Startup {
					menu.GnomonEndPoint()
				}

				if rpc.Daemon.IsConnected() && menu.Gnomes.IsInitialized() {
					connect_box.Disconnect.SetChecked(true)
					if menu.Gnomes.IsRunning() {
						menu.DisableIndexControls(false)
						menu.Gnomes.IndexContains()
						scids := " Indexed SCIDs: " + strconv.Itoa(int(menu.Gnomes.SCIDS))
						menu.Assets.Gnomes_index.Text = scids
						menu.Assets.Gnomes_index.Refresh()
						if menu.Gnomes.HasIndex(2) {
							menu.Gnomes.Checked(true)
						}
					}

					menu.Assets.Balances.Refresh()
					if rpc.Wallet.IsConnected() {
						menu.Assets.Swap.Show()
					} else {
						menu.Assets.Swap.Hide()
					}

					if menu.Gnomes.Indexer.LastIndexedHeight >= menu.Gnomes.Indexer.ChainHeight-3 {
						menu.Gnomes.Synced(true)
					} else {
						menu.Gnomes.Synced(false)
						menu.Gnomes.Checked(false)
					}
				} else {
					menu.DisableIndexControls(true)
					connect_box.Disconnect.SetChecked(false)
				}

				if rpc.Daemon.IsConnected() {
					rpc.Startup = false
				}

				d.SignalChannel()

			case <-d.Closing(): // exit
				logger.Printf("[%s] Closing...", app_tag)
				if menu.Gnomes.Icon_ind != nil {
					menu.Gnomes.Icon_ind.Stop()
				}
				ticker.Stop()
				d.CloseAllDapps()
				time.Sleep(time.Second)
				done <- struct{}{}
				return
			}
		}
	}()

	go func() {
		time.Sleep(450 * time.Millisecond)
		w.SetContent(container.NewStack(d.Background, container.NewStack(bundle.NewAlpha180(), tabs), container.NewVBox(layout.NewSpacer(), connect_box.Container)))
	}()
	w.ShowAndRun()
	<-done
	logger.Printf("[%s] Closed\n", app_tag)
}
