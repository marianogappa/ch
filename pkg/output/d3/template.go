package d3

import (
	"html/template"
	"log"
)

var Templates = map[string]*template.Template{}

func init() {
	var err error
	for name, content := range templateStrings {
		if Templates[name], err = template.New(name).Parse(content); err != nil {
			log.Fatalf("d3.templates: error parsing text template: %v", err)
		}
	}
}

var templateStrings = map[string]string{
	"base": `<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <title>{{.Title}}</title>
    <script src="https://d3js.org/d3.v7.min.js"></script>
    <style>
        body { font-family: sans-serif; margin: 20px; }
        .chart-container { max-width: 900px; margin: 0 auto; }
        .axis-label { font-size: 12px; font-weight: bold; }
        .tooltip {
            position: absolute;
            text-align: center;
            padding: 6px;
            font: 12px sans-serif;
            background: white;
            border: 1px solid #ccc;
            border-radius: 4px;
            pointer-events: none;
            opacity: 0;
            box-shadow: 0 2px 4px rgba(0,0,0,0.1);
        }
    </style>
</head>
<body>
    <div class="chart-container">
        <h1>{{.Title}}</h1>
        <div id="chart"></div>
    </div>
    <script>
        const data = {{.Data}};
        const config = {{.Config}};
        
        const margin = {top: 40, right: 40, bottom: 60, left: 60};
        const width = 800 - margin.left - margin.right;
        const height = 500 - margin.top - margin.bottom;

        const svg = d3.select("#chart")
            .append("svg")
            .attr("width", width + margin.left + margin.right)
            .attr("height", height + margin.top + margin.bottom)
            .append("g")
            .attr("transform", "translate(" + margin.left + "," + margin.top + ")");

        // Tooltip
        const tooltip = d3.select("body").append("div")
            .attr("class", "tooltip");

        const showTooltip = function(event, d, text) {
            tooltip.transition().duration(200).style("opacity", .9);
            tooltip.html(text)
                .style("left", (event.pageX + 10) + "px")
                .style("top", (event.pageY - 28) + "px");
        };
        const hideTooltip = function(d) {
            tooltip.transition().duration(500).style("opacity", 0);
        };

        {{.ChartScript}}
    </script>
</body>
</html>`,

	"bar": `
        // Bar Chart
        const x = d3.scaleBand()
            .range([0, width])
            .padding(0.1);
        const y = d3.scaleLinear()
            .range([height, 0]);

        x.domain(data.map(d => d.label));
        y.domain([0, d3.max(data, d => d.value)]);

        svg.append("g")
            .attr("transform", "translate(0," + height + ")")
            .call(d3.axisBottom(x))
            .selectAll("text")
            .style("text-anchor", "end")
            .attr("dx", "-.8em")
            .attr("dy", ".15em")
            .attr("transform", "rotate(-45)");

        svg.append("g")
            .call(d3.axisLeft(y));

        svg.selectAll(".bar")
            .data(data)
            .enter().append("rect")
            .attr("class", "bar")
            .attr("x", d => x(d.label))
            .attr("width", x.bandwidth())
            .attr("y", d => y(d.value))
            .attr("height", d => height - y(d.value))
            .attr("fill", config.color || "steelblue")
            .on("mouseover", function(event, d) { showTooltip(event, d, d.label + ": " + d.value); })
            .on("mouseout", hideTooltip);
            
        // Labels
        if (config.xLabel) {
            svg.append("text")
                .attr("class", "axis-label")
                .attr("text-anchor", "middle")
                .attr("x", width/2)
                .attr("y", height + margin.bottom - 5)
                .text(config.xLabel);
        }
        if (config.yLabel) {
            svg.append("text")
                .attr("class", "axis-label")
                .attr("text-anchor", "middle")
                .attr("transform", "rotate(-90)")
                .attr("y", -margin.left + 20)
                .attr("x", -height/2)
                .text(config.yLabel);
        }
    `,

	"pie": `
        // Pie Chart
        const radius = Math.min(width, height) / 2;
        const pieSvg = svg.append("g")
            .attr("transform", "translate(" + width / 2 + "," + height / 2 + ")");

        const color = d3.scaleOrdinal(d3.schemeCategory10);

        const pie = d3.pie()
            .value(d => d.value);

        const path = d3.arc()
            .outerRadius(radius - 10)
            .innerRadius(0);

        const label = d3.arc()
            .outerRadius(radius - 40)
            .innerRadius(radius - 40);

        const arc = pieSvg.selectAll(".arc")
            .data(pie(data))
            .enter().append("g")
            .attr("class", "arc");

        arc.append("path")
            .attr("d", path)
            .attr("fill", d => color(d.data.label))
            .on("mouseover", function(event, d) { showTooltip(event, d, d.data.label + ": " + d.data.value + " (" + Math.round((d.endAngle - d.startAngle)/(2*Math.PI)*100) + "%)"); })
            .on("mouseout", hideTooltip);

        arc.append("text")
            .attr("transform", d => "translate(" + label.centroid(d) + ")")
            .attr("dy", "0.35em")
            .text(d => d.data.label);
    `,

	"scatter": `
        // Scatter Plot
        const x = d3.scaleLinear()
            .range([0, width]);
        const y = d3.scaleLinear()
            .range([height, 0]);

        x.domain(d3.extent(data, d => d.x)).nice();
        y.domain(d3.extent(data, d => d.y)).nice();

        svg.append("g")
            .attr("transform", "translate(0," + height + ")")
            .call(d3.axisBottom(x));

        svg.append("g")
            .call(d3.axisLeft(y));

        svg.selectAll(".dot")
            .data(data)
            .enter().append("circle")
            .attr("class", "dot")
            .attr("r", 3.5)
            .attr("cx", d => x(d.x))
            .attr("cy", d => y(d.y))
            .style("fill", config.color || "steelblue")
            .on("mouseover", function(event, d) { showTooltip(event, d, "(" + d.x + ", " + d.y + ")"); })
            .on("mouseout", hideTooltip);

        // Labels
        if (config.xLabel) {
            svg.append("text")
                .attr("class", "axis-label")
                .attr("text-anchor", "middle")
                .attr("x", width/2)
                .attr("y", height + margin.bottom - 5)
                .text(config.xLabel);
        }
        if (config.yLabel) {
            svg.append("text")
                .attr("class", "axis-label")
                .attr("text-anchor", "middle")
                .attr("transform", "rotate(-90)")
                .attr("y", -margin.left + 20)
                .attr("x", -height/2)
                .text(config.yLabel);
        }
    `,

	"histogram": `
        // Histogram
        const x = d3.scaleLinear()
            .domain(d3.extent(data, d => d.value))
            .range([0, width]);
            
        svg.append("g")
            .attr("transform", "translate(0," + height + ")")
            .call(d3.axisBottom(x));

        const histogram = d3.histogram()
            .value(d => d.value)
            .domain(x.domain())
            .thresholds(x.ticks(20));

        const bins = histogram(data);

        const y = d3.scaleLinear()
            .range([height, 0]);
            
        y.domain([0, d3.max(bins, d => d.length)]);

        svg.append("g")
            .call(d3.axisLeft(y));

        svg.selectAll("rect")
            .data(bins)
            .enter().append("rect")
            .attr("x", 1)
            .attr("transform", d => "translate(" + x(d.x0) + "," + y(d.length) + ")")
            .attr("width", d => Math.max(0, x(d.x1) - x(d.x0) - 1))
            .attr("height", d => height - y(d.length))
            .style("fill", config.color || "steelblue")
            .on("mouseover", function(event, d) { showTooltip(event, d, "Range: " + d.x0 + " - " + d.x1 + "<br>Count: " + d.length); })
            .on("mouseout", hideTooltip);
            
        // Labels
        if (config.xLabel) {
            svg.append("text")
                .attr("class", "axis-label")
                .attr("text-anchor", "middle")
                .attr("x", width/2)
                .attr("y", height + margin.bottom - 5)
                .text(config.xLabel);
        }
        if (config.yLabel) {
            svg.append("text")
                .attr("class", "axis-label")
                .attr("text-anchor", "middle")
                .attr("transform", "rotate(-90)")
                .attr("y", -margin.left + 20)
                .attr("x", -height/2)
                .text(config.yLabel);
        }
    `,
}
