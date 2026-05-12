import type { Actions } from './$types';
import { load as dataload } from '../+page.server';
import type { Contact } from '$lib/data';

function contact(id: string, name: string, email: string): Contact {
	return { id, name, email };
}

export const load = async ({ params }) => {
	return {
		template: (await dataload()).data[Number(params.edit)],
		contacts: [contact('0', 'First', 'first@mail.ru'), contact('1', 'Second', 'second@mail.ru')]
	};
};

export const actions = {
	default: async ({ request, params }) => {
		const id = params.edit;
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
