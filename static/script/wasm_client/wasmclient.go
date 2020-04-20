// +build js, wasm

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strings"
	"syscall/js"
	"wasm/bindings/chartjs"
	"wasm/bindings/dom"
	"wasm/utils"
	"wowstatistician/models"
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
	stats := models.Stats{}
	json.NewDecoder(resp.Body).Decode(&stats)
	return stats
}

func genLabelsExpanded(stats models.Stats) []string {
	labels := []string{}
	for _, distribution := range stats.Distributions {
		for _, spec := range distribution.Specs {
			label := fmt.Sprintf("%v - %v", spec.Spec, distribution.Class)
			labels = append(labels, label)
		}
	}
	return labels
}

func genLabelsMerged(stats models.Stats) []string {
	labels := []string{}
	for _, distribution := range stats.Distributions {
		labels = append(labels, distribution.Class)
	}
	return labels
}

func genDataExpanded(stats models.Stats) []interface{} {
	data := []interface{}{}
	for _, distribution := range stats.Distributions {
		for _, spec := range distribution.Specs {
			data = append(data, spec.Count)
		}
	}
	return data
}

func genDataMerged(stats models.Stats) []interface{} {
	data := []interface{}{}
	for _, distribution := range stats.Distributions {
		data = append(data, distribution.Total)
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

func genColorsExpanded(labels []string) []string {
	colors := []string{}
	for _, v := range labels {
		classraw := strings.Split(v, " - ")
		class := classraw[len(classraw)-1]
		classkey := strings.Replace(class, " ", "_", -1)
		colors = append(colors, classColors[classkey])
	}
	return colors
}

func genColorsMerged(labels []string) []string {
	colors := []string{}
	for _, v := range labels {
		classkey := strings.Replace(v, " ", "_", -1)
		colors = append(colors, classColors[classkey])
	}
	return colors
}

func makeConfig(stats models.Stats, merged bool) *chartjs.Config {
	config := &chartjs.Config{
		Type: "bar",
		Data: makeData(stats, merged),
		Options: &chartjs.Options{
			Scales: &chartjs.Scales{
				YAxes: []*chartjs.Axe{
					{
						Type: "linear",
						GridLines: &chartjs.GridLines{
							Display: utils.Bool(false),
						},
						ScaleLabel: &chartjs.ScaleLabel{
							Display: utils.Bool(false),
						},
						Ticks: &chartjs.Ticks{
							Display:     utils.Bool(false),
							BeginAtZero: utils.Bool(true),
						},
					},
				},
				XAxes: []*chartjs.Axe{
					{
						Type: "category",
						GridLines: &chartjs.GridLines{
							Display:         utils.Bool(true),
							DrawBorder:      utils.Bool(true),
							DrawOnChartArea: utils.Bool(false),
							DrawTicks:       utils.Bool(true),
						},
					},
				},
			},
			Responsive:          utils.Bool(true),
			MaintainAspectRatio: utils.Bool(false),
			Title: &chartjs.Title{
				Display: utils.Bool(true),
				Text:    fmt.Sprintf("Total players for: %v - %v", strings.Title(stats.Source), stats.Overall),
			},
		},
	}
	return config
}

func makeData(stats models.Stats, merged bool) *chartjs.Data {
	labels := []string{}
	data := []interface{}{}
	colors := []string{}
	if merged {
		labels = genLabelsMerged(stats)
		data = genDataMerged(stats)
		labels, data = minToMax(labels, data)
		colors = genColorsMerged(labels)
	} else {
		labels = genLabelsExpanded(stats)
		data = genDataExpanded(stats)
		labels, data = minToMax(labels, data)
		colors = genColorsExpanded(labels)
	}
	return &chartjs.Data{
		Labels: labels,
		Datasets: []*chartjs.Dataset{
			{
				Label:           "Number of players",
				Data:            data,
				BackgroundColor: colors,
			},
		},
	}
}

func makeChart(source string, merged bool) {
	chart := chartjs.GetChart("statsChart")
	ctx := dom.Document().GetElementById("stats").GetContext("2d")
	if chart.Value.IsUndefined() {
		stats := getStats(source)
		config := makeConfig(stats, merged)
		chart = chartjs.NewChart(ctx, config)
		chart.Register("statsChart")
	} else {
		chart.Detroy()
		stats := getStats(source)
		config := makeConfig(stats, merged)
		chart = chartjs.NewChart(ctx, config)
		chart.Register("statsChart")
	}
}

func isMerged() bool {
	var text string
	radios := dom.Document().GetElementsByName("merge")
	for i := 0; i < len(radios); i++ {
		if radios[i].Get("checked").Bool() {
			text = radios[i].Get("value").String()

		}
	}
	switch text {
	case "spec":
		return false
	case "class":
		return true
	default:
		return false
	}
}

func setDate(source string) {
	stats := getStats(source)
	dom.Document().GetElementById("syncdate").SetInnerHTML(stats.SyncDate)
}

func dropDownCallback(this js.Value, args []js.Value) interface{} {
	go func() {
		dropDown := dom.Document().GetElementById("source")
		selected := dropDown.Get("selectedIndex").Int()
		source := dropDown.Get("options").Index(selected).Get("text").String()
		merged := isMerged()
		makeChart(strings.ToLower(source), merged)
		setDate(strings.ToLower(source))
	}()
	return nil
}

func main() {
	dropDownCallback(js.Null(), nil)
	dom.Document().GetElementById("source").AddEventListener("change", dropDownCallback)
	radios := dom.Document().GetElementsByName("merge")
	for i := 0; i < len(radios); i++ {
		radios[i].AddEventListener("change", dropDownCallback)
	}
	select {}
}
