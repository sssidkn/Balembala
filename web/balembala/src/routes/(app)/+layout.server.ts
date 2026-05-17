import { redirect } from '@sveltejs/kit';
import type { LayoutServerLoad } from './$types';
import { getAuth } from '$lib/auth';

export const load: LayoutServerLoad = async ({ cookies }) => {
	const auth = getAuth(cookies);
	if (!auth) {
		throw redirect(303, '/login');
	}
};
