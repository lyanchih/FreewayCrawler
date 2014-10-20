package freeway

import (
  "github.com/lyanchih/goamf"
)

func checkString(v interface{}) string {
  if str, ok := v.(string); ok {
    return str
  }
  return ""
}

func RequestGraphSecs(freewayId, locId string) ([]*GraphSection, error) {
  pack, err := goamf.NewAmfPacket(goamf.AMF3)
  if err != nil {
    return nil, err
  }
  
  err = setGraphSections(pack, freewayId, locId)
  if err != nil {
    return nil, err
  }
  
  resp, err := sendRequest(pack)
  if err != nil {
    return nil, err
  }
  
  graphSections := make([]*GraphSection, 0)
  for _, message := range resp.Messages {
    obj, ok := message.Value.(*goamf.AMF3Object)
    if !ok {
      continue
    }
    
    body, ok := obj.Values["body"].(*goamf.AMF3Object)
    if !ok {
      continue
    }
    
    arr, ok := body.Values["secArray"].(*goamf.AMF3Array)
    if !ok || arr.DenseValues == nil {
      continue
    }
    
    sectionDatas := make([]*SectionData, 0, len(arr.DenseValues))
    for _, v := range arr.DenseValues {
      obj, ok := v.(*goamf.AMF3Object)
      if !ok {
        continue
      }
      
      sectionDatas = append(sectionDatas, &SectionData{
        checkString(obj.Values["locationType"]),
        checkString(obj.Values["freewayId"]),
        checkString(obj.Values["mainDirection"]),
        checkString(obj.Values["locationId"]),
        checkString(obj.Values["locationName"]),
        checkString(obj.Values["fmileAge"]),
        checkString(obj.Values["tmileAge"]),
        checkString(obj.Values["ftype"]),
        checkString(obj.Values["ttype"]),
        checkString(obj.Values["numberOfLane"]),
      })
    }
    
    graphSections = append(graphSections, &GraphSection{
      checkString(obj.Values["timestamp"]),
      sectionDatas,
    })
  }
  return graphSections, nil
}


func RequestLocationInfos(locIds []string) ([]*LocationInfo, error) {
  pack, err := goamf.NewAmfPacket(goamf.AMF3)
  if err != nil {
    return nil, err
  }
  
  err = setLocationInfos(pack, locIds)
  if err != nil {
    return nil, err
  }
  
  resp, err := sendRequest(pack)
  if err != nil {
    return nil, err
  }

  locationInfos := make([]*LocationInfo, 0, 1)
  for _, message := range resp.Messages {
    obj, ok := message.Value.(*goamf.AMF3Object)
    if !ok {
      continue
    }
    
    body, ok := obj.Values["body"].(*goamf.AMF3Object)
    if !ok {
      continue
    }
    
    cctvArr, ok := body.Values["cctvs"].(*goamf.AMF3Array)
    if !ok {
      continue
    }
    
    cctvs := make([]*Cctv, 0, len(cctvArr.DenseValues))
    for _, cctv := range cctvArr.DenseValues {
      obj, ok := cctv.(*goamf.AMF3Object)
      if !ok {
        continue
      }
      
      cctvs = append(cctvs, &Cctv{
        checkString(obj.Values["id"]),
        checkString(obj.Values["ch"]),
        checkString(obj.Values["type"]),
        checkString(obj.Values["milepost"]),
        checkString(obj.Values["direction"]),
        checkString(obj.Values["locationId"]),
        checkString(obj.Values["gifUrl"]),
        checkString(obj.Values["videoUrl"]),
        checkString(obj.Values["enabled"]),
      })
    }
    
    speedArr, ok := body.Values["avgSpds"].(*goamf.AMF3Array)
    if !ok || speedArr.DenseValues == nil {
      continue
    }
    
    speeds := make([]*LocationSpeed, 0, len(speedArr.DenseValues))
    for _, speed := range speedArr.DenseValues {
      arr, ok := speed.(*goamf.AMF3Array)
      if !ok || arr.AssocValues == nil {
        continue
      }
      
      speeds = append(speeds, &LocationSpeed{
        checkString(arr.AssocValues["directionId"]),
        checkString(arr.AssocValues["averageSpeed"]),
        checkString(arr.AssocValues["hov_speed"]),
        checkString(arr.AssocValues["none_hov_speed"]),
      })
    }
    
    locationInfos = append(locationInfos, &LocationInfo{
      checkString(obj.Values["timestamp"]),
      checkString(body.Values["locationId"]),
      checkString(body.Values["freewayId"]),
      checkString(body.Values["direction"]),
      checkString(body.Values["fmile"]),
      checkString(body.Values["tmile"]),
      checkString(body.Values["displayStartMile"]),
      checkString(body.Values["displayEndMile"]),
      cctvs,
      speeds,
    })
  }
  
  return locationInfos, nil
}
