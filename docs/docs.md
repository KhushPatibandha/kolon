# Documentation

## Introduction

I, created this language to learn more about programming languages and how they work under the hood. I have tried my best to make it simple to understand and write, while also being unique in its own way.

You might have already noticed—or will notice later when going through the docs—the extensive use of `kolon (:)`. That's about it! I think that's unique enough, right?

Before we start, you should have Kolon installed on your system, if you don't already. You can do so by building it from the source or installing the binary. Check out [installation guide](https://github.com/KhushPatibandha/Kolon/blob/main/docs/installation.md) for more details.

## Hello, World!!

Kolon is a functional language, and hence everything is organized into functions. In Kolon, the entry point of a program is, you guessed it, `main()`. Here’s how you can write your very first program in Kolon:

```kolon
fun: main() {
    println("Hello World!!");
}
```

## Running Kolon

You can run code written in Kolon like this:

```
kolon run: <path-to-file>
```

## Comments

To comment a line, you can use `//`, just like in many other languages.

```kolon
fun: main() {
    // println("Hello, World!!");
}
```

## Variables and Types

Kolon is a strongly and statically typed language, and hence the type of a variable must be declared when defining the variable and that will be checked before runtime.

### Data Types

Here are the data types supported by Kolon: `int`, `float`, `string`, `char`, `bool`. Note that there is NO concept of `long` or `double`, so both `int` and `float` are `64-bit`.

### Declaring Variables

You can declare a variable using the `var` keyword:

```kolon
fun: main() {
    var someInt: int = 10;
    var someFloat: float = 1.1;
    var someString: string = "Hello, World";
    var someChar: char = 'c';
    var someBool: bool = true;
}
```

Variables declared with the `var` keyword are mutable and can be changed, but only to values of the same type as the one assigned.

### Declaring Constant Variables

Kolon also supports constant variables, which cannot be changed once declared (non-mutable). You can declare them using the `const` keyword:

```kolon
fun: main() {
    const someInt: int = 10;
    const someFloat: float = 1.1;
    const someString: string = "Hello, World";
    const someChar: char = 'c';
    const someBool: bool = true;

    const someInt: int = 100; // Error! Cannot re-declare a constant
    someInt = 100; // Error! Cannot re-assign a constant
}
```

### Default Values

You can skip defining a value when declaring a variable with the `var` keyword. In such cases, a default value is assigned:

```kolon
fun: main() {
    var someInt: int; // default = 0
    var someFloat: float; // default = 0.0
    var someString: string; // default = ""
    var someChar: char; // default = ''
    var someBool: bool; // default = false
}
```

On the other hand, a variable declared with the `const` keyword MUST be initialized with a value:

```kolon
fun: main() {
    const someInt: int = 10; // Won't throw an error
    const someOtherInt: int; // Will throw an error
}
```

### Multi-Value Assignment

You can assign values to multiple variables or identifiers of different types on the left side at the same time.

```kolon
fun: main() {
    var a: int;
    var c: string;
    a, var b: bool, c = 10, true, "hello!";
}
```

Note:

- The sequence is very important in this. The values on the right side must match the variables on the left side in the correct order.

## Data Structures

### List/Arrays

You can define an array by adding `[]` after the type in a variable. An array can only have elements of that type. You can define the elements in the array using `[]`.

Also, an empty array must be defined with `[]`.

```kolon
fun: main() {
    var a: int[] = [1, 2, 3, 4, 5];
    var b: int[]; // not valid, will throw an error
    var c: int[] = []; // Won't throw an error
}
```

#### Accessing Array Elements

You can use `[]` to access an element in an array. The index must be greater than or equal to 0 and less than the length of the array.

```kolon
fun: main() {
    var a: int[] = [1, 2, 3];
    println(a[0]); // 1
    println(a[1]); // 2
    println(a[4]); // error!!
}
```

### HashMaps

You can define a hashmap by adding `[type]` after the type in a variable. All key-value pairs in the hashmap must follow this type rule. Hashmaps are represented with `{}`.

Additionally, an empty map must be defined with `{}`.

```kolon
fun: main() {
    var a: string[int] = {"kolon": 1, "hello": 2};
    var b: string[int]; // Not valid, will throw an error
    var c: string[int] = {}; // Won't throw and error
}
```

#### Accessing Map Elements

You can use `[]` to access a value associated with a key.

```kolon
fun: main() {
    var a: string[int] = {"kolon": 1, "hello": 2};
    println(a["kolon"]) // 1
    println(a["someKey"]) // Key doesn't exist error!
}
```

## if - else if - else

Like most languages, Kolon also has conditionals. All conditions in `if` and `else if` statements must evaluate to a boolean value (`true` or `false`).

```kolon
if: (a > b): {
    var d: int = 40;
} else if: (b > c): {
    var e: int = 50;
} else: {
    var f: int = 60;
}
```

## Loops

Kolon supports two type of loops, `for` loop and `while` loop.

### For Loop

```kolon
for: (var i: int = 0; i < 3; i++): {
    println(i);
}
```

`for` loop expects three arguments:

- A `var` statement to declare and initialize the loop variable or an assignment expression in case the variable is already declared.
- An infix operation that evaluates to a boolean value (the loop condition).
- A postfix operation (e.g., increment or decrement) or an assignment expression (e.g., +=, -=, etc...) to update the loop variable.

Note:

- The first and third argument MUST always result in an integer value.
- The second argument MUST be an infix operation that results in a boolean value.

### While Loop

```kolon
var i: int = 0;
while: (i < 10): {
    println(i);
    i++;
}
```

`while` loop expects a single argument:

- Break condition that MUST always evaluates to a boolean value.

### Continue and Break

You can also use `continue` to skip to next loop iteration or `break` to break out of the loop

```kolon
for: (var i: int = 0; i < 10; i++): {
    if: (i == 5): {
        continue;
    } else if: (i == 7): {
        break;
    }
    print(i);
}
// 012346 will be printed out
```

## Functions

### Defining Functions

- The keyword to define a function is `fun`, because why not? :)
- A function has a few components: Name, Parameters, Return Types, and a Body.
- A function will only be executed if it is called from the `main` function.
- The `main` function must NOT take any parameters and must NOT return anything.

```kolon
// main function
fun: main() {
    // Statements
}

// function with Parameters and Return types
fun: someName(a: int, b: int, c: string): (int, bool) {
    // Statements
}

// function with no Parameters and Return Types
fun: SomeOtherName() {
    // Statements
}
```

As you may have noticed, you can take multiple arguments of different types in your function and return multiple values with different data types. You can do this with the help of the `return` keyword. More on that later.

### Calling Functions

You can call a function by simply writing its name, followed by parentheses, and passing parameters inside the parentheses if needed.

```kolon
fun: main() {
    var d: float = 1.1;
    var b: int;
    var a: string, b, var c: bool, var e: int[], var f: string[int] = callMe(d);
    println(a); // Kolon
    println(b); // 10
    println(c); // false
    callMeAgain(); // hello
}
fun: callMe(var1: float): (string, int, bool, int[], string[int]) {
    return: ("Kolon", 10, false, [1, 2, 3, 4, 5], {"kolon": 1, "hello": 2});
}
fun: callMeAgain() {
    println("hello")
}
```

### Builtin Functions

Kolon has many built-in functions that you can use without needing to define them. You can simply call them.

#### print()

| **Num of Args** | **Type of Args**                         | **Returns** | **Description**                                    |
| --------------- | ---------------------------------------- | ----------- | -------------------------------------------------- |
| 1               | int/float/bool/char/string/array/hashmap | -           | Prints to the console without going to a new line. |

```kolon
fun: main() {
    print(1);
    print("hello!! " + toString(1));
    print("hello!! " + toString(1.1) + " ");
    print("hehe!! " + toString(true));
    print(1.1);
    print(true);
    print('c');
    // this program will print:
    // 1hello!! 1hello!! 1.1 hehe!! true1.1truec
}
```

#### println()

| **Num of Args** | **Type of Args**                         | **Returns** | **Description**                                  |
| --------------- | ---------------------------------------- | ----------- | ------------------------------------------------ |
| 1               | int/float/bool/char/string/array/hashmap | -           | Prints to the console and goes to the next line. |

```kolon
fun: main() {
    println(1);
    println("hello!! " + toString(1));
    println("hello!! " + toString(1.1) + " ");
    println("hehe!! " + toString(true));
    println(1.1);
    println(true);
    println('c');
    // this program will print:
    // 1
    // hello!! 1
    // hello!! 1.1
    // hehe!! true
    // 1.1
    // true
    // c
}
```

#### len()

| **Num of Args** | **Type of Args**     | **Returns** | **Description**                              |
| --------------- | -------------------- | ----------- | -------------------------------------------- |
| 1               | string/array/hashmap | int         | Returns the length of the provided argument. |

```kolon
fun: main() {
    var a: string = "hello";
    var b: int[] = [1, 2, 3];
    var c: string[int] = {"kolon": 1, "hello": 2};
    println(len(a)); // 5
    println(len(b)); // 3
    println(len(c)); // 2
    println(len("hehe")) // 4
    println(len([1, 2, 3, 4, 5])) // 5
    println(len({"kolon": 1, "hello": 2})) // 2
}
```

#### toString()

| **Num of Args** | **Type of Args**                         | **Returns** | **Description**                                    |
| --------------- | ---------------------------------------- | ----------- | -------------------------------------------------- |
| 1               | int/float/bool/char/string/array/hashmap | string      | Converts the provided argument to its string form. |

```kolon
fun: main() {
    println("hello!! " + toString(1)); // hello!! 1
    println("hello!! " + toString(1.1)); // hello!! 1.1
    println("hehe!! " + toString(true)); // hehe!! true
}
```

#### toFloat()

| **Num of Args** | **Type of Args** | **Returns** | **Description**                                   |
| --------------- | ---------------- | ----------- | ------------------------------------------------- |
| 1               | int/float/string | float       | Converts the provided argument to its float form. |

```kolon
fun: main() {
    var a: int = 10;
    println(toFloat(a)); // 10.0
    println(typeOf(toFloat(a))); // float
}
```

#### toInt()

| **Num of Args** | **Type of Args**      | **Returns** | **Description**                                 |
| --------------- | --------------------- | ----------- | ----------------------------------------------- |
| 1               | int/float/string/char | int         | Converts the provided argument to its int form. |

```kolon
fun: main() {
    println(toInt(11.9)); // 11 [always returns floor value]
    println(toInt("10")); // 10
    println(toInt('A')); // 65
}
```

#### typeOf()

| **Num of Args** | **Type of Args**                         | **Returns** | **Description**                        |
| --------------- | ---------------------------------------- | ----------- | -------------------------------------- |
| 1               | int/float/bool/char/string/array/hashmap | string      | Returns the type of the given argument |

```kolon
fun: main() {
    var a: int = 10;
    var b: float = 10.1;
    var c: string = "someString";
    var d: char = 'a';
    var e: bool = true;
    var f: int[] = [1, 2, 3, 4];
    var g: string[int] = {"Khush": 1};
    println(typeOf(a)); // int
    println(typeOf(b)); // float
    println(typeOf(c)); // string
    println(typeOf(d)); // char
    println(typeOf(e)); // bool
    println(typeOf(f)); // int[]
    println(typeOf(g)); // string[int]
}
```

#### scan()

| **Num of Args** | **Type of Args** | **Returns** | **Description**                                                                                                                                                          |
| --------------- | ---------------- | ----------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| 0               | -                | string      | Waits for user to input 1 or more than 1 lines.                                                                                                                          |
| 1               | string           | string      | Waits for user to input 1 or more than 1 lines. First arg is prompt, will be printed when the function is called. By default WON'T take input from next line             |
| 2               | string, bool     | string      | Waits for user to input 1 or more than 1 lines. First arg is prompt, will be printed when the function is called. Second arg if true will take user input from next line |

#### scanln()

| **Num of Args** | **Type of Args** | **Returns** | **Description**                                                                                                                                          |
| --------------- | ---------------- | ----------- | -------------------------------------------------------------------------------------------------------------------------------------------------------- |
| 0               | -                | string      | Waits for user to input 1 line.                                                                                                                          |
| 1               | string           | string      | Waits for user to input 1 line. First arg is prompt, will be printed when the function is called. By default WON'T take input from next line             |
| 2               | string, bool     | string      | Waits for user to input 1 line. First arg is prompt, will be printed when the function is called. Second arg if true will take user input from next line |

#### push()

| **Data Structure** | **Num of Args** | **Type of Args**                                | **Returns** | **Format**             | **Description**                          |
| ------------------ | --------------- | ----------------------------------------------- | ----------- | ---------------------- | ---------------------------------------- |
| Array              | 2               | array, int/float/bool/string/char               | array       | push(array, element);  | Adds an element to the end of the array. |
| HashMap            | 3               | hashmap, int/float/bool/string/char(key, value) | hashmap     | push(map, key, value); | Adds a key-value pair to the hashmap.    |

#### pop()

| **Data Structure** | **Num of Args** | **Type of Args**   | **Returns**                | **Format**           | **Description**                                         |
| ------------------ | --------------- | ------------------ | -------------------------- | -------------------- | ------------------------------------------------------- |
| Array              | 1               | array              | int/float/bool/string/char | pop(array);          | Removes the last element from the array and returns it. |
| Array              | 2               | array, int (index) | int/float/bool/string/char | pop(array, element); | Removes the element at the given index and returns it.  |

#### insert()

| **Data Structure** | **Num of Args** | **Type of Args**                               | **Returns** | **Format**                     | **Description**                                                               |
| ------------------ | --------------- | ---------------------------------------------- | ----------- | ------------------------------ | ----------------------------------------------------------------------------- |
| Array              | 3               | array, int (index), int/float/string/bool/char | array       | insert(array, index, element); | Inserts the given element at the specified index in the array and returns it. |

#### remove()

| **Data Structure** | **Num of Args** | **Type of Args**                    | **Returns**                | **Format**              | **Description**                                                                                     |
| ------------------ | --------------- | ----------------------------------- | -------------------------- | ----------------------- | --------------------------------------------------------------------------------------------------- |
| Array              | 2               | array, int/float/string/bool/char   | array                      | remove(array, element); | Removes the first occurrence of the specified element from the array and returns the updated array. |
| HashMap            | 2               | hashmap, int/float/string/bool/char | int/float/string/bool/char | remove(map, key);       | Removes the key-value pair for the specified key from the hashmap and returns the removed value.    |

#### getIndex()

| **Data Structure** | **Num of Args** | **Type of Args**                  | **Returns** | **Format**                | **Description**                                                                                                          |
| ------------------ | --------------- | --------------------------------- | ----------- | ------------------------- | ------------------------------------------------------------------------------------------------------------------------ |
| Array              | 2               | array, int/float/string/bool/char | int         | getIndex(array, element); | Returns the index of the first occurrence of the specified element in the array. Returns -1 if the element is not found. |

#### keys()

| **Data Structure** | **Num of Args** | **Type of Args** | **Returns**     | **Format** | **Description**                                                 |
| ------------------ | --------------- | ---------------- | --------------- | ---------- | --------------------------------------------------------------- |
| Hash               | 1               | hash             | array (of keys) | keys(map); | Returns an array containing all the keys of the given hash map. |

#### values()

| **Data Structure** | **Num of Args** | **Type of Args** | **Returns**       | **Format**   | **Description**                                                   |
| ------------------ | --------------- | ---------------- | ----------------- | ------------ | ----------------------------------------------------------------- |
| Hash               | 1               | hash             | array (of values) | values(map); | Returns an array containing all the values of the given hash map. |

#### containsKey()

| **Data Structure** | **Num of Args** | **Type of Args**                       | **Returns** | **Format**             | **Description**                                                            |
| ------------------ | --------------- | -------------------------------------- | ----------- | ---------------------- | -------------------------------------------------------------------------- |
| Hash               | 2               | hash, key (int/float/string/bool/char) | bool        | containsKey(map, key); | Returns `true` if the given key exists in the hash map, otherwise `false`. |

#### slice()

| **Data Structure** | **Num of Args** | **Type of Args**      | **Returns** | **Format**                       | **Description**                                                                                   |
| ------------------ | --------------- | --------------------- | ----------- | -------------------------------- | ------------------------------------------------------------------------------------------------- |
| Array              | 3               | array, int, int       | array       | slice(array, start, end);        | Returns an array with element from start index (inclusive) to end index (exclusive)               |
| Array              | 4               | array, int, int, int  | array       | slice(array, start, end, step);  | Returns an array with element from start index (inclusive) to end index (exclusive) with steps    |
| string             | 3               | string, int, int      | string      | slice(string, start, end);       | Returns a string with characters from start index (inclusive) to end index (exclusive)            |
| string             | 4               | string, int, int, int | string      | slice(string, start, end, step); | Returns a string with characters from start index (inclusive) to end index (exclusive) with steps |

### Overriding Built-in functions

You can override built-in functions by simply defining them in your file. If a function with the same name exists both in the file and as a built-in, Kolon will give preference to the version defined in the file.

```kolon
fun: main() {
    var a: int = len("hi");
    println(a); // 100
}
fun: len(a: string): (int) {
    return: 100;
}
```

## Return Statements

The return statement is used to exit from a function and return to where it was called. In the case of the main function, the return statement must empty.

### Return Statement in Functions with No Return Type

For functions that do not have a return type, the return statement simply stops the function’s execution and exits:

```kolon
fun: main() {
    var a: int = 1;
    if: (a == 1): {
        println(true); // Will be printed
        return; // Exits the fun :(
    }
    println(false); // Won't be printed
}
```

### Return Statement in Functions with a Return Type

For functions with a specified return type, the type of the value(s) returned must match the defined return type in the function signature.

#### Returning Multiple Values

If the function returns multiple values, the values should be wrapped inside `(...)`

```kolon
fun: main() {
    var a: int, var b: bool = callMe();
}
fun: callMe(): (int, bool) {
    return: (100, true);
}
```

#### Returning a Single Value

If the function returns only a single value, `()` should not be used:

```kolon
fun: main() {
    var a: int = callMe();
}
fun: callMe(): (int) {
    return: 100;
}
```

## Prefix Operation

Kolon supports two prefix symbols: `-` (Minus) and `!` (Not). These symbols can be used with specific data types:

- `-` is used for `int` and `float` types to negate the value.
- `!` is used for `bool` types to negate the boolean value.

```kolon
fun: main() {
    var a: int = -10; // -10
    var c: int = -a; // 10
    var b: bool = !true; // false
}
```

## Postfix Operation

Kolon supports two prefix symbols: `++` and `--`.

- Both of these symbols can only be used with `int` and `float`

```kolon
fun: main() {
    var a: int = 10; // 10
    var b: int = a++; // 11
    b--; // 10
    var c: float = 2.14;
    c++; // 3.14
}
```

## Infix Operation

An infix operation has three inputs:

- Left operand
- Operator
- right operand

NOTE: If the left and right operands are call expressions that return a value, their length must be equal to 1.

### Left: `int`, Right: `int`

| Operator | Description      | Example | Output Type |
| -------- | ---------------- | ------- | ----------- |
| +        | Addition         | 5 + 3   | Integer     |
| -        | Subtraction      | 5 - 3   | Integer     |
| /        | Division         | 6 / 2   | Integer     |
| \*       | Multiplication   | 5 \* 3  | Integer     |
| %        | Modulus          | 5 % 3   | Integer     |
| &        | Bitwise AND      | 5 & 3   | Integer     |
| \|       | Bitwise OR       | 5 \| 3  | Integer     |
| >        | Greater than     | 5 > 3   | Boolean     |
| <        | Less than        | 3 < 5   | Boolean     |
| >=       | Greater or equal | 5 >= 5  | Boolean     |
| <=       | Less or equal    | 3 <= 5  | Boolean     |
| ==       | Equal            | 5 == 5  | Boolean     |
| !=       | Not equal        | 5 != 3  | Boolean     |

### Left: `float`, Right: `float`

| Operator | Description      | Example    | Output Type |
| -------- | ---------------- | ---------- | ----------- |
| +        | Addition         | 5.2 + 3.1  | Float       |
| -        | Subtraction      | 5.2 - 3.1  | Float       |
| /        | Division         | 6.4 / 2.0  | Float       |
| \*       | Multiplication   | 5.5 \* 3.2 | Float       |
| >        | Greater than     | 5.5 > 3.2  | Boolean     |
| <        | Less than        | 3.2 < 5.5  | Boolean     |
| >=       | Greater or equal | 5.5 >= 5.5 | Boolean     |
| <=       | Less or equal    | 3.2 <= 5.5 | Boolean     |
| ==       | Equal            | 5.5 == 5.5 | Boolean     |
| !=       | Not equal        | 5.5 != 3.2 | Boolean     |

### Left: `bool`, Right: `bool`

| Operator | Description | Example         | Output Type |
| -------- | ----------- | --------------- | ----------- |
| ==       | Equal       | true == true    | Boolean     |
| !=       | Not equal   | true != false   | Boolean     |
| &&       | Logical AND | true && false   | Boolean     |
| \|\|     | Logical OR  | true \|\| false | Boolean     |

### Left: `string`, Right: `string`

| Operator | Description   | Example        | Output Type |
| -------- | ------------- | -------------- | ----------- |
| +        | Concatenation | "foo" + "bar"  | String      |
| ==       | Equal         | "foo" == "foo" | Boolean     |
| !=       | Not equal     | "foo" != "bar" | Boolean     |

### Left: `char`, Right: `char`

| Operator | Description   | Example    | Output Type |
| -------- | ------------- | ---------- | ----------- |
| +        | Concatenation | 'a' + 'a'  | String      |
| ==       | Equal         | 'a' == 'a' | Boolean     |
| !=       | Not equal     | 'a' != 'a' | Boolean     |

### Left: array, Right: array

| Operator | Description   | Example                | Output Type |
| -------- | ------------- | ---------------------- | ----------- |
| +        | Concatenation | [1, 2, 3] + [4, 5, 6]  | Array       |
| ==       | Equal         | [1, 2, 3] == [4, 5, 6] | Boolean     |
| !=       | Not equal     | [1, 2, 3] != [4, 5, 6] | Boolean     |

### Left: `float`, Right: `int` || Left: `int`, Right: `float`

Variable with `int` type will be converted to `float` type. Than the operation will be performed.

## Assignment Operation

An assignment operation has three inputs:

- Left operand
- Operator
- right operand

Operators supported: `=`, `+=`, `-=`, `*=`, `/=`, `%=`

Example:

- `a = 10` - Stores 10 in variable `a`
- `a += 10` - Adds 10 to the value of `a` and stores it in `a`
- `a -= 10` - Subtracts 10 from the value of `a` and stores it in `a`
- `a *= 10` - Multiplies the value of `a` by 10 and stores it in `a`
- `a /= 10` - Divides the value of `a` by 10 and stores it in `a`
- `a %= 10` - Takes the modulus of the value of `a` by 10 and stores it in `a`

Assignment operation can be interpreted as `a = a + 10` or `a = a - 10` or `a = a * 10` or `a = a / 10` or `a = a % 10`. Hence all the operations in infix operation can be performed using assignment operation. Only EXCEPTION is left = `int` and right = `float`
