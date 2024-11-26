package main

import (
    "bufio"
    "encoding/json"
    "fmt"
    "time"
    "gonum.org/v1/plot"
    "gonum.org/v1/plot/plotter"
    "gonum.org/v1/plot/vg"
    "golang.org/x/image/font"
    "golang.org/x/image/font/basicfont"
    "golang.org/x/image/math/fixed"
    "image"
    "image/color"
    "image/draw"
    "image/png"
    "log"
    "os"
    "sort"
)

type Data struct {
    Seq       int    `json:"seq"`
    Code      int    `json:"code"`
    Latency   int    `json:"latency"`
    Timestamp string `json:"timestamp"`
}

type Metrics struct {
    successRate float64
    p99Latency  float64
    avgLatency  float64
    minLatency  float64
    maxLatency  float64
}

func main() {
    if len(os.Args) != 4 {
        log.Fatalf("Usage: %s <go_metrics.txt> <node_metrics.txt> <title>\n", os.Args[0])
    }
    goData := readMetrics(os.Args[1])
    nodeData := readMetrics(os.Args[2])
    title := os.Args[3]
    generateLineChart(goData, nodeData, title)
    generateBarCharts(goData, nodeData, title)
}

func readMetrics(filePath string) []Data {
    file, err := os.Open(filePath)
    if err != nil {
        log.Fatalf("Failed to open file: %v", err)
    }
    defer file.Close()
    var data []Data
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        var d Data
        if err := json.Unmarshal(scanner.Bytes(), &d); err != nil {
            log.Fatalf("Error parsing JSON: %v", err)
        }
        data = append(data, d)
    }
    return data
}

func averageBySecond(data []Data) map[int64]float64 {
    if len(data) == 0 {
        return nil
    }

    // Parse the initial timestamp
    initialTime, err := time.Parse(time.RFC3339Nano, data[0].Timestamp)
    if err != nil {
        log.Fatalf("Error parsing initial timestamp: %v", err)
    }

    secondMap := make(map[int64][]int)
    for _, d := range data {
        t, err := time.Parse(time.RFC3339Nano, d.Timestamp)
        if err != nil {
            log.Fatalf("Error parsing timestamp: %v", err)
        }

        // Calculate the seconds elapsed since the initial timestamp
        sec := int64(t.Sub(initialTime).Seconds())
        secondMap[sec] = append(secondMap[sec], d.Latency)
    }

    averages := make(map[int64]float64)
    for sec, latencies := range secondMap {
        sum := 0
        for _, latency := range latencies {
            sum += latency
        }
        avg := float64(sum) / float64(len(latencies))
        averages[sec] = avg
    }

    return averages
}

func generateLineChart(goData, nodeData []Data, title string) {
    p := plot.New()
    p.Title.Text = "Latency over seconds"
    p.X.Label.Text = "Seconds elapsed"
    p.Y.Label.Text = "Latency (ms)"

    groupedGoData := averageBySecond(goData)
    groupedNodeData := averageBySecond(nodeData)
    fmt.Printf("%d %d", len(groupedGoData), len(groupedNodeData))

    smallerDataSet:= min(len(groupedNodeData), len(groupedNodeData))

    goPoints := make(plotter.XYs, smallerDataSet)//len(goData))
    nodePoints := make(plotter.XYs, smallerDataSet)

    for i := 0; i < smallerDataSet; i++ {
        goPoints[i].X = float64(i)
        goLatency := float64(groupedGoData[int64(i)] / 1e6)
        goPoints[i].Y = goLatency

        nodePoints[i].X = float64(i)
        nodeLatency := float64(groupedNodeData[int64(i)] / 1e6)
        fmt.Printf("Node: %.2f, Go: %.2f\n", nodeLatency, goLatency)
        nodePoints[i].Y = nodeLatency
    }

    // Create line plots for Go and Node
    goLine, err := plotter.NewLine(goPoints)
    if err != nil {
        log.Fatalf("Failed to create Go line plot: %v", err)
    }
    goLine.Color = color.RGBA{0, 0, 255, 255} // Blue color for Go

    nodeLine, err := plotter.NewLine(nodePoints)
    if err != nil {
        log.Fatalf("Failed to create Node line plot: %v", err)
    }
    nodeLine.Color = color.RGBA{0, 128, 0, 255} // Green color for Node

    // Add line plots to the plot
    p.Add(goLine, nodeLine)

    // Set the legend
    p.Legend.Add("Golang", goLine)
    p.Legend.Add("Node", nodeLine)
    p.Legend.Top = true

    // Save the plot to a PNG file
    if err := p.Save(8*vg.Inch, 4*vg.Inch, "latencies_over_request_counter.png"); err != nil {
        log.Fatalf("Failed to save chart: %v", err)
    }
}

func generateBarCharts(goData, nodeData []Data, title string) {
    goMetrics := calculateMetrics(goData)
    nodeMetrics := calculateMetrics(nodeData)
    paths := []string{"success_rate.png", "p99_latency.png", "avg_latency.png", "min_latency.png", "max_latency.png"}
    createBarChart(paths[0], "Success Rate", []string{"Node.js", "Golang"}, []float64{nodeMetrics.successRate, goMetrics.successRate}, []color.Color{color.RGBA{0, 128, 0, 255}, color.RGBA{0, 0, 255, 255}}, true)
    createBarChart(paths[1], "p99 Latency", []string{"Node.js", "Golang"}, []float64{nodeMetrics.p99Latency, goMetrics.p99Latency}, []color.Color{color.RGBA{0, 128, 0, 255}, color.RGBA{0, 0, 255, 255}}, false)
    createBarChart(paths[2], "Average Latency", []string{"Node.js", "Golang"}, []float64{nodeMetrics.avgLatency, goMetrics.avgLatency}, []color.Color{color.RGBA{0, 128, 0, 255}, color.RGBA{0, 0, 255, 255}}, false)
    createBarChart(paths[3], "Minimum Latency", []string{"Node.js", "Golang"}, []float64{nodeMetrics.minLatency, goMetrics.minLatency}, []color.Color{color.RGBA{0, 128, 0, 255}, color.RGBA{0, 0, 255, 255}}, false)
    createBarChart(paths[4], "Maximum Latency", []string{"Node.js", "Golang"}, []float64{nodeMetrics.maxLatency, goMetrics.maxLatency}, []color.Color{color.RGBA{0, 128, 0, 255}, color.RGBA{0, 0, 255, 255}}, false)
    combineImagesWithTitle(paths, "output.png", title)
    for _, file := range paths {
        os.Remove(file)
    }
    fmt.Printf("%s\n\nNode\nSuccess Rate: %.2f%%\np99 Latency: %.2fms\nAverage latency: %.2fms\nMinimum Latency: %.2fms\nMaximum Latency: %.2fms\n\n", title, nodeMetrics.successRate, nodeMetrics.p99Latency, nodeMetrics.avgLatency, nodeMetrics.minLatency, nodeMetrics.maxLatency)
    fmt.Printf("Go\nSuccess Rate: %.2f%%\np99 Latency: %.2fms\nAverage latency: %.2fms\nMinimum Latency: %.2fms\nMaximum Latency: %.2fms", goMetrics.successRate, goMetrics.p99Latency, goMetrics.avgLatency, goMetrics.minLatency, goMetrics.maxLatency)
}

func createBarChart(filename, title string, labels []string, values []float64, barColors []color.Color, hasLegend bool) {
    p := plot.New()
    p.Title.Text = title

    w := vg.Points(20) // Width of the bars
    for i, value := range values {
        // Each value needs to be in a slice to create a bar chart
        bars, err := plotter.NewBarChart(plotter.Values{value}, w)
        if err != nil {
            log.Fatalf("Failed to create bar chart: %v", err)
        }

        // Set the color for each bar
        bars.Color = barColors[i]

        if i == 0 {
            bars.Offset = 0
        } else {
            bars.Offset = w + 10
        }
        p.Add(bars)
        if hasLegend {
            p.Legend.Add(labels[i], bars)
        }
    }
    if hasLegend {
        p.Legend.Top = true
    }
    p.X.Min = 0
    p.X.Max = 1
    p.Y.Label.Text = "Latency (ms)"
    p.NominalX("")

    // Save the plot to a PNG file
    if err := p.Save(4*vg.Inch, 4*vg.Inch, filename); err != nil {
        log.Fatalf("Failed to save chart: %v", err)
    }
}


func calculateMetrics(data []Data) (metrics Metrics) {
    var latencies []float64
    for _, d := range data {
        latencies = append(latencies, float64(d.Latency)/1e6)
        if d.Code >= 200 && d.Code < 300 {
            metrics.successRate++
        }
    }
    metrics.successRate = (metrics.successRate / float64(len(data))) * 100
    sort.Float64s(latencies)
    metrics.p99Latency = percentile(latencies, 99)
    metrics.avgLatency = average(latencies)
    metrics.minLatency = latencies[0]
    metrics.maxLatency = latencies[len(latencies)-1]
    return metrics
}

func percentile(sortedData []float64, perc float64) float64 {
    if len(sortedData) == 0 {
        return 0
    }
    index := perc / 100 * float64(len(sortedData)-1)
    lower := sortedData[int(index)]
    upper := sortedData[int(index)+1]
    return lower + (upper-lower)*(index-float64(int(index)))
}

func average(data []float64) float64 {
    sum := 0.0
    for _, v := range data {
        sum += v
    }
    return sum / float64(len(data))
}


func loadImage(filePath string) (image.Image, error) {
    file, err := os.Open(filePath)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    img, err := png.Decode(file)
    if err != nil {
        return nil, err
    }

    return img, nil
}

func addLabel(img *image.RGBA, label string) {
    col := color.RGBA{0, 0, 0, 255} // Black color for the title
    face := basicfont.Face7x13

    // Calculate the width of the text
    textWidth := font.MeasureString(face, label).Ceil()

    // Calculate the starting position (centered)
    x := (img.Bounds().Dx() - textWidth) / 2
    y := face.Metrics().Ascent.Ceil() + 20 // 10 for padding from the top

    point := fixed.Point26_6{fixed.Int26_6(x * 64), fixed.Int26_6(y * 64)}

    d := &font.Drawer{
        Dst:  img,
        Src:  image.NewUniform(col),
        Face: face,
        Dot:  point,
    }
    d.DrawString(label)
}

func drawLine(img *image.RGBA, start, end image.Point, col color.Color) {
    for x := start.X; x < end.X; x++ {
        for y := start.Y; y < end.Y; y++ {
            img.Set(x, y, col)
        }
    }
}

func min(values ...int) int {
    minVal := values[0]
    for _, v := range values[1:] {
        if v < minVal {
            minVal = v
        }
    }
    return minVal
}

func max(values ...int) int {
    maxVal := values[0]
    for _, v := range values[1:] {
        if v > maxVal {
            maxVal = v
        }
    }
    return maxVal
}

func combineImagesWithTitle(imagePaths []string, outputPath, title string) error {
    images := make([]image.Image, 5)
    var err error
    for i, path := range imagePaths {
        images[i], err = loadImage(path)
        if err != nil {
            return err
        }
    }

    // Spaces, title height, and line width
    space := 10 // Space between images and for lines
    titleHeight := 40 // Space for the title
    lineWidth := 2 // Width of the line

    // Calculate the width for rows 2 and 3, including spaces
    row2Width := images[1].Bounds().Dx() + space + images[2].Bounds().Dx()
    row3Width := images[3].Bounds().Dx() + space + images[4].Bounds().Dx()
    canvasWidth := max(images[0].Bounds().Dx(), max(row2Width, row3Width))

    // Calculate total height considering spaces, title, and lines
    canvasHeight := images[0].Bounds().Dy() + max(images[1].Bounds().Dy(), images[2].Bounds().Dy()) + max(images[3].Bounds().Dy(), images[4].Bounds().Dy()) + space*4 + titleHeight + lineWidth*2

    combinedImg := image.NewRGBA(image.Rect(0, 0, canvasWidth, canvasHeight))

    // Fill the background with white
    draw.Draw(combinedImg, combinedImg.Bounds(), &image.Uniform{color.White}, image.Point{}, draw.Src)

    // Draw the title on top
    addLabel(combinedImg, title)

    // Draw image A
    startY := titleHeight + space
    draw.Draw(combinedImg, image.Rect(0, startY, images[0].Bounds().Dx(), startY+images[0].Bounds().Dy()), images[0], image.Point{}, draw.Src)

    // Draw horizontal line after image A
    drawLine(combinedImg, image.Point{0, startY + images[0].Bounds().Dy()}, image.Point{canvasWidth, startY + images[0].Bounds().Dy() + lineWidth}, color.Black)

    // Draw images B and C
    yStartRow2 := startY + images[0].Bounds().Dy() + space + lineWidth
    draw.Draw(combinedImg, image.Rect(0, yStartRow2, images[1].Bounds().Dx(), yStartRow2+images[1].Bounds().Dy()), images[1], image.Point{}, draw.Src)

    // Draw vertical line between B and C
    drawLine(combinedImg, image.Point{images[1].Bounds().Dx(), yStartRow2 - space}, image.Point{images[1].Bounds().Dx() + lineWidth, yStartRow2 + images[1].Bounds().Dy()}, color.Black)
    draw.Draw(combinedImg, image.Rect(images[1].Bounds().Dx() + space + lineWidth, yStartRow2, row2Width, yStartRow2+images[2].Bounds().Dy()), images[2], image.Point{}, draw.Src)

    // Draw horizontal line after row 2
    yEndRow2 := yStartRow2 + max(images[1].Bounds().Dy(), images[2].Bounds().Dy())
    drawLine(combinedImg, image.Point{0, yEndRow2}, image.Point{canvasWidth, yEndRow2 + lineWidth}, color.Black)

    // Draw images D and E
    yStartRow3 := yEndRow2 + space + lineWidth
    draw.Draw(combinedImg, image.Rect(0, yStartRow3, images[3].Bounds().Dx(), yStartRow3+images[3].Bounds().Dy()), images[3], image.Point{}, draw.Src)

    // Draw vertical line between D and E
    drawLine(combinedImg, image.Point{images[3].Bounds().Dx(), yStartRow3 - space}, image.Point{images[3].Bounds().Dx() + lineWidth, yStartRow3 + images[3].Bounds().Dy()}, color.Black)
    draw.Draw(combinedImg, image.Rect(images[3].Bounds().Dx() + space + lineWidth, yStartRow3, row3Width, yStartRow3+images[4].Bounds().Dy()), images[4], image.Point{}, draw.Src)

    outputFile, err := os.Create(outputPath)
    if err != nil {
        return err
    }
    defer outputFile.Close()

    return png.Encode(outputFile, combinedImg)
}

