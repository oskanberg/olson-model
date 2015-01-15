package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type AgentRecord struct {
	Position  Position
	NextPoint Vector2D
	AgentType string
	Fitness   int
	Sensors   string
}

type Step struct {
	Positions []AgentRecord
}

type Record struct {
	Steps       []Step
	currentStep int
}

func (s *Record) NewStep() {
	s.Steps = append(s.Steps, Step{})
	s.currentStep += 1
}

func (s *Record) AddRecordToCurrentStep(agent Agent) {
	agentRecord := AgentRecord{
		Position: Position{
			Location:  *agent.GetLocation(),
			Direction: *agent.GetDirection(),
		},
		Fitness: agent.GetFitness(),
		Sensors: agent.GetSensors(),
	}

	switch agent.(type) {
	case *Prey:
		agentRecord.AgentType = "Prey"
		agentRecord.NextPoint = *(agent.GetLocation().Add(agent.GetDirection().Multiplied(PreyViewDistance)))
	case *Predator:
		agentRecord.AgentType = "Predator"
		agentRecord.NextPoint = *(agent.GetLocation().Add(agent.GetDirection().Multiplied(PredatorViewDistance)))
	}
	s.Steps[s.currentStep].Positions = append(s.Steps[s.currentStep].Positions, agentRecord)
}

func (s *Record) WriteToFile(filename string) {
	jsonEnc, err := json.Marshal(s)
	if err != nil {
		fmt.Println(err)
	}
	f, err := os.Create("output/" + filename + ".json")
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()
	_, err = f.Write(jsonEnc)
	if err != nil {
		fmt.Println(err)
	}
}

func NewRecord() *Record {
	return &Record{
		Steps:       make([]Step, 1),
		currentStep: 0,
	}
}
