package encoder

import (
	"bytes"
	"encoding/json"
)

var b bytes.Buffer
var jsonEnc *json.Encoder

func init() {
	jsonEnc = json.NewEncoder(&b)
	jsonEnc.SetEscapeHTML(false)
}

func EncodeJSON(v interface{}) ([]byte, error) {
	err := jsonEnc.Encode(v)
	if err != nil {
		return []byte{}, err
	}

	bs := bytes.TrimSuffix(b.Bytes(), []byte{'\n'})
	b.Reset()
	return bs, nil
}
