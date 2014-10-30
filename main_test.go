package main

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestBuildUDPPayload(t *testing.T) {
	assert.Equal(t, "sonar.metrics.xlmsp.coverage:5", buildUDPPayload("5", "coverage", "xlmsp"))
}

func TestBuildQueryUrl(t *testing.T) {
	assert.Equal(t, "sonar.local.netconomy.net/api/resources?format=json&includetrends=true&resource=xlmsp&metrics=coverage", buildQueryUrl("sonar.local.netconomy.net", "xlmsp", "coverage"))
}

func TestGetStatistics(t *testing.T) {
	resultJson := []map[string]interface{}{}

	jsonResp := []byte(`
		[
			{
				"msr": [
					{
						"val": 5.0	
					}	
				]	
			}	
		]
	`)
	json.Unmarshal(jsonResp, &resultJson)
	assert.Equal(t, 5, getStatistics(resultJson))
}

func TestGetUDPAddressFromConfig(t *testing.T) {
	cfg := map[string]interface{}{}
	json.Unmarshal([]byte(`
		{
			"udpserver": {
				"ip": "127.0.0.1/24",
				"port": 1234
			}	
		}	
	`), &cfg)
	result, err := getUDPAddressFromConfig(cfg)
	expectedIP, _, _ := net.ParseCIDR("127.0.0.1/24")
	assert.Equal(t, net.UDPAddr{IP: expectedIP, Port: 1234}, result)
	assert.Nil(t, err)
}

func TestGetSonarUrlFromConfig(t *testing.T) {
	cfg := map[string]interface{}{}
	json.Unmarshal([]byte(`
		{
			"sonar": {
				"url": "sonar.local.netconomy.net"	
			}	
		}	
	`), &cfg)
	assert.Equal(t, "sonar.local.netconomy.net", getSonarUrlFromConfig(cfg))
}

func TestGetMetricsFromConfig(t *testing.T) {
	cfg := map[string]interface{}{}
	json.Unmarshal([]byte(`
		{
			"metrics": [
				"coverage",
				"some",
				"things",
				"should",
				"be here"
			]	
		}	
	`), &cfg)
	assert.Equal(t, "coverage", getMetricsFromConfig(cfg)[0])
	assert.Equal(t, 5, len(getMetricsFromConfig(cfg)))
}

func TestHandleMetric(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, `
		[
			{
				"msr": [
					{
						"val": 5.0	
					}	
				]	
			}	
		]`)
	}))
	defer ts.Close()

	result, err := handleMetric("coverage", "xlmsp", ts.URL)
	assert.Equal(t, "sonar.metrics.xlmsp.coverage:5", result)
	assert.Nil(t, err)
}
