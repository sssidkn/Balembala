import { redirect } from '@sveltejs/kit';
import type { Actions } from './$types';
import { resolve } from '$app/paths';

export const load = async () => {
	return {
		id: '',
		name: 'New contact',
		email: ''
	};
};

export const actions = {
	default: async ({ request }) => {
		const data = await request.formData();
		console.log('creating new contact');
		const contact = {
			name: data.get('name'),
			email: data.get('email')
		};
		// TODO: Подрубить бэк
		console.log(contact);
		redirect(303, resolve('/contacts'));
	}
} satisfies Actions;
