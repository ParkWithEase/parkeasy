import type {
    ErrorResponseJSON,
    HttpMethod,
    PathsWithMethod,
    RequestBodyJSON,
    RequiredKeysOf,
    SuccessResponseJSON,
} from 'openapi-typescript-helpers';

import { URL, URLSearchParams } from 'https://jslib.k6.io/url/1.0.0/index.js';
import type { Response } from 'k6/http';
import http from 'k6/http';
import { CONTENT_TYPE_HEADER } from '../constants';

export interface Client<Paths extends {}> {
    get: ClientMethod<Paths, 'get'>;
    post: ClientMethod<Paths, 'post'>;
}

export type ClientMethod<
    Paths extends Record<string, Record<HttpMethod, {}>>,
    Method extends HttpMethod,
> = <Path extends PathsWithMethod<Paths, Method>>(
    url: Path,
    opts: RequestOptions<Paths[Path][Method]>,
) => RequestResult<Paths[Path][Method]>;

export interface RequestResult<Operation extends Record<string | number, any>> {
    data?: SuccessResponseJSON<Operation>;
    error?: ErrorResponseJSON<Operation>;
    response: Response;
}

export type RequestOptions<Operation extends Record<string, any>> =
    OptParam<Operation> & OptBody<Operation> & {};

type OptParam<Operation> = Operation extends { parameters: any }
    ? RequiredKeysOf<Operation['parameters']> extends never
        ? { params?: Operation['parameters'] }
        : { params: Operation['parameters'] }
    : {
          params: {
              query?: Record<string, unknown>;
          };
      };

type OptBody<Operation> = RequestBodyJSON<Operation> extends never
    ? { body?: never }
    : { body: RequestBodyJSON<Operation> };

export interface ClientOptions {
    baseUrl: string;
}

function buildURL(
    base: string,
    path: string,
    params?: {
        query?: Record<string, unknown>;
        path?: Record<string, unknown>;
    },
): URL {
    const result = new URL(path, base);

    if (params?.query !== undefined) {
        const searchParams = new URLSearchParams(
            Object.entries(params.query).map(([key, val]): [string, string] => [
                key,
                String(val),
            ]),
        );
        result.search = searchParams.toString();
    }

    if (params?.path !== undefined) {
        const path = params.path;
        result.pathname = result.pathname.replace(
            /\{(.*)\}/,
            (match, key): string =>
                path[key] !== undefined ? String(path[key]) : match,
        );
    }

    return result;
}

function handleResponse<Operation extends object>(
    resp: Response,
): RequestResult<Operation> {
    const hasBody =
        (typeof resp.body === 'string' || resp.body instanceof String) &&
        resp.body.length > 0;
    if (resp.status >= 200 && resp.status < 300) {
        return {
            data: hasBody ? JSON.parse(resp.body as string) : undefined,
            error: undefined,
            response: resp,
        };
    }
    return {
        data: undefined,
        error: hasBody ? JSON.parse(resp.body as string) : undefined,
        response: resp,
    };
}

export function createClient<Paths extends Record<any, any>>(
    opts: ClientOptions,
): Client<Paths> {
    const cleanBaseUrl = opts.baseUrl.replace(/\/+$/, '');

    return {
        get: (url, ropts) => {
            const rurl = buildURL(cleanBaseUrl, String(url), ropts.params);
            const resp = http.get(rurl.toString());

            return handleResponse(resp);
        },
        post: (url, ropts) => {
            const rurl = buildURL(cleanBaseUrl, String(url), ropts.params);
            const body =
                ropts.body !== undefined
                    ? JSON.stringify(ropts.body)
                    : undefined;
            const resp = http.post(rurl.toString(), body, {
                headers: CONTENT_TYPE_HEADER,
            });

            return handleResponse(resp);
        },
    };
}
