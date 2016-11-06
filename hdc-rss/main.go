package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"time"

	"github.com/tominescu/hdc-rss/config"
	"github.com/tominescu/hdc-rss/rss"
)

const LINK_RETRY_TIMES = 3

func main() {
	c := loadConfig()
	timeout := c.Timeout
	if timeout <= 0 || timeout > 600 {
		timeout = 30
	}

	httpClient := http.Client{
		Timeout: time.Duration(timeout) * time.Second,
	}

	finish := make(chan int, 1)
	for _, site := range c.Sites {
		go handleSite(site, &httpClient, finish)
	}

	// Wait for finish
	for _ = range c.Sites {
		<-finish
	}
	log.Println("Exit because all goroutines are down")
}

func handleSite(site config.Site, httpClient *http.Client, finish chan<- int) {
	logPrefix := fmt.Sprintf("[%s]\t", site.Name)
	log := log.New(os.Stdout, logPrefix, log.LstdFlags)
	log.Printf("Start with url %s", site.Rss)

	// create download dir
	if err := os.MkdirAll(site.DownloadDir, 0755); err != nil {
		log.Printf("Error create download directory: %s", err)
		finish <- 1
	}

	linkCountMap := make(map[string]int, 10)
	firstTurn := true
	for {
		if firstTurn {
			firstTurn = false
		} else {
			time.Sleep(time.Duration(site.Interval) * time.Second)
		}

		rsp, err := httpClient.Get(site.Rss)
		if err != nil {
			log.Printf("Error fetching rss:%s", err)
			continue
		}
		data, err := ioutil.ReadAll(rsp.Body)
		rsp.Body.Close()
		if err != nil {
			log.Printf("Error read response", err)
			continue
		}
		r, err := rss.NewRss(data)
		if err != nil {
			log.Printf("Error parse rss")
			continue
		}
		links := r.GetLinks()
		var newLinks []string

		for _, link := range links {
			if linkCountMap[link] < LINK_RETRY_TIMES {
				newLinks = append(newLinks, link)
			}
		}
		log.Printf("Get %d links include %d new links", len(links), len(newLinks))
		for _, link := range newLinks {
			if err = handleLink(&site, link, log, httpClient); err != nil {
				linkCountMap[link]++
			} else {
				linkCountMap[link] = LINK_RETRY_TIMES
			}
		}
	}
}

func handleLink(site *config.Site, link string, log *log.Logger, httpClient *http.Client) error {
	log.Printf("Handled link %s", link)
	rsp, err := httpClient.Get(link)
	if err != nil {
		log.Printf("Error download link %s : %s", link, err)
		return err
	}
	defer rsp.Body.Close()
	data, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		log.Printf("Error read link response %s : %s", link, err)
		return err
	}

	tmpVal := rsp.Header.Get("Content-Disposition")
	value, err := url.QueryUnescape(tmpVal)
	if err != nil {
		value = tmpVal
	}
	filename := getFileName(value)

	log.Printf("Get file %s", filename)
	filePath := path.Join(site.DownloadDir, filename)

	err = ioutil.WriteFile(filePath, data, 0644)
	if err != nil {
		log.Printf("Save File %s failed: %s", filename, err)
		return err
	}
	log.Printf("Torrent save to %s", filePath)

	return nil
}

func getFileName(value string) string {
	for _, str := range strings.Split(value, ";") {
		str = strings.TrimSpace(str)
		if strings.HasPrefix(str, "filename=") {
			return strings.TrimPrefix(str, "filename=")
		}
	}
	return ""
}

func loadConfig() *config.Config {
	configPath := flag.String("c", "config.json", "The config file path")
	help := flag.Bool("h", false, "Print this usage")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [-c config.json] [-h]\n\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()
	if *help {
		flag.Usage()
		os.Exit(0)
	}

	configData, err := ioutil.ReadFile(*configPath)
	if err != nil {
		loadConfigErr(err)
	}

	c, err := config.NewConfig(configData)
	if err != nil {
		loadConfigErr(err)
	}
	return c
}

func loadConfigErr(err error) {
	fmt.Fprintf(os.Stderr, "Error Loading config: %s\n\n", err)
	flag.Usage()
	os.Exit(1)
}
