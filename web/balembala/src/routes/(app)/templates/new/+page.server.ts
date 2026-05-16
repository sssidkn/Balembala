import { redirect } from '@sveltejs/kit';
import type { Actions } from './$types';
import { resolve } from '$app/paths';

export const load = async () => {
	return {
		id: '',
		title: 'New template',
		message: 'Sample message'
	};
};

export const actions = {
	default: async ({ request }) => {
		const data = await request.formData();
		console.log('creating new template');
		const template = {
			title: data.get('title'),
			message: data.get('message')
		};
		// TODO: Подрубить бэк
		console.log(template);
		redirect(303, resolve('/templates'));
	}
} satisfies Actions;
