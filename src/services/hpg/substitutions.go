package hpg

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

var weekdays = [5]string{"Mo", "Di", "Mi", "Do", "Fr"}

// Returns true if the given string begins with a weekday
func beginsWithAWeekday(s string) bool {
	b := strings.Fields(s)[0]
	for _, weekday := range weekdays {
		if b == weekday {
			return true
		}
	}
	return false
}

func GetSubstituationOfStudent(authid, authpw string) (map[string][]string, error) {

	// Request the HTML page.
	res, err := http.PostForm(fmt.Sprintf("https://vertretungsplan.hpg-speyer.de/pmwiki/pmwiki.php?n=Main.%s", strings.ToLower(authid)),
		url.Values{
			"authid": {authid},
			"authpw": {authpw},
		},
	)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	defer res.Body.Close()

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	if !strings.Contains(doc.Text(), "abmelden") {
		return nil, WrongCredentialsError
	}

	// Find the review items
	s := doc.Find("table") // if s.Length()=4, there are new substituations

	// Check if there are substitutions
	substitutionTableLength := doc.Find("#wikitext").Find("div").First().Find("table").Length()

	if substitutionTableLength < 1 {
		return map[string][]string{}, nil
	}

	sp := s.Eq(1)

	spMap := map[string][]string{}

	weekday := ""
	sp.Find("tr").Each(func(i int, s *goquery.Selection) {
		textToAdd := ""

		txt := strings.ReplaceAll(s.Text(), "\n", "")
		if beginsWithAWeekday(txt) {
			weekday = txt
			return
		}

		s.Find("td").Each(func(j int, t *goquery.Selection) {
			textToAdd += strings.TrimSpace(t.Text()) + " "
		})

		spMap[weekday] = append(spMap[weekday], strings.ReplaceAll(textToAdd, "\n", ""))
	})

	return spMap, nil
}
