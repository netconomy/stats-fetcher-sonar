package main

import (
	"fmt"
	"github.com/docopt/docopt-go"
	"net/http"
	"os"
	"strconv"
)

func main() {
	var port int
	usage := `SCC Stats Fetcher.

	Usage:
		scc-statsfetcher serve <port>

	Options:
		-h --help     Show this screen.`

	arguments, _ := docopt.Parse(usage, nil, true, "SCC Statsfetcher 1.0", false)
	fmt.Println(arguments)
	if arguments["serve"] == true {
		port, _ = strconv.Atoi(arguments["<port>"].(string))
		fmt.Println(port)
		if err := http.ListenAndServe(fmt.Sprintf("127.0.0.1:%v", port), nil); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}
