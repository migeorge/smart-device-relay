package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	port := getDefaultEnvVar("HTTP_PORT", "8080")
	lifx_token := getRequiredEnvVar("LIFX_TOKEN")

	lifx := LIFX_Commander{token: lifx_token}

	http.HandleFunc("/toggle/lifx/", func(w http.ResponseWriter, r *http.Request) {
		duration := "0"

		dataStr := strings.TrimPrefix(r.RequestURI, "/toggle/lifx/")
		segments := strings.Split(dataStr, "/")

		if len(segments) > 1 {
			duration = segments[1]
		}
		id := segments[0]

		lifx.ToggleDeviceByID(id, duration)

		w.WriteHeader(http.StatusOK)
	})

	httpServeAddr := fmt.Sprintf(":%s", port)
	log.Printf("serving http on %s\n", httpServeAddr)
	err := http.ListenAndServe(httpServeAddr, nil)
	if err != nil {
		log.Printf("unable to serve http. err: %s\n", err.Error())
	}
}

func getDefaultEnvVar(name, defaultValue string) string {
	value := os.Getenv(name)
	if value == "" {
		value = defaultValue
	}

	return value
}

func getRequiredEnvVar(name string) string {
	value := os.Getenv(name)
	if value == "" {
		log.Fatalf("FATAL: Unable to retrieve required environment variable: %s", name)
	}

	return value
}

type LIFX_Commander struct {
	token string
}

// ToggleDeviceByID fires a request to lifx, the response is not validated beyond a successful status
func (lc *LIFX_Commander) ToggleDeviceByID(id, duration string) {
	fmt.Printf("toggling LIFX device by id (%s) duration (%s)\n", id, duration)

	url := fmt.Sprintf("https://api.lifx.com/v1/lights/id%%3A%s/toggle", id)

	payload := strings.NewReader(fmt.Sprintf("{\"duration\":%s}", duration))

	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", lc.token))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("request failed. err: %s\n", err.Error())
	}

	defer res.Body.Close()
}
