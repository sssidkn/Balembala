import type { Actions } from './$types';
import { API_URL } from '$env/static/private';
import { error } from '@sveltejs/kit';
import type { Contact } from '$lib/data.js';

export const load = async ({ params, fetch }) => {
	const response = await fetch(API_URL + '/contact/' + params.id, {
		method: 'GET'
	});
	if (!response.ok) {
		error(response.status, response.statusText);
	}
	const body = await response.json();
	return { data: body.contact as Contact };
};

export const actions = {
	default: async ({ request, params, fetch }) => {
		const data = await request.formData();
		const contact = {
			name: data.get('name'),
			email: data.get('email')
		};
		const response = await fetch(API_URL + '/contact/' + params.id, {
			method: 'PUT',
			body: JSON.stringify(contact)
		});
		if (!response.ok) {
			error(response.status, response.statusText);
		}
		return { success: true };
	}
} satisfies Actions;
