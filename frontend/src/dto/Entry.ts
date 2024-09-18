export class Entry {
    bucket: boolean;
    key: Key;
    value: Value;
}

export class Key {
    hex: string;
    str: string;
}

export class Value {
    hex: string;
    pretty: Pretty;
}

export class Pretty {
    content_type: string;
    value: string;
}
