import { getAuth } from '$lib/auth';
import type { HandleFetch } from '@sveltejs/kit';

export const handleFetch: HandleFetch = async ({ request, event, fetch }) => {
	const auth = getAuth(event.cookies);
	if (auth) {
		request.headers.set('Authorization', 'Bearer ' + auth);
	}
	return fetch(request);
};
