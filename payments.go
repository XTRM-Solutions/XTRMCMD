package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

func shortPause() {
	v := time.Duration(5*1000*1000) + time.Duration(rand.Int63n(2*1000*1000*100)) +
		time.Duration(rand.Int63n(2*1000*1000*100)) +
		time.Duration(rand.Int63n(2*1000*1000*100))
	time.Sleep(v)
}

func mockPayments() {

	rand.Seed(time.Now().UTC().UnixNano())
	fmt.Println("")
	for firstName, lastName := range patronNames {
		amount, fee, total := newDetails()
		shortPause()
		fmt.Printf("Success! TransactionID %s for %s (%s transferred, %s fee) %s to recipient %s (%s %s %s)\n",
			newTransactionString(),
			total, amount, fee, "USD",
			getPat(),
			firstName, lastName,
			strings.ToLower(firstName+"."+lastName+"@mythosseed.com"))

	}
}

var patronList = map[int64]bool{0: true}

func getBool(k int64) (b bool) {
	_, b = patronList[k]
	return b
}

func getPat() string {
	var ix int64
	for ix = rand.Int63n(19999999) + rand.Int63n(99999) + 678; getBool(ix); {
		// collision? Try again!
		ix = rand.Int63n(5555) + rand.Int63n(55) + 678
	}
	patronList[ix] = true
	return "PAT" + strconv.FormatInt(ix, 10)

}

var tx int64 = 11921

func newTransactionString() (rx string) {
	tx += rand.Int63n(5) + rand.Int63n(5) + rand.Int63n(5) + 1
	return strconv.FormatInt(tx, 10)
}

func newDetails() (amount string, fee string, total string) {

	num := rand.Int31n(51) + 11
	dec := rand.Int31n(100)
	amount = fmt.Sprintf("%d.%02d", int(num), int(dec))
	v, _ := strconv.ParseFloat(amount, 64)
	w := v*.0125 + 0.97 // fee
	fee = fmt.Sprintf("%.02f", w)
	total = fmt.Sprintf("%.02f", v+w)

	return amount, fee, total
}

var patronNames = map[string]string{
	"Jessie":     "Orr",
	"Juliet":     "Vega",
	"Casey":      "Copeland",
	"Damien":     "Oneal",
	"Jamya":      "Stuart",
	"Davion":     "Levy",
	"Pierre":     "Clarke",
	"Kayleigh":   "Castaneda",
	"Azul":       "Bowen",
	"Gillian":    "Santana",
	"Garrett":    "Trujillo",
	"Marie":      "Mcfarland",
	"Cassius":    "Scott",
	"Paloma":     "Rivers",
	"Hudson":     "Ibarra",
	"Anna":       "Hartman",
	"Savion":     "Kelley",
	"Konnor":     "Thompson",
	"Aylin":      "Hodges",
	"Ashanti":    "Taylor",
	"Krista":     "Swanson",
	"Andrea":     "Shaw",
	"Zaire":      "Evans",
	"William":    "Bentley",
	"Sophie":     "Gregory",
	"Carissa":    "House",
	"Aleah":      "Carpenter",
	"Jeremiah":   "Nash",
	"Corinne":    "Madden",
	"Harry":      "Lane",
	"Meghan":     "Rodgers",
	"Sylvia":     "Combs",
	"Trystan":    "Meyers",
	"Justus":     "Daugherty",
	"Aileen":     "Moses",
	"Whitney":    "Small",
	"Lauren":     "Aguirre",
	"Ella":       "Perry",
	"Gaige":      "Mcconnell",
	"Giovani":    "Pitts",
	"Barthold":   "Mccullough",
	"Carleigh":   "Horne",
	"Johanna":    "Wilkinson",
	"Devyn":      "Maddox",
	"Matteo":     "Humphrey",
	"Joy":        "Holmes",
	"Jordin":     "Ballard",
	"Nash":       "Williams",
	"Camren":     "Jefferson",
	"Nora":       "Leon",
	"Willie":     "Morton",
	"Valerie":    "Bishop",
	"Emanuel":    "Michael",
	"Griffin":    "Baker",
	"Eddie":      "Werner",
	"Jaylen":     "Blackwell",
	"Andy":       "Carter",
	"Finnegan":   "Shields",
	"Reece":      "Conner",
	"Sheldon":    "Reid",
	"Aditya":     "Frye",
	"Kaylah":     "Holloway",
	"Marco":      "Shields",
	"August":     "Phelps",
	"Marvin":     "Grimes",
	"Addyson":    "Wilcox",
	"Haylee":     "Montes",
	"Noah":       "Moran",
	"Anabelle":   "Hill",
	"Scarlet":    "Flowers",
	"Elliot":     "Brewer",
	"Aryan":      "Black",
	"Amaya":      "Pratt",
	"Alisa":      "Davies",
	"Kellen":     "Richards",
	"Kristina":   "Wilkinson",
	"Amir":       "Collins",
	"Aurora":     "Tapia",
	"Jaime":      "Baxter",
	"Cherish":    "Sawyer",
	"Lukas":      "Dunn",
	"Ariana":     "Dickson",
	"Kadence":    "Paul",
	"Kierra":     "Bryant",
	"Rebekah":    "Carpenter",
	"Brisa":      "Mccann",
	"Addison":    "Krueger",
	"Kolby":      "Jackson",
	"Kaleb":      "Rivera",
	"Mckayla":    "Stephenson",
	"Dorian":     "Barnett",
	"Alexzander": "Cole",
	"Ryder":      "Leblanc",
	"Leticia":    "Cannon",
	"Jack":       "Hayden",
	"Ryleigh":    "Huffman",
	"Alexander":  "Ortiz",
	"Karter":     "Doyle",
	"Ayana":      "Caldwell",
	"Jovanni":    "Boyle",
}
