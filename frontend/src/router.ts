import Vue from 'vue';
import Router from 'vue-router';

import Browse from '@/views/Browse.vue';

Vue.use(Router);

export default new Router({
    mode: 'history',
    base: process.env.BASE_URL,
    routes: [
        {
            path: '/*',
            name: 'browse-children',
            component: Browse,
        },
        {
            path: '/',
            name: 'browse',
            component: Browse,
        },
        {
            path: '*',
            redirect: {
                name: 'browse',
            }
        },
    ],
});
