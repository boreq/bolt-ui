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

    unmarshal(path: string): Key[] {
        const buf = new StringBuffer(path);
        const result: Key[] = [];

        for (let state: stateFunc = stateStart; state;) {
            state = state(buf, result);
        }

        for (const key of result) {
            if (!key.hex) {
                key.hex = this.hexEncode(key.str);
            }
        }

        return result;
    }

    private hexEncode(s: string): string {
        let result = '';
        for (let i = 0; i < s.length; i++) {
            const hex = s.charCodeAt(i).toString(16);
            result += hex;
        }
        return result;
    }

}

class StringBuffer {

    constructor(private s: string) {
    }

    next(): string {
        if (this.s.length === 0) {
            return tokenNull;
        }

        const v = this.s[0];
        this.s = this.s.slice(1);
        return v;
    }

}

type stateFunc = (buf: StringBuffer, result: Key[]) => stateFunc;

const tokenNull = null;
const tokenSpace = ' ';
const tokenSeparator = '/';

const tokenStrStart = '"';
const tokenStrEnd = '"';

const tokenHexStartFirst = '0';
const tokenHexStartSecondSmall = 'x';
const tokenHexStartSecondCapital = 'X';

const errInvalid = 'invalid path';

function stateStart(buf: StringBuffer): stateFunc {
    const next = buf.next();
    switch (next) {
        case tokenSpace:
            return stateStart; // keep consuming spaces
        case tokenStrStart:
            return stateBeforeStr;
        case tokenHexStartFirst:
            return stateHexStartSecond;
        case tokenNull:
            return null; // this is an empty path
        default:
            throw errInvalid;
    }
}

function stateBeforeSeparator(buf: StringBuffer): stateFunc {
    const next = buf.next();
    switch (next) {
        case tokenSpace:
            return stateBeforeSeparator; // keep consuming spaces
        case tokenSeparator:
            return stateAfterSeparator;
        case tokenNull:
            return null; // we read the full path
        default:
            throw errInvalid;
    }
}

function stateAfterSeparator(buf: StringBuffer): stateFunc {
    const next = buf.next();
    switch (next) {
        case tokenSpace:
            return stateAfterSeparator; // keep consuming spaces
        case tokenStrStart:
            return stateBeforeStr;
        case tokenHexStartFirst:
            return stateHexStartSecond;
        default:
            throw errInvalid;
    }
}

function stateBeforeStr(_: StringBuffer, result: Key[]): stateFunc {
    result.push(
        {
            hex: null,
            str: '',
        },
    );
    return stateStr;
}

function stateStr(buf: StringBuffer, result: Key[]): stateFunc {
    const next = buf.next();
    switch (next) {
        case tokenStrEnd:
            return stateBeforeSeparator;
        case tokenNull:
            throw errInvalid;
        default:
            result[result.length - 1].str += next;
            return stateStr;
    }
}

function stateHexStartSecond(buf: StringBuffer): stateFunc {
    const next = buf.next();
    switch (next) {
        case tokenHexStartSecondSmall:
        case tokenHexStartSecondCapital:
            return stateBeforeHex;
        default:
            throw errInvalid;
    }
}

function stateBeforeHex(_: StringBuffer, result: Key[]): stateFunc {
    result.push(
        {
            hex: '',
            str: null,
        },
    );
    return stateHex;
}


function stateHex(buf: StringBuffer, result: Key[]): stateFunc {
    const hex = ['0', '1', '2', '3', '4', '5', '6', '7', '8', '9', 'a', 'b', 'c', 'd', 'e', 'f'];
    const next = buf.next();

    switch (next) {
        case tokenSpace:
            return stateBeforeSeparator;
        case tokenSeparator:
            return stateAfterSeparator;
        case tokenNull:
            return null; // we read the whole hex
    }

    if (hex.includes(next.toLowerCase())) {
        result[result.length - 1].hex += next;
        return stateHex;
    }

    throw errInvalid; 
}
