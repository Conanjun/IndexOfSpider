package main

import (
	"flag"
	"fmt"
	"github.com/antchfx/htmlquery"
	"github.com/gocolly/colly"
	"net/url"
	"os"
	"strings"
	"time"
)


var indefOfUrl string
var help bool

func main() {
	flag.StringVar(&indefOfUrl,"indexofurl", "", "index of url")

	flag.Parse()

	if help{
		flag.Usage()
		return
	}

	if indefOfUrl==""{
		flag.Usage()
		return
	}

	targetIndexOfUrl:= indefOfUrl

	u,err:=url.Parse(targetIndexOfUrl)
	if err!=nil{
		fmt.Println(err)
	}
	host:=u.Host
	fmt.Println(host)
	err=os.MkdirAll(host,os.ModePerm)
	if err!=nil{
		fmt.Println(err)
	}

	c := colly.NewCollector(
		colly.AllowedDomains(host),
		colly.MaxDepth(1),
	)


	c.OnRequest(func(r *colly.Request){
		fmt.Println("Visiting",r.URL.String())
	})

	c.Limit(&colly.LimitRule{
		RandomDelay: 1*time.Second,
	})


	c.OnResponse(func(r *colly.Response) {
		doc,err:=htmlquery.Parse(strings.NewReader(string(r.Body)))
		if err!=nil{
			fmt.Println(err)
		}
		nodes:=htmlquery.Find(doc,`/html/body/table/tbody/tr/td/a`)

		//如果nodes长度为0则保存
		if len(nodes)==0{
			absurl, _ := url.Parse(r.Request.URL.String())
			abspath:=host+absurl.Path
			r.Save(abspath)
		}

		for _,node:=range nodes{
			link:=htmlquery.InnerText(node)
			if link=="Parent Directory"{
				continue
			}
			if strings.HasSuffix(link,"/"){
				absurl, _ :=url.Parse(r.Request.AbsoluteURL(link))
				abspath:=host+absurl.Path
				fmt.Println(abspath)
				err=os.MkdirAll(abspath,os.ModePerm)
				if err!=nil{
					fmt.Println(err)
				}
				c.Visit(r.Request.AbsoluteURL(link))
			} else {
				c.Visit(r.Request.AbsoluteURL(link))
			}
		}

	})


	err=c.Visit(targetIndexOfUrl)
	if err!=nil{
		fmt.Println(err)
	}
}
