
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

        if (this.entry.value.pretty) {
            return this.entry.value.pretty.content_type;
        }

        return 'unknown';
    }

    get formatTooltip(): string {
        if (!this.entry.value) {
            return 'The value is empty.';
        }

        if (this.entry.value.pretty) {
            return `Recognized content type ${this.entry.value.pretty.content_type} for pretty printing.`;
        }

        return 'Pretty printing is unavailable due to unrecognized content type of this value.';
    }

    get valuePretty(): string {
        if (!this.entry.value) {
            return null;
        }

        if (this.entry.value.pretty) {
            return this.entry.value.pretty.value;
        }

        return null;
    }

    get valueHex(): string {
        if (!this.entry.value) {
            return null;
        }

        return this.entry.value.hex;
    }
}
