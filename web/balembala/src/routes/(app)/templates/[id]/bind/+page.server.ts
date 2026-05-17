import { error, redirect } from '@sveltejs/kit';
import type { Actions } from './$types';
import { resolve } from '$app/paths';
import { API_URL } from '$env/static/private';
import type { Contact } from '$lib/data';

export const load = async ({ fetch, params }) => {
	const responses = [
		await fetch(API_URL + '/contacts', {
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
	const contacts = new Set(bodies[1].contacts.map((c: Contact) => c.id));
	const diff: Set<Contact> = new Set();
	for (const contact of bodies[0].contacts as Contact[]) {
		if (!contacts.has(contact.id)) {
			diff.add(contact);
		}
	}
	return {
		data: diff
	};
};

export const actions = {
	default: async ({ request, params, fetch }) => {
		const id = params.id;
		const data = await request.formData();
		const contacts_id = data.get('contacts_id');
		const response = await fetch(API_URL + '/template/' + params.id, {
			method: 'POST',
			body: contacts_id
		});
		if (!response.ok) {
			error(response.status, response.statusText);
		}
		redirect(308, resolve(`/templates/${id}`));
	}
} satisfies Actions;
