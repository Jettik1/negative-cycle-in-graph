package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanLines)

	n, m := readDimensions(scanner)
	if n <= 0 || m < 0 {
		fmt.Println("NO")
		return
	}

	maxVertex, dist, next, valid := readEdges(scanner, m)
	if !valid {
		fmt.Println("NO")
		return
	}

	floydWarshall(dist, next, maxVertex)

	if cycle := findNegativeCycle(dist, next, n); len(cycle) > 0 {
		fmt.Println("YES")
		printCycle(cycle)
	} else {
		fmt.Println("NO")
	}
}

func readDimensions(scanner *bufio.Scanner) (int, int) {
	if !scanner.Scan() {
		return 0, 0
	}
	parts := strings.Fields(scanner.Text())
	if len(parts) != 2 {
		return 0, 0
	}
	n, _ := strconv.Atoi(parts[0])
	m, _ := strconv.Atoi(parts[1])
	return n, m
}

func initMatrices(n int) ([][]int, [][]int) {
	dist := make([][]int, n+1)
	next := make([][]int, n+1)
	for i := 1; i <= n; i++ {
		dist[i] = make([]int, n+1)
		next[i] = make([]int, n+1)
		for j := 1; j <= n; j++ {
			if i == j {
				dist[i][j] = 0
			} else {
				dist[i][j] = math.MaxInt32
			}
			next[i][j] = -1
		}
	}
	return dist, next
}

func readEdges(scanner *bufio.Scanner, m int) (int, [][]int, [][]int, bool) {
	maxVertex := 0
	
	edges := make([][3]int, 0, m)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "" {
			break
		}
		parts := strings.Fields(scanner.Text())
		if len(parts) != 3 {
			fmt.Println("Некорректная строка рёбер, пропускаем.")
			continue
		}

		u, errU := strconv.Atoi(parts[0])
		v, errV := strconv.Atoi(parts[1])
		w, errW := strconv.Atoi(parts[2])

		if errU != nil || errV != nil || errW != nil {
			fmt.Println("Некорректные данные, пропускаем строку.")
			continue
		}

		edges = append(edges, [3]int{u,v,w}) 

		if u > maxVertex {
			maxVertex = u
		}

		if v > maxVertex {
			maxVertex = v
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Ошибка ввода:", err)
		return 0, nil, nil, false
	}

	if len(edges) == 0 {
		return 0, nil, nil, false
	}

	dist, next := initMatrices(maxVertex)

		for _, edge := range edges {
			u, v, w := edge[0], edge[1], edge[2]
			if dist[u][v] > w {
				dist[u][v] = w
				next[u][v] = v
			}
		}
	return maxVertex, dist, next, true
}

func floydWarshall(dist [][]int, next [][]int, n int) {
	for k := 1; k <= n; k++ {
		for i := 1; i <= n; i++ {
			for j := 1; j <= n; j++ {
				if dist[i][k] < math.MaxInt32 && dist[k][j] < math.MaxInt32 {
					if newDist := dist[i][k] + dist[k][j]; newDist < dist[i][j] {
						dist[i][j] = newDist
						next[i][j] = next[i][k]
					}
				}
			}
		}
	}
}

func findNegativeCycle(dist [][]int, next [][]int, n int) []int {
	for i := 1; i <= n; i++ {
		if dist[i][i] < 0 {
			return restoreCycle(next, i)
		}
	}
	return nil
}

func restoreCycle(next [][]int, start int) []int {
	cycle := []int{}
	visited := make([]bool, len(next))
	current := start

	for i := 0; i < len(next); i++ {
		current = next[current][start]
	}

	for !visited[current] {
		visited[current] = true
		cycle = append(cycle, current)
		current = next[current][start]
	}

	cycle = append(cycle, current)
	return cycle
}

func printCycle(cycle []int) {
	for i, v := range cycle {
		if i > 0 {
			fmt.Print(" ")
		}
		fmt.Print(v)
	}
	fmt.Println()
}
