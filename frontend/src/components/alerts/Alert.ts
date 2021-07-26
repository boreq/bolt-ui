import { Component, Vue, Prop } from 'vue-property-decorator';

@Component
export default class Alert extends Vue {

    @Prop()
    kind: Kind;

}

enum Kind {
    Warning = 'warning',
}
