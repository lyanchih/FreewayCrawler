package main

import (
  "io"
  "os"
  "fmt"
  "log"
  "path"
  "time"
  "bytes"
  "runtime"
  "errors"
  "strconv"
  "strings"
  "./freeway"
  "github.com/robfig/cron"
)

var saveFolder = ""

func init() {
  if folder := os.Getenv("SaveFolder"); folder != "" {
    saveFolder = folder
  } else {
    saveFolder = "data"
  }
}

func timeToFilename(date time.Time) string {
  return path.Join(saveFolder, fmt.Sprintf("%4d%02d%02d%02d%02d.csv", date.Year(), date.Month(), date.Day(), date.Hour(), (int)(date.Minute()/5)*5))
}

func unixToFilename(sec int64) string {
  return timeToFilename(time.Unix(sec, 0).UTC())
}

func dumpLocationInfos(h *freeway.Highway) error {
  locIds := make([]string, 0, len(h.Locations))
  for _, l := range h.Locations {
    locIds = append(locIds, l.Id)
  }
  
  infos, err := freeway.RequestLocationInfos(locIds)
  if err != nil {
    return err
  }
  
  if len(infos) == 0 {
    return err
  }
  
  unixTime := infos[0].Timestamp.Unix()
  f, err := os.OpenFile(unixToFilename(unixTime), os.O_WRONLY | os.O_CREATE | os.O_TRUNC, 0600)
  if err != nil {
    return err
  }
  defer f.Close()
  
  buffer := &bytes.Buffer{}
  _, err = buffer.WriteString(csvHeader)
  if err != nil {
    return err
  }
  for _, info := range infos {
    strs := []string{strconv.FormatInt(unixTime, 10), info.FreewayId, info.LocationId}
    for _, speed := range info.Speeds {
      strs = append(strs, speed.DirectionId, speed.AverageSpeed)
    }
    _, err := buffer.WriteString(string(append([]byte(strings.Join(strs, ",")), '\r', '\n')))
    if err != nil {
      return err
    }
  }
  
  _, err = io.Copy(f, buffer)
  return err
}

func main() {
  var err error
  h, _ := freeway.LoadJson("highway.json")
  
  if h == nil {
    h, err = freeway.ParseHighway(true)
    if err != nil {
      log.Print(err)
    }
  }
  
  handleDumpLocationError := retryFactory(h)
  f := func() {
    defer func() {
      if r := recover(); r != nil {
        log.Printf("Run time panic %v", r)
        if err, ok := r.(runtime.Error); ok {
          handleDumpLocationError(err)
        } else if s, ok := r.(string); ok {
          handleDumpLocationError(errors.New(s))
        } else if err, ok := r.(error); ok {
          handleDumpLocationError(err)
        } else {
          log.Println("The run time panic can't be handled")
        }
      }
    }()
    
    log.Println("Fetch location info")
    if err := dumpLocationInfos(h); err != nil {
      handleDumpLocationError(err)
      return
    }
    log.Println("Save location info")
  }
  
  c := cron.New()
  go func() {
    f()
    c.AddFunc("@every 5m", f)
    c.Start()
  }()
  
  select {}
}

func retryFactory(h *freeway.Highway) func(err error) {
  retryChan := make(chan *retryInfo, retryBufferCapability)
  
  go func() {
    for {
      select {
      case info, ok := <-retryChan:
        if !ok {
          return
        }
        
        if info.retryCount == maxRetryCount {
          log.Printf("%s speed information is catch fail\n", info.date.String())
          continue
        }
        
        after := time.After(info.delayTime)
        <-after
        log.Printf("Refetch %s Count: %d\n", info.date.String(), info.retryCount)
        if err := dumpLocationInfos(h); err != nil {
          log.Println(err)
          info.delayTime = info.delayTime * 2
          info.retryCount = info.retryCount + 1
          retryChan<- info
          continue
        }
        log.Println("Refetch location information success")
      }
    }
  }()
  
  return func(err error) {
    if err == nil {
      return
    }
    
    log.Println(err)
    retryChan<- &retryInfo{time.Now(), retryBaseDelayTime, 1}
  }
}
