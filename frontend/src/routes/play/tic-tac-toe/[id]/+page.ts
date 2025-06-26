import type { PageLoad } from './$types';

export const load: PageLoad = async ({ params, fetch }) => {
	const roomId = params.id;
	const res = await fetch(import.meta.env.VITE_SOCKET_URL + `/api/v1/rooms/${roomId}`);
	const room = await res.json();

	return { room };
};
