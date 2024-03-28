package holdero

import (
	"sort"
	"strconv"
	"strings"
	"time"

	dreams "github.com/dReam-dApps/dReams"
	"github.com/dReam-dApps/dReams/gnomes"
	"github.com/dReam-dApps/dReams/menu"
	"github.com/dReam-dApps/dReams/rpc"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

// Check if wallet owns Holdero table
func checkTableOwner(scid string) bool {
	if len(scid) != 64 || !gnomon.IsReady() {
		return false
	}

	check := strings.Trim(scid, " 0123456789")
	if check == "Holdero Tables:" {
		return false
	}

	owner, _ := gnomon.GetSCIDValuesByKey(scid, "owner:")
	if owner != nil {
		return owner[0] == rpc.Wallet.Address
	}

	return false
}

// Check if Holdero table is a tournament table
func checkHolderoContract(scid string) bool {
	if len(scid) != 64 || !gnomon.IsReady() {
		return false
	}

	_, deck := gnomon.GetSCIDValuesByKey(scid, "Deck Count:")
	_, version := gnomon.GetSCIDValuesByKey(scid, "V:")
	_, tourney := gnomon.GetSCIDValuesByKey(scid, "Tournament")
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
	_, v := gnomon.GetSCIDValuesByKey(scid, "V:")

	if v != nil && v[0] >= 100 {
		return v[0]
	}
	return 0
}

// Make list of public and owned tables with icon images for lists
func createTableList(progress *widget.ProgressBar) {
	if gnomon.IsReady() && !creating {
		creating = true
		defer func() {
			creating = false
		}()
		var owner bool
		var newPublic, newOwned, newFavorites []tableInfo

		tables := gnomon.GetAllOwnersAndSCIDs()
		if progress != nil {
			progress.Max = float64(len(tables))
		}

		for scid := range tables {
			if !gnomon.IsReady() || !rpc.Wallet.IsConnected() {
				newOwned = []tableInfo{}
				newPublic = []tableInfo{}
				newFavorites = []tableInfo{}
				break
			}

			if _, valid := gnomon.GetSCIDValuesByKey(scid, "Deck Count:"); valid != nil {
				_, version := gnomon.GetSCIDValuesByKey(scid, "V:")
				if version != nil {
					var info tableInfo
					headers := gnomes.GetSCHeaders(scid)
					if headers.Name != "" {
						info.name = headers.Name
						if headers.Description != "" {
							info.desc = headers.Description
						}

						if headers.IconURL != "" {
							if img, err := dreams.DownloadCanvas(headers.IconURL, headers.Name); err == nil {
								img.SetMinSize(fyne.NewSize(66, 66))
								info.image = &img
							} else {
								logger.Errorln("[Holdero]", err)
							}
						}
					}

					if _, last := gnomon.GetSCIDValuesByKey(scid, "Last"); last != nil {
						since := time.Since(time.Unix(int64(last[0]), 0))
						info.last = since.Truncate(time.Second).String()
					} else {
						info.last = "?"
					}

					var hidden bool
					_, restrict := gnomon.GetSCIDValuesByKey(rpc.RatingSCID, "restrict")
					_, rating := gnomon.GetSCIDValuesByKey(rpc.RatingSCID, scid)

					if restrict != nil && rating != nil {
						// TODO maybe replace with gnomes.SC
						menu.Control.Lock()
						menu.Control.Ratings[scid] = rating[0]
						menu.Control.Unlock()
						info.rating = rating[0]
						if rating[0] <= restrict[0] {
							hidden = true
						}
					}

					d := valid[0]
					v := version[0]

					info.scid = scid
					info.version = strconv.Itoa(int(v))

					if _, s := gnomon.GetSCIDValuesByKey(scid, "Seats at Table:"); s != nil {
						if s[0] > 1 {
							sit := 1
							if p2, _ := gnomon.GetSCIDValuesByKey(scid, "Player2 ID:"); p2 != nil {
								sit++
							}

							if p3, _ := gnomon.GetSCIDValuesByKey(scid, "Player3 ID:"); p3 != nil {
								sit++
							}

							if p4, _ := gnomon.GetSCIDValuesByKey(scid, "Player4 ID:"); p4 != nil {
								sit++
							}

							if p5, _ := gnomon.GetSCIDValuesByKey(scid, "Player5 ID:"); p5 != nil {
								sit++
							}

							if p6, _ := gnomon.GetSCIDValuesByKey(scid, "Player6 ID:"); p6 != nil {
								sit++
							}

							info.seats = "Seats: " + strconv.Itoa(int(s[0])-sit)
						}

						if owner, _ := gnomon.GetSCIDValuesByKey(scid, "owner:"); owner != nil {
							info.owner = owner[0]
						}

						if chips, _ := gnomon.GetSCIDValuesByKey(scid, "Chips"); chips != nil {
							if chips[0] == "ASSET" {
								if c, _ := gnomon.GetSCIDValuesByKey(scid, "HGC"); c != nil {
									info.chips = "Playing with: HGC"
								} else {
									info.chips = "Playing with: dReams"
								}
							} else {
								info.chips = "Playing with: DERO"
							}
						}

						if _, bb := gnomon.GetSCIDValuesByKey(scid, "BB:"); bb != nil {
							if _, sb := gnomon.GetSCIDValuesByKey(scid, "SB:"); bb != nil {
								info.blinds = "Blinds: " + blindString(rpc.Float64Type(bb[0]), rpc.Float64Type(sb[0]))
							}
						}

					} else {
						info.chips = "Table Closed"
					}

					if d >= 1 && v == 110 && !hidden {
						newPublic = append(newPublic, info)
					}

					if d >= 1 && v >= 100 {
						if checkTableOwner(scid) {
							newOwned = append(newOwned, info)
							table.unlock.Hide()
							table.new.Show()
							owner = true
							table.owner.valid = true
						}
					}
				}
			}
			if progress != nil {
				progress.SetValue(progress.Value + 1)
			}
		}

		// Sort public tables
		sort.Slice(newPublic, func(i, j int) bool {
			if newPublic[i].rating > newPublic[j].rating {
				return true
			}

			if newPublic[i].rating == newPublic[j].rating && newPublic[i].name > newPublic[j].name {
				return true
			}

			return false
		})

		// Sort owned tables
		sort.Slice(newOwned, func(i, j int) bool {
			if newOwned[i].rating > newOwned[j].rating {
				return true
			}

			if newOwned[i].rating == newOwned[j].rating && newOwned[i].name > newOwned[j].name {
				return true
			}

			return false
		})

		publicTables = newPublic
		ownedTables = newOwned

		for _, sc := range GetAccount().(accountData).Tables {
			for _, t := range publicTables {
				if t.scid == sc {
					newFavorites = append(newFavorites, t)
					break
				}
			}
		}

		// Sort fav tables
		sort.Slice(newFavorites, func(i, j int) bool {
			if newFavorites[i].rating > newFavorites[j].rating {
				return true
			}

			if newFavorites[i].rating == newFavorites[j].rating && newFavorites[i].name > newFavorites[j].name {
				return true
			}

			return false
		})

		favoriteTables = newFavorites

		if !owner {
			table.unlock.Show()
			table.new.Hide()
			table.owner.valid = false
		}

		table.Favorites.List.Refresh()
		table.Public.List.Refresh()
		table.Owned.List.Refresh()
	}
}
