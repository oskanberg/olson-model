package main

import "math"

const SimulationSpaceSize = 1000

const TotalSimulationSteps = 2000
const TotalGenerations = 20

const NumberOfPrey = 10
const NumberOfPredators = 10

const MaximumOutNodes = 4
const MaximumInNodes = 4

const MinimumOutNodes = 2
const MinimumInNodes = 1

const NumRetinaSlices = 12
const NumActuators = 2
const NumHiddenNodes = 6

//                    two layers øf retina in prey
const NumTotalNodes = NumRetinaSlices*2 + NumActuators + NumHiddenNodes

const RetinaSliceWidthRadians = math.Pi / NumRetinaSlices
const PreyViewDistance = 100
const PredatorViewDistance = 200

const DegToRad = math.Pi / 180
const PreyTurnAmountRadians = 8 * DegToRad
const PredatorTurnAmountRadians = 6 * DegToRad

const MutationRate = 0.01
const DuplicationLikelihood = 0.05
const DeletionLikelihood = 0.02

const InitialGenomeLength = 500
const ArtificialStartCodons = 4

const EatingDistance = 10

var AgentViewAngle = math.Cos(math.Pi / 4)
