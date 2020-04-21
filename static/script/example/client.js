const classColors = {
	Warrior: "rgba(229, 130, 80, 1)",
	Paladin: "rgba(216, 119, 159, 1)",
	Hunter: "rgba(154, 178, 107, 1)",
	Rogue: "rgba(255, 204, 102, 1)",
	Priest: "rgba(206, 214, 229, 1)",
	Death_Knight: "rgba(229, 103, 103, 1)",
	Shaman: "rgba(81, 132, 204, 1)",
	Mage: "rgba(109, 186, 242, 1)",
	Warlock: "rgba(142, 122, 204, 1)",
	Druid: "rgba(255, 178, 102, 1)",
	Monk: "rgba(71, 178, 169, 1)",
	Demon_Hunter: "rgba(178, 107, 178, 1)",
};

function getStats(source) {
	return new Promise((resolve) => {
		fetch(`/stats/${source}`)
			.then((response) => {
				return response.json();
			})
			.then((json) => {
				resolve(json);
			});
	});
}

function genLabels(json) {
	let labels = [];
	json.distributions.forEach((distribution) => {
		distribution.specs.forEach((spec) => {
			labels.push(`${spec.spec} - ${distribution.class}`);
		});
	});
	return labels;
}

function genData(json) {
	let data = [];
	json.distributions.forEach((distribution) => {
		distribution.specs.forEach((spec) => {
			data.push(spec.count);
		});
	});
	return data;
}

function minToMax(labels, data) {
	let objArray = labels.map((value, index) => {
		return {
			label: value,
			data: data[index] || 0,
		};
	});
	objArray.sort((a, b) => {
		return a.data - b.data;
	});
	let sortedLabels = [];
	let sortedData = [];
	objArray.forEach((obj) => {
		sortedLabels.push(obj.label);
		sortedData.push(obj.data);
	});
	return {
		labels: sortedLabels,
		data: sortedData,
	};
}

function genColors(labels) {
	let colors = [];
	labels.forEach((label) => {
		classraw = label.split(" - ").pop();
		classkey = classraw.replace(/\s/g, "_");
		colors.push(classColors[classkey]);
	});
	return colors;
}

function makeConfig(stats) {
	let data = makeData(stats);
	return {
		type: "bar",
		data: data,
		options: {
			scales: {
				yAxes: [
					{
						type: "linear",
						gridLines: {
							display: false,
						},
						scaleLabel: {
							display: false,
						},
						ticks: {
							display: false,
							beginAtZero: true,
						},
					},
				],
				xAxes: [
					{
						type: "category",
						gridLines: {
							display: true,
							drawBorder: true,
							drawOnChartArea: false,
							drawTicks: true,
						},
					},
				],
			},
			responsive: true,
			maintainAspectRatio: false,
			title: {
				display: true,
				text: `Total players for: ${stats.source} - ${stats.overall}`,
			},
		},
	};
}

function makeData(stats) {
	let labels = genLabels(stats);
	let data = genData(stats);
	let sorted = minToMax(labels, data);
	let colors = genColors(sorted.labels);
	return {
		labels: sorted.labels,
		datasets: [
			{
				label: "Number of players",
				data: sorted.data,
				backgroundColor: colors,
			},
		],
	};
}

function makeTitle(stats) {
	return `Total players for: ${stats.source} - ${stats.overall}`;
}

async function makeChart(source) {
	if (document.statsChart) {
		stats = await getStats(source);
		document.statsChart.data = makeData(stats);
		document.statsChart.options.title.text = makeTitle(stats);
		document.statsChart.update();
	} else {
		stats = await getStats(source);
		let ctx = document.getElementById("stats").getContext("2d");
		let config = makeConfig(stats);
		document.statsChart = new Chart(ctx, config);
	}
}

function dropDownCallback() {
	let dropDown = document.getElementById("source");
	let selected = dropDown.selectedIndex;
	let source = dropDown.options[selected].text;
	makeChart(source.toLowerCase());
}

dropDownCallback();
document.getElementById("source").addEventListener("change", dropDownCallback);
