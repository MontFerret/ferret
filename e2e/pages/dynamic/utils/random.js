export default function random(min = 1000, max = 5000) {
    const val = Math.random() * 1000 * 10;

    if (val < min) {
        return min;
    }

    if (val > max) {
        return max;
    }

    return val;
}