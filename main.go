package main

import (
	"encoding/json"
	"math"
	"math/rand"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
)

//Income is to accept the post request from Python
type Income struct {
	XData        []float64
	YData        []float64
	PositionMax  []float64
	PositionMin  []float64
	NumStepWC1C2 []float64
}
type particle struct {
	position          []float64
	velocity          []float64
	localBestPosition []float64
	localBestLoss     float64
}

//triexponetialLoss returns the Mean Square Error between the triexponetial model and the experiment curve
func triexponetialLoss(params []float64, xData []float64, yData []float64) float64 {
	a1, a2, a3, e1, e2, e3 := params[0], params[1], params[2], params[3], params[4], params[5]
	yPredict := make([]float64, len(xData))
	Loss := 0.0
	for ix := range xData {
		yPredict[ix] = a1*math.Exp(-xData[ix]/e1) + a2*math.Exp(-xData[ix]/e2) + a3*math.Exp(-xData[ix]/e3)
		Loss += math.Pow(yData[ix]-yPredict[ix], 2)
	}
	return Loss
}
func newParticle(dimension int, positionRange []float64, xData []float64, yData []float64) *particle {

	position := make([]float64, dimension)
	velocity := make([]float64, dimension)
	localBestPosition := make([]float64, dimension)
	for ix := 0; ix < dimension; ix++ {
		//We want to randomly distribute each particle to the parameter space
		position[ix] = rand.Float64() * positionRange[ix]
		//Initial velocity is -0.1*Position
		velocity[ix] = position[ix] * -0.1
	}
	copy(localBestPosition, position)
	localBestLoss := triexponetialLoss(position, xData, yData)
	return &particle{position, velocity, localBestPosition, localBestLoss}
}
func (p *particle) evolution(dimension int, w float64, c1 float64, c2 float64, globalBestPosition []float64, xData []float64, yData []float64, positionMax []float64, positionMin []float64, velocityMax []float64, velocityMin []float64) {
	//v=w*v+c1*rand*(local_best_p-p)+c2*rand*(global_best_p-p)
	for ix := 0; ix < dimension; ix++ {
		p.velocity[ix] = w*p.velocity[ix] + c1*rand.Float64()*(p.localBestPosition[ix]-p.position[ix]) + c2*rand.Float64()*(globalBestPosition[ix]-p.position[ix])
		//To limit the maximum of velocity
		if p.velocity[ix] >= velocityMax[ix] {
			p.velocity[ix] = velocityMax[ix]
		}
		if p.velocity[ix] <= velocityMin[ix] {
			p.velocity[ix] = velocityMin[ix]
		}
	}
	//p = p+v
	for ix := 0; ix < dimension; ix++ {
		p.position[ix] = p.position[ix] + p.velocity[ix]
		//To avoid the particle escaping the bound
		if p.position[ix] >= positionMax[ix] {
			p.position[ix] = positionMax[ix]
		}
		if p.position[ix] <= positionMin[ix] {
			p.position[ix] = positionMin[ix]
		}
	}
	//Update the localBestPosition and localBestLoss
	newLoss := triexponetialLoss(p.position, xData, yData)
	if newLoss <= p.localBestLoss {
		p.localBestLoss = newLoss
		for ix := 0; ix < dimension; ix++ {
			p.localBestPosition[ix] = p.position[ix]
		}
	}
}

type swarm struct {
	particles          []*particle
	globalBestPosition []float64
	globalBestLoss     float64
	positionMax        []float64
	positionMin        []float64
	velocityMax        []float64
	velocityMin        []float64
	dimension          int
	numofParticles     int
}

func newSwarm(numofParticles int, positionMax []float64, positionMin []float64, xData []float64, yData []float64) *swarm {
	//Calculate the position and velocity limits
	dimension := len(positionMax)
	positionRange := make([]float64, dimension)
	velocityMax := make([]float64, dimension)
	velocityMin := make([]float64, dimension)
	for ix := 0; ix < dimension; ix++ {
		positionRange[ix] = positionMax[ix] - positionMin[ix]
		velocityMax[ix] = positionRange[ix] / 10
		velocityMin[ix] = velocityMax[ix] * -1
	}
	//Construct the particles
	particles := make([]*particle, numofParticles)
	globalBestPosition := make([]float64, dimension)
	var globalBestLoss float64
	for ix := 0; ix < numofParticles; ix++ {
		aNewParticle := newParticle(dimension, positionRange, xData, yData)
		if ix == 0 {
			globalBestPosition = aNewParticle.localBestPosition
			globalBestLoss = aNewParticle.localBestLoss
		} else {
			if aNewParticle.localBestLoss <= globalBestLoss {
				globalBestPosition = aNewParticle.localBestPosition
				globalBestLoss = aNewParticle.localBestLoss
			}
		}
		particles[ix] = aNewParticle
	}
	return &swarm{particles, globalBestPosition, globalBestLoss, positionMax, positionMin, velocityMax, velocityMin, dimension, numofParticles}
}
func (s *swarm) evolution(w float64, c1 float64, c2 float64, xData []float64, yData []float64) {
	chanIn := make(chan int)
	numofThreads := 40
	for th := 0; th < numofThreads; th++ {
		go func() {
			for ix := range chanIn {
				s.particles[ix].evolution(s.dimension, w, c1, c2, s.globalBestPosition, xData, yData, s.positionMax, s.positionMin, s.velocityMax, s.velocityMin)
			}
		}()
	}
	for ix := 0; ix < s.numofParticles; ix++ {
		chanIn <- ix
	}
	close(chanIn)
	// for ix := 0; ix < s.numofParticles; ix++ {
	// 	s.particles[ix].evolution(s.dimension, w, c1, c2, s.globalBestPosition, xData, yData, s.positionMax, s.positionMin, s.velocityMax, s.velocityMin)
	// }
	for ix := 0; ix < s.numofParticles; ix++ {
		if s.particles[ix].localBestLoss <= s.globalBestLoss {
			s.globalBestPosition = s.particles[ix].localBestPosition
			s.globalBestLoss = s.particles[ix].localBestLoss
		}
	}
}

func main() {
	runtime.GOMAXPROCS(8)
	rand.Seed(time.Now().UnixNano())
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	router.POST("/", func(c *gin.Context) {
		incomeJSON := c.PostForm("data")
		var income Income
		_ = json.Unmarshal([]byte(incomeJSON), &income)
		numofParticles := int(income.NumStepWC1C2[0])
		steps := int(income.NumStepWC1C2[1])
		w := income.NumStepWC1C2[2]
		c1 := income.NumStepWC1C2[3]
		c2 := income.NumStepWC1C2[4]
		xData := income.XData
		yData := income.YData
		positionMax := income.PositionMax
		positionMin := income.PositionMin
		aNewSwarm := newSwarm(numofParticles, positionMax, positionMin, xData, yData)
		for step := 0; step < steps; step++ {
			aNewSwarm.evolution(w, c1, c2, xData, yData)
		}

		c.JSON(200, gin.H{"MSE": aNewSwarm.globalBestLoss, "param_opt": aNewSwarm.globalBestPosition})
	})
	router.Run(":8080")
}
