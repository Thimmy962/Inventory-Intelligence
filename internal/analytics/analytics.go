package analytics

import (
	"context"
	"database/sql"
	"errors"
	// "log"
	"main/internal/database"
	"math"
)

/**
	Decay constant or lambda: λ=−ln(1−α) 0.2 for 10 days
	Decay Factor: w=e^(−λ⋅Δt)
	Time decayed EWMA: S(t) = ((1-w) * V(t)) + (w * S(t-1))

*/
var lambda = 0.2 
func Analytics(ctx context.Context, queries *database.Queries, sales_id int32) {
	// get sales time
	timeT, _ := queries.GetSalesTime(ctx, sales_id)
	// get all sales_items associated with this sale
	salesItemIDs, _ := queries.GetSalesItems(ctx, sales_id)



	for _, items := range salesItemIDs {
		ewma, nerr := queries.GetRecentProductEWMA(ctx, items.ProductID)
		// subtract lastupdated from now

		if errors.Is(nerr, sql.ErrNoRows) {
			// Timestamp source: sale.CreatedAt (not worker execution time).
			//
			// Rationale:
			//   EWMA weight = e^(-λ * Δt)
			//   Δt = time elapsed since previous sale.
			//
			// Using time.Now() would inflate Δt by the queue/processing delay,
			// causing this sale to be incorrectly discounted before the next one arrives.
			// The decay clock should start when the event *happened*, not when we *processed* it.
			queries.AddNewEWMA(ctx, database.AddNewEWMAParams{
				ProductID: items.ProductID, Ewma: float64(items.QuantitySold), RecordedAt: timeT, 
			})
		}  else {
			// get time difference
			timeDiff := timeT.Sub(ewma.RecordedAt)
			// convert time diff to days i.e. Δt
			deltaTime := timeDiff.Seconds() / 86400.0
			decayFactor := math.Exp(-lambda * deltaTime)
			newEWMA :=  float64(items.QuantitySold) + (decayFactor * ewma.Ewma)
			queries.AddNewEWMA(ctx, database.AddNewEWMAParams{
				ProductID: items.ProductID,Ewma: newEWMA, RecordedAt: timeT,
			})
		}
	}

}

// calculates the ewma for a product being sold the first time
func NewEWMA() {

}