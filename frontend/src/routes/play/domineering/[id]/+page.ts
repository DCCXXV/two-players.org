import type { PageLoad } from './$types';

export const load: PageLoad = async ({ params, fetch, url }) => {
	try {
		const roomId = params.id;
		const storedRoom = sessionStorage.getItem(`room_${roomId}`);

		if (storedRoom) {
			console.log('Room data found in sessionStorage');
			const room = JSON.parse(storedRoom);
			return {
				room,
				error: null,
				meta: {
					title: `Domineering Game | Room ${roomId}`,
					description: 'Join this Domineering game and challenge your friend!',
					imageUrl: `${url.origin}/screenshots/tic-tac-toe.png`,
					url: url.href
				}
			};
		}

		const res = await fetch(import.meta.env.VITE_SOCKET_URL + `/api/v1/rooms/${roomId}`);

		if (!res.ok) {
			const errorData = await res
				.json()
				.catch(() => ({ message: `HTTP error! status: ${res.status}` }));
			return { room: null, error: errorData.error || errorData.message };
		}

		const room = await res.json();
		return {
			room,
			error: null,
			meta: {
				title: `Domineering Game | Room ${roomId}`,
				description: 'Join this Domineering game and challenge your friend!',
				imageUrl: `${url.origin}/screenshots/tic-tac-toe.png`,
				url: url.href
			}
		};
	} catch (e) {
		return { room: null, error: e instanceof Error ? e.message : 'An unknown error occurred' };
	}
};
