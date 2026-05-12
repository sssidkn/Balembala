import type { Actions } from './$types';
import { load as dataload } from '../+page.server';

export const load = async ({ params }) => {
	// TODO: back
	return (await dataload()).data[Number(params.id)];
};

export const actions = {
	default: async ({ request, params }) => {
		const data = await request.formData();
		const id = params.id;
		console.log('Editing contact');
		const contact = {
			id,
			name: data.get('name'),
			email: data.get('email')
		};
		// TODO: Подрубить бэк
		console.log(contact);
		return { success: true };
	}
} satisfies Actions;
