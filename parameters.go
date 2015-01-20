package main

import "math"

const SimulationSpaceSize = 1000

const TotalSimulationSteps = 2000
const TotalGenerations = 2000
const RoundsPerGeneration = 3

const NumberOfPrey = 50
const NumberOfPredators = 20

const NumberOfPredatorClones = 1
const NumberOfPreyClones = 25

const MaximumOutNodes = 4
const MaximumInNodes = 4

const MinimumOutNodes = 2
const MinimumInNodes = 1

// Number of slices per species (i.e. actual is *2)
const NumRetinaSlices = 12

const NumActuators = 2
const NumHiddenNodes = 6

//                    two layers øf retina in prey
const NumTotalNodes = NumRetinaSlices*2 + NumActuators + NumHiddenNodes

const NodeWidth = 255 / NumTotalNodes

const PreyViewDistance = 100
const PredatorViewDistance = 200

const DegToRad = math.Pi / 180

const PreyTurnAmountRadians = 8 * DegToRad
const PredatorTurnAmountRadians = 6 * DegToRad
const PredatorSpeedMultiplier = 3

const MutationRate = 0.01
const DuplicationLikelihood = 0.05
const DeletionLikelihood = 0.02

const InitialGenomeLength = 500
const ArtificialStartCodons = 4

const EatingDistance = 10

var AgentViewAngleRadians = math.Pi
var HalfAgentViewAngleRadians = AgentViewAngleRadians / 2

var CosHalfAgentViewAngle = math.Cos(AgentViewAngleRadians / 2)
var RetinaSliceWidth = AgentViewAngleRadians / NumRetinaSlices

const PredatorConfusion = true
const EatCooldown = true

const SavePredators = false
const SeedPredators = true

const SavePrey = true
const SeedPrey = true

const RigMarkovNetwork = false
const Model = "LinearWeights"

// const Model = "MarkovNetwork"
// const Model = "Olson"

const PreyHeadStart = 100
