import { Component, Vue, Watch } from 'vue-property-decorator';
import { Mutation } from '@/store';
import { Entry as EntryDTO, Key as KeyDTO } from '@/dto/Entry';
import { NavigationService } from '@/services/NavigationService';
import { PathService } from '@/services/PathService';

import Notifications from '@/components/Notifications.vue';
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
    selectedValueKey: KeyDTO = null;
    selectedValue: EntryDTO = null;

    editingSelectedPath = false;
    editedPath: string = null;

    private readonly navigationService = new NavigationService();
    private readonly pathService = new PathService();

    private readonly numVisibleTrees = 3;

    isTreeVisible(index: number): boolean {
        let minIndex = this.paths.length - this.numVisibleTrees;
        if (this.selectedValueKey) {
            minIndex++;
        }

        return index >= minIndex;
    }

    get selectedPath(): KeyDTO[] {
        if (this.paths.length === 0) {
            return null;
        }
        const path = [
            ...this.paths[this.paths.length - 1],
        ];
        if (this.selectedValueKey) {
            path.push(this.selectedValueKey);
        }
        return path;
    }

    @Watch('$route')
    onRouteChanged(): void {
        this.setToken();
        this.loadFromRoute();
    }

    created(): void {
        this.setToken();
        this.loadFromRoute();

        document.body.addEventListener('click', this.cancelEditing);
    }

    destroyed(): void {
        document.body.removeEventListener('click', this.cancelEditing);
    }

    treeKey(path: KeyDTO[]): string {
        return path.map(v => v.hex).join('-');
    }

    onHeaderClick(): void {
        this.loadBlank();
    }

    onEntry(path: KeyDTO[], entry: EntryDTO): void {
        const index = this.paths.indexOf(path);
        if (index >= 0) {
            this.paths.length = index + 1;
        }

        if (entry.bucket) {
            const childPath = [...path, entry.key];
            this.paths.push(childPath);
            this.selectedValueKey = null;

            const next = this.navigationService.getBrowse(childPath, null);
            this.$router.push(next);
        } else {
            const shouldNavigate = this.selectedValueKey?.hex !== entry.key?.hex;
            this.selectedValue = entry;
            this.selectedValueKey = entry.key;

            if (shouldNavigate) {
                const next = this.navigationService.getBrowse(path, entry.key);
                this.$router.push(next);
            }
        }
    }

    onPath(path: KeyDTO[]): void {
        const index = path.length;
        for (let i = 0; i < path.length; i++) {
            this.paths[index][i].str = path[i].str;
        }
    }

    startEditing(): void {
        if (this.paths.length > 0) {
            this.editedPath = this.pathService.marshal(
                this.paths[this.paths.length - 1],
                this.selectedValueKey,
            );
        }

        this.editingSelectedPath = true;
    }

    finishEditing(): void {
        try {
            const newPath = this.pathService.unmarshal(this.editedPath);

            this.loadBlank();

            for (let i = 1; i <= newPath.path.length; i++) {
                this.paths.push(
                    newPath.path.slice(0, i),
                );
            }

            this.selectedValueKey = newPath.value;
            this.editingSelectedPath = false;
        }
        catch (e) {
            Notifications.pushError(this, 'Invalid path.', e);
        }
    }

    cancelEditing(): void {
        this.editingSelectedPath = false;
    }

    private setToken(): void {
        const token = this.$route.query.token;
        this.$store.commit(Mutation.SetToken, token);
    }

    private loadBlank(): void {
        this.paths = [
            [],
        ];
        this.selectedValueKey = null;
        this.selectedValue = null;
    }

    private loadFromRoute(): void {
        this.loadBlank();

        const path: KeyDTO[] = this.$route.params.pathMatch
            .split('/')
            .filter(v => v !== "")
            .map(
                (v: string): KeyDTO => {
                    return {
                        hex: v,
                        str: null,
                    };
                }
            );

        for (let i = 1; i <= path.length; i++) {
            this.paths.push(
                path.slice(0, i),
            );
        }

        if (this.$route.query.value) {
            this.selectedValueKey = {
                hex: this.$route.query.value as string,
                str: null,
            };
        }
    }
}
