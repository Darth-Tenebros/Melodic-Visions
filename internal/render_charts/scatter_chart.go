package render_charts

import (
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

func generateScatterItems(featureSet []float64) []opts.ScatterData {
	items := make([]opts.ScatterData, 0)

	for _, feature := range featureSet {
		items = append(items, opts.ScatterData{
			Value:        feature,
			Symbol:       "roundRect",
			SymbolSize:   5,
			SymbolRotate: 10,
		})
	}
	return items
}
func ScatterBasic(tracksIds []string, acousticness, danceability, valence []float64) *charts.Scatter {
	scatter := charts.NewScatter()
	scatter.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{Title: "basic scatter example"}),
		charts.WithXAxisOpts(opts.XAxis{Show: opts.Bool(false)}),
	)

	scatter.SetXAxis(tracksIds).
		AddSeries("acousticness", generateScatterItems(acousticness)).
		AddSeries("danceability", generateScatterItems(danceability)).
		AddSeries("valence", generateScatterItems(valence))

	return scatter
}
