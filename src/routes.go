package main

import (
	"encoding/json"
	"net/http"
)

// --------------------------------------------------------
//   Routing functions
// --------------------------------------------------------

func apiRoute(w http.ResponseWriter, r *http.Request) {

	path := r.URL.Path[5:] // strip api prefix

	switch r.Method {

	case http.MethodPost:

		user := createUser()
		type responseType struct {
			Id  string `json:"id"`
			Key string `json:"key"`
		}
		response := responseType{user.id, user.key}
		sendResponseGood(w, response)

	case http.MethodDelete:

		if validUserFormat(path) == false {
			sendResponseBad(w, "invalid credentials")
			return
		}

		type requestType struct {
			Key string `json:"key"`
		}
		var req requestType
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			sendResponseBad(w, "JSON")
			return
		}
		if validUserKey(req.Key) == false {
			sendResponseBad(w, "invalid credentials")
			return
		}

		if err := deactivateUser(path, req.Key); err != nil {
			sendResponseBad(w, "data-"+err.Error())
			return
		}
		type responseType struct {
			Status string `json:"status"`
		}
		response := responseType{"success"}
		sendResponseGood(w, response)

	case http.MethodPut:

		if validUserFormat(path) == false {
			sendResponseBad(w, "invalid credentials")
			return
		}
		type requestType struct {
			Key string `json:"key"`
			Lon string `json:"lon"`
			Lat string `json:"lat"`
			Acc string `json:"acc"`
		}
		var req requestType
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			sendResponseBad(w, "JSON")
			return
		}
		if validUserKey(req.Key) == false {
			sendResponseBad(w, "invalid credentials")
			return
		}
		if validLocation(req.Lon) == false {
			sendResponseBad(w, "invalid longitude")
			return
		}
		if validLocation(req.Lat) == false {
			sendResponseBad(w, "invalid latitude")
			return
		}
		if validNumber(req.Acc) == false {
			sendResponseBad(w, "invalid accuracy")
			return
		}

		location := UserLocation{req.Lon, req.Lat, req.Acc, 0}
		modTime, err := setUserLocation(path, req.Key, location)
		if err != nil {
			sendResponseBad(w, "data-"+err.Error())
			return
		}
		type responseType struct {
			Mod int64 `json:"mod"`
		}
		response := responseType{modTime}
		sendResponseGood(w, response)

	case http.MethodGet:

		if validUserFormat(path) == false {
			sendResponseBad(w, "invalid credentials")
			return
		}
		user, err := getUserLocation(path)
		if err {
			sendResponseBad(w, "data-expired")
			return
		}

		tagReq := r.Header.Get("If-None-Match")
		if len(tagReq) > 0 {
			const digits = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-_"
			tagIdx := user.mod % 262143 // (1 << 6*3) - 1 == 262143
			etag := string(digits[tagIdx%64]) +
				string(digits[(tagIdx>>6)%64]) +
				string(digits[(tagIdx>>12)%64])
			if tagReq == etag {
				sendResponseUnmodified(w)
				return
			}
		}

		type responseType struct {
			Lon string `json:"lon"`
			Lat string `json:"lat"`
			Acc string `json:"acc"`
			Mod int64  `json:"mod"`
		}
		response := responseType{user.lon, user.lat, user.acc, user.mod}
		sendResponseGood(w, response)

	case http.MethodOptions:

		w.Header().Set("Allow", "OPTIONS, GET, POST, PUT, DELETE")
		w.Header().Set("Cache-Control", "max-age=604800")
		w.Header().Set("Access-Control-Allow-Origin", "https://mapon.me")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, X-Requested-With, x-api-key")
		w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, GET, POST, PUT, DELETE")
		w.WriteHeader(http.StatusOK)

	default:

		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// --------------------------------------------------------
//   Uniform response headers
// --------------------------------------------------------

func sendResponseGood(w http.ResponseWriter, data any) {

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}

func sendResponseBad(w http.ResponseWriter, err string) {

	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte("{\"error\":\"" + err + "\"}"))
}

func sendResponseUnmodified(w http.ResponseWriter) {

	w.WriteHeader(http.StatusNotModified)
	w.Write([]byte(""))
}

// --------------------------------------------------------
//   Validation functions
// --------------------------------------------------------

// Allow exactly 8 characters, [0-9a-Z\-=]
func validUserFormat(s string) bool {

	if len(s) != 8 {
		return false
	}

	for _, c := range s {
		if c == '-' || c == '=' {
			continue
		}
		if c > 47 && c < 58 {
			continue
		}
		if c > 64 && c < 91 {
			continue
		}
		if c > 96 && c < 123 {
			continue
		}
		return false
	}

	return true
}

// Allow exactly 24 characters, [0-9a-Z+/]
func validUserKey(s string) bool {

	if len(s) != 24 {
		return false
	}

	for _, c := range s {
		if c == '+' || c == '/' {
			continue
		}
		if c > 47 && c < 58 {
			continue
		}
		if c > 64 && c < 91 {
			continue
		}
		if c > 96 && c < 123 {
			continue
		}
		return false
	}

	return true
}

// Allow 15 digits, optional minus, required dot
func validLocation(s string) bool {

	if len(s) < 1 {
		return false
	}

	foundDot := false
	supSz := 0
	subSz := 0
	i := 0

	if s[0] == 45 { // "-"
		i += 1
	}

	for ; i < len(s); i++ {
		if s[i] == 46 { // "."
			foundDot = true
			i++
			break
		}
		if s[i] > 47 && s[i] < 58 {
			supSz += 1
			continue
		}
		return false
	}

	if foundDot == false {
		return false
	}

	if supSz < 1 || supSz > 3 {
		return false
	}

	for ; i < len(s); i++ {
		if s[i] > 47 && s[i] < 58 {
			subSz += 1
			continue
		}
		return false
	}

	if subSz > 15 {
		return false
	}

	return true
}

// Allow 1-16 digits, any order but nothing else
func validNumber(s string) bool {

	if len(s) < 1 {
		return false
	}

	if len(s) > 16 {
		return false
	}

	for _, c := range s {
		if c > 47 && c < 58 {
			continue
		}
		return false
	}

	return true
}
