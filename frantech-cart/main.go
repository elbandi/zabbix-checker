package main

import (
	"errors"
	"fmt"
	"github.com/Elbandi/zabbix-checker/common/lld"
	"github.com/antchfx/htmlquery"
	"github.com/urfave/cli/v2"
	"golang.org/x/net/html"
	"golang.org/x/net/html/charset"
	"log"
	"net/http"
	"os"
	"strings"
)

func fetchPage(url string) (*html.Node, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "qalandar-fetcher/1.0")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	r, err := charset.NewReader(resp.Body, resp.Header.Get("Content-Type"))
	if err != nil {
		return nil, err
	}
	return html.Parse(r)
}

func getAvailable(n *html.Node) string {
	if n == nil {
		return "99"
	}
	s := htmlquery.InnerText(n)
	s = strings.TrimSpace(strings.ReplaceAll(s, "Available", ""))
	return s
}

func init() {
	cli.VersionFlag = &cli.BoolFlag{
		Name:    "print-version",
		Aliases: []string{"V"},
		Usage:   "print only the version",
	}
}

func main() {
	app := &cli.App{
		Name:        "frantech-cart-fetcher",
		Version:     "v1.0",
		Usage:       "fetch frantech cart availability data",
		Description: "",
		Authors: []*cli.Author{
			{
				Name:  "Elbandi",
				Email: "elso.andras@gmail.com",
			},
		},
		Action: actFetch,
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func actFetch(context *cli.Context) error {
	first := context.Args().First()
	if first == "" {
		return cli.Exit("first argument is required", 2)
	}
	doc, err := fetchPage("https://my.frantech.ca/cart.php?gid=" + first)
	if err != nil {
		return err
	}
	list, err := htmlquery.QueryAll(doc, "//div[@id='products']//div[contains(@class, 'package') and starts-with(@id , 'product')]")
	if err != nil {
		return err
	}
	if len(list) < 1 {
		return errors.New("Error fetch produts")
	}
	d := make(lld.DiscoveryData, 0)
	for _, pack := range list {
		//		log.Println("-----")
		//		log.Println(htmlquery.OutputHTML(pack, true))
		nameNode, err := htmlquery.Query(pack, ".//h3")
		if err != nil || nameNode == nil {
			continue
		}
		name := htmlquery.InnerText(nameNode)
		availNode, err := htmlquery.Query(pack, ".//div[contains(@class, 'package-qty')]")
		if err != nil {
			continue
		}
		item := make(lld.DiscoveryItem, 0)
		item["ID"] = name
		item["AVAILABILITY"] = getAvailable(availNode)
		d = append(d, item)
	}
	fmt.Println(d.JsonLine())
	return nil
}
