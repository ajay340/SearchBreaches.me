package populate

import (
	"fmt"
	"regexp"

	"log"

	"github.com/PuerkitoBio/goquery"
	"github.com/ajay340/SearchBreaches.me/database"
	"github.com/remeh/sizedwaitgroup"
)

type WikiBreachstruct struct {
	Entity  string
	Year    string
	Records string
	OrgType string
	Method  string
}

func getWikiBreaches() []WikiBreachstruct {
	re, err := regexp.Compile(`,`)
	if err != nil {
		log.Fatal(err)
	}
	breaches := []WikiBreachstruct{}
	doc, _ := goquery.NewDocument("https://en.wikipedia.org/wiki/List_of_data_breaches")
	doc.Find("table.wikitable tbody tr").Each(func(i int, s *goquery.Selection) {
		breaches = append(breaches, WikiBreachstruct{
			Entity:  s.Find("td").Eq(0).Text(),
			Year:    s.Find("td").Eq(1).Text(),
			Records: re.ReplaceAllString(s.Find("td").Eq(2).Text(), ""),
			OrgType: s.Find("td").Eq(3).Text(),
			Method:  s.Find("td").Eq(4).Text(),
		})
	})
	return breaches
}

func PopulateWKBreaches() {
	jobs := 6
	wg := sizedwaitgroup.New(jobs)
	for _, breach := range getWikiBreaches() {
		b := breach
		wg.Add()
		go func() {
			defer wg.Done()
			if database.FindRowBreach(b.Entity, "Name_of_Covered_Entity").ID == "" && b.Entity != "" {
				database.AddBreach(b.Entity, "UNKNOWN", "UNKNOWN", b.Records, "UNKNOWN", b.Method, "UNKNOWN", "UNKNOWN", "Breach of "+b.Entity+" by "+b.Method, "UNKNOWN", "UNKNOWN", b.Year, b.OrgType)
				fmt.Println("Breach added " + b.Entity + " from Wikipedia")
			}
		}()
	}
	wg.Wait()
}
