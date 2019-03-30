package storage

import (
	"strconv"
	"testing"
	"time"
)

func TestFormJson(t *testing.T) {
	inputString := []string{"empty", "AI1", "UI2", "AI2", "UI3", "AI3"}
	resp := formJson(inputString)
	t.Log(resp)
}

func TestAddtoMongo(t *testing.T) {
	inputString := []string{"empty", "AI1", "UI2", "AI2", "UI3", "AI3"}
	err := AddToMongo(time.Now(), inputString)
	if err != nil {
		t.Fatal("got error inserting data in mongo", err)
	}
	t.Log("Success!")
}

func BenchmarkAddtoMongo10000(b *testing.B) {
	for i := 0; i < b.N; i++ {
		inputString := []string{"empty", "AI1" + strconv.Itoa(i), "UI2" + strconv.Itoa(i), "AI2", "UI3" + strconv.Itoa(i), "AI3" + strconv.Itoa(i)}
		err := AddToMongo(time.Now(), inputString)
		if err != nil {
			b.Fatal("got error inserting data in mongo", err)
		}
	}
}
