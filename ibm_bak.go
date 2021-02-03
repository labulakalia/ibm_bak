package main

import (
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
    "os"
    "strings"
    "time"

    md "github.com/JohannesKaufmann/html-to-markdown"
    "github.com/PuerkitoBio/goquery"
)

func savePage(url string) error {
    //url := "https://developer.ibm.com/zh/articles/db-redis-from-basic-to-advanced-01/"

    resp, err := http.Get(url)
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    doc, err := goquery.NewDocumentFromReader(resp.Body)
    if err != nil {
        return err
    }
    title := doc.Find("#content > div.page_content > div > div.code_postcontent > div.cpt-hero > div > div > div > div > h1")

    title2 := doc.Find("#content > div.page_content > div > div.code_postcontent > div.cpt-hero > div > div > div > div > h2")

    content := doc.Find("#content > div.page_content > div > div.code_postcontent > div.cpt-body > div > div > div.bx--col-md-7.bx--col-lg-7")

    tag := doc.Find("#pane1")

    converter := md.NewConverter("", true, nil)
    converter.Before(func(selec *goquery.Selection) {
        selec.Find("img").Each(func(i int, s *goquery.Selection) {
            _, ok := s.Attr("src")
            if ok {
                return
            }
            _, err = os.Stat("ibm_articles_img")
            if os.IsNotExist(err) {
                err = os.Mkdir("ibm_articles_img", 0755)
                if err != nil {
                    log.Println(err)
                    return
                }
            }
            src, ok := s.Attr("data-src")
            if !ok {
                return
            }

            resp, err := http.Get(src)
            if err != nil {
                log.Println(err)
                return
            }
            defer resp.Body.Close()
            img_path := "ibm_articles_img/" + strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(src, "https://developer.ibm.com/developer/default/articles/", ""), "/", "_"), "images/", "")
            // os.OpenFile(img_path, os, perm os.FileMode)
    
            // "https://developer.ibm.com/developer/default/articles/images/policy-based-governance-security.png"
            data, err := ioutil.ReadAll(resp.Body)
            if err != nil {
                log.Println(err)
                return
            }
            ioutil.WriteFile(img_path, data, 0644)
            s.SetAttr("src", "../"+img_path)
        })
    })
    tags := []string{}
    tag.Find("a").Each(func(i int, s *goquery.Selection) {
        tags = append(tags, s.Text())
    })
    mdcontent := fmt.Sprintf("# %s\n%s\n\n**标签:** %s\n\n[原文链接](%s)\n\n", title.Text(), title2.Text(), strings.Join(tags, ","), url)
    mdcontent += converter.Convert(content)

    _, err = os.Stat("ibm_articles")
    if os.IsNotExist(err) {
        err = os.Mkdir("ibm_articles", 0755)
        if err != nil {
            return err
        }
    }
    t := strings.ReplaceAll(title.Text(), "/", "_")
    t = strings.ReplaceAll(t, " ", "")
    return ioutil.WriteFile(fmt.Sprintf("ibm_articles/%s.md", t), []byte(mdcontent), 0644)
}

func SavePages(urls chan string) {
    for url := range urls {
        log.Printf("start save page url:%s\n", url)
        err := savePage(url)
        if err != nil {
            log.Printf("save url:%s failed: %v\n", url, err)
        }
        log.Println("success save")
    }
}

func main() {
    var urls = make(chan string)
    pageurl := `https://developer.ibm.com/zh/articles/page/%d/?fa=date%3ADESC&fb`
    go SavePages(urls)
    for i := 2; i <= 73; i++ {
        var url = fmt.Sprintf(pageurl, i)
        resp, err := http.Get(url)
        if err != nil {
            log.Fatalln(err)
        }
        if resp.StatusCode == 301 {
            url = resp.Header.Get("location")

            resp, err = http.Get(url)
            if err != nil {
                log.Fatalln(err)
            }
        }
        doc, err := goquery.NewDocumentFromReader(resp.Body)
        if err != nil {
            log.Fatalln(err)
        }
        sec := doc.Find("#content > div > div > div.cpt-archive-new > div.search-results > div > div:nth-child(2) > div.bx--row.bx--no-gutter--right.code-card-grid > div")
        //fmt.Println(sec.Html())
        sec.Each(func(i int, s *goquery.Selection) {

            url, ok := s.Find("div > a").Attr("href")
            if !ok {
                return
            }

            urls <- "https://developer.ibm.com" + url
        })
    }
    time.Sleep(time.Second * 10)

}

