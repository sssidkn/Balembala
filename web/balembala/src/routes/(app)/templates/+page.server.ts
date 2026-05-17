import { API_URL } from '$env/static/private';
import type { Template } from '$lib/data';
import { error } from '@sveltejs/kit';
import type { Actions } from './$types';

export const load = async ({ fetch }) => {
	const response = await fetch(API_URL + '/templates', {
		method: 'GET'
	});
	if (!response.ok) {
		error(response.status, response.statusText);
	}
	const body = await response.json();
	return { data: body.templates as Template[] };
};

export const actions = {
	delete: async ({ request, fetch }) => {
		const data = await request.formData();
		const id = data.get('id');
		const response = await fetch(API_URL + '/template/' + id, {
			method: 'DELETE'
		});
		if (!response.ok) {
			error(response.status, response.statusText);
		}
		return { success: true };
	},
	send: async ({ request, fetch }) => {
		const data = await request.formData();
		const id = data.get('id');
		const response = await fetch(API_URL + '/send/' + id, {
			method: 'POST'
		});
		if (!response.ok) {
			error(response.status, response.statusText);
		}
		return { success: true };
	}
} satisfies Actions;
