package core

import "math/rand"

type Swarm struct {
	Particles          []*Particle
	globalBestPosition []float64
}

//func newSwarm (upperBound []float64, lowerBound []float64, numOfParticles uint, xData []float64, yData []float64) *Swarm{
//	particles := make([]*Particle, numOfParticles)
//	var globalBestLoss float64
//	for ix := range particles{
//		aNewParticle := newParticle(upperBound,lowerBound)
//		MSEOfNewParticle := MSELoss(xData,yData,aNewParticle.localBestPosition)
//		if ix == 0 {
//			globalBestLoss = MSEOfNewParticle
//			//globalBestPosition :=
//		}else{
//			//if glo
//		}
//	}
//}
type Particle struct {
	position          []float64
	velocity          []float64
	localBestPosition []float64
}

func newParticle(upperBound []float64, lowerBound []float64) *Particle {

	dimensionOfParams := len(upperBound)
	position := make([]float64, dimensionOfParams)
	velocity := make([]float64, dimensionOfParams)
	localBestPosition := make([]float64, dimensionOfParams)
	for ix := range upperBound {
		positionRange := upperBound[ix] - lowerBound[ix]
		pos := rand.Float64()*positionRange + lowerBound[ix]
		position[ix] = pos
		velocity[ix] = -pos * 0.1
		localBestPosition[ix] = pos
	}
	return &Particle{
		position:          position,
		velocity:          velocity,
		localBestPosition: localBestPosition,
	}
}
