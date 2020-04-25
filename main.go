package main

import (
	"log"
	"net/http"
	"os"
	"path"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/quhar/bme280"
	"golang.org/x/exp/io/i2c"
)

var (
	tempGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "bme280_temperature",
			Help: "Temperature in celsius degree",
		},
		[]string{"device", "location", "register"},
	)

	pressGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "bme280_pressure",
			Help: "Barometric pressure in hPa",
		},
		[]string{"device", "location", "register"},
	)

	humidityGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "bme280_humidity",
			Help: "Humidity in percentage of relative humidity",
		},
		[]string{"device", "location", "register"},
	)

	bus = &bme280.BME280{}

	device             string
	location           string
	registerAddressStr string
)

func main() {
	r := mux.NewRouter()
	r.Use(bme280collector)
	r.Handle("/metrics", promhttp.Handler())

	if err := prometheus.Register(tempGauge); err != nil {
		panic(err)
	}
	if err := prometheus.Register(pressGauge); err != nil {
		panic(err)
	}
	if err := prometheus.Register(humidityGauge); err != nil {
		panic(err)
	}

	listenAddress := os.Getenv("LISTEN_ADDRESS")
	if listenAddress == "" {
		listenAddress = ":8080"
	}

	device = os.Getenv("BME280_DEVICE")
	if device == "" {
		device = "i2c-1"
	}

	location = os.Getenv("BME280_LOCATION")
	if location == "" {
		location = ""
	}

	registerAddressStr = os.Getenv("BME280_REGISTER_ADDRESS")
	if registerAddressStr == "" {
		registerAddressStr = "0x77"
	}

	registerAddress, err := strconv.ParseInt(registerAddressStr, 0, 32)
	if err != nil {
		log.Println(err.Error())
	}

	d, err := i2c.Open(&i2c.Devfs{Dev: path.Join("/dev/", device)}, int(registerAddress))
	if err != nil {
		log.Println(err.Error())
	}

	bus = bme280.New(d)
	err = bus.Init()

	if err != nil {
		log.Println(err.Error())
		panic(err)
	}

	defer d.Close()

	log.Printf("Listenning on '%s'. Fetching data from device '%s' on register address '%s'", listenAddress, device, registerAddressStr)

	log.Fatal(http.ListenAndServe(listenAddress, r))
}

func bme280collector(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		t, p, h, err := bus.EnvData()
		if err != nil {
			log.Println(err.Error())
			return
		}

		tempGauge.WithLabelValues(device, location, registerAddressStr).Set(t)
		pressGauge.WithLabelValues(device, location, registerAddressStr).Set(p)
		humidityGauge.WithLabelValues(device, location, registerAddressStr).Set(h)

		next.ServeHTTP(w, r)
	})
}
