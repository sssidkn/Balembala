import { fail, redirect } from '@sveltejs/kit';
import type { Actions } from './$types';
import { ACCESS_TOKEN } from '$lib/constants.server';
export const actions = {
	default: async ({ request, cookies }) => {
		const data = await request.formData();
		const username = data.get('username');
		const email = data.get('email');
		const password = data.get('password');
		if (username === 'error') {
			return fail(400, { username, email, error: 'Generic error' });
		}
		console.log({ username, email, password });
		// TODO: Подрубить бэк
		cookies.set(ACCESS_TOKEN, 'jwt goes here', { path: '/' });
		redirect(303, '/');
	}
} satisfies Actions;
