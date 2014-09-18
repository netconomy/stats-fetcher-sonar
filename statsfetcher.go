package main

import (
	"encoding/json"
	"fmt"
	"github.com/docopt/docopt-go"
	"github.com/gorilla/mux"
	"io"
	"io/ioutil"
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
	fmt.Printf("unmarshal: %v, %v\n", cfg, cfg["config"])
	return cfg, nil
}

func getMetricsFromConfig(cfg map[string]interface{}) []interface{} {
	return (cfg["metrics"]).([]interface{})
}

func getSonarUrlFromConfig(cfg map[string]interface{}) string {
	return ((cfg["sonar"]).(map[string]interface{})["url"]).(string)
}

func buildQueryUrl(baseUrl string, project string, metric string) string {
	return baseUrl + "/api/resources?format=json&includetrends=true&resource=" + project + "&metrics=" + metric
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
	sonarUrl := getSonarUrlFromConfig(cfg)
	metrics := getMetricsFromConfig(cfg)

	if arguments["serve"] == true {
		r := mux.NewRouter()
		r.HandleFunc("/fetch/{project}", func(resp http.ResponseWriter, req *http.Request) {
			vars := mux.Vars(req)
			project := vars["project"]

			res, _ := http.Get(buildQueryUrl(sonarUrl, project, metrics[0].(string)))
			body, _ := ioutil.ReadAll(res.Body)
			fmt.Println(string(body))

			io.WriteString(resp, "hi\n"+"project: "+project+" response: "+string(body))
		})
		http.Handle("/", r)
		port, _ = strconv.Atoi(arguments["<port>"].(string))
		fmt.Println(port)
		if err := http.ListenAndServe(fmt.Sprintf("127.0.0.1:%v", port), nil); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}
