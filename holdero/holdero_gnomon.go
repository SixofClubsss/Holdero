package holdero

import (
	"sort"
	"strconv"
	"strings"

	dreams "github.com/dReam-dApps/dReams"
	"github.com/dReam-dApps/dReams/menu"
	"github.com/dReam-dApps/dReams/rpc"

	"fyne.io/fyne/v2/canvas"
)

// Check if wallet owns Holdero table
func checkTableOwner(scid string) bool {
	if len(scid) != 64 || !menu.Gnomes.IsReady() {
		return false
	}

	check := strings.Trim(scid, " 0123456789")
	if check == "Holdero Tables:" {
		return false
	}

	owner, _ := menu.Gnomes.GetSCIDValuesByKey(scid, "owner:")
	if owner != nil {
		return owner[0] == rpc.Wallet.Address
	}

	return false
}

// Check if Holdero table is a tournament table
func checkHolderoContract(scid string) bool {
	if len(scid) != 64 || !menu.Gnomes.IsReady() {
		return false
	}

	_, deck := menu.Gnomes.GetSCIDValuesByKey(scid, "Deck Count:")
	_, version := menu.Gnomes.GetSCIDValuesByKey(scid, "V:")
	_, tourney := menu.Gnomes.GetSCIDValuesByKey(scid, "Tournament")
	if deck != nil && version != nil && version[0] >= 100 {
		signals.contract = true
	}

	if tourney != nil && tourney[0] == 1 {
		return true
	}

	return false
}

// Check Holdero table version
func checkTableVersion(scid string) uint64 {
	_, v := menu.Gnomes.GetSCIDValuesByKey(scid, "V:")

	if v != nil && v[0] >= 100 {
		return v[0]
	}
	return 0
}

// Make list of public and owned tables
func createTableList() {
	if menu.Gnomes.IsReady() {
		var owner bool
		list := []string{}
		owned := []string{}
		tables := menu.Gnomes.GetAllOwnersAndSCIDs()
		for scid := range tables {
			if !menu.Gnomes.IsReady() {
				break
			}

			if _, valid := menu.Gnomes.GetSCIDValuesByKey(scid, "Deck Count:"); valid != nil {
				_, version := menu.Gnomes.GetSCIDValuesByKey(scid, "V:")

				if version != nil {
					d := valid[0]
					v := version[0]

					headers := menu.GetSCHeaders(scid)
					name := "?"
					desc := "?"
					if headers != nil {
						if headers[1] != "" {
							desc = headers[1]
						}

						if headers[0] != "" {
							name = " " + headers[0]
						}
					}

					var hidden bool
					_, restrict := menu.Gnomes.GetSCIDValuesByKey(rpc.RatingSCID, "restrict")
					_, rating := menu.Gnomes.GetSCIDValuesByKey(rpc.RatingSCID, scid)

					if restrict != nil && rating != nil {
						menu.Control.Lock()
						menu.Control.Contract_rating[scid] = rating[0]
						menu.Control.Unlock()
						if rating[0] <= restrict[0] {
							hidden = true
						}
					}

					if d >= 1 && v == 110 && !hidden {
						list = append(list, name+"   "+desc+"   "+scid)
					}

					if d >= 1 && v >= 100 {
						if checkTableOwner(scid) {
							owned = append(owned, name+"   "+desc+"   "+scid)
							table.unlock.Hide()
							table.new.Show()
							owner = true
							table.owner.valid = true
						}
					}
				}
			}
		}

		if !owner {
			table.unlock.Show()
			table.new.Hide()
			table.owner.valid = false
		}

		t := len(list)
		list = append(list, "  Holdero Tables: "+strconv.Itoa(t))
		sort.Strings(list)
		table.Public.SCIDs = list

		sort.Strings(owned)
		table.Owned.SCIDs = owned

		table.Public.List.Refresh()
		table.Owned.List.Refresh()
	}
}

// Get current Holdero table menu stats
func getTableStats(scid string, single bool) {
	if menu.Gnomes.IsReady() && len(scid) == 64 {
		table.stats.version.Show()
		table.stats.last.Show()
		table.stats.name.Show()
		table.stats.desc.Show()
		if single {
			if h := menu.GetSCHeaders(scid); h != nil {
				table.stats.name.Text = (" Name: " + h[0])
				table.stats.name.Refresh()
				table.stats.desc.Text = (" Description: " + h[1])
				table.stats.desc.Refresh()
				if len(h[2]) > 6 {
					table.stats.image, _ = dreams.DownloadFile(h[2], h[0])
				} else {
					table.stats.image = *canvas.NewImageFromImage(nil)
				}

			} else {
				table.stats.name.Text = (" Name: ?")
				table.stats.name.Refresh()
				table.stats.desc.Text = (" Description: ?")
				table.stats.desc.Refresh()
				table.stats.image = *canvas.NewImageFromImage(nil)
			}
		}

		if _, v := menu.Gnomes.GetSCIDValuesByKey(scid, "V:"); v != nil {
			table.stats.version.Text = (" Table Version: " + strconv.Itoa(int(v[0])))
			table.stats.version.Refresh()
		} else {
			table.stats.version.Text = (" Table Version: ?")
			table.stats.version.Refresh()
		}

		if _, l := menu.Gnomes.GetSCIDValuesByKey(scid, "Last"); l != nil {
			time, _ := rpc.MsToTime(strconv.Itoa(int(l[0]) * 1000))
			table.stats.last.Text = (" Last Move: " + time.String())
			table.stats.last.Refresh()
		} else {
			table.stats.last.Text = (" Last Move: ?")
			table.stats.last.Refresh()
		}

		table.stats.owner.Hide()
		table.stats.chips.Hide()
		table.stats.blinds.Hide()

		if _, s := menu.Gnomes.GetSCIDValuesByKey(scid, "Seats at Table:"); s != nil {
			if s[0] > 1 {
				sit := 1
				if p2, _ := menu.Gnomes.GetSCIDValuesByKey(scid, "Player2 ID:"); p2 != nil {
					sit++
				}

				if p3, _ := menu.Gnomes.GetSCIDValuesByKey(scid, "Player3 ID:"); p3 != nil {
					sit++
				}

				if p4, _ := menu.Gnomes.GetSCIDValuesByKey(scid, "Player4 ID:"); p4 != nil {
					sit++
				}

				if p5, _ := menu.Gnomes.GetSCIDValuesByKey(scid, "Player5 ID:"); p5 != nil {
					sit++
				}

				if p6, _ := menu.Gnomes.GetSCIDValuesByKey(scid, "Player6 ID:"); p6 != nil {
					sit++
				}

				table.stats.seats.Text = (" Seats at Table: " + strconv.Itoa(int(s[0])-sit))
				table.stats.seats.Refresh()
			}

			table.stats.owner.Show()
			table.stats.chips.Show()
			table.stats.blinds.Show()

			table.stats.owner.Text = (" Owner: " + round.p1.name)
			table.stats.owner.Refresh()

			if round.asset {
				table.stats.chips.Text = (" Playing with: " + rpc.GetAssetSCIDName(round.assetID))
			} else {
				table.stats.chips.Text = (" Playing with: Dero")
			}
			table.stats.chips.Refresh()

			table.stats.blinds.Text = (" Blinds: " + blindString(rpc.Float64Type(round.BB), rpc.Float64Type(round.SB)))
			table.stats.blinds.Refresh()

		} else {
			table.stats.seats.Text = (" Table Closed")
			table.stats.seats.Refresh()
		}
	}
}
