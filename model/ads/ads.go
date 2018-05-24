package ads

import "github.com/vininta-solution/bid/model/placement"

type Ads struct {
	Id       int     `json:"id"`
	Bid      float64 `json:"bid"`
	Category int     `json:"category"`
}

func (ad *Ads) IsMatch(p placement.Placement) bool {
	if ad.Category == p.Category {
		return true
	}
	return false
}
