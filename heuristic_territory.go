
package main

import (
	"github.com/Battle-Bunker/cyphid-snake/agent"
)

// HeuristicTerritory rewards snakes for controlling more territory
func HeuristicTerritory(snapshot agent.GameSnapshot) float64 {
	board := snapshot.Board()
	you := snapshot.You()
	yourID := you.ID()
	
	// Get your territory size
	yourTerritory, exists := board.Territories[yourID]
	if !exists {
		return 0.0 // No territory found
	}
	
	yourTerritorySize := float64(len(yourTerritory))
	
	// Calculate total accessible territory on the board
	totalAccessibleCells := 0
	for _, cells := range board.Territories {
		totalAccessibleCells += len(cells)
	}
	
	// If no cells are accessible, avoid division by zero
	if totalAccessibleCells == 0 {
		return 0.0
	}
	
	// Calculate territory ratio (how much of the accessible board you control)
	territoryRatio := yourTerritorySize / float64(totalAccessibleCells)
	
	// Simple scoring as requested: 100 * territoryRatio
	return 100.0 * territoryRatio
}
