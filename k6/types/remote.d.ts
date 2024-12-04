declare module 'https://*';

declare module 'https://jslib.k6.io/k6-utils/1.6.0/index.js' {
    export function randomString(length: number, charset?: string): string;
}

declare module 'https://jslib.k6.io/url/1.0.0/index.js' {
    export { URL, URLSearchParams } from 'url';
}
