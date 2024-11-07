package evaluator

import (
	"errors"
	"fmt"

	"github.com/KhushPatibandha/Kolon/src/ast"
	"github.com/KhushPatibandha/Kolon/src/object"
)

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

func Eval(node ast.Node) (object.Object, bool, error) {
	switch node := node.(type) {
	case *ast.Program:
		return evalStatements(node.Statements)
	case *ast.ExpressionStatement:
		return Eval(node.Expression)
	case *ast.IntegerValue:
		return &object.Integer{Value: node.Value}, false, nil
	case *ast.FloatValue:
		return &object.Float{Value: node.Value}, false, nil
	case *ast.StringValue:
		return &object.String{Value: node.Value}, false, nil
	case *ast.CharValue:
		return &object.Char{Value: node.Value}, false, nil
	case *ast.BooleanValue:
		if node.Value {
			return TRUE, false, nil
		} else {
			return FALSE, false, nil
		}
	case *ast.PrefixExpression:
		if postfix, ok := node.Right.(*ast.PostfixExpression); ok {
			left, hasErr, err := Eval(postfix.Left)
			if err != nil {
				return left, hasErr, err
			}
			operator := postfix.Operator
			res, hasErr, err := evalPostfixExpression(operator, left)
			if err != nil {
				return res, hasErr, err
			}

			var resVal interface{}

			switch {
			case res.Type() == object.INTEGER_OBJ:
				resVal = res.(*object.Integer).Value
			case res.Type() == object.FLOAT_OBJ:
				resVal = res.(*object.Float).Value
			default:
				return NULL, true, errors.New("Only Integer and Float datatypes supported with Postfix operation. got: " + string(res.Type()))
			}

			switch val := resVal.(type) {
			case int64:
				if operator == "++" {
					return evalPrefixExpression(node.Operator, &object.Integer{Value: val - 2})
				} else {
					return evalPrefixExpression(node.Operator, &object.Integer{Value: val + 2})
				}
			case float64:
				if operator == "++" {
					return evalPrefixExpression(node.Operator, &object.Float{Value: val - 2})
				} else {
					return evalPrefixExpression(node.Operator, &object.Float{Value: val + 2})
				}
			}
		}

		right, hasErr, err := Eval(node.Right)
		if err != nil {
			return right, hasErr, err
		}
		return evalPrefixExpression(node.Operator, right)
	case *ast.InfixExpression:
		left, hasErr, err := Eval(node.Left)
		if err != nil {
			return left, hasErr, err
		}
		right, hasErr, err := Eval(node.Right)
		if err != nil {
			return right, hasErr, err
		}
		return evalInfixExpression(node.Operator, left, right)
	case *ast.PostfixExpression:
		left, hasErr, err := Eval(node.Left)
		if err != nil {
			return left, hasErr, err
		}
		return evalPostfixExpression(node.Operator, left)
	case *ast.FunctionBody:
		return evalStatements(node.Statements)
	case *ast.IfStatement:
		return evalIfStatements(node)
	case *ast.ReturnStatement:
		var val []object.Object

		for i := 0; i < len(node.Value); i++ {
			rsObj, hasErr, err := evalReturnValue(node, i)
			if err != nil {
				return NULL, hasErr, err
			}
			val = append(val, rsObj)
		}

		return &object.ReturnValue{Value: val}, false, nil
	}
	return nil, true, errors.New("No Eval function for given node type. got: " + string(node.String()))
}

func evalStatements(stmts []ast.Statement) (object.Object, bool, error) {
	var result object.Object
	var hasErr bool
	var err error
	for _, statement := range stmts {
		result, hasErr, err = Eval(statement)

		if result != nil && result.Type() == object.RETURN_VALUE_OBJ {
			return result, hasErr, err
		}
	}
	return result, hasErr, err
}

func evalReturnValue(rs *ast.ReturnStatement, idx int) (object.Object, bool, error) {
	currNode := rs.Value[idx]
	switch currNode.(type) {
	case *ast.IntegerValue, *ast.FloatValue, *ast.BooleanValue, *ast.StringValue, *ast.CharValue, *ast.PrefixExpression, *ast.PostfixExpression, *ast.InfixExpression:
		return Eval(currNode)
	default:
		return NULL, true, errors.New("Can Only return expressions and datatypes. got: " + fmt.Sprintf("%T", currNode))
	}
}

func evalIfStatements(node *ast.IfStatement) (object.Object, bool, error) {
	condition, hasErr, err := Eval(node.Value)
	if err != nil {
		return condition, hasErr, err
	}

	conditionRes, hasErr, err := execIf(condition)
	if err != nil {
		return NULL, hasErr, err
	}

	if conditionRes {
		return Eval(node.Body)
	} else if node.MultiConseq != nil {
		for i := 0; i < len(node.MultiConseq); i++ {
			resObj, hasErr, err := evalElseIfStatement(node.MultiConseq[i])
			if err != nil {
				return resObj, hasErr, err
			}
			if resObj != nil {
				return resObj, hasErr, err
			}
		}
	}

	if node.Consequence != nil {
		return evalStatements(node.Consequence.Body.Statements)
	}
	return nil, false, nil
}

func evalElseIfStatement(node *ast.ElseIfStatement) (object.Object, bool, error) {
	condition, hasErr, err := Eval(node.Value)
	if err != nil {
		return condition, hasErr, err
	}
	conditionRes, hasErr, err := execIf(condition)
	if err != nil {
		return NULL, hasErr, err
	}
	if conditionRes {
		return Eval(node.Body)
	}
	return nil, false, nil
}

func execIf(obj object.Object) (bool, bool, error) {
	switch obj {
	case TRUE:
		return true, false, nil
	case FALSE:
		return false, false, nil
	default:
		return false, true, errors.New("Conditions for 'if' and 'else if' statements must result in a boolean. got: " + string(obj.Type()))
	}
}

func evalPrefixExpression(operator string, right object.Object) (object.Object, bool, error) {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusOperatorExpression(right)
	default:
		return NULL, true, errors.New("Only 2 Prefix Operator's supported. !(Bang) and -(Dash/Minus). got: " + operator)
	}
}

func evalMinusOperatorExpression(right object.Object) (object.Object, bool, error) {
	if right.Type() != object.INTEGER_OBJ && right.Type() != object.FLOAT_OBJ {
		return NULL, true, errors.New("Dash/Minus(-) operator can be only used with Integers(int) and Floats(float) entities. got: " + string(right.Type()))
	}

	if right.Type() == object.FLOAT_OBJ {
		value := right.(*object.Float).Value
		return &object.Float{Value: -value}, false, nil
	}
	value := right.(*object.Integer).Value
	return &object.Integer{Value: -value}, false, nil
}

func evalBangOperatorExpression(right object.Object) (object.Object, bool, error) {
	switch right := right.(type) {
	case *object.Boolean:
		if right == TRUE {
			return FALSE, false, nil
		} else {
			return TRUE, false, nil
		}
	default:
		return NULL, true, errors.New("Bang Operator(!) can be only used with boolean(bool) entities. got: " + string(right.Type()))
	}
}

func evalPostfixExpression(operator string, left object.Object) (object.Object, bool, error) {
	switch {
	case left.Type() == object.INTEGER_OBJ:
		return evalIntegerPostfixExpression(operator, left)
	case left.Type() == object.FLOAT_OBJ:
		return evalFloatPostfixExpression(operator, left)
	default:
		return NULL, true, errors.New("Only Integer and Float datatypes supported with Postfix operation. got: " + string(left.Type()))
	}
}

func evalFloatPostfixExpression(operator string, left object.Object) (object.Object, bool, error) {
	leftVal := left.(*object.Float).Value

	switch operator {
	case "++":
		return &object.Float{Value: leftVal + 1}, false, nil
	case "--":
		return &object.Float{Value: leftVal - 1}, false, nil
	default:
		return NULL, true, errors.New("Only ++ and -- Postfix operation supported. got: " + operator)
	}
}

func evalIntegerPostfixExpression(operator string, left object.Object) (object.Object, bool, error) {
	leftVal := left.(*object.Integer).Value

	switch operator {
	case "++":
		return &object.Integer{Value: leftVal + 1}, false, nil
	case "--":
		return &object.Integer{Value: leftVal - 1}, false, nil
	default:
		return NULL, true, errors.New("Only ++ and -- Postfix operation supported. got: " + operator)
	}
}

func evalInfixExpression(operator string, left object.Object, right object.Object) (object.Object, bool, error) {
	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(operator, left, right)
	case left.Type() == object.BOOLEAN_OBJ && right.Type() == object.BOOLEAN_OBJ:
		return evalBooleanInfixExpression(operator, left, right)
	case left.Type() == object.FLOAT_OBJ && right.Type() == object.FLOAT_OBJ:
		return evalFloatInfixExpression(operator, left, right)
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.FLOAT_OBJ || left.Type() == object.FLOAT_OBJ && right.Type() == object.INTEGER_OBJ:
		leftVal := 0.0
		rightVal := 0.0
		if left.Type() == object.INTEGER_OBJ {
			leftVal = float64(left.(*object.Integer).Value)
			rightVal = right.(*object.Float).Value
		} else {
			leftVal = left.(*object.Float).Value
			rightVal = float64(right.(*object.Integer).Value)
		}
		return evalFloatInfixExpression(operator, &object.Float{Value: leftVal}, &object.Float{Value: rightVal})
	case left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ:
		return evalStringInfixExpression(operator, left, right)
	case left.Type() == object.CHAR_OBJ && right.Type() == object.CHAR_OBJ:
		return evalCharInfixExpression(operator, left, right)
	default:
		return NULL, true, errors.New("Invalid operation with variable types on left and right.")
	}
}

func evalCharInfixExpression(operator string, left object.Object, right object.Object) (object.Object, bool, error) {
	leftVal := left.(*object.Char).Value
	rightVal := right.(*object.Char).Value
	leftVal = leftVal[1 : len(leftVal)-1]
	rightVal = rightVal[1 : len(rightVal)-1]

	switch operator {
	case "+":
		return &object.String{Value: "\"" + leftVal + rightVal + "\""}, false, nil
	case "==":
		if leftVal == rightVal {
			return TRUE, false, nil
		}
		return FALSE, false, nil
	case "!=":
		if leftVal != rightVal {
			return TRUE, false, nil
		}
		return FALSE, false, nil
	default:
		return NULL, true, errors.New("Can only perform \"+\", \"==\", \"!=\" with chars. got: " + operator)
	}
}

func evalStringInfixExpression(operator string, left object.Object, right object.Object) (object.Object, bool, error) {
	leftVal := left.(*object.String).Value
	rightVal := right.(*object.String).Value
	leftVal = leftVal[1 : len(leftVal)-1]
	rightVal = rightVal[1 : len(rightVal)-1]

	switch operator {
	case "+":
		return &object.String{Value: "\"" + leftVal + rightVal + "\""}, false, nil
	case "==":
		if leftVal == rightVal {
			return TRUE, false, nil
		}
		return FALSE, false, nil
	case "!=":
		if leftVal != rightVal {
			return TRUE, false, nil
		}
		return FALSE, false, nil
	default:
		return NULL, true, errors.New("Can only perform \"+\", \"==\", \"!=\" with strings. got: " + operator)
	}
}

func evalFloatInfixExpression(operator string, left object.Object, right object.Object) (object.Object, bool, error) {
	leftVal := left.(*object.Float).Value
	rightVal := right.(*object.Float).Value

	switch operator {
	case "+":
		return &object.Float{Value: leftVal + rightVal}, false, nil
	case "-":
		return &object.Float{Value: leftVal - rightVal}, false, nil
	case "/":
		return &object.Float{Value: leftVal / rightVal}, false, nil
	case "*":
		return &object.Float{Value: leftVal * rightVal}, false, nil
	case "%":
		return &object.Float{Value: float64(int(leftVal) % int(rightVal))}, false, nil
	case ">":
		if leftVal > rightVal {
			return TRUE, false, nil
		}
		return FALSE, false, nil
	case "<":
		if leftVal < rightVal {
			return TRUE, false, nil
		}
		return FALSE, false, nil
	case "==":
		if leftVal == rightVal {
			return TRUE, false, nil
		}
		return FALSE, false, nil
	case "!=":
		if leftVal != rightVal {
			return TRUE, false, nil
		}
		return FALSE, false, nil
	case "<=":
		if leftVal <= rightVal {
			return TRUE, false, nil
		}
		return FALSE, false, nil
	case ">=":
		if leftVal >= rightVal {
			return TRUE, false, nil
		}
		return FALSE, false, nil
	default:
		return NULL, true, errors.New("Can only perform \"+\", \"-\", \"/\", \"*\", \"%\", \">\", \"<\", \"<=\", \">=\", \"!=\", \"==\" with 2 Float var. got: " + operator)
	}
}

func evalIntegerInfixExpression(operator string, left object.Object, right object.Object) (object.Object, bool, error) {
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value

	switch operator {
	case "+":
		return &object.Integer{Value: leftVal + rightVal}, false, nil
	case "-":
		return &object.Integer{Value: leftVal - rightVal}, false, nil
	case "/":
		return &object.Integer{Value: leftVal / rightVal}, false, nil
	case "*":
		return &object.Integer{Value: leftVal * rightVal}, false, nil
	case "%":
		return &object.Integer{Value: leftVal % rightVal}, false, nil
	case "&":
		return &object.Integer{Value: leftVal & rightVal}, false, nil
	case "|":
		return &object.Integer{Value: leftVal | rightVal}, false, nil
	case ">":
		if leftVal > rightVal {
			return TRUE, false, nil
		}
		return FALSE, false, nil
	case "<":
		if leftVal < rightVal {
			return TRUE, false, nil
		}
		return FALSE, false, nil
	case "==":
		if leftVal == rightVal {
			return TRUE, false, nil
		}
		return FALSE, false, nil
	case "!=":
		if leftVal != rightVal {
			return TRUE, false, nil
		}
		return FALSE, false, nil
	case "<=":
		if leftVal <= rightVal {
			return TRUE, false, nil
		}
		return FALSE, false, nil
	case ">=":
		if leftVal >= rightVal {
			return TRUE, false, nil
		}
		return FALSE, false, nil
	default:
		return NULL, true, errors.New("Can only perform \"+\", \"-\", \"/\", \"*\", \"%\", \"&\", \"|\", \">\", \"<\", \"<=\", \">=\", \"!=\", \"==\" with 2 Integers. got: " + operator)
	}
}

func evalBooleanInfixExpression(operator string, left object.Object, right object.Object) (object.Object, bool, error) {
	leftVal := left.(*object.Boolean).Value
	rightVal := right.(*object.Boolean).Value

	switch operator {
	case "==":
		if leftVal == rightVal {
			return TRUE, false, nil
		}
		return FALSE, false, nil
	case "!=":
		if leftVal != rightVal {
			return TRUE, false, nil
		}
		return FALSE, false, nil
	case "&&":
		if leftVal && rightVal {
			return TRUE, false, nil
		}
		return FALSE, false, nil
	case "||":
		if leftVal || rightVal {
			return TRUE, false, nil
		}
		return FALSE, false, nil
	default:
		return NULL, true, errors.New("Can only perform \"&&\", \"||\", \"!=\", \"==\" with 2 Boolean values. got: " + operator)
	}
}
