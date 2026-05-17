import { fail, redirect } from '@sveltejs/kit';
import type { Actions } from './$types';
import { setAuth } from '$lib/auth';
import { API_URL } from '$env/static/private';
export const actions = {
	default: async ({ request, cookies, fetch }) => {
		const data = await request.formData();
		const email = data.get('email');
		const password = data.get('password');
		const response = await fetch(API_URL + '/login', {
			method: 'POST',
			body: JSON.stringify({ email, password })
		});
		if (!response.ok) {
			return fail(response.status, { email, error: response.statusText });
		}
		const body = await response.json();
		setAuth(cookies, body.token);
		redirect(303, '/');
	}
} satisfies Actions;
