package main

import (
	"encoding/json"
	"fmt"
	"github.com/docopt/docopt-go"
	"github.com/gorilla/mux"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strconv"
)

func readConfig() (map[string]interface{}, error) {
	file, e := ioutil.ReadFile("./config.json")
	if e != nil {
		fmt.Printf("File error: %v\n", e)
	}
	var cfg map[string]interface{}
	json.Unmarshal(file, &cfg)
	return cfg, nil
}

func getSonarUrlFromConfig(cfg map[string]interface{}) string {
	return ((cfg["sonar"]).(map[string]interface{})["url"]).(string)
}

func getUDPAddressFromConfig(cfg map[string]interface{}) net.UDPAddr {
	udpserver := cfg["udpserver"].(map[string]interface{})
	parsedIp, _, _ := net.ParseCIDR(udpserver["ip"].(string))
	return net.UDPAddr{IP: parsedIp, Port: int(udpserver["port"].(float64))}
}

func getMetricsFromConfig(cfg map[string]interface{}) []interface{} {
	return (cfg["metrics"]).([]interface{})
}

func handleMetric(metric string, project string, sonarUrl string, conn net.UDPConn, udpAddr net.UDPAddr) {
	res, _ := http.Get(buildQueryUrl(sonarUrl, project, metric))
	body, _ := ioutil.ReadAll(res.Body)

	resultJson := []map[string]interface{}{}
	json.Unmarshal(body, &resultJson)

	if len(resultJson) > 0 {
		statistics := getStatistics(resultJson)
		udpPayload := buildUDPPayload(strconv.FormatFloat(statistics, 'f', -1, 64), metric, project)

		conn.WriteToUDP([]byte(udpPayload), &udpAddr)
	}
}

func buildQueryUrl(baseUrl string, project string, metric string) string {
	return baseUrl + "/api/resources?format=json&includetrends=true&resource=" + project + "&metrics=" + metric
}

func getStatistics(json []map[string]interface{}) float64 {
	msr := json[0]["msr"]
	return msr.([]interface{})[0].(map[string]interface{})["val"].(float64)
}

func buildUDPPayload(value string, metric string, project string) string {
	return "sonar.metrics." + project + "." + metric + ":" + value
}

func main() {
	var port int
	readConf, _ := readConfig()
	cfg := readConf["config"].(map[string]interface{})
	usage := `SCC Stats Fetcher.

	Usage:
		scc-statsfetcher serve <port>

	Options:
		-h --help     Show this screen.`

	arguments, _ := docopt.Parse(usage, nil, true, "SCC Statsfetcher 1.0", false)

	if arguments["serve"] == true {
		sonarUrl := getSonarUrlFromConfig(cfg)
		udpAddr := getUDPAddressFromConfig(cfg)
		metrics := getMetricsFromConfig(cfg)
		conn, _ := net.ListenUDP("udp", &net.UDPAddr{IP: net.IPv4zero, Port: 0})

		r := mux.NewRouter()
		r.HandleFunc("/fetch/{project}", func(resp http.ResponseWriter, req *http.Request) {
			project := mux.Vars(req)["project"]

			for _, element := range metrics {
				handleMetric(element.(string), project, sonarUrl, *conn, udpAddr)
			}

			io.WriteString(resp, "metrics tracked for "+project)
		})
		http.Handle("/", r)

		port, _ = strconv.Atoi(arguments["<port>"].(string))
		if err := http.ListenAndServe(fmt.Sprintf("127.0.0.1:%v", port), nil); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}
