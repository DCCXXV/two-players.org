<div class="mx-auto max-w-4xl">
	<div class="mb-6 border-b-1 border-zinc-700 pb-2">
		<h1 class="lora-700 text-3xl text-zinc-200">Privacy Policy</h1>
		<p class="mt-2 text-sm text-zinc-500">Last updated: November 2025</p>
	</div>

	<div class="space-y-6">
		<section>
			<h2 class="mb-2 text-xl font-semibold text-zinc-300">What Data We Collect</h2>
			<p class="mb-3 text-zinc-400">
				We collect minimal data to facilitate gameplay:
			</p>
			<ul class="list-disc space-y-2 pl-6 text-zinc-400">
				<li>
					<strong class="text-zinc-300">Display names:</strong> Automatically generated names (e.g., "Alice52891")
					assigned when you connect. You can change your display name, but it must be unique among active players.
				</li>
				<li>
					<strong class="text-zinc-300">Session data:</strong> Your connection status, current room, and last activity
					timestamp stored temporarily in our PostgreSQL database.
				</li>
				<li>
					<strong class="text-zinc-300">Game room data:</strong> Room names you create, game type, and privacy settings.
				</li>
				<li>
					<strong class="text-zinc-300">Browser storage:</strong> Room information temporarily stored in your browser's
					session storage (deleted when you close the tab), and sound preference in local storage.
				</li>
			</ul>
		</section>

		<section>
			<h2 class="mb-2 text-xl font-semibold text-zinc-300">Analytics</h2>
			<p class="mb-3 text-zinc-400">
				We use Cloudflare Web Analytics to understand how visitors use our site. Cloudflare Analytics is
				privacy-first and does not use cookies or track individual users. It collects:
			</p>
			<ul class="list-disc space-y-1 pl-6 text-zinc-400">
				<li>Aggregated page views and visitor counts</li>
				<li>Referrer information (where visitors came from)</li>
				<li>Browser and device type (anonymized)</li>
				<li>Country-level location data (derived from IP, not stored)</li>
			</ul>
			<p class="mt-3 text-zinc-400">
				Cloudflare does not track you across sites or sell your data. Learn more at
				<a
					href="https://www.cloudflare.com/web-analytics/"
					target="_blank"
					rel="noopener noreferrer"
					class="text-blue-400 hover:underline"
				>
					Cloudflare Web Analytics
				</a>.
			</p>
		</section>

		<section>
			<h2 class="mb-2 text-xl font-semibold text-zinc-300">What We Don't Collect</h2>
			<p class="mb-3 text-zinc-400">
				We do not collect or store:
			</p>
			<ul class="list-disc space-y-1 pl-6 text-zinc-400">
				<li>Email addresses, phone numbers, or any personal contact information</li>
				<li>Personal IP addresses or device identifiers</li>
				<li>Payment or financial information</li>
				<li>Browsing history or tracking cookies</li>
				<li>Cross-site tracking data or advertising profiles</li>
			</ul>
		</section>

		<section>
			<h2 class="mb-2 text-xl font-semibold text-zinc-300">How We Use Your Data</h2>
			<p class="text-zinc-400">
				Display names and session data are used exclusively to:
			</p>
			<ul class="list-disc space-y-1 pl-6 text-zinc-400">
				<li>Identify players within game rooms</li>
				<li>Manage WebSocket connections for real-time gameplay</li>
				<li>Display who is in each game room</li>
			</ul>
			<p class="mt-3 text-zinc-400">
				We do not use your data for advertising, tracking, profiling, or any purpose beyond facilitating gameplay.
			</p>
		</section>

		<section>
			<h2 class="mb-2 text-xl font-semibold text-zinc-300">Chat Messages</h2>
			<p class="text-zinc-400">
				Chat messages are stored temporarily in memory only during your game session. They are automatically
				deleted when you leave the room or the room closes. Messages are automatically filtered for profanity
				and rate-limited to 5 messages per 10 seconds to prevent spam.
			</p>
		</section>

		<section>
			<h2 class="mb-2 text-xl font-semibold text-zinc-300">Data Retention and Deletion</h2>
			<p class="mb-3 text-zinc-400">
				All data is ephemeral and automatically deleted:
			</p>
			<ul class="list-disc space-y-1 pl-6 text-zinc-400">
				<li>Display names are deleted immediately when you disconnect</li>
				<li>Game rooms are deleted when the host leaves or all players disconnect</li>
				<li>Chat messages are cleared from memory when the room closes</li>
				<li>Stale connections are automatically cleaned up after 5 minutes of inactivity</li>
			</ul>
			<p class="mt-3 text-zinc-400">
				Browser session storage is automatically cleared when you close the browser tab.
			</p>
		</section>

		<section>
			<h2 class="mb-2 text-xl font-semibold text-zinc-300">Third-Party Services</h2>
			<p class="mb-3 text-zinc-400">
				We use the following third-party services:
			</p>
			<ul class="list-disc space-y-1 pl-6 text-zinc-400">
				<li>
					<strong class="text-zinc-300">Google Fonts:</strong> For web font delivery via Google's CDN
				</li>
				<li>
					<strong class="text-zinc-300">Cloudflare Web Analytics:</strong> Privacy-first analytics without cookies
					or cross-site tracking
				</li>
			</ul>
			<p class="mt-3 text-zinc-400">
				We do not use advertising networks, tracking pixels, or invasive analytics.
			</p>
		</section>

		<section>
			<h2 class="mb-2 text-xl font-semibold text-zinc-300">Your Rights (GDPR)</h2>
			<p class="text-zinc-400">
				Under GDPR, you have rights regarding your personal data. Given the ephemeral nature of our data
				collection (no accounts, no persistent tracking), these rights are automatically satisfied by
				disconnecting from the service, which immediately deletes all your data.
			</p>
		</section>

		<section>
			<h2 class="mb-2 text-xl font-semibold text-zinc-300">Contact Us</h2>
			<p class="text-zinc-400">
				If you have any questions about this Privacy Policy, please contact us at
				<a href="mailto:twoplayers@tutamail.com" class="text-blue-400 hover:underline">
				    admin@two-players.org
				</a>.
			</p>
		</section>
	</div>
</div>
