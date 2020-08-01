package crypto

import "strings"

func EasyEncode(data [][]byte) string {
	var tmp []string
	for _, v := range data {
		tmp = append(tmp, Base64Encode(v))
	}
	return strings.Join(tmp, ",")
}

func EasyDecode(data string) ([][]byte, error) {
	tmp := strings.Split(data, ",")

	var res [][]byte
	for _, v := range tmp {
		r, err := Base64Decode(v)
		if err != nil {
			return nil, err
		}
		res = append(res, r)
	}
	return res, nil
}
