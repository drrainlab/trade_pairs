package main

import (
	"encoding/json"
	"flag"
	"log"
	"math/rand"
	"net/http"
	"time"
)

var pairs = []string{
	"NEAR/USDT", "XLM/USDT", "XRP/USDT", "DOGE/USDT", "ADA/USDT", "SOL/USDT",
	"XMR/USDT", "MRK/USDT", "SAND/USDT", "XTZ/USDT", "CRV/USDT", "DASH/USDT",
	"SC/USDT", "AUDIO/USDT", "DODO/USDT",
}

var lastUsedPairs = make([]string, 0)

func getThreeRandomPairs() []string {
	rand.Seed(time.Now().UnixNano())

	availablePairs := make([]string, len(pairs))
	copy(availablePairs, pairs)

	for _, usedPair := range lastUsedPairs {
		for i, availablePair := range availablePairs {
			if usedPair == availablePair {
				availablePairs = append(availablePairs[:i], availablePairs[i+1:]...)
				break
			}
		}
	}

	selectedPairs := make([]string, 0, 3)
	for i := 0; i < 3; i++ {
		index := rand.Intn(len(availablePairs))
		selectedPairs = append(selectedPairs, availablePairs[index])
		availablePairs = append(availablePairs[:index], availablePairs[index+1:]...)
	}

	lastUsedPairs = selectedPairs
	return selectedPairs
}

func handlePairs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	response := struct {
		Pairs         []string `json:"pairs"`
		RefreshPeriod int      `json:"refresh_period"`
	}{
		Pairs:         getThreeRandomPairs(),
		RefreshPeriod: 1800,
	}

	json.NewEncoder(w).Encode(response)
}

func main() {
	port := flag.String("port", "8090", "port to listen on")
	flag.Parse()

	http.HandleFunc("/pairs", handlePairs)
	log.Printf("Server started on port %s\n", *port)
	log.Fatal(http.ListenAndServe(":"+*port, nil))
}
