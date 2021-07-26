import { Component, Vue } from 'vue-property-decorator';

import Notifications from '@/components/Notifications.vue';


@Component({
    components: {
        Notifications,
    },
})
export default class App extends Vue {
}
