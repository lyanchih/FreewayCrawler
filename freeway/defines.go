package freeway

const (
  host = "http://1968.freeway.gov.tw"
  freewayUrl = "http://1968.freeway.gov.tw/common/getfrees?id=sec_selt&df=1"
  interchangeFromUrl = "http://1968.freeway.gov.tw/common/getnodsecs/fid/%s?id=from_selt"
  interchangeToUrl = "http://1968.freeway.gov.tw/common/getnodsecs/fid/%s?lc=1&id=end_selt"
  locationUrl = "http://1968.freeway.gov.tw/traffic/index/fid/%s"
  geolocationUrl = "http://1968.freeway.gov.tw/section/getlocationinfo/?loc=%s"
  maxHttpClient = 10
)

type Highway struct {
  Freeways []*Freeway `json:"freeways"`
  Interchanges []*Interchange `json:"interchanges"`
  Locations []*Location `json:"locations"`
}

type Freeway struct {
  Name string `json:"name"`
  Id string `json:"id"`
  Direction bool `json:"direction"`
  Locations []string `json:"locations"`
}

type Interchange struct {
  Name string `json:"name"`
  Id string `json:"id"`
  FreewayId string `json:"freeway_id"`
  Locations []string `json:"locations"`
}

type Location struct {
  Id string `json:"id"`
  FreewayId string `json:"freeway_id"`
  Interchanges []string `json:"interchanges"`
  *Geolocation
}

type Geolocation struct {
  FromX float32 `json:"from_x"`
  FromY float32 `json:"from_y"`
  ToX float32 `json:"to_x"`
  ToY float32 `json:"to_y"`
  locId string
}
