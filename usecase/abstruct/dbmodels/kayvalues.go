package dbmodels

type KeyValue[X comparable, Y any] struct {
	key     X
	value   Y
	isEmpty bool
}

func NewKeyValue[X comparable, Y any](k X, v Y) KeyValue[X, Y] {
	return KeyValue[X, Y]{
		key:     k,
		value:   v,
		isEmpty: false,
	}
}

func EmptyKeyValue[X comparable, Y any]() KeyValue[X, Y] {
	return KeyValue[X, Y]{
		isEmpty: true,
	}
}

func KeyValueToMap[X comparable, Y any](kv []KeyValue[X, Y]) map[X]Y {
	res := make(map[X]Y)

	for _, item := range kv {
		if item.isEmpty {
			continue
		}
		res[item.key] = item.value
	}

	return res
}
