public class Count {
    public static void main(String[] args) {
        long res = 0;
        for (int i = 1; i <= 1_000_000; i++) {
            res += i;
        }
        System.out.println(res);
    }
}
