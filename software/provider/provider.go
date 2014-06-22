package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/tarm/goserial"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

var readingChannel = make(chan float64)
var extract_wg sync.WaitGroup
var pusher_wg sync.WaitGroup

type Reading struct {
	Timestamp time.Time
	Value     float64
}

type ErrorMessage struct {
	Id      string
	Message string
}

func extract_readings(r io.Reader) {
	defer extract_wg.Done()
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
		elements := strings.Split(line, ";")
		if elements[0] == "F" {
			frequency, err := strconv.ParseFloat(elements[1], 32)
			if err != nil {
				log.Println("Received broken frequency: " + line)
				continue
			}
			if frequency < 48 || frequency > 52 {
				log.Println("Frequency out of plausible range: " + line)
				continue
			}
			readingChannel <- frequency
		} else if elements[0] == "I" {
			log.Println("Info message: " + line)
		} else {
			log.Println("Received unknown data: " + line)
		}
	}
}

func pusher() {
	defer pusher_wg.Done()
	for frequency := range readingChannel {
		log.Println("Frequency: " + strconv.FormatFloat(frequency, 'f', 5, 32))
		client := &http.Client{}
		body := Reading{time.Now(), frequency}
		bodyBytes, _ := json.Marshal(body)
		reqUrl := fmt.Sprintf("%s/api/submit/%s", Cfg.Network.Host,
			Cfg.API.Meter)
		req, err := http.NewRequest("POST", reqUrl, bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Add("X-API-Key", Cfg.API.Key)
		resp, err := client.Do(req)
		if err != nil {
			log.Println("Error posting data: ", err.Error())
			continue
		}
		defer resp.Body.Close()

		response, rerr := ioutil.ReadAll(resp.Body)
		if rerr != nil {
			log.Println("Error getting post result data: ", err.Error())
			continue
		}
		if resp.StatusCode != http.StatusOK {
			decoder := json.NewDecoder(bytes.NewReader(response))
			var errorMessage ErrorMessage
			err := decoder.Decode(&errorMessage)
			if err != nil {
				log.Println("Failed to decode error message: " + err.Error())
			} else {
				log.Println(resp.StatusCode, errorMessage.Id, ":", errorMessage.Message)
			}
		}
	}
}

var configFile = flag.String("config", "defluxio-provider.conf", "configuration file")

func init() {
	flag.Parse()
	loadConfiguration(*configFile)
}

func main() {
	c := &serial.Config{Name: Cfg.Device.Path, Baud: Cfg.Device.Baudrate}
	s, err := serial.OpenPort(c)
	if err != nil {
		log.Fatal(err)
	}
	// just one reader, since it is a serial connection
	extract_wg.Add(1)
	go extract_readings(s)
	for c := 0; c < runtime.NumCPU(); c++ {
		pusher_wg.Add(1)
		go pusher()
	}

	extract_wg.Wait()
	close(readingChannel)
	pusher_wg.Wait()
}
