package freeway

import (
  "fmt"
  "log"
  "sync"
  "regexp"
)

func (h *Highway) parseLocations() {
  wg := &sync.WaitGroup{}
  ch := make(chan interface{}, len(h.Freeways))
  
  for _, f := range h.Freeways {
    wg.Add(1)
    go parseLocation(wg, ch, f.Id, h.Interchanges)
  }
  
  waitData(wg, ch, func(data interface{}) {
    if l, ok := data.(*Location); ok {
      for _, interchangeId := range l.Interchanges {
        for _, i := range h.Interchanges {
          if interchangeId == i.Id {
            i.Locations = append(i.Locations, l.Id)
          }
        }
      }
      for _, f := range h.Freeways {
        if f.Id == l.FreewayId {
          f.Locations = append(f.Locations, l.Id)
        }
      }
      h.Locations = append(h.Locations, l)
    }
  })
}

func parseLocation(wg *sync.WaitGroup, ch chan<- interface{}, freewayId string, interchanges []*Interchange) {
  defer func() {
    wg.Done()
  }()
  
  maps := make(map[string]interface{})
  
  doc, err := parseHtml(fmt.Sprintf(locationUrl, freewayId))
  if err != nil {
    log.Println(err)
    return
  }
  defer doc.Free()
  
  rows, err := doc.Search("//tbody[@id='secs_body']/tr")
  if err != nil {
    log.Println(err)
    return
  }
  
  nameReg, err := regexp.Compile(`^(.+)\([^)]+\) - (.+)\([^)]+\)$`)
  if err != nil {
    log.Println(err)
    return
  }
  linkReg, err := regexp.Compile(`/loc/(\d+)$`)
  if err != nil {
    log.Println(err)
    return
  }
  
  for _, row := range rows {
    arr, err := row.Search("td[@class='sec_name'] | td[@class='sec_detail']/a/@href")
    
    if err != nil || len(arr) < 2 {
      log.Println("row xpath search fail")
      continue
    }
    
    names := nameReg.FindStringSubmatch(arr[0].Content())
    if names == nil || len(names) < 3 {
      log.Println("location name regex match fail", arr[0].Content())
      continue
    }
    
    interchangeIds := make([]string, 2)
    for index, name := range names[1:] {
      for _, i := range interchanges {
        if i.Name == name && i.FreewayId == freewayId {
          interchangeIds[index] = i.Id
        }
      }
    }
    if len(interchangeIds) != 2 {
      log.Println("location can't find their interchanges")
      continue
    }
    
    link := linkReg.FindStringSubmatch(arr[1].Content())
    if link == nil || len(link) < 2 {
      continue
    }
    id := link[1]
    
    ch<- &Location{id, freewayId, interchangeIds, nil}
    maps[id] = nil
  }
}
