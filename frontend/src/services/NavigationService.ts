import { Location } from 'vue-router';

export class NavigationService {

    getBrowse(): Location {
        return {
            name: 'browse',
        };
    }

}
