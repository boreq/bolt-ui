import Vue from 'vue';
import Vuex from 'vuex';

Vue.use(Vuex);

export enum Mutation {
    SetToken = 'setToken',
}

export class State {
    token: string;
}

export default new Vuex.Store<State>({
    state: {
        token: undefined,
    },
    mutations: {
        [Mutation.SetToken](state: State, token: string): void {
            state.token = token;
        },
    },
});

