package main

import (
	"bufio"
	"github.com/tarm/goserial"
	"log"
	"strconv"
	"strings"
)

func main() {
	c := &serial.Config{Name: "/dev/cu.usbserial-A800ftEg", Baud: 9600}
	s, err := serial.OpenPort(c)
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(s)
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
			}
			log.Println("Frequency: " + strconv.FormatFloat(frequency, 'f', 10,
				32))
		} else {
			log.Println("Received unknown data: " + line)
		}
	}
}
