export interface Contact {
	id: string;
	name: string;
	email: string;
}

export interface Template {
	id: string;
	title: string;
	message: string;
}

function template(id: string, title: string, message: string): Template {
	return { id, title, message };
}

export const templates = [
	template('0', 'First', 'message1'),
	template('1', 'Second', 'message2'),
	template('2', 'Third', 'message3'),
	template('3', 'Forth', 'message4'),
	template('4', 'Fifth', 'message5')
];

function contact(id: string, name: string, email: string): Contact {
	return { id, name, email };
}

export const contacts = [
	contact('0', 'First', 'first@mail.ru'),
	contact('1', 'Second', 'second@mail.ru'),
	contact('2', 'Third', 'third@mail.ru'),
	contact('3', 'Forth', 'forth@mail.ru'),
	contact('4', 'Fifth', 'fifth@mail.ru')
];

export const bound_contacts = [contacts[0], contacts[1]];
