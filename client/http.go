package main

import (
	"fmt"
	"net/http"
	"net/url"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"strconv"
	"errors"
)

// Registers server using /register web service
func registerSenzor() error{
	// Post params	
	data := url.Values{}
    data.Add("username", clientName)
    data.Add("latitude", strconv.FormatFloat(lat, 'f', 6, 64))
	data.Add("longitude", strconv.FormatFloat(long, 'f', 6, 64))
	data.Add("port", strconv.Itoa(port))
	data.Add("ip", ip.String())

	// Request
	req, err := http.NewRequest("POST", "http://localhost:8888/register", bytes.NewBufferString(data.Encode()))
	if err != nil {
		return err
	}

	// Content header
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// Sending request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	// Reading body
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil{
			return err
	}

	// Parsing json
	var body map[string]string
	err = json.Unmarshal(bodyBytes, &body)
	if err != nil {
		return err
	}

	// Based on statuc code notifies user
	if resp.Status != "200 OK"{
		return errors.New("Senzor registration failed with code: " + resp.Status)
	}else{
		fmt.Println("Senzor successfuly registrated\n")
	}

	return nil
}



// Fetchec closes neighbour from /searchNeighbour web service
func getNeighbour() (string, string){

	// Post params	
	data := url.Values{}
    data.Add("username", clientName)

	// Request
	req, err := http.NewRequest("POST", "http://localhost:8888/searchNeighbour", bytes.NewBufferString(data.Encode()))
	if err != nil {
		fmt.Println("Failed to fetch senzor neighbour")
		return "", ""
	}

	// Content header
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// Sending request
	client := &http.Client{}
	resp, _ := client.Do(req)

	defer resp.Body.Close()

	// Reading body
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil{
		fmt.Println("Body malformed")
		return "", ""
	}

	// Parsing json
	var body map[string]string
	err = json.Unmarshal(bodyBytes, &body)
	if err != nil {
		fmt.Println("error:", err)
	}

	// Based on statuc code notifies user
	if resp.Status != "200 OK"{
		fmt.Println("Fetch neighbour failed: Error: " + body["message"] + "\n")
		return "", ""
	}else{
		fmt.Printf("Neighbour senzor fetched successfuly. IP: %s, port %s\n\n", body["ip"], body["port"])
	}

	return body["ip"], body["port"]
}

func sendMeasurements(measurementData [6]string) error{
	postParams := []string{"username", "temperature", "pressure", "humidity", "CO", "NO2", "SO2"}

	data := url.Values{}
	for index, element := range postParams{
		if element == "username"{
			data.Add(element, clientName)
		}else{
			data.Add(element, measurementData[index - 1])
		}
	}

	// Request
	req, err := http.NewRequest("POST", "http://localhost:8888/storeMeasurement", bytes.NewBufferString(data.Encode()))
	if err != nil {
		return err
	}

	// Content header
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// Sending request
	client := &http.Client{}
	resp, _ := client.Do(req)

	defer resp.Body.Close()

	// Reading body
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil{
		return err
	}

	// Parsing json
	var body map[string]string
	err = json.Unmarshal(bodyBytes, &body)
	if err != nil {
		return err
	}

	// Based on statuc code notifies user
	if resp.Status != "200 OK"{
		return errors.New("Error: " + body["message"])
	}else{
		fmt.Printf("Measurements for senzor named '" + clientName + "' stored successfuly\n\n")
	}

	return nil
}
