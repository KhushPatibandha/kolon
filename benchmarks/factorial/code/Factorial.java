public class Factorial {
    public static long factorial(int n) {
        if (n == 0 || n == 1)
            return 1;
        return n * factorial(n - 1);
    }

    public static void main(String[] args) {
        long result = 0;
        int p = 0;
        for (int i = 0; i < 1_000_000; i++) {
            result = factorial(20);
            p = i;
        }
        System.out.println(result);
        System.out.println(p);
    }
}
