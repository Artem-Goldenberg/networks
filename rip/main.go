package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"text/tabwriter"
	"time"
)

var maxIP = 256

type Pair struct {
	srcIP, dstIP string
}

type Value struct {
	next     string
	distance int
}

type RoutingTable map[Pair]Value

type Node struct {
	ip          string
	lock        *sync.Mutex
	connections map[string]Node
	// key not exists ==> distance is infinity
	table RoutingTable
}

func createEdge(g []Node, i, j int) {
	_, ok := g[i].connections[g[j].ip]
	if ok {
		// already has this edge, check synchronization
		_, ok = g[j].connections[g[i].ip]
		assert(ok)
		return
	}
	// add connection on both sides
	g[i].connections[g[j].ip] = g[j]
	g[j].connections[g[i].ip] = g[i]

	// distance between neighbours is 1
	g[i].table[Pair{g[i].ip, g[j].ip}] = Value{g[j].ip, 1}
	g[j].table[Pair{g[j].ip, g[i].ip}] = Value{g[i].ip, 1}
}

func generateGraph(n int, edgeProbability float64) []Node {
	g := make([]Node, n)
	for i := 0; i < n; i++ {
		// chance to repeat ~ C(n, 2) * 1/256^4, n = 100 (>max) ==> chance = 0,0000011525
		// a.b.c.d
		a := rand.Intn(maxIP)
		b := rand.Intn(maxIP)
		c := rand.Intn(maxIP)
		d := rand.Intn(maxIP)
		g[i].ip = fmt.Sprintf("%d.%d.%d.%d", a, b, c, d)
		g[i].lock = &sync.Mutex{}
		g[i].connections = make(map[string]Node)
		g[i].table = make(RoutingTable)
		// distance to itself is 0
		// g[i].table[Pair{g[i].ip, g[i].ip}] = Value{g[i].ip, 0}
	}
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			if rand.Float64() < edgeProbability {
				createEdge(g, i, j)
			}
		}
	}
	return g
}

func printGraph(g []Node) {
	for _, node := range g {
		fmt.Printf("*: %s, connected to %d nodes:\n", node.ip, len(node.connections))
		for _, conn := range node.connections {
			fmt.Printf("%s, ", conn.ip)
		}
		fmt.Printf("\n\n")
	}
}

var wg sync.WaitGroup

func (n *Node) receiveUpdateFrom(ip string) {
	wg.Add(1)
	n.lock.Lock()
	locked := true
	// distance to the neighbor, should always be one
	d := n.table[Pair{n.ip, ip}].distance
	assert(d == 1)
	updated := false
	for k, v := range n.connections[ip].table {
		if k.dstIP == n.ip {
			continue
		}
		key := Pair{n.ip, k.dstIP}
		current, ok := n.table[key]
		if ok && current.distance <= d+v.distance {
			continue
		}
		n.table[key] = Value{next: ip, distance: d + v.distance}
		updated = true
	}
	if updated {
		fmt.Printf("State of node %s has changed:\n", n.ip)
		n.print()
		fmt.Println()
		n.lock.Unlock()
		locked = false
		for _, conn := range n.connections {
			// asynchronous work
			go conn.receiveUpdateFrom(n.ip)
		}
	}
	if locked { 
		n.lock.Unlock()
	}
	wg.Done()
}

func (n *Node) print() {
	writer := tabwriter.NewWriter(os.Stdout, 0, 4, 3, ' ', 0)
	fmt.Fprintln(writer, "Source IP\tDestination IP\tNext Hop\tDistance")
	for k, v := range n.table {
		fmt.Fprintf(writer, "%s\t%s\t%s\t      %d\n", k.srcIP, k.dstIP, v.next, v.distance)
	}
	writer.Flush()
}

func main() {
	n, err := strconv.Atoi(os.Args[1])
	assert(err == nil)
	p, err := strconv.ParseFloat(os.Args[2], 64)
	assert(err == nil)

	g := generateGraph(n, p)
	fmt.Println("Net configuration:")
	printGraph(g)
	fmt.Println()

	initialPrint(g)

	fmt.Printf("---------- Starting RIP -----------\n\n")
	start(g)
	time.Sleep(time.Second)
	wg.Wait()
	fmt.Printf("---------- Finished RIP -----------\n\n")
	finalPrint(g)
}

// to start the algorithm everyone will tell everyone that their's table has changed
func start(g []Node) {
	for _, node := range g {
		for _, conn := range node.connections {
			go conn.receiveUpdateFrom(node.ip)
		}
	}
}

func initialPrint(g []Node) {
	for _, node := range g {
		fmt.Printf("Initial state of node %s\n", node.ip)
		node.print()
		fmt.Println()
	}
}

func finalPrint(g []Node) {
	for _, node := range g {
		fmt.Printf("Final state of node %s\n", node.ip)
		node.print()
		fmt.Println()
	}
}

func assert(cond bool) {
	if !cond {
		panic("Assert failed")
	}
}
