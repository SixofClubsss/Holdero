package holdero

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"time"

	"github.com/dReam-dApps/dReams/rpc"
	"github.com/deroproject/derohe/cryptography/crypto"
	dero "github.com/deroproject/derohe/rpc"
)

const (
	TourneySCID    = "c2e1ec16aed6f653aef99a06826b2b6f633349807d01fbb74cc0afb5ff99c3c7"
	Holdero100     = "95e69b382044ddc1467e030a80905cf637729612f65624e8d17bf778d4362b8d"
	HolderoSCID    = "e3f37573de94560e126a9020c0a5b3dfc7a4f3a4fbbe369fba93fbd219dc5fe9"
	pHolderoSCID   = "896834d57628d3a65076d3f4d84ddc7c5daf3e86b66a47f018abda6068afe2e6"
	HGCHolderoSCID = "efe646c48977fd776fee73cdd3df147a2668d3b7d965cdb7a187dda4d23005d8"
)

type displayStrings struct {
	seats    string
	pot      string
	blinds   string
	ante     string
	dealer   string
	playerId string
	readout  string
	results  string
}

type player struct {
	name   string
	url    string
	folded bool
}

type holderoValues struct {
	Version  int
	Contract string
	ID       int
	Players  int
	Turn     int
	Last     int64
	Pot      uint64
	BB       uint64
	SB       uint64
	Ante     uint64
	Wager    uint64
	Raised   uint64
	seed     string
	winner   string
	flop     bool
	localEnd bool
	asset    bool
	printed  bool
	notified bool
	tourney  bool
	assetID  string
	display  displayStrings
	p1       player
	p2       player
	p3       player
	p4       player
	p5       player
	p6       player
	bettor   string
	raiser   string
	cards    struct {
		CardSpecs
		flop1  int
		flop2  int
		flop3  int
		turn   int
		river  int
		Local1 string
		Local2 string
		P1C1   string
		P1C2   string
		P2C1   string
		P2C2   string
		P3C1   string
		P3C2   string
		P4C1   string
		P4C2   string
		P5C1   string
		P5C2   string
		P6C1   string
		P6C2   string

		Key1 string
		Key2 string
		Key3 string
		Key4 string
		Key5 string
		Key6 string
	}
	winningHand []int
	first       bool
	delay       bool
	trigger     struct {
		local bool
		flop  bool
		turn  bool
		river bool
	}
}

var round holderoValues

// Get Holdero SC data
func fetchHolderoSC() {
	if rpc.Daemon.IsConnected() && signals.contract {
		rpcClientD, ctx, cancel := rpc.SetDaemonClient(rpc.Daemon.Rpc)
		defer cancel()

		var result *dero.GetSC_Result
		params := dero.GetSC_Params{
			SCID:      round.Contract,
			Code:      false,
			Variables: true,
		}

		if err := rpcClientD.CallFor(ctx, &result, "DERO.GetSC", params); err != nil {
			logger.Errorln("[fetchHolderoSC]", err)
			return
		}

		var Pot_jv uint64
		Seats_jv := result.VariableStringKeys["Seats at Table:"]
		V_jv := result.VariableStringKeys["V:"]

		if V_jv != nil {
			round.Version = rpc.IntType(V_jv)
		}

		if Seats_jv != nil && rpc.IntType(Seats_jv) > 0 {
			// Count_jv := result.VariableStringKeys["Counter:"]
			Ante_jv := result.VariableStringKeys["Ante:"]
			BigBlind_jv := result.VariableStringKeys["BB:"]
			SmallBlind_jv := result.VariableStringKeys["SB:"]
			Turn_jv := result.VariableStringKeys["Player:"]
			OneId_jv := result.VariableStringKeys["Player1 ID:"]
			TwoId_jv := result.VariableStringKeys["Player2 ID:"]
			ThreeId_jv := result.VariableStringKeys["Player3 ID:"]
			FourId_jv := result.VariableStringKeys["Player4 ID:"]
			FiveId_jv := result.VariableStringKeys["Player5 ID:"]
			SixId_jv := result.VariableStringKeys["Player6 ID:"]
			Dealer_jv := result.VariableStringKeys["Dealer:"]
			Wager_jv := result.VariableStringKeys["Wager:"]
			Raised_jv := result.VariableStringKeys["Raised:"]
			FlopBool_jv := result.VariableStringKeys["Flop"]
			// TurnBool_jv = result.VariableStringKeys["Turn"]
			// RiverBool_jv = result.VariableStringKeys["River"]
			RevealBool_jv := result.VariableStringKeys["Reveal"]
			// Bet_jv := result.VariableStringKeys["Bet"]
			Full_jv := result.VariableStringKeys["Full"]
			// Open_jv := result.VariableStringKeys["Open"]
			Seed_jv := result.VariableStringKeys["HandSeed"]
			Face_jv := result.VariableStringKeys["Face:"]
			//Back_jv := result.VariableStringKeys["Back:"]
			Flop1_jv := result.VariableStringKeys["FlopCard1"]
			Flop2_jv := result.VariableStringKeys["FlopCard2"]
			Flop3_jv := result.VariableStringKeys["FlopCard3"]
			TurnCard_jv := result.VariableStringKeys["TurnCard"]
			RiverCard_jv := result.VariableStringKeys["RiverCard"]
			P1C1_jv := result.VariableStringKeys["Player1card1"]
			P1C2_jv := result.VariableStringKeys["Player1card2"]
			P2C1_jv := result.VariableStringKeys["Player2card1"]
			P2C2_jv := result.VariableStringKeys["Player2card2"]
			P3C1_jv := result.VariableStringKeys["Player3card1"]
			P3C2_jv := result.VariableStringKeys["Player3card2"]
			P4C1_jv := result.VariableStringKeys["Player4card1"]
			P4C2_jv := result.VariableStringKeys["Player4card2"]
			P5C1_jv := result.VariableStringKeys["Player5card1"]
			P5C2_jv := result.VariableStringKeys["Player5card2"]
			P6C1_jv := result.VariableStringKeys["Player6card1"]
			P6C2_jv := result.VariableStringKeys["Player6card2"]
			P1F_jv := result.VariableStringKeys["0F"]
			P2F_jv := result.VariableStringKeys["1F"]
			P3F_jv := result.VariableStringKeys["2F"]
			P4F_jv := result.VariableStringKeys["3F"]
			P5F_jv := result.VariableStringKeys["4F"]
			P6F_jv := result.VariableStringKeys["5F"]
			P1Out_jv := result.VariableStringKeys["0SO"]
			// P2Out_jv = result.VariableStringKeys["1SO"]
			// P3Out_jv = result.VariableStringKeys["2SO"]
			// P4Out_jv = result.VariableStringKeys["3SO"]
			// P5Out_jv = result.VariableStringKeys["4SO"]
			// P6Out_jv = result.VariableStringKeys["5SO"]
			Key1_jv := result.VariableStringKeys["Player1Key"]
			Key2_jv := result.VariableStringKeys["Player2Key"]
			Key3_jv := result.VariableStringKeys["Player3Key"]
			Key4_jv := result.VariableStringKeys["Player4Key"]
			Key5_jv := result.VariableStringKeys["Player5Key"]
			Key6_jv := result.VariableStringKeys["Player6Key"]
			End_jv := result.VariableStringKeys["End"]
			Chips_jv := result.VariableStringKeys["Chips"]
			Tourney_jv := result.VariableStringKeys["Tournament"]
			Last_jv := result.VariableStringKeys["Last"]

			if Last_jv != nil {
				round.Last = int64(rpc.Float64Type(Last_jv))
			} else {
				round.Last = 0
			}

			if Tourney_jv == nil {
				round.tourney = false
				if Chips_jv != nil {
					if rpc.HexToString(Chips_jv) == "ASSET" {
						round.asset = true
						if _, ok := result.VariableStringKeys["dReams"].(string); ok {
							Pot_jv = result.Balances[rpc.DreamsSCID]
							round.assetID = rpc.DreamsSCID
						} else if _, ok = result.VariableStringKeys["HGC"].(string); ok {
							Pot_jv = result.Balances[rpc.HgcSCID]
							round.assetID = rpc.HgcSCID
						}
					} else {
						round.asset = false
						round.assetID = ""
						Pot_jv = result.Balances["0000000000000000000000000000000000000000000000000000000000000000"]
					}
				} else {
					round.asset = false
					round.assetID = ""
					Pot_jv = result.Balances["0000000000000000000000000000000000000000000000000000000000000000"]
				}
			} else {
				round.tourney = true
				if Chips_jv != nil {
					if rpc.HexToString(Chips_jv) == "ASSET" {
						round.asset = true
						Pot_jv = result.Balances[TourneySCID]
					} else {
						round.asset = false
						round.assetID = ""
						Pot_jv = result.Balances["0000000000000000000000000000000000000000000000000000000000000000"]
					}
				} else {
					round.asset = false
					round.assetID = ""
					Pot_jv = result.Balances["0000000000000000000000000000000000000000000000000000000000000000"]
				}
			}

			round.Ante = rpc.Uint64Type(Ante_jv)
			round.BB = rpc.Uint64Type(BigBlind_jv)
			round.SB = rpc.Uint64Type(SmallBlind_jv)
			round.Pot = Pot_jv

			hasFolded(P1F_jv, P2F_jv, P3F_jv, P4F_jv, P5F_jv, P6F_jv)
			allFolded(P1F_jv, P2F_jv, P3F_jv, P4F_jv, P5F_jv, P6F_jv, Seats_jv)

			if !round.localEnd {
				getCommCardValues(Flop1_jv, Flop2_jv, Flop3_jv, TurnCard_jv, RiverCard_jv)
				getPlayerCardValues(P1C1_jv, P1C2_jv, P2C1_jv, P2C2_jv, P3C1_jv, P3C2_jv, P4C1_jv, P4C2_jv, P5C1_jv, P5C2_jv, P6C1_jv, P6C2_jv)
			}

			if !rpc.Startup {
				setHolderoName(OneId_jv, TwoId_jv, ThreeId_jv, FourId_jv, FiveId_jv, SixId_jv)
				setSignals(Pot_jv, P1Out_jv)
			}

			if P1Out_jv != nil {
				signals.Out1 = true
			} else {
				signals.Out1 = false
			}

			tableOpen(Seats_jv, Full_jv, TwoId_jv, ThreeId_jv, FourId_jv, FiveId_jv, SixId_jv)

			if FlopBool_jv != nil {
				round.flop = true
				rpc.Wallet.KeyLock = false
			} else {
				round.flop = false
			}

			round.display.playerId = checkPlayerId(getAvatar(1, OneId_jv), getAvatar(2, TwoId_jv), getAvatar(3, ThreeId_jv), getAvatar(4, FourId_jv), getAvatar(5, FiveId_jv), getAvatar(6, SixId_jv))

			if Wager_jv != nil {
				if round.bettor == "" {
					round.bettor = findBettor(Turn_jv)
				}
				round.Wager = rpc.Uint64Type(Wager_jv)
			} else {
				round.bettor = ""
				round.Wager = 0
			}

			if Raised_jv != nil {
				if round.raiser == "" {
					round.raiser = findBettor(Turn_jv)
				}
				round.Raised = rpc.Uint64Type(Raised_jv)
			} else {
				round.raiser = ""
				round.Raised = 0
			}

			if round.ID == rpc.IntType(Turn_jv)+1 {
				signals.myTurn = true
			} else if round.ID == 1 && Turn_jv == Seats_jv {
				signals.myTurn = true
			} else {
				signals.myTurn = false
			}

			round.display.pot = rpc.FromAtomic(Pot_jv, 5)
			round.display.seats = fmt.Sprint(Seats_jv)
			round.display.ante = rpc.FromAtomic(Ante_jv, 5)
			round.display.blinds = blindString(BigBlind_jv, SmallBlind_jv)
			round.display.dealer = rpc.AddOne(Dealer_jv)

			round.seed = fmt.Sprint(Seed_jv)

			if face, ok := Face_jv.(string); ok {
				if face != "nil" {
					var c = &CardSpecs{}
					if err := json.Unmarshal([]byte(rpc.HexToString(face)), c); err == nil {
						round.cards.Faces.Name = c.Faces.Name
						round.cards.Backs.Name = c.Backs.Name
						round.cards.Faces.Url = c.Faces.Url
						round.cards.Backs.Url = c.Backs.Url
					}
				}
			} else {
				round.cards.Faces.Name = ""
				round.cards.Backs.Name = ""
				round.cards.Faces.Url = ""
				round.cards.Backs.Url = ""
			}

			if round.ID != 1 {
				Settings.faces.URL = round.cards.Faces.Url
				Settings.backs.URL = round.cards.Backs.Url
			}

			// // Unused at moment
			// if back, ok := Back_jv.(string); ok {
			// 	if back != "nil" {
			// 		var a = &TableSpecs{}
			// 		json.Unmarshal([]byte(FromHexToString(back)), a)
			// 	}
			// }

			if RevealBool_jv != nil && !signals.reveal && !round.localEnd {
				if rpc.AddOne(Turn_jv) == round.display.playerId {
					signals.clicked = true
					signals.height = rpc.Wallet.Height
					signals.reveal = true
					go RevealKey(rpc.Wallet.ClientKey)
				}
			}

			if Turn_jv != Seats_jv {
				round.display.readout = turnReadout(Turn_jv)
				if turn, ok := Turn_jv.(float64); ok {
					round.Turn = int(turn) + 1
				}
			} else {
				round.Turn = 1
				round.display.readout = turnReadout(float64(0))
			}

			if End_jv != nil {
				round.cards.Key1 = fmt.Sprint(Key1_jv)
				round.cards.Key2 = fmt.Sprint(Key2_jv)
				round.cards.Key3 = fmt.Sprint(Key3_jv)
				round.cards.Key4 = fmt.Sprint(Key4_jv)
				round.cards.Key5 = fmt.Sprint(Key5_jv)
				round.cards.Key6 = fmt.Sprint(Key6_jv)
				signals.end = true

			}

			if round.Version >= 110 && round.ID == 1 && signals.times.kick > 0 && !signals.myTurn && round.Pot > 0 && !round.localEnd && !signals.end {
				if round.Last != 0 {
					now := time.Now().Unix()
					if now > round.Last+int64(signals.times.kick)+18 {
						if rpc.Wallet.Height > signals.times.block+3 {
							TimeOut()
							signals.times.block = rpc.Wallet.Height
						}
					}
				}
			}

			winningHand(End_jv)
		} else {
			closedTable()
		}

		potIsEmpty(Pot_jv)
		allFoldedWinner()
	}
}

// Submit playerId, name, avatar and sit at Holdero table
//   - name and av are for name and avatar in player id string
func SitDown(name, av string) {
	rpcClientW, ctx, cancel := rpc.SetWalletClient(rpc.Wallet.Rpc, rpc.Wallet.UserPass)
	defer cancel()

	var player playerId
	player.Id = rpc.Wallet.IdHash
	player.Name = name
	player.Avatar = av

	mar, _ := json.Marshal(player)
	hx := hex.EncodeToString(mar)

	arg1 := dero.Argument{Name: "entrypoint", DataType: "S", Value: "PlayerEntry"}
	arg2 := dero.Argument{Name: "address", DataType: "S", Value: hx}
	args := dero.Arguments{arg1, arg2}
	txid := dero.Transfer_Result{}

	t1 := dero.Transfer{
		Destination: "dero1qyr8yjnu6cl2c5yqkls0hmxe6rry77kn24nmc5fje6hm9jltyvdd5qq4hn5pn",
		Amount:      0,
		Burn:        0,
	}

	t := []dero.Transfer{t1}
	fee := rpc.GasEstimate(round.Contract, "[Holdero]", args, t, rpc.HighLimitFee)
	params := &dero.Transfer_Params{
		Transfers: t,
		SC_ID:     round.Contract,
		SC_RPC:    args,
		Ringsize:  2,
		Fees:      fee,
	}

	if err := rpcClientW.CallFor(ctx, &txid, "transfer", params); err != nil {
		rpc.PrintError("[Holdero] Sit Down: %s", err)
		return
	}

	rpc.PrintLog("[Holdero] Sit Down TX: %s", txid)
}

// Leave Holdero seat on players turn
func Leave() {
	rpcClientW, ctx, cancel := rpc.SetWalletClient(rpc.Wallet.Rpc, rpc.Wallet.UserPass)
	defer cancel()

	checkoutId := rpc.StringToInt(round.display.playerId)
	singleNameClear(checkoutId)
	arg1 := dero.Argument{Name: "entrypoint", DataType: "S", Value: "PlayerLeave"}
	arg2 := dero.Argument{Name: "id", DataType: "U", Value: checkoutId}
	args := dero.Arguments{arg1, arg2}
	txid := dero.Transfer_Result{}

	t1 := dero.Transfer{
		Destination: "dero1qyr8yjnu6cl2c5yqkls0hmxe6rry77kn24nmc5fje6hm9jltyvdd5qq4hn5pn",
		Amount:      0,
		Burn:        0,
	}

	t := []dero.Transfer{t1}
	fee := rpc.GasEstimate(round.Contract, "[Holdero]", args, t, rpc.LowLimitFee)
	params := &dero.Transfer_Params{
		Transfers: t,
		SC_ID:     round.Contract,
		SC_RPC:    args,
		Ringsize:  2,
		Fees:      fee,
	}

	if err := rpcClientW.CallFor(ctx, &txid, "transfer", params); err != nil {
		rpc.PrintError("[Holdero] Leave: %s", err)
		return
	}

	rpc.PrintLog("[Holdero] Leave TX: %s", txid)
}

// Owner table settings for Holdero
//   - seats defines max players at table
//   - bb, sb and ante define big blind, small blind and antes. Ante can be 0
//   - chips defines if tables is using Dero or assets
//   - name and av are for name and avatar in owners id string
func SetTable(seats int, bb, sb, ante uint64, chips, name, av string) {
	rpcClientW, ctx, cancel := rpc.SetWalletClient(rpc.Wallet.Rpc, rpc.Wallet.UserPass)
	defer cancel()

	var player playerId
	player.Id = rpc.Wallet.IdHash
	player.Name = name
	player.Avatar = av

	mar, _ := json.Marshal(player)
	hx := hex.EncodeToString(mar)

	var args dero.Arguments
	arg1 := dero.Argument{Name: "entrypoint", DataType: "S", Value: "SetTable"}
	arg2 := dero.Argument{Name: "seats", DataType: "U", Value: seats}
	arg3 := dero.Argument{Name: "bigBlind", DataType: "U", Value: bb}
	arg4 := dero.Argument{Name: "smallBlind", DataType: "U", Value: sb}
	arg5 := dero.Argument{Name: "ante", DataType: "U", Value: ante}
	arg6 := dero.Argument{Name: "address", DataType: "S", Value: hx}
	txid := dero.Transfer_Result{}

	if round.Version < 110 {
		args = dero.Arguments{arg1, arg2, arg3, arg4, arg5, arg6}
	} else if round.Version == 110 {
		arg7 := dero.Argument{Name: "chips", DataType: "S", Value: chips}
		args = dero.Arguments{arg1, arg2, arg3, arg4, arg5, arg6, arg7}
	}

	t1 := dero.Transfer{
		Destination: "dero1qyr8yjnu6cl2c5yqkls0hmxe6rry77kn24nmc5fje6hm9jltyvdd5qq4hn5pn",
		Amount:      0,
		Burn:        0,
	}

	t := []dero.Transfer{t1}
	fee := rpc.GasEstimate(round.Contract, "[Holdero]", args, t, rpc.HighLimitFee)
	params := &dero.Transfer_Params{
		Transfers: t,
		SC_ID:     round.Contract,
		SC_RPC:    args,
		Ringsize:  2,
		Fees:      fee,
	}

	if err := rpcClientW.CallFor(ctx, &txid, "transfer", params); err != nil {
		rpc.PrintError("[Holdero] Set Table: %s", err)
		return
	}

	rpc.PrintLog("[Holdero] Set Table TX: %s", txid)
}

// Submit blinds/ante to deal Holdero hand
func DealHand() (tx string) {
	rpcClientW, ctx, cancel := rpc.SetWalletClient(rpc.Wallet.Rpc, rpc.Wallet.UserPass)
	defer cancel()

	if !rpc.Wallet.KeyLock {
		rpc.Wallet.ClientKey = generateKey()
	}

	arg1 := dero.Argument{Name: "entrypoint", DataType: "S", Value: "DealHand"}
	arg2 := dero.Argument{Name: "pcSeed", DataType: "H", Value: rpc.Wallet.ClientKey}
	args := dero.Arguments{arg1, arg2}
	txid := dero.Transfer_Result{}

	var amount uint64

	if round.Pot == 0 {
		amount = round.Ante + round.SB
	} else if round.Pot == round.SB || round.Pot == round.Ante+round.SB {
		amount = round.Ante + round.BB
	} else {
		amount = round.Ante
	}

	t := []dero.Transfer{}
	if round.asset {
		t1 := dero.Transfer{
			Destination: "dero1qyr8yjnu6cl2c5yqkls0hmxe6rry77kn24nmc5fje6hm9jltyvdd5qq4hn5pn",
			Amount:      500,
			Burn:        0,
		}

		if round.tourney {
			t2 := dero.Transfer{
				SCID:        crypto.HashHexToHash(TourneySCID),
				Destination: "dero1qyr8yjnu6cl2c5yqkls0hmxe6rry77kn24nmc5fje6hm9jltyvdd5qq4hn5pn",
				Burn:        amount,
			}
			t = append(t, t1, t2)
		} else {
			t2 := rpc.GetAssetSCIDforTransfer(amount, round.assetID)
			if t2.Destination == "" {
				rpc.PrintError("[Holdero] Deal: err getting asset SCID for transfer")
				return
			}
			t = append(t, t1, t2)
		}
	} else {
		t1 := dero.Transfer{
			Destination: "dero1qyr8yjnu6cl2c5yqkls0hmxe6rry77kn24nmc5fje6hm9jltyvdd5qq4hn5pn",
			Amount:      500,
			Burn:        amount,
		}
		t = append(t, t1)
	}

	fee := rpc.GasEstimate(round.Contract, "[Holdero]", args, t, rpc.LowLimitFee)
	params := &dero.Transfer_Params{
		Transfers: t,
		SC_ID:     round.Contract,
		SC_RPC:    args,
		Ringsize:  2,
		Fees:      fee,
	}

	if err := rpcClientW.CallFor(ctx, &txid, "transfer", params); err != nil {
		rpc.PrintError("[Holdero] Deal: %s", err)
		return
	}

	round.display.results = ""
	updateStatsWager(float64(amount) / 100000)
	rpc.PrintLog("[Holdero] Deal TX: %s", txid)

	return txid.TXID
}

// Make Holdero bet
func Bet(amt string) (tx string) {
	rpcClientW, ctx, cancel := rpc.SetWalletClient(rpc.Wallet.Rpc, rpc.Wallet.UserPass)
	defer cancel()

	arg1 := dero.Argument{Name: "entrypoint", DataType: "S", Value: "Bet"}
	args := dero.Arguments{arg1}
	txid := dero.Transfer_Result{}

	var t1 dero.Transfer
	if round.asset {
		if round.tourney {
			t1 = dero.Transfer{
				SCID:        crypto.HashHexToHash(TourneySCID),
				Destination: "dero1qyr8yjnu6cl2c5yqkls0hmxe6rry77kn24nmc5fje6hm9jltyvdd5qq4hn5pn",
				Burn:        rpc.ToAtomic(amt, 1),
			}
		} else {
			t1 = rpc.GetAssetSCIDforTransfer(rpc.ToAtomic(amt, 1), round.assetID)
			if t1.Destination == "" {
				rpc.PrintError("[Holdero] Bet: err getting asset SCID for transfer")
				return
			}
		}
	} else {
		t1 = dero.Transfer{
			Destination: "dero1qyr8yjnu6cl2c5yqkls0hmxe6rry77kn24nmc5fje6hm9jltyvdd5qq4hn5pn",
			Amount:      0,
			Burn:        rpc.ToAtomic(amt, 1),
		}
	}

	t := []dero.Transfer{t1}
	fee := rpc.GasEstimate(round.Contract, "[Holdero]", args, t, rpc.LowLimitFee)
	params := &dero.Transfer_Params{
		Transfers: t,
		SC_ID:     round.Contract,
		SC_RPC:    args,
		Ringsize:  2,
		Fees:      fee,
	}

	if err := rpcClientW.CallFor(ctx, &txid, "transfer", params); err != nil {
		rpc.PrintError("[Holdero] Bet: %s", err)
		return
	}

	if f, err := strconv.ParseFloat(amt, 64); err == nil {
		updateStatsWager(f)
	}

	round.display.results = ""
	signals.placedBet = true
	rpc.PrintLog("[Holdero] Bet TX: %s", txid)

	return txid.TXID
}

// Holdero check and fold
func Check() (tx string) {
	rpcClientW, ctx, cancel := rpc.SetWalletClient(rpc.Wallet.Rpc, rpc.Wallet.UserPass)
	defer cancel()

	arg1 := dero.Argument{Name: "entrypoint", DataType: "S", Value: "Bet"}
	args := dero.Arguments{arg1}
	txid := dero.Transfer_Result{}

	var t1 dero.Transfer
	if !round.asset {
		t1 = dero.Transfer{
			Destination: "dero1qyr8yjnu6cl2c5yqkls0hmxe6rry77kn24nmc5fje6hm9jltyvdd5qq4hn5pn",
			Amount:      0,
			Burn:        0,
		}
	} else {
		if round.tourney {
			t1 = dero.Transfer{
				SCID:        crypto.HashHexToHash(TourneySCID),
				Destination: "dero1qyr8yjnu6cl2c5yqkls0hmxe6rry77kn24nmc5fje6hm9jltyvdd5qq4hn5pn",
				Burn:        0,
			}
		} else {
			t1 = rpc.GetAssetSCIDforTransfer(0, round.assetID)
			if t1.Destination == "" {
				rpc.PrintError("[Holdero] Check/Fold: err getting asset SCID for transfer")
				return
			}
		}
	}

	t := []dero.Transfer{t1}
	fee := rpc.GasEstimate(round.Contract, "[Holdero]", args, t, rpc.LowLimitFee)
	params := &dero.Transfer_Params{
		Transfers: t,
		SC_ID:     round.Contract,
		SC_RPC:    args,
		Ringsize:  2,
		Fees:      fee,
	}

	if err := rpcClientW.CallFor(ctx, &txid, "transfer", params); err != nil {
		rpc.PrintError("[Holdero] Check/Fold: %s", err)
		return
	}

	round.display.results = ""
	rpc.PrintLog("[Holdero] Check/Fold TX: %s", txid)

	return txid.TXID
}

// Holdero single winner payout
//   - w defines which player the pot is going to
func PayOut(w string) string {
	rpcClientW, ctx, cancel := rpc.SetWalletClient(rpc.Wallet.Rpc, rpc.Wallet.UserPass)
	defer cancel()

	arg1 := dero.Argument{Name: "entrypoint", DataType: "S", Value: "Winner"}
	arg2 := dero.Argument{Name: "whoWon", DataType: "S", Value: w}
	args := dero.Arguments{arg1, arg2}
	txid := dero.Transfer_Result{}

	t1 := dero.Transfer{
		Destination: "dero1qyr8yjnu6cl2c5yqkls0hmxe6rry77kn24nmc5fje6hm9jltyvdd5qq4hn5pn",
		Amount:      0,
		Burn:        0,
	}

	t := []dero.Transfer{t1}
	fee := rpc.GasEstimate(round.Contract, "[Holdero]", args, t, rpc.LowLimitFee)
	params := &dero.Transfer_Params{
		Transfers: t,
		SC_ID:     round.Contract,
		SC_RPC:    args,
		Ringsize:  2,
		Fees:      fee,
	}

	if err := rpcClientW.CallFor(ctx, &txid, "transfer", params); err != nil {
		rpc.PrintError("[Holdero] Payout: %s", err)
		return ""
	}

	rpc.PrintLog("[Holdero] Payout TX: %s", txid)

	return txid.TXID
}

// Holdero split winners payout
//   - Pass in ranker from hand and folded bools to determine split
func PayoutSplit(r ranker, f1, f2, f3, f4, f5, f6 bool) string {
	rpcClientW, ctx, cancel := rpc.SetWalletClient(rpc.Wallet.Rpc, rpc.Wallet.UserPass)
	defer cancel()

	ways := 0
	splitWinners := [6]string{"Zero", "Zero", "Zero", "Zero", "Zero", "Zero"}

	if r.p1HighCardArr[0] > 0 && !f1 {
		ways = 1
		splitWinners[0] = "Player1"
	}

	if r.p2HighCardArr[0] > 0 && !f2 {
		ways++
		splitWinners[1] = "Player2"
	}

	if r.p3HighCardArr[0] > 0 && !f3 {
		ways++
		splitWinners[2] = "Player3"
	}

	if r.p4HighCardArr[0] > 0 && !f4 {
		ways++
		splitWinners[3] = "Player4"
	}

	if r.p5HighCardArr[0] > 0 && !f5 {
		ways++
		splitWinners[4] = "Player5"
	}

	if r.p6HighCardArr[0] > 0 && !f6 {
		ways++
		splitWinners[5] = "Player6"
	}

	sort.Strings(splitWinners[:])

	arg1 := dero.Argument{Name: "entrypoint", DataType: "S", Value: "SplitWinner"}
	arg2 := dero.Argument{Name: "div", DataType: "U", Value: ways}
	arg3 := dero.Argument{Name: "split1", DataType: "S", Value: splitWinners[0]}
	arg4 := dero.Argument{Name: "split2", DataType: "S", Value: splitWinners[1]}
	arg5 := dero.Argument{Name: "split3", DataType: "S", Value: splitWinners[2]}
	arg6 := dero.Argument{Name: "split4", DataType: "S", Value: splitWinners[3]}
	arg7 := dero.Argument{Name: "split5", DataType: "S", Value: splitWinners[4]}
	arg8 := dero.Argument{Name: "split6", DataType: "S", Value: splitWinners[5]}

	args := dero.Arguments{arg1, arg2, arg3, arg4, arg5, arg6, arg7, arg8}
	txid := dero.Transfer_Result{}

	t1 := dero.Transfer{
		Destination: "dero1qyr8yjnu6cl2c5yqkls0hmxe6rry77kn24nmc5fje6hm9jltyvdd5qq4hn5pn",
		Amount:      0,
		Burn:        0,
	}

	t := []dero.Transfer{t1}
	fee := rpc.GasEstimate(round.Contract, "[Holdero]", args, t, rpc.LowLimitFee)
	params := &dero.Transfer_Params{
		Transfers: t,
		SC_ID:     round.Contract,
		SC_RPC:    args,
		Ringsize:  2,
		Fees:      fee,
	}

	if err := rpcClientW.CallFor(ctx, &txid, "transfer", params); err != nil {
		rpc.PrintError("[Holdero] Split Winner: %s", err)
		return ""
	}

	rpc.PrintLog("[Holdero] Split Winner TX: %s", txid)

	return txid.TXID
}

// Reveal Holdero hand key for showdown
func RevealKey(key string) {
	time.Sleep(6 * time.Second)
	rpcClientW, ctx, cancel := rpc.SetWalletClient(rpc.Wallet.Rpc, rpc.Wallet.UserPass)
	defer cancel()

	arg1 := dero.Argument{Name: "entrypoint", DataType: "S", Value: "RevealKey"}
	arg2 := dero.Argument{Name: "pcSeed", DataType: "H", Value: key}
	args := dero.Arguments{arg1, arg2}
	txid := dero.Transfer_Result{}

	t1 := dero.Transfer{
		Destination: "dero1qyr8yjnu6cl2c5yqkls0hmxe6rry77kn24nmc5fje6hm9jltyvdd5qq4hn5pn",
		Amount:      0,
		Burn:        0,
	}

	t := []dero.Transfer{t1}
	fee := rpc.GasEstimate(round.Contract, "[Holdero]", args, t, rpc.LowLimitFee)
	params := &dero.Transfer_Params{
		Transfers: t,
		SC_ID:     round.Contract,
		SC_RPC:    args,
		Ringsize:  2,
		Fees:      fee,
	}

	if err := rpcClientW.CallFor(ctx, &txid, "transfer", params); err != nil {
		rpc.PrintError("[Holdero] Reveal: %s", err)
		return
	}

	round.display.results = ""
	rpc.PrintLog("[Holdero] Reveal TX: %s", txid)
}

// Owner can shuffle deck for Holdero, clean above 0 can retrieve balance
func CleanTable(amt uint64) {
	rpcClientW, ctx, cancel := rpc.SetWalletClient(rpc.Wallet.Rpc, rpc.Wallet.UserPass)
	defer cancel()

	arg1 := dero.Argument{Name: "entrypoint", DataType: "S", Value: "CleanTable"}
	arg2 := dero.Argument{Name: "amount", DataType: "U", Value: amt}
	args := dero.Arguments{arg1, arg2}
	txid := dero.Transfer_Result{}

	t1 := dero.Transfer{
		Destination: "dero1qyr8yjnu6cl2c5yqkls0hmxe6rry77kn24nmc5fje6hm9jltyvdd5qq4hn5pn",
		Amount:      0,
		Burn:        0,
	}

	t := []dero.Transfer{t1}
	fee := rpc.GasEstimate(round.Contract, "[Holdero]", args, t, rpc.LowLimitFee)
	params := &dero.Transfer_Params{
		Transfers: t,
		SC_ID:     round.Contract,
		SC_RPC:    args,
		Ringsize:  2,
		Fees:      fee,
	}

	if err := rpcClientW.CallFor(ctx, &txid, "transfer", params); err != nil {
		rpc.PrintError("[Holdero] Clean Table: %s", err)
		return
	}

	rpc.PrintLog("[Holdero] Clean Table TX: %s", txid)
}

// Owner can timeout a player at Holdero table
func TimeOut() {
	rpcClientW, ctx, cancel := rpc.SetWalletClient(rpc.Wallet.Rpc, rpc.Wallet.UserPass)
	defer cancel()

	arg1 := dero.Argument{Name: "entrypoint", DataType: "S", Value: "TimeOut"}
	args := dero.Arguments{arg1}
	txid := dero.Transfer_Result{}

	t1 := dero.Transfer{
		Destination: "dero1qyr8yjnu6cl2c5yqkls0hmxe6rry77kn24nmc5fje6hm9jltyvdd5qq4hn5pn",
		Amount:      0,
		Burn:        0,
	}

	t := []dero.Transfer{t1}
	fee := rpc.GasEstimate(round.Contract, "[Holdero]", args, t, rpc.LowLimitFee)
	params := &dero.Transfer_Params{
		Transfers: t,
		SC_ID:     round.Contract,
		SC_RPC:    args,
		Ringsize:  2,
		Fees:      fee,
	}

	if err := rpcClientW.CallFor(ctx, &txid, "transfer", params); err != nil {
		rpc.PrintError("[Holdero] Timeout: %s", err)
		return
	}

	rpc.PrintLog("[Holdero] Timeout TX: %s", txid)
}

// Owner can force start a Holdero table with empty seats
func ForceStat() {
	rpcClientW, ctx, cancel := rpc.SetWalletClient(rpc.Wallet.Rpc, rpc.Wallet.UserPass)
	defer cancel()

	arg1 := dero.Argument{Name: "entrypoint", DataType: "S", Value: "ForceStart"}
	args := dero.Arguments{arg1}
	txid := dero.Transfer_Result{}

	t1 := dero.Transfer{
		Destination: "dero1qyr8yjnu6cl2c5yqkls0hmxe6rry77kn24nmc5fje6hm9jltyvdd5qq4hn5pn",
		Amount:      0,
		Burn:        0,
	}

	t := []dero.Transfer{t1}
	fee := rpc.GasEstimate(round.Contract, "[Holdero]", args, t, rpc.LowLimitFee)
	params := &dero.Transfer_Params{
		Transfers: t,
		SC_ID:     round.Contract,
		SC_RPC:    args,
		Ringsize:  2,
		Fees:      fee,
	}

	if err := rpcClientW.CallFor(ctx, &txid, "transfer", params); err != nil {
		rpc.PrintError("[Holdero] Force Start: %s", err)
		return
	}

	rpc.PrintLog("[Holdero] Force Start TX: %s", txid)
}

// Share asset url at Holdero table
//   - face and back are the names of assets
//   - faceUrl and backUrl are the Urls for those assets
func SharedDeckUrl(face, faceUrl, back, backUrl string) {
	rpcClientW, ctx, cancel := rpc.SetWalletClient(rpc.Wallet.Rpc, rpc.Wallet.UserPass)
	defer cancel()

	var cards CardSpecs
	if face != "" && back != "" {
		cards.Faces.Name = face
		cards.Faces.Url = faceUrl
		cards.Backs.Name = back
		cards.Backs.Url = backUrl
	}

	if mar, err := json.Marshal(cards); err == nil {
		arg1 := dero.Argument{Name: "entrypoint", DataType: "S", Value: "Deck"}
		arg2 := dero.Argument{Name: "face", DataType: "S", Value: string(mar)}
		arg3 := dero.Argument{Name: "back", DataType: "S", Value: "nil"}
		args := dero.Arguments{arg1, arg2, arg3}
		txid := dero.Transfer_Result{}

		t1 := dero.Transfer{
			Destination: "dero1qyr8yjnu6cl2c5yqkls0hmxe6rry77kn24nmc5fje6hm9jltyvdd5qq4hn5pn",
			Amount:      0,
			Burn:        0,
		}

		t := []dero.Transfer{t1}
		fee := rpc.GasEstimate(round.Contract, "[Holdero]", args, t, rpc.LowLimitFee)
		params := &dero.Transfer_Params{
			Transfers: t,
			SC_ID:     round.Contract,
			SC_RPC:    args,
			Ringsize:  2,
			Fees:      fee,
		}

		if err := rpcClientW.CallFor(ctx, &txid, "transfer", params); err != nil {
			rpc.PrintError("[Holdero] Shared: %s", err)
			return
		}

		rpc.PrintLog("[Holdero] Shared TX: %s", txid)
	}
}

// Deposit tournament chip bal with name to leader board SC
func TourneyDeposit(bal uint64, name string) {
	if bal > 0 {
		rpcClientW, ctx, cancel := rpc.SetWalletClient(rpc.Wallet.Rpc, rpc.Wallet.UserPass)
		defer cancel()

		arg1 := dero.Argument{Name: "entrypoint", DataType: "S", Value: "Deposit"}
		arg2 := dero.Argument{Name: "name", DataType: "S", Value: name}
		args := dero.Arguments{arg1, arg2}
		txid := dero.Transfer_Result{}

		t1 := dero.Transfer{
			SCID:        crypto.HashHexToHash(TourneySCID),
			Destination: "dero1qyr8yjnu6cl2c5yqkls0hmxe6rry77kn24nmc5fje6hm9jltyvdd5qq4hn5pn",
			Amount:      0,
			Burn:        bal,
		}

		t := []dero.Transfer{t1}
		fee := rpc.GasEstimate(TourneySCID, "[Holdero]", args, t, rpc.LowLimitFee)
		params := &dero.Transfer_Params{
			Transfers: t,
			SC_ID:     TourneySCID,
			SC_RPC:    args,
			Ringsize:  2,
			Fees:      fee,
		}

		if err := rpcClientW.CallFor(ctx, &txid, "transfer", params); err != nil {
			rpc.PrintError("[Holdero] Tournament Deposit: %s", err)
			return
		}

		rpc.PrintLog("[Holdero] Tournament Deposit TX: %s", txid)

	} else {
		rpc.PrintError("[Holdero] Tournament Deposit: no chips")
	}
}

// Code latest SC code for Holdero public or private SC
//   - version defines which type of Holdero contract
//   - 0 for standard public
//   - 1 for standard private
//   - 2 for HGC
func GetHolderoCode(version int) string {
	if rpc.Daemon.IsConnected() {
		rpcClientD, ctx, cancel := rpc.SetDaemonClient(rpc.Daemon.Rpc)
		defer cancel()

		var result *dero.GetSC_Result
		var params dero.GetSC_Params
		switch version {
		case 0:
			params = dero.GetSC_Params{
				SCID:      HolderoSCID,
				Code:      true,
				Variables: false,
			}
		case 1:
			params = dero.GetSC_Params{
				SCID:      pHolderoSCID,
				Code:      true,
				Variables: false,
			}
		case 2:
			params = dero.GetSC_Params{
				SCID:      HGCHolderoSCID,
				Code:      true,
				Variables: false,
			}
		default:

		}

		if err := rpcClientD.CallFor(ctx, &result, "DERO.GetSC", params); err != nil {
			logger.Errorln("[GetHolderoCode]", err)
			return ""
		}

		return result.Code

	}

	return ""
}

var unlockFee = uint64(300000)

// Contract unlock transfer
func OwnerT3(o bool) (t *dero.Transfer) {
	if o {
		t = &dero.Transfer{
			Destination: "dero1qyr8yjnu6cl2c5yqkls0hmxe6rry77kn24nmc5fje6hm9jltyvdd5qq4hn5pn",
			Amount:      0,
		}
	} else {
		if fee, ok := rpc.GetStringKey(rpc.RatingSCID, "ContractUnlock", rpc.Daemon.Rpc).(float64); ok {
			unlockFee = uint64(fee)
		} else {
			logger.Println("[FetchFees] Could not get current contract unlock fee, using default")
		}

		t = &dero.Transfer{
			Destination: "dero1qyr8yjnu6cl2c5yqkls0hmxe6rry77kn24nmc5fje6hm9jltyvdd5qq4hn5pn",
			Amount:      unlockFee,
		}
	}

	return
}

// Install new Holdero SC
//   - pub defines public or private SC
func uploadHolderoContract(pub int) {
	if rpc.IsReady() {
		rpcClientW, ctx, cancel := rpc.SetWalletClient(rpc.Wallet.Rpc, rpc.Wallet.UserPass)
		defer cancel()

		code := GetHolderoCode(pub)
		if code == "" {
			rpc.PrintError("[Holdero] Upload: could not get SC code")
			return
		}

		args := dero.Arguments{}
		txid := dero.Transfer_Result{}

		params := &dero.Transfer_Params{
			Transfers: []dero.Transfer{*OwnerT3(table.owner.valid)},
			SC_Code:   code,
			SC_Value:  0,
			SC_RPC:    args,
			Ringsize:  2,
		}

		if err := rpcClientW.CallFor(ctx, &txid, "transfer", params); err != nil {
			rpc.PrintError("[Holdero] Upload: %s", err)
			return
		}

		rpc.PrintLog("[Holdero] Upload TX: %s", txid)
	}
}
