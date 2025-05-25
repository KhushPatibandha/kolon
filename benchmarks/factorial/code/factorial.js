function factorial(n) {
    if (n === 0 || n === 1) return 1;
    return n * factorial(n - 1);
}

let result;
let p;
for (let i = 0; i < 1_000_000; i++) {
    result = factorial(20);
    p = i;
}

console.log(result);
console.log(p);
