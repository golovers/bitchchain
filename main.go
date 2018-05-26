package main

import (
	"fmt"
	"time"
	"crypto/sha256"
	"encoding/hex"
	"strings"
	"sync"
	"github.com/davecgh/go-spew/spew"
)

type Bitch struct {
	Header  Header
	Bitches []BitchDetail
}

type Header struct {
	Version    int
	Hash       string
	PrevHash   string
	MerkleRoot string
	TimeStamp  time.Time
	Difficulty int
	Nonce      string
}

type BitchDetail struct {
	Version int
	Details string
}

var chain = make([]*Bitch, 0)
var difficulty = 5
var lock = sync.Mutex{}
var diffCircle = 10
var lastDiffCheck = 0

func main() {
	makeTheMotherOfTheBitches()
	// simulate 2 dudes compete each other to find new bitches :))
	go func() {
		for i := 0; i < 2;i++ {
			j := i
			go func() {
				findingBitches(fmt.Sprintf("dude %v", j))
			}()
		}
	}()

	handleDifficultyAdjustment()
	done  := make(chan string)
	<-done
}

// makeTheMotherOfTheBitches create the mother of the bitches...
func makeTheMotherOfTheBitches() {
	genesis := &Bitch{
		Header: Header{
			Version:    1,
			TimeStamp:  time.Now(),
			PrevHash:   "",
			MerkleRoot: "",
		},
		Bitches:[]BitchDetail{BitchDetail{Version:1, Details:"mother of all beautiful bitches :))"}},
	}
	genesis.Header.MerkleRoot = genesis.MerkleRoot()
	genesis.Header.Nonce = proofOfBitch(genesis.Hash(), difficulty)
	genesis.Header.Hash = genesis.Hash()
	chain = append(chain, genesis)
}

// handleDifficultyAdjustment check and ensure the rate of found a new bitch is under control :)
func handleDifficultyAdjustment() {
	t := time.NewTicker(1 * time.Second)
	go func() {
		for _ = range t.C {
			adjustDifficulty(2)
		}
	}()
}

// findingBitches simulate how a dude finds a new beautiful bitch for himself...
func findingBitches(dude string) {
	for {
		bitch := newBitch([]BitchDetail{BitchDetail{Version:1, Details: "beautiful bitch found at "+time.Now().String()}})
		bitch.Header.Nonce = proofOfBitch(bitch.Hash(), difficulty)

		currChain := chain
		tellTheWorldNewBitchFound(dude, currChain, bitch)
	}
}

//	tellTheWorldNewBitchFound tell the world a new bitch was found...
func tellTheWorldNewBitchFound(dude string, currChain []*Bitch, bitch *Bitch) {
	lock.Lock()
	if len(currChain) >= len(chain) {
		chain = currChain
		if isBitchBeautiful(chain[len(chain)-1], bitch) {
			chain = append(chain, bitch)
			fmt.Printf("%v found new bitch...\n", dude)
			spew.Dump(bitch)
		}
	}
	for i := 1; i < len(chain); i++{
		if !isBitchBeautiful(chain[i-1], chain[i]) {
			panic("err")
		}
	}
	lock.Unlock()
}

func newBitch(bitches []BitchDetail) *Bitch {
	bitch := &Bitch{
		Header: Header{
			Difficulty: difficulty,
			Version:    1,
			PrevHash:   chain[len(chain) -1].Hash(),
			TimeStamp:  time.Now(),
		},
		Bitches:bitches,
	}
	bitch.Header.MerkleRoot = bitch.MerkleRoot()
	bitch.Header.Hash = bitch.Hash()
	return bitch
}

// isBitchBeautiful check if the new bitch is beautiful :)
func isBitchBeautiful(old, new *Bitch) bool {
	if new.Header.PrevHash != old.Hash() {
		return false
	}
	blockHash := hashit(new.Hash() + new.Header.Nonce)
	return strings.HasPrefix(blockHash, strings.Repeat("0", new.Header.Difficulty))
}

func (b *Bitch) MerkleRoot() string {
	if len(b.Bitches) == 0 {
		return hashit(hashit(b.Bitches[0].Hash() + b.Bitches[0].Hash()))
	}
	merkle := make([]string, 0)
	for _, b := range b.Bitches {
		merkle = append(merkle, b.Hash())
	}
	for len(merkle) > 1 {
		newMerkle := make([]string, 0)
		if len(merkle) % 2 != 0 {
			merkle = append(merkle, merkle[len(merkle) -1])
		}
		for i := 0; i < len(merkle); i++{
			v := merkle[i] + merkle[i+1]
			newMerkle = append(newMerkle, hashit(v))
			i++
		}
		merkle = newMerkle
	}
	return merkle[0]
}

// proofOfBitch prove a new beautiful bitch is found by given lots of efforts :)
func proofOfBitch(hash string, difficulty int) string{
	prefix := strings.Repeat("0", difficulty)
	for i := 0; ;i++ {
		nonce  := string(fmt.Sprintf("%x", i))
		h := hashit(hash + nonce)
		if strings.HasPrefix(h, prefix) {
			return nonce
		}
	}
	return ""
}

// adjustDifficulty ensure finding new bitch not too fast :)
func adjustDifficulty(threshold float64)  {
	if lastDiffCheck == len(chain) {
		return
	}
	lastDiffCheck = len(chain)
	if len(chain) < diffCircle {
		return
	}
	oldDiff := difficulty
	delta := chain[len(chain)-1].Header.TimeStamp.Sub(chain[len(chain) - diffCircle].Header.TimeStamp)
	if delta.Seconds()/float64(diffCircle) < threshold {
		difficulty = difficulty + 1
	} else {
		difficulty = difficulty - 1
	}
	if oldDiff != difficulty {
		fmt.Printf("adjust difficulty from %v to %v\n", oldDiff, difficulty)
	}
}

func hashit(v string) string {
	h := sha256.New()
	h.Write([]byte(v))
	return hex.EncodeToString(h.Sum(nil))
}

func (b *Bitch) Hash() string {
	return hashit(fmt.Sprintf("%v%v%v%v%v%v", b.Header.MerkleRoot, b.Header.PrevHash, b.Header.TimeStamp, b.Header.Version, b.Header.Difficulty))
}

func (t *BitchDetail) Hash() string {
	return hashit(fmt.Sprintf("%v%v", t.Version, t.Details))
}