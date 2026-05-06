import { fail, redirect } from '@sveltejs/kit';
import type { Actions } from './$types';
import { ACCESS_TOKEN } from '$lib/constants.server';
export const actions = {
	default: async ({ request, cookies }) => {
		const data = await request.formData();
		const email = data.get('email');
		const password = data.get('password');
		if (password === 'error') {
			return fail(400, { email, error: 'Generic error' });
		}
		console.log({ email, password });
		// TODO: Подрубить бэк
		cookies.set(ACCESS_TOKEN, 'jwt goes here', { path: '/' });
		redirect(303, '/');
	}
} satisfies Actions;
