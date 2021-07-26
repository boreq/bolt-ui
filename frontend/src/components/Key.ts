import { Component, Vue, Prop } from 'vue-property-decorator';
import { Key as KeyDTO } from '@/dto/Entry';

@Component
export default class Key extends Vue {

    @Prop()
    k: KeyDTO;

}
