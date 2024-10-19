import type { Client, FetchOptions } from 'openapi-fetch';
import type { PathsWithMethod, FilterKeys } from 'openapi-typescript-helpers';
import type { paths } from '$lib/sdk/schema';

function extractNextCursor(linkHeader: string | null): string | null {
    if (!linkHeader) {
        return null;
    }

    const regex = /<(.*)>; rel="next"/g;
    let result: string | null = null;
    for (const link of linkHeader.split(',')) {
        const nextURI = link.match(regex)?.[1];
        if (nextURI) {
            result = new URL(nextURI).searchParams.get('after');
            if (result) {
                break;
            }
        }
    }
    return result;
}

export default async function* <
    Path extends PathsWithMethod<paths, 'get'>,
    Options extends FetchOptions<FilterKeys<paths[Path], 'get'>>
>(client: Client<paths, `${string}/${string}`>, path: Path, options: Options) {
    let { data, error, response } = await client.GET(path, options);

    const clonedOpts = { params: {}, ...options };
    clonedOpts.params.query ??= {};
    let nextCursor = extractNextCursor(response.headers.get('Link'));
    yield { data, error, response, hasNext: !!nextCursor };

    while (nextCursor) {
        clonedOpts.params.query.after = nextCursor;
        ({ data, error, response } = await client.GET(path, clonedOpts));
        nextCursor = extractNextCursor(response.headers.get('Link'));
        yield { data, error, response, hasNext: !!nextCursor };
    }
}
