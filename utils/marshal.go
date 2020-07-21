package utils

import (
  "reflect"
  "encoding/xml"
)

// UnmarshalXML ...
func UnmarshalXML(bytes []byte, container interface{}) error {
  out := map[string]interface{}{}
  containerType := reflect.TypeOf(container)
  containerValue := reflect.ValueOf(container)

  var parsedValue map[string]*xml.RawMessage
  err := xml.Unmarshal(bytes, &parsedValue)
  if err != nil {
    return err
  }

  for i := 0; i < containerType.NumField(); i++ {
      field := containerType.Field(i)
      name := field.Tag.Get("systembolagetXML")
      if name == "" {
        name = field.Name
      }

      var value interface{}
      err = xml.Unmarshal(*parsedValue[name], &value)
      if err != nil {
        return err
      }

      field.Set(value)
  }
}
