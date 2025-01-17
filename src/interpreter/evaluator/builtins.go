package evaluator

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/KhushPatibandha/Kolon/src/interpreter/object"
)

var builtins = map[string]*object.Builtin{
	"len": {
		Fn: func(args ...object.Object) (object.Object, bool, error) {
			if len(args) != 1 {
				return NULL, true, errors.New("Wrong number of arguments. got=" + strconv.Itoa(len(args)) + ", want=1")
			}
			switch arg := args[0].(type) {
			case *object.String:
				return &object.Integer{Value: int64(len(arg.Value)) - 2}, false, nil
			case *object.Array:
				return &object.Integer{Value: int64(len(arg.Elements))}, false, nil
			case *object.Hash:
				return &object.Integer{Value: int64(len(arg.Pairs))}, false, nil
			default:
				return NULL, true, errors.New("Argument to `len` not supported, got " + string(args[0].Type()))
			}
		},
	},
	"toString": {
		Fn: func(args ...object.Object) (object.Object, bool, error) {
			if len(args) != 1 {
				return NULL, true, errors.New("Wrong number of arguments. got=" + strconv.Itoa(len(args)) + ", want=1")
			}
			switch arg := args[0].(type) {
			case *object.Integer:
				newStr := strconv.FormatInt(arg.Value, 10)
				newStr = "\"" + newStr + "\""
				return &object.String{Value: newStr}, false, nil
			case *object.Float:
				newStr := strconv.FormatFloat(arg.Value, 'f', -1, 64)
				newStr = "\"" + newStr + "\""
				return &object.String{Value: newStr}, false, nil
			case *object.Boolean:
				newStr := strconv.FormatBool(arg.Value)
				newStr = "\"" + newStr + "\""
				return &object.String{Value: newStr}, false, nil
			case *object.Char:
				newStr := arg.Value[1 : len(arg.Value)-1]
				newStr = "\"" + newStr + "\""
				return &object.String{Value: newStr}, false, nil
			case *object.String:
				return arg, false, nil
			case *object.Array:
				return &object.String{Value: arg.Inspect()}, false, nil
			case *object.Hash:
				return &object.String{Value: arg.Inspect()}, false, nil
			default:
				return NULL, true, errors.New("Argument to `toString` not supported, got " + string(args[0].Type()))
			}
		},
	},
	"print": {
		Fn: func(args ...object.Object) (object.Object, bool, error) {
			if len(args) != 1 {
				return NULL, true, errors.New("Wrong number of arguments. got=" + strconv.Itoa(len(args)) + ", want=1")
			}
			switch arg := args[0].(type) {
			case *object.String:
				newStr := arg.Value[1 : len(arg.Value)-1]
				fmt.Print(newStr)
				return NULL, false, nil
			case *object.Char:
				newStr := arg.Value[1 : len(arg.Value)-1]
				fmt.Print(newStr)
				return NULL, false, nil
			case *object.Integer:
				newStr := strconv.FormatInt(arg.Value, 10)
				fmt.Print(newStr)
				return NULL, false, nil
			case *object.Float:
				newStr := strconv.FormatFloat(arg.Value, 'f', -1, 64)
				fmt.Print(newStr)
				return NULL, false, nil
			case *object.Boolean:
				fmt.Print(arg.Value)
				return NULL, false, nil
			case *object.Array:
				fmt.Print(arg.Inspect())
				return NULL, false, nil
			case *object.Hash:
				fmt.Print(arg.Inspect())
				return NULL, false, nil
			case *object.Null:
				return NULL, false, nil
			default:
				return NULL, true, errors.New("Argument to `print` not supported, can only print strings, got " + string(args[0].Type()) + " convert " + string(args[0].Type()) + " to string using toString()")
			}
		},
	},
	"println": {
		Fn: func(args ...object.Object) (object.Object, bool, error) {
			if len(args) != 1 {
				return NULL, true, errors.New("Wrong number of arguments. got=" + strconv.Itoa(len(args)) + ", want=1")
			}
			switch arg := args[0].(type) {
			case *object.String:
				newStr := arg.Value[1 : len(arg.Value)-1]
				fmt.Println(newStr)
				return NULL, false, nil
			case *object.Char:
				newStr := arg.Value[1 : len(arg.Value)-1]
				fmt.Println(newStr)
				return NULL, false, nil
			case *object.Integer:
				newStr := strconv.FormatInt(arg.Value, 10)
				fmt.Println(newStr)
				return NULL, false, nil
			case *object.Float:
				newStr := strconv.FormatFloat(arg.Value, 'f', -1, 64)
				fmt.Println(newStr)
				return NULL, false, nil
			case *object.Boolean:
				fmt.Println(arg.Value)
				return NULL, false, nil
			case *object.Array:
				fmt.Println(arg.Inspect())
				return NULL, false, nil
			case *object.Hash:
				fmt.Println(arg.Inspect())
				return NULL, false, nil
			case *object.Null:
				return NULL, false, nil
			default:
				return NULL, true, errors.New("Argument to `println` not supported, can only print strings, got " + string(args[0].Type()) + " convert " + string(args[0].Type()) + " to string using toString()")
			}
		},
	},
	"push": {
		Fn: func(args ...object.Object) (object.Object, bool, error) {
			switch arg := args[0].(type) {
			case *object.Array:
				if len(args) != 2 {
					return NULL, true, errors.New("Wrong number of arguments for `push` on array. got=" + strconv.Itoa(len(args)) + ", want=2. `push(array, element)`")
				}
				arrayType := arg.TypeOf
				switch arrayType {
				case "int":
					arg2, ok := args[1].(*object.Integer)
					if !ok {
						return NULL, true, errors.New("Argument type mismatch. Expected int, got " + string(args[1].Type()))
					}
					arg.Elements = append(arg.Elements, arg2)
					return arg, false, nil
				case "float":
					arg2, ok := args[1].(*object.Float)
					if !ok {
						return NULL, true, errors.New("Argument type mismatch. Expected float, got " + string(args[1].Type()))
					}
					arg.Elements = append(arg.Elements, arg2)
					return arg, false, nil
				case "bool":
					arg2, ok := args[1].(*object.Boolean)
					if !ok {
						return NULL, true, errors.New("Argument type mismatch. Expected bool, got " + string(args[1].Type()))
					}
					arg.Elements = append(arg.Elements, arg2)
					return arg, false, nil
				case "char":
					arg2, ok := args[1].(*object.Char)
					if !ok {
						return NULL, true, errors.New("Argument type mismatch. Expected char, got " + string(args[1].Type()))
					}
					arg.Elements = append(arg.Elements, arg2)
					return arg, false, nil
				case "string":
					arg2, ok := args[1].(*object.String)
					if !ok {
						return NULL, true, errors.New("Argument type mismatch. Expected string, got " + string(args[1].Type()))
					}
					arg.Elements = append(arg.Elements, arg2)
					return arg, false, nil
				default:
					return NULL, true, errors.New("Array type not supported. got=" + arrayType)
				}
			case *object.Hash:
				if len(args) != 3 {
					return NULL, true, errors.New("Wrong number of arguments for `push` on hash. got=" + strconv.Itoa(len(args)) + ", want=3. `push(map, key, value)`")
				}
				keyType := arg.KeyType
				valueType := arg.ValueType
				var arg2 object.Object
				var arg3 object.Object
				var ok bool
				switch keyType {
				case "int":
					arg2, ok = args[1].(*object.Integer)
					if !ok {
						return NULL, true, errors.New("Key type mismatch. Expected int, got " + string(args[1].Type()))
					}
				case "float":
					arg2, ok = args[1].(*object.Float)
					if !ok {
						return NULL, true, errors.New("Key type mismatch. Expected float, got " + string(args[1].Type()))
					}
				case "bool":
					arg2, ok = args[1].(*object.Boolean)
					if !ok {
						return NULL, true, errors.New("Key type mismatch. Expected bool, got " + string(args[1].Type()))
					}
				case "char":
					arg2, ok = args[1].(*object.Char)
					if !ok {
						return NULL, true, errors.New("Key type mismatch. Expected char, got " + string(args[1].Type()))
					}
				case "string":
					arg2, ok = args[1].(*object.String)
					if !ok {
						return NULL, true, errors.New("Key type mismatch. Expected string, got " + string(args[1].Type()))
					}
				default:
					return NULL, true, errors.New("Key type not supported. got=" + keyType)
				}

				switch valueType {
				case "int":
					arg3, ok = args[2].(*object.Integer)
					if !ok {
						return NULL, true, errors.New("Value type mismatch. Expected int, got " + string(args[2].Type()))
					}
				case "float":
					arg3, ok = args[2].(*object.Float)
					if !ok {
						return NULL, true, errors.New("Value type mismatch. Expected float, got " + string(args[2].Type()))
					}
				case "bool":
					arg3, ok = args[2].(*object.Boolean)
					if !ok {
						return NULL, true, errors.New("Value type mismatch. Expected bool, got " + string(args[2].Type()))
					}
				case "char":
					arg3, ok = args[2].(*object.Char)
					if !ok {
						return NULL, true, errors.New("Value type mismatch. Expected char, got " + string(args[2].Type()))
					}
				case "string":
					arg3, ok = args[2].(*object.String)
					if !ok {
						return NULL, true, errors.New("Value type mismatch. Expected string, got " + string(args[2].Type()))
					}
				default:
					return NULL, true, errors.New("Value type not supported. got=" + valueType)
				}

				hashKey, ok := arg2.(object.Hashable)
				if !ok {
					return NULL, true, errors.New("Key type not hashable")
				}
				hashed := hashKey.HashKey()
				arg.Pairs[hashed] = object.HashPair{Key: arg2, Value: arg3}
				return arg, false, nil
			default:
				return NULL, true, errors.New("Data structure not supported by `push`. got=" + string(args[0].Type()))
			}
		},
	},
	"pop": {
		Fn: func(args ...object.Object) (object.Object, bool, error) {
			switch arg := args[0].(type) {
			case *object.Array:
				if len(args) != 1 && len(args) != 2 {
					return NULL, true, errors.New("Wrong number of arguments for `pop` on array. got=" + strconv.Itoa(len(args)) + ", want=1 or 2. `pop(array)` or `pop(array, index)`")
				}
				if len(args) == 1 {
					if len(arg.Elements) == 0 {
						return NULL, false, nil
					}
					popped := arg.Elements[len(arg.Elements)-1]
					arg.Elements = arg.Elements[:len(arg.Elements)-1]
					return popped, false, nil
				} else {
					index, ok := args[1].(*object.Integer)
					if !ok {
						return NULL, true, errors.New("Index must be an integer")
					}
					if index.Value < 0 || index.Value >= int64(len(arg.Elements)) {
						return NULL, true, errors.New("Index out of bounds")
					}
					popped := arg.Elements[index.Value]
					arg.Elements = append(arg.Elements[:index.Value], arg.Elements[index.Value+1:]...)
					return popped, false, nil
				}
			default:
				return NULL, true, errors.New("Data structure not supported by `pop`. got=" + string(args[0].Type()))
			}
		},
	},
	"insert": {
		Fn: func(args ...object.Object) (object.Object, bool, error) {
			switch arg := args[0].(type) {
			case *object.Array:
				if len(args) != 3 {
					return NULL, true, errors.New("Wrong number of arguments for `add` on array. got=" + strconv.Itoa(len(args)) + ", want=3. `insert(array, index, element)`")
				}
				index, ok := args[1].(*object.Integer)
				if !ok {
					return NULL, true, errors.New("Index must be an integer")
				}
				if index.Value < 0 || index.Value >= int64(len(arg.Elements)) {
					return NULL, true, errors.New("Index out of bounds")
				}
				arrayType := arg.TypeOf
				var arg2 object.Object
				switch arrayType {
				case "int":
					arg2, ok = args[2].(*object.Integer)
					if !ok {
						return NULL, true, errors.New("Argument type mismatch. Expected int, got " + string(args[2].Type()))
					}
				case "float":
					arg2, ok = args[2].(*object.Float)
					if !ok {
						return NULL, true, errors.New("Argument type mismatch. Expected float, got " + string(args[2].Type()))
					}
				case "bool":
					arg2, ok = args[2].(*object.Boolean)
					if !ok {
						return NULL, true, errors.New("Argument type mismatch. Expected bool, got " + string(args[2].Type()))
					}
				case "char":
					arg2, ok = args[2].(*object.Char)
					if !ok {
						return NULL, true, errors.New("Argument type mismatch. Expected char, got " + string(args[2].Type()))
					}
				case "string":
					arg2, ok = args[2].(*object.String)
					if !ok {
						return NULL, true, errors.New("Argument type mismatch. Expected string, got " + string(args[2].Type()))
					}
				default:
					return NULL, true, errors.New("Array type not supported. got=" + arrayType)
				}
				temp := make([]object.Object, len(arg.Elements)+1)
				copy(temp, arg.Elements[:index.Value])
				temp[index.Value] = arg2
				copy(temp[index.Value+1:], arg.Elements[index.Value:])
				arg.Elements = temp
				return arg, false, nil
			default:
				return NULL, true, errors.New("Data structure not supported by `insert`. got=" + string(args[0].Type()))
			}
		},
	},
	"remove": {
		Fn: func(args ...object.Object) (object.Object, bool, error) {
			switch arg := args[0].(type) {
			case *object.Array:
				if len(args) != 2 {
					return NULL, true, errors.New("Wrong number of arguments for `remove` on array. got=" + strconv.Itoa(len(args)) + ", want=2. `remove(array, element)`")
				}
				arrayType := arg.TypeOf
				switch arrayType {
				case "int":
					arg2, ok := args[1].(*object.Integer)
					if !ok {
						return NULL, true, errors.New("Argument type mismatch. Expected int, got " + string(args[1].Type()))
					}
					for i, element := range arg.Elements {
						if element.(*object.Integer).Value == arg2.Value {
							arg.Elements = append(arg.Elements[:i], arg.Elements[i+1:]...)
							return arg, false, nil
						}
					}
					return NULL, false, nil
				case "float":
					arg2, ok := args[1].(*object.Float)
					if !ok {
						return NULL, true, errors.New("Argument type mismatch. Expected float, got " + string(args[1].Type()))
					}
					for i, element := range arg.Elements {
						if element.(*object.Float).Value == arg2.Value {
							arg.Elements = append(arg.Elements[:i], arg.Elements[i+1:]...)
							return arg, false, nil
						}
					}
					return NULL, false, nil
				case "bool":
					arg2, ok := args[1].(*object.Boolean)
					if !ok {
						return NULL, true, errors.New("Argument type mismatch. Expected bool, got " + string(args[1].Type()))
					}
					for i, element := range arg.Elements {
						if element.(*object.Boolean).Value == arg2.Value {
							arg.Elements = append(arg.Elements[:i], arg.Elements[i+1:]...)
							return arg, false, nil
						}
					}
					return NULL, false, nil
				case "char":
					arg2, ok := args[1].(*object.Char)
					if !ok {
						return NULL, true, errors.New("Argument type mismatch. Expected char, got " + string(args[1].Type()))
					}
					for i, element := range arg.Elements {
						if element.(*object.Char).Value == arg2.Value {
							arg.Elements = append(arg.Elements[:i], arg.Elements[i+1:]...)
							return arg, false, nil
						}
					}
					return NULL, false, nil
				case "string":
					arg2, ok := args[1].(*object.String)
					if !ok {
						return NULL, true, errors.New("Argument type mismatch. Expected string, got " + string(args[1].Type()))
					}
					for i, element := range arg.Elements {
						if element.(*object.String).Value == arg2.Value {
							arg.Elements = append(arg.Elements[:i], arg.Elements[i+1:]...)
							return arg, false, nil
						}
					}
					return NULL, false, nil
				default:
					return NULL, true, errors.New("Array type not supported. got=" + arrayType)
				}
			case *object.Hash:
				if len(args) != 2 {
					return NULL, true, errors.New("Wrong number of arguments for `remove` on hash. got=" + strconv.Itoa(len(args)) + ", want=2. `remove(map, key)`")
				}
				keyType := arg.KeyType
				var arg2 object.Object
				var ok bool

				switch keyType {
				case "int":
					arg2, ok = args[1].(*object.Integer)
					if !ok {
						return NULL, true, errors.New("Key type mismatch. Expected int, got " + string(args[1].Type()))
					}
				case "float":
					arg2, ok = args[1].(*object.Float)
					if !ok {
						return NULL, true, errors.New("Key type mismatch. Expected float, got " + string(args[1].Type()))
					}
				case "bool":
					arg2, ok = args[1].(*object.Boolean)
					if !ok {
						return NULL, true, errors.New("Key type mismatch. Expected bool, got " + string(args[1].Type()))
					}
				case "char":
					arg2, ok = args[1].(*object.Char)
					if !ok {
						return NULL, true, errors.New("Key type mismatch. Expected char, got " + string(args[1].Type()))
					}
				case "string":
					arg2, ok = args[1].(*object.String)
					if !ok {
						return NULL, true, errors.New("Key type mismatch. Expected string, got " + string(args[1].Type()))
					}
				default:
					return NULL, true, errors.New("Key type not supported. got=" + keyType)
				}

				hashKey, ok := arg2.(object.Hashable)
				if !ok {
					return NULL, true, errors.New("Key type not hashable")
				}
				hashed := hashKey.HashKey()
				value, ok := arg.Pairs[hashed]
				if !ok {
					return NULL, true, errors.New("Key not found. can't remove pair that doesn't exist")
				}
				delete(arg.Pairs, hashed)
				return value.Value, false, nil
			default:
				return NULL, true, errors.New("Data structure not supported by `remove`. got=" + string(args[0].Type()))
			}
		},
	},
	"getIndex": {
		Fn: func(args ...object.Object) (object.Object, bool, error) {
			switch arg := args[0].(type) {
			case *object.Array:
				if len(args) != 2 {
					return NULL, true, errors.New("Wrong number of arguments for `getIndex` on array. got=" + strconv.Itoa(len(args)) + ", want=2. `getIndex(array, element)`")
				}
				arrayType := arg.TypeOf
				switch arrayType {
				case "int":
					arg2, ok := args[1].(*object.Integer)
					if !ok {
						return NULL, true, errors.New("Argument type mismatch. Expected int, got " + string(args[1].Type()))
					}
					for i, element := range arg.Elements {
						if element.(*object.Integer).Value == arg2.Value {
							return &object.Integer{Value: int64(i)}, false, nil
						}
					}
					return &object.Integer{Value: -1}, false, nil
				case "float":
					arg2, ok := args[1].(*object.Float)
					if !ok {
						return NULL, true, errors.New("Argument type mismatch. Expected float, got " + string(args[1].Type()))
					}
					for i, element := range arg.Elements {
						if element.(*object.Float).Value == arg2.Value {
							return &object.Integer{Value: int64(i)}, false, nil
						}
					}
					return &object.Integer{Value: -1}, false, nil
				case "bool":
					arg2, ok := args[1].(*object.Boolean)
					if !ok {
						return NULL, true, errors.New("Argument type mismatch. Expected bool, got " + string(args[1].Type()))
					}
					for i, element := range arg.Elements {
						if element.(*object.Boolean).Value == arg2.Value {
							return &object.Integer{Value: int64(i)}, false, nil
						}
					}
					return &object.Integer{Value: -1}, false, nil
				case "char":
					arg2, ok := args[1].(*object.Char)
					if !ok {
						return NULL, true, errors.New("Argument type mismatch. Expected char, got " + string(args[1].Type()))
					}
					for i, element := range arg.Elements {
						if element.(*object.Char).Value == arg2.Value {
							return &object.Integer{Value: int64(i)}, false, nil
						}
					}
					return &object.Integer{Value: -1}, false, nil
				case "string":
					arg2, ok := args[1].(*object.String)
					if !ok {
						return NULL, true, errors.New("Argument type mismatch. Expected string, got " + string(args[1].Type()))
					}
					for i, element := range arg.Elements {
						if element.(*object.String).Value == arg2.Value {
							return &object.Integer{Value: int64(i)}, false, nil
						}
					}
					return &object.Integer{Value: -1}, false, nil
				default:
					return NULL, true, errors.New("Array type not supported. got=" + arrayType)
				}
			default:
				return NULL, true, errors.New("Data structure not supported by `getIndex`. got=" + string(args[0].Type()))
			}
		},
	},
	"keys": {
		Fn: func(args ...object.Object) (object.Object, bool, error) {
			switch arg := args[0].(type) {
			case *object.Hash:
				if len(args) != 1 {
					return NULL, true, errors.New("Wrong number of arguments for `keys` on hash. got=" + strconv.Itoa(len(args)) + ", want=1. `keys(map)`")
				}
				var keys []object.Object
				for _, pair := range arg.Pairs {
					keys = append(keys, pair.Key)
				}
				return &object.Array{Elements: keys, TypeOf: arg.KeyType}, false, nil
			default:
				return NULL, true, errors.New("Data structure not supported by `keys`. got=" + string(args[0].Type()))
			}
		},
	},
	"values": {
		Fn: func(args ...object.Object) (object.Object, bool, error) {
			switch arg := args[0].(type) {
			case *object.Hash:
				if len(args) != 1 {
					return NULL, true, errors.New("Wrong number of arguments for `values` on hash. got=" + strconv.Itoa(len(args)) + ", want=1. `values(map)`")
				}
				var values []object.Object
				for _, pair := range arg.Pairs {
					values = append(values, pair.Value)
				}
				return &object.Array{Elements: values, TypeOf: arg.ValueType}, false, nil
			default:
				return NULL, true, errors.New("Data structure not supported by `values`. got=" + string(args[0].Type()))
			}
		},
	},
	"containsKey": {
		Fn: func(args ...object.Object) (object.Object, bool, error) {
			switch arg := args[0].(type) {
			case *object.Hash:
				if len(args) != 2 {
					return NULL, true, errors.New("Wrong number of arguments for `containsKey` on hash. got=" + strconv.Itoa(len(args)) + ", want=2. `containsKey(map, key)`")
				}
				keyType := arg.KeyType
				var arg2 object.Object
				var ok bool
				switch keyType {
				case "int":
					arg2, ok = args[1].(*object.Integer)
					if !ok {
						return NULL, true, errors.New("Key type mismatch. Expected int, got " + string(args[1].Type()))
					}
				case "float":
					arg2, ok = args[1].(*object.Float)
					if !ok {
						return NULL, true, errors.New("Key type mismatch. Expected float, got " + string(args[1].Type()))
					}
				case "bool":
					arg2, ok = args[1].(*object.Boolean)
					if !ok {
						return NULL, true, errors.New("Key type mismatch. Expected bool, got " + string(args[1].Type()))
					}
				case "char":
					arg2, ok = args[1].(*object.Char)
					if !ok {
						return NULL, true, errors.New("Key type mismatch. Expected char, got " + string(args[1].Type()))
					}
				case "string":
					arg2, ok = args[1].(*object.String)
					if !ok {
						return NULL, true, errors.New("Key type mismatch. Expected string, got " + string(args[1].Type()))
					}
				default:
					return NULL, true, errors.New("Key type not supported. got=" + keyType)
				}
				hashKey, ok := arg2.(object.Hashable)
				if !ok {
					return NULL, true, errors.New("Key type not hashable")
				}
				hashed := hashKey.HashKey()
				_, ok = arg.Pairs[hashed]
				if ok {
					return TRUE, false, nil
				}
				return FALSE, false, nil
			default:
				return NULL, true, errors.New("Data structure not supported by `containsKey`. got=" + string(args[0].Type()))
			}
		},
	},
}
