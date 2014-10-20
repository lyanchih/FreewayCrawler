package freeway

import (
  "sync"
  "io/ioutil"
  "net/http"
  "github.com/moovweb/gokogiri"
  "github.com/moovweb/gokogiri/html"
)

func parseHtml(url string) (*html.HtmlDocument, error) {
  resp, err := http.Get(url)
  if err != nil {
    return nil, err
  }
  defer resp.Body.Close()
  
  content, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    return nil, err
  }
  
  return gokogiri.ParseHtml(content)
}

func waitData(wg *sync.WaitGroup, ch chan interface{}, f func(interface{})) {
  go func() {
    wg.Wait()
    close(ch)
  }()
  
  loop := true
  for loop {
    select {
    case data, ok := <-ch:
      if ok {
        f(data)
      } else {
        loop = false
      }
    }
  }
}

func (h *Highway) parseFreeways() error {
  doc, err := parseHtml(freewayUrl)
  if err != nil {
    return err
  }
  defer doc.Free()
  
  fs, err := doc.Search("//select[@id='sec_selt']/option")
  if err != nil {
    return err
  }
  
  for _, f := range fs {
    h.Freeways = append(h.Freeways, &Freeway{f.Content(), f.Attr("value"), false, make([]string, 0)})
  }
  return nil
}
