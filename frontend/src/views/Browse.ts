import { Component, Vue } from 'vue-property-decorator';
import { Mutation } from '@/store';
import { Entry as EntryDTO, Key as KeyDTO } from '@/dto/Entry';

import Tree from '@/components/Tree.vue';
import Value from '@/components/Value.vue';
import Key from '@/components/Key.vue';


@Component({
    components: {
        Tree,
        Value,
        Key,
    },
})
export default class Browse extends Vue {

    paths: KeyDTO[][] = [];
    selected: KeyDTO[] = null;
    selectedEntry: EntryDTO = null;

    private readonly numVisibleTrees = 3;

    get visiblePaths(): Tree[] {
        const paths = [];

        let minIndex = this.paths.length - this.numVisibleTrees;
        if (this.selectedEntry) {
            minIndex++;
        }

        this.paths.forEach((path, index) => {
            if (index >= minIndex) {
                paths.push(path);
            }
        });

        return paths;
    }

    created(): void {
        this.setToken();
        this.load();
    }

    treeKey(path: KeyDTO[]): string {
        return path.map(v => v.hex).join('-');
    }

    onHeaderClick(): void {
        this.load();
    }

    onEntry(path: KeyDTO[], entry: EntryDTO): void {
        const index = this.paths.indexOf(path);
        if (index >= 0) {
            this.paths.length = index + 1;
        }

        const childPath = [...path, entry.key];
        this.selected = childPath;

        if (entry.bucket) {
            this.paths.push(childPath);
            this.selectedEntry = null;
        } else {
            this.selectedEntry = entry;
        }
    }

    private setToken(): void {
        const token = this.$route.query.token;
        this.$store.commit(Mutation.SetToken, token);
    }

    private load(): void {
        this.paths = [
            [],
        ];
        this.selected = null;
    }
}
