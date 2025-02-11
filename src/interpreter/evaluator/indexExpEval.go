package evaluator

import (
	"errors"
	"strconv"

	"github.com/KhushPatibandha/Kolon/src/object"
)

func evalIndexExpression(left object.Object, index object.Object) (object.Object, bool, error) {
	if left.Type() == object.ARRAY_OBJ && index.Type() == object.INTEGER_OBJ {
		return evalArrayIndexExpression(left, index)
	} else {
		return evalHashIndexExpression(left, index)
	}
}

func evalArrayIndexExpression(array object.Object, index object.Object) (object.Object, bool, error) {
	arrayObj := array.(*object.Array)
	idx := index.(*object.Integer).Value
	maxIdx := int64(len(arrayObj.Elements) - 1)
	if idx < 0 || idx > maxIdx {
		return NULL, true, errors.New("index out of range, index: " + strconv.FormatInt(idx, 10) + " max index: " + strconv.FormatInt(maxIdx, 10))
	}
	return arrayObj.Elements[idx], false, nil
}

func evalHashIndexExpression(hash object.Object, index object.Object) (object.Object, bool, error) {
	hashObj := hash.(*object.Hash)

	key, ok := index.(object.Hashable)
	if !ok {
		return NULL, true, errors.New("unusable as hash key: " + string(index.Type()))
	}
	pair, ok := hashObj.Pairs[key.HashKey()]
	if !ok {
		return NULL, true, errors.New("key not found: " + index.Inspect())
	}
	return pair.Value, false, nil
}
