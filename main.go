package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"

	"github.com/albenik/go-serial"
)

type Data struct {
	Id     uint8
	Val    float32
	Day    uint8
	Month  uint8
	Year   uint16
	Hour   uint8
	Minute uint8
	Second uint8
}

func scanLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, nil
}

func bfloat32(bytes []byte) float32 {
	bits := binary.LittleEndian.Uint32(bytes)
	float := math.Float32frombits(bits)
	return float
}

func bint16(bytes []byte) uint16 {
	var value uint16
	value |= uint16(bytes[0]) //Little endian
	value |= uint16(bytes[1]) << 8
	return value
}

func buffCorrect(buff []byte) bool {
	return buff[0] == 255 && len(buff) == 13 && buff[12] <= 60
}

func saveData(data *Data, buff []byte) {
	data.Id = buff[1]
	bVal := []byte{buff[2], buff[3], buff[4], buff[5]}
	data.Val = bfloat32(bVal)
	data.Day = buff[6]
	data.Month = buff[7]
	bYear := []byte{buff[8], buff[9]}
	data.Year = bint16(bYear)
	data.Hour = buff[10]
	data.Minute = buff[11]
	data.Second = buff[12]
}

func postData(telemetry map[string]string) {
	json_data, _ := json.Marshal(telemetry)
	http.Post(TbUrl+TbToken+"/telemetry", "application/json", bytes.NewBuffer(json_data))
}

var TbToken string
var TbUrl string

func init() {
	//Load tokens in memory
	lines, err := scanLines("./tokens.txt")
	if err != nil {
		panic(err)
	}
	TbUrl = lines[0]
	TbToken = lines[1]
}

func main() {
	data := Data{}
	sensorId := [4]string{"temperature", "turbidity", "conductivity", "acidity"}

	// Retrieve the port list
	ports, err := serial.GetPortsList()
	if err != nil {
		log.Fatal(err)
	}
	if len(ports) == 0 {
		log.Fatal("No serial ports found!")
	}

	mode := &serial.Mode{
		BaudRate: 115200,
		Parity:   serial.NoParity,
		DataBits: 8,
		StopBits: serial.OneStopBit,
	}
	port, err := serial.Open(ports[0], mode)
	if err != nil {
		log.Fatal(err)
	}

	//Make the buffer as big as the payload
	buff := make([]byte, 13)
	//Infinite loop
	for {
		for {
			n, err := port.Read(buff)
			if err == nil && n != 0 && buffCorrect(buff) {
				saveData(&data, buff)
				postData(map[string]string{sensorId[data.Id]: fmt.Sprintf("%v", data.Val)})
				fmt.Println(data)
			}
		}
	}
}
