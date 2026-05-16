import { redirect } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';
import { ACCESS_TOKEN } from '$lib/constants.server';

export const load: PageServerLoad = async () => {
	throw redirect(308, '/templates');
};

// Когда-нибудь я вырасту, и подключу бэк
export const actions = {
	logout: async ({ cookies }) => {
		cookies.delete(ACCESS_TOKEN, { path: '/' });
		redirect(303, '/login');
	}
};
