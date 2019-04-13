package main

import (
	"testing"

	"github.com/TexaProject/texajson"
)

func TestSlabtoJSON(t *testing.T) {
	input := []texajson.SlabPage{{"news", 3, 0, 4}, {"sports", 2, 0, 4}}
	texajson.SlabToJson(input)
}

//[{Pandi [{news 999}]} {AI [{news 0} {sports 2}]}

func TestCattoJSON(t *testing.T) {
	input := []texajson.CatPage{{"Pandi", []texajson.CatValArray{{"news", 999}}}, {"AI", []texajson.CatValArray{{"news", 0}, {"sports", 2}}}}
	texajson.CatToJson(input)
}
