package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ynqa/word-embedding/builder"
	"gonum.org/v1/gonum/mat"
)

// Test uid:  5abc5152c2d9048b32bfc917

// MatchValue represents each porspecive user in their distance from the caller
type MatchValue struct {
	Usr  Profile      `json:"profile"`
	Dist float64      `json:"distance"`
	Loc  *Geolocation `json:"location"`
}

// WordModel is the vector representation of the words in the corpus file
var WordModel map[string][]float64

// WordModelContextWindow -
var WordModelContextWindow = 20

// WordModelDimension -
var WordModelDimension = 8

// WordModelRandomParam - number of words considered for similarity
var WordModelRandomParam = 50

// GetMatches returns an ordered list of user uid's from closest to furthest to the caller
func GetMatches(UserID string, neighbors []User) (MatchResponse, error) {
	callingUser, err := GetUser(UserID)
	if err != nil {
		return MatchResponse{}, errors.New("Unable to fetch calling user")
	}
	order := GetOrder(callingUser, neighbors, WordModel)
	return order, nil
}

// GetOrder - preprocess for sortMap
func GetOrder(caller User, prospUsers []User, model map[string][]float64) MatchResponse {
	callerStr := userToString(caller)
	callerVec := parToVector(callerStr, model)
	prospUsers = removeCaller(caller, prospUsers)
	var mr MatchResponse
	mr.Matches = []MatchValue{}
	start := time.Now()
	for _, pu := range prospUsers {
		prospStr := userToString(pu)
		prospVec := parToVector(prospStr, model)
		distance := nestedDistance(callerVec, prospVec)
		mr.Matches = append(mr.Matches, MatchValue{pu.Profile, distance, pu.Location})
	}
	elapsed := time.Since(start)
	fmt.Println("Destance Calculation took: " + elapsed.String())
	return mr
}

// // sortMap - returns uid's with shortest distance first
// func sortMap(m map[string]float64) []string {
// 	reverseMap := map[float64]string{}
// 	distances := []float64{}
// 	for uid, d := range m {
// 		reverseMap[d] = uid
// 		distances = append(distances, d)
// 	}
// 	sort.Float64s(distances)
// 	res := []string{}
// 	for _, d := range distances {
// 		res = append(res, reverseMap[d])
// 	}
// 	fmt.Println(res)
// 	return res
// }

// nestedDistance - distance metric between par vectors
func nestedDistance(src []*mat.VecDense, dst []*mat.VecDense) float64 {
	d := 0.0
	for _, si := range src {
		for _, di := range dst {
			temp := mat.NewVecDense(WordModelDimension, nil)
			temp.SubVec(si, di)
			d += flattenVector(WordModelDimension, temp)
		}
	}
	return d
}
func flattenVector(rows int, vec mat.Matrix) float64 {
	res := 0.0
	for i := 0; i < rows; i++ {
		res += math.Abs(vec.At(i, 0))
	}
	return res
}

// removeCaller takes the calling user out of prospective match list
func removeCaller(caller User, prospUsers []User) []User {
	s := -1
	for i, u := range prospUsers {
		if u.ID == caller.ID {
			s = i
		}
	}
	if s < 0 {
		return prospUsers
	}
	return append(prospUsers[:s], prospUsers[s+1:]...)

}

// stripStopWords -
func stripStopWords(str string) string {
	return ""
}

// parToVector converts a string representation of user to numeric vector using
// the given word embeddings model
func parToVector(userStr string, model map[string][]float64) []*mat.VecDense {
	res := []*mat.VecDense{}
	par := strings.Split(userStr, " ")
	n := WordModelRandomParam
	i := 0
	for i < n {
		l := len(par)
		randomWord := par[random(0, l)]
		randomWord = strings.TrimSpace(strings.ToLower(randomWord))
		if val, found := model[randomWord]; found {
			if len(val) == WordModelDimension {
				vec := mat.NewVecDense(WordModelDimension, val)
				res = append(res, vec)
				i++
			} // else {
			// 	fmt.Printf("[Meetover model warning!!] length of "+
			// 		"vector in model: %d for word : %s\n", len(val), randomWord)
			// }
		}
	}
	return res
}

// userToString converts the user object to a paragraph for vector translation
func userToString(u User) string {
	var res string
	res = u.Profile.Greeting + " " + u.Profile.Headline + " " + u.Profile.Summary + " " + u.Profile.Industry
	for _, pos := range u.Profile.Positions.Values {
		res += pos.Company.Industry + " "
		res += pos.Company.Name + " "
		res += pos.Summary + " "
		res += pos.Title
	}
	return res
}

// InitMLModel check if model has been created or creates it
func InitMLModel(windowSize int, wordDimensions int) {
	modelFile := "./ml/meetOver.model"
	if _, err := os.Stat(modelFile); os.IsNotExist(err) {
		fmt.Println("Model does not exist. Creating Model")
		corpusFile := "./ml/corpus.dat"
		createModel(modelFile, corpusFile, windowSize, wordDimensions)
	}
	WordModel = readModel(modelFile)
}

// createModel uses the word2vec algo to create word embeddings
func createModel(destinationFileName string, corpusFile string, windowSize int, wordDimensions int) {
	if _, err := os.Stat(corpusFile); os.IsNotExist(err) {
		fmt.Println("[-] Corpus file not found. No model created")
		return
	}
	b := builder.NewWord2VecBuilder()
	b.SetDimension(wordDimensions).
		SetWindow(windowSize).
		SetModel("skip-gram").
		SetOptimizer("ns").
		SetNegativeSampleSize(15).
		SetVerbose()
	m, err := b.Build()
	if err != nil {
		fmt.Println("[-] Unable to build word2vec neural net")
	}
	inputFile1, _ := os.Open(corpusFile)
	f1, err := m.Preprocess(inputFile1)
	if err != nil {
		fmt.Println("Failed to Preprocess.")
	}
	// Start to Train.
	m.Train(f1)
	f1.Close()
	// Save word vectors to a text file.
	m.Save(destinationFileName)
}

// readModel converts the generated model to an in-memory object
func readModel(modelFile string) map[string][]float64 {
	content, err := ioutil.ReadFile(modelFile)
	if err != nil {
		log.Fatal(err)
	}
	model := make(map[string][]float64)
	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		vector := strings.Split(line, " ")
		word := vector[0]
		n := len(vector)
		if n < 2 {
			vector = vector[1:]
		} else {
			vector = vector[1 : len(vector)-1]
		}
		floatVector := make([]float64, len(vector))
		for jj := range vector {
			floatVector[jj], err = strconv.ParseFloat(vector[jj], 64)
		}
		model[word] = floatVector
	}
	return model
}
