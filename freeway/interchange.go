package freeway

import (
  "fmt"
  "sync"
)

func (h *Highway) parseInterchanges() {
  wg := &sync.WaitGroup{}
  ch := make(chan interface{}, len(h.Freeways))
  for _, f := range h.Freeways {
    wg.Add(1)
    go parseInterchange(wg, ch, f.Id)
  }
  
  waitData(wg, ch, func(data interface{}) {
    if i, ok := data.(*Interchange); ok {
      h.Interchanges = append(h.Interchanges, i)
    }
  })
}

func parseInterchange(wg *sync.WaitGroup, ch chan<- interface{}, freewayId string) {
  defer func() {
    wg.Done()
  }()
  
  maps := make(map[string]interface{})
  f := func(url string) {
    doc, err := parseHtml(url)
    if err != nil {
      return
    }
    defer doc.Free()
    is, err := doc.Search("//select/option")
    if err != nil {
      return
    }
    
    for _, i := range is {
      name := i.Content()
      value := i.Attr("value")
      id := freewayId + "." + value
      if _, ok := maps[id]; ok {
        continue
      } 
      
      ch<- &Interchange{name, id, freewayId, make([]string, 0)}
      maps[id] = nil
    }
  }
  
  f(fmt.Sprintf(interchangeFromUrl, freewayId))
  f(fmt.Sprintf(interchangeToUrl, freewayId))
}
