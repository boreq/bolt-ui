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
    selectedValue: EntryDTO = null;

    private readonly numVisibleTrees = 3;

    get visiblePaths(): Tree[] {
        const paths = [];

        let minIndex = this.paths.length - this.numVisibleTrees;
        if (this.selectedValue) {
            minIndex++;
        }

        this.paths.forEach((path, index) => {
            if (index >= minIndex) {
                paths.push(path);
            }
        });

        return paths;
    }

    get selectedPath(): KeyDTO[] {
        if (this.paths.length === 0) {
            return null;
        }
        return this.paths[this.paths.length - 1];
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

        if (entry.bucket) {
            const childPath = [...path, entry.key];
            this.paths.push(childPath);
            this.selectedValue = null;
        } else {
            this.selectedValue = entry;
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
        this.selectedValue = null;
    }
}
