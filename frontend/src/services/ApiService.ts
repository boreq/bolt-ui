import { Vue } from 'vue-property-decorator';
import axios, { AxiosResponse } from 'axios'; // do not add { }, some webshit bs?
import { Mutation } from '@/store';
import { Tree } from '@/dto/Tree';

const authTokenHeaderName = 'Access-Token';

export class ApiService {

    private readonly axios = axios.create();

    constructor(private vue: Vue) {
        this.axios.interceptors.request.use(
            config => {
                const token = this.vue.$store.state.token;
                if (token) {
                    config.headers[authTokenHeaderName] = token;
                }
                return config;
            },
            error => {
                return Promise.reject(error);
            },
        );

        this.axios.interceptors.response.use(
            response => {
                return response;
            },
            error => {
                if (error.response && error.response.status === 401) {
                    this.vue.$store.commit(Mutation.SetToken, null);
                }
                return Promise.reject(error);
            });
    }

    browse(path: string, before: string, after: string, from: string): Promise<AxiosResponse<Tree>> {
        const url = path ? `browse/${path}` : `browse/`;
        return this.axios.get<Tree>(
            process.env.VUE_APP_API_PREFIX + url,
            {
                params: this.browseParams(before, after, from),
            },
        );
    }

    private browseParams(before: string, after: string, from: string): { before: string } | { after: string } | { from: string} | null {
        if (before) {
            return {
                before: before,
            };
        }

        if (after) {
            return {
                after: after,
            };
        }

        if (from) {
            return {
                from: from,
            };
        }

        return null;
    }

}
