package main

import (
	"flag"
	"fmt"
	"net/http"
	"strings"

	"github.com/spf13/viper"
)

func main() {
	netinterface := flag.String("interface", "", "http network interface")
	port := flag.Int("port", 8000, "http serve port")
	flag.Parse()

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/etc/smart-device-relay/")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	lifx := LIFX_Commander{token: viper.GetString("LIFX_TOKEN")}

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

	httpServeAddr := fmt.Sprintf("%s:%d", *netinterface, *port)
	fmt.Printf("serving http on %s\n", httpServeAddr)
	err = http.ListenAndServe(httpServeAddr, nil)
	if err != nil {
		fmt.Printf("unable to serve http. err: %s\n", err.Error())
	}
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
