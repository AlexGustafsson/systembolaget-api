package systembolaget

// Source ...
type Source interface {
  ConvertToXML(pretty bool) ([]byte, error)
  ConvertToJSON(pretty bool) ([]byte, error)
}
