package main

import (
	"encoding/json"
	"fmt"
	"js/bindings/chartjs"
	"log"
	"net/http"
	"sort"
	"strings"
	"wowstatistician/models"

	"github.com/gopherjs/gopherjs/js"
)

var classColors = map[string]string{
	"Warrior":      "rgba(229, 130, 80, 1)",
	"Paladin":      "rgba(216, 119, 159, 1)",
	"Hunter":       "rgba(154, 178, 107, 1)",
	"Rogue":        "rgba(255, 204, 102, 1)",
	"Priest":       "rgba(206, 214, 229, 1)",
	"Death_Knight": "rgba(229, 103, 103, 1)",
	"Shaman":       "rgba(81, 132, 204, 1)",
	"Mage":         "rgba(109, 186, 242, 1)",
	"Warlock":      "rgba(142, 122, 204, 1)",
	"Druid":        "rgba(255, 178, 102, 1)",
	"Monk":         "rgba(71, 178, 169, 1)",
	"Demon_Hunter": "rgba(178, 107, 178, 1)",
}

func getStats(source string) models.Stats {
	resp, err := http.Get("/stats/" + source)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()
	var stats models.Stats
	json.NewDecoder(resp.Body).Decode(&stats)
	return stats
}

func genLabels(stats models.Stats) []string {
	var labels []string
	for _, distribution := range stats.Distributions {
		for _, spec := range distribution.Specs {
			label := fmt.Sprintf("%v - %v", spec.Spec, distribution.Class)
			labels = append(labels, label)
		}
	}
	return labels
}

func genData(stats models.Stats) []interface{} {
	var data []interface{}
	for _, distribution := range stats.Distributions {
		for _, spec := range distribution.Specs {
			data = append(data, spec.Count)
		}
	}
	return data
}

func minToMax(labels []string, data []interface{}) ([]string, []interface{}) {
	type ObjArray struct {
		Label string
		Count interface{}
	}
	objArray := []ObjArray{}
	for i := range labels {
		objArray = append(objArray, ObjArray{Label: labels[i], Count: data[i]})
	}
	sort.SliceStable(objArray, func(i, j int) bool {
		switch objArray[i].Count.(type) {
		case int:
			return objArray[i].Count.(int) < objArray[j].Count.(int)
		case float64:
			return objArray[i].Count.(float64) < objArray[j].Count.(float64)
		default:
			return false
		}
	})
	sortedLabels := []string{}
	sortedData := []interface{}{}
	for _, v := range objArray {
		sortedLabels = append(sortedLabels, v.Label)
		sortedData = append(sortedData, v.Count)
	}
	return sortedLabels, sortedData
}

func genColors(labels []string) []string {
	colors := []string{}
	for _, v := range labels {
		classraw := strings.Split(v, " - ")
		class := classraw[len(classraw)-1]
		classkey := strings.Replace(class, " ", "_", -1)
		colors = append(colors, classColors[classkey])
	}
	return colors
}

func makeConfig(stats models.Stats) *chartjs.Config {
	data := makeData(stats)
	yTicks := chartjs.NewTicks().SetBeginAtZero(true).SetDisplay(false)
	yScaleLabel := chartjs.NewScaleLabel().SetDisplay(false)
	yGridLines := chartjs.NewGridLines().SetDisplay(false)
	yAxe := chartjs.NewAxe("linear").AddTicks(yTicks).AddScaleLabel(yScaleLabel).AddGridLines(yGridLines)
	xGridlines := chartjs.NewGridLines().SetDisplay(true).SetDrawBorder(true).SetDrawTicks(true).SetDrawOnChartArea(false)
	xAxe := chartjs.NewAxe("category").AddGridLines(xGridlines)
	scales := chartjs.NewScales().AddYAxe(yAxe).AddXAxe(xAxe)
	title := chartjs.NewTitle(fmt.Sprintf("Total players for: %v - %v", stats.Source, stats.Overall)).SetDisplay(true)
	options := chartjs.NewOptions().AddScales(scales).AddTitle(title).SetResponsive(true).SetMaintainAspectRatio(false)
	config := chartjs.NewConfig("bar", data, options)
	return config
}

func makeData(stats models.Stats) *chartjs.Data {
	labels := genLabels(stats)
	data := genData(stats)
	sortedLabels, sortedData := minToMax(labels, data)
	colors := genColors(sortedLabels)
	dataset := chartjs.NewDataset("Number of players", sortedData).SetBackgroundColor(colors)
	return chartjs.NewData().SetLabels(sortedLabels).AddDataset(dataset)
}

func makeChart(source string) {
	chart := chartjs.GetChart("statsChart")
	ctx := js.Global.Get("document").Call("getElementById", "stats").Call("getContext", "2d")
	if chart.Object == js.Undefined {
		stats := getStats(source)
		config := makeConfig(stats)
		chart := chartjs.NewChart(ctx, config)
		js.Global.Get("document").Set("statsChart", chart)
	} else {
		chart.Destroy()
		stats := getStats(source)
		config := makeConfig(stats)
		chart := chartjs.NewChart(ctx, config)
		js.Global.Get("document").Set("statsChart", chart)
	}
}

func dropDownCallback(event *js.Object) {
	go func() {
		dropDown := js.Global.Get("document").Call("getElementById", "source")
		selected := dropDown.Get("selectedIndex").Int()
		source := dropDown.Get("options").Index(selected).Get("text").String()
		makeChart(strings.ToLower(source))
	}()
}

func main() {
	dropDownCallback(nil)
	js.Global.Get("document").Call("getElementById", "source").Call("addEventListener", "change", dropDownCallback)
}
