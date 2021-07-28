import { Component, Vue, Prop } from 'vue-property-decorator';
import { Tree as TreeDTO } from '@/dto/Tree';
import { Entry as EntryDTO, Key as KeyDTO } from '@/dto/Entry';
import { ApiService } from '@/services/ApiService';

import Notifications from '@/components/Notifications.vue';
import Entries from '@/components/Entries.vue'
import Spinner from '@/components/Spinner.vue'

@Component({
    components: {
        Entries,
        Spinner,
    },
})
export default class Tree extends Vue {

    @Prop()
    path: KeyDTO[];

    @Prop()
    selected: KeyDTO[];

    tree: TreeDTO = null;
 
    private readonly apiService = new ApiService(this);

    get selectedInTree(): EntryDTO {
        if (!this.selected || !this.tree) {
            return null;
        }

        for (const entry of this.tree.entries) {
            const entryPath = [
                ...this.path,
                entry.key,
            ]

            if (this.pathHasPrefix(this.selected, entryPath)) {
                return entry;
            }
        }

        return null;
    }

    created(): void {
        this.load();
    }


    onEntry(entry: EntryDTO): void {
        this.$emit('entry', entry);
    }

    private load(): void {
        this.tree = null;

        const stringPath  = this.path.map(key => key.hex).join('/');

        this.apiService.browse(stringPath, null, null)
            .then(
                result => {
                    this.tree = result.data;
                },
                error => {
                    Notifications.pushError(this, 'Could not query the backend.', error);
                },
            );
    }

    private pathHasPrefix(path: KeyDTO[], prefix: KeyDTO[]): boolean { 
        if (prefix.length > path.length) {
            return false;
        }

        for (let i = 0; i < prefix.length; i++) {
            if (prefix[i].hex !== path[i].hex) {
                return false;
            }
        }
        return true;
    }

}
