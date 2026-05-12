import type { Actions } from './$types';
import { load as dataload } from '../+page.server';

export const load = async ({ params }) => {
	if (params.edit == 'new') {
		return {
			id: '',
			name: 'New contact',
			email: ''
		};
	} else {
		// TODO: back
		return (await dataload()).data[Number(params.edit)];
	}
};

export const actions = {
	default: async ({ request, params }) => {
		const data = await request.formData();
		const id = params.edit;
		if (id == 'new') {
			console.log('creating new contact');
		} else {
			console.log('Editing contact');
		}
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
