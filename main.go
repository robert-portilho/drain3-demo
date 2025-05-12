package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/jaeyo/go-drain3/pkg/drain3"
)

type LogEntry struct {
	Timestamp time.Time
	Message   string
}

func main() {
	drain, _ := drain3.NewDrain(drain3.WithExtraDelimiter([]string{"_"}))

	miner := drain3.NewTemplateMiner(drain, drain3.NewMemoryPersistence())

	timeRegex := regexp.MustCompile(`^\[(.*?)\]`)
	layout := "Mon Jan 02 15:04:05 2006"

	counts := make(map[string]map[int64]int)

	file, err := os.Open("exemplo.log")
	if err != nil {
		fmt.Println("Erro ao abrir o arquivo:", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	ctx := context.Background()
	for scanner.Scan() {
		log := scanner.Text()
		timeMatch := timeRegex.FindStringSubmatch(log)
		if len(timeMatch) < 2 {
			continue
		}

		message := strings.TrimSpace(strings.ReplaceAll(log, timeMatch[0], ""))

		_, cluster, _, _, _ := miner.AddLogMessage(ctx, message)

		logTime, _ := time.Parse(layout, timeMatch[1])
		hour := logTime.Format("2006-01-02 15")

		if _, ok := counts[hour]; !ok {
			counts[hour] = make(map[int64]int)
		}

		counts[hour][cluster.ClusterId]++

	}

	templates := make(map[int64]string)
	for _, cluster := range drain.GetClusters() {
		templates[cluster.ClusterId] = cluster.GetTemplate()
	}

	keys := make([]string, 0, len(counts))

	for k := range counts {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	for _, hour := range keys {
		fmt.Printf("Hora: %s:00\n", hour)
		for clusterID, count := range counts[hour] {
			fmt.Printf("Cluster (%d) %s: %d\n", clusterID, templates[clusterID], count)
		}
	}

	fmt.Println("----------------------- TOTAL ------------------------")
	for _, cluster := range drain.GetClusters() {
		templates[cluster.ClusterId] = cluster.GetTemplate()
		fmt.Printf("Cluster (%d) %s: %d\n", cluster.ClusterId, cluster.GetTemplate(), cluster.Size)
	}
}
