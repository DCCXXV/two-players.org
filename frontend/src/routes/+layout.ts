import type { LayoutLoad } from './$types';

export const load: LayoutLoad = async ({ url }) => {
	const title = 'Two Players | Online Combinatorial Games';
	const description =
		'Challenge your friends to classic and new combinatorial games like Tic-Tac-Toe, Domineering, and more. No registration required, just create a room and share the link!';
	const imageUrl = `${url.origin}/screenshots/tic-tac-toe.png`;

	return {
		meta: {
			title,
			description,
			imageUrl,
			url: url.origin
		}
	};
};
