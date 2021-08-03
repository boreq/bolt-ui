import { Key } from '@/dto/Entry';

export class PathService {

    marshal(path: Key[]): string {
        const elements: string[] = [];

        for (const key of path) {
            if (key.str) {
                elements.push(`"${key.str}"`);
            } else {
                elements.push(`0x${key.hex}`);
            }
        }

        return elements.join(' / ');
    }

    //unmarshal(path: string): Key[] {
    //}

}
