import createClient from "openapi-fetch";
import type { ClientOptions } from "openapi-fetch";
import type {paths} from "$lib/sdk/schema";
import { BACKEND_SERVER } from '$lib/constants';

const defaultOpts: ClientOptions = {
    baseUrl: BACKEND_SERVER,
    credentials: "include",
}

export function newClient(options?: ClientOptions) {
    const mergedOptions = {...defaultOpts, ...options};
    return createClient<paths>(mergedOptions);
}