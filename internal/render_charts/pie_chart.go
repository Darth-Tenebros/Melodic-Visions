package render_charts

import (
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

func generatePieItems(data map[string]int) []opts.PieData {
	items := make([]opts.PieData, 0)
	for key, value := range data {
		items = append(items, opts.PieData{Name: key, Value: value})
	}
	return items
}

func PieBasic(data map[string]int) *charts.Pie {
	pie := charts.NewPie()
	pie.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{Title: "basic pie example"}),
	)

	pie.AddSeries("pie", generatePieItems(data)).
		SetSeriesOptions(
			charts.WithLabelOpts(
				opts.Label{
					Show:      opts.Bool(true),
					Formatter: "{b}: {c}",
				},
			),
		)
	return pie
}
