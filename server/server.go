package main

import (
	"net"
	"net/http"
	"fmt"
	"strconv"
	"encoding/json"
	"math"
	"time"
)

type Senzor struct {
	Name string
	Latitude float64
	Longitude float64
	IPaddress string
	Port int

	Measurements map[string]string
}

type LogEntry struct {
	Time time.Time
	Request string
	Data string
	Response string
}

var senzors map[string]Senzor

var log []LogEntry

/* /register POST

Registers senzor and stores it to senzors map

Post params:

username String		Senzor name - identifier
latitude Float 
longitude Float
ip String 			IP address of senzor, must be valid IP
port Int 			Senzor port

Returns:

Success
200 
{
	"message":"Senzor successfully added"	
}


Failure
4xx
{
	"message":Error message
}
*/
func register(w http.ResponseWriter, r *http.Request) {
	// Checks post params
	if  err, msg := validateRequest(r, []string{"username", "latitude", "longitude", "ip", "port"}); !err{
		writeResponse(w, 400, map[string]string{"message":"Field '" + msg + "' is required"})
		return
	}

	// Checks types
	lat, err := strconv.ParseFloat(r.FormValue("latitude"), 64)
	long, err := strconv.ParseFloat(r.FormValue("longitude"), 64)
	port, err := strconv.Atoi(r.FormValue("port"))

	if err != nil{
		w.WriteHeader(400)
		fmt.Fprintf(w, "{\"message\":\"Wrong parameter type\"}")
		return
	}

	// Checks ip
	ip := net.ParseIP(r.FormValue("ip"))
	if ip.To4() == nil {
		writeResponse(w, 400, map[string]string{"message":"Malformed IP"})
		return
    }

	// Inits senzor map
	if senzors == nil{
		senzors = make(map[string]Senzor)
	}

	// Stores new senzor in senzor map
	senzors[r.FormValue("username")] = Senzor{r.FormValue("username"), lat, long, r.FormValue("ip"), port, make(map[string]string)}

	// On success writes a response
	writeResponse(w, 200, map[string]string{"message":"Senzor successfully added"})
}


/* /searchNeighbour POST

Searches all senzors and returns closest one

Post params:

username String		Senzor name - identifier

Returns:

Success
200 
{
	"ip":Neighbours ip adress
	"port":Neighbours port
}


Failure
4xx
{
	"message":Error message
}
*/
func searchNeighbour(w http.ResponseWriter, r *http.Request) {
	if err, msg := validateRequest(r, []string{"username"}); !err{
		writeResponse(w, 400, map[string]string{"message":"Field '" + msg + "' is required"})
		return
	}

	// Looks for senzor with given name. If found searches closest neighbour
	if senzor, ok := senzors[r.FormValue("username")]; ok {
		if len(senzors) < 2{
			writeResponse(w, 400, map[string]string{"message":"Zero other seznors registered."})
			return
		}

		// Calculates closest senzor
		min := -1.0
		calculations := "Distance calculated - "
		var closestSenzor Senzor

		for key, value := range senzors{
			if key != senzor.Name{
				dlat := value.Latitude - senzor.Latitude
				dlon := value.Longitude - senzor.Longitude

				a := math.Pow((math.Sin(dlat/2)),2) + (math.Cos(senzor.Latitude) * math.Cos(value.Latitude) * math.Pow((math.Sin(dlon/2)),2))
				c := 2 * math.Atan2( math.Sqrt(a), math.Sqrt(1-a))
				d := 6373 * c

				calculations += key + ": " + strconv.FormatFloat(d, 'f', 3, 64) + ", "
				if min == -1.0 || d < min{
					min = d
					closestSenzor = value
				}
			}
		}

		calculations += "value chosen: " + closestSenzor.Name

		log = append(log, LogEntry{time.Now(), "searchNeighbour for name " + r.FormValue("username"), calculations, "ip: " + closestSenzor.IPaddress + ", port: " + strconv.Itoa(closestSenzor.Port)})

		writeResponse(w, 200, map[string]string{"ip": closestSenzor.IPaddress, "port":strconv.Itoa(closestSenzor.Port)})
		return
	}else{
		writeResponse(w, 400, map[string]string{"message":"Senzor with name " + r.FormValue("username") +  " does not exist"})
		return
	}
}

/* /register POST

Registers senzor and stores it to senzors map

Post params:

username String		Senzor name - identifier
temperature String
pressure String
humidity String
CO String
NO2 String
SO2 String

Returns:

Success
200 
{
	"message":"Measurements for senzor with name a successfully added"
}


Failure
4xx
{
	"message":Error message
}
*/
func storeMeasurement(w http.ResponseWriter, r *http.Request) {
	postParams := []string{"username", "temperature", "pressure", "humidity", "CO", "NO2", "SO2"}

	// Checks post params
	if  err, msg := validateRequest(r, postParams); !err{
		writeResponse(w, 400, map[string]string{"message":"Field '" + msg + "' is required"})
		return
	}

	// Looks for senzor with given name. If found searches closest neighbour
	if senzor, ok := senzors[r.FormValue("username")]; ok {
		// Calculates closest senzor

		message := "Data stored - "
		for _, element := range postParams{
			if element != "username"{
				message += element + ": " + r.FormValue(element) + ", "
				senzor.Measurements[element] = r.FormValue(element)
			}
		}

		log = append(log, LogEntry{time.Now(), "storeMeasurement for name " + r.FormValue("username"), message, "Measurements for senzor with name '" + r.FormValue("username") + "' successfully added"})

		writeResponse(w, 200, map[string]string{"message":"Measurements for senzor with name '" + r.FormValue("username") + "' successfully added"})
		return
	}else{
		writeResponse(w, 400, map[string]string{"message":"Senzor with name '" + r.FormValue("username") +  "' is not registered"})
		return
	}
}

// Get server log
// /getLog
func getLog(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	jsonRes, _ := json.Marshal(log)
	fmt.Fprintf(w,"%s", jsonRes)

}

func main() {
	log = []LogEntry{}

    http.HandleFunc("/register", register)
    http.HandleFunc("/searchNeighbour", searchNeighbour)
    http.HandleFunc("/storeMeasurement", storeMeasurement)
	http.HandleFunc("/getLog", getLog)

    http.ListenAndServe(":8888", nil)
}



func validateRequest(r *http.Request, keys []string) (bool, string){

	for key := range keys{
		if r.FormValue(keys[key]) == ""{
			return false, keys[key]
		}
	}

	return true, ""
}

func writeResponse(w http.ResponseWriter, statusCode int, response map[string]string){
	w.WriteHeader(statusCode)
	jsonRes, _ := json.Marshal(response)
	fmt.Fprintf(w,"%s", jsonRes)
}
