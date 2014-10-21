package main

import (
  "io"
  "os"
  "log"
  "fmt"
  "strconv"
  "strings"
  "./freeway"
  "github.com/robfig/cron"
)

const (
  csvHeader = "timestamp,freeway_id,location_id,direction_1,speed_1,direction_2,speed_2,\n"
)

func dumpLocationInfos(h *freeway.Highway) {
  locIds := make([]string, 0, len(h.Locations))
  for _, l := range h.Locations {
    locIds = append(locIds, l.Id)
  }
  
  infos, err := freeway.RequestLocationInfos(locIds)
  if err != nil {
    log.Fatal(err)
  }
  
  if len(infos) == 0 {
    log.Fatal("Do not response any location info")
  }
  
  f, err := os.OpenFile(fmt.Sprintf("data/%s.csv", strconv.FormatInt(infos[0].Timestamp.Unix(), 10)), os.O_WRONLY | os.O_CREATE | os.O_TRUNC, 0600)
  if err != nil {
    log.Fatal(err)
  }
  defer f.Close()
  
  _, err = io.WriteString(f, csvHeader)
  if err != nil {
    log.Fatal(err)
  }
  for _, info := range infos {
    strs := []string{strconv.FormatInt(info.Timestamp.Unix(), 10), info.FreewayId, info.LocationId}
    for _, speed := range info.Speeds {
      strs = append(strs, speed.DirectionId, speed.AverageSpeed)
    }
    strs = append(strs, "\n")
    _, err := io.WriteString(f, strings.Join(strs, ","))
    if err != nil {
      log.Fatal(err)
    }
  }
}

func main() {
  var err error
  h, _ := freeway.LoadJson("highway.json")
  
  if h == nil {
    h, err = freeway.ParseHighway(true)
    if err != nil {
      log.Fatal(err)
    }
  }

  f := func() {
    log.Println("Fetch location info")
    dumpLocationInfos(h)
    log.Println(fmt.Sprintf("Save location info"))
  }
  
  c := cron.New()
  go func() {
    f()
    c.AddFunc("@every 5m", f)
    c.Start()
  }()
  
  select {}
}
