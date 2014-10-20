package freeway

import (
  "io"
  "os"
  "bytes"
  "encoding/json"
)

func ParseHighway(autoSave bool) (*Highway, error) {
  h := &Highway{make([]*Freeway, 0), make([]*Interchange, 0), make([]*Location, 0)}
  err := h.parseFreeways()
  if err != nil {
    return nil, err
  }
  
  h.parseInterchanges()
  h.parseLocations()
  h.parseGeolocations()
  
  if autoSave {
    h.SaveJson("highway.json")
  }
  
  return h, nil
}

func LoadJson(path string) (h *Highway, err error) {
  fd, err := os.Open(path)
  if err != nil {
    return nil, err
  }
  defer fd.Close()
  
  dec := json.NewDecoder(fd)
  err = dec.Decode(&h)
  return
}

func (h *Highway) SaveJson(path string) error {
  fd, err := os.OpenFile(path, os.O_WRONLY | os.O_CREATE | os.O_TRUNC, 0600)
  if err != nil {
    return err
  }
  defer fd.Close()

  bs, err := json.MarshalIndent(h, "", "  ")
  if err != nil {
    return err
  }

  io.Copy(fd, bytes.NewReader(bs))
  
  return err
}
