import { ACCESS_TOKEN } from '$lib/constants.server';
import type { Handle } from '@sveltejs/kit';

export const handle: Handle = async ({ event, resolve }) => {
	const token = event.cookies.get(ACCESS_TOKEN);

	if (token) {
		event.locals.logged_in = true;
	} else {
		event.locals.logged_in = false;
	}
	return resolve(event);
};
