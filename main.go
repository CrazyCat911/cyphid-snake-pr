package main

import (
	"github.com/Battle-Bunker/cyphid-snake/agent"
	"github.com/Battle-Bunker/cyphid-snake/server"
	"github.com/BattlesnakeOfficial/rules/client"
)

func main() {

	metadata := client.SnakeMetadataResponse{
		APIVersion: "1",
		Author:     "zuthan",
		Color:      "#FF7F7F",
		Head:       "evil",
		Tail:       "nr-booster",
	}

	portfolio := agent.NewPortfolio(
		agent.NewHeuristic(1.0, "health", HeuristicHealth),
		agent.NewHeuristic(1.0, "food", HeuristicFood),
		agent.NewHeuristic(1.0, "space", HeuristicSpace),
		agent.NewHeuristic(1.0, "territory", HeuristicTerritory),
	)

	snakeAgent := agent.NewSnakeAgent(portfolio, metadata, 
		agent.WithTemperature(5.0),
		agent.WithPerformanceLogging(true))
	server := server.NewServer(snakeAgent)

	server.Start()
}
