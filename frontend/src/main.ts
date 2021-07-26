import Vue from 'vue';
import App from './App.vue';
import router from './router';
import store from './store';
import VTooltip from 'v-tooltip';
import VueTimeago from 'vue-timeago';
import { distanceFilter, durationFilter, dateFilter } from './filters';

Vue.use(VTooltip);

Vue.use(VueTimeago, {});

Vue.config.productionTip = false;

new Vue({
    router,
    store,
    render: (h) => h(App),
    filters: {
        distanceFilter,
        durationFilter,
        dateFilter,
    },
}).$mount('#app');
