package ads

import "github.com/vininta-solution/bid/model/placement"

type Ads struct {
	Id               int     `json:"id"`
	Bid              float64 `json:"bid"`
	Category         []int   `json:"category"`
	OnlyPlacement    []int   `json:"onlyPlacement"`
	ExcludePlacement []int   `json:"excludePlacement"`
}

func (ad *Ads) IsMatch(p placement.Placement) bool {
	var logic bool

	if len(ad.Category) > 0 {
		logic = false
		// At least one category matched
		for _, placementCategory := range p.Category {
			for _, adCategory := range ad.Category {
				if adCategory == placementCategory {
					logic = true
				}
			}
		}
		if logic == false {
			return false
		}
	}

	if len(ad.OnlyPlacement) > 0 {

	}

	return true
}
