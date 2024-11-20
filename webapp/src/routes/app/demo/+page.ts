// import { newClient } from '$lib/utils/client'; import { redirect, error } from '@sveltejs/kit'; import { modelToError } from '$lib/utils/error-adapters'; import type { PageLoad } from './$types'; import { Button } from 'carbon-components-svelte';

//             export const load: PageLoad = async ({ fetch }) => {
//                     const client = newClient({ fetch });
//                         const { data, error: err, response } = await client.GET('/demo');

//                                                                                     if (err) {
//                  switch (err.status || response.status) {
//                         case 401:
//                 redirect(307, '/auth/login');
//                             break;
//                     default:
//                             error(err.status ?? 500, modelToError(err));
//                     }
//         }
//             console.log(data);

// return {message: data};
//                         };
