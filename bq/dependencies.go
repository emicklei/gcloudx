package bq

import (
	"context"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"cloud.google.com/go/bigquery"
	"github.com/emicklei/dot"
)

type BigQueryArguments struct {
	Verbose      bool
	TableSources []string
	Output       string
}

func ExportViewDepencyGraph(args BigQueryArguments) error {
	g := dot.NewGraph(dot.Directed)
	visited := map[string]bool{}
	for _, each := range args.TableSources {
		p, d, v := tokenize(each)
		if err := addDependencies(p, d, v, g, visited); err != nil {
			return err
		}
	}
	return os.WriteFile(args.Output, []byte(g.String()), os.ModePerm)
}

func addDependencies(project string, dataset string, table string, root *dot.Graph, visited map[string]bool) error {
	key := fmt.Sprintf("%s.%s.%s", project, dataset, table)
	if _, ok := visited[key]; ok {
		return nil
	}
	log.Println("reading definition of ", key)
	visited[key] = true
	fromNode, _ := ensureTableNode(root, project, dataset, table)
	ctx := context.Background()
	client, err := bigquery.NewClient(ctx, project)
	if err != nil {
		log.Println(err)
		return err
	}
	defer client.Close()
	ds := client.Dataset(dataset)
	vw := ds.Table(table)
	m, err := vw.Metadata(ctx)
	if err != nil {
		log.Println(err)
		return err
	}
	if len(m.ViewQuery) == 0 { // table not view
		return nil
	}
	// recurse
	for _, each := range deps(m.ViewQuery) {
		p, d, t := tokenize(each)
		toNode, _ := ensureTableNode(root, p, d, t)
		fromNode.Edge(toNode)
		addDependencies(p, d, t, root, visited)
	}
	return nil
}

func ensureTableNode(root *dot.Graph, project string, dataset string, table string) (dot.Node, bool) {
	pg, ok := root.FindSubgraph(project)
	if !ok {
		pg = root.Subgraph(project, dot.ClusterOption{})
		pg.Label(project)
		dg := pg.Subgraph(dataset)
		n := dg.Node(table)
		// modify label
		n.Label(wrap(table))
		return n, true
	}
	dg, ok := pg.FindSubgraph(dataset)
	if !ok {
		dg = pg.Subgraph(dataset, dot.ClusterOption{})
		dg.Label(dataset)
		n := dg.Node(table)
		// modify label
		n.Label(wrap(table))
		return n, true
	}
	tn, ok := dg.FindNodeById(table)
	if !ok {
		n := dg.Node(table)
		// modify label
		n.Label(wrap(table))
		return n, true
	}
	return tn, false
}

func wrap(s string) string {
	if len(s) < 24 {
		return s
	}
	return s[0:len(s)/2] + "\n" + s[len(s)/2:]
}

// deduplicated
func deps(vw string) []string {
	r := regexp.MustCompile("`.+`")
	quoted := r.FindAllString(vw, -1)
	for i, each := range quoted {
		quoted[i] = strings.Trim(each, "`")
	}
	includes := func(list []string, element string) bool {
		for _, each := range list {
			if each == element {
				return true
			}
		}
		return false
	}
	dedup := []string{}
	for _, each := range quoted {
		if !includes(dedup, each) {
			dedup = append(dedup, each)
		}
	}
	return dedup
}

func tokenize(t string) (string, string, string) {
	uncolon := strings.ReplaceAll(t, ":", ".")
	parts := strings.Split(uncolon, ".")
	if len(parts) != 3 {
		log.Fatal("failed to parse a full qualified table|view name, expected PROJECT(.|:)DATASET.VIEW:", t)
	}
	return parts[0], parts[1], parts[2]
}
