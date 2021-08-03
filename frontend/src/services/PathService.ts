import { Key } from '@/dto/Entry';

export class UnmarshalResult {
    path: Key[]
    value: Key;
}

export class PathService {

    marshal(path: Key[], value: Key): string {
        const elements: string[] = [];

        for (const key of path) {
            if (key.str) {
                elements.push(
                    tokenStrStart + key.str + tokenStrEnd
                );
            } else {
                elements.push(
                    tokenHexStartFirst + tokenHexStartSecondSmall + key.hex
                );
            }
        }

        let result = elements.join(tokenSpace + tokenSeparator + tokenSpace);
        if (value) {
            const prefix = tokenSpace + tokenValueSeparator + tokenSpace;
            if (value.str) {
                result +=  prefix + value.str;
            } else {
                result +=  prefix + tokenHexStartFirst + tokenHexStartSecondSmall + value.hex;
            }
        }

        return result;
    }

    unmarshal(path: string): UnmarshalResult {
        const buf = new StringBuffer(path);
        const result: ResultElement[] = [];

        for (let state: stateFunc = stateStart; state;) {
            state = state(buf, result);
        }

        return this.convert(result);
    }

    private convert(parsed: ResultElement[]): UnmarshalResult {
        const result: UnmarshalResult = {
            path: [],
            value: null,
        }

        let foundValueSeparator = false;

        for (const element of parsed){
            if (!isSeparator(element)) {
                if (result.value) {
                    throw 'Encountered bucket after value.'; 
                }

                if (!element.hex) {
                    element.hex = this.hexEncode(element.str);
                }

                if (foundValueSeparator) {
                    result.value = element;
                } else {
                    result.path.push(element);
                }
            } else {
                if (!element.bucket) {
                    foundValueSeparator = true;
                }
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

    private last: string = null;

    constructor(private s: string) {
    }

    next(): string {
        if (this.s.length === 0) {
            return tokenNull;
        }

        this.last = this.s[0];
        this.s = this.s.slice(1);
        return this.last;
    }

    unread(): void {
        if (this.s) {
            this.s = this.last + this.s;
        } else {
            this.s = this.last;
        }
    }

}

class Separator {
    bucket: boolean;
}

type ResultElement = Key | Separator;

function isSeparator(v: Key | Separator): v is Separator {
   return (<Separator>v).bucket !== undefined;
}

type stateFunc = (buf: StringBuffer, result: ResultElement[]) => stateFunc;

const tokenNull = null;
const tokenSpace = ' ';
const tokenSeparator = '/';

const tokenStrStart = '"';
const tokenStrEnd = '"';

const tokenHexStartFirst = '0';
const tokenHexStartSecondSmall = 'x';
const tokenHexStartSecondCapital = 'X';

const tokenValueSeparator = "-"

const errInvalid = 'invalid path';

function stateStart(buf: StringBuffer): stateFunc {
    const next = buf.next();
    switch (next) {
        case tokenSpace:
            return stateStart; // keep consuming spaces
        case tokenStrStart:
            return stateBeforeBucketStr;
        case tokenHexStartFirst:
            return stateHexStartSecond;
        case tokenNull:
            return null; // this is an empty path
        default:
            throw errInvalid;
    }
}

function stateAfterElement(buf: StringBuffer): stateFunc {
    const next = buf.next();
    switch (next) {
        case tokenSpace:
            return stateAfterElement; // keep consuming spaces
        case tokenSeparator:
            return stateSeparator;
        case tokenValueSeparator:
            return stateValueSeparator;
        case tokenNull:
            return null; // we read the full path
        default:
            throw errInvalid;
    }
}

function stateSeparator(_: StringBuffer, result: ResultElement[]): stateFunc {
    result.push(
        {
            bucket: true,
        }
    );

    return stateBeforeElement;
}

function stateBeforeElement(buf: StringBuffer): stateFunc {
    const next = buf.next();
    switch (next) {
        case tokenSpace:
            return stateBeforeElement; // keep consuming spaces
        case tokenStrStart:
            return stateBeforeBucketStr;
        case tokenHexStartFirst:
            return stateHexStartSecond;
        default:
            throw errInvalid;
    }
}

function stateBeforeBucketStr(_: StringBuffer, result: ResultElement[]): stateFunc {
    result.push(
        {
            hex: null,
            str: '',
        },
    );
    return stateStr;
}

function stateStr(buf: StringBuffer, result: ResultElement[]): stateFunc {
    const next = buf.next();
    switch (next) {
        case tokenStrEnd:
            return stateAfterElement;
        case tokenNull:
            throw errInvalid;
        default:
            (result[result.length - 1] as Key).str += next;
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

function stateBeforeHex(_: StringBuffer, result: ResultElement[]): stateFunc {
    result.push(
        {
            hex: '',
            str: null,
        },
    );
    return stateHex;
}

const hex = ['0', '1', '2', '3', '4', '5', '6', '7', '8', '9', 'a', 'b', 'c', 'd', 'e', 'f'];

function stateHex(buf: StringBuffer, result: ResultElement[]): stateFunc {
    const next = buf.next();

    switch (next) {
        case tokenNull:
            return null; // we read the whole hex
    }

    if (hex.includes(next.toLowerCase())) {
        (result[result.length - 1] as Key).hex += next;
        return stateHex;
    }

    buf.unread();
    return stateAfterElement;
}

function stateValueSeparator(_: StringBuffer, result: ResultElement[]): stateFunc {
    result.push(
        {
            bucket: false,
        }
    );

    return stateBeforeElement;
}

