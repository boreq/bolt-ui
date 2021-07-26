
import { Component, Vue, Prop } from 'vue-property-decorator';
import { Entry as Entry } from '@/dto/Entry';

import Key from '@/components/Key.vue';

@Component({
    components: {
        Key,
    },
})
export default class Value extends Vue {

    @Prop()
    entry: Entry;

    get format(): string {
        if (!this.entry.value) {
            return 'nil';
        }

        if (this.entry.value.str) {
            return 'string';
        }

        return 'hex';
    }

    get formatTooltip(): string {
        if (!this.entry.value) {
            return 'The value is empty.';
        }

        if (this.entry.value.str) {
            return 'Displaying the value as string.';
        }

        return 'Display the bytes using hexadecimal encoding.';
    }

    get valueString(): string {
        if (!this.entry.value) {
            return null;
        }

        if (this.entry.value.str) {
            return this.entry.value.str;
        }

        return this.entry.value.hex;
    }

}
