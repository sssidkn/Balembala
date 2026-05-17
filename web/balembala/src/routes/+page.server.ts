import { redirect } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';
import { deleteAuth } from '$lib/auth';

export const load: PageServerLoad = async () => {
	throw redirect(308, '/templates');
};

export const actions = {
	logout: async ({ cookies }) => {
		deleteAuth(cookies);
		redirect(303, '/login');
	}
};
