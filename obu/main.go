package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/gorilla/websocket"
	"github.com/kushal0511-not/toll_calculator/types"
)

const sendInterval = 5 * time.Second
const wsEndpoint = "ws://127.0.0.1:3000/ws"

func genLatLong() (float64, float64) {
	return genCoord(), genCoord()
}
func genCoord() float64 {
	n := rand.Intn(100) + 1
	f := rand.Float64()

	return float64(n) + f
}

func genOBUIDs(n int) []int {
	obusids := make([]int, n)
	for i := 0; i < n; i++ {
		obusids[i] = rand.Intn(100000)
	}
	return obusids
}

func main() {
	obuIDs := genOBUIDs(20)
	conn, _, err := websocket.DefaultDialer.Dial(wsEndpoint, nil)
	if err != nil {
		log.Fatal(err)
	}
	for {
		for i := 0; i < len(obuIDs); i++ {
			lat, long := genLatLong()
			data := types.OBUData{
				OBUID: obuIDs[i],
				Lat:   lat,
				Long:  long,
			}
			fmt.Printf("%+v\n", data)
			if err := conn.WriteJSON(data); err != nil {
				log.Fatal(err)
			}
		}
		time.Sleep(sendInterval)
	}
}

func init() {
	rand.New(rand.NewSource(int64(time.Now().Nanosecond())))
}
