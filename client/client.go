package main

import (
	"fmt"
	"net"
	"math/rand"
	"time"
	"io/ioutil"
	"strings"
	"bufio"
	"os"
	"errors"
)

var clientName string
var lat, long float64
var port int
var ip net.IP

var startTime time.Time

var inputText []string

var listener net.Listener

func main(){

	startTime = time.Now()

	if err := prepInputFile(); err != nil{
		fmt.Println("Input file rading error: ", err)
		return
	}

	if err := initSenzor(); err != nil{
		fmt.Println("Error while initializing senzor: ", err)
        return
	}else{
		fmt.Printf("Senzor initialized. Name: %s, lat: %f, long: %f, ip: %s, port: %d\n\n", clientName, lat, long, ip, port)
	}

	if err := registerSenzor(); err != nil{
		fmt.Println("Senzor registration failed: ", err)
		return
	}

	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Press any key to start data reading...\n")
		reader.ReadString('\n')

		if err := initTCPClient(); err != nil{
			fmt.Println("Failed to init TCP client: ", err)
		}
	}

	defer listener.Close()
}

// Opens and reads input file
func prepInputFile() error {
	dat, err := ioutil.ReadFile("mjerenja.csv")
    if err != nil {
        return err
	}

    inputText = strings.Split(string(dat),"\n")

	return nil
}

// Inits senzor with needed data
func initSenzor() error{
	// Init client
	rand.Seed( time.Now().UTC().UnixNano())

	clientName = RandStringBytes(10)

	lat = (rand.Float64() - 0.5) * 180
	long = (rand.Float64() - 0.5) * 360

	// Init client tcp port
	if listener, err := initListener(); !err{
		return errors.New("Failed to init TCP listener")
	}else{
		ip = listener.Addr().(*net.TCPAddr).IP
		port = listener.Addr().(*net.TCPAddr).Port
	}

	return nil
}

// Random string generator
const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandStringBytes(n int) string {
    b := make([]byte, n)
    for i := range b {
        b[i] = letterBytes[rand.Intn(len(letterBytes))]
    }
    return string(b)
}
