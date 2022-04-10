package controller

import (
	"fmt"
	"github.com/1340691923/xwl_bi/platform-basic-libs/util"
	"github.com/PuerkitoBio/goquery"
	"github.com/valyala/fasthttp"
	"strings"
)

type Word struct {
	Danci        string   `json:"danci"`
	ChineseFanyi []string `json:"chinese_fanyi"`
	YingYingBiao string   `json:"ying_ying_biao"`
	MeiYingBiao  string   `json:"mei_ying_biao"`
	Additional   string   `json:"additional"`
	WebTrans     struct {
		TWebTrans map[string]string `json:"t_web_trans"`
		WebPhrase map[string]string `json:"web_phrase"`
		TPETrans  map[string][]struct {
			Title      string `json:"title"`
			Additional string `json:"additional"`
			Source     string `json:"source"`
			Trans      string `json:"trans"`
		} `json:"tpe_trans"`
		TEETrans map[string][]struct {
			Def      string `json:"def"`
			Juzi     string `json:"juzi"`
			TongYiCi string `json:"tong_yi_ci"`
		} `json:"tee_trans"`
		TEETransName map[string]string `json:"tee_trans_name"`
	} `json:"web_trans"`
	WordGroup    map[string]string `json:"word_group"`
	JinYiCi      map[string]string `json:"jin_yi_ci"`
	TongGengCi   map[string]string `json:"tong_geng_ci"`
	CiYuBianXi   map[string]string `json:"ci_yu_bian_xi"`
	ShuangYuLiJu map[string]string `json:"shuang_yu_li_ju"`
	YingAudio    string            `json:"ying_audio"`
	MeiAudio     string            `json:"mei_audio"`
}

func GetWordParse(ctx *fasthttp.RequestCtx) {

	if string(ctx.FormValue("word")) == "" {
		util.WriteJSON(ctx, map[string]interface{}{
			"code": 500,
			"msg":  "请点击有效单词",
		})
		return
	}

	word := new(Word)
	word.Danci = string(ctx.FormValue("word"))
	uri := "https://dict.youdao.com/search?q=" + word.Danci + "&keyfrom=new-fanyi.smartResult"

	doc, err := goquery.NewDocument(uri)
	if err != nil {
		panic(err)
	}

	doc.Find("#results-contents #phrsListTab .trans-container ul li").Each(func(i int, selection *goquery.Selection) {
		word.ChineseFanyi = append(word.ChineseFanyi, selection.Text())
	})

	word.Additional = strings.Replace(doc.Find("#results-contents #phrsListTab .additional").Text(), "\n", "", -1)

	doc.Find(".baav .pronounce").Each(func(i int, selection *goquery.Selection) {
		tmp := selection.Find(".phonetic").Text()
		if i == 0 {
			word.YingYingBiao = tmp
			word.YingAudio = fmt.Sprintf("http://dict.youdao.com/dictvoice?audio=%s&type=1", word.Danci)
		} else if i == 1 {
			word.MeiYingBiao = tmp
			word.MeiAudio = fmt.Sprintf("http://dict.youdao.com/dictvoice?audio=%s&type=2", word.Danci)
		}
	})

	word.WebTrans.TWebTrans = map[string]string{}
	doc.Find("#webTrans .wt-container").Each(func(i int, selection *goquery.Selection) {
		key := selection.Find(".title ").Text()
		val, _ := selection.Find(".collapse-content").Eq(0).Html()
		word.WebTrans.TWebTrans[strings.Replace(key, "\n", "", -1)] = strings.Replace(val, "\n", "", -1)
	})
	word.WebTrans.WebPhrase = map[string]string{}
	doc.Find("#webPhrase .wordGroup").Each(func(i int, selection *goquery.Selection) {
		key := selection.Find(".contentTitle").Text()
		selection.Find(".contentTitle").Empty()

		val := strings.Replace(selection.Text(), " ", "", -1)
		// 去除换行符
		val = strings.Replace(val, "\n", "", -1)

		word.WebTrans.WebPhrase[key] = val

	})

	word.WebTrans.TPETrans = map[string][]struct {
		Title      string `json:"title"`
		Additional string `json:"additional"`
		Source     string `json:"source"`
		Trans      string `json:"trans"`
	}{}
	doc.Find("#tPETrans-type-list a").Each(func(i int, selection *goquery.Selection) {

		rel, found := selection.Attr("rel")
		if found {
			tmp := doc.Find("#tPETrans-all-trans ." + rel)
			tmp.Find(".items").Each(func(i int, trans *goquery.Selection) {
				word.WebTrans.TPETrans[selection.Text()] = append(word.WebTrans.TPETrans[selection.Text()], struct {
					Title      string `json:"title"`
					Additional string `json:"additional"`
					Source     string `json:"source"`
					Trans      string `json:"trans"`
				}(struct {
					Title      string
					Additional string
					Source     string
					Trans      string
				}{
					Title:      trans.Find(".title").Text(),
					Additional: trans.Find(".additional").Text(),
					Source:     trans.Find(".source").Text(),
					Trans:      trans.Find(".trans").Text(),
				}))
			})

		}
	})

	TEETransNameDoc := doc.Find("#tEETrans h4")

	TEETransNameVal := TEETransNameDoc.Find(".phonetic").Text()
	TEETransNameDoc.Find(".phonetic").Empty()
	TEETransNameKey := TEETransNameDoc.Text()

	TEETransNameKey = strings.Replace(TEETransNameKey, " ", "", -1)
	// 去除换行符
	TEETransNameKey = strings.Replace(TEETransNameKey, "\n", "", -1)

	word.WebTrans.TEETransName = map[string]string{TEETransNameKey: TEETransNameVal}

	word.WebTrans.TEETrans = map[string][]struct {
		Def      string `json:"def"`
		Juzi     string `json:"juzi"`
		TongYiCi string `json:"tong_yi_ci"`
	}{}
	doc.Find("#tEETrans li").Each(func(i int, selection *goquery.Selection) {

		pos := selection.Find(".pos").Text()
		word.WebTrans.TEETrans[pos] = nil

		if selection.Find(".ol").Text() != "" {
			selection.Find(".ol li").Each(func(i int, selection *goquery.Selection) {

				word.WebTrans.TEETrans[pos] = append(word.WebTrans.TEETrans[pos],
					struct {
						Def      string `json:"def"`
						Juzi     string `json:"juzi"`
						TongYiCi string `json:"tong_yi_ci"`
					}(struct {
						Def      string `json:"def"`
						Juzi     string `json:"juzi"`
						TongYiCi string `json:"tong_yi_ci"`
					}{
						Def:      strings.ReplaceAll(selection.Find(".def").Text(), `"`, ""),
						Juzi:     strings.ReplaceAll(selection.Find("p em").Text(), `"`, ""),
						TongYiCi: selection.Find(".gray a").Text(),
					}))
			})
		} else {
			word.WebTrans.TEETrans[pos] = append(word.WebTrans.TEETrans[pos],
				struct {
					Def      string `json:"def"`
					Juzi     string `json:"juzi"`
					TongYiCi string `json:"tong_yi_ci"`
				}{
					Def:      strings.ReplaceAll(selection.Find(".def").Text(), `"`, ""),
					Juzi:     strings.ReplaceAll(selection.Find("p em").Text(), `"`, ""),
					TongYiCi: selection.Find(".gray a").Text(),
				})
		}
	})

	word.WordGroup = map[string]string{}

	doc.Find("#wordGroup .wordGroup").Each(func(i int, selection *goquery.Selection) {
		a := selection.Find("a")
		key := a.Text()
		a.Empty()

		val := strings.Replace(selection.Text(), " ", "", -1)
		// 去除换行符
		val = strings.Replace(val, "\n", "", -1)

		word.WordGroup[key] = val
	})
	tmp := []string{}
	word.JinYiCi = map[string]string{}
	doc.Find("#synonyms ul li").Each(func(i int, selection *goquery.Selection) {
		word.JinYiCi[selection.Text()] = ""
		tmp = append(tmp, selection.Text())
	})

	doc.Find("#synonyms ul .wordGroup").Each(func(i int, selection *goquery.Selection) {
		list := []string{}
		selection.Find("a").Each(func(i int, selection *goquery.Selection) {
			list = append(list, selection.Text())
		})

		word.JinYiCi[tmp[i]] = strings.Join(list, ";")
	})

	word.CiYuBianXi = map[string]string{}
	doc.Find("#discriminate .collapse-content .wordGroup p").Each(func(i int, selection *goquery.Selection) {
		contentTitle := selection.Find(".contentTitle")

		key := contentTitle.Text()

		key = strings.Replace(key, " ", "", -1)
		// 去除换行符
		key = strings.Replace(key, "\n", "", -1)

		contentTitle.Empty()

		val := strings.Replace(selection.Text(), " ", "", -1)
		// 去除换行符
		val = strings.Replace(val, "\n", "", -1)
		word.CiYuBianXi[key] = val
	})
	word.TongGengCi = map[string]string{}
	doc.Find("#relWordTab .wordGroup").Each(func(i int, selection *goquery.Selection) {
		contentTitle := selection.Find(".contentTitle")

		key := contentTitle.Text()

		key = strings.Replace(key, " ", "", -1)
		// 去除换行符
		key = strings.Replace(key, "\n", "", -1)

		contentTitle.Empty()

		val := strings.Replace(selection.Text(), " ", "", -1)
		// 去除换行符
		val = strings.Replace(val, "\n", "", -1)
		word.TongGengCi[key] = val
	})
	word.ShuangYuLiJu = map[string]string{}
	doc.Find("#bilingual li").Each(func(i int, selection *goquery.Selection) {
		list := []string{}
		selection.Find("p").RemoveClass("example-via").Each(func(i int, selection *goquery.Selection) {

			val := strings.Replace(selection.Text(), "\n", "", -1)
			list = append(list, val)
		})
		if len(list) > 2 {
			word.ShuangYuLiJu[list[0]] = list[1]
		}
	})

	util.WriteJSON(ctx, map[string]interface{}{
		"code": 0,
		"msg":  "ok",
		"data": word,
	})
}
