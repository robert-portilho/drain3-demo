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

const (
	logFile     = "exemplo.log"
	timeLayout  = "Mon Jan 02 15:04:05 2006"
	timePattern = `^\[(.*?)\]`
)

func main() {
	// Inicializa Drain com delimitador extra
	drain, err := drain3.NewDrain(drain3.WithExtraDelimiter([]string{"_"}))
	if err != nil {
		fmt.Println("Erro ao inicializar Drain:", err)
		return
	}

	miner := drain3.NewTemplateMiner(drain, drain3.NewMemoryPersistence())

	counts := make(map[string]map[int64]int)
	templates := make(map[int64]string)

	file, err := os.Open(logFile)
	if err != nil {
		fmt.Println("Erro ao abrir o arquivo:", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	timeRegex := regexp.MustCompile(timePattern)
	ctx := context.Background()

	// Processamento linha a linha
	for scanner.Scan() {
		logLine := scanner.Text()
		timestamp, message, err := parseLogLine(logLine, timeRegex)
		if err != nil {
			continue // Ignora linhas inválidas
		}

		cluster, err := addToCluster(miner, ctx, message)
		if err != nil {
			fmt.Println("Erro ao processar log:", err)
			continue
		}

		hourKey := timestamp.Format("2006-01-02 15")
		if counts[hourKey] == nil {
			counts[hourKey] = make(map[int64]int)
		}
		counts[hourKey][cluster.ClusterId]++
	}

	// Mapeia todos os templates conhecidos
	for _, cluster := range drain.GetClusters() {
		templates[cluster.ClusterId] = cluster.GetTemplate()
	}

	printHourlyCounts(counts, templates)
	printTotalSummary(drain)
}

// parseLogLine extrai o timestamp e a mensagem da linha do log
func parseLogLine(line string, re *regexp.Regexp) (time.Time, string, error) {
	matches := re.FindStringSubmatch(line)
	if len(matches) < 2 {
		return time.Time{}, "", fmt.Errorf("timestamp não encontrado")
	}
	timestamp, err := time.Parse(timeLayout, matches[1])
	if err != nil {
		return time.Time{}, "", err
	}
	message := strings.TrimSpace(strings.ReplaceAll(line, matches[0], ""))
	return timestamp, message, nil
}

// addToCluster adiciona a mensagem ao minerador e retorna o cluster
func addToCluster(miner *drain3.TemplateMiner, ctx context.Context, message string) (*drain3.LogCluster, error) {
	_, cluster, _, _, err := miner.AddLogMessage(ctx, message)
	return cluster, err
}

// printHourlyCounts exibe os clusters agrupados por hora
func printHourlyCounts(counts map[string]map[int64]int, templates map[int64]string) {
	hours := make([]string, 0, len(counts))
	for hour := range counts {
		hours = append(hours, hour)
	}
	sort.Strings(hours)

	for _, hour := range hours {
		fmt.Printf("Hora: %s:00\n", hour)
		for clusterID, count := range counts[hour] {
			fmt.Printf("  Cluster (%d) %s: %d\n", clusterID, templates[clusterID], count)
		}
	}
}

// printTotalSummary exibe o total de ocorrências por cluster
func printTotalSummary(drain *drain3.Drain) {
	fmt.Println("\n----------------------- TOTAL ------------------------")
	for _, cluster := range drain.GetClusters() {
		fmt.Printf("Cluster (%d) %s: %d\n", cluster.ClusterId, cluster.GetTemplate(), cluster.Size)
	}
}
