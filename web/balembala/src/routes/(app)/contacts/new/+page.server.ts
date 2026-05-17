import { error, redirect } from '@sveltejs/kit';
import type { Actions } from './$types';
import { resolve } from '$app/paths';
import { API_URL } from '$env/static/private';

export const load = async () => {
	return {
		id: '',
		name: 'New contact',
		email: ''
	};
};

export const actions = {
	default: async ({ request, fetch }) => {
		const data = await request.formData();
		const contact = {
			name: data.get('name'),
			email: data.get('email')
		};
		const response = await fetch(API_URL + '/contacts', {
			method: 'POST',
			body: JSON.stringify(contact)
		});
		if (!response.ok) {
			error(response.status, response.statusText);
		}
		redirect(303, resolve('/contacts'));
	}
} satisfies Actions;
