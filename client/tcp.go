package main

import (
	"fmt"
	"time"
	"strings"
	"net"
	"math"
	"bufio"
	"strconv"
	"errors"
)

// Inits TCP listener (server), accepts new connections and concurently proceses them
func initListener() (net.Listener, bool){
	listener, err := net.Listen("tcp", "localhost:0")
    if err != nil {
		fmt.Println("Listener creation error: ", err)
		return listener, false
	}

	// Inits tcp listen
	go  func(){
		for {
			conn, err := listener.Accept()

			if err != nil {
				fmt.Println("Error:", err)
			}

			go handleTCPRequest(conn)
		}
	}()

	return listener, true
}

// Inits TCP client and performs TCP connection
func initTCPClient() error{
	ip, port := getNeighbour()

	if ip == ""{
		return errors.New("Neighbour senzor not found")
	}
	
	conn, err := net.Dial("tcp", ip + ":" + port)
	if err != nil{
		return err
	}

	for {
		// Reads from file
		firstReading := strings.Split(inputText[getTime() + 2], ",")
		fmt.Printf("Senzor read - temperature: %s, pressure: %s, humidity: %s, CO: %s, NO2: %s, SO2: %s\n",firstReading[0], firstReading[1], firstReading[2], firstReading[3], firstReading[4], firstReading[5])

		// Connects to neighbour
		_, err = conn.Write([]byte("umjeravanje"))
		if err != nil{
			return err
		}

		fmt.Println("Fetching data from neighbour senzor")

		// Reads neighbour response
		message, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil{
			return err
		}

		// Formats response
		secondReading := strings.Split(message, ",")
		fmt.Printf("Data from other senzor - temperature: %s, pressure: %s, humidity: %s, CO: %s, NO2: %s, SO2: %s\n\n", secondReading[0], secondReading[1], secondReading[2], secondReading[3], secondReading[4], secondReading[5])

	    var finalReading [6]string
		for index, _ := range finalReading{
			val1, _ := strconv.Atoi(firstReading[index])
			val2, _ := strconv.Atoi(secondReading[index])

			if firstReading[index] == "" || secondReading[index] == ""{
				finalReading[index] = strconv.FormatFloat(float64(val1 + val2), 'f', 1, 64)
			}else{
				finalReading[index] = strconv.FormatFloat(float64(val1 + val2)/2.0, 'f', 1, 64)
			}
		}

		// Prints final resoult
		fmt.Printf("Final data - temperature: %s, pressure: %s, humidity: %s, CO: %s, NO2: %s, SO2: %s\n\n", finalReading[0], finalReading[1], finalReading[2], finalReading[3], finalReading[4], finalReading[5])

		// Sending measurements
		if err = sendMeasurements(finalReading); err != nil{
			fmt.Println("Measurement storage failed: ", err)
		}

		time.Sleep(15 * time.Second)
	}
}

// TCP request handler
func handleTCPRequest(conn net.Conn){
	defer conn.Close()

	fmt.Println("\nRecived new TCP connection")
	for {
		inBuffer := make([]byte, 11)
		_, err := conn.Read(inBuffer)
		if err != nil{
			fmt.Println("Server TCP read error: ", err)
			return
		}

		fmt.Println("Processing request")

		if strings.Compare(string(inBuffer), "umjeravanje") != 0{
			fmt.Println("TCP request not recognized")
			return
		}

		outBuffer := []byte(inputText[getTime() + 2] + "\n")
		_, err = conn.Write(outBuffer)
		if err != nil{
			fmt.Println("Server TCP write error: ", err)
			return
		}

		fmt.Println("Data sent\n")
	}
}

func getTime() int{
	duration := time.Since(startTime)
	math.Floor(duration.Seconds())

	return int(math.Floor(duration.Seconds()))
}
