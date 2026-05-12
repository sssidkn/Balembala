import type { Actions } from './$types';
import { bound_contacts, templates } from '$lib/data';

export const load = async ({ params }) => {
	return {
		template: templates[Number(params.id)],
		contacts: bound_contacts
	};
};

export const actions = {
	default: async ({ request, params }) => {
		const id = params.id;
		const data = await request.formData();
		console.log(`editing template with ${id}`);
		const template = {
			id,
			title: data.get('title'),
			message: data.get('message')
		};
		// TODO: Подрубить бэк
		console.log(template);
		return { success: true };
	}
} satisfies Actions;
