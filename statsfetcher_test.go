package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBuildUDPPayload(t *testing.T) {
	assert.Equal(t, "sonar.metrics.xlmsp.coverage:5", buildUDPPayload("5", "coverage", "xlmsp"))
}

func TestBuildQueryUrl(t *testing.T) {
	assert.Equal(t, "sonar.local.netconomy.net/api/resources?format=json&includetrends=true&resource=xlmsp&metrics=coverage", buildQueryUrl("sonar.local.netconomy.net", "xlmsp", "coverage"))
}

func TestGetStatistics(t *testing.T) {
	resultJson := [1]map[string]interface{}{}
	firstItem := map[string]interface{}{}
	msr := [1]interface{}{}
	firstValue := map[string]interface{}{}
	firstValue["val"] = 5.0
	msr[0] = firstValue
	firstItem["msr"] = msr[0:1]
	resultJson[0] = firstItem

	assert.Equal(t, 5, getStatistics(resultJson[0:1]))
}
