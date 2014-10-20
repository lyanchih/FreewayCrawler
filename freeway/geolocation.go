package freeway

import (
  "io"
  "log"
  "fmt"
  "sync"
  "errors"
  "strconv"
  "bufio"
  "regexp"
  "net/http"
)

var geolocationRegexp, _ = regexp.Compile(`fromX=(\d+\.\d+);\s*fromY=(\d+\.\d+);\s*toX=(\d+\.\d+);\s*toY=(\d+\.\d+)`)

func (h *Highway) parseGeolocations() (error) {
  wg := &sync.WaitGroup{}
  ch := make(chan *Geolocation, len(h.Locations) >> 1)
  queue := make(chan bool, maxHttpClient)
  
  for _, l := range h.Locations {
    wg.Add(1)
    go func(locId string) {
      defer func() {
        wg.Done()
        <-queue
      }()
      queue<- true
      geolocation, err := parseGeolocation(locId)
      if err != nil {
        log.Println(err)
        return
      }
      ch<- geolocation
    }(l.Id)
  }
  
  go func() {
    for {
      select {
      case geolocation, ok := <-ch:
        if ok {
          for _, l := range h.Locations {
            if l.Id == geolocation.locId {
              l.Geolocation = geolocation
              break
            }
          }
        } else {
          return
        }
      }
    }
  }()
  
  wg.Wait()
  close(ch)
  return nil
}

func parseGeolocation(locId string) (*Geolocation, error) {
  url := fmt.Sprintf(geolocationUrl, locId)
  resp, err := http.Get(url);
  if err != nil {
    return nil, err
  }
  defer resp.Body.Close()
  r := bufio.NewReader(resp.Body)
  for {
    line, _, err := r.ReadLine()
    if err != nil {
      if err == io.EOF {
        err = errors.New(url + " can't find location's latitide and longitude")
      }
      return nil, err
    }
    
    result := geolocationRegexp.FindSubmatch(line)
    if result == nil {
      continue
    }
    fromX, err := strconv.ParseFloat(string(result[1]), 32)
    if err != nil {
      log.Println(err)
      return nil, err
    }
    fromY, err := strconv.ParseFloat(string(result[2]), 32)
    if err != nil {
      log.Println(err)
      return nil, err
    }
    toX, err := strconv.ParseFloat(string(result[3]), 32)
    if err != nil {
      log.Println(err)
      return nil, err
    }
    toY, err := strconv.ParseFloat(string(result[4]), 32)
    if err != nil {
      log.Println(err)
      return nil, err
    }
    return &Geolocation{float32(fromX), float32(fromY), float32(toX), float32(toY), locId}, nil
  }
}
