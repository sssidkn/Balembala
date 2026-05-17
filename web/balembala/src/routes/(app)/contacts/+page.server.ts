import { API_URL } from '$env/static/private';
import { error } from '@sveltejs/kit';
import type { Actions } from './$types';
import type { Contact } from '$lib/data.js';

export const load = async ({ fetch }) => {
	const response = await fetch(API_URL + '/contacts', {
		method: 'GET'
	});
	if (!response.ok) {
		error(response.status, response.statusText);
	}
	const body = await response.json();
	return { data: body.contacts as Contact[] };
};

export const actions = {
	delete: async ({ request, fetch }) => {
		const data = await request.formData();
		const id = data.get('id');
		const response = await fetch(API_URL + '/contact/' + id, {
			method: 'DELETE'
		});
		if (!response.ok) {
			error(response.status, response.statusText);
		}
		return { success: true };
	}
} satisfies Actions;
