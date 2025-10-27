package evaluator

import "github.com/KhushPatibandha/Kolon/src/object"

func deepCopy(o object.Object) object.Object {
	switch obj := o.(type) {
	case *object.Integer:
		return &object.Integer{Value: obj.Value}
	case *object.Float:
		return &object.Float{Value: obj.Value}
	case *object.Bool:
		return &object.Bool{Value: obj.Value}
	case *object.String:
		return &object.String{Value: obj.Value}
	case *object.Char:
		return &object.Char{Value: obj.Value}
	case *object.Array:
		copyEle := make([]object.Object, len(obj.Elements))
		for i, ele := range obj.Elements {
			copyEle[i] = deepCopy(ele)
		}
		return &object.Array{Elements: copyEle}
	default:
		newPairs := make(map[object.HashKey]object.HashPair)
		for k, v := range obj.(*object.HashMap).Pairs {
			newPairs[k] = object.HashPair{
				Key:   deepCopy(v.Key),
				Value: deepCopy(v.Value),
			}
		}
		return &object.HashMap{Pairs: newPairs}
	}
}
