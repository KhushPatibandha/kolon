package evaluator

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/KhushPatibandha/Kolon/src/object"
)

var builtins = map[string]*object.Builtin{
	"len": {
		Fn: func(args ...object.Object) (object.Object, bool, error) {
			if len(args) != 1 {
				return NULL, true, errors.New("wrong number of arguments for `len`, got: " + strconv.Itoa(len(args)) + ", want: 1")
			}
			switch arg := args[0].(type) {
			case *object.String:
				return &object.Integer{Value: int64(len(arg.Value)) - 2}, false, nil
			case *object.Array:
				return &object.Integer{Value: int64(len(arg.Elements))}, false, nil
			case *object.Hash:
				return &object.Integer{Value: int64(len(arg.Pairs))}, false, nil
			default:
				return NULL, true, errors.New("argument for `len` not supported, got: " + string(args[0].Type()) + ", want: array, hashmap or `string`")
			}
		},
	},
	"toString": {
		Fn: func(args ...object.Object) (object.Object, bool, error) {
			if len(args) != 1 {
				return NULL, true, errors.New("wrong number of arguments for `toString`, got: " + strconv.Itoa(len(args)) + ", want: 1")
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
				return NULL, true, errors.New("argument for `toString` not supported, got: " + string(args[0].Type()) + ", want: array, hashmap, `int`, `float`, `bool`, `char` or `string`")
			}
		},
	},
	"toFloat": {
		Fn: func(args ...object.Object) (object.Object, bool, error) {
			if len(args) != 1 {
				return NULL, true, errors.New("wrong number of arguments for `toFloat`, got: " + strconv.Itoa(len(args)) + ", want: 1")
			}
			switch arg := args[0].(type) {
			case *object.Integer:
				return &object.Float{Value: float64(arg.Value)}, false, nil
			case *object.Float:
				return arg, false, nil
			case *object.String:
				newStr := arg.Value[1 : len(arg.Value)-1]
				floatValue, err := strconv.ParseFloat(newStr, 64)
				if err != nil {
					return NULL, true, errors.New("Error converting string to float, can't convert: " + newStr)
				}
				return &object.Float{Value: floatValue}, false, nil
			default:
				return NULL, true, errors.New("argument for `toFloat` not supported, got: " + string(args[0].Type()) + ", want: `int`, `float` or `string`")
			}
		},
	},
	"toInt": {
		Fn: func(args ...object.Object) (object.Object, bool, error) {
			if len(args) != 1 {
				return NULL, true, errors.New("wrong number of arguments for `toInt`, got: " + strconv.Itoa(len(args)) + ", want: 1")
			}
			switch arg := args[0].(type) {
			case *object.Integer:
				return arg, false, nil
			case *object.Float:
				return &object.Integer{Value: int64(arg.Value)}, false, nil
			case *object.Char:
				newStr := arg.Value[1 : len(arg.Value)-1]
				return &object.Integer{Value: int64(newStr[0])}, false, nil
			case *object.String:
				newStr := arg.Value[1 : len(arg.Value)-1]
				intValue, err := strconv.ParseInt(newStr, 10, 64)
				if err != nil {
					return NULL, true, errors.New("Error converting string to int, can't convert: " + newStr)
				}
				return &object.Integer{Value: intValue}, false, nil
			default:
				return NULL, true, errors.New("argument for `toInt` not supported, got: " + string(args[0].Type()) + ", want: `int`, `float`, `string` or `char`")
			}
		},
	},
	"print": {
		Fn: func(args ...object.Object) (object.Object, bool, error) {
			if len(args) != 1 {
				return NULL, true, errors.New("wrong number of arguments for `print`, got: " + strconv.Itoa(len(args)) + ", want: 1")
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
				if !strings.Contains(newStr, ".") {
					newStr += ".0"
				}
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
				return NULL, true, errors.New("argument to `print` not supported, got: " + string(args[0].Type()) + ", want: array, hashmap, `int`, `float`, `bool`, `char` or `string`. use `toString` to convert to `string` in case of using other datatypes with `string`")
			}
		},
	},
	"println": {
		Fn: func(args ...object.Object) (object.Object, bool, error) {
			if len(args) != 1 {
				return NULL, true, errors.New("wrong number of arguments for `println`, got: " + strconv.Itoa(len(args)) + ", want: 1")
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
				if !strings.Contains(newStr, ".") {
					newStr += ".0"
				}
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
				return NULL, true, errors.New("argument to `println` not supported, got: " + string(args[0].Type()) + ", want: array, hashmap, `int`, `float`, `bool`, `char` or `string`. use `toString` to convert to `string` in case of using other datatypes with `string`")
			}
		},
	},
	"scanln": {
		Fn: func(args ...object.Object) (object.Object, bool, error) {
			if len(args) != 0 && len(args) != 1 && len(args) != 2 {
				return NULL, true, errors.New("wrong number of arguments for `scanln`, got: " + strconv.Itoa(len(args)) + ", want: 0, 1 or 2")
			}
			if len(args) != 0 {
				strToPrint := args[0].(*object.String).Value[1 : len(args[0].(*object.String).Value)-1]
				if len(args) == 2 && args[1].Inspect() == "true" {
					fmt.Println(strToPrint)
				} else {
					fmt.Print(strToPrint)
				}
			}
			reader := bufio.NewReader(os.Stdin)
			var input string
			rawInput, err := reader.ReadString('\n')
			if err != nil {
				return NULL, true, errors.New("error reading input: " + err.Error())
			}
			input = strings.TrimSpace(rawInput)
			return &object.String{Value: "\"" + input + "\""}, false, nil
		},
	},
	"scan": {
		Fn: func(args ...object.Object) (object.Object, bool, error) {
			if len(args) != 0 && len(args) != 1 && len(args) != 2 {
				return NULL, true, errors.New("wrong number of arguments for `scan`, got: " + strconv.Itoa(len(args)) + ", want: 0, 1 or 2")
			}
			if len(args) != 0 {
				strToPrint := args[0].(*object.String).Value[1 : len(args[0].(*object.String).Value)-1]
				if len(args) == 2 && args[1].Inspect() == "true" {
					fmt.Println(strToPrint)
				} else {
					fmt.Print(strToPrint)
				}
			}
			reader := bufio.NewReader(os.Stdin)
			var input []string
			for {
				line, err := reader.ReadString('\n')
				if err != nil {
					return NULL, true, errors.New("error reading input: " + err.Error())
				}
				line = strings.TrimSpace(line)
				if line == "" {
					break
				}
				input = append(input, line)
			}
			res := strings.Join(input, "\n")
			return &object.String{Value: "\"" + res + "\""}, false, nil
		},
	},
	"push": {
		Fn: func(args ...object.Object) (object.Object, bool, error) {
			switch arg := args[0].(type) {
			case *object.Array:
				if len(args) != 2 {
					return NULL, true, errors.New("wrong number of arguments for `push` for array, got: " + strconv.Itoa(len(args)) + ", want=2. `push(array, element)`")
				}
				arrayType := arg.TypeOf
				switch arrayType {
				case "int":
					arg2, ok := args[1].(*object.Integer)
					if !ok {
						return NULL, true, errors.New("argument type mismatch for `push`, expected element to be int, got: " + string(args[1].Type()))
					}
					arg.Elements = append(arg.Elements, arg2)
					return arg, false, nil
				case "float":
					arg2, ok := args[1].(*object.Float)
					if !ok {
						return NULL, true, errors.New("argument type mismatch for `push`, expected element to be float, got: " + string(args[1].Type()))
					}
					arg.Elements = append(arg.Elements, arg2)
					return arg, false, nil
				case "bool":
					arg2, ok := args[1].(*object.Boolean)
					if !ok {
						return NULL, true, errors.New("argument type mismatch for `push`, expected element to be bool, got: " + string(args[1].Type()))
					}
					arg.Elements = append(arg.Elements, arg2)
					return arg, false, nil
				case "char":
					arg2, ok := args[1].(*object.Char)
					if !ok {
						return NULL, true, errors.New("argument type mismatch for `push`, expected element to be char, got: " + string(args[1].Type()))
					}
					arg.Elements = append(arg.Elements, arg2)
					return arg, false, nil
				case "string":
					arg2, ok := args[1].(*object.String)
					if !ok {
						return NULL, true, errors.New("argument type mismatch for `push`, expected element to be string, got: " + string(args[1].Type()))
					}
					arg.Elements = append(arg.Elements, arg2)
					return arg, false, nil
				default:
					return NULL, true, errors.New("array type not supported for `push`, got: " + arrayType + ", want: `int`, `float`, `string`, `char` or `bool`")
				}
			case *object.Hash:
				if len(args) != 3 {
					return NULL, true, errors.New("wrong number of arguments for `push` for hashmap, got: " + strconv.Itoa(len(args)) + ", want: 3. `push(map, key, value)`")
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
						return NULL, true, errors.New("key type mismatch for `push`, expected key to be int, got: " + string(args[1].Type()))
					}
				case "float":
					arg2, ok = args[1].(*object.Float)
					if !ok {
						return NULL, true, errors.New("key type mismatch for `push`, expected key to be float, got: " + string(args[1].Type()))
					}
				case "bool":
					arg2, ok = args[1].(*object.Boolean)
					if !ok {
						return NULL, true, errors.New("key type mismatch for `push`, expected key to be bool, got: " + string(args[1].Type()))
					}
				case "char":
					arg2, ok = args[1].(*object.Char)
					if !ok {
						return NULL, true, errors.New("key type mismatch for `push`, expected key to be char, got: " + string(args[1].Type()))
					}
				case "string":
					arg2, ok = args[1].(*object.String)
					if !ok {
						return NULL, true, errors.New("key type mismatch for `push`, expected key to be string, got: " + string(args[1].Type()))
					}
				default:
					return NULL, true, errors.New("key type not supported, got: " + keyType + ", want: `int`, `float`, `string`, `char` or `bool`")
				}

				switch valueType {
				case "int":
					arg3, ok = args[2].(*object.Integer)
					if !ok {
						return NULL, true, errors.New("value type mismatch for `push`, expected value to be int, got: " + string(args[2].Type()))
					}
				case "float":
					arg3, ok = args[2].(*object.Float)
					if !ok {
						return NULL, true, errors.New("value type mismatch for `push`, expected value to be float, got: " + string(args[2].Type()))
					}
				case "bool":
					arg3, ok = args[2].(*object.Boolean)
					if !ok {
						return NULL, true, errors.New("value type mismatch for `push`, expected value to be bool, got: " + string(args[2].Type()))
					}
				case "char":
					arg3, ok = args[2].(*object.Char)
					if !ok {
						return NULL, true, errors.New("value type mismatch for `push`, expected value to be char, got: " + string(args[2].Type()))
					}
				case "string":
					arg3, ok = args[2].(*object.String)
					if !ok {
						return NULL, true, errors.New("value type mismatch for `push`, expected value to be string, got: " + string(args[2].Type()))
					}
				default:
					return NULL, true, errors.New("value type not supported, got: " + valueType + ", want: `int`, `float`, `string`, `char` or `bool`")
				}

				hashKey, ok := arg2.(object.Hashable)
				if !ok {
					return NULL, true, errors.New("key type `" + keyType + "` not hashable")
				}
				hashed := hashKey.HashKey()
				arg.Pairs[hashed] = object.HashPair{Key: arg2, Value: arg3}
				return arg, false, nil
			default:
				return NULL, true, errors.New("data structure not supported by `push`, got: " + string(args[0].Type()) + ", want: array or hashmap")
			}
		},
	},
	"pop": {
		Fn: func(args ...object.Object) (object.Object, bool, error) {
			switch arg := args[0].(type) {
			case *object.Array:
				if len(args) != 1 && len(args) != 2 {
					return NULL, true, errors.New("wrong number of arguments for `pop` for array, got: " + strconv.Itoa(len(args)) + ", want: 1 or 2. `pop(array)` or `pop(array, index)`")
				}
				if len(args) == 1 {
					if len(arg.Elements) == 0 {
						return NULL, true, errors.New("array is empty, can't pop any elements")
					}
					popped := arg.Elements[len(arg.Elements)-1]
					arg.Elements = arg.Elements[:len(arg.Elements)-1]
					return popped, false, nil
				} else {
					index, ok := args[1].(*object.Integer)
					if !ok {
						return NULL, true, errors.New("index must be an integer for `pop`, got: " + string(args[1].Type()))
					}
					if index.Value < 0 || index.Value >= int64(len(arg.Elements)) {
						return NULL, true, errors.New("index out of bounds")
					}
					popped := arg.Elements[index.Value]
					arg.Elements = append(arg.Elements[:index.Value], arg.Elements[index.Value+1:]...)
					return popped, false, nil
				}
			default:
				return NULL, true, errors.New("data structure not supported by `pop`, got: " + string(args[0].Type()) + ", want: array")
			}
		},
	},
	"insert": {
		Fn: func(args ...object.Object) (object.Object, bool, error) {
			switch arg := args[0].(type) {
			case *object.Array:
				if len(args) != 3 {
					return NULL, true, errors.New("wrong number of arguments for `insert` for array, got: " + strconv.Itoa(len(args)) + ", want: 3. `insert(array, index, element)`")
				}
				index, ok := args[1].(*object.Integer)
				if !ok {
					return NULL, true, errors.New("index must be an integer for `insert`, got: " + string(args[1].Type()))
				}
				if index.Value < 0 || index.Value >= int64(len(arg.Elements)) {
					return NULL, true, errors.New("index out of bounds")
				}
				arrayType := arg.TypeOf
				var arg2 object.Object
				switch arrayType {
				case "int":
					arg2, ok = args[2].(*object.Integer)
					if !ok {
						return NULL, true, errors.New("argument type mismatch for `insert`, expected element to be int, got: " + string(args[2].Type()))
					}
				case "float":
					arg2, ok = args[2].(*object.Float)
					if !ok {
						return NULL, true, errors.New("argument type mismatch for `insert`, expected element to be float, got: " + string(args[2].Type()))
					}
				case "bool":
					arg2, ok = args[2].(*object.Boolean)
					if !ok {
						return NULL, true, errors.New("argument type mismatch for `insert`, expected element to be bool, got: " + string(args[2].Type()))
					}
				case "char":
					arg2, ok = args[2].(*object.Char)
					if !ok {
						return NULL, true, errors.New("argument type mismatch for `insert`, expected element to be char, got: " + string(args[2].Type()))
					}
				case "string":
					arg2, ok = args[2].(*object.String)
					if !ok {
						return NULL, true, errors.New("argument type mismatch for `insert`, expected element to be string, got: " + string(args[2].Type()))
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
				return NULL, true, errors.New("data structure not supported by `insert`, got: " + string(args[0].Type()) + ", want: array")
			}
		},
	},
	"remove": {
		Fn: func(args ...object.Object) (object.Object, bool, error) {
			switch arg := args[0].(type) {
			case *object.Array:
				if len(args) != 2 {
					return NULL, true, errors.New("wrong number of arguments for `remove` for array, got: " + strconv.Itoa(len(args)) + ", want: 2. `remove(array, element)`")
				}
				arrayType := arg.TypeOf
				switch arrayType {
				case "int":
					arg2, ok := args[1].(*object.Integer)
					if !ok {
						return NULL, true, errors.New("argument type mismatch for `remove`, expected element to be int, got: " + string(args[1].Type()))
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
						return NULL, true, errors.New("argument type mismatch for `remove`, expected element to be float, got: " + string(args[1].Type()))
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
						return NULL, true, errors.New("argument type mismatch for `remove`, expected element to be bool, got: " + string(args[1].Type()))
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
						return NULL, true, errors.New("argument type mismatch for `remove`, expected element to be char, got: " + string(args[1].Type()))
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
						return NULL, true, errors.New("argument type mismatch for `remove`, expected element to be string, got: " + string(args[1].Type()))
					}
					for i, element := range arg.Elements {
						if element.(*object.String).Value == arg2.Value {
							arg.Elements = append(arg.Elements[:i], arg.Elements[i+1:]...)
							return arg, false, nil
						}
					}
					return NULL, false, nil
				default:
					return NULL, true, errors.New("array type not supported for `remove`, got: " + arrayType + ", want: `int`, `float`, `string`, `char` or `bool`")
				}
			case *object.Hash:
				if len(args) != 2 {
					return NULL, true, errors.New("wrong number of arguments for `remove` for hashmap, got: " + strconv.Itoa(len(args)) + ", want: 2. `remove(map, key)`")
				}
				keyType := arg.KeyType
				var arg2 object.Object
				var ok bool

				switch keyType {
				case "int":
					arg2, ok = args[1].(*object.Integer)
					if !ok {
						return NULL, true, errors.New("key type mismatch for `remove`, expected key to be int, got: " + string(args[1].Type()))
					}
				case "float":
					arg2, ok = args[1].(*object.Float)
					if !ok {
						return NULL, true, errors.New("key type mismatch for `remove`, expected key to be float, got: " + string(args[1].Type()))
					}
				case "bool":
					arg2, ok = args[1].(*object.Boolean)
					if !ok {
						return NULL, true, errors.New("key type mismatch for `remove`, expected key to be bool, got: " + string(args[1].Type()))
					}
				case "char":
					arg2, ok = args[1].(*object.Char)
					if !ok {
						return NULL, true, errors.New("key type mismatch for `remove`, expected key to be char, got: " + string(args[1].Type()))
					}
				case "string":
					arg2, ok = args[1].(*object.String)
					if !ok {
						return NULL, true, errors.New("key type mismatch for `remove`, expected key to be string, got: " + string(args[1].Type()))
					}
				default:
					return NULL, true, errors.New("key type not supported, got: " + keyType + ", want: `int`, `float`, `string`, `char` or `bool`")
				}

				hashKey, ok := arg2.(object.Hashable)
				if !ok {
					return NULL, true, errors.New("key type not hashable")
				}
				hashed := hashKey.HashKey()
				value, ok := arg.Pairs[hashed]
				if !ok {
					return NULL, true, errors.New("key not found. can't remove pair that doesn't exist")
				}
				delete(arg.Pairs, hashed)
				return value.Value, false, nil
			default:
				return NULL, true, errors.New("data structure not supported by `remove`, got: " + string(args[0].Type()) + ", want: array or hashmap")
			}
		},
	},
	"getIndex": {
		Fn: func(args ...object.Object) (object.Object, bool, error) {
			switch arg := args[0].(type) {
			case *object.Array:
				if len(args) != 2 {
					return NULL, true, errors.New("wrong number of arguments for `getIndex` for array, got: " + strconv.Itoa(len(args)) + ", want: 2. `getIndex(array, element)`")
				}
				arrayType := arg.TypeOf
				switch arrayType {
				case "int":
					arg2, ok := args[1].(*object.Integer)
					if !ok {
						return NULL, true, errors.New("argument type mismatch for `getIndex`, expected element to be int, got: " + string(args[1].Type()))
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
						return NULL, true, errors.New("argument type mismatch for `getIndex`, expected element to be float, got: " + string(args[1].Type()))
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
						return NULL, true, errors.New("argument type mismatch for `getIndex`, expected element to be bool, got: " + string(args[1].Type()))
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
						return NULL, true, errors.New("argument type mismatch for `getIndex`, expected element to be char, got: " + string(args[1].Type()))
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
						return NULL, true, errors.New("argument type mismatch for `getIndex`, expected element to be string, got: " + string(args[1].Type()))
					}
					for i, element := range arg.Elements {
						if element.(*object.String).Value == arg2.Value {
							return &object.Integer{Value: int64(i)}, false, nil
						}
					}
					return &object.Integer{Value: -1}, false, nil
				default:
					return NULL, true, errors.New("array type not supported for `getIndex`, got: " + arrayType + ", want: `int`, `float`, `string`, `char` or `bool`")
				}
			default:
				return NULL, true, errors.New("data structure not supported by `getIndex`, got: " + string(args[0].Type()) + ", want: array")
			}
		},
	},
	"keys": {
		Fn: func(args ...object.Object) (object.Object, bool, error) {
			switch arg := args[0].(type) {
			case *object.Hash:
				if len(args) != 1 {
					return NULL, true, errors.New("wrong number of arguments for `keys` for hashmap, got: " + strconv.Itoa(len(args)) + ", want: 1. `keys(map)`")
				}
				var keys []object.Object
				for _, pair := range arg.Pairs {
					keys = append(keys, pair.Key)
				}
				return &object.Array{Elements: keys, TypeOf: arg.KeyType}, false, nil
			default:
				return NULL, true, errors.New("data structure not supported by `keys`, got: " + string(args[0].Type()) + ", want: hashmap")
			}
		},
	},
	"values": {
		Fn: func(args ...object.Object) (object.Object, bool, error) {
			switch arg := args[0].(type) {
			case *object.Hash:
				if len(args) != 1 {
					return NULL, true, errors.New("wrong number of arguments for `values` for hashmap, got: " + strconv.Itoa(len(args)) + ", want: 1. `values(map)`")
				}
				var values []object.Object
				for _, pair := range arg.Pairs {
					values = append(values, pair.Value)
				}
				return &object.Array{Elements: values, TypeOf: arg.ValueType}, false, nil
			default:
				return NULL, true, errors.New("data structure not supported by `values`, got: " + string(args[0].Type()) + ", want: hashmap")
			}
		},
	},
	"containsKey": {
		Fn: func(args ...object.Object) (object.Object, bool, error) {
			switch arg := args[0].(type) {
			case *object.Hash:
				if len(args) != 2 {
					return NULL, true, errors.New("wrong number of arguments for `containsKey` for hashmap, got: " + strconv.Itoa(len(args)) + ", want: 2. `containsKey(map, key)`")
				}
				keyType := arg.KeyType
				var arg2 object.Object
				var ok bool
				switch keyType {
				case "int":
					arg2, ok = args[1].(*object.Integer)
					if !ok {
						return NULL, true, errors.New("key type mismatch for `containsKey`, expected key to be int, got: " + string(args[1].Type()))
					}
				case "float":
					arg2, ok = args[1].(*object.Float)
					if !ok {
						return NULL, true, errors.New("key type mismatch for `containsKey`, expected key to be float, got: " + string(args[1].Type()))
					}
				case "bool":
					arg2, ok = args[1].(*object.Boolean)
					if !ok {
						return NULL, true, errors.New("key type mismatch for `containsKey`, expected key to be bool, got: " + string(args[1].Type()))
					}
				case "char":
					arg2, ok = args[1].(*object.Char)
					if !ok {
						return NULL, true, errors.New("key type mismatch for `containsKey`, expected key to be char, got: " + string(args[1].Type()))
					}
				case "string":
					arg2, ok = args[1].(*object.String)
					if !ok {
						return NULL, true, errors.New("key type mismatch for `containsKey`, expected key to be string, got: " + string(args[1].Type()))
					}
				default:
					return NULL, true, errors.New("key type not supported, got: " + keyType + ", want: `int`, `float`, `string`, `char` or `bool`")
				}
				hashKey, ok := arg2.(object.Hashable)
				if !ok {
					return NULL, true, errors.New("key type not hashable")
				}
				hashed := hashKey.HashKey()
				_, ok = arg.Pairs[hashed]
				if ok {
					return TRUE, false, nil
				}
				return FALSE, false, nil
			default:
				return NULL, true, errors.New("data structure not supported by `containsKey`, got: " + string(args[0].Type()) + ", want: hashmap")
			}
		},
	},
	"typeOf": {
		Fn: func(args ...object.Object) (object.Object, bool, error) {
			if len(args) != 1 {
				return NULL, true, errors.New("wrong number of arguments for `typeOf`, got: " + strconv.Itoa(len(args)) + ", want: 1")
			}
			switch arg := args[0].(type) {
			case *object.Array:
				res := "\"" + arg.TypeOf + "[]\""
				return &object.String{Value: res}, false, nil
			case *object.Hash:
				res := "\"" + arg.KeyType + "[" + arg.ValueType + "]\""
				return &object.String{Value: res}, false, nil
			case *object.Integer:
				return &object.String{Value: "\"int\""}, false, nil
			case *object.Float:
				return &object.String{Value: "\"float\""}, false, nil
			case *object.Boolean:
				return &object.String{Value: "\"bool\""}, false, nil
			case *object.String:
				return &object.String{Value: "\"string\""}, false, nil
			case *object.Char:
				return &object.String{Value: "\"char\""}, false, nil
			default:
				return NULL, true, errors.New("data structure not supported by `typeOf`, got: " + string(args[0].Type()))
			}
		},
	},
	"slice": {
		Fn: func(args ...object.Object) (object.Object, bool, error) {
			if len(args) != 3 && len(args) != 4 {
				return NULL, true, errors.New("wrong number of arguments for `slice`, got: " + strconv.Itoa(len(args)) + ", want: 3 or 4. `slice(array/string, start, end)` or `slice(array/string, start, end, step)`")
			}
			switch arg := args[0].(type) {
			case *object.Array:
				start, ok := args[1].(*object.Integer)
				if !ok {
					return NULL, true, errors.New("start index must be an `int` for `slice`, got: " + string(args[1].Type()))
				}
				end, ok := args[2].(*object.Integer)
				if !ok {
					return NULL, true, errors.New("end index must be an `int` for `slice`, got: " + string(args[2].Type()))
				}
				if start.Value < 0 || start.Value >= int64(len(arg.Elements)) || end.Value < 0 || end.Value > int64(len(arg.Elements)) {
					return NULL, true, errors.New("index out of bounds for `slice` operation")
				}
				if start.Value > end.Value {
					return NULL, true, errors.New("start index must be less than end index for `slice` operation")
				}
				if len(args) == 3 {
					newArr := &object.Array{TypeOf: arg.TypeOf}
					slicedElement := arg.Elements[start.Value:end.Value]
					newArr.Elements = slicedElement
					return newArr, false, nil
				} else {
					step, ok := args[3].(*object.Integer)
					if !ok {
						return NULL, true, errors.New("step must be an `int` for `slice`, got: " + string(args[3].Type()))
					}
					if step.Value == 0 {
						return NULL, true, errors.New("step cannot be zero for `slice` operation")
					}
					if step.Value < 0 {
						return NULL, true, errors.New("step cannot be negative for `slice` operation")
					}
					newArr := &object.Array{TypeOf: arg.TypeOf}
					var slicedElement []object.Object
					for i := start.Value; i < end.Value; i += step.Value {
						slicedElement = append(slicedElement, arg.Elements[i])
					}
					newArr.Elements = slicedElement
					return newArr, false, nil
				}
			case *object.String:
				start, ok := args[1].(*object.Integer)
				if !ok {
					return NULL, true, errors.New("start index must be an `int` for `slice`, got: " + string(args[1].Type()))
				}
				end, ok := args[2].(*object.Integer)
				if !ok {
					return NULL, true, errors.New("end index must be an `int` for `slice`, got: " + string(args[2].Type()))
				}
				if start.Value < 0 || start.Value >= int64(len(arg.Value)) || end.Value < 0 || end.Value > int64(len(arg.Value)) {
					return NULL, true, errors.New("index out of bounds for `slice` operation")
				}
				if start.Value > end.Value {
					return NULL, true, errors.New("start index must be less than end index for `slice` operation")
				}
				if len(args) == 3 {
					oldStr := arg.Value
					oldStr = oldStr[1 : len(oldStr)-1]
					newStr := oldStr[start.Value:end.Value]
					newStr = "\"" + newStr + "\""
					return &object.String{Value: newStr}, false, nil
				} else {
					step, ok := args[3].(*object.Integer)
					if !ok {
						return NULL, true, errors.New("step must be an `int` for `slice`, got: " + string(args[3].Type()))
					}
					if step.Value == 0 {
						return NULL, true, errors.New("step cannot be zero for `slice` operation")
					}
					if step.Value < 0 {
						return NULL, true, errors.New("step cannot be negative for `slice` operation")
					}
					oldStr := arg.Value
					oldStr = oldStr[1 : len(oldStr)-1]
					newStr := ""
					for i := start.Value; i < end.Value; i += step.Value {
						newStr += string(oldStr[i])
					}
					newStr = "\"" + newStr + "\""
					return &object.String{Value: newStr}, false, nil
				}
			default:
				return NULL, true, errors.New("data structure not supported by `slice`, got: " + string(args[0].Type()) + ", want: array or string")
			}
		},
	},
}
