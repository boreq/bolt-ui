import { Location } from 'vue-router';
import { Key as KeyDTO } from '@/dto/Entry';

export class NavigationService {

    getBrowse(path: KeyDTO[], value: KeyDTO): Location {
        const query = this.getQuery(value);

        if (path.length === 0) {
            return {
                name: 'browse',
                query: query,
            };
        }

        const stringPath = path.map(key => key.hex).join('/');
        return {
            name: 'browse-children',
            params: {
                pathMatch: stringPath,
            },
            query: query,
        };
    }

    private getQuery(value: KeyDTO): { value: string } {
        if (!value) {
            return null;
        }

        return {
            value: value.hex,
        };
    }

}
