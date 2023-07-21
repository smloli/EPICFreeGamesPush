package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"loli/wxPusher"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type AppliedRulesInfo struct {
	EndDate string
}

type LineOffersInfo struct {
	AppliedRules []AppliedRulesInfo
}

type KeyImagesInfo struct {
	Type string
	Url  string
}

type GamesInfo struct {
	Title       string
	Description string
	ProductSlug string
	Price       struct {
		TotalPrice struct {
			FmtPrice struct {
				DiscountPrice string
			}
		}
		LineOffers []LineOffersInfo
	}
	KeyImages []KeyImagesInfo
}

type Games struct {
	Data struct {
		Catalog struct {
			SearchStore struct {
				Elements []GamesInfo
			}
		}
	}
}

type Config struct {
	AppToken string   `json:"appToken"`
	TopicIds []int    `json:"topicIds"`
	Uids     []string `json:"uids"`
}

func get(url string) ([]byte, error) {
	c := http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/115.0.0.0 Safari/537.36")
	req.Header.Set("Referer", "https://store.epicgames.com/zh-CN/")
	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, _ := io.ReadAll(resp.Body)
	return data, nil
}

func (c *Config) read() error {
	path := filepath.Join(filepath.Dir(os.Args[0]), "config.json")
	if _, openErr := os.Stat(path); openErr != nil {
		c.TopicIds = make([]int, 1)
		c.Uids = make([]string, 1)
		data, _ := json.MarshalIndent(c, "", "    ")
		f, createErr := os.Create(path)
		if createErr != nil {
			log.Println(createErr)
			return createErr
		}
		defer f.Close()
		f.Write(data)
		log.Println("请在config.json中填写配置参数")
		return openErr
	}
	f, err := os.Open(path)
	if err != nil {
		log.Println(err)
		return err
	}
	defer f.Close()
	data, _ := io.ReadAll(f)
	json.Unmarshal(data, c)
	return nil
}

func main() {
	var g Games
	var c Config
	var text string
	if err := c.read(); err != nil {
		return
	}
	url := "https://store-site-backend-static-ipv4.ak.epicgames.com/freeGamesPromotions?locale=zh-CN&country=CN&allowCountries=CN"
	data, err := get(url)
	if err != nil {
		panic(err)
	}
	json.Unmarshal(data, &g)
	for _, v := range g.Data.Catalog.SearchStore.Elements {
		if len(v.Price.LineOffers[0].AppliedRules) == 0 || v.Price.TotalPrice.FmtPrice.DiscountPrice != "0" {
			continue
		}
		var imageUrl string
		for _, k := range v.KeyImages {
			if k.Type == "Thumbnail" {
				imageUrl = k.Url
				break
			}
		}
		t, _ := time.Parse(time.RFC3339, v.Price.LineOffers[0].AppliedRules[0].EndDate)
		text += fmt.Sprintf("![](%s)\n**<center>%s</center>**\n>##### 游戏简介\n> %s\n>##### 结束时间\n> %s\n>##### 领取地址\n> https://store.epicgames.com/zh-CN/p/%s\n\n", imageUrl, v.Title, v.Description, t.Format(time.DateTime), v.ProductSlug)
	}
	resp, err := wxPusher.Send(c.AppToken, text, "EPIC本周免费游戏推送", 3, c.TopicIds, c.Uids)
	if err != nil {
		return
	}
	fmt.Println(string(resp))
}
