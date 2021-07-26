import { Vue } from 'vue-property-decorator';
import axios, { AxiosResponse } from 'axios'; // do not add { }, some webshit bs?
import { Mutation } from '@/store';
import { AuthService } from '@/services/AuthService';
import { Tree } from '@/dto/Tree';

const authTokenHeaderName = 'Access-Token';

export class ApiService {

    private readonly axios = axios.create();
    private readonly authService = new AuthService();

    constructor(private vue: Vue) {
        this.axios.interceptors.request.use(
            config => {
                const token = this.authService.getToken();
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
                    this.authService.clearToken();
                    this.vue.$store.commit(Mutation.SetToken, null);
                }
                return Promise.reject(error);
            });
    }

    browse(path: string, before: string, after: string): Promise<AxiosResponse<Tree>> {
        const url = path ? `browse/${path}` : `browse/`;
        return this.axios.get<Tree>(
            process.env.VUE_APP_API_PREFIX + url,
            {
                params: this.browseParams(before, after),
            },
        );
    }

    private browseParams(before: string, after: string): { before: string } | { after: string } | null {
        if (before && after) {
            throw new Error('defined both before and after');
        }

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

        return null;
    }

}
