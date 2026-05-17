import { API_URL } from '$env/static/private';
import type { Actions } from './$types';
import type { Contact, Template } from '$lib/data.js';
import { error } from '@sveltejs/kit';

export const load = async ({ params, fetch }) => {
	const responses = [
		await fetch(API_URL + '/template/' + params.id, {
			method: 'GET'
		}),
		await fetch(API_URL + '/contacts/' + params.id, {
			method: 'GET'
		})
	];
	for (const response of responses) {
		if (!response.ok) {
			error(response.status, response.statusText);
		}
	}
	const bodies = [await responses[0].json(), await responses[1].json()];
	return {
		template: bodies[0].template as Template,
		contacts: bodies[1].contacts as Contact[]
	};
};

export const actions = {
	default: async ({ request, params, fetch }) => {
		const data = await request.formData();
		const template = {
			title: data.get('title'),
			message: data.get('message')
		};
		const response = await fetch(API_URL + '/template/' + params.id, {
			method: 'PUT',
			body: JSON.stringify(template)
		});
		if (!response.ok) {
			error(response.status, response.statusText);
		}
		return { success: true };
	}
} satisfies Actions;
