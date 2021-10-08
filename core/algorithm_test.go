package core

import "testing"

func TestNewParticle(t *testing.T) {
	upperBound := []float64{2.0, 3.0}
	lowerBound := []float64{-2.0, -3.0}
	t.Log(newParticle(upperBound, lowerBound))
}
