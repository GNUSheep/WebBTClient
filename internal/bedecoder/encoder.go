package bedecoder

import (
	"bytes"
	"errors"
	"fmt"
	"sort"
)

func Encode(buffer *bytes.Buffer, message any) error {
	switch elemt := message.(type) {
	case int64:
		fmt.Fprintf(buffer, "i:%de", elemt)
	case []any:
		buffer.WriteByte('l')
		for _, item := range elemt {
			if err := Encode(buffer, item); err != nil {
				return err
			}
		}
		buffer.WriteByte('e')
	case map[string]any:
		buffer.WriteByte('d')

		keys := make([]string, 0, len(elemt))
		for key := range elemt {
			keys = append(keys, key)
		}
		sort.Strings(keys)

		for _, key := range keys {
			fmt.Fprintf(buffer, "%d:%s", len(key), key)
			if err := Encode(buffer, elemt[key]); err != nil {
				return err
			}
		}
		buffer.WriteByte('e')
	case string:
		fmt.Fprintf(buffer, "%d:%s", len(elemt), elemt)
	default:
		return errors.New(fmt.Sprintf("Unsupported type in becoded message: %T", elemt))
	}

	return nil
}
