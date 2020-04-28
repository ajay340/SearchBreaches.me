package populate

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ajay340/SearchBreaches.me/database"
	"github.com/remeh/sizedwaitgroup"
)

type HIBPBreach struct {
	Name         string   `json:Name`
	Title        string   `json:Title`
	Domain       string   `json:Domain`
	BreachDate   string   `json:BreachDate`
	AddedDate    string   `json:AddedDate`
	ModifiedDate string   `json:ModifiedDate`
	PwnCount     int      `json:PwnCount`
	Description  string   `json:Description`
	LogoPath     string   `json:LogoPath`
	DataClasses  []string `json:DataClasses`
	IsVerified   bool     `json:IsVerified`
	IsFabricated bool     `json:IsFabricated`
	IsSensitive  bool     `json:IsSensitive`
	IsRetired    bool     `json:IsRetired`
	IsSpamList   bool     `json:IsSpamList`
}

func getHIBPBreaches() []HIBPBreach {
	url := "https://haveibeenpwned.com/api/v3/breaches"
	apiRequest := http.Client{
		Timeout: time.Second * 30,
	}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("User-Agent", "searchbreaches.me")
	res, getErr := apiRequest.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}
	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}
	var breaches []HIBPBreach
	jsonErr := json.Unmarshal(body, &breaches)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}
	return breaches
}

func DetermineIndustry(description string) string {
	healthcare := []string{"healthcare", "health", "fitness", "wellness", "hospital", "doctor", "nurse", "EMT", "Surgeon", "medical practice", "clinical", "medicine", "treatment", "patients", "therapy", "chiropractic", "pharmacy", "patient", "pediatric", "prescription drugs", "diagnosis", "dentistry"}
	gaming := []string{"gaming", "gamer", "video game", "video games", "game developer", "game development", "game dev", "game", "players", "play"}
	finance := []string{"finance", "economics", "investment", "bank", "banking", "fund", "credit", "business", "corporate", "equity", "tax", "accounting", "investing", "fiscal", "money", "loan", "stonk", "stock", "budget", "debt", "asset", "pay", "equity"}
	insurance := []string{"insurance", "life insurance", "reinsurance", "claims-adjuster", "policy", "coverage", "insurer", "retirement", "mortgage", "health insurance", "property insurance", "loss ratio", "insurance fraud", "disability insurance", "tax", "coverage", "medicare", "premium", "pension"}
	government := []string{"government", "politics", "administration", "governing", "judiciary", "political science", "congress", "constitutional", "political party", "politics", "election", "democracy", "leaders", "ministry", "law", "parliment"}
	web := []string{"web", "site", "social", "social media"}
	technology := []string{"technology", "hosting", "electronics", "information technology", "information system", "engineering", "informatics", "hardware", "cyberspace", "computer programming", "robots", "intel", "HP", "AWS", "lenovo", "NVIDIA", "radeon"}
	industries := [][]string{gaming, finance, healthcare, insurance, government, technology, web}

	for _, industry := range industries {
		for _, term := range industry {
			if strings.Contains(description, term) {
				return industry[0]
			}
		}
	}
	return "Other"
}

func PopulateHIBPBreaches() {
	jobs := 6
	wg := sizedwaitgroup.New(jobs)
	for _, breach := range getHIBPBreaches() {
		b := breach
		wg.Add()
		go func() {
			defer wg.Done()
			if database.FindRowBreach(b.Name, "Name_of_Covered_Entity").ID == "" && b.Name != "" {
				industry := DetermineIndustry(b.Description)
				database.AddBreach(b.Name, "UNKNOWN", "UNKNOWN", strconv.Itoa(b.PwnCount), b.BreachDate, "Hacking", "UNKNOWN", b.AddedDate, b.Description, "UNKNOWN", "UNKNOWN", b.BreachDate[0:4], industry)
				fmt.Println("Breach added " + b.Name + " from HaveIBeenPwned.com")
			}
		}()
	}
	wg.Wait()
}
