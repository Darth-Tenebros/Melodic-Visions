package render_charts

import (
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

func generateBarItems(values []int) []opts.BarData {
	items := make([]opts.BarData, 0)
	for i := 0; i < len(values); i++ {
		items = append(items, opts.BarData{Value: values[i]})
	}
	return items
}

func BarBasic(keys []string, values []int) *charts.Bar {
	bar := charts.NewBar()
	bar.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{Title: "basic bar example", Subtitle: "This is the subtitle."}),
	)

	bar.SetXAxis(keys).
		AddSeries("Artists", generateBarItems(values))
	return bar
}
