package main

import "math"

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

const DegToRad = math.Pi / 180
const PreyTurnAmountRadians = 8 * DegToRad

var AgentViewAngle = math.Cos(math.Pi)
