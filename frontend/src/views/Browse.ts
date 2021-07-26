import { Component, Vue } from 'vue-property-decorator';
import { ApiService } from '@/services/ApiService';
import { Tree } from '@/dto/Tree';
import { Entry as EntryDTO, Key as KeyDTO } from '@/dto/Entry';

import Notifications from '@/components/Notifications.vue';
import Entries from '@/components/Entries.vue';
import Value from '@/components/Value.vue';
import Key from '@/components/Key.vue';


@Component({
    components: {
        Entries,
        Value,
        Key,
    },
})
export default class Browse extends Vue {

    trees: Tree[] = [];
    selected: EntryDTO = null;

    private readonly apiService = new ApiService(this);
    private readonly numVisibleTrees = 3;

    get selectedPath(): KeyDTO[] {
        if (this.trees.length === 0) {
            return [];
        }

        const path =  [
            ...this.trees[this.trees.length - 1].path,
        ];

        if (this.selected) {
            path.push(this.selected.key);
        }

        return path;
    }

    get visibleTrees(): Tree[] {
        const trees = [];

        let minIndex = this.trees.length - this.numVisibleTrees;
        if (this.selected) {
            minIndex++;
        }

        this.trees.forEach((tree, index) => {
            if (index >= minIndex) {
                trees.push(tree);
            }
        });

        return trees;
    }

    created(): void {
        this.load();
    }

    treeKey(tree: Tree): string {
        return tree.path.map(v => v.hex).join('-');
    }

    onHeaderClick(): void {
        this.load();
    }

    onEntry(tree: Tree, entry: EntryDTO): void {
        const index = this.trees.indexOf(tree);
        if (index >= 0) {
            this.trees.length = index + 1;
        }

        this.selected = null;

        if (entry.bucket) {
            const path = [...tree.path, entry.key];
            const stringPath  = path.map(key => key.hex).join('/');
            this.apiService.browse(stringPath, null, null)
                .then(
                    result => {
                        this.trees.push(result.data);
                    },
                    error => {
                        Notifications.pushError(this, 'Could not query the backend.', error);
                    },

                );
        } else {
            this.selected = entry;
        }
    }

    selectedInTree(tree: Tree): EntryDTO {
        const index = this.trees.indexOf(tree);
        if (index < 0) {
            return null;
        }

        const selectedPath =  this.selectedPath;
        if (index > selectedPath.length - 1) {
            return null;
        }

        return tree.entries.find(entry => entry.key.hex === selectedPath[index].hex);
    }

    private load(): void {
        this.trees = [];
        this.selected = null;

        this.apiService.browse(null, null, null)
            .then(
                result => {
                    this.trees = [
                        result.data,
                    ];
                },
                error => {
                    Notifications.pushError(this, 'Could not query the backend.', error);
                },

            );
    }
}
