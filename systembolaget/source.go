package systembolaget

// Source represents the common base for each API source.
type Source interface {
  Download() error
  ParseFromXML(bytes []byte) error
  ParseFromJSON(bytes []byte) error
  ConvertToXML(pretty bool) ([]byte, error)
  ConvertToJSON(pretty bool) ([]byte, error)
}
