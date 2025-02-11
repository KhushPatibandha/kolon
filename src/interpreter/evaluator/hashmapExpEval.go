package evaluator

import (
	"errors"

	"github.com/KhushPatibandha/Kolon/src/ast"
	"github.com/KhushPatibandha/Kolon/src/object"
)

func evalHashMap(node *ast.HashMap, env *object.Environment) (object.Object, bool, error) {
	pairs := make(map[object.HashKey]object.HashPair)

	for keyNode, valueNode := range node.Pairs {
		key, hasErr, err := Eval(keyNode, env, inTesting)
		if err != nil {
			return NULL, hasErr, err
		}

		hashKey, ok := key.(object.Hashable)
		if !ok {
			return NULL, true, errors.New("unusable as hash key: " + string(key.Type()))
		}

		value, hasErr, err := Eval(valueNode, env, inTesting)
		if err != nil {
			return NULL, hasErr, err
		}

		hashed := hashKey.HashKey()
		pairs[hashed] = object.HashPair{Key: key, Value: value}
	}

	if node.KeyType != nil && node.ValueType != nil {
		return &object.Hash{Pairs: pairs, KeyType: node.KeyType.Value, ValueType: node.ValueType.Value}, false, nil
	}

	return &object.Hash{Pairs: pairs}, false, nil
}
