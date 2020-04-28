package recommender

import (
	"sort"
	"strconv"

	"github.com/ajay340/SearchBreaches.me/database"
)

type keyPair struct {
	Key   int
	Value int
}

func RecommendTopBreach(id string) [4]int {
	var recommendationsMap = map[int]int{}
	recommend := func(id string, breachFieldInput string, searchField string, additionalScore int) {
		for _, recommendation := range database.FindRowsBreach(breachFieldInput, searchField) {
			if id != recommendation.ID {
				recommendation_id, _ := strconv.Atoi(recommendation.ID)
				score := recommendationsMap[recommendation_id]
				recommendationsMap[recommendation_id] = score + additionalScore
			}
		}
	}
	givenBreach := database.FindRowBreach(id, "ID")
	recommend(id, givenBreach.Name_of_Covered_Entity, "Name_of_Covered_Entity", 15)
	recommend(id, givenBreach.Industry[0:3], "Industry", 35)
	recommend(id, givenBreach.Type_of_Breach[0:3], "Type_of_Breach", 20)
	recommend(id, givenBreach.State, "State", 5)
	recommend(id, givenBreach.Location_of_Breached_Information, "Location_of_Breached_Information", 5)

	var sortedRecommendations []keyPair
	for k, v := range recommendationsMap {
		sortedRecommendations = append(sortedRecommendations, keyPair{k, v})
	}

	sort.Slice(sortedRecommendations, func(i, j int) bool {
		return sortedRecommendations[i].Value > sortedRecommendations[j].Value
	})

	recommendations := new([4]int)
	for i, recommendation := range sortedRecommendations {
		if i == 4 {
			break
		}
		recommendations[i] = recommendation.Key
	}
	return *recommendations
}
