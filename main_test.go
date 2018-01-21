package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSendAlert(t *testing.T) {
	newAlert := sendAlert(true, false)
	assert.Equal(t, newAlert, true, "We should send an alert when it is new")

	currentlyAlerting := sendAlert(true, true)
	assert.Equal(t, currentlyAlerting, false, "We should not send an alert when we are already in alert")

	needsRecovery := sendAlert(false, true)
	assert.Equal(t, needsRecovery, false, "We should not send an alert when we are not in a bad place")

	allGood := sendAlert(false, false)
	assert.Equal(t, allGood, false, "No need to alert when we are all good")
}

func TestRecoverAlert(t *testing.T) {
	newAlert := recoverAlert(true, false)
	assert.Equal(t, newAlert, false, "We should not recover an alert when it is new")

	currentlyAlerting := recoverAlert(true, true)
	assert.Equal(t, currentlyAlerting, false, "We should not recovery an alert when we are still in a bad spot")

	needsRecovery := recoverAlert(false, true)
	assert.Equal(t, needsRecovery, true, "We should recover an alert when we are not in a bad place")

	allGood := recoverAlert(false, false)
	assert.Equal(t, allGood, false, "No need to recover when we are all good")

}

func TestOverThreshold(t *testing.T) {
	over := overThreshold(10.0, 1.0)
	assert.Equal(t, over, true, "We are over threshold when the value is larger than the threshold")

	under := overThreshold(1.0, 10.0)
	assert.Equal(t, under, false, "We are not over threshold when the value is smaller than the threshold")

	equal := overThreshold(10.0, 10.0)
	assert.Equal(t, equal, true, "We are over threshold when the value is equal to the threshold")
}
