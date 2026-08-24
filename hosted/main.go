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

	tree.Mount("/sys/mem", newNode("metrics"))
	tree.nodes["/sys/mem"].data = "used=40%\n"

	tree.Mount("/sys/cpu", newNode("metrics"))
	tree.nodes["/sys/cpu"].data = "load=0.1\n"

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

// tickers keeps metrics/net resources alive so watchers see a living system.
func tickers(t *Tree) {
	mem := t.nodes["/sys/mem"]
	cpu := t.nodes["/sys/cpu"]
	net := t.nodes["/dev/net0"]
	i := 0
	for range time.Tick(1 * time.Second) {
		i++
		mem.Write([]byte(fmt.Sprintf("used=%d%%\n", 40+i%30)))
		cpu.Write([]byte(fmt.Sprintf("load=%.1f\n", float64(i%10)/10.0)))
		net.Write([]byte(fmt.Sprintf("frame seq=%d src=tcp!node%d!9pp\n", i, i%3)))
	}
}
