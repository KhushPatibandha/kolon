def factorial(n):
    if n == 0 or n == 1:
        return 1
    return n * factorial(n - 1)

for i in range(1_000_000):
    result = factorial(20)
    p = i

print(result)
print(i)
