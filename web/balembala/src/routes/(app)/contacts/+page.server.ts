import type { Contact } from '$lib/data';
import type { Actions } from './$types';

function contact(id: string, name: string, email: string): Contact {
	return { id, name, email };
}

const data = [
	contact('0', 'First', 'first@mail.ru'),
	contact('1', 'Second', 'second@mail.ru'),
	contact('2', 'Third', 'third@mail.ru'),
	contact('3', 'Forth', 'forth@mail.ru'),
	contact('4', 'Fifth', 'fifth@mail.ru')
];

export const load = async () => {
	return { data };
};

export const actions = {
	delete: async ({ request }) => {
		const data = await request.formData();
		const id = data.get('id');
		// TODO: Подрубить бэк
		console.log(`Deleting contact with id: ${id}`);
		return { success: true };
	}
} satisfies Actions;
