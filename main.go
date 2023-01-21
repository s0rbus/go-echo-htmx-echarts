package main

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"math/rand"
	"net/http"
	"time"

	//"github.com/glebarez/sqlite"
	//"github.com/urfave/cli/v2"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/go-echarts/v2/types"

	"github.com/labstack/echo/v4"
)

// define a Teamplate type and implement echo.Renderer
type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

var CHART_TYPES = map[string]string {
    "page_views": "Total Page Views",
    "unique_visitors": "Unique Visitors",
    "signups": "Signups",
}

func main() {
	e := echo.New()
	e.HideBanner = true
  e.Debug = true

	// Routes (using templates)
	t := &Template{
		templates: template.Must(template.New("").Funcs(template.FuncMap{
    "dict": func(values ...interface{}) (map[string]interface{}, error) {
        if len(values)%2 != 0 {
            return nil, fmt.Errorf("invalid dict call")
        }
        dict := make(map[string]interface{}, len(values)/2)
        for i := 0; i < len(values); i+=2 {
            key, ok := values[i].(string)
            if !ok {
                return nil, fmt.Errorf("dict keys must be strings")
            }
            dict[key] = values[i+1]
        }
        return dict, nil
    },
}).ParseGlob("./templates/*.gohtml")),
	}
	e.Renderer = t
	e.GET("/", dashboard)
	e.GET("/charts", chart_view_hx)
  e.Static("/", "assets")

  // Start server
	err := e.Start(":8888")
  if err != nil {
    log.Println(err)
  }
}

// Handlers

func dashboard(c echo.Context) error {
  type Data struct {
    Title string
    Welcome string
    Charts []*charts.Line
  }
  chart_data := fake_chart_data(7, "Total Page Views")
  chart_id := c.FormValue("chart_id")
  // create a new line instance
	line := charts.NewLine()
	// set some global options like Title/Legend/ToolTip or anything else
	line.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{Theme: types.ThemeWesteros,ChartID: chart_id}),
		charts.WithTitleOpts(opts.Title{
			Title:    "Total Page Views",
		}))

	// Put data into instance
	line.SetXAxis(chart_data.X).
	AddSeries("Category A", chart_data.Y).
	SetSeriesOptions(charts.WithLineChartOpts(opts.LineChart{Smooth: true}))
  
  data := Data {
    Title: "htmx and echarts",
    Welcome: "Demo of using HTMX with echarts",
    Charts: []*charts.Line{line},
  }
  return c.Render(http.StatusOK,"dashboard", data)
}

//type ChartRenderer struct {
  
//}

//func (c *ChartRenderer) Render(w io.Writer)

func chart_view_hx(c echo.Context) error {
  //Returns chart options for echarts
  period := c.FormValue("period")
  if period == "" {
    period = "week"
  }
  chart_id := c.FormValue("chart_id")
  chart_type := c.FormValue("chart_type")
  if chart_type == "" {
    chart_type = "page_views"
  }

  days_in_period := make(map[string]int)
  days_in_period["week"] = 7
  days_in_period["month"] = 30
  filter_by := days_in_period[period]

  //simulate fetching this from your database
  chart_title := CHART_TYPES[chart_type]
  chart_data := fake_chart_data(filter_by, chart_title)

  // create a new line instance
	line := charts.NewLine()
	// set some global options like Title/Legend/ToolTip or anything else
	line.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{Theme: types.ThemeWesteros, ChartID: chart_id}),
		charts.WithTitleOpts(opts.Title{
			Title:    chart_title,
		}))

	// Put data into instance
	line.SetXAxis(chart_data.X).
	AddSeries("Category A", chart_data.Y).
	SetSeriesOptions(charts.WithLineChartOpts(opts.LineChart{Smooth: true}))
	line.Render(io.Writer(c.Response().Writer))
  return nil
}

/*
def fake_chart_data(num_days:int, chart_title:str):

  base = datetime.datetime.today()
  dates = [base - datetime.timedelta(days=x) for x in range(num_days)]
  dates = [d.strftime("%d %b") for d in dates]

  return {
    "x": dates,
    "y": [randint(10, 60) for x in range(num_days)],
     "chart_title": chart_title
  }
*/

type FakeData struct {
  X []string
  Y []opts.LineData
  Title string
}
func fake_chart_data(num_days int, chart_title string) FakeData {
  base := time.Now()
  var dates []string
  var y []opts.LineData
  for i:=0;i<num_days;i++ {
    dates = append(dates, base.AddDate(0,0,-1).Format("02 Jan"))
    y = append(y, opts.LineData{Value:rand.Intn(50) + 10})
  }
  
  return FakeData{
    X: dates,
    Y: y,
    Title: chart_title,
  }
}


