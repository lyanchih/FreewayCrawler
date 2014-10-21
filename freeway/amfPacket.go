package freeway

import (
  "io"
  "fmt"
  "time"
  "bytes"
  "errors"
  "net/http"
  "github.com/lyanchih/goamf"
  "code.google.com/p/go-uuid/uuid"
)

type LocationSpeed struct {
  DirectionId, AverageSpeed, HovSpeed, NoneHovSpeed string
}

type SectionData struct {
  LocationType, FreewayId, Direction, LocationId, LocationName, FmileAge, TmileAge, Ftype, Ttype, NumberOfLane string
}

type GraphSection struct {
  Timestamp time.Time
  SectionDatas []*SectionData
}

type Cctv struct {
  Id, Ch, Type, Milepost, Direction, LocationId, GifUrl, VideoUrl, Enabled string
}

type LocationInfo struct {
  Timestamp time.Time
  LocationId, FreewayId, Direction, Fmile, TMile, DisplayStartMile, DisplayEndMile string
  Cctvs []*Cctv
  Speeds []*LocationSpeed
}

func setGraphSections(pack *goamf.Packet, freewayId, locId string) (error) {
  if pack == nil {
    return errors.New("goamf Packet is nil")
  }
  
  obj2 := goamf.NewAMF3Object("", true)
  obj2.AddDynValue("DSEndpoint", uint(0))
  obj2.AddDynValue("DSId", "nil")
  
  arr := goamf.NewAMF3Array(2)
  arr.AddDenseValue(freewayId)
  arr.AddDenseValue(locId)
  
  obj := goamf.NewAMF3Object("flex.messaging.messages.RemotingMessage", false)
  obj.AddValue("source", "SectionService")
  obj.AddValue("operation", "getGraphSecs")
  obj.AddValue("clientId", nil)
  obj.AddValue("destination", "zend")
  obj.AddValue("messageId", uuid.New())
  obj.AddValue("timeToLive", uint(0))
  obj.AddValue("handleErroraders", obj2)
  obj.AddValue("timestamp", uint(0))
  obj.AddValue("body", arr)
  
  return pack.AddMessage("null", "/2", []interface{}{obj})
}

func setLocationInfos(p *goamf.Packet, locIds []string) error {
  if p == nil {
    return errors.New("goamf Packet is nil")
  }
  for index, locId := range locIds {
    obj2 := goamf.NewAMF3Object("", true)
    obj2.AddDynValue("DSEndpoint", nil)
    obj2.AddDynValue("DSId", "nil")
    
    arr := goamf.NewAMF3Array(1)
    arr.AddDenseValue(locId)
    
    obj := goamf.NewAMF3Object("flex.messaging.messages.RemotingMessage", false)
    obj.AddValue("source", "SectionService")
    obj.AddValue("operation", "getLocationInfo")
    obj.AddValue("clientId", nil)
    obj.AddValue("destination", "zend")
    obj.AddValue("messageId", uuid.New())
    obj.AddValue("timeToLive", uint(0))
    obj.AddValue("handleErroraders", obj2)
    obj.AddValue("timestamp", uint(0))
    obj.AddValue("body", arr)
    
    err := p.AddMessage("null", fmt.Sprintf("/%d", index+1), []interface{}{obj})
    if err != nil {
      return err
    }
  }
  return nil
}

func sendRequest(p *goamf.Packet) (*goamf.Packet, error) {
  if p == nil {
    return nil, errors.New("goamf Packet is nil")
  }
  
  bs, err := goamf.MarshalAmf0(p)
  if err != nil {
    return nil, err
  }
  
  resp, err := http.Post("http://1968.freeway.gov.tw/gateway", "application/x-amf", bytes.NewReader(bs))
  if err != nil {
    return nil, err
  }
  
  defer resp.Body.Close()
  buf := bytes.Buffer{}
  _, err = io.Copy(&buf, resp.Body)
  if err != nil {
    return nil, err
  }
  
  pack, err := goamf.UnmarshalAmf0(buf.Bytes())
  if err != nil {
    return nil, err
  }

  if p, ok := pack.(*goamf.Packet); ok {
    return p, nil
  }
  
  return nil, errors.New("The response obj can't convert to Amf Packet")
}
