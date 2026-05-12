import { templates } from '$lib/data';
import type { Actions } from './$types';

export const load = async () => {
	return { data: templates };
};

export const actions = {
	delete: async ({ request }) => {
		const data = await request.formData();
		const id = data.get('id');
		// TODO: Подрубить бэк
		console.log(`Deleting contact with id: ${id}`);
		return { success: true };
	},
	send: async ({ request }) => {
		const data = await request.formData();
		const id = data.get('id');
		// TODO: Подрубить бэк
		console.log(`Sending message to: ${id}`);
		return { success: true };
	}
} satisfies Actions;
