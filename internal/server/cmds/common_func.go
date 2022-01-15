package cmds

import (
	"car-rental/internal/server/domain"
	"time"

	log "github.com/sirupsen/logrus"
)

/*
Build time frame for rents search
*/
func buildFromToFilter(from string, to string, isExists bool) domain.SearchParams {
	var searchParams domain.SearchParams
	fromExists := false
	if len(from) > 0 {
		fromExists = true
	}
	if fromExists && len(to) > 0 {
		brokenTime := false
		covertFrom, err := time.Parse(domain.TimeLayout, from)
		if err != nil {
			brokenTime = true
			log.Error(err)
		}
		covertTo, err := time.Parse(domain.TimeLayout, to)
		if err != nil {
			brokenTime = true
			log.Error(err)
		}
		if !covertFrom.Before(covertTo) {
			log.Error("Please provide correct dates, from must be less than too")
			brokenTime = true
		}
		if !brokenTime && !isExists {
			searchParams.DateFilter += " (to_time IS NULL or to_time<'" + from + "') or"
			searchParams.DateFilter += " (from_time IS NULL or from_time>'" + to + "')"
		} else if !brokenTime && isExists {
			searchParams.DateFilter += " (to_time IS NOT NULL AND from_time IS NOT NULL) and"
			searchParams.DateFilter += " ((from_time between '" + from + "' and '" + to + "') or (to_time between '" + from + "' and '" + to + "')) or "
			searchParams.DateFilter += " (from_time <='" + from + "' and to_time >= '" + to + "')"
		}

	}
	return searchParams
}
