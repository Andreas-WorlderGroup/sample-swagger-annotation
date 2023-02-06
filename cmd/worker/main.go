package main

import (
	sensor "andreas/internal/handler"
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/go-co-op/gocron"
)

func insertToBackend() {
	rand.Seed(time.Now().UnixNano())
	letter := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	now := time.Now().Format("Mon 01/02/2006-15:04:05")

	data := sensor.SensorData{
		Value:     rand.Intn(100-0) + 0,
		ID1:       rand.Intn(10-0) + 0,
		ID2:       string(letter[rand.Intn(len(letter))]),
		Timestamp: now,
	}

	bdata, _ := json.Marshal(data)

	request, err := http.NewRequest("POST", "http://localhost:8080/data", bytes.NewBuffer(bdata))
	if err != nil {
		fmt.Println("error request ", err)
	}
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")

	client := &http.Client{}
	response, error := client.Do(request)
	if error != nil {
		panic(error)
	}

	fmt.Println(data)

	defer response.Body.Close()
}

func main() {
	// test code
	// loc, _ := time.LoadLocation("Asia/Jakarta")
	// fmt.Println(time.Unix(1670157149, 0))
	// tm, _ := time.ParseInLocation("Mon 01/02/2006-15:04:05", "Sun 12/04/2022-19:32:29", loc)
	// fmt.Println(tm.Unix())

	// Worker Scheduler for inserting data to backend via api
	var cron = gocron.NewScheduler(time.Local)
	cron.Every(1).Seconds().Tag("insert_request").Do(insertToBackend)
	cron.StartBlocking()
}
