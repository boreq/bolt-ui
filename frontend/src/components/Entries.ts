import { Component, Vue, Prop } from 'vue-property-decorator';
import { Entry } from '@/dto/Entry';

import Key from '@/components/Key.vue'

@Component({
    components: {
        Key,
    },
})
export default class Entries extends Vue {

    @Prop()
    entries: Entry[];

    @Prop()
    selected: Entry;

    get isEmpty(): boolean {
        return this.entries && this.entries.length === 0;
    }

    onClick(entry: Entry): void {
        this.$emit('entry', entry);
    }

}
