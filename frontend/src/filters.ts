import Vue from 'vue';

export const dateFilter = Vue.filter('date', (value: string): string => {
    value = value.replace('T', ' ');
    value = value.replace('Z', ' UTC');
    return value;
});

export const distanceFilter = Vue.filter('distance', (value: number): string => {
    return round(value / 1000, 1) + 'km';
});

export const durationFilter = Vue.filter('duration', (value: number): string => {
    const units = ['s', 'm', 'h', 'd'];
    const multipliers = [60, 60, 24];

    const results = [];
    for (let i = 0; i < multipliers.length; i++) {
        const nextValue = value / multipliers[i];
        results.splice(0, 0, [value % multipliers[i], units[i]]);
        if (i === multipliers.length - 1) {
            results.splice(0, 0, [nextValue, units[i + 1]]);
        }
        value = nextValue;
    }

    return results
        .filter(v => v[0] >= 1)
        .filter((_, index) => index < 2)
        .map(v => Math.floor(v[0]) + v[1])
        .join(' ');
});

function round(value, precision) {
    const multiplier = Math.pow(10, precision || 0);
    return Math.round(value * multiplier) / multiplier;
}
