import { Component, Vue } from 'vue-property-decorator';

import Alert from '@/components/alerts/Alert.vue';

@Component({
    components: {
        Alert,
    },
})
export default class WarningAlert extends Vue {
}
