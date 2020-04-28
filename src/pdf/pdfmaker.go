package pdf

import (
	//"fmt"
	"bytes"

	"github.com/ajay340/SearchBreaches.me/database"

	"github.com/johnfercher/maroto/pkg/consts"
	"github.com/johnfercher/maroto/pkg/pdf"
	"github.com/johnfercher/maroto/pkg/props"
)

func PDFgen(breach *database.Breach) (bytes.Buffer, error) {
	m := pdf.NewMaroto(consts.Portrait, consts.Letter)

	m.Row(8, func() {
		m.Col(func() {
			m.Text("Name: "+breach.Name_of_Covered_Entity, props.Text{
				Size: 15,
			})
		})
	})

	m.Row(6, func() {
		m.Col(func() {
			m.Text("State: "+breach.State, props.Text{
				Size: 12,
			})
		})
	})

	m.Line(6)

	m.Row(40, func() {
		m.Col(func() {
			m.Text("Summary: " + breach.Summary)
		})
	})

	var top float64 = 5

	m.Row(top, func() {
		m.Col(func() {
			m.Text("Business Associate Involved: " + breach.Business_Associate_Involved)
		})
	})
	m.Row(top, func() {
		m.Col(func() {
			m.Text("Individuals Affected: " + breach.Individuals_Affected)
		})
	})
	m.Row(top, func() {
		m.Col(func() {
			m.Text("Date of Breach: " + breach.Date_of_Breach)
		})
	})
	m.Row(top, func() {
		m.Col(func() {
			m.Text("Type of Breach: " + breach.Type_of_Breach)
		})
	})
	m.Row(top, func() {
		m.Col(func() {
			m.Text("Location of Breached Information: " + breach.Location_of_Breached_Information)
		})
	})
	m.Row(top, func() {
		m.Col(func() {
			m.Text("Date Posted or Updated: " + breach.Date_Posted_or_Updated)
		})
	})
	m.Row(top, func() {
		m.Col(func() {
			m.Text("Breach Start: " + breach.Breach_start)
		})
	})
	m.Row(top, func() {
		m.Col(func() {
			m.Text("Breach End: " + breach.Breach_end)
		})
	})
	m.Row(top, func() {
		m.Col(func() {
			m.Text("Year: " + breach.Year)
		})
	})
	m.Row(top, func() {
		m.Col(func() {
			m.Text("Industry: " + breach.Industry)
		})
	})

	return m.Output()
}
