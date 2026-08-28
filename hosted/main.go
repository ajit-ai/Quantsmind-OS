package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

func main() {
	tree := NewTree()

	tree.Mount("/srv/status", newNode(""))
	tree.nodes["/srv/status"].data = "quantis hosted proto up\n"

	tree.Mount("/srv/version", newNode(""))
	tree.nodes["/srv/version"].data = Product + " " + Version + "\n"

	tree.Mount("/sys/mem", newNode("metrics"))
	tree.Mount("/sys/cpu", newNode("metrics"))

	tree.Mount("/dev/net0", newNode("net"))
	tree.nodes["/dev/net0"].data = "link up\n"

	tree.Mount("/ai/models/qm-demo/infer", newNode("infer"))
	tree.nodes["/ai/models/qm-demo/infer"].data =
		"POST a prompt here; tokens stream on ./watch\n"

	mux := http.NewServeMux()
	mux.HandleFunc("/ns", tree.HandleNamespace)
	mux.HandleFunc("/r/", func(w http.ResponseWriter, r *http.Request) {
		p := strings.TrimPrefix(r.URL.Path, "/r/")
		if !strings.HasPrefix(p, "/") {
			p = "/" + p
		}
		if strings.HasSuffix(p, "/watch") {
			tree.HandleWatch(w, r, strings.TrimSuffix(p, "/watch"))
			return
		}
		switch r.Method {
		case http.MethodGet:
			tree.HandleRead(w, p)
		case http.MethodPost:
			tree.HandleWrite(w, r, p)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	go tickers(tree)

	addr := ":9090"
	log.Printf("quantisd listening on %s  (namespace: http://localhost%s/ns)", addr, addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}

// tickers feeds live system stats (real on Windows/Linux, simulated elsewhere)
// plus synthetic network frames, so watchers see a living system.
func tickers(t *Tree) {
	mem := t.nodes["/sys/mem"]
	cpu := t.nodes["/sys/cpu"]
	net := t.nodes["/dev/net0"]
	i := 0
	simCPU, simMem := -1, -1
	for range time.Tick(1 * time.Second) {
		i++
		c, m, src := sampleSystem()
		if c < 0 {
			simCPU = (simCPU + 3) % 71
			if simCPU < 0 {
				simCPU = 0
			}
			c = simCPU
			src = "simulated"
		}
		if m < 0 {
			simMem = (i * 7) % 60
			m = simMem
		}
		mem.Write([]byte(fmt.Sprintf("used=%d%% (%s)\n", m, src)))
		cpu.Write([]byte(fmt.Sprintf("load=%d%% (%s)\n", c, src)))
		net.Write([]byte(fmt.Sprintf("frame seq=%d src=tcp!node%d!9pp\n", i, i%3)))
	}
}
