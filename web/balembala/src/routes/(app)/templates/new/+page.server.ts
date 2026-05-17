import { error, redirect } from '@sveltejs/kit';
import type { Actions } from './$types';
import { resolve } from '$app/paths';
import { API_URL } from '$env/static/private';

export const load = async () => {
	return {
		id: '',
		title: 'New template',
		message: 'Sample message'
	};
};

export const actions = {
	default: async ({ request, fetch }) => {
		const data = await request.formData();
		const template = {
			title: data.get('title'),
			message: data.get('message')
		};
		const response = await fetch(API_URL + '/templates', {
			method: 'POST',
			body: JSON.stringify(template)
		});
		if (!response.ok) {
			error(response.status, response.statusText);
		}
		redirect(303, resolve('/templates'));
	}
} satisfies Actions;
